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
	"github.com/fogfish/gouldian/optics"
)

//
type pathArrow func(*Context, string) error

/*

Path is an endpoint to match URL of HTTP request. The function takes a path
pattern as arguments. The pattern is sequence of either literals or lenses,
where each term corresponds to the path segment. The function do not match
if length of path is not equal to the length of pattern or segment do not
match to pattern

  e := µ.GET( µ.Path("foo") )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Path(segments ...interface{}) Endpoint {
	x := mkPathMatcher(segments)
	return seth(x).mkPathEndpointA
}

// segment is custom implementation of strings.Split
func segment(path string, a int) (int, string) {
	for i := a + 1; i < len(path); i++ {
		if path[i] == '/' {
			return i, path[a+1 : i]
		}
	}
	return len(path) - 1, path[a+1:]
}

type seth []pathArrow

func (segments seth) mkPathEndpointA(ctx *Context) (err error) {
	return ErrNoMatch
}

func (segments seth) mkPathEndpoint(ctx *Context) (err error) {
	path := ctx.Request.URL.Path
	last := len(path) - 1

	/* */
	if len(ctx.segments) == 0 {
		// fmt.Println("building " + ctx.Request.RequestURI)
		for hd := 0; hd < last; {
			tl, segment := segment(path, hd)
			ctx.segments = append(ctx.segments, segment)
			hd = tl
		}
	}
	// fmt.Println("using " + ctx.Request.RequestURI)
	// fmt.Println(ctx.segments)

	if len(segments) != len(ctx.segments) {
		return ErrNoMatch
	}

	for i, f := range segments {
		seg := ctx.segments[i]
		if err = f(ctx, seg); err != nil {
			return err
		}
	}
	return nil
	/* */

	/*
		hd := 0
		for at, f := range segments {
			tl, segment := segment(path, hd)
			if hd == tl {
				return NoMatch{}
			}
			if err := f(ctx, segment); err != nil {
				return err
			}
			hd = tl

			// url resource path is shorter than pattern
			if hd == last && at != len(segments)-1 {
				return NoMatch{}
			}
		}

		// url resource path is not consumed by the pattern
		if hd != last {
			return NoMatch{}
		}

		return nil
	*/
}

/*

PathSeq is like Path but last element in the pattern must be lens that lifts
the tail of path.

  e := µ.GET( µ.PathSeq("foo", suffix) )
  e(mock.Input(mock.URL("/foo/bar"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func PathSeq(arrows ...interface{}) Endpoint {
	return nil
	// return mkPathSeqEndpoint(mkPathMatcher(arrows))
}

// func mkPathSeqEndpoint(segments []pathArrow) Endpoint {
// 	return func(ctx *Context) error {
// 		path := ctx.Request.URL.Path
// 		// last := len(path) - 1

// 		hd := 0
// 		last := len(segments) - 1
// 		for at := 0; at < last; at++ {
// 			tl, segment := segment(path, hd)
// 			if hd == tl {
// 				return NoMatch{}
// 			}
// 			if err := segments[at](ctx, segment); err != nil {
// 				return err
// 			}
// 			hd = tl
// 		}

// 		// url resource path is not consumed by the pattern
// 		if hd == len(path)-1 {
// 			return NoMatch{}
// 		}

// 		// url resource path consume suffix
// 		if err := segments[last](ctx, path[hd+1:]); err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// }

func mkPathMatcher(arrows []interface{}) []pathArrow {
	seq := make([]pathArrow, len(arrows))

	for i, arrow := range arrows {
		switch v := arrow.(type) {
		case string:
			switch v {
			case Any:
				seq[i] = pathAny()
			default:
				seq[i] = pathIs(v)
			}
		case optics.Lens:
			seq[i] = pathTo(v)
		default:
			seq[i] = pathNone()
		}
	}

	return seq
}

/*

Is matches a path segment to defined literal
  e := µ.GET( µ.Path("foo") )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func pathIs(val string) pathArrow {
	// fn := func(ctx *Context, segment string) error {
	// 	if len(segment) != len(val) {
	// 		return ErrNoMatch
	// 	}

	// 	// if segment[0] != val[0] {
	// 	// 	return ErrNoMatch
	// 	// }
	// 	// if segment[1] != val[1] {
	// 	// 	return ErrNoMatch
	// 	// }

	// 	if segment == val {
	// 		return nil
	// 	}
	// 	return ErrNoMatch
	// }

	// return fn
	return pth(val).pis
}

type pth string

func (val pth) pis(ctx *Context, segment string) error {
	// if len(segment) != len(val) {
	// 	return ErrNoMatch
	// }

	// if segment == string(val) {
	// 	return nil
	// }
	return ErrNoMatch
}

/*

None matches nothing
*/
func pathNone() pathArrow {
	return func(*Context, string) error {
		return ErrNoMatch
	}
}

/*

Any is a wildcard matcher of path segment
  e := µ.GET( µ.Path(path.Any) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) == nil
*/
func pathAny() pathArrow {
	return func(*Context, string) error {
		return nil
	}
}

/*

Lifts the path segment to lens
*/
func pathTo(l optics.Lens) pathArrow {
	return func(ctx *Context, segment string) error {
		return ctx.Put(l, segment)
	}
}
