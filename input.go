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

/*

Pattern is a union type of allowed params to matcher functions
*/
type Pattern interface{ string | Lens }

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

Token is a container for access token
*/
type Token map[string]string

// Jti is unique JWT token identity
func (t Token) Jti() string { return t["jti"] }

// Iss -uer of token
func (t Token) Iss() string { return t["iss"] }

// Exp -ires after
func (t Token) Exp() string { return t["exp"] }

// Sub -ject of token
func (t Token) Sub() string { return t["sub"] }

// Scope of the token
func (t Token) Scope() string { return t["scope"] }

// Username associated with token
func (t Token) Username() string { return t["username"] }

// ClientID associated with token
func (t Token) ClientID() string { return t["client_id"] }

/*

NewToken creates access token object
*/
func NewToken(raw map[string]interface{}) Token {
	asString := func(id string) string {
		if val, ok := raw[id]; ok {
			return val.(string)
		}
		return ""
	}

	return Token{
		"jti":       asString("jti"),
		"iss":       asString("iss"),
		"exp":       asString("exp"),
		"sub":       asString("sub"),
		"scope":     asString("scope"),
		"username":  asString("username"),
		"client_id": asString("client_id"),
	}
}
