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

import "github.com/fogfish/gouldian/optics"

const (
	// Any constant matches any term
	Any = "_"
)

/*

DELETE composes product Endpoint match HTTP DELETE request.
  e := µ.DELETE()
  e(mock.Input(mock.Method("DELETE"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func DELETE(arrows ...Endpoint) Endpoint {
	return Method("DELETE").Then(Join(arrows...))
}

/*

GET composes product Endpoint match HTTP GET request.
  e := µ.GET()
  e(mock.Input(mock.Method("GET"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func GET(arrows ...Endpoint) Endpoint {
	return Method("GET").Then(Join(arrows...))
}

/*

PATCH composes product Endpoint match HTTP PATCH request.
  e := µ.PATCH()
  e(mock.Input(mock.Method("PATCH"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PATCH(arrows ...Endpoint) Endpoint {
	return Method("PATCH").Then(Join(arrows...))
}

/*

POST composes product Endpoint match HTTP POST request.
  e := µ.POST()
  e(mock.Input(mock.Method("POST"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func POST(arrows ...Endpoint) Endpoint {
	return Method("POST").Then(Join(arrows...))
}

/*

PUT composes product Endpoint match HTTP PUT request.
  e := µ.PUT()
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PUT(arrows ...Endpoint) Endpoint {
	return Method("PUT").Then(Join(arrows...))
}

/*

ANY composes product Endpoint match HTTP PUT request.
  e := µ.ANY()
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) == nil
*/
func ANY(arrows ...Endpoint) Endpoint {
	return Method(Any).Then(Join(arrows...))
}

// Method is an endpoint to match HTTP verb request
func Method(verb string) Endpoint {
	if verb == Any {
		return func(req *Input) error {
			req.Context.Free()
			return nil
		}
	}

	return func(req *Input) error {
		if req.Method == verb {
			req.Context.Free()
			return nil
		}
		return NoMatch{}
	}
}

// Body decodes HTTP request body to struct
func Body(lens optics.Lens) Endpoint {
	return func(req *Input) error {
		if len(req.Payload) != 0 || req.Stream != nil {
			if err := req.ReadAll(); err != nil {
				return err
			}

			return req.Context.Put(lens, req.Payload)
		}

		return NoMatch{}
	}
}

/*

FMap applies clojure to matched HTTP request,
taking the execution context as the input to closure
*/
// func FMap(f func(Context) error) Endpoint {
// 	return func(req *Input) error { return f(req.Context) }
// }
