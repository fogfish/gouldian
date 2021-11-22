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
	"net/textproto"
	"net/url"
	"strings"

	µ "github.com/fogfish/gouldian"
)

// Mock is an option type to customize mock event
type Mock func(*µ.Input) *µ.Input

// Input mocks HTTP request, takes mock options to customize event
func Input(spec ...Mock) *µ.Input {
	input := &µ.Input{
		Context:  µ.NewContext(context.Background()),
		Method:   "GET",
		Resource: µ.Segments{},
		Params:   µ.Params{},
		Headers:  µ.Headers{},
	}

	for _, f := range spec {
		input = f(input)
	}
	return input
}

// Method changes the verb of mocked HTTP request
func Method(verb string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.Method = verb
		return mock
	}
}

// URL changes URL of mocked HTTP request
func URL(httpURL string) Mock {
	return func(mock *µ.Input) *µ.Input {
		uri, err := url.Parse(httpURL)
		if err != nil {
			panic(err)
		}

		segments := strings.Split(uri.Path, "/")[1:]
		if len(segments) == 1 && segments[0] == "" {
			segments = []string{}
		}
		mock.Resource = segments

		params := µ.Params{}
		for key, val := range uri.Query() {
			params[key] = []string{strings.Join(val, "")}
		}
		mock.Params = params

		return mock
	}
}

// Param add raw param string to mocked HTTP request
func Param(key, val string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.Params[key] = []string{val}
		return mock
	}
}

// Header adds Header to mocked HTTP request
func Header(header string, value string) Mock {
	return func(mock *µ.Input) *µ.Input {
		head := textproto.CanonicalMIMEHeaderKey(header)
		mock.Headers[head] = []string{value}
		return mock
	}
}

// JSON adds payload to mocked HTTP request
func JSON(val interface{}) Mock {
	return func(mock *µ.Input) *µ.Input {
		body, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		mock.Headers["Content-Type"] = []string{"application/json"}
		mock.Stream = bytes.NewReader(body)
		return mock
	}
}

// Text adds payload to mocked HTTP request
func Text(val string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.Stream = strings.NewReader(val)
		return mock
	}
}

// JWT adds JWT token to mocked HTTP request
func JWT(token µ.JWT) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.JWT = token
		return mock
	}
}
