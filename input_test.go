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

package gouldian_test

import (
	"encoding/json"
	"testing"

	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/it"
)

func TestAccessToken(t *testing.T) {
	raw := `{
		"jti": "jti",
		"iss": "iss",
		"exp": "exp",
		"sub": "sub",
		"scope": "scope",
		"username": "username",
		"client_id": "client_id"
	}`

	var jwt map[string]interface{}
	err := json.Unmarshal([]byte(raw), &jwt)
	it.Ok(t).If(err).Must().Equal(nil)

	token := µ.NewToken(jwt)
	it.Ok(t).
		If(token.Jti()).Equal("jti").
		If(token.Iss()).Equal("iss").
		If(token.Exp()).Equal("exp").
		If(token.Sub()).Equal("sub").
		If(token.Scope()).Equal("scope").
		If(token.Username()).Equal("username").
		If(token.ClientID()).Equal("client_id")
}

func TestErrNoMatch(t *testing.T) {
	it.Ok(t).If(µ.ErrNoMatch.Error()).Should().Equal("No Match")
}
