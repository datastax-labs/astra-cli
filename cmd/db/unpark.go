//  Copyright 2022 DataStax
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

// Package db is where the Astra DB commands are
package db

import (
	"fmt"
	"os"

	"github.com/datastax-labs/astra-cli/pkg"
	"github.com/spf13/cobra"
)

// UnparkCmd provides unparking support for classic database tiers in Astra
var UnparkCmd = &cobra.Command{
	Use:   "unpark <id>",
	Short: "parks the database specified, does not work with serverless",
	Long:  `parks the database specified, only works on classic tier databases and can take a very long time to park (20-30 minutes)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		err := executeUnpark(args, creds.Login)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// executeUnpark unparks the database with the specified ID. If no ID is provided
// the command will error out
func executeUnpark(args []string, makeClient func() (pkg.Client, error)) error {
	client, err := makeClient()
	if err != nil {
		return fmt.Errorf("unable to login with error %v", err)
	}
	id := args[0]
	fmt.Printf("starting to unpark database %v\n", id)
	if err := client.Unpark(id); err != nil {
		return fmt.Errorf("unable to unpark '%s' with error %v", id, err)
	}
	fmt.Printf("database %v unparked\n", id)
	return nil
}
