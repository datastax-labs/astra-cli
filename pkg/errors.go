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
	"fmt"
	"strings"
)

//ParseError is used to indicate there is an error in the command line args
type ParseError struct {
	Args []string
	Err  error
}

//Error outpus the error with the args provided, if there are no args that becomes the error
func (p *ParseError) Error() string {
	if len(p.Args) == 0 {
		return "no args provided"
	}
	return fmt.Sprintf("Unable to parse command line with args: %v. Nested error was '%v'", strings.Join(p.Args, ", "), p.Err)
}

//JSONParseError when unable to read JSON
type JSONParseError struct {
	Original string
	Err      error
}

func (j *JSONParseError) Error() string {
	return fmt.Sprintf("JSON parsing error: %s. Original file %s", j.Err, j.Original)
}

//FileNotFoundError when unable to read file
type FileNotFoundError struct {
	Path string
	Err  error
}

func (j *FileNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find file error: %s. Path to file %s", j.Err, j.Path)
}
