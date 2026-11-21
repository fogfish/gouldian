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
	"slices"
	"strings"
)

/*
JWTClaim is function type to extract claims from token
*/
type JWTClaim func(Token) string

/*
JWT combinator defines primitives to match JWT token in the HTTP requests.

	  endpoint := µ.GET(
	    µ.JWT(µ.Token.Username, "joedoe"),
	  )

	  endpoint(
	    mock.Input(
				mock.JWT(µ.Token{"username": "joedoe"})
	    )
	  ) == nil
*/
func JWT[T Pattern](claim JWTClaim, val T) Endpoint {
	switch v := any(val).(type) {
	case string:
		return jwtClaim(claim).Is(v)
	case Lens:
		return jwtClaim(claim).To(v)
	default:
		panic("type system failure")
	}
}

type jwtClaim JWTClaim

/*
Is matches a key of JWT to defined literal value
*/
func (claim jwtClaim) Is(val string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return NoMatchJWT
		}

		if claim(ctx.JWT) != val {
			return NoMatchJWT
		}

		return nil
	}
}

/*
To matches key of JWT value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint fails if the claim
value cannot be decoded to the target type. See optics.Lens type for details.
*/
func (claim jwtClaim) To(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return NoMatchJWT
		}

		if val := claim(ctx.JWT); val != "" {
			return ctx.put(SourceJWT, "", lens, val)
		}

		return NoMatchJWT
	}
}

/*
JWTMaybe matches key of JWT to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint does not fail if the claim
is absent, a malformed claim is a failure. See optics.Lens type for details.

	type MyT struct{ Username string }

	username := µ.Optics1[MyT, string]()
	e := µ.GET( µ.JWTMaybe(µ.JWT.Sub).Maybe(username) )
	e(mock.Input(mock.JWT(µ.JWT{"username": "joedoe"}))) == nil
*/
func JWTMaybe(claim JWTClaim, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return NoMatchJWT
		}

		if val := claim(ctx.JWT); val != "" {
			return ctx.put(SourceJWT, "", lens, val)
		}

		return nil
	}
}

/*
JWTOneOf matches a key of JWT if it contains one of the tokens.
The claim is a space-delimited list (e.g. OAuth scope), tokens are compared
as whole words, `ro` does not match `prod:write`.

	µ.GET( µ.JWTOneOf(µ.JWT.Scope, "ro", "rw") )
*/
func JWTOneOf(claim JWTClaim, vals ...string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return NoMatchJWT
		}

		val := strings.Fields(claim(ctx.JWT))
		for _, x := range vals {
			if slices.Contains(val, x) {
				return nil
			}
		}

		return NoMatchJWT
	}
}

/*
JWTAllOf matches a key of JWT if it contains all of the tokens.
The claim is a space-delimited list (e.g. OAuth scope), tokens are compared
as whole words, `ro` does not match `prod:write`.

	µ.GET( µ.JWTAllOf(µ.JWT.Scope, "ro", "rw") )
*/
func JWTAllOf(claim JWTClaim, vals ...string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return NoMatchJWT
		}

		val := strings.Fields(claim(ctx.JWT))
		for _, x := range vals {
			if !slices.Contains(val, x) {
				return NoMatchJWT
			}
		}

		return nil
	}
}
