//   Copyright 2022 DataStax
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
	"fmt"
	"os"

	"github.com/datastax-labs/astra-cli/pkg"
	"github.com/spf13/cobra"
)

// DeleteCmd provides the delete database command
var DeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete database by databaseID",
	Long:  `deletes a database from your Astra account by ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		msg, err := executeDelete(args, creds.Login)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, msg)
	},
}

func executeDelete(args []string, makeClient func() (pkg.Client, error)) (string, error) {
	client, err := makeClient()
	if err != nil {
		return "", fmt.Errorf("unable to login with error '%v'", err)
	}
	id := args[0]
	fmt.Printf("starting to delete database %v\n", id)
	if err := client.Terminate(id, false); err != nil {
		return "", fmt.Errorf("unable to delete '%s' with error %v", id, err)
	}
	return fmt.Sprintf("database %v deleted", id), nil
}
