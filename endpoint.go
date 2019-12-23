package gouldian

import "errors"

// Endpoint is a composable function that abstract HTTP endpoint.
// The function takes HTTP request and returns value of some type:
// `Input => Output`.
//
// ↣ `Input` is a wrapper over Lambda AWS Gateway Event with additional
// context.
//
// ↣ `Output` is sum type that represents if it is matched on a given input
// or not. The library uses `error` type to represent both valid and invalid
// variants.
//
// Any `Endpoint A` can be composed with `Endpoint B` into new `Endpoint C`.
// It supports two combinators: and-then, or-else.
//
// ↣ Use `and-then` to build product Endpoint. The product type matches Input
// if each composed function successfully matches it.
//
// ↣ Use `or-else` to build co-product Endpoint. The co-product is also known
// as sum-type matches first successful function.
//
// Endpoint life-cycle - each incoming HTTP request is wrapped with `Input`
// and applied to an endpoint. A returned error-like results is checked
// against successful Output or NoMatch error. All these machinery is handled
// by the libray, you should only dare to declare Endpoint from ready made
// primitives.
//
// gouldian library delivers set of built-in endpoints to deal with HTTP
// request processing.
type Endpoint func(*Input) error

// Then build product Endpoint.
func (a Endpoint) Then(b Endpoint) Endpoint {
	return func(http *Input) (err error) {
		if err = a(http); err == nil {
			return b(http)
		}
		return err
	}
}

// Or build co-product Endpoint.
func (a Endpoint) Or(b Endpoint) Endpoint {
	return func(http *Input) (err error) {
		if err = a(http); errors.Is(err, NoMatch{}) {
			return b(http)
		}
		return err
	}
}

// JoinOr joins sequence of Endpoint(s) to co-product Endpoint.
func JoinOr(seq ...Endpoint) Endpoint {
	if len(seq) == 1 {
		return seq[0]
	}

	a := seq[0]
	for _, b := range seq[1:] {
		a = a.Or(b)
	}
	return a
}
