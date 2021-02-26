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
	"errors"
	"testing"
)

func TestReadLogin(t *testing.T) {
	c, err := ReadLogin("testdata/sa.json")
	if err != nil {
		t.Fatal(err)
	}
	name := "me@example.com"
	if c.ClientName != name {
		t.Errorf("expected %v but was %v", name, c.ClientName)
	}
	id := "deeb55bd-2a55-4988-a345-d8fdddd0e0c9"
	if c.ClientID != id {
		t.Errorf("expected %v but was %v", id, c.ClientID)
	}
	secret := "6ae15bff-1435-430f-975b-9b3d9914b698"
	if c.ClientSecret != secret {
		t.Errorf("expected %v but was %v", secret, c.ClientSecret)
	}
}

func TestReadLoginWithNoFile(t *testing.T) {
	_, err := ReadLogin("testdata/not-a-real-file.json")
	if err == nil {
		t.Fatal("expected an error but there was none")
	}
	var e *FileNotFoundError
	if !errors.As(err, &e) {
		t.Errorf("expected %T but was %T", e, err)
	}
}

func TestReadLoginWithEmptyFile(t *testing.T) {
	_, err := ReadLogin("testdata/empty.json")
	if err == nil {
		t.Fatal("expected an error but there was none")
	}
	var e *JSONParseError
	if !errors.As(err, &e) {
		t.Errorf("expected %T but was %T", e, err)
	}
}
