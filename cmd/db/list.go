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

var listCmd = flag.NewFlagSet("list", flag.ExitOnError)
var limitFlag = listCmd.Int("limit", 10, "limit of databases retrieved")
var includeFlag = listCmd.String("include", "", "the type of filter to apply")
var providerFlag = listCmd.String("provider", "", "provider to filter by")
var startingAfterFlag = listCmd.String("startingAfter", "", "timestamp filter, ie only show databases created after this timestamp")
var listFmt = listCmd.String("format", "text", "Output format for report default is json")

// ListUsage shows the help for the List command
func ListUsage() string {
	var out strings.Builder
	out.WriteString("\tlist\n")
	listCmd.VisitAll(func(f *flag.Flag) {
		out.WriteString(fmt.Sprintf("\t\t%v\n", f.Usage))
	})
	return out.String()
}

// ExecuteList lists databases in astra
func ExecuteList(args []string, client *astraops.AuthenticatedClient) error {
	if err := listCmd.Parse(args); err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  err,
		}
	}
	var dbs []astraops.DataBase
	var err error
	if dbs, err = client.ListDb(*includeFlag, *providerFlag, *startingAfterFlag, int32(*limitFlag)); err != nil {
		return fmt.Errorf("unable to get list of dbs with error %v", err)
	}
	switch *listFmt {
	case "text":
		var rows [][]string
		rows = append(rows, []string{"name", "id", "status"})
		for _, db := range dbs {
			rows = append(rows, []string{db.Info.Name, db.ID, db.Status})
		}
		for _, row := range pkg.PadColumns(rows) {
			fmt.Println(strings.Join(row, " "))
		}
	case "json":
		b, err := json.MarshalIndent(dbs, "", "  ")
		if err != nil {
			return fmt.Errorf("unexpected error marshaling to json: '%v', Try -format text instead", err)
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("-format %q is not valid option", *getFmt)
	}
	return nil
}
