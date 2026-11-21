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
	"fmt"
	"strings"
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
	path, lens := segmentsToLens(segments)
	endpoint := segmentsToEndpoint(lens)
	return func() Route {
		return Route{Path: path, Endpoint: endpoint}
	}
}

/*
Path is an endpoint to match a single URL segment of HTTP request.
The function takes a path pattern as arguments. The pattern is either literal
or lens. The function do not match if segment do not match to pattern

	e := µ.GET( µ.URI(µ.Path("foo")) )
	e(mock.Input(mock.URL("/foo"))) == nil
	e(mock.Input(mock.URL("/bar"))) != nil

A segment that cannot be decoded by the lens does not match the route,
it allows to overload the path (/user/:id by int and by string).
*/
func Path[T Pattern](segment T) Segment {
	switch v := any(segment).(type) {
	case string:
		return Segment{path: v}
	case Lens:
		return Segment{optics: &v, path: ":"}
	default:
		panic("type system failure")
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

func segmentsToLens(segments []Segment) ([]string, []Lens) {
	lens := make([]Lens, 0)
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
			lens = append(lens, *segment.optics)
		}
	}

	return path, lens
}

func segmentsToEndpoint(lens []Lens) Endpoint {
	return func(ctx *Context) error {
		if len(ctx.values) != len(lens) {
			return NoMatchPath
		}

		for i, l := range lens {
			if err := ctx.putPath(l, ctx.values[i]); err != nil {
				return err
			}
		}

		return nil
	}
}

/*
URIf builds Routable from a URI pattern, the "f" is for "format": the
pattern string is followed by the lenses it binds, like Printf. It is the string
form of URI. The syntax is a subset of Go 1.22 http.ServeMux patterns (braces,
not % verbs), extended with query parameters (RFC 6570).

	µ.GET(
	  µ.URIf("/users/{id}/files/{path...}{?limit,q}", lID, lPath, lLimit, lQ),
	)

↣ `{name}` matches a path segment, it is bound to the next lens;

↣ `{name...}` matches the rest of the path (the last segment only);

↣ `{_}` or `{}` matches any segment without binding;

↣ `{?a,b}` optional query parameters, each is bound to the next lens. A param
that is not defined is ignored, a malformed one is Failure (see ParamMaybe).

Literals are matched as is. The pattern is checked when it is constructed:
the number of placeholders must be equal to the number of lenses, otherwise it
panics with the pattern in the message.
*/
func URIf(uri string, lenses ...Lens) Routable {
	panicf := func(format string, args ...any) {
		panic(fmt.Errorf("gouldian: invalid URIf pattern %q: %s", uri, fmt.Sprintf(format, args...)))
	}

	pattern, query, hasQuery := strings.Cut(uri, "{?")
	if hasQuery {
		if !strings.HasSuffix(query, "}") {
			panicf("query part is not closed")
		}
		query = strings.TrimSuffix(query, "}")
	}

	if !strings.HasPrefix(pattern, "/") {
		panicf("must start with /")
	}

	var segments []Segment
	used := 0
	nextLens := func() Lens {
		if used >= len(lenses) {
			panicf("has more placeholders than lenses (%d)", len(lenses))
		}
		used++
		return lenses[used-1]
	}

	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	if pattern == "/" {
		parts = nil
	}
	for i, part := range parts {
		switch {
		case part == "":
			panicf("empty segment")
		case strings.HasPrefix(part, "{") && strings.HasSuffix(part, "...}"):
			if i != len(parts)-1 {
				panicf("{%s} must be the last segment", part[1:len(part)-1])
			}
			segments = append(segments, PathAll(nextLens()))
		case strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}"):
			name := part[1 : len(part)-1]
			if name == "" || name == "_" {
				segments = append(segments, PathAny())
			} else {
				segments = append(segments, Path(nextLens()))
			}
		case strings.ContainsAny(part, "{}"):
			panicf("placeholder must be the entire segment, got %q", part)
		default:
			segments = append(segments, Path(part))
		}
	}

	var params Endpoints
	if hasQuery {
		for _, key := range strings.Split(query, ",") {
			key = strings.TrimSpace(key)
			if key == "" {
				panicf("empty query parameter")
			}
			params = append(params, ParamMaybe(key, nextLens()))
		}
	}

	if used != len(lenses) {
		panicf("has %d placeholders but %d lenses", used, len(lenses))
	}

	uri_ := URI(segments...)
	if len(params) == 0 {
		return uri_
	}

	return func() Route {
		route := uri_()
		route.Endpoint = append(Endpoints{route.Endpoint}, params...).Join
		return route
	}
}
