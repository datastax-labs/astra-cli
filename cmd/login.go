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

// Package cmd is the entry point for all of the commands
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/spf13/cobra"
)

var clientID string
var clientName string
var clientSecret string
var clientJSON string
var authToken string

func init() {
	loginCmd.Flags().StringVarP(&authToken, "token", "t", "", "authtoken generated with enough rights to perform the devops actions. Generated from the Astra site")
	loginCmd.Flags().StringVarP(&clientJSON, "json", "j", "", "copy the json for service account from the Astra site")
	loginCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "clientSecret from service account. Ignored if -json flag is used.")
	loginCmd.Flags().StringVarP(&clientName, "name", "n", "", "clientName from service account. Ignored if -json flag is used.")
	loginCmd.Flags().StringVarP(&clientID, "id", "i", "", "clientId from service account. Ignored if -json flag is used.")
}

const (
	CriticalError  = 1
	WriteError     = 2
	CannotFindHome = 3
	JSONError      = 4
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Stores credentials for the cli to use in other commands to operate on the Astra DevOps API",
	Long:  `Token or service account is saved in .config/astra/ for use by the other commands`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		exitCode, err := executeLogin(args, func() (string, pkg.ConfFiles, error) {
			return pkg.GetHome(os.UserHomeDir)
		}, cobraCmd.Usage)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	},
}

func executeLogin(args []string, getHome func() (string, pkg.ConfFiles, error), usageFunc func() error) (int, error) {
	if len(args) == 0 {
		if err := usageFunc(); err != nil {
			return CriticalError, fmt.Errorf("cannot show usage %v", err)
		}
		return 0, nil
	}
	confDir, confFiles, err := getHome()
	if err != nil {
		return CannotFindHome, err
	}
	switch {
	case authToken != "":
		if err := makeConf(confDir, confFiles.TokenPath, authToken); err != nil {
			return WriteError, err
		}
		return 0, nil
	case clientJSON != "":
		return executeLoginJSON(args, confDir, confFiles)
	default:
		if clientID == "" {
			return JSONError, &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientId missing"),
			}
		}
		if clientName == "" {
			return JSONError, &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientName missing"),
			}
		}
		if clientSecret == "" {
			return JSONError, &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientSecret missing"),
			}
		}
		clientJSON = fmt.Sprintf("{\"clientId\":\"%v\",\"clientName\":\"%v\",\"clientSecret\":\"%v\"}", clientID, clientName, clientSecret)
		if err := makeConf(confDir, confFiles.SaPath, clientJSON); err != nil {
			return WriteError, err
		}
		return 0, nil
	}
}

func executeLoginJSON(args []string, confDir string, confFiles pkg.ConfFiles) (int, error) {
	var clientInfo pkg.ClientInfo
	err := json.Unmarshal([]byte(clientJSON), &clientInfo)
	if err != nil {
		return JSONError, fmt.Errorf("unable to serialize the json into a valid login due to error %s", err)
	}
	if len(clientInfo.ClientName) == 0 {
		return JSONError, &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("clientName missing"),
		}
	}
	if len(clientInfo.ClientID) == 0 {
		return JSONError, &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("clientId missing"),
		}
	}
	if len(clientInfo.ClientSecret) == 0 {
		return JSONError, &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("clientSecret missing"),
		}
	}
	if err := makeConf(confDir, confFiles.SaPath, clientJSON); err != nil {
		return WriteError, err
	}
	return 0, nil
}

func makeConf(confDir, confFile, content string) error {
	if err := os.MkdirAll(confDir, 0700); err != nil {
		return fmt.Errorf("unable to get make config directory with error %s", err)
	}
	f, err := os.Create(confFile)
	if err != nil {
		return fmt.Errorf("unable to create the login file due to error %s", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Printf("failed unable to write file with error %s\n", err)
		}
	}()
	writer := bufio.NewWriter(f)
	// safe to write after validation
	_, err = writer.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("error writing file")
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error finishing file")
	}
	fmt.Println("Login information saved")
	return nil
}
