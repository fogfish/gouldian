/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

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

package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	µ "github.com/fogfish/gouldian"
)

// Mock is an option type to customize mock event
type Mock func(*µ.Context) *µ.Context

// Endpoint mock Route
func Endpoint(route µ.Routable) µ.Endpoint {
	return µ.NewRoutes(route).Endpoint()
}

// Input mocks HTTP request, takes mock options to customize event
func Input(spec ...Mock) *µ.Context {
	input := µ.NewContext(context.Background())

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	input.Request = req

	for _, f := range spec {
		input = f(input)
	}
	return input
}

// Method changes the verb of mocked HTTP request
func Method(verb string) Mock {
	return func(mock *µ.Context) *µ.Context {
		mock.Request.Method = verb
		return mock
	}
}

// URL changes URL of mocked HTTP request
func URL(httpURL string) Mock {
	uri, err := url.Parse(httpURL)
	if err != nil {
		panic(err)
	}

	return func(mock *µ.Context) *µ.Context {
		mock.Request.URL = uri
		return mock
	}
}

// Header adds Header to mocked HTTP request
func Header(header string, value string) Mock {
	return func(mock *µ.Context) *µ.Context {
		mock.Request.Header.Set(header, value)
		return mock
	}
}

// JSON adds payload to mocked HTTP request
func JSON(val interface{}) Mock {
	return func(mock *µ.Context) *µ.Context {
		body, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		mock.Request.Header.Set("Content-Type", "application/json")
		mock.Request.Body = io.NopCloser(bytes.NewReader(body))
		return mock
	}
}

// Text adds payload to mocked HTTP request
func Text(val string) Mock {
	return func(mock *µ.Context) *µ.Context {
		mock.Request.Body = io.NopCloser(strings.NewReader(val))
		return mock
	}
}

// JWT adds JWT token to mocked HTTP request
func JWT(token µ.JWT) Mock {
	return func(mock *µ.Context) *µ.Context {
		mock.JWT = token
		return mock
	}
}
