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
	astraops "github.com/datastax/astra-client-go/v2/astra"
	"github.com/spf13/cobra"
)

var createDbName string
var createDbKeyspace string
var createDbRegion string
var createDbTier string
var createDbCloudProvider string

func init() {
	CreateCmd.Flags().StringVarP(&createDbName, "name", "n", "", "name to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbKeyspace, "keyspace", "k", "", "keyspace user to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbRegion, "region", "r", "us-east1", "region to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbTier, "tier", "t", "serverless", "tier to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbCloudProvider, "cloudProvider", "l", "GCP", "cloud provider flag to give to the Astra Database")
}

// CreateCmd creates a database in Astra
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a database by id",
	Long:  `creates a database by id`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		err := executeCreate(creds.Login)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func executeCreate(makeClient func() (pkg.Client, error)) error {
	client, err := makeClient()
	if err != nil {
		return fmt.Errorf("unable to login with error %v", err)
	}
	createDb := astraops.DatabaseInfoCreate{
		Name:          createDbName,
		Keyspace:      createDbKeyspace,
		CapacityUnits: 1, // we only support 1 CU on initial creation as of Feb 14 2022
		Region:        createDbRegion,
		Tier:          astraops.Tier(createDbTier),
		CloudProvider: astraops.CloudProvider(createDbCloudProvider),
	}
	db, err := client.CreateDb(createDb)
	if err != nil {
		return fmt.Errorf("unable to create '%v' with error %v", createDb, err)
	}
	fmt.Printf("database %v created\n", db.Id)
	return nil
}
