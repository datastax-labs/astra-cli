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
	"strconv"

	"github.com/datastax-labs/astra-cli/pkg"
	"github.com/spf13/cobra"
)

const noRequiredArgs = 2

// ResizeCmd provides the resize database command
var ResizeCmd = &cobra.Command{
	Use:   "resize <id> <capacity unit>",
	Short: "Resizes a database by id with the specified capacity unit",
	Long:  "Resizes a database by id with the specified capacity unit. Note does not work on serverless.",
	Args:  cobra.ExactArgs(noRequiredArgs),
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		err := executeResize(args, creds.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to resize with error %v\n", err)
			os.Exit(1)
		}
	},
}

// executeResize resizes the database with the specified ID with the specified size. If no ID is provided
// the command will error out
func executeResize(args []string, makeClient func() (pkg.Client, error)) error {
	client, err := makeClient()
	if err != nil {
		return fmt.Errorf("unable to login with error %v", err)
	}
	id := args[0]
	capacityUnitRaw := args[1]
	defaultCapacity := 10
	bits := 32
	capacityUnit, err := strconv.ParseInt(capacityUnitRaw, defaultCapacity, bits)
	if err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  fmt.Errorf("unable to parse capacity unit '%s' with error %v", capacityUnitRaw, err),
		}
	}
	if err := client.Resize(id, int(capacityUnit)); err != nil {
		return fmt.Errorf("unable to resize '%s' with error %v", id, err)
	}
	fmt.Printf("resize database %v submitted with size %v\n", id, capacityUnit)
	return nil
}
