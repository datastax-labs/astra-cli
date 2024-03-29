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

// Package pkg is the top level package for shared libraries
package pkg

import (
	"bytes"
	"strings"
	"testing"
)

func TestTabWriterLayout(t *testing.T) {
	w := bytes.NewBufferString("")
	rows := [][]string{
		{
			"abc", "def", "ghi",
		},
		{
			"", "123456", "1",
		},
	}
	err := WriteRows(w, rows)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected1 := "abc def    ghi"
	expected2 := "    123456 1"
	expected := strings.Join([]string{expected1, expected2}, "\n")
	if w.String() != expected {
		t.Errorf("expected/actual \n'%v'\n'%v'", expected, w.String())
	}
}
