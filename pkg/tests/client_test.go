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

// Package test is for test utilies and mocks
package test

import (
	"errors"
	"testing"
	astraops "github.com/datastax/astra-client-go/v2/astra"
)

func TestGetError(t *testing.T) {
	client := &MockClient{
		ErrorQueue: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	err := client.getError()
	if err.Error() != "error 1" {
		t.Errorf("expected 'error 1' but was '%v'", err.Error())
	}
	err = client.getError()
	if err.Error() != "error 2" {
		t.Errorf("expected 'error 2' but was '%v'", err.Error())
	}
	err = client.getError()
	if err.Error() != "error 3" {
		t.Errorf("expected 'error 3' but was '%v'", err.Error())
	}
	err = client.getError()
	if err != nil {
		t.Errorf("expected nil but was '%v'", err.Error())
	}
}

func TestGetDB(t *testing.T) {
	client := &MockClient{
		Databases: []astraops.Database{
			{Id: "1"},
			{Id: "2"},
			{Id: "3"},
		},
	}
	id := client.getDb().Id
	if id != "1" {
		t.Errorf("expected '1' but was '%v'", id)
	}
	id = client.getDb().Id
	if id != "2" {
		t.Errorf("expected '2' but was '%v'", id)
	}
	id = client.getDb().Id
	if id != "3" {
		t.Errorf("expected '3' but was '%v'", id)
	}
	id = client.getDb().Id
	if id != "" {
		t.Errorf("expected '' but was '%v'", id)
	}
}

func TestPark(t *testing.T) {
	client := &MockClient{}
	id := "123"
	err := client.Park(id)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if client.Call(0) != id {
		t.Errorf("expected '%v' but was '%v'", id, client.Call(0))
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestUnpark(t *testing.T) {
	client := &MockClient{}
	id := "parkid"
	err := client.Unpark(id)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if client.Call(0) != id {
		t.Errorf("expected '%v' but was '%v'", id, client.Call(0))
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestTerminate(t *testing.T) {
	client := &MockClient{}
	id := "termid"
	err := client.Terminate(id, false)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if client.Call(0) != id {
		t.Errorf("expected '%v' but was '%v'", id, client.Call(0))
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestGetSecurteBundleId(t *testing.T) {
	url := "myurl"
	client := &MockClient{
		Bundle: astraops.CredsURL{
			DownloadURL: url,
		},
	}
	id := "secid"
	bundle, err := client.GetSecureBundle(id)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if bundle.DownloadURL != url {
		t.Errorf("expected '%v' but was '%v'", url, bundle.DownloadURL)
	}
	if client.Call(0) != id {
		t.Errorf("expected '%v' but was '%v'", id, client.Call(0))
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestFindDb(t *testing.T) {
	id := "DSQ"

	client := &MockClient{
		Databases: []astraops.Database{
			{Id: id},
			{Id: "fakeid"},
		},
	}
	db, err := client.FindDb(id)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if db.Id != id {
		t.Errorf("expected '%v' but was '%v'", id, db.Id)
	}
	if client.Call(0) != id {
		t.Errorf("expected '%v' but was '%v'", id, client.Call(0))
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestCreateDb(t *testing.T) {
	id := "DSQ"

	client := &MockClient{
		Databases: []astraops.Database{
			{Id: id},
			{Id: "fakeid"},
		},
	}
	db, err := client.CreateDb(astraops.DatabaseInfoCreate{
		Name: "myname",
	})
	if err != nil {
		t.Fatal("unexpected error")
	}
	if db.Id != id {
		t.Errorf("expected '%v' but was '%v'", id, db.Id)
	}
	if client.Call(0).(astraops.DatabaseInfoCreate).Name != "myname" {
		t.Errorf("expected '%v' but was '%v'", "myname", client.Call(0).(astraops.DatabaseInfoCreate).Name)
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestResize(t *testing.T) {
	client := &MockClient{}
	id := "987"
	size := 10
	err := client.Resize(id, size)
	if err != nil {
		t.Fatal("unexpected error")
	}
	actual := client.Call(0).([]interface{})
	if actual[0].(string) != id {
		t.Errorf("expected '%v' but was '%v'", id, actual[0])
	}
	if actual[1].(int) != size {
		t.Errorf("expected '%v' but was '%v'", size, actual[1])
	}
	if len(client.Calls()) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(client.Calls()))
	}
}

func TestTiers(t *testing.T) {
	client := &MockClient{
		Tiers: []astraops.AvailableRegionCombination{
			{Tier: "abc"},
		},
	}
	tiers, err := client.GetTierInfo()
	if err != nil {
		t.Fatal("unexpected error")
	}

	if tiers[0].Tier != "abc" {
		t.Errorf("expected '%v' but was '%v'", "abc", tiers[0].Tier)
	}

	if len(client.Calls()) != 0 {
		t.Errorf("expected '%v' but was '%v'", 0, len(client.Calls()))
	}
}

func TestListdDb(t *testing.T) {
	id1 := "1"
	id2 := "2"
	include := "filter"
	provider := "gcp"
	starting := "today"
	limit := 1000
	client := &MockClient{
		Databases: []astraops.Database{
			{Id: id1},
			{Id: id2},
		},
	}
	dbs, err := client.ListDb(include, provider, starting, limit)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if len(dbs) != 2 {
		t.Errorf("expected '%v' but was '%v'", 2, len(dbs))
	}
	calls := client.Calls()
	if len(calls) != 1 {
		t.Errorf("expected '%v' but was '%v'", 1, len(calls))
	}
	args := calls[0].([]interface{})
	actualInclude := args[0].(string)
	if actualInclude != include {
		t.Errorf("expected '%v' but was '%v'", include, actualInclude)
	}
	actualProvider := args[1].(string)
	if actualProvider != provider {
		t.Errorf("expected '%v' but was '%v'", provider, actualProvider)
	}
	actualStarting := args[2].(string)
	if actualStarting != starting {
		t.Errorf("expected '%v' but was '%v'", starting, actualStarting)
	}
	actualLimit := args[3].(int)
	if actualLimit != limit {
		t.Errorf("expected '%v' but was '%v'", limit, actualLimit)
	}
}
