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

// Package db is where the Astra DB commands are
package db

import (
	"errors"
	"testing"

	"github.com/datastax-labs/astra-cli/pkg"
	tests "github.com/datastax-labs/astra-cli/pkg/tests"
)

func TestPark(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	id := "abcd"
	msg, err := executePark([]string{id}, func() (pkg.Client, error) {
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
	expected := "database abcd parked"
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestParkFailedLogin(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	id := "abcd"
	msg, err := executePark([]string{id}, func() (pkg.Client, error) {
		return mockClient, errors.New("bad login")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to login with error bad login"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
	expected := ""
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestParkFailed(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{errors.New("unable to park")}
	id := "123"
	msg, err := executePark([]string{id}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to park '123' with error unable to park"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
	expected := ""
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}
