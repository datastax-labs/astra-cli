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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// ClientInfo provides access to
type ClientInfo struct {
	ClientSecret string
	ClientName   string
	ClientID     string
}

// ConfFiles supports both formats of credentials and will say if the token one is present
type ConfFiles struct {
	TokenPath string
	SaPath    string
}

// HasServiceAccount returns true if there is a service account file present and accessible
func (c ConfFiles) HasServiceAccount() (bool, error) {
	if _, err := os.Stat(c.SaPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("warning error of %v is unexpected", err)
	}
	return true, nil
}

// HasToken returns true if there is a token file present and accessible
func (c ConfFiles) HasToken() (bool, error) {
	if _, err := os.Stat(c.TokenPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("warning error of %v is unexpected", err)
	}
	return true, nil
}

// GetHome returns the configuration directory and file
// error will return if there is no user home folder
func GetHome(getHome func() (string, error)) (confDir string, confFiles ConfFiles, err error) {
	var home string
	home, err = getHome()
	if err != nil {
		return "", ConfFiles{}, fmt.Errorf("unable to get user home directory with error '%s'", err)
	}
	confDir = path.Join(home, ".config", "astra")

	tokenFile := path.Join(confDir, PathWithEnv("token"))
	saFile := path.Join(confDir, PathWithEnv("sa.json"))
	return confDir, ConfFiles{
		TokenPath: tokenFile,
		SaPath:    saFile,
	}, nil
}

func PathWithEnv(f string) string {
	if strings.Contains(f, string(os.PathSeparator)) {
		tokens := strings.Split(f, string(os.PathSeparator))
		tokenLen := len(tokens)
		if tokenLen > 0 {
			last := tokens[tokenLen-1]
			tokens[tokenLen-1] = Env + "_" + last
			return strings.Join(tokens, string(os.PathSeparator))
		}
	}
	return Env + "_" + f
}

// ReadToken retrieves the login from the specified json file
func ReadToken(tokenFile string) (string, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return "", &FileNotFoundError{
			Path: tokenFile,
			Err:  fmt.Errorf("unable to read login file with error '%w'", err),
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning unable to close %v with error '%v'", tokenFile, err)
		}
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("unable to read login file '%s' with error '%w'", tokenFile, err)
	}
	if len(b) == 0 {
		return "", fmt.Errorf("token file '%s' is empty", tokenFile)
	}
	token := strings.Trim(string(b), "\n")
	if !strings.HasPrefix(token, "AstraCS") {
		return "", fmt.Errorf("missing prefix 'AstraCS' in token file '%s'", tokenFile)
	}
	return token, nil
}

// ReadLogin retrieves the login from the specified json file
func ReadLogin(saJSONFile string) (ClientInfo, error) {
	f, err := os.Open(saJSONFile)
	if err != nil {
		return ClientInfo{}, &FileNotFoundError{
			Path: saJSONFile,
			Err:  fmt.Errorf("unable to read login file with error %w", err),
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning unable to close %v with error %v", saJSONFile, err)
		}
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		return ClientInfo{}, fmt.Errorf("unable to read login file %s with error %w", saJSONFile, err)
	}
	var clientInfo ClientInfo
	err = json.Unmarshal(b, &clientInfo)
	if err != nil {
		return ClientInfo{}, &JSONParseError{
			Original: string(b),
			Err:      fmt.Errorf("unable to parse json from login file %s with error %s", saJSONFile, err),
		}
	}
	if clientInfo.ClientID == "" {
		return ClientInfo{}, fmt.Errorf("Invalid service account: Client ID for service account is empty for file '%v'", saJSONFile)
	}
	if clientInfo.ClientName == "" {
		return ClientInfo{}, fmt.Errorf("Invalid service account: Client name for service account is empty for file '%v'", saJSONFile)
	}
	if clientInfo.ClientSecret == "" {
		return ClientInfo{}, fmt.Errorf("Invalid service account: Client secret for service account is empty for file '%v'", saJSONFile)
	}
	return clientInfo, err
}
