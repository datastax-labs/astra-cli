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
	"flag"
	"fmt"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

var createCmd = flag.NewFlagSet("create", flag.ExitOnError)
var createDbNameFlag = createCmd.String("name", "", "name to give to the Astra Database")
var createDbKeyspaceFlag = createCmd.String("keyspace", "", "keyspace user to give to the Astra Database")
var createDbUserFlag = createCmd.String("user", "", "user password to give to the Astra Database")
var createDbPasswordFlag = createCmd.String("password", "", "db password to give to the Astra Database")
var createDbRegionFlag = createCmd.String("region", "us-east1", "region to give to the Astra Database")
var createDbTierFlag = createCmd.String("tier", "serverless", "tier to give to the Astra Database")
var createDbCapacityUnitFlag = createCmd.Int("capacityUnit", 1, "capacityUnit flag to give to the Astra Database")
var createDbCloudProviderFlag = createCmd.String("cloudProvider", "GCP", "cloud provider flag to give to the Astra Database")

// CreateUsage shows the help for the create command
func CreateUsage() string {
	return pkg.PrintFlags(createCmd, "create", "creates a database by id")
}

// ExecuteCreate submits a new database to astra
func ExecuteCreate(args []string, client *astraops.AuthenticatedClient) error {
	if err := createCmd.Parse(args); err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  err,
		}
	}
	createDb := astraops.CreateDb{
		Name:          *createDbNameFlag,
		Keyspace:      *createDbKeyspaceFlag,
		CapacityUnits: *createDbCapacityUnitFlag,
		Region:        *createDbRegionFlag,
		User:          *createDbUserFlag,
		Password:      *createDbPasswordFlag,
		Tier:          *createDbTierFlag,
		CloudProvider: *createDbCloudProviderFlag,
	}
	id, _, err := client.CreateDb(createDb)
	if err != nil {
		return fmt.Errorf("unable to create '%v' with error %v", createDb, err)
	}
	fmt.Printf("database %v created\n", id)
	return nil
}
