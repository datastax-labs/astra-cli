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

// Package cmd contains all fo the commands for the cli
package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/rsds143/astra-cli/pkg"
)

const testJSON = `{"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"me@example.com","clientSecret":"6ae15bff-1435-430f-975b-9b3d9914b698"}`
const testSecret = "jljlajef"
const testName = "me@me.com"
const testID = "abd278332"

func usageFunc() error {
	return nil
}
func TestLoginCmdJson(t *testing.T) {
	clientJSON = testJSON
	defer func() {
		clientJSON = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %v", exitCode)
	}
	fd, err := os.Open(f)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	b, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if string(b) == "" {
		t.Error("login file was empty")
	}
	if clientJSON != string(b) {
		t.Errorf("expected '%v' but was '%v'", clientJSON, string(b))
	}
}

func TestLoginCmdJsonInvalidPerms(t *testing.T) {
	clientJSON = testJSON
	authToken = ""
	clientID = ""
	clientName = ""
	clientSecret = ""
	defer func() {
		clientJSON = ""
	}()
	dir := t.TempDir()
	inaccessible := path.Join(dir, "inaccessible")
	err := os.Mkdir(inaccessible, 0400)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer os.RemoveAll(inaccessible)
	f := path.Join(inaccessible, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Error("expected error")
	}
	if exitCode != WriteError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
	expected := fmt.Sprintf("unable to create the login file due to error open %v: permission denied", f)
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestArgs(t *testing.T) {
	clientID = `deeb55bd-2a55-4988-a345-d8fdddd0e0c9`
	clientName = "me@example.com"
	clientSecret = "fortestargs"
	defer func() {
		clientID = ""
		clientSecret = ""
		clientName = ""
	}()

	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--clientId", clientID, "--clientName", clientName, "--clientSecret", clientSecret}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %v", exitCode)
	}
	fd, err := os.Open(pkg.PathWithEnv(f))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	b, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if string(b) == "" {
		t.Error("login file was empty")
	}
	expected := `{"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"me@example.com","clientSecret":"fortestargs"}`
	if expected != string(b) {
		t.Errorf("expected\n%v\nactual:\n%v", expected, string(b))
	}
}
func TestArgsWithNoPermission(t *testing.T) {
	clientJSON = ""
	authToken = ""
	clientID = `deeb55bd-2a55-4988-a345-d8fdddd0e0c9`
	clientName = "me@example.com"
	clientSecret = "fortestargsnoperm"
	defer func() {
		clientID = ""
		clientSecret = ""
		clientName = ""
	}()

	dir := t.TempDir()
	inaccessible := path.Join(dir, "inaccessible")
	err := os.Mkdir(inaccessible, 0400)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer os.RemoveAll(inaccessible)
	f := path.Join(inaccessible, "mytempFile")
	exitCode, err := executeLogin([]string{"--clientId", clientID, "--clientName", clientName, "--clientSecret", clientSecret}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Error("expected error")
	}
	if exitCode != WriteError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
	expected := fmt.Sprintf("unable to create the login file due to error open %v: permission denied", pkg.PathWithEnv(f))
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestLoginArgsMissingId(t *testing.T) {
	clientJSON = ""
	authToken = ""
	clientName = testName
	clientSecret = testSecret
	clientID = ""
	defer func() {
		clientSecret = ""
		clientName = ""
	}()
	exitCode, err := executeLogin([]string{"--clientId", clientID, "--clientName", clientName, "--clientSecret", clientSecret}, func() (string, pkg.ConfFiles, error) {
		return "", pkg.ConfFiles{}, nil
	}, usageFunc)
	if err == nil {
		t.Error("expected error")
	}
	expected := `Unable to parse command line with args: --clientId, , --clientName, me@me.com, --clientSecret, jljlajef. Nested error was 'clientId missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginArgsMissingName(t *testing.T) {
	clientJSON = ""
	authToken = ""
	clientName = ""
	clientSecret = testSecret
	clientID = testID
	defer func() {
		clientSecret = ""
		clientID = ""
	}()
	exitCode, err := executeLogin([]string{"--clientId", clientID, "--clientName", clientName, "--clientSecret", clientSecret}, func() (string, pkg.ConfFiles, error) {
		return "", pkg.ConfFiles{}, nil
	}, usageFunc)
	if err == nil {
		t.Errorf("expected error")
	}
	expected := `Unable to parse command line with args: --clientId, abd278332, --clientName, , --clientSecret, jljlajef. Nested error was 'clientName missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginArgsMissingSecret(t *testing.T) {
	clientJSON = ""
	authToken = ""
	clientName = testName
	clientSecret = ""
	clientID = testID
	defer func() {
		clientName = ""
		clientID = ""
	}()
	exitCode, err := executeLogin([]string{"--clientId", clientID, "--clientName", clientName, "--clientSecret", clientSecret}, func() (string, pkg.ConfFiles, error) {
		return "", pkg.ConfFiles{}, nil
	}, usageFunc)
	if err == nil {
		t.Error("expected error")
	}
	expected := `Unable to parse command line with args: --clientId, abd278332, --clientName, me@me.com, --clientSecret, . Nested error was 'clientSecret missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginHomeError(t *testing.T) {
	clientJSON = "invalidjson"
	defer func() { clientJSON = "" }()
	exitCode, err := executeLogin([]string{}, func() (string, pkg.ConfFiles, error) {
		return "", pkg.ConfFiles{}, fmt.Errorf("big error")
	}, usageFunc)
	if err == nil {
		t.Logf("expected error there was none and exit code was %v", exitCode)
		t.FailNow()
	}
	expected := "big error"
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != CannotFindHome {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginCmdJsonMissignId(t *testing.T) {
	clientJSON = `{"clientId":"","clientName":"me@example.com","clientSecret":"6ae15bff-1435-430f-975b-9b3d9914b698"}`
	defer func() {
		clientJSON = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Error("expected error")
	}
	expected := `Unable to parse command line with args: --json, {"clientId":"","clientName":"me@example.com","clientSecret":"6ae15bff-1435-430f-975b-9b3d9914b698"}. Nested error was 'clientId missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginCmdJsonMissignName(t *testing.T) {
	clientJSON = `{"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"","clientSecret":"6ae15bff-1435-430f-975b-9b3d9914b698"}`
	defer func() {
		clientJSON = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Errorf("expected error")
	}
	expected := `Unable to parse command line with args: --json, {"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"","clientSecret":"6ae15bff-1435-430f-975b-9b3d9914b698"}. Nested error was 'clientName missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginCmdJsonMissignSecret(t *testing.T) {
	clientJSON = `{"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"me@example.com","clientSecret":""}`
	defer func() {
		clientJSON = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Error("expected error")
	}
	expected := `Unable to parse command line with args: --json, {"clientId":"deeb55bd-2a55-4988-a345-d8fdddd0e0c9","clientName":"me@example.com","clientSecret":""}. Nested error was 'clientSecret missing'`
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginCmdJsonInvalid(t *testing.T) {
	clientJSON = `invalidtext`
	defer func() {
		clientJSON = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--json", clientJSON}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			SaPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Errorf("expected error")
	}
	if exitCode != JSONError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
}

func TestLoginToken(t *testing.T) {
	authToken = `6ae15bff-1435-430f-975b-9b3d9914b698`
	defer func() {
		authToken = ""
	}()
	dir := path.Join(t.TempDir(), "config")
	f := path.Join(dir, "mytempFile")
	exitCode, err := executeLogin([]string{"--token", authToken}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			TokenPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %v", exitCode)
	}
	fd, err := os.Open(pkg.PathWithEnv(f))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	b, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if string(b) == "" {
		t.Error("login file was empty")
	}
	if authToken != string(b) {
		t.Errorf("expected '%v' but was '%v'", authToken, string(b))
	}
}

func TestLoginTokenInvalidPerms(t *testing.T) {
	authToken = `6ae15bff-1435-430f-975b-9b3d9914b698`
	defer func() {
		authToken = ""
	}()
	dir := t.TempDir()
	inaccessible := path.Join(dir, "inaccessible")
	err := os.Mkdir(inaccessible, 0400)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer os.RemoveAll(inaccessible)
	f := path.Join(inaccessible, "mytempFile")
	exitCode, err := executeLogin([]string{"--token", authToken}, func() (string, pkg.ConfFiles, error) {
		return dir, pkg.ConfFiles{
			TokenPath: f,
		}, nil
	}, usageFunc)
	defer os.RemoveAll(dir)
	if err == nil {
		t.Error("expected error")
	}
	if exitCode != WriteError {
		t.Errorf("unexpected exit code %v", exitCode)
	}
	expected := fmt.Sprintf("unable to create the login file due to error open %v: permission denied", pkg.PathWithEnv(f))
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestMakeConf(t *testing.T) {
	content := "mycontent"
	dir := t.TempDir()
	f := path.Join(dir, "mytempFile")
	err := makeConf(dir, f, content)
	if err != nil {
		t.Fatalf("unable to create conf with error %v", err)
	}
	fd, err := os.Open(f)
	if err != nil {
		t.Fatalf("unable to read conf with error %v", err)
	}
	b, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("unable to read conf with error %v", err)
	}
	str := string(b)
	if str != content {
		t.Errorf("expected '%v' but was '%v'", content, str)
	}
}

func TestMakeConfWithInaccessibleDir(t *testing.T) {
	dir := t.TempDir()
	inaccessible := path.Join(dir, "inaccessible")
	err := os.Mkdir(inaccessible, 0400)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer os.RemoveAll(inaccessible)
	f := path.Join(inaccessible, "mytempFile")
	err = makeConf(inaccessible, f, "mycontent")
	if err == nil {
		t.Fatalf("expected error")
	}
	expected := fmt.Sprintf("unable to create the login file due to error open %v: permission denied", f)
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}

func TestMakeConfWithNonMakeableDir(t *testing.T) {
	dir := t.TempDir()
	inaccessible := path.Join(dir, "inaccessible")
	err := os.Mkdir(inaccessible, 0400)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer os.RemoveAll(inaccessible)
	newDir := path.Join(inaccessible, "new")
	f := path.Join(newDir, "mytempFile")
	err = makeConf(newDir, f, "mycontent")
	if err == nil {
		t.Fatalf("expected error")
	}
	expected := fmt.Sprintf("unable to get make config directory with error mkdir %v: permission denied", newDir)
	if err.Error() != expected {
		t.Errorf("expected '%v' but was '%v'", expected, err.Error())
	}
}
