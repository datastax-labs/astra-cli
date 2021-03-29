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
	"errors"
	"fmt"
	"testing"
)

func TestParseErrorNoArgs(t *testing.T) {
	parseError := ParseError{
		Err: errors.New("bogus error"),
	}

	expected := "no args provided"
	if parseError.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, parseError.Error())
	}
}

func TestParseError(t *testing.T) {
	parseError := ParseError{
		Args: []string{"a", "b"},
		Err:  errors.New("bogus error"),
	}

	expected := "Unable to parse command line with args: a, b. Nested error was 'bogus error'"
	if parseError.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, parseError.Error())
	}
}

func TestFileNotFoundError(t *testing.T) {
	fileErr := FileNotFoundError{
		Path: "/a/b/C",
		Err:  fmt.Errorf("Bogus Error"),
	}
	expected := "Unable to find file '/a/b/C' with error: 'Bogus Error'"
	if fileErr.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, fileErr.Error())
	}
}

func TestJSONParseError(t *testing.T) {
	fileErr := JSONParseError{
		Original: "invalid string",
		Err:      fmt.Errorf("Bogus Error"),
	}
	expected := "JSON parsing error for json 'invalid string' with error 'Bogus Error'"
	if fileErr.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, fileErr.Error())
	}
}
