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

//Package db provides the sub-commands for the db command
package db

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
	"github.com/spf13/cobra"
)

var tiersFmt string

func init() {
	TiersCmd.Flags().StringVarP(&tiersFmt, "output", "o", "text", "Output format for report default is json")
}

//TiersCmd is the command to list availability data in Astra
var TiersCmd = &cobra.Command{
	Use:   "tiers",
	Short: "List all available tiers on the Astra DevOps API",
	Long:  `List all available tiers on the Astra DevOps API. Each tier is a combination of costs, size, region, and name`,
	Run: func(cmd *cobra.Command, args []string) {
		var tiers []astraops.TierInfo
		client, err := pkg.LoginClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to login with error %v\n", err)
			os.Exit(1)
		}
		if tiers, err = client.GetTierInfo(); err != nil {
			fmt.Fprintf(os.Stderr, "unable to get tiers with error %v\n", err)
			os.Exit(1)
		}
		switch tiersFmt {
		case "text":
			var rows [][]string
			rows = append(rows, []string{"name", "cloud", "region", "db (used)/(limit)", "cap (used)/(limit)", "cost per month", "cost per minute"})
			for _, tier := range tiers {
				costMonthRaw := tier.Cost.CostPerMonthCents
				var costMonth float64
				if costMonthRaw > 0 {
					costMonth = costMonthRaw / 100.0
				}
				costMinRaw := tier.Cost.CostPerMinCents
				var costMin float64
				if costMinRaw > 0 {
					costMin = costMinRaw / 100.0
				}
				rows = append(rows, []string{
					tier.Tier,
					tier.CloudProvider,
					tier.Region,
					fmt.Sprintf("%v/%v", tier.DatabaseCountUsed, tier.DatabaseCountLimit),
					fmt.Sprintf("%v/%v", tier.CapacityUnitsUsed, tier.CapacityUnitsLimit),
					fmt.Sprintf("$%.2f", costMonth),
					fmt.Sprintf("$%.2f", costMin)})
			}
			for _, row := range pkg.PadColumns(rows) {
				fmt.Println(strings.Join(row, " "))
			}
		case "json":
			b, err := json.MarshalIndent(tiers, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "unexpected error marshaling to json: '%v', Try -format text instead\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
		default:
			fmt.Fprintf(os.Stderr, "-o %q is not valid option.\n", tiersFmt)
			os.Exit(1)
		}
	},
}
