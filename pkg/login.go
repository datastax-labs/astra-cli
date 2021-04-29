//  Copyright 2021 Ryan Svihla
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package pkg is the top level package for shared libraries
package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"net/http"
	"net"
	"encoding/json"
	"strings"

	"github.com/rsds143/astra-cli/pkg/env"
	astraops "github.com/datastax/astra-client-go/v2/astra"
)


func closeBody(res *http.Response) {
	if err := res.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to close request body '%v'", err)
	}
}

// LoginService provides interface to implement logins and produce an Client
type LoginService interface {
	Login() (Client, error)
}

// Client is the abstraction for client interactions. Allows alternative db management clients
type Client interface {
	CreateDb(astraops.DatabaseInfoCreate) (astraops.Database, error)
	Terminate(string, bool) error
	FindDb(string) (astraops.Database, error)
	ListDb(string, string, string, int) ([]astraops.Database, error)
	Park(string) error
	Unpark(string) error
	Resize(string, int) error
	GetSecureBundle(string) (astraops.CredsURL, error)
	GetTierInfo() ([]astraops.AvailableRegionCombination, error)
}

// AuthenticatedClient is the abstraction for client interactions. Allows alternative db management clients
type AuthenticatedClient struct {
	token string
	client astraops.ClientInterface
	verbose bool
}

func (a *AuthenticatedClient) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", a.token)
	req.Header.Set("Content-Type", "application/json")
}


func (a *AuthenticatedClient) requestEditorForAPI(rctx context.Context, req *http.Request) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", a.token)
	req.Header.Set("Content-Type", "application/json")
	return nil
}

// WaitUntil will keep checking the database for the requested status until it is available. Eventually it will timeout if the operation is not
// yet complete.
// * @param id string - the database id to find
// * @param tries int - number of attempts
// * @param intervalSeconds int - seconds to wait between tries
// * @param status StatusEnum - status to wait for
// @returns (Database, error)
func (a *AuthenticatedClient) WaitUntil(id string, tries int, intervalSeconds int, status astraops.StatusEnum) (astraops.Database, error) {
	for i := 0; i < tries; i++ {
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
		db, err := a.FindDb(id)
		if err != nil {
			if a.verbose {
				log.Printf("db %s not able to be found with error '%v' trying again %v more times", id, err, tries-i-1)
			} else {
				log.Printf("waiting")
			}
			continue
		}
		if db.Status == status {
			return db, nil
		}
		if a.verbose {
			log.Printf("db %s in state %v but expected %v trying again %v more times", id, db.Status, status, tries-i-1)
		} else {
			log.Printf("waiting")
		}
	}
	return astraops.Database{}, fmt.Errorf("unable to find db id %s with status %s after %v seconds", id, status, intervalSeconds*tries)
}

func (a* AuthenticatedClient) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10 * time.Second)
}

// CreateDb creates a database in Astra waits until it is in a created state
// * @param createDb Definition of new database
// @return (Database, error)
func (a* AuthenticatedClient) CreateDb(db astraops.DatabaseInfoCreate) (astraops.Database, error) {
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.CreateDatabase(ctx, astraops.CreateDatabaseJSONRequestBody(db), a.requestEditorForAPI)
	if err != nil {
		return astraops.Database{}, fmt.Errorf("failed creating database with: %w", err)
	}
	defer closeBody(res)
	if res.StatusCode != 201 {
		return astraops.Database{}, readErrorFromResponse(res, 201)
	}
    id := res.Header.Get("location")
	newDb, err := a.WaitUntil(id, 30, 30, astraops.StatusEnum_ACTIVE)
	if err != nil {
		return newDb, fmt.Errorf("create db failed because '%v'", err)
	}
    return newDb, nil
}

// Terminate deletes the database at the specified id and will block until it shows up as deleted or is removed from the system
// * @param databaseID string representation of the database ID
// * @param "PreparedStateOnly" -  For internal use only.  Used to safely terminate prepared databases
// @return error
func (a* AuthenticatedClient) Terminate(databaseID string, internalOnly bool) error {
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.TerminateDatabase(ctx, astraops.DatabaseIdParam(databaseID), nil, a.requestEditorForAPI)
	if err != nil {
		return fmt.Errorf("failed to terminate database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != 202 {
		return readErrorFromResponse(res, 202)
	}
	tries := 30
	intervalSeconds := 10
	var lastResponse string
	var lastStatusCode int
	for i := 0; i < tries; i++ {
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
		serviceURL := "https://api.astra.datastax.com/v2/databases"
		req, err := http.Get(fmt.Sprintf("%s/%s", serviceURL, databaseID))
		if err != nil {
			return fmt.Errorf("failed creating request to find db with id %s with: %w", databaseID, err)
		}
		defer closeBody(req)
		lastStatusCode = res.StatusCode
		if res.StatusCode == 401 {
			return nil
		}
		if res.StatusCode == 200 {
			var db astraops.Database
			err = json.NewDecoder(res.Body).Decode(&db)
			if err != nil {
				return fmt.Errorf("critical error trying to get status of database not deleted, unable to decode response with error: %v", err)
			}
			if db.Status == astraops.StatusEnum_TERMINATED || db.Status == astraops.StatusEnum_TERMINATING {
				if a.verbose {
					log.Printf("delete status is %v for db %v and is therefore successful, we are going to exit now", db.Status, databaseID)
				}
				return nil
			}
			if a.verbose {
				log.Printf("db %s not deleted yet expected status code 401 or a 200 with a db Status of %v or %v but was 200 with a db status of %v. trying again", databaseID, astraops.StatusEnum_TERMINATED, astraops.StatusEnum_TERMINATING, db.Status)
			} else {
				log.Printf("waiting")
			}
			continue
		}
		lastResponse = fmt.Sprintf("%v", readErrorFromResponse(res, 200, 401))
		if a.verbose {
			log.Printf("db %s not deleted yet expected status code 401 or a 200 with a db Status of %v or %v but was: %v and error was '%v'. trying again", databaseID, astraops.StatusEnum_TERMINATED,  astraops.StatusEnum_TERMINATING, res.StatusCode, lastResponse)
		} else {
			log.Printf("waiting")
		}
	}
	return fmt.Errorf("delete of db %s not complete. Last response from finding db was '%v' and last status code was %v", databaseID, lastResponse, lastStatusCode)
}

// FindDb Returns specified database
// * @param databaseID string representation of the database ID
// @return (Database, error)
func (a* AuthenticatedClient) FindDb(databaseID string) (astraops.Database, error) {
	var dbs astraops.Database
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.GetDatabase(ctx, astraops.DatabaseIdParam(databaseID), a.requestEditorForAPI)
	if err != nil {
		return dbs, fmt.Errorf("failed get database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != 200 {
		return dbs, readErrorFromResponse(res, 200)
	}
	err = json.NewDecoder(res.Body).Decode(&dbs)
	if err != nil {
		return astraops.Database{}, fmt.Errorf("unable to decode response with error: %w", err)
	}
	return dbs, nil 
}

// ListDb find all databases that match the parameters
// * @param "include" (optional.string) -  Allows filtering so that databases in listed states are returned
// * @param "provider" (optional.string) -  Allows filtering so that databases from a given provider are returned
// * @param "startingAfter" (optional.string) -  Optional parameter for pagination purposes. Used as this value for starting retrieving a specific page of results
// * @param "limit" (optional.int32) -  Optional parameter for pagination purposes. Specify the number of items for one page of data
// @return ([]Database, error)
func (a* AuthenticatedClient) 	ListDb(include string, provider string, startingAfter string, limit int) ([]astraops.Database, error) {
	var dbs []astraops.Database
	params := astraops.ListDatabasesParams{}
	if len(include) > 0 {
		params.Include = &include
	}
	if len(provider) > 0 {
		params.Provider = &provider
	}
	if len(startingAfter) > 0 {
		params.StartingAfter = &startingAfter
	}
	if limit > 0 {
		params.Limit = &limit
	}
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.ListDatabases(ctx, &params, a.requestEditorForAPI)
	if err != nil {
		return dbs, fmt.Errorf("failed listing databases with: %v", err)
	}
	defer closeBody(res)
	if res.StatusCode != 200 {
		return dbs, readErrorFromResponse(res, 200)
	}
	err = json.NewDecoder(res.Body).Decode(&dbs)
	if err != nil {
		return []astraops.Database{}, fmt.Errorf("unable to decode response with error: %v", err)
	}
	return dbs, nil
}

// Park parks the database at the specified id and will block until the database is parked
// * @param databaseID string representation of the database ID
// @return error
func (a* AuthenticatedClient) Park(databaseID string) error {
	ctx, cancel := a.getContext()
	res, err := a.client.ParkDatabase(ctx, astraops.DatabaseIdParam( databaseID), a.requestEditorForAPI)
	defer cancel()
	if err != nil {
		return fmt.Errorf("failed to park database id %s with: %w",  databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != 202 {
		return readErrorFromResponse(res, 202)
	}
	_, err = a.WaitUntil(databaseID, 30, 30, astraops.StatusEnum_PARKED)
	if err != nil {
		return fmt.Errorf("unable to check status for park db because of error '%v'", err)
	}
	return nil
}

// Unpark unparks the database at the specified id and will block until the database is unparked
// * @param databaseID String representation of the database ID
// @return error
func (a* AuthenticatedClient) Unpark(databaseID string) error {
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.UnparkDatabase(ctx, astraops.DatabaseIdParam (databaseID), a.requestEditorForAPI)
	if err != nil {
		return fmt.Errorf("failed to unpark database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != 202 {
		return readErrorFromResponse(res, 202)
	}
	_, err = a.WaitUntil(databaseID, 60, 30, astraops.StatusEnum_ACTIVE)
	if err != nil {
		return fmt.Errorf("unable to check status for unpark db because of error '%v'", err)
	}
	return nil
}

// Resize a database. Total number of capacity units desired should be specified. Reducing a size of a database is not supported at this time. Note you cannot resize a serverless database
// * @param databaseID string representation of the database ID
// * @param capacityUnits int32 containing capacityUnits key with a value greater than the current number of capacity units (max increment of 3 additional capacity units)
// @return error
func (a* AuthenticatedClient) 	Resize(databaseID string, capacityUnits int) error {
	ctx, cancel := a.getContext()
	defer cancel()
	units := astraops.ResizeDatabaseJSONRequestBody{
		CapacityUnits: &capacityUnits,
	}
	res, err := a.client.ResizeDatabase(ctx, astraops.DatabaseIdParam(databaseID), units, a.requestEditorForAPI)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		var resObj ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&resObj)
		if err != nil {
			return fmt.Errorf("unable to decode error response with error: %w", err)
		}
		return fmt.Errorf("expected status code 2xx but had: %v with error(s) - %v", res.StatusCode, FormatErrors(resObj.Errors))
	}
	return nil
}

// GetSecureBundle Returns a temporary URL to download a zip file with certificates for connecting to the database.
// The URL expires after five minutes.&lt;p&gt;There are two types of the secure bundle URL: &lt;ul&gt
// * @param databaseID string representation of the database ID
// @return (SecureBundle, error)
func (a* AuthenticatedClient) 	GetSecureBundle(databaseID string) (astraops.CredsURL, error) {
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.GenerateSecureBundleURL(ctx, astraops.DatabaseIdParam(databaseID), a.requestEditorForAPI)
	if err != nil {
		return astraops.CredsURL{}, fmt.Errorf("failed get secure bundle for database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != 200 {
		return astraops.CredsURL{}, readErrorFromResponse(res, 200)
	}
	var sb astraops.CredsURL
	err = json.NewDecoder(res.Body).Decode(&sb)
	if err != nil {
		return astraops.CredsURL{}, fmt.Errorf("unable to decode response with error: %w", err)
	}
	return sb, nil
}

// GetTierInfo Returns all supported tier, cloud, region, count, and capacitity combinations
// @return ([]TierInfo, error)
func (a* AuthenticatedClient) 	GetTierInfo() ([]astraops.AvailableRegionCombination, error) {
	ctx, cancel := a.getContext()
	defer cancel()
	res, err := a.client.ListAvailableRegions(ctx, a.requestEditorForAPI)
	if err != nil {
		return []astraops.AvailableRegionCombination{}, fmt.Errorf("failed listing tier info with: %w", err)
	}
	defer closeBody(res)
	if res.StatusCode != 200 {
		return []astraops.AvailableRegionCombination{}, readErrorFromResponse(res, 200)
	}
	var ti []astraops.AvailableRegionCombination
	err = json.NewDecoder(res.Body).Decode(&ti)
	if err != nil {
		return []astraops.AvailableRegionCombination{}, fmt.Errorf("unable to decode response with error: %w", err)
	}
	return ti, nil
}


// Creds knows how handle and store credentials
type Creds struct {
	GetHomeFunc func() (string, error) // optional. If not specified os.UserHomeDir is used for log base directory to find creds
}

// Login logs into the Astra DevOps API using the local configuration provided by the 'astra-cli login' command
func (c *Creds) Login() (Client, error) {
	getHome := c.GetHomeFunc
	if getHome == nil {
		getHome = os.UserHomeDir
	}
	confDir, confFile, err := GetHome(getHome)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to read conf dir with error '%v'", err)
	}
	hasToken, err := confFile.HasToken()
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to read token file '%v' with error '%v'", confFile.TokenPath, err)
	}
	var client *AuthenticatedClient
	if hasToken {
		token, err := ReadToken(confFile.TokenPath)
		if err != nil {
			return &AuthenticatedClient{}, fmt.Errorf("found token at '%v' but unable to read token with error '%v'", confFile.TokenPath, err)
		}
		return AuthenticateToken(token, env.Verbose)
	}
	hasSa, err := confFile.HasServiceAccount()
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to read service account file '%v' with error '%v'", confFile.SaPath, err)
	}
	if !hasSa {
		return &AuthenticatedClient{}, fmt.Errorf("unable to access any file for directory `%v`, run astra-cli login first", confDir)
	}
	clientInfo, err := ReadLogin(confFile.SaPath)
	if err != nil {
		return &AuthenticatedClient{}, err
	}
	client, err = Authenticate(clientInfo, env.Verbose)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("authenticate failed with error %v", err)
	}
	return client, nil
}


func readErrorFromResponse(res *http.Response, expectedCodes ...int) error {
	var resObj ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&resObj)
	if err != nil {
		return fmt.Errorf("unable to decode error response with error: '%v'. status code was %v", err, res.StatusCode)
	}
	var statusSuffix string
	if len(expectedCodes) > 0 {
		statusSuffix = "s"
	}
	var errorSuffix string
	if len(resObj.Errors) > 0 {
		errorSuffix = "s"
	}
	var codeString []string
	for _, c := range expectedCodes {
		codeString = append(codeString, fmt.Sprintf("%v", c))
	}
	formattedCodes := strings.Join(codeString, ", ")
	return fmt.Errorf("expected status code%v %v but had: %v error with error%v - %v", statusSuffix, formattedCodes, res.StatusCode, errorSuffix, FormatErrors(resObj.Errors))
}

// FormatErrors puts the API errors into a well formatted text output
func FormatErrors(es []Error) string {
	var formatted []string
	for _, e := range es {
		formatted = append(formatted, fmt.Sprintf("ID: %v Text: '%v'", e.ID, e.Message))
	}
	return strings.Join(formatted, ", ")
}

// ErrorResponse when the API has an error
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error when the api has an error this is the structure
type Error struct {
	// API specific error code
	ID int32 `json:"ID,omitempty"`
	// User-friendly description of error
	Message string `json:"message"`
}

// AuthenticateToken returns a client
// * @param token string - token generated for login in the astra UI
// * @param verbose bool - if true the logging is much more verbose
// @returns (*AuthenticatedClient , error)
func AuthenticateToken(token string, verbose bool) (*AuthenticatedClient, error) {
	astraClient, err := astraops.NewClient(astraops.ServerURL)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to create astra client %v", err)
	}
	return &AuthenticatedClient{
		client:  astraClient,
		token:   fmt.Sprintf("Bearer %s", token),
		verbose: verbose,
	}, nil
}

// Authenticate returns a client using legacy Service Account. This is not deprecated but one should move to AuthenticateToken
// * @param clientInfo - classic service account from legacy Astra
// * @param verbose bool - if true the logging is much more verbose
// @returns (*AuthenticatedClient , error)
func Authenticate(clientInfo ClientInfo, verbose bool) (*AuthenticatedClient, error) {
	url := "https://api.astra.datastax.com/v2/authenticateServiceAccount"
	body, err := json.Marshal(clientInfo)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to marshal JSON object with: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("failed creating request with: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c := newHTTPClient()
	res, err := c.Do(req)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("failed listing databases with: %w", err)
	}
	defer closeBody(res)
	if res.StatusCode != 200 {
		return &AuthenticatedClient{}, readErrorFromResponse(res, 200)
	}
	var tokenResponse TokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to decode response with error: %w", err)
	}
	if tokenResponse.Token == "" {
		return &AuthenticatedClient{}, errors.New("empty token in token response")
	}
	astraClient, err := astraops.NewClient(astraops.ServerURL)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unable to create astra client %v", err)
	}
	return &AuthenticatedClient{
		client:  astraClient,
		token:   fmt.Sprintf("Bearer %s", tokenResponse.Token),
		verbose: verbose,
	}, nil
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 10,
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}


// TokenResponse comes from the classic service account auth
type TokenResponse struct {
	Token  string  `json:"token"`
	Errors []Error `json:"errors"`
}