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

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/it/v2"
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
	it.Then(t).
		Must(it.Equal(err, nil))

	token := µ.NewToken(jwt)
	it.Then(t).
		Should(it.Equal(token.Jti(), "jti")).
		Should(it.Equal(token.Iss(), "iss")).
		Should(it.Equal(token.Exp(), "exp")).
		Should(it.Equal(token.Sub(), "sub")).
		Should(it.Equal(token.Scope(), "scope")).
		Should(it.Equal(token.Username(), "username")).
		Should(it.Equal(token.ClientID(), "client_id"))
}

func TestErrNoMatch(t *testing.T) {
	it.Then(t).
		Should(it.Equal(µ.ErrNoMatch.Error(), "No Match"))
}
