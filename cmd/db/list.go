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

var limit int
var include string
var provider string
var startingAfter string
var listFmt string

func init() {
	ListCmd.Flags().IntVarP(&limit, "limit", "l", 10, "limit of databases retrieved")
	ListCmd.Flags().StringVarP(&include, "include", "i", "", "the type of filter to apply")
	ListCmd.Flags().StringVarP(&provider, "provider", "p", "", "provider to filter by")
	ListCmd.Flags().StringVarP(&startingAfter, "startingAfter", "a", "", "timestamp filter, ie only show databases created after this timestamp")
	ListCmd.Flags().StringVarP(&listFmt, "output", "o", "text", "Output format for report default is json")
}

// ListCmd provides the list databases command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all databases",
	Long:  `lists all databases in your Astra account`,
	Run: func(cmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		msg, err := executeList(creds.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
	},
}

func executeList(login func() (pkg.Client, error)) (string, error) {
	client, err := login()
	if err != nil {
		return "", fmt.Errorf("unable to login with error '%v'", err)
	}
	var dbs []astraops.Database
	if dbs, err = client.ListDb(include, provider, startingAfter, int32(limit)); err != nil {
		return "", fmt.Errorf("unable to get list of dbs with error '%v'", err)
	}
	switch listFmt {
	case pkg.TextFormat:
		var rows [][]string
		rows = append(rows, []string{"name", "id", "status"})
		for _, db := range dbs {
			rows = append(rows, []string{db.Info.Name, db.ID, string(db.Status)})
		}
		var out bytes.Buffer
		err = pkg.WriteRows(&out, rows)
		if err != nil {
			return "", fmt.Errorf("unexpected error writing text output '%v'", err)
		}
		return out.String(), nil
	case pkg.JSONFormat:
		b, err := json.MarshalIndent(dbs, "", "  ")
		if err != nil {
			return "", fmt.Errorf("unexpected error marshaling to json: '%v', Try -output text instead", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("-o %q is not valid option", listFmt)
	}
}
