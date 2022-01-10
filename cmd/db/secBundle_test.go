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
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	astraops "github.com/datastax/astra-client-go/v2/astra"
	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
)

func TestSecBundle(t *testing.T) {
	id := "secId123"
	secBundleLoc = "my_loc"
	secBundleFmt = "json"
	bundle := astraops.CredsURL{
		DownloadURL:                       "abcd",
		DownloadURLInternal:               astraops.StringPtr("wyz"),
		DownloadURLMigrationProxy:         astraops.StringPtr("opu"),
		DownloadURLMigrationProxyInternal: astraops.StringPtr("zert"),
	}
	jsonTxt, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Bundle: bundle,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	// after we went to the newer api with it's heavy use of pointers we lost easy comparison, here I convert
	// the struct into json text for comparison
	bundleTxt, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		t.Fatalf("unexpected error with json %v", err)
	}
	if string(bundleTxt) != jsonTxt {
		t.Errorf("expected '%v' but was '%v", string(bundleTxt), jsonTxt)
	}
}

func TestSecBundleZip(t *testing.T) {
	zipContent := "zip file content"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, zipContent)
	}))
	defer ts.Close()
	tmpDir := t.TempDir()
	zipFile := path.Join(tmpDir, "bundle.zip")
	defer func() {
		if err := os.Remove(zipFile); err != nil {
			t.Logf("unable to remove '%v' in test due to error '%v'", zipFile, err)
		}
	}()
	id := "abc"
	secBundleLoc = zipFile
	secBundleFmt = "zip"
	bundle := astraops.CredsURL{
		DownloadURL:                       ts.URL,
		DownloadURLInternal:               astraops.StringPtr("wyz"),
		DownloadURLMigrationProxy:         astraops.StringPtr("opu"),
		DownloadURLMigrationProxyInternal: astraops.StringPtr("zert"),
	}
	msg, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Bundle: bundle,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := fmt.Sprintf("file %v saved 17 bytes written", zipFile)
	if msg != expected {
		t.Errorf("expected '%v' but was '%v'", expected, msg)
	}
}

func TestSecBundleInvalidFmt(t *testing.T) {
	id := "abc"
	secBundleFmt = "ham"
	bundle := astraops.CredsURL{
		DownloadURL:                       "url",
		DownloadURLInternal:               astraops.StringPtr("wyz"),
		DownloadURLMigrationProxy:         astraops.StringPtr("opu"),
		DownloadURLMigrationProxyInternal: astraops.StringPtr("zert"),
	}
	_, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Bundle: bundle,
		}, nil
	})
	if err == nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := "-o \"ham\" is not valid option"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestSecBundleFailed(t *testing.T) {
	// setting package variables by hand, there be dragons
	mockClient := &tests.MockClient{}
	mockClient.ErrorQueue = []error{errors.New("no db")}
	id := "12390"
	_, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return mockClient, nil
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "unable to get '12390' with error no db"
	if err.Error() != expectedErr {
		t.Errorf("expected '%v' but was '%v'", expectedErr, err)
	}
}
