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
	"errors"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
)

func TestResize(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	id := "resizeId1"
	size := "100"
	err := executeResize([]string{id, size}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}

	if len(mockClient.Calls()) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.Calls()))
	}
	actualID := mockClient.Call(0).([]interface{})[0]
	if id != actualID {
		t.Errorf("expected '%v' but was '%v'", id, actualID)
	}
	actualSize := mockClient.Call(0).([]interface{})[1].(int32)
	if int32(100) != actualSize {
		t.Errorf("expected '%v' but was '%v'", size, actualSize)
	}
}

func TestResizeParseError(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	id := "resizeparseId"
	size := "poppaoute"
	err := executeResize([]string{id, size}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	expectedError := "Unable to parse command line with args: resizeparseId, poppaoute. Nested error was 'unable to parse capacity unit 'poppaoute' with error strconv.ParseInt: parsing \"poppaoute\": invalid syntax'"
	if err.Error() != expectedError {
		t.Errorf("expected '%v' but was '%v'", expectedError, err.Error())
	}
	if len(mockClient.Calls()) != 0 {
		t.Fatalf("expected 0 call but was %v", len(mockClient.Calls()))
	}
}

func TestResizeFailed(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{errors.New("no db")}
	id := "12389"
	err := executeResize([]string{id, "100"}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to resize '12389' with error no db"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
}

func TestResizeFailedLogin(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{}
	id := "12390"
	err := executeResize([]string{id, "100"}, func() (pkg.Client, error) {
		return mockClient, errors.New("no db")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to login with error no db"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
}
