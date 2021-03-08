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

// UnparkUsage shows the help for the unpark command
func UnparkUsage() string {
	return "\tunpark <id> #parks a database by id\n"
}

// ExecuteUnpark unparks the database with the specified ID. If no ID is provided
// the command will error out
func ExecuteUnpark(args []string, client *astraops.AuthenticatedClient) error {
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for unparking the database"),
		}
	}
	id := args[0]
	fmt.Printf("starting to unpark database %v\n", id)
	if err := client.Unpark(id); err != nil {
		return fmt.Errorf("unable to unpark '%s' with error %v\n", id, err)
	}
	fmt.Printf("database %v unparked\n", id)
	return nil
}
