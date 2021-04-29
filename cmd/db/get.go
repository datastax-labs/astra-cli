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

// Package db provides the sub-commands for the db command
package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	astraops "github.com/rsds143/astra-cli/pkg/swagger"
	"github.com/spf13/cobra"
)

var getFmt string

func init() {
	GetCmd.Flags().StringVarP(&getFmt, "output", "o", "text", "Output format for report default is text")
}

// GetCmd provides the get database command
var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "get database by databaseID",
	Long:  `gets a database from your Astra account by ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		txt, err := executeGet(args, creds.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to login with error %v\n", err)
			os.Exit(1)
		}
		fmt.Println(txt)
	},
}

func executeGet(args []string, login func() (pkg.Client, error)) (string, error) {
	client, err := login()
	if err != nil {
		return "", fmt.Errorf("unable to login with error %v", err)
	}
	id := args[0]
	var db astraops.Database
	if db, err = client.FindDb(id); err != nil {
		return "", fmt.Errorf("unable to get '%s' with error %v", id, err)
	}
	switch getFmt {
	case pkg.TextFormat:
		var rows [][]string
		rows = append(rows, []string{"name", "id", "status"})
		rows = append(rows, []string{db.Info.Name, db.ID, string(db.Status)})
		var buf bytes.Buffer
		err = pkg.WriteRows(&buf, rows)
		if err != nil {
			return "", fmt.Errorf("unexpected error writing out text %v", err)
		}
		return buf.String(), nil
	case pkg.JSONFormat:
		b, err := json.MarshalIndent(db, "", "  ")
		if err != nil {
			return "", fmt.Errorf("unexpected error marshaling to json: '%v', Try -output text instead", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("-o %q is not valid option", getFmt)
	}
}
