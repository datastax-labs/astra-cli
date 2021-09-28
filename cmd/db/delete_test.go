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
	"fmt"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
)

func TestDelete(t *testing.T) {
	mockClient := &tests.MockClient{}
	id := "123"
	msg, err := executeDelete([]string{id}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	if len(mockClient.Calls()) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.Calls()))
	}
	if id != mockClient.Call(0) {
		t.Errorf("expected '%v' but was '%v'", id, mockClient.Call(0))
	}
	expected := "database 123 deleted"
	if expected != msg {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestDeleteLoginError(t *testing.T) {
	mockClient := &tests.MockClient{}
	id := "123"
	msg, err := executeDelete([]string{id}, func() (pkg.Client, error) {
		return mockClient, fmt.Errorf("unable to login")
	})
	if err == nil {
		t.Fatalf("should be returning an error and is not")
	}
	expectedError := "unable to login with error 'unable to login'"
	if err.Error() != expectedError {
		t.Errorf("expected '%v' but was '%v'", expectedError, err)
	}
	if len(mockClient.Calls()) != 0 {
		t.Fatalf("expected no calls but was %v", len(mockClient.Calls()))
	}

	if "" != msg {
		t.Errorf("expected empty but was '%v'", msg)
	}
}

func TestDeleteError(t *testing.T) {
	mockClient := &tests.MockClient{
		ErrorQueue: []error{fmt.Errorf("timeout error")},
	}
	id := "123"
	msg, err := executeDelete([]string{id}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatal("expected error but none came")
	}
	if len(mockClient.Calls()) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.Calls()))
	}
	if id != mockClient.Call(0) {
		t.Errorf("expected '%v' but was '%v'", id, mockClient.Call(0))
	}
	expected := ""
	if expected != msg {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}
