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
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const connections = 10
const standardTimeOut = 5 * time.Second
const dialTimeout = 10 * time.Second
const expectContinueResponse = 1 * time.Second

// NewHTTPClient fires up client with 'better' defaults
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: standardTimeOut,
		Transport: &http.Transport{
			MaxIdleConns:        connections,
			MaxConnsPerHost:     connections,
			MaxIdleConnsPerHost: connections,
			Dial: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: dialTimeout,
			}).Dial,
			TLSHandshakeTimeout:   standardTimeOut,
			ResponseHeaderTimeout: standardTimeOut,
			ExpectContinueTimeout: expectContinueResponse,
		},
	}
}

// DownloadZip pulls down the URL listed and saves it to the specified location
func DownloadZip(downloadURL string, secBundleLoc string) (int64, error) {
	httpClient := NewHTTPClient()
	res, err := httpClient.Get(downloadURL)
	if err != nil {
		return -1, fmt.Errorf("unable to download zip with error %v", err)
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warn: error closing http response body %v\n for request %v with status code %v\n", err, downloadURL, res.StatusCode)
		}
	}()
	f, err := os.Create(secBundleLoc)
	if err != nil {
		return -1, fmt.Errorf("unable to create file to save too %v", err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warn: error closing file %v for file %v\n", err, secBundleLoc)
		}
	}()
	i, err := io.Copy(f, res.Body)
	if err != nil {
		return -1, fmt.Errorf("unable to copy downloaded file to %v", err)
	}
	return i, nil
}
