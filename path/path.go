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

/*

Package path defines primitives to match URL of HTTP requests.

	import "github.com/fogfish/gouldian/path"

	endpoint := µ.GET( µ.Path(path.Is("foo"), ...) )
	endpoint(mock.Input(mock.URL("/foo"))) == nil

*/
package path

import (
	"strconv"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
)

// Then is a product of path match arrows. It helps to build a group of
// composable patterns for sup path(es)
/*
func (arrow ArrowPath) Then(x ArrowPath) ArrowPath {
	return func(ctx µ.Context, segments []string) error {
		sz := len(segments)
		at := 0

		for _, f := range []ArrowPath{arrow, x} {
			if sz <= at {
				return NoMatch{}
			}
			switch err := f(segments[at:]).(type) {
			case nil:
				at++
			case Match:
				at = at + err.N
			default:
				return err
			}
		}

		return Match{N: at}
	}
}
*/

// Or is a co-product of path match arrows.
/*
func (arrow ArrowPath) Or(x ArrowPath) ArrowPath {
	return func(segments []string) error {
		for _, f := range []ArrowPath{arrow, x} {
			if err := f(segments); !errors.Is(err, NoMatch{}) {
				return err
			}
		}
		return NoMatch{}
	}
}
*/

/*

Is matches a path segment to defined literal
  e := µ.GET( µ.Path(path.Is("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Is(val string) µ.ArrowPath {
	if val == "*" {
		return func(µ.Context, µ.Segments) error { return nil }
	}

	return func(ctx µ.Context, segments µ.Segments) error {
		if segments[0] == val {
			return nil
		}
		return µ.NoMatch{}
	}
}

/*

Any is a wildcard matcher of path segment
  e := µ.GET( µ.Path(path.Any()) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) == nil
*/
func Any() µ.ArrowPath {
	return func(µ.Context, µ.Segments) error {
		return nil
	}
}

/*

String matches a path segment to context variable of string type.
  e := µ.GET( µ.Path(path.String(FOO)) )
  e(mock.Input(mock.URL("/foo"))) == nil && *ctx.String(FOO) == "foo"
  e(mock.Input(mock.URL("/1"))) == nil && *ctx.String(FOO) == "1"
*/
func String(lens optics.Lens) µ.ArrowPath {
	return func(ctx µ.Context, segments µ.Segments) error {
		ctx.Put(lens, segments[0])
		return nil
	}
}

/*

Int matches a path segment to context variable of int type
  const FOO µ.Symbol = iota
  e := µ.GET( µ.Path(path.Int(FOO)) )
  e(mock.Input(mock.URL("/1"))) == nil && *ctx.Int(FOO) == 1
  e(mock.Input(mock.URL("/foo"))) != nil
*/
func Int(lens optics.Lens) µ.ArrowPath {
	return func(ctx µ.Context, segments µ.Segments) error {
		if value, err := strconv.Atoi(segments[0]); err == nil {
			ctx.Put(lens, value)
			return nil
		}
		return µ.NoMatch{}
	}
}

/*

Float matches a path segment to context variable of float type
  const FOO µ.Symbol = iota
  e := µ.GET( µ.Path(path.Float(FOO)) )
  e(mock.Input(mock.URL("/1.0"))) == nil && *ctx.Float(FOO) == 1.0
  e(mock.Input(mock.URL("/foo"))) != nil
*/
func Float(lens optics.Lens) µ.ArrowPath {
	return func(ctx µ.Context, segments µ.Segments) error {
		if value, err := strconv.ParseFloat(segments[0], 64); err == nil {
			ctx.Put(lens, value)
			return nil
		}
		return µ.NoMatch{}
	}
}

/*

Seq matches path 1 or N segments to closed slice
  var seq []string
  e := µ.GET( µ.Path(path.Seq(seq)) )
  e(mock.Input(mock.URL("/a/b"))) == nil && seq == []string{"a", "b"}
*/
func Seq(lens optics.Lens) µ.ArrowPath {
	return func(ctx µ.Context, segments µ.Segments) error {
		ctx.Put(lens, append([]string{}, segments...))
		return µ.Match{N: len(segments)}
	}
}
