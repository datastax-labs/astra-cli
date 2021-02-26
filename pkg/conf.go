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

package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
	"io"
	"os"
	"path"
)

// GetHome returns the configuration directory and file
// error will return if there is no user home folder
func GetHome() (confDir string, confFile string, err error) {
	var home string
	home, err = os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("unable to get user home directory with error %s", err)
	}
	confDir = path.Join(home, ".config", "astra")
	confFile = path.Join(confDir, "sa.json")
	return confDir, confFile, nil
}

// ReadLogin retrieves the login from the specified json file
func ReadLogin(saJSONFile string) (astraops.ClientInfo, error) {
	f, err := os.Open(saJSONFile)
	if err != nil {
		return astraops.ClientInfo{}, &FileNotFoundError{
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
		return astraops.ClientInfo{}, fmt.Errorf("unable to read login file %s with error %w", saJSONFile, err)
	}
	var clientInfo astraops.ClientInfo
	err = json.Unmarshal(b, &clientInfo)
	if err != nil {
		return astraops.ClientInfo{}, &JSONParseError{
			Original: string(b),
			Err:      fmt.Errorf("unable to parse json from login file %s with error %s", saJSONFile, err),
		}
	}
	return clientInfo, err
}
