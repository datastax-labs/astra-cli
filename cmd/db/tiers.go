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
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rsds143/astra-cli/pkg"
	astraops "github.com/datastax/astra-client-go/v2/astra"
	"github.com/spf13/cobra"
)

var tiersFmt string

func init() {
	TiersCmd.Flags().StringVarP(&tiersFmt, "output", "o", "text", "Output format for report default is json")
}

// TiersCmd is the command to list availability data in Astra
var TiersCmd = &cobra.Command{
	Use:   "tiers",
	Short: "List all available tiers on the Astra DevOps API",
	Long:  `List all available tiers on the Astra DevOps API. Each tier is a combination of costs, size, region, and name`,
	Run: func(cmd *cobra.Command, args []string) {
		creds := &pkg.Creds{}
		msg, err := executeTiers(creds.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
	},
}

func executeTiers(login func() (pkg.Client, error)) (string, error) {
	var tiers []astraops.AvailableRegionCombination
	client, err := login()
	if err != nil {
		return "", fmt.Errorf("unable to login with error %v", err)
	}
	if tiers, err = client.GetTierInfo(); err != nil {
		return "", fmt.Errorf("unable to get tiers with error %v", err)
	}
	switch tiersFmt {
	case pkg.TextFormat:
		var rows [][]string
		rows = append(rows, []string{"name", "cloud", "region", "db (used)/(limit)", "cap (used)/(limit)", "cost per month", "cost per minute"})
		for _, tier := range tiers {
			var costMonthRaw float64
			var costMinRaw float64
			var emtpyCosts astraops.Costs 
			if tier.Cost != emtpyCosts {
				costMonthRaw = astraops.Float64Value(tier.Cost.CostPerMonthCents)
				costMinRaw = astraops.Float64Value(tier.Cost.CostPerMinCents)
			}
			divisor := 100.0
			var costMonth float64
			if costMonthRaw > 0.0 {
				costMonth = costMonthRaw / divisor
			}
			var costMin float64
			if costMinRaw > 0.0 {
				costMin = costMinRaw / divisor
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
		var buf bytes.Buffer
		err = pkg.WriteRows(&buf, rows)
		if err != nil {
			return "", fmt.Errorf("unexpected error writing text output %v", err)
		}
		return buf.String(), nil
	case pkg.JSONFormat:
		b, err := json.MarshalIndent(tiers, "", "  ")
		if err != nil {
			return "", fmt.Errorf("unexpected error marshaling to json: '%v', Try -format text instead", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("-o %q is not valid option", tiersFmt)
	}
}
