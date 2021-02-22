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

import "fmt"

// PadColumns pads all columns to be the same width across all rows
// rather than modify rows a new set is returned
func PadColumns(rows [][]string) [][]string {
	// find amount of columns in the biggest row
	var totalColumns int
	for _, row := range rows {
		rowLen := len(row)
		if rowLen > totalColumns {
			totalColumns = rowLen
		}
	}
	//now find the largest column size in each column
	biggestColumnSize := make([]int, totalColumns)
	for _, row := range rows {
		for ci, col := range row {
			longest := biggestColumnSize[ci]
			if len(col) > longest {
				biggestColumnSize[ci] = len(col)
			}
		}
	}
	//now pad each column, go ahead and add a column if it does not exist in the row and pad it
	paddedArray := make([][]string, len(rows))
	for i, row := range rows {
		newRow := make([]string, totalColumns)
		rowLength := len(row)
		for ci := 0; ci < totalColumns; ci++ {
			maxLength := biggestColumnSize[ci]
			col := ""
			if rowLength > ci {
				col = row[ci]
			}
			newRow[ci] = fmt.Sprintf("%-*s", maxLength, col)
		}
		paddedArray[i] = newRow
	}
	return paddedArray
}
