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

	astraops "github.com/datastax/astra-client-go/v2/astra"
	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
)

func TestCreateGetsId(t *testing.T) {
	expectedID := "createID1234"
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{
		Databases: []astraops.Database{
			{
				Id: expectedID,
			},
		},
	}
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}

	if len(mockClient.Calls()) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.Calls()))
	}
}
func TestCreateLoginFails(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, fmt.Errorf("service down")
	})
	if err == nil {
		t.Fatal("expected error")
	}

	expected := "unable to login with error service down"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if len(mockClient.Calls()) != 0 {
		t.Fatalf("expected 0 call but was %v", len(mockClient.Calls()))
	}
}

func TestCreateFails(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{
		ErrorQueue: []error{fmt.Errorf("service down")},
	}
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatal("expected error")
	}

	if len(mockClient.Calls()) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.Calls()))
	}
}

func TestCreateSetsName(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbName = "mydb"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.Name != createDbName {
		t.Errorf("expected '%v' but was '%v'", arg0.Name, createDbName)
	}
}

func TestCreateSetsKeyspace(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbKeyspace = "myKeyspace"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.Keyspace != createDbKeyspace {
		t.Errorf("expected '%v' but was '%v'", arg0.Keyspace, createDbKeyspace)
	}
}

func TestCreateSetsCapacityUnit(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbCapacityUnit = 10000
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.CapacityUnits != createDbCapacityUnit {
		t.Errorf("expected '%v' but was '%v'", arg0.CapacityUnits, createDbCapacityUnit)
	}
}

func TestCreateSetsRegion(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbRegion = "EU-West1"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.Region != createDbRegion {
		t.Errorf("expected '%v' but was '%v'", arg0.Region, createDbRegion)
	}
}

func TestCreateSetsUser(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbUser = "john@james.com"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.User != createDbUser {
		t.Errorf("expected '%v' but was '%v'", arg0.User, createDbUser)
	}
}

func TestCreateSetsPass(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbUser = "afdfdf"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.Password != createDbPassword {
		t.Errorf("expected '%v' but was '%v'", arg0.Password, createDbPassword)
	}
}

func TestCreateSetsTier(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbTier = "afdfdf"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.Tier != astraops.DatabaseInfoCreateTier(createDbTier) {
		t.Errorf("expected '%v' but was '%v'", arg0.Tier, createDbTier)
	}
}

func TestCreateSetsProvider(t *testing.T) {
	mockClient := &tests.MockClient{}
	createDbCloudProvider = "ryanscloud"
	err := executeCreate(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.Call(0).(astraops.DatabaseInfoCreate)
	if arg0.CloudProvider != astraops.DatabaseInfoCreateCloudProvider(createDbCloudProvider) {
		t.Errorf("expected '%v' but was '%v'", arg0.CloudProvider, createDbCloudProvider)
	}
}
