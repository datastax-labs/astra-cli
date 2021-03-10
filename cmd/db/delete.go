//   Copyright 2021 Ryan Svihla
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

//Package db provides the sub-commands for the db command
package db

import (
    "os"
	"fmt"

    "github.com/spf13/cobra"
	"github.com/rsds143/astra-cli/pkg"
)


//DeleteCmd provides the delete database command
var DeleteCmd =  &cobra.Command{
  Use:   "delete <id>",
  Short: "delete database by databaseID",
  Long: `deletes a database from your Astra account by ID`,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    client, err := pkg.LoginClient()
	if err != nil {
	    fmt.Fprintln(os.Stderr, fmt.Sprintf("unable to login with error %v", err))
        os.Exit(1)
    }
    id := args[0]
	fmt.Printf("starting to delete database %v\n", id)
	if err := client.Terminate(id, false); err != nil {
	    fmt.Fprintln(os.Stderr, fmt.Errorf("unable to delete '%s' with error %v", id, err))
	    os.Exit(1)
	}
	fmt.Printf("database %v deleted\n",id)
  },
}
