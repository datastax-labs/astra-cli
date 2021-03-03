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
	"fmt"
	"strings"

	"github.com/rsds143/astra-cli/cmd/db"
	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

func DBUsage() string {
	return strings.Join([]string{
		"\tastra-cli db <subcommands>",
		"\tdb subcommands:",
		fmt.Sprintf("\t%v", db.CreateUsage()),
		fmt.Sprintf("\t%v", db.DeleteUsage()),
		fmt.Sprintf("\t%v", db.GetUsage()),
		fmt.Sprintf("\t%v", db.ListUsage()),
		fmt.Sprintf("\t%v", db.ParkUsage()),
		fmt.Sprintf("\t%v", db.UnparkUsage()),
		fmt.Sprintf("\t%v", db.ResizeUsage()),
		fmt.Sprintf("\t%v", db.TiersUsage()),
	}, "\n")
}

// ExecuteDB launches several different subcommands and as of today is the main entry point
// into automation of Astra
func ExecuteDB(args []string, confFile string, verbose bool) error {
	clientInfo, err := pkg.ReadLogin(confFile)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	client, err := astraops.Authenticate(clientInfo, verbose)
	if err != nil {
		return fmt.Errorf("authenticate failed with error %v", err)
	}
	if len(args) == 0 {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("there is no standalone db command"),
		}
	}
	switch args[0] {
	case "create":
		return db.ExecuteCreate(args[1:], client)
	case "delete":
		return db.ExecuteDelete(args[1:], client)
	case "park":
		return db.ExecutePark(args[1:], client)
	case "unpark":
		return db.ExecuteUnpark(args[1:], client)
	case "resize":
		return db.ExecuteResize(args[1:], client)
	case "get":
		return db.ExecuteGet(args[1:], client)
	case "list":
		return db.ExecuteList(args[1:], client)
	case "tiers":
		return db.ExecuteTiers(args[1:], client)
	}
	return nil
}
