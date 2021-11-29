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

↣ `Output` is sum type that represents if it is matched on a given input
or not. The library uses `error` type to represent both valid and invalid
variants.

Any `Endpoint A` can be composed with `Endpoint B` into new `Endpoint C`.
It supports two combinators: and-then, or-else.

↣ Use `and-then` to build product Endpoint. The product type matches Input
if each composed function successfully matches it.

↣ Use `or-else` to build co-product Endpoint. The co-product is also known
as sum-type matches first successful function.

Endpoint life-cycle - each incoming HTTP request is wrapped with `Input`
and applied to an endpoint. A returned error-like results is checked
against successful Output or NoMatch error. All these machinery is handled
by the libray, you should only dare to declare Endpoint from ready made
primitives.

gouldian library delivers set of built-in endpoints to deal with HTTP
request processing.

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
	return func(http *Context) (err error) {
		switch err := a(http).(type) {
		case NoMatch:
			return b(http)
		default:
			return err
		}
	}
}

/*

Routable is endpoint with routing metadata
*/
type Routable func() ([]string, Endpoint)

/*

Router is data structure that holds routing information,
convertable to Endpoint
*/
type Router interface {
	Endpoint() Endpoint
}

// NoMatch is returned by Endpoint if Context is not matched.
type NoMatch int

func (err NoMatch) Error() string {
	return "No Match"
}

// ErrNoMatch constant
var ErrNoMatch error = NoMatch(255)

/*

Endpoints is sequence of Endpoints
*/
type Endpoints []Endpoint

/*

Join builds product endpoint from sequence
*/
func (seq Endpoints) Join(ctx *Context) (err error) {
	for _, f := range seq {
		err = f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

/*

Or builds co-product endpoint from sequence
*/
func (seq Endpoints) Or(ctx *Context) (err error) {
	for _, f := range seq {
		x := f(ctx)
		switch err := x.(type) {
		case NoMatch:
			continue
		default:
			return err
		}
	}
	return ErrNoMatch
}
