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

// Package db is where the Astra DB commands are
package db

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
	astraops "github.com/datastax/astra-client-go/v2/astra"
	tests "github.com/rsds143/astra-cli/pkg/tests"
)

func TestList(t *testing.T) {
	listFmt = pkg.JSONFormat
	dbs := []astraops.Database{
		{Id: "1"},
		{Id: "2"},
	}
	jsonTxt, err := executeList(func() (pkg.Client, error) {
		return &tests.MockClient{
			Databases: dbs,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var fromServer []astraops.Database
	err = json.Unmarshal([]byte(jsonTxt), &fromServer)
	if err != nil {
		t.Fatalf("unexpected error with json %v with text %v", err, jsonTxt)
	}
	if len(fromServer) != len(dbs) {
		t.Errorf("expected '%v' but was '%v'", len(dbs), len(fromServer))
	}
	if fromServer[0].Id != dbs[0].Id {
		t.Errorf("expected '%v' but was '%v'", dbs[0].Id, fromServer[0].Id)
	}
	if fromServer[1].Id != dbs[1].Id {
		t.Errorf("expected '%v' but was '%v'", dbs[1].Id, fromServer[1].Id)
	}
}

func TestListText(t *testing.T) {
	listFmt = pkg.TextFormat
	name1 := "A"
	name2 := "B"
	dbs := []astraops.Database{
		{
			Id: "1",
			Info: astraops.DatabaseInfo{
				Name: &name1,
			},
			Status: astraops.StatusEnum_ACTIVE,
		},
		{
			Id: "2",
			Info: astraops.DatabaseInfo{
				Name: &name2,
			},
			Status: astraops.StatusEnum_TERMINATING,
		},
	}
	txt, err := executeList(func() (pkg.Client, error) {
		return &tests.MockClient{
			Databases: dbs,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := strings.Join([]string{
		"name id status",
		"A    1  ACTIVE",
		"B    2  TERMINATING",
	},
		"\n")
	if txt != expected {
		t.Errorf("expected '%v' but was '%v'", expected, txt)
	}
}

func TestListInvalidFmt(t *testing.T) {
	listFmt = "listham"
	_, err := executeList(func() (pkg.Client, error) {
		return &tests.MockClient{}, nil
	})
	if err == nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := "-o \"listham\" is not valid option"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestListFails(t *testing.T) {
	getFmt = pkg.JSONFormat
	dbs := []astraops.Database{}
	jsonTxt, err := executeList(func() (pkg.Client, error) {
		return &tests.MockClient{
			Databases:  dbs,
			ErrorQueue: []error{errors.New("cant find db")},
		}, nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	expected := "unable to get list of dbs with error 'cant find db'"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if jsonTxt != "" {
		t.Errorf("expected '%v' but was '%v'", "", jsonTxt)
	}
}

func TestListFailedLogin(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{errors.New("no db")}
	msg, err := executeList(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to get list of dbs with error 'no db'"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
	expected := ""
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}
