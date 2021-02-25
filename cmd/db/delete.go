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

// DeleteUsage shows the help for the delete command
func DeleteUsage() string {
	return "\tdelete <id>\n\t\tdeletes a database by id\n"
}

// ExecuteDelete removes the database with the specified ID. If no ID is provided
// the command will error out
func ExecuteDelete(args []string, client *astraops.AuthenticatedClient) error {
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for deleting the database"),
		}
	}
	id := args[0]
	fmt.Printf("starting to delete database %v\n", id)
	if err := client.Terminate(id, false); err != nil {
		return fmt.Errorf("unable to delete '%s' with error %v", id, err)
	}
	fmt.Printf("database %v deleted\n", id)
	return nil
}
