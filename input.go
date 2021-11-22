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

package gouldian

import (
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"
)

/*

Segments of path in HTTP request
*/
type Segments []string

/*

Params of path query in HTTP request
*/
type Params map[string][]string

// Get parameter by key
func (params Params) Get(key string) (string, bool) {
	v, exists := params[key]
	if !exists {
		return "", exists
	}

	return v[0], exists
}

/*

Headers of HTTP request
*/
type Headers map[string][]string

// Get header by its value
func (headers Headers) Get(key string) (string, bool) {
	header := textproto.CanonicalMIMEHeaderKey(key)
	v, exists := headers[header]
	if !exists {
		// Note: required due to browser behavior
		v, exists = headers[strings.ToLower(header)]
		if !exists {
			return "", exists
		}
		return v[0], exists
	}
	return v[0], exists
}

/*

Input is the HTTP request
*/
type Input struct {
	Context Context

	Method   string
	Resource Segments
	Params   Params
	Headers  Headers
	JWT      JWT
	Payload  string
	Stream   io.Reader
}

// ReadAll is helper function to consume Stream multiple types
func (in *Input) ReadAll() error {
	if in.Stream != nil {
		buf, err := ioutil.ReadAll(in.Stream)
		if err != nil {
			return err
		}
		// This is copied from runtime. It relies on the string
		// header being a prefix of the slice header!
		// in.Payload = *(*string)(unsafe.Pointer(&buf))
		in.Payload = string(buf)
		in.Stream = nil
		return nil
	}

	return nil
}

/*

JWT is a container for access token
*/
type JWT map[string]string

// Jti is unique JWT token identity
func (t JWT) Jti() string { return t["jti"] }

// Iss -uer of token
func (t JWT) Iss() string { return t["iss"] }

// Exp -ires after
func (t JWT) Exp() string { return t["exp"] }

// Sub -ject of token
func (t JWT) Sub() string { return t["sub"] }

// Scope of the token
func (t JWT) Scope() string { return t["scope"] }

// Username associated with token
func (t JWT) Username() string { return t["username"] }

// ClientID associated with token
func (t JWT) ClientID() string { return t["client_id"] }

/*

NewJWT creates access token object
*/
func NewJWT(raw map[string]interface{}) JWT {
	asString := func(id string) string {
		if val, ok := raw[id]; ok {
			return val.(string)
		}
		return ""
	}

	return JWT{
		"jti":       asString("jti"),
		"iss":       asString("iss"),
		"exp":       asString("exp"),
		"sub":       asString("sub"),
		"scope":     asString("scope"),
		"username":  asString("username"),
		"client_id": asString("client_id"),
	}
}
