package gouldian

import "errors"

// Endpoint abstarcts HTTP endpoint as a function.
// It takes HTTP request and returns value of some type.
//
// The composition is an essential part of the type.
// Endpoint A and B can be composed into new Endpoint C.
//
// The type supports two combinators: and-then, or-else.
type Endpoint func(*Input) error

// Then combines
func (a Endpoint) Then(b Endpoint) Endpoint {
	return func(http *Input) (err error) {
		if err = a(http); err == nil {
			return b(http)
		}
		return err
	}
}

// Or combines
func (a Endpoint) Or(b Endpoint) Endpoint {
	return func(http *Input) (err error) {
		if err = a(http); errors.Is(err, NoMatch{}) {
			return b(http)
		}
		return err
	}
}
