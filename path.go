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
func Path(segments ...interface{}) Routable {
	return func() ([]string, Endpoint) {
		path, lens := segmentsTo(segments, true)

		return path, func(ctx *Context) error {
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
}

//
func PathSeq(segments ...interface{}) Routable {
	return func() ([]string, Endpoint) {
		path, lens := segmentsTo(segments, false)

		return path, func(ctx *Context) error {
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
}

func segmentsTo(segments []interface{}, strict bool) ([]string, []optics.Lens) {
	lens := make([]optics.Lens, 0)
	path := make([]string, 0, len(segments))
	for i, segment := range segments {
		switch v := segment.(type) {
		case string:
			path = append(path, v)
		case optics.Lens:
			if i == len(segments)-1 && !strict {
				path = append(path, "*")
			} else {
				path = append(path, ":")
			}
			lens = append(lens, v)
		}
	}

	return path, lens
}
