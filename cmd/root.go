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

//Package cmd contains all fo the commands for the cli
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/rsds143/astra-cli/pkg/env"
)

func init() {
  rootCmd.PersistentFlags().BoolVarP(&env.Verbose, "verbose", "v", false, "turns on verbose logging")
  rootCmd.AddCommand(loginCmd)
  rootCmd.AddCommand(dbCmd)
}

var rootCmd = &cobra.Command{
  Use:   "astra-cli",
  Short: "An easy to use client for automating DataStax Astra",
  Long: `Manage and provision databases on DataStax Astra
                Complete documentation is available at https://github.com/rsds143/astra-cli`,
  Run: func(cmd *cobra.Command, args []string) {
  },
}
