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
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-cli/pkg/httputils"
	astraops "github.com/datastax/astra-client-go/v2/astra"
	"github.com/spf13/cobra"
)

var secBundleFmt string
var secBundleLoc string

func init() {
	SecBundleCmd.Flags().StringVarP(&secBundleFmt, "output", "o", "zip", "Output format for report default is zip")
	SecBundleCmd.Flags().StringVarP(&secBundleLoc, "location", "l", "secureBundle.zip", "location of bundle to download to if using zip format. ignore if using json")
}

// SecBundleCmd  provides the secBundle database command
var SecBundleCmd = &cobra.Command{
	Use:   "secBundle <id>",
	Short: "get secure bundle by databaseID",
	Long:  `gets the secure connetion bundle for the database from your Astra account by ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		out, err := executeSecBundle(args, creds.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(out)
	},
}

func executeSecBundle(args []string, login func() (pkg.Client, error)) (string, error) {
	client, err := login()
	if err != nil {
		return "", fmt.Errorf("unable to login with error %v", err)
	}
	id := args[0]
	var secBundle astraops.CredsURL
	if secBundle, err = client.GetSecureBundle(id); err != nil {
		return "", fmt.Errorf("unable to get '%s' with error %v", id, err)
	}
	switch secBundleFmt {
	case "zip":
		bytesWritten, err := httputils.DownloadZip(secBundle.DownloadURL, secBundleLoc)
		if err != nil {
			return "", fmt.Errorf("error outputing zip format '%v'", err)
		}
		return fmt.Sprintf("file %v saved %v bytes written", secBundleLoc, bytesWritten), nil
	case pkg.JSONFormat:
		b, err := json.MarshalIndent(secBundle, "", "  ")
		if err != nil {
			return "", fmt.Errorf("unexpected error marshaling to json: '%v', Try -output text instead", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("-o %q is not valid option", secBundleFmt)
	}
}
