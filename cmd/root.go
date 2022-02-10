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

// Package cmd contains all fo the commands for the cli
package cmd

import (
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-cli/pkg/env"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&env.Verbose, "verbose", "v", false, "turns on verbose logging")
	RootCmd.PersistentFlags().StringVarP(&pkg.Env, "env", "e", "prod", "environment to automate, other options are test and dev")
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(dbCmd)
}

// RootCmd is the entry point for the whole app
var RootCmd = &cobra.Command{
	Use:   "astra-cli",
	Short: "An easy to use client for automating DataStax Astra",
	Long: `Manage and provision databases on DataStax Astra
                Complete documentation is available at https://github.com/rsds143/astra-cli`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		if err := executeRoot(cobraCmd.Usage); err != nil {
			os.Exit(1)
		}
	},
}

func executeRoot(usage func() error) error {
	if err := usage(); err != nil {
		return fmt.Errorf("warn unable to show usage %v", err)
	}
	return nil
}
