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
	tests "github.com/rsds143/astra-cli/pkg/tests"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

func TestGet(t *testing.T) {
	getFmt = pkg.JSONFormat
	dbs := []astraops.Database{
		{ID: "1"},
		{ID: "2"},
	}
	jsonTxt, err := executeGet([]string{"1"}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Databases: dbs,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var fromServer astraops.Database
	err = json.Unmarshal([]byte(jsonTxt), &fromServer)
	if err != nil {
		t.Fatalf("unexpected error with json %v with text %v", err, jsonTxt)
	}
	if fromServer.ID != dbs[0].ID {
		t.Errorf("expected '%v' but was '%v'", dbs[0].ID, fromServer.ID)
	}
}

func TestGetFindDbFails(t *testing.T) {
	getFmt = pkg.JSONFormat
	dbs := []astraops.Database{}
	jsonTxt, err := executeGet([]string{"1"}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Databases:  dbs,
			ErrorQueue: []error{errors.New("cant find db")},
		}, nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	expected := "unable to get '1' with error cant find db"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if jsonTxt != "" {
		t.Errorf("expected '%v' but was '%v'", "", jsonTxt)
	}
}

func TestGetFailedLogin(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{}
	id := "12345"
	msg, err := executeGet([]string{id}, func() (pkg.Client, error) {
		return mockClient, errors.New("no db")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := tests.LoginError
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
	expected := ""
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestGetText(t *testing.T) {
	getFmt = pkg.TextFormat
	dbs := []astraops.Database{
		{
			ID: "1",
			Info: astraops.DatabaseInfo{
				Name: "A",
			},
			Status: astraops.ACTIVE,
		},
		{
			ID: "2",
			Info: astraops.DatabaseInfo{
				Name: "B",
			},
			Status: astraops.TERMINATING,
		},
	}
	txt, err := executeGet([]string{"1"}, func() (pkg.Client, error) {
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
	},
		"\n")
	if txt != expected {
		t.Errorf("expected '%v' but was '%v'", expected, txt)
	}
}

func TestGetInvalidFmt(t *testing.T) {
	getFmt = "badham"
	_, err := executeGet([]string{"abc"}, func() (pkg.Client, error) {
		return &tests.MockClient{}, nil
	})
	if err == nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := "-o \"badham\" is not valid option"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}
