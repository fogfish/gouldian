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
Endpoint is a composable function that abstract HTTP endpoint.
The function takes HTTP request and returns value of some type:
`Context => Output`.

↣ `Context` is a wrapper over HTTP request with additional context.

↣ `Output` is sum type that represents the outcome of the endpoint:

	nil       matched, continue the product
	NoMatch   not my route, Or tries the next endpoint (zero allocation)
	*Failure  my route, bad input (malformed value, empty body, ...),
	          terminal, Or does not try other endpoints
	*Output   terminal response of the handler
	error     any other error, it is 500

Any `Endpoint A` can be composed with `Endpoint B` into new `Endpoint C`.
It supports two combinators: and-then, or-else.

↣ Use `and-then` to build product Endpoint. The product type matches Input
if each composed function successfully matches it.

↣ Use `or-else` to build co-product Endpoint. The co-product is also known
as sum-type matches first successful function. A branch that does not match
leaves no traces in the context.
*/
type Endpoint func(*Context) error

// Then builds product Endpoint
func (a Endpoint) Then(b Endpoint) Endpoint {
	return func(http *Context) (err error) {
		if err = a(http); err == nil {
			return b(http)
		}
		return err
	}
}

// Or builds co-product Endpoint
func (a Endpoint) Or(b Endpoint) Endpoint {
	return func(http *Context) error {
		// bindings captured by a rejected branch must not reach the next one
		mark := len(http.morphism)
		err := a(http)
		na, isNoMatch := err.(NoMatch)
		if !isNoMatch {
			return err
		}

		http.rollback(mark)
		err = b(http)
		if nb, ok := err.(NoMatch); ok && na > nb {
			return na
		}
		return err
	}
}

// Endpoints is sequence of Endpoints
type Endpoints []Endpoint

// Join builds product endpoint from sequence
func (seq Endpoints) Join(ctx *Context) (err error) {
	for _, f := range seq {
		err = f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Or builds co-product endpoint from sequence
func (seq Endpoints) Or(ctx *Context) error {
	// bindings captured by a rejected branch must not reach the next one
	mark := len(ctx.morphism)
	reason := NoMatchAny
	for _, f := range seq {
		err := f(ctx)
		no, isNoMatch := err.(NoMatch)
		if !isNoMatch {
			return err
		}
		ctx.rollback(mark)
		reason = max(reason, no)
	}
	return reason
}

// Join builds product endpoint from sequence
func Join(seq ...Endpoint) Endpoint { return Endpoints(seq).Join }

// Or builds co-product endpoint from sequence
func Or(seq ...Endpoint) Endpoint { return Endpoints(seq).Or }
