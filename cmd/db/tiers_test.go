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
	"reflect"
	"strings"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

func TestTiers(t *testing.T) {
	tiersFmt = "json"
	tier1 := astraops.TierInfo{
		Tier: "abd",
	}
	tier2 := astraops.TierInfo{
		Tier: "xyz",
	}
	jsonTxt, err := executeTiers(func() (pkg.Client, error) {
		return &tests.MockClient{
			Tiers: []astraops.TierInfo{
				tier1,
				tier2,
			},
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var fromServer []astraops.TierInfo
	err = json.Unmarshal([]byte(jsonTxt), &fromServer)
	if err != nil {
		t.Fatalf("unexpected error with json %v", err)
	}
	expected := []astraops.TierInfo{
		tier1,
		tier2,
	}
	if !reflect.DeepEqual(fromServer, expected) {
		t.Errorf("expected '%v' but was '%v'", expected, fromServer)
	}
}

func TestTiersText(t *testing.T) {
	tiersFmt = "text"
	tier1 := astraops.TierInfo{
		Tier:               "tier1",
		CloudProvider:      "cloud1",
		Region:             "region1",
		DatabaseCountUsed:  1,
		DatabaseCountLimit: 1,
		Cost: &astraops.Costs{
			CostPerMonthCents: 10,
			CostPerMinCents:   1,
		},
		CapacityUnitsUsed:  1,
		CapacityUnitsLimit: 1,
	}
	tier2 := astraops.TierInfo{
		Tier:          "tier2",
		CloudProvider: "cloud2",
		Region:        "region2",
		Cost: &astraops.Costs{
			CostPerMonthCents: 20,
			CostPerMinCents:   2,
		},
		CapacityUnitsUsed:  2,
		CapacityUnitsLimit: 2,
	}
	msg, err := executeTiers(func() (pkg.Client, error) {
		return &tests.MockClient{
			Tiers: []astraops.TierInfo{
				tier1,
				tier2,
			},
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := strings.Join([]string{
		"name  cloud  region  db (used)/(limit) cap (used)/(limit) cost per month cost per minute",
		"tier1 cloud1 region1 1/1               1/1                $0.10          $0.01",
		"tier2 cloud2 region2 0/0               2/2                $0.20          $0.02",
	}, "\n")
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestTiersTextWithNoCost(t *testing.T) {
	tiersFmt = "text"
	tier1 := astraops.TierInfo{
		Tier:               "tier1",
		CloudProvider:      "cloud1",
		Region:             "region1",
		DatabaseCountUsed:  1,
		DatabaseCountLimit: 1,
		CapacityUnitsUsed:  1,
		CapacityUnitsLimit: 1,
	}
	tier2 := astraops.TierInfo{
		Tier:               "tier2",
		CloudProvider:      "cloud2",
		Region:             "region2",
		CapacityUnitsUsed:  2,
		CapacityUnitsLimit: 2,
	}
	msg, err := executeTiers(func() (pkg.Client, error) {
		return &tests.MockClient{
			Tiers: []astraops.TierInfo{
				tier1,
				tier2,
			},
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := strings.Join([]string{
		"name  cloud  region  db (used)/(limit) cap (used)/(limit) cost per month cost per minute",
		"tier1 cloud1 region1 1/1               1/1                $0.00          $0.00",
		"tier2 cloud2 region2 0/0               2/2                $0.00          $0.00",
	}, "\n")
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestTiersnvalidFmt(t *testing.T) {
	tiersFmt = "ham"
	_, err := executeTiers(func() (pkg.Client, error) {
		return &tests.MockClient{}, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expected := "-o \"ham\" is not valid option"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestTiersFailedLogin(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{}
	_, err := executeTiers(func() (pkg.Client, error) {
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

func TestTiersFailed(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{errors.New("no db")}
	_, err := executeTiers(func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to get tiers with error no db"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
}
