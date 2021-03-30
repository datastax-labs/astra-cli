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

//Package db is where the Astra DB commands are
package db

import (
	"testing"

	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

//MockClient is used for testing
type MockClient struct {
	errorStack []error
	calls      []interface{}
	databases  []astraops.Database
	tiers      []astraops.TierInfo
	bundle     astraops.SecureBundle
}

func (c *MockClient) getError() error {
	var err error
	if len(c.errorStack) > 0 {
		n := len(c.errorStack) - 1
		err = c.errorStack[n]
		c.errorStack = c.errorStack[:n]
	}
	return err
}

func (c *MockClient) getDb() astraops.Database {
	var db astraops.Database
	if len(c.databases) > 0 {
		n := len(c.databases) - 1
		db = c.databases[n]
		c.databases = c.databases[:n]
	}
	return db
}

func (c *MockClient) CreateDb(db astraops.CreateDb) (astraops.Database, error) {
	c.calls = append(c.calls, db)
	return c.getDb(), c.getError()
}

func (c *MockClient) Terminate(id string, internal bool) error {
	c.calls = append(c.calls, id)
	return c.getError()
}

func (c *MockClient) FindDb(id string) (astraops.Database, error) {
	c.calls = append(c.calls, id)
	return c.getDb(), c.getError()
}

func (c *MockClient) ListDb(include string, provider string, startingAfter string, limit int32) ([]astraops.Database, error) {

	c.calls = append(c.calls, []interface{}{
		include,
		provider,
		startingAfter,
		limit,
	})
	return c.databases, c.getError()
}

func (c *MockClient) Unpark(id string) error {
	c.calls = append(c.calls, id)
	return c.getError()
}
func (c *MockClient) Park(id string) error {
	c.calls = append(c.calls, id)
	return c.getError()
}

func (c *MockClient) Resize(id string, size int32) error {
	c.calls = append(c.calls, []interface{}{id, size})
	return c.getError()
}

func (c *MockClient) GetSecureBundle(id string) (astraops.SecureBundle, error) {
	c.calls = append(c.calls, id)
	return c.bundle, c.getError()
}

func (c *MockClient) GetTierInfo() ([]astraops.TierInfo, error) {
	return c.tiers, c.getError()
}

func TestCreateGetsId(t *testing.T) {
	expectedID := "abcd"
	//setting package variables by hand, there be dragons
	mockClient := &MockClient{
		databases: []astraops.Database{
			astraops.Database{
				ID: expectedID,
			},
		},
	}
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}

	if len(mockClient.calls) != 1 {
		t.Fatalf("expected 1 call but was %v", len(mockClient.calls))
	}
}

func TestCreateSetsName(t *testing.T) {
	mockClient := &MockClient{}
	createDbName = "mydb"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.Name != createDbName {
		t.Errorf("expected '%v' but was '%v'", arg0.Name, createDbName)
	}
}

func TestCreateSetsKeyspace(t *testing.T) {
	mockClient := &MockClient{}
	createDbKeyspace = "myKeyspace"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.Keyspace != createDbKeyspace {
		t.Errorf("expected '%v' but was '%v'", arg0.Keyspace, createDbKeyspace)
	}
}

func TestCreateSetsCapacityUnit(t *testing.T) {
	mockClient := &MockClient{}
	createDbCapacityUnit = 10000
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.CapacityUnits != int32(createDbCapacityUnit) {
		t.Errorf("expected '%v' but was '%v'", arg0.CapacityUnits, createDbCapacityUnit)
	}
}

func TestCreateSetsRegion(t *testing.T) {
	mockClient := &MockClient{}
	createDbRegion = "EU-West1"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.Region != createDbRegion {
		t.Errorf("expected '%v' but was '%v'", arg0.Region, createDbRegion)
	}
}

func TestCreateSetsUser(t *testing.T) {
	mockClient := &MockClient{}
	createDbUser = "john@james.com"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.User != createDbUser {
		t.Errorf("expected '%v' but was '%v'", arg0.User, createDbUser)
	}
}

func TestCreateSetsPass(t *testing.T) {
	mockClient := &MockClient{}
	createDbUser = "afdfdf"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.Password != createDbPassword {
		t.Errorf("expected '%v' but was '%v'", arg0.Password, createDbPassword)
	}
}

func TestCreateSetsTier(t *testing.T) {
	mockClient := &MockClient{}
	createDbTier = "afdfdf"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.Tier != createDbTier {
		t.Errorf("expected '%v' but was '%v'", arg0.Tier, createDbTier)
	}
}

func TestCreateSetsProvider(t *testing.T) {
	mockClient := &MockClient{}
	createDbCloudProvider = "ryanscloud"
	err := executeCreate(mockClient)
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
	arg0 := mockClient.calls[0].(astraops.CreateDb)
	if arg0.CloudProvider != createDbCloudProvider {
		t.Errorf("expected '%v' but was '%v'", arg0.CloudProvider, createDbCloudProvider)
	}
}
