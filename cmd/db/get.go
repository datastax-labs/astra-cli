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
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

var getCmd = flag.NewFlagSet("get", flag.ExitOnError)
var getFmt = getCmd.String("format", "text", "Output format for report default is json")

// GetUsage shows the help for the get command
func GetUsage() string {
	var out strings.Builder
	out.WriteString("\tget <id>\n")
	getCmd.VisitAll(func(f *flag.Flag) {
		out.WriteString(fmt.Sprintf("\t\t%v\n", f.Usage))
	})
	return out.String()
}

// ExecuteGet get the database with the specified ID. If no ID is provided
// the command will error out
func ExecuteGet(args []string, client *astraops.AuthenticatedClient) error {
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no id provided for parking the database"),
		}
	}
	if err := getCmd.Parse(args[1:]); err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  err,
		}
	}
	id := args[0]
	var db astraops.DataBase
	var err error
	if db, err = client.FindDb(id); err != nil {
		return fmt.Errorf("unable to get '%s' with error %v\n", id, err)
	}
	fmt.Println(*getFmt)
	switch *getFmt {
	case "text":
		var rows [][]string
		rows = append(rows, []string{"name", "id", "status"})
		rows = append(rows, []string{db.Info.Name, db.ID, db.Status})
		for _, row := range pkg.PadColumns(rows) {
			fmt.Println(strings.Join(row, " "))
		}
	case "json":
		b, err := json.MarshalIndent(db, "", "  ")
		if err != nil {
			return fmt.Errorf("unexpected error marshaling to json: '%v', Try -format text instead", err)
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("-format %q is not valid option.", *getFmt)
	}
	return nil
}
