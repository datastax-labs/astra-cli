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

//Package cmd is the entry point for all of the commands
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Stores credentials for the cli to use in other commands to operate on the Astra DevOps API",
	Long:  `Token or service account is saved in .config/astra/ for use by the other commands`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		confDir, confFiles, err := pkg.GetHome(os.UserHomeDir)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(3)
		}
		if authToken != "" {
			if err := makeConf(confDir, confFiles.TokenPath, authToken); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		} else if clientJSON != "" {
			var clientInfo astraops.ClientInfo
			err := json.Unmarshal([]byte(clientJSON), &clientInfo)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("unable to serialize the json into a valid login due to error %s", err))
				os.Exit(1)
			}
			if len(clientInfo.ClientName) == 0 {
				fmt.Fprintln(os.Stderr, pkg.ParseError{
					Args: args,
					Err:  fmt.Errorf("clientName missing"),
				})
				os.Exit(1)
			}
			if len(clientInfo.ClientID) == 0 {
				fmt.Fprintln(os.Stderr, pkg.ParseError{
					Args: args,
					Err:  fmt.Errorf("clientId missing"),
				})
				os.Exit(1)
			}
			if len(clientInfo.ClientSecret) == 0 {
				fmt.Fprintln(os.Stderr, pkg.ParseError{
					Args: args,
					Err:  fmt.Errorf("clientSecret missing"),
				})
				os.Exit(1)
			}
			if err := makeConf(confDir, confFiles.SaPath, clientJSON); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		} else {
			clientJSON = fmt.Sprintf("{\"clientId\":\"%v\",\"clientName\":\"%v\",\"clientSecret\":\"%v:\"}", clientID, clientName, clientSecret)
			if err := makeConf(confDir, confFiles.SaPath, clientJSON); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			return
		}
	},
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
	//safe to write after validation
	_, err = writer.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("error writing file")
	}
	writer.Flush()
	fmt.Println("Login information saved")
	return nil
}
