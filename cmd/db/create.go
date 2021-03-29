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
package db

import (
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
	"github.com/spf13/cobra"
)

var createDbName string
var createDbKeyspace string
var createDbUser string
var createDbPassword string
var createDbRegion string
var createDbTier string
var createDbCapacityUnit int
var createDbCloudProvider string

func init() {
	CreateCmd.Flags().StringVarP(&createDbName, "name", "n", "", "name to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbKeyspace, "keyspace", "k", "", "keyspace user to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbUser, "user", "u", "", "user password to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbPassword, "password", "p", "", "db password to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbRegion, "region", "r", "us-east1", "region to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbTier, "tier", "t", "serverless", "tier to give to the Astra Database")
	CreateCmd.Flags().IntVarP(&createDbCapacityUnit, "capacityUnit", "c", 1, "capacityUnit flag to give to the Astra Database")
	CreateCmd.Flags().StringVarP(&createDbCloudProvider, "cloudProvider", "l", "GCP", "cloud provider flag to give to the Astra Database")

}

//CreateCmd creates a database in Astra
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a database by id",
	Long:  ``,
	Run: func(cobraCmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		client, err := creds.Login()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to login with error %v\n", err)
			os.Exit(1)
		}

		err = executeCreate(client)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func executeCreate(client *astraops.AuthenticatedClient) error {
	capacity := int32(createDbCapacityUnit)
	createDb := astraops.CreateDb{
		Name:          createDbName,
		Keyspace:      createDbKeyspace,
		CapacityUnits: capacity,
		Region:        createDbRegion,
		User:          createDbUser,
		Password:      createDbPassword,
		Tier:          createDbTier,
		CloudProvider: createDbCloudProvider,
	}
	db, err := client.CreateDb(createDb)
	if err != nil {
		return fmt.Errorf("unable to create '%v' with error %v", createDb, err)
	}
	id := db.ID
	fmt.Printf("database %v created\n", id)
	return nil
}
