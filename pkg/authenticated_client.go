//  Copyright 2022 DataStax
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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/datastax-labs/astra-cli/pkg/env"
	"github.com/datastax/astra-client-go/v2/astra"
)

// Error when the api has an error this is the structure
type Error struct {
	// API specific error code
	ID int
	// User-friendly description of error
	Message string
}

// ErrorResponse when the API has an error
type ErrorResponse struct {
	Errors []Error
}

func closeBody(res *http.Response) {
	if err := res.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to close request body '%v'", err)
	}
}

func handleErrors(body []byte, statusCode string) error {
	var errorStrings []string

	var resObj ErrorResponse
	err := json.Unmarshal(body, &resObj)
	if err != nil {
		return fmt.Errorf("CRITIAL ERROR unable to decode error response with error: '%v'. body was '%v' and http status was %v", err, string(body), statusCode)
	}
	for _, e := range resObj.Errors {
		errorString := fmt.Sprintf("(%v:%v)", e.ID, e.Message)
		errorStrings = append(errorStrings, errorString)
	}
	return fmt.Errorf("%v with status code: %v", strings.Join(errorStrings, ", "), statusCode)
}

// FormatErrors puts the API errors into a well formatted text output
func FormatErrors(es []Error) string {
	var formatted []string
	for _, e := range es {
		formatted = append(formatted, fmt.Sprintf("ID: %v Text: '%v'", e.ID, e.Message))
	}
	return strings.Join(formatted, ", ")
}

// AuthenticatedClient has a token and the methods to query the Astra DevOps API
type AuthenticatedClient struct {
	token          string
	client         *http.Client
	astraclient    *astra.ClientWithResponses
	timeoutSeconds int
	verbose        bool
}

func newHTTPClient() *http.Client {
	expectTimeout := 1
	defaultTimeout := 10
	connections := 10
	return &http.Client{
		Timeout: time.Duration(defaultTimeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        connections,
			MaxConnsPerHost:     connections,
			MaxIdleConnsPerHost: connections,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(defaultTimeout) * time.Second,
				KeepAlive: time.Duration(defaultTimeout) * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   time.Duration(defaultTimeout) * time.Second,
			ResponseHeaderTimeout: time.Duration(defaultTimeout) * time.Second,
			ExpectContinueTimeout: time.Duration(expectTimeout) * time.Second,
		},
	}
}

func timeoutContext(timeSeconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Duration(timeSeconds)*time.Second),
	)
}

func AuthenticateToken(token string, verbose bool) (*AuthenticatedClient, error) {
	astraClient, err := astra.NewClientWithResponses(apiURL(), func(c *astra.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return nil
		})
		return nil
	})
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error setting up devops api client: %v", err)
	}
	timeout := 10
	authenticatedClient := &AuthenticatedClient{
		verbose:        verbose,
		timeoutSeconds: timeout,
		astraclient:    astraClient,
		client:         newHTTPClient(),
		token:          fmt.Sprintf("Bearer %s", token),
	}
	return authenticatedClient, nil
}

func Authenticate(clientInfo ClientInfo, verbose bool) (*AuthenticatedClient, error) {
	timeout := 10
	tokenInput := astra.AuthenticateServiceAccountTokenJSONRequestBody{
		ClientId:     clientInfo.ClientID,
		ClientName:   clientInfo.ClientName,
		ClientSecret: clientInfo.ClientSecret,
	}
	astraClientTmp, err := astra.NewClientWithResponses(apiURL())
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error setting up devops api client: %v", err)
	}
	ctx, cancel := timeoutContext(timeout)
	defer cancel()
	response, err := astraClientTmp.AuthenticateServiceAccountTokenWithResponse(ctx, tokenInput)
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error logging into devops api client: %v", err)
	}
	if response.StatusCode() != http.StatusOK {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error logging into devops api client: %v - %v", response.StatusCode(), response.Status())
	}
	token := response.JSON200.Token
	astraClient, err := astra.NewClientWithResponses(apiURL(), func(c *astra.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
			return nil
		})
		return nil
	})

	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error logging into devops api client: %v", err)
	}
	authenticatedClient := &AuthenticatedClient{
		token:          fmt.Sprintf("Bearer %s", *token),
		verbose:        verbose,
		timeoutSeconds: timeout,
		astraclient:    astraClient,
		client:         newHTTPClient(),
	}
	if err != nil {
		return &AuthenticatedClient{}, fmt.Errorf("unexpected error authenticating: %v", err)
	}

	return authenticatedClient, nil
}

func apiURL() string {
	if env.Verbose {
		log.Printf("env is %v", Env)
	}
	var url string
	switch Env {
	case "dev":
		url = "https://api.dev.cloud.datastax.com"
	case "test":
		url = "https://api.test.cloud.datastax.com"
	default:
		url = "https://api.astra.datastax.com"
	}
	if env.Verbose {
		log.Printf("api url is %v", url)
	}
	return url
}

func dbURL() string {
	url := fmt.Sprintf("%v/v2/databases", apiURL())
	if env.Verbose {
		log.Printf("db url is %v", url)
	}
	return url
}

var Env = "prod"

func (a *AuthenticatedClient) ctx() (context.Context, context.CancelFunc) {
	return timeoutContext(a.timeoutSeconds)
}

func (a *AuthenticatedClient) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", a.token)
	req.Header.Set("Content-Type", "application/json")
}

// WaitUntil will keep checking the database for the requested status until it is available. Eventually it will timeout if the operation is not
// yet complete.
// * @param id string - the database id to find
// * @param tries int - number of attempts
// * @param intervalSeconds int - seconds to wait between tries
// * @param status StatusEnum - status to wait for
// @returns (Database, error)
func (a *AuthenticatedClient) WaitUntil(id string, tries int, intervalSeconds int, status ...astra.StatusEnum) (astra.Database, error) {
	for i := 0; i < tries; i++ {
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
		db, err := a.FindDb(id)
		if err != nil {
			if a.verbose {
				log.Printf("db %s not able to be found with error '%v' trying again %v more times", id, err, tries-i-1)
			} else {
				fmt.Print(".")
			}
			continue
		}

		if db.Status == astra.StatusEnumERROR {
			return db, fmt.Errorf("database %v in error status, exiting", id)
		}
		var statusStrings []string
		for _, s := range status {
			if db.Status == s {
				return db, nil
			}
			statusStrings = append(statusStrings, fmt.Sprintf("%v", s))
		}
		if a.verbose {
			log.Printf("db %s in state(s) %v but expected %v trying again %v more times", id, db.Status, strings.Join(statusStrings, ", "), tries-i-1)
		} else {
			fmt.Print(".")
		}
	}
	return astra.Database{}, fmt.Errorf("unable to find db id %s with status %s after %v seconds", id, status, intervalSeconds*tries)
}

// ListDb find all databases that match the parameters
// * @param "include" (optional.string) -  Allows filtering so that databases in listed states are returned
// * @param "provider" (optional.string) -  Allows filtering so that databases from a given provider are returned
// * @param "startingAfter" (optional.string) -  Optional parameter for pagination purposes. Used as this value for starting retrieving a specific page of results
// * @param "limit" (optional.int) -  Optional parameter for pagination purposes. Specify the number of items for one page of data
// @return ([]Database, error)
func (a *AuthenticatedClient) ListDb(include string, provider string, startingAfter string, limit int) ([]astra.Database, error) {
	var params astra.ListDatabasesParams
	if len(include) > 0 {
		astraInclude := astra.ListDatabasesParamsInclude(include)
		params.Include = &astraInclude
	}
	if len(provider) > 0 {
		astraProvider := astra.ListDatabasesParamsProvider(provider)
		params.Provider = &astraProvider
	}
	if len(startingAfter) > 0 {
		params.StartingAfter = astra.StringPtr(startingAfter)
	}
	if limit > 0 {
		limitInt := limit
		params.Limit = &limitInt
	}
	ctx, cancel := a.ctx()
	defer cancel()
	dbs, err := a.astraclient.ListDatabasesWithResponse(ctx, &params)
	if err != nil {
		return []astra.Database{}, fmt.Errorf("unexpected error listing databases '%v'", err)
	}
	if dbs.StatusCode() != http.StatusOK {
		return []astra.Database{}, handleErrors(dbs.Body, dbs.Status())
	}

	return *dbs.JSON200, nil
}

// CreateDb creates a database in Astra, username and password fields are required only on legacy tiers and waits until it is in a created state
// * @param createDb Definition of new database
// @return (Database, error)
func (a *AuthenticatedClient) CreateDb(createDb astra.DatabaseInfoCreate) (astra.Database, error) {
	ctx, cancel := a.ctx()
	defer cancel()
	response, err := a.astraclient.CreateDatabaseWithResponse(ctx, astra.CreateDatabaseJSONRequestBody(createDb))
	if err != nil {
		return astra.Database{}, err
	}
	if response.StatusCode() != http.StatusCreated {
		return astra.Database{}, handleErrors(response.Body, response.Status())
	}
	id := response.HTTPResponse.Header.Get("location")

	tries := 120
	interval := 30
	db, err := a.WaitUntil(id, tries, interval, astra.StatusEnumACTIVE)
	if err != nil {
		return db, fmt.Errorf("waiting for status check on create db failed because '%v'", err)
	}
	return db, nil
}

// FindDb Returns specified database
// * @param databaseID string representation of the database ID
// @return (Database, error)
func (a *AuthenticatedClient) FindDb(databaseID string) (astra.Database, error) {
	ctx, cancel := a.ctx()
	defer cancel()
	dbs, err := a.astraclient.GetDatabaseWithResponse(ctx, astra.DatabaseIdParam(databaseID))
	if err != nil {
		return astra.Database{}, fmt.Errorf("failed creating request to find db with id %s with: %w", databaseID, err)
	}
	if dbs.StatusCode() != http.StatusOK {
		return astra.Database{}, handleErrors(dbs.Body, dbs.Status())
	}
	return *dbs.JSON200, nil
}

// AddKeyspaceToDb Adds keyspace into database
// * @param databaseID string representation of the database ID
// * @param keyspaceName Name of database keyspace
// @return error
func (a *AuthenticatedClient) AddKeyspaceToDb(databaseID string, keyspaceName string) error {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.AddKeyspaceWithResponse(ctx, astra.DatabaseIdParam(databaseID), astra.KeyspaceNameParam(keyspaceName))
	if err != nil {
		return fmt.Errorf("failed creating request to add keyspace to db with id %s with: %w", databaseID, err)
	}
	if res.StatusCode() != http.StatusOK {
		return handleErrors(res.Body, res.Status())
	}
	return nil
}

// GetSecureBundle Returns a temporary URL to download a zip file with certificates for connecting to the database.
// The URL expires after five minutes.&lt;p&gt;There are two types of the secure bundle URL: &lt;ul&gt
// * @param databaseID string representation of the database ID
// @return (SecureBundle, error)
func (a *AuthenticatedClient) GetSecureBundle(databaseID string) (astra.CredsURL, error) {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.GenerateSecureBundleURLWithResponse(ctx, astra.DatabaseIdParam(databaseID))

	if err != nil {
		return astra.CredsURL{}, fmt.Errorf("failed get secure bundle for database id %s with: %w", databaseID, err)
	}
	if res.StatusCode() != http.StatusOK {
		return astra.CredsURL{}, handleErrors(res.Body, res.Status())
	}
	return *res.JSON200, nil
}

// Terminate deletes the database at the specified id and will block until it shows up as deleted or is removed from the system
// * @param databaseID string representation of the database ID
// * @param "PreparedStateOnly" -  For internal use only.  Used to safely terminate prepared databases
// @return error
func (a *AuthenticatedClient) Terminate(id string, preparedStateOnly bool) error {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.TerminateDatabaseWithResponse(ctx, astra.DatabaseIdParam(id), &astra.TerminateDatabaseParams{
		PreparedStateOnly: &preparedStateOnly,
	})
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusAccepted {
		return handleErrors(res.Body, res.Status())
	}
	tries := 30
	interval := 10
	_, err = a.WaitUntil(id, tries, interval, astra.StatusEnumTERMINATED, astra.StatusEnumTERMINATING, astra.StatusEnumUNKNOWN)
	return err
}

// ParkAsync parks the database at the specified id. Note you cannot park a serverless database
// * @param databaseID string representation of the database ID
// @return error
func (a *AuthenticatedClient) ParkAsync(databaseID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/park", dbURL(), databaseID), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed creating request to park db with id %s with: %w", databaseID, err)
	}
	a.setHeaders(req)
	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to park database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != http.StatusAccepted {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body for park operation to db %v with error '%v'. http status of request was %v", databaseID, err, res.StatusCode)
		}
		return handleErrors(b, res.Status)
	}
	return nil
}

// Park parks the database at the specified id and will block until the database is parked
// * @param databaseID string representation of the database ID
// @return error
func (a *AuthenticatedClient) Park(databaseID string) error {
	err := a.ParkAsync(databaseID)
	if err != nil {
		return fmt.Errorf("park db failed because '%v'", err)
	}
	tries := 30
	interval := 30
	_, err = a.WaitUntil(databaseID, tries, interval, astra.StatusEnumPARKED)
	if err != nil {
		return fmt.Errorf("unable to check status for park db because of error '%v'", err)
	}
	return nil
}

// UnparkAsync unparks the database at the specified id. NOTE you cannot unpark a serverless database
// * @param databaseID String representation of the database ID
// @return error
func (a *AuthenticatedClient) UnparkAsync(databaseID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/unpark", dbURL(), databaseID), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed creating request to unpark db with id %s with: %w", databaseID, err)
	}
	a.setHeaders(req)
	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to unpark database id %s with: %w", databaseID, err)
	}
	defer closeBody(res)
	if res.StatusCode != http.StatusAccepted {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body for unpark operation to db %v with error '%v'. http status of request was %v", databaseID, err, res.StatusCode)
		}
		return handleErrors(b, res.Status)
	}
	return nil
}

// Unpark unparks the database at the specified id and will block until the database is unparked
// * @param databaseID String representation of the database ID
// @return error
func (a *AuthenticatedClient) Unpark(databaseID string) error {
	err := a.UnparkAsync(databaseID)
	if err != nil {
		return fmt.Errorf("unpark db failed because '%v'", err)
	}
	tries := 60
	interval := 30
	_, err = a.WaitUntil(databaseID, tries, interval, astra.StatusEnumACTIVE)
	if err != nil {
		return fmt.Errorf("unable to check status for unpark db because of error '%v'", err)
	}
	return nil
}

// Resize a database. Total number of capacity units desired should be specified. Reducing a size of a database is not supported at this time. Note you cannot resize a serverless database
// * @param databaseID string representation of the database ID
// * @param capacityUnits int32 containing capacityUnits key with a value greater than the current number of capacity units (max increment of 3 additional capacity units)
// @return error
func (a *AuthenticatedClient) Resize(databaseID string, capacityUnits int) error {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.ResizeDatabaseWithResponse(ctx, astra.DatabaseIdParam(databaseID), astra.ResizeDatabaseJSONRequestBody{
		CapacityUnits: &capacityUnits,
	})
	if err != nil {
		return fmt.Errorf("failed to resize database for database id %s with: %w", databaseID, err)
	}
	if res.StatusCode() != http.StatusAccepted {
		return handleErrors(res.Body, res.Status())
	}
	return nil
}

// ResetPassword changes the password for the database at the specified id
// * @param databaseID string representation of the database ID
// * @param username string containing username
// * @param password string containing password. The specified password will be updated for the specified database user
// @return error
func (a *AuthenticatedClient) ResetPassword(databaseID, username, password string) error {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.ResetPasswordWithResponse(ctx, astra.DatabaseIdParam(databaseID), astra.ResetPasswordJSONRequestBody{
		Username: astra.StringPtr(username),
		Password: astra.StringPtr(password),
	})
	if err != nil {
		return fmt.Errorf("failed to reset password for database id %s with: %w", databaseID, err)
	}

	if res.StatusCode() != http.StatusOK {
		return handleErrors(res.Body, res.Status())
	}
	return nil
}

// GetTierInfo Returns all supported tier, cloud, region, count, and capacitity combinations
// @return ([]TierInfo, error)
func (a *AuthenticatedClient) GetTierInfo() ([]astra.AvailableRegionCombination, error) {
	ctx, cancel := a.ctx()
	defer cancel()
	res, err := a.astraclient.ListAvailableRegionsWithResponse(ctx)
	if err != nil {
		return []astra.AvailableRegionCombination{}, fmt.Errorf("failed listing tier info with: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		return []astra.AvailableRegionCombination{}, handleErrors(res.Body, res.Status())
	}
	return *res.JSON200, nil
}
