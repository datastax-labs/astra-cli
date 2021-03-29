/**
   Copyright 2021 Ryan Svihla

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/
package pkg

import (
	"fmt"
	"path"
	"testing"
)

func TestUnableToReadHomeDir(t *testing.T) {
	noPath := func() (string, error) { return "", fmt.Errorf("unexpected error") }
	creds := &Creds{
		GetHomeFunc: noPath,
	}
	_, err := creds.Login()
	if err == nil {
		t.Fatal("expected an error on an empty path")
	}
	expected := "unable to read conf dir with error 'unable to get user home directory with error 'unexpected error''"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}
func TestMissingConfigFolder(t *testing.T) {
	noPath := func() (string, error) { return "", nil }
	creds := &Creds{
		GetHomeFunc: noPath,
	}
	_, err := creds.Login()
	if err == nil {
		t.Fatal("expected an error on an empty path")
	}
	expected := "unable to access any file for directory `.config/astra`, run astra-cli login first"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestLoginWithInvalidTokenFile(t *testing.T) {
	invalid := func() (string, error) { return path.Join("testdata", "with_invalid_token"), nil }
	creds := &Creds{
		GetHomeFunc: invalid,
	}
	_, err := creds.Login()
	if err == nil {
		t.Fatal("expected an error on an empty path")
	}
	expected := "found token at 'testdata/with_invalid_token/.config/astra/token' but unable to read token with error 'missing prefix 'AstraCS' in token file 'testdata/with_invalid_token/.config/astra/token''"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestLoginWithEmptyTokenFile(t *testing.T) {
	invalid := func() (string, error) { return path.Join("testdata", "with_empty_token"), nil }
	creds := &Creds{
		GetHomeFunc: invalid,
	}
	_, err := creds.Login()
	if err == nil {
		t.Fatal("expected an error on an empty path")
	}
	expected := "found token at 'testdata/with_empty_token/.config/astra/token' but unable to read token with error 'token file 'testdata/with_empty_token/.config/astra/token' is empty'"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestLoginValidToken(t *testing.T) {
	valid := func() (string, error) { return path.Join("testdata", "with_token"), nil }
	creds := &Creds{
		GetHomeFunc: valid,
	}
	_, err := creds.Login()
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}
}

func TestLoginWithInvalidSA(t *testing.T) {
	invalid := func() (string, error) { return path.Join("testdata", "with_invalid_sa"), nil }
	creds := &Creds{
		GetHomeFunc: invalid,
	}
	_, err := creds.Login()
	if err == nil {
		t.Fatal("expected an error on an empty path")
	}
	expected := "Invalid service account: Client ID for service account is empty for file 'testdata/with_invalid_sa/.config/astra/sa.json'"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}
