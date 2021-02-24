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
package db

import (
	"fmt"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

// ParkUsage shows the help for the park command
func ParkUsage() string {
	return "\tpark <id>\n\t\tparks a database by id\n"
}

// ExecutePark parks the database with the specified ID. If no ID is provided
// the command will error out
func ExecutePark(args []string, client *astraops.AuthenticatedClient) error {
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for parking the database"),
		}
	}
	id := args[1]
	fmt.Printf("starting to park database %v\n", id)
	if err := client.Park(id); err != nil {
		return fmt.Errorf("unable to park '%s' with error %v\n", id, err)
	}
	fmt.Printf("database %v parked\n", id)
	return nil
}
