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

// Package cmd contains all fo the commands for the cli
package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

func TestDBUsageFails(t *testing.T) {
	fails := func() error {
		return errors.New("error showing usage")
	}
	err := executeDB(fails)
	if err == nil {
		t.Fatal("there is supposed to be an error")
	}
	expected := "warn unable to show usage error showing usage"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestDBUsage(t *testing.T) {
	fails := func() error {
		return nil
	}
	err := executeDB(fails)
	if err != nil {
		t.Fatalf("unexpected eror %v", err)
	}
}

func TestDBShowHelp(t *testing.T) {
	clientJSON = ""
	authToken = ""
	clientName = ""
	clientSecret = ""
	clientID = ""
	originalOut := RootCmd.OutOrStderr()
	defer func() {
		RootCmd.SetOut(originalOut)
		RootCmd.SetArgs([]string{})
	}()
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"db"})
	err := RootCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error '%v'", err)
	}
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	expected := dbCmd.UsageString()

	if string(out) != expected {
		t.Errorf("expected\n'%q'\nbut was\n'%q'", expected, string(out))
	}
}
