/**
   Copyright 2021 Ryan Svihla

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
	"os"
	"strings"
)

var loginCmd = flag.NewFlagSet("login", flag.ExitOnError)
var clientIDFlag = loginCmd.String("id", "", "clientId from service account. Ignored if -json flag is used.")
var clientNameFlag = loginCmd.String("name", "", "clientName from service account. Ignored if -json flag is used.")
var clientSecretFlag = loginCmd.String("secret", "", "clientSecret from service account. Ignored if -json flag is used.")
var clientJSONFlag = loginCmd.String("json", "", "copy the json for service account from the Astra page")

//LoginUsage returns the usage text for login
func LoginUsage() string {
	var out strings.Builder
	out.WriteString("\tastra-cli login\n")
	loginCmd.VisitAll(func(f *flag.Flag) {
		out.WriteString(fmt.Sprintf("\t\t%v\n", f.Usage))
	})
	return out.String()
}

//ExecuteLogin logs into Astra
func ExecuteLogin(args []string, confDir string, confFile string) error {
	if err := loginCmd.Parse(args); err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("incorrect options with error %v", err),
		}
	}
	var clientJSON string
	if clientJSONFlag != nil {
		clientJSON = *clientJSONFlag
		var clientInfo astraops.ClientInfo
		err := json.Unmarshal([]byte(clientJSON), &clientInfo)
		if err != nil {
			return fmt.Errorf("unable to serialize the json into a valid login due to error %s", err)
		}

		if len(clientInfo.ClientName) == 0 {
			return &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientName missing"),
			}
		}
		if len(clientInfo.ClientID) == 0 {
			return &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientId missing"),
			}
		}
		if len(clientInfo.ClientSecret) == 0 {
			return &pkg.ParseError{
				Args: args,
				Err:  fmt.Errorf("clientSecret missing"),
			}
		}

	} else {
		clientID := *clientIDFlag
		clientName := *clientNameFlag
		clientSecret := *clientSecretFlag
		clientJSON = fmt.Sprintf("{\"clientId\":\"%v\",\"clientName\":\"%v\",\"clientSecret\":\"%v:\"}", clientID, clientName, clientSecret)
	}
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
	_, err = writer.Write([]byte(clientJSON))
	if err != nil {
		return fmt.Errorf("error writing file")
	}
	writer.Flush()
	fmt.Println("Login information saved")
	return nil
}
