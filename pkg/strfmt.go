//  Copyright 2021 Ryan Svihla
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

// Package pkg is the top level package for shared libraries
package pkg

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// WriteRows outputs a flexiable right aligned tabwriter
func WriteRows(w io.Writer, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	for i, row := range rows {
		rowStr := strings.Join(row, "\t")
		if i > 0 {
			fmt.Fprint(tw, "\n")
		}
		fmt.Fprint(tw, rowStr)
	}
	return tw.Flush()
}
