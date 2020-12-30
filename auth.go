//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian

// AccessToken is a container for user identity
type AccessToken struct {
	Jti      string
	Iss      string
	Exp      string
	Sub      string
	Scope    string
	UserID   string
	ClientID string
}

//
func mkAccessToken(raw map[string]interface{}) AccessToken {
	asString := func(id string) string {
		if val, ok := raw[id]; ok {
			return val.(string)
		}
		return ""
	}

	return AccessToken{
		Jti:      asString("jti"),
		Iss:      asString("iss"),
		Exp:      asString("exp"),
		Sub:      asString("sub"),
		Scope:    asString("scope"),
		UserID:   asString("username"),
		ClientID: asString("client_id"),
	}
}
