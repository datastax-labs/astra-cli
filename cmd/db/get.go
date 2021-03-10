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
	"encoding/json"
	"os"
	"fmt"
	"strings"

    "github.com/spf13/cobra"
	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

var getFmt string

func init () {
    GetCmd.Flags().StringVarP(&getFmt, "output", "o", "text", "Output format for report default is text")
}

//GetCmd provides the get database command
var GetCmd =  &cobra.Command{
  Use:   "get <id>",
  Short: "get database by databaseID",
  Long: `gets a database from your Astra account by ID`,
  Args: cobra.ExactArgs(1),
  Run: func(cobraCmd *cobra.Command, args []string) {
    client, err := pkg.LoginClient()
	if err != nil {
	    fmt.Fprintln(os.Stderr, fmt.Sprintf("unable to login with error %v", err))
        os.Exit(1)
    }
    id := args[0]
	var db astraops.Database
	if db, err = client.FindDb(id); err != nil {
	    fmt.Fprintln(os.Stderr, fmt.Sprintf("unable to get '%s' with error %v\n", id, err))
	    os.Exit(1)
	}
	switch getFmt {
	case "text":
		var rows [][]string
		rows = append(rows, []string{"name", "id", "status"})
		rows = append(rows, []string{db.Info.Name, db.ID, string(db.Status)})
		for _, row := range pkg.PadColumns(rows) {
			fmt.Println(strings.Join(row, " "))
		}
	case "json":
		b, err := json.MarshalIndent(db, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("unexpected error marshaling to json: '%v', Try -output text instead", err))
	    os.Exit(1)
		}
		fmt.Println(string(b))
	default:
	    fmt.Fprintln(os.Stderr, fmt.Sprintf("-output %q is not valid option.", getFmt))
	    os.Exit(1)
	}
    },
}
