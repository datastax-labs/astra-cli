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
	"flag"
	"fmt"
	"github.com/rsds143/astra-cli/cmd"
	"github.com/rsds143/astra-cli/pkg"
	"os"
)

var verbose = flag.Bool("v", false, "turns on verbose logging")

func usage() {
	flag.Usage()
	fmt.Println("commands:")
	fmt.Println(cmd.LoginUsage())
	fmt.Println(cmd.DBUsage())
}
func main() {
	flag.Parse()
	confDir, confFiles, err := pkg.GetHome()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(3)
	}
	if flag.NArg() == 1 {
		usage()
		os.Exit(1)
	}
	switch flag.Arg(0) {
	case "login":
		err = cmd.ExecuteLogin(flag.Args()[1:], confDir, confFiles)
	case "db":
		err = cmd.ExecuteDB(flag.Args()[1:], confFiles, *verbose)
	default:
		fmt.Printf("%q is not valid command.\n", flag.Arg(1))
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
