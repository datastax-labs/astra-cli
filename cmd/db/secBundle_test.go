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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
	tests "github.com/rsds143/astra-cli/pkg/tests"
	"github.com/rsds143/astra-devops-sdk-go/astraops"
)

func TestSecBundle(t *testing.T) {
	id := "abc"
	secBundleLoc = "my_loc"
	secBundleFmt = "json"
	bundle := astraops.SecureBundle{
		DownloadURL:                       "abcd",
		DownloadURLInternal:               "wyz",
		DownloadURLMigrationProxy:         "opu",
		DownloadURLMigrationProxyInternal: "zert",
	}
	jsonTxt, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Bundle: bundle,
		}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var fromServer astraops.SecureBundle
	err = json.Unmarshal([]byte(jsonTxt), &fromServer)
	if err != nil {
		t.Fatalf("unexpected error with json %v", err)
	}
	if fromServer != bundle {
		t.Errorf("expected '%v' but was '%v'", bundle, fromServer)
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
	bundle := astraops.SecureBundle{
		DownloadURL:                       ts.URL,
		DownloadURLInternal:               "wyz",
		DownloadURLMigrationProxy:         "opu",
		DownloadURLMigrationProxyInternal: "zert",
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
	bundle := astraops.SecureBundle{
		DownloadURL:                       "url",
		DownloadURLInternal:               "wyz",
		DownloadURLMigrationProxy:         "opu",
		DownloadURLMigrationProxyInternal: "zert",
	}
	_, err := executeSecBundle([]string{id}, func() (pkg.Client, error) {
		return &tests.MockClient{
			Bundle: bundle,
		}, nil
	})
	if err == nil {
		t.Fatalf("unexpected error %v", err)
	}
	expected := "-output \"ham\" is not valid option"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}
