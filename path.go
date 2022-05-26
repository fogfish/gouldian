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
	"github.com/fogfish/gouldian/internal/optics"
)

// Segment union type, make URI type safe
type Segment struct {
	optics *Lens
	path   string
}

/*

URI is an endpoint to match URL of HTTP request. The function takes a sequence
of segment patterns as input. These patterns are either literals or lenses,
where each term corresponds to the path segment. The function do not match
if length of path is not equal to the length of pattern or segment do not
match to pattern

  e := µ.GET( µ.URI(µ.Path("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func URI(segments ...Segment) Routable {
	return func() ([]string, Endpoint) {
		path, lens := segmentsToLens(segments, true)
		return path, segmentsToEndpoint(path, lens)
	}
}

/*

Path is an endpoint to match a single URL segment of HTTP request.
The function takes a path pattern as arguments. The pattern is either literal
or lens. The function do not match if segment do not match to pattern

  e := µ.GET( µ.URI(µ.Path("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Path[T Pattern](segment T) Segment {
	switch v := any(segment).(type) {
	case string:
		return Segment{path: v}
	case Lens:
		return Segment{optics: &v, path: ":"}
	default:
		panic("")
	}
}

/*

PathAny is a synonym of µ.Path("_"), it matches any segments
*/
func PathAny() Segment {
	return Segment{path: "_"}
}

/*

PathAll is an endpoint to match entire remaining path of URI
*/
func PathAll(segment Lens) Segment {
	return Segment{optics: &segment, path: "*"}
}

//
func segmentsToLens(segments []Segment, strict bool) ([]string, []optics.Lens) {
	lens := make([]optics.Lens, 0)
	path := make([]string, 0, len(segments))
	for i, segment := range segments {
		if segment.optics == nil {
			path = append(path, segment.path)
		} else {
			if i == len(segments)-1 && segment.path == "*" {
				path = append(path, "*")
			} else {
				path = append(path, ":")
			}
			lens = append(lens, segment.optics.Lens)
		}
	}

	return path, lens
}

//
func segmentsToEndpoint(path []string, lens []optics.Lens) Endpoint {
	return func(ctx *Context) error {
		if len(ctx.values) != len(lens) {
			return ErrNoMatch
		}

		for i, l := range lens {
			if err := ctx.Put(l, ctx.values[i]); err != nil {
				return err
			}
		}

		return nil
	}
}
