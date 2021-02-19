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

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rsds143/astra-mgmt-go/astraops"
	"io"
	"os"
	"path"
	"strconv"
)

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	clientIDFlag := loginCmd.String("id", "", "clientId from service account. Ignored if -json flag is used.")
	clientNameFlag := loginCmd.String("name", "", "clientName from service account. Ignored if -json flag is used.")
	clientSecretFlag := loginCmd.String("secret", "", "clientSecret from service account. Ignored if -json flag is used.")
	clientJSONFlag := loginCmd.String("json", "", "copy the json for service account from the Astra page")
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createDbNameFlag := createCmd.String("name", "", "name to give to the Astra Database")
	createDbKeyspaceFlag := createCmd.String("keyspace", "", "keyspace user to give to the Astra Database")
	createDbUserFlag := createCmd.String("user", "", "user password to give to the Astra Database")
	createDbPasswordFlag := createCmd.String("password", "", "db password to give to the Astra Database")
	createDbRegionFlag := createCmd.String("region", "us-east1", "region to give to the Astra Database")
	createDbTierFlag := createCmd.String("tier", "free", "tier to give to the Astra Database")
	createDbCapacityUnitFlag := createCmd.Int("capacityUnit", 1, "capacityUnit flag to give to the Astra Database")
	createDbCloudProviderFlag := createCmd.String("cloudProvider", "GCP", "cloud provider flag to give to the Astra Database")
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getFmt := getCmd.String("format", "text", "Output format for report default is json")
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	limitFlag := listCmd.Int("limit", 10, "limit of databases retrieved")
	includeFlag := listCmd.String("include", "", "the type of filter to apply")
	providerFlag := listCmd.String("provider", "", "provider to filter by")
	startingAfterFlag := listCmd.String("startingAfter", "", "timestamp filter, ie only show databases created after this timestamp")
	listFmt := listCmd.String("format", "text", "Output format for report default is json")
	tiersCmd := flag.NewFlagSet("tiers", flag.ExitOnError)
	tiersFmt := tiersCmd.String("format", "text", "Output format for report default is json")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("unable to get user home directory with error %s\n", err)
		os.Exit(2)
	}
	confDir := path.Join(home, ".config", "astra")
	confFile := path.Join(confDir, "sa.json")
	showDBUsage := func() {
		fmt.Println("db create")
		createCmd.PrintDefaults()
		fmt.Println("db delete <id>")
		fmt.Println("db park <id>")
		fmt.Println("db unpark <id>")
		fmt.Println("db resize <id> <capacity unit>")
		fmt.Println("db get <id>")
		getCmd.PrintDefaults()
		fmt.Println("db list")
		listCmd.PrintDefaults()
		fmt.Println("db tiers")
		tiersCmd.PrintDefaults()
	}
	if len(os.Args) == 1 {
		flag.Usage()
		fmt.Println("login <json> note:json is optional")
		loginCmd.PrintDefaults()
		showDBUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "login":
		if err := loginCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("login <json> //note:json is optional")
			loginCmd.PrintDefaults()
			fmt.Printf("incorrect options with error %v\n", err)
			os.Exit(2)
		}
		var clientJSON string
		if clientJSONFlag != nil {
			clientJSON = *clientJSONFlag
			var clientInfo astraops.ClientInfo
			err = json.Unmarshal([]byte(clientJSON), &clientInfo)
			if err != nil {
				fmt.Printf("unable to serialize the json into a valid login due to error %s\n", err)
				os.Exit(2)
			}

			if len(clientInfo.ClientName) == 0 {
				fmt.Println("clientName missing")
				os.Exit(2)
			}
			if len(clientInfo.ClientID) == 0 {
				fmt.Println("clientId missing")
				os.Exit(2)
			}
			if len(clientInfo.ClientSecret) == 0 {
				fmt.Println("clientSecret missing")
				os.Exit(2)
			}

		} else {
			clientID := *clientIDFlag
			clientName := *clientNameFlag
			clientSecret := *clientSecretFlag
			clientJSON = fmt.Sprintf("{\"clientId\":\"%v\",\"clientName\":\"%v\",\"clientSecret\":\"%v:\"}", clientID, clientName, clientSecret)
		}
		if err = os.MkdirAll(confDir, 0600); err != nil {
			fmt.Printf("unable to get make config directory with error %s\n", err)
			os.Exit(2)
		}
		f, err := os.Create(confFile)
		if err != nil {
			fmt.Printf("unable to create the login file due to error %s\n", err)
			os.Exit(2)
		}
		defer func() {
			if err = f.Close(); err != nil {
				fmt.Printf("failed unable to write file with error %s\n", err)
			}
		}()
		writer := bufio.NewWriter(f)
		//safe to write after validation
		_, err = writer.Write([]byte(clientJSON))
		if err != nil {
			fmt.Printf("error writing file\n")
			os.Exit(2)
		}
		writer.Flush()
		fmt.Println("Login information saved")
	case "db":
		clientInfo, err := ReadLogin(confFile)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		client, err := astraops.Authenticate(clientInfo.ClientName, clientInfo.ClientID, clientInfo.ClientSecret)
		if err != nil {
			fmt.Printf("authenticate failed with error %v", err)
			os.Exit(2)
		}
		if len(os.Args) == 2 {
			flag.Usage()
			showDBUsage()
			os.Exit(1)
		}
		switch os.Args[2] {
		case "create":
			if err := createCmd.Parse(os.Args[3:]); err != nil {
				fmt.Println(err)
				os.Exit(2)
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
				fmt.Printf("unable to create '%v' with error %v\n", createDb, err)
				os.Exit(2)
			}
			fmt.Printf("database %v created\n", id)
		case "delete":
			id := os.Args[3]
			fmt.Printf("starting to delete database %v\n", id)
			if err := client.Terminate(id, false); err != nil {
				fmt.Printf("unable to delete '%s' with error %v\n", id, err)
				os.Exit(2)
			}
			fmt.Printf("database %v deleted\n", id)
		case "park":
			id := os.Args[3]
			fmt.Printf("starting to park database %v\n", id)
			if err := client.Park(id); err != nil {
				fmt.Printf("unable to park '%s' with error %v\n", id, err)
				os.Exit(2)
			}
			fmt.Printf("database %v parked\n", id)
		case "unpark":
			id := os.Args[3]
			fmt.Printf("starting to unpark database %v\n", id)
			if err := client.UnPark(id); err != nil {
				fmt.Printf("unable to unpark '%s' with error %v\n", id, err)
				os.Exit(2)
			}
			fmt.Printf("database %v unparked\n", id)
		case "resize":
			id := os.Args[3]
			capacityUnitRaw := os.Args[4]
			capacityUnit, err := strconv.Atoi(capacityUnitRaw)
			if err != nil {
				fmt.Printf("unable to parse capacity unit '%s' with error %v\n", capacityUnitRaw, err)
				os.Exit(3)
			}
			if err := client.Resize(id, int32(capacityUnit)); err != nil {
				fmt.Printf("unable to unpark '%s' with error %v\n", id, err)
				os.Exit(2)
			}
			fmt.Printf("resize database %v submitted with size %v\n", id, capacityUnit)
		case "get":
			if err := getCmd.Parse(os.Args[3:]); err != nil {
				fmt.Println(err)
				os.Exit(2)

			}
			id := os.Args[3]
			var db astraops.DataBase
			if db, err = client.FindDb(id); err != nil {
				fmt.Printf("unable to get '%s' with error %v\n", id, err)
				os.Exit(2)
			}
			switch *getFmt {
			case "text":
				fmt.Println("name\tid\tstatus")
				fmt.Printf("%v\t%v\t%v\n", db.Info.Name, db.ID, db.Status)
			case "json":
				b, err := json.MarshalIndent(db, "", "  ")
				if err != nil {
					fmt.Printf("unexpected error marshaling to json: '%v', Try -format text instead", err)
					os.Exit(2)
				}
				fmt.Println(string(b))
			default:
				fmt.Printf("-format %q is not valid option.\n", *getFmt)
				os.Exit(2)
			}
		case "list":
			if err := listCmd.Parse(os.Args[3:]); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			var dbs []astraops.DataBase
			if dbs, err = client.ListDb(*includeFlag, *providerFlag, *startingAfterFlag, int32(*limitFlag)); err != nil {
				fmt.Printf("unable to get list of dbs with error %v\n", err)
				os.Exit(2)
			}
			switch *listFmt {
			case "text":
				fmt.Println("name\tid\tstatus")
				for _, db := range dbs {
					fmt.Printf("%v\t%v\t%v\n", db.Info.Name, db.ID, db.Status)
				}
			case "json":
				b, err := json.MarshalIndent(dbs, "", "  ")
				if err != nil {
					fmt.Printf("unexpected error marshaling to json: '%v', Try -format text instead", err)
					os.Exit(2)
				}
				fmt.Println(string(b))
			default:
				fmt.Printf("-format %q is not valid option.\n", *getFmt)
				os.Exit(2)
			}

		case "tiers":
			if err := tiersCmd.Parse(os.Args[3:]); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			var tiers []astraops.TierInfo
			if tiers, err = client.GetTierInfo(); err != nil {
				fmt.Printf("unable to get tiers with error %v\n", err)
				os.Exit(2)
			}

			switch *tiersFmt {
			case "text":
				fmt.Println("name\tcloud\tregion\tdb (used)/(limit)\tcap (used)/(limit)")
				for _, tier := range tiers {
					fmt.Printf("%v\t%v\t%v\t%v/%v\t%v/%v\n", tier.Tier, tier.CloudProvider, tier.RegionDisplay, tier.DatabaseCountUsed, tier.DatabaseCountLimit, tier.CapacityUnitsUsed, tier.CapacityUnitsLimit)
				}
			case "json":
				b, err := json.MarshalIndent(tiers, "", "  ")
				if err != nil {
					fmt.Printf("unexpected error marshaling to json: '%v', Try -format text instead", err)
					os.Exit(2)
				}
				fmt.Println(string(b))
			default:
				fmt.Printf("-format %q is not valid option.\n", *tiersFmt)
				os.Exit(2)
			}
		}
	default:
		fmt.Printf("db %q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

// ReadLogin retrieves the login from the specified json file
func ReadLogin(saJsonFile string) (astraops.ClientInfo, error) {
	f, err := os.Open(saJsonFile)
	if err != nil {
		return astraops.ClientInfo{}, fmt.Errorf("unable to read login file %s with error %s", saJsonFile, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning unable to close %v with error %v", saJsonFile, err)
		}
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		return astraops.ClientInfo{}, fmt.Errorf("unable to read login file %s with error %s", saJsonFile, err)
	}
	var clientInfo astraops.ClientInfo
	err = json.Unmarshal(b, &clientInfo)
	if err != nil {
		return astraops.ClientInfo{}, fmt.Errorf("unable to parse json from login file %s with error %s", saJsonFile, err)
	}
	return clientInfo, err
}
