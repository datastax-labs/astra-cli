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
package pkg

import (
	"flag"
	"fmt"
	"strings"
)

//PrintFlags outputs all of the flag parameters for a flagset
func PrintFlags(flagSet *flag.FlagSet, name string, desc string) string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("\t%v #%v\n", name, desc))
	var flags [][]string
	flagSet.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != "" {
			flags = append(flags, []string{
				fmt.Sprintf("-%v", f.Name),
				f.Usage,
				fmt.Sprintf("default: %v", f.Value),
			})
		} else {
			flags = append(flags, []string{
				fmt.Sprintf("-%v", f.Name),
				f.Usage,
				"",
			})
		}
	})
	for _, row := range PadColumns(flags) {
		out.WriteString(fmt.Sprintf("\t\t  %v\n", strings.Join(row, "  ")))
	}
	return out.String()
}
