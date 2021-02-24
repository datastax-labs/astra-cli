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
	"errors"
	"fmt"
	"github.com/rsds143/astra-cli/cmd"
	"github.com/rsds143/astra-cli/pkg"
	"os"
)

func usage() {
	fmt.Println("usage: astra-cli <cmd>")
	fmt.Println("commands:")
	fmt.Println(cmd.LoginUsage())
	fmt.Println(cmd.DBUsage())
}
func main() {
	confDir, confFile, err := pkg.GetHome()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(3)
	}
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "login":
		err = cmd.ExecuteLogin(os.Args[2:], confDir, confFile)
	case "db":
		err = cmd.ExecuteDB(os.Args[2:], confFile)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(1)
	}
	var e *pkg.ParseError
	if errors.As(err, &e) {
		fmt.Println(err)
		usage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}
}
