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

/*

Server is any data structure convertable to Endpoint
*/
type Server interface {
	Endpoint() Endpoint
}

type Builder func(*Node) *Node

//
func (n *Node) EvalRoot(http *Context) (err error) {
	path := http.Request.URL.Path
	i, node := RGet(n, path)

	// fmt.Println("s > ", path, i, node)

	// panic('x')

	if len(path) == i && node.Endpoint != nil {
		// fmt.Println("s > ", stack)
		return node.Endpoint(http)
	}

	// fmt.Println(i, path, node)
	// panic('x')

	// for _, cn := range n.Children {
	// 	switch v := cn.Eval(http).(type) {
	// 	case NoMatch:
	// 		continue
	// 	case *NoMatch:
	// 		continue
	// 	default:
	// 		return v
	// 	}
	// }
	return ErrNoMatch
}

/*
func (n *Node) Eval(http *Context) (err error) {
	x := n.Endpoint(http)
	if x == nil {
		if len(n.Children) == 0 {
			return nil
		}

		for _, cn := range n.Children {
			x := cn.Eval(http)
			if x == nil {
				return x
			}
		}

		return ErrNoMatch
	}
	return x

	// switch err := x.(type) {
	// case nil:
	// 	if len(n.Children) == 0 {
	// 		return nil
	// 	}

	// 	for _, cn := range n.Children {
	// 		x := cn.Eval(http)
	// 		if x == nil {
	// 			return x
	// 		}
	// 		// switch v := cn.Eval(http).(type) {
	// 		// case NoMatch:
	// 		// 	continue
	// 		// case *NoMatch:
	// 		// 	continue
	// 		// default:
	// 		// 	return v
	// 		// }
	// 	}
	// 	return ErrNoMatch
	// case NoMatch:
	// 	return err
	// case *NoMatch:
	// 	return err
	// default:
	// 	return err
	// }
}
*/

// NoMatch is returned by Endpoint if Context is not matched.
type NoMatch int

func (err NoMatch) Error() string {
	return "No Match"
}

// ErrNoMatch constant
var ErrNoMatch error = NoMatch(255) // &NoMatch{}

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
		case *NoMatch:
			return b(http)
		default:
			return err
		}
		// if err = a(http); !errors.Is(err, NoMatch{}) {
		// 	return err
		// }
		// return b(http)
	}
}

func JoinN(seq ...Builder) Builder {
	return func(n *Node) *Node {
		a := n
		for _, f := range seq {
			a = f(a)
		}

		return n
	}
}

// Join builds a product endpoint from sequence
func Join(seq ...Endpoint) Endpoint {
	return tseq(seq).j
	// if len(seq) == 1 {
	// 	return seq[0]
	// }

	// return func(http *Context) (err error) {
	// 	for _, f := range seq {
	// 		if err = f(http); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }
}

func Join2(a, b Endpoint) Endpoint {
	return func(http *Context) (err error) {
		if err = a(http); err != nil {
			return err
		}
		if err = b(http); err != nil {
			return err
		}
		return nil
	}
}

type tseq []Endpoint

func (seq tseq) j(http *Context) (err error) {
	for _, f := range seq {
		err = f(http)
		if err != nil {
			return err
		}
	}
	return nil
}

func (seq tseq) r(http *Context) (err error) {
	for _, f := range seq {
		x := f(http)
		switch err := x.(type) {
		case NoMatch:
			continue
		case *NoMatch:
			continue
		default:
			return err
		}
	}
	return ErrNoMatch
}

// Or joins sequence of Endpoint(s) to co-product Endpoint.
func Or(seq ...Endpoint) Endpoint {
	return tseq(seq).r
	// if len(seq) == 1 {
	// 	return seq[0]
	// }

	// return func(http *Context) (err error) {
	// 	for _, f := range seq {
	// 		x := f(http)
	// 		switch err := x.(type) {
	// 		case NoMatch:
	// 			continue
	// 		case *NoMatch:
	// 			continue
	// 		default:
	// 			return err
	// 		}
	// 	}
	// 	return ErrNoMatch
	// }
}
