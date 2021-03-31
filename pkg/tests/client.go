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

import "github.com/rsds143/astra-devops-sdk-go/astraops"

// MockClient is used for testing
type MockClient struct {
	ErrorQueue []error
	calls      []interface{}
	Databases  []astraops.Database
	Tiers      []astraops.TierInfo
	Bundle     astraops.SecureBundle
}

// getError pops the next error stored off the stack
func (c *MockClient) getError() error {
	var err error
	if len(c.ErrorQueue) > 0 {
		err = c.ErrorQueue[0]
		c.ErrorQueue[0] = nil
		c.ErrorQueue = c.ErrorQueue[1:]
	}
	return err
}

// getError pops the next db object stored off the stack
func (c *MockClient) getDb() astraops.Database {
	var db astraops.Database
	if len(c.Databases) > 0 {
		db = c.Databases[0]
		c.Databases = c.Databases[1:]
	}
	return db
}

// Call returns a call at the specified index
func (c *MockClient) Call(index int) interface{} {
	return c.calls[index]
}

// Calls returns all calls made in order
func (c *MockClient) Calls() []interface{} {
	return c.calls
}

// CreateDb returns the next error and the next db created
func (c *MockClient) CreateDb(db astraops.CreateDb) (astraops.Database, error) {
	c.calls = append(c.calls, db)
	return c.getDb(), c.getError()
}

// Terminate returns the next error and stores the id used, internal is ignored
func (c *MockClient) Terminate(id string, internal bool) error {
	c.calls = append(c.calls, id)
	return c.getError()
}

// FindDb returns the next database and next error, the id call is stored
func (c *MockClient) FindDb(id string) (astraops.Database, error) {
	c.calls = append(c.calls, id)
	return c.getDb(), c.getError()
}

// ListDb returns all databases and stores the arguments as an interface array
func (c *MockClient) ListDb(include string, provider string, startingAfter string, limit int32) ([]astraops.Database, error) {
	c.calls = append(c.calls, []interface{}{
		include,
		provider,
		startingAfter,
		limit,
	})
	return c.Databases, c.getError()
}

// Unpark returns the next error, the id call is stored
func (c *MockClient) Unpark(id string) error {
	c.calls = append(c.calls, id)
	return c.getError()
}

// Park returns the next error, the id call is stored
func (c *MockClient) Park(id string) error {
	c.calls = append(c.calls, id)
	return c.getError()
}

// Resize returns the next error, the id call and size is stored
func (c *MockClient) Resize(id string, size int32) error {
	c.calls = append(c.calls, []interface{}{id, size})
	return c.getError()
}

// GetSecureBundle returns the next error, the secured bundle stored, and the id call is stored
func (c *MockClient) GetSecureBundle(id string) (astraops.SecureBundle, error) {
	c.calls = append(c.calls, id)
	return c.Bundle, c.getError()
}

// GetTierInfo returns the next error, and the tierinfo objects stored
func (c *MockClient) GetTierInfo() ([]astraops.TierInfo, error) {
	return c.Tiers, c.getError()
}
