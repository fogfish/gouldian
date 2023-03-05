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

package main

import (
	µ "github.com/fogfish/gouldian/v2"
	ø "github.com/fogfish/gouldian/v2/output"
	"github.com/fogfish/gouldian/v2/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			text(),
			json(),
			form(),
		),
	)
}

/*

curl -v http://localhost:8080/echo \
  -H "Content-Type: text/plain" \
  -d 'Hello World.'

curl -v http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"string": "Hello World.", "number": 1010}'

curl -v http://localhost:8080/echo \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'string=Hello+World.&number=1010'

*/

/*
matches string payload
*/
type reqText struct {
	Value string
}

var lensText = µ.Optics1[reqText, string]()

func text() µ.Routable {
	return µ.POST(
		µ.URI(µ.Path("echo")),
		µ.ContentType.TextPlain,
		µ.Body(lensText),
		µ.FMap(func(ctx *µ.Context, req *reqText) error {
			return ø.Status.OK(
				ø.ContentType.TextPlain,
				ø.Send(req.Value),
			)
		}),
	)
}

/*
matches JSON payload
*/
type reqJSON struct {
	Value myJSON
}

type myJSON struct {
	String string `json:"string,omitempty"`
	Number int    `json:"number,omitempty"`
}

var lensJSON = µ.Optics1[reqJSON, myJSON]()

func json() µ.Routable {
	return µ.POST(
		µ.URI(µ.Path("echo")),
		µ.ContentType.ApplicationJSON,
		µ.Body(lensJSON),
		µ.FMap(func(ctx *µ.Context, req *reqJSON) error {
			return ø.Status.OK(
				ø.ContentType.TextPlain,
				ø.Send(req.Value),
			)
		}),
	)
}

/*
matches Form payload
*/
type reqForm struct {
	Value myJSON `content:"form"`
}

var lensForm = µ.Optics1[reqForm, myJSON]()

func form() µ.Routable {
	return µ.POST(
		µ.URI(µ.Path("echo")),
		µ.ContentType.Form,
		µ.Body(lensForm),
		µ.FMap(func(ctx *µ.Context, req *reqForm) error {
			return ø.Status.OK(
				ø.ContentType.TextPlain,
				ø.Send(req.Value),
			)
		}),
	)
}
