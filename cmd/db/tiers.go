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
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/rsds143/astra-cli/pkg"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

var tiersCmd = flag.NewFlagSet("tiers", flag.ExitOnError)
var tiersFmt = tiersCmd.String("format", "text", "Output format for report default is json")

// TiersUsage shows the help for the Tiers command
func TiersUsage() string {
	var out strings.Builder
	out.WriteString("\ttiers\n")
	tiersCmd.VisitAll(func(f *flag.Flag) {
		out.WriteString(fmt.Sprintf("\t\t%v\n", f.Usage))
	})
	return out.String()
}

// ExecuteTiers lists tiers available to this login in astra
func ExecuteTiers(args []string, client *astraops.AuthenticatedClient) error {
	if err := tiersCmd.Parse(args); err != nil {
		return &pkg.ParseError{
			Args: args,
			Err:  err,
		}
	}
	var tiers []astraops.TierInfo
	var err error
	if tiers, err = client.GetTierInfo(); err != nil {
		return fmt.Errorf("unable to get tiers with error %v", err)
	}

	switch *tiersFmt {
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
				tier.RegionDisplay,
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
			return fmt.Errorf("unexpected error marshaling to json: '%v', Try -format text instead", err)
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("-format %q is not valid option.\n", *tiersFmt)
	}
	return nil
}
