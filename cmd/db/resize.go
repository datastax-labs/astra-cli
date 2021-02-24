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
	"strconv"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

// ResizeUsage shows the help for the delete command
func ResizeUsage() string {
	return "\tresize <id> <capacity unit>\n\t\tresizes a database by id with the specified capacity unit\n"
}

// ExecuteResize resizes the database with the specified ID with the specified size. If no ID is provided
// the command will error out
func ExecuteResize(args []string, client *astraops.AuthenticatedClient) error {
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for resizing the database"),
		}
	}
	if len(args) == 1 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for resizing the database"),
		}
	}
	id := args[0]
	capacityUnitRaw := args[1]
	capacityUnit, err := strconv.Atoi(capacityUnitRaw)
	if err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("unable to parse capacity unit '%s' with error %v\n", capacityUnitRaw, err),
		}
	}
	if err := client.Resize(id, int32(capacityUnit)); err != nil {
		return fmt.Errorf("unable to unpark '%s' with error %v\n", id, err)
	}
	fmt.Printf("resize database %v submitted with size %v\n", id, capacityUnit)
	return nil
}
