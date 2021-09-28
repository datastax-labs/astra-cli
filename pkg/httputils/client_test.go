//   Copyright 2021 Ryan Svihla
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

// Package httputils provides common http functions and utilities
package httputils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestDownloadUrl(t *testing.T) {
	zipContent := "zip file content"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, zipContent)
	}))
	defer ts.Close()
	tmpDir := t.TempDir()
	zipFile := path.Join(tmpDir, "bundle.zip")
	bytesWritten, err := DownloadZip(ts.URL, zipFile)
	if bytesWritten == 0 {
		t.Fatal("Expected bytes to be written but none were")
	}

	if err != nil {
		t.Fatalf("Unexpected error test '%v'", err)
	}

	b, err := os.ReadFile(zipFile)
	if err != nil {
		t.Fatalf("Unexpected error reading file '%v'", err)
	}
	expectedZipContent := fmt.Sprintf("%v\n", zipContent)
	if expectedZipContent != string(b) {
		t.Errorf("expected/actual \n'%q'\n'%q'", expectedZipContent, string(b))
	}
}
