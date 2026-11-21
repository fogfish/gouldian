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
	"encoding/json"
	"strconv"
	"strings"
)

// Pattern is a union type of allowed params to matcher functions
type Pattern interface{ string | Lens }

// Query of HTTP request
type Query map[string][]string

// Get returns the first value of the parameter
func (query Query) Get(key string) (string, bool) {
	v, exists := query[key]
	if !exists || len(v) == 0 {
		return "", false
	}

	return v[0], true
}

// GetAll returns all values of the (repeated) parameter
func (query Query) GetAll(key string) ([]string, bool) {
	v, exists := query[key]
	return v, exists && len(v) > 0
}

// Token is a container for claims of access token. The claims keep the
// type they have in the token (string, number, bool, array, object).
// Use typed accessors to read them.
type Token map[string]any

// Jti is unique JWT token identity
func (t Token) Jti() string { return t.String("jti") }

// Iss -uer of token
func (t Token) Iss() string { return t.String("iss") }

// Exp -ires after
func (t Token) Exp() string { return t.String("exp") }

// Sub -ject of token
func (t Token) Sub() string { return t.String("sub") }

// Scope of the token
func (t Token) Scope() string { return t.String("scope") }

// Username associated with token
func (t Token) Username() string { return t.String("username") }

// ClientID associated with token
func (t Token) ClientID() string { return t.String("client_id") }

// Get claim as is
func (t Token) Get(claim string) (any, bool) {
	v, ok := t[claim]
	return v, ok
}

// String returns claim as a string. Numbers and booleans are converted to
// their string form (no exponent), array of strings is joined with space
// (e.g. scopes). Objects and null are empty.
func (t Token) String(claim string) string {
	switch x := t[claim].(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case json.Number:
		return x.String()
	case bool:
		return strconv.FormatBool(x)
	case []string:
		return strings.Join(x, " ")
	case []any:
		seq := make([]string, 0, len(x))
		for _, v := range x {
			if s, ok := v.(string); ok {
				seq = append(seq, s)
			}
		}
		return strings.Join(seq, " ")
	default:
		return ""
	}
}

// Number returns numeric claim (e.g. exp, iat, nbf).
func (t Token) Number(claim string) (float64, bool) {
	switch x := t[claim].(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case json.Number:
		v, err := x.Float64()
		return v, err == nil
	case string:
		v, err := strconv.ParseFloat(x, 64)
		return v, err == nil
	default:
		return 0, false
	}
}

// Bool returns boolean claim
func (t Token) Bool(claim string) (bool, bool) {
	switch x := t[claim].(type) {
	case bool:
		return x, true
	case string:
		v, err := strconv.ParseBool(x)
		return v, err == nil
	default:
		return false, false
	}
}

// Strings returns claim as list of strings, a string is split by space.
func (t Token) Strings(claim string) []string {
	return strings.Fields(t.String(claim))
}

// NewToken creates access token object from claims. It never panics,
// the claims keep their types. Null claims are ignored.
func NewToken(raw map[string]interface{}) Token {
	t := make(Token, len(raw))
	for k, v := range raw {
		if v != nil {
			t[k] = v
		}
	}

	return t
}
