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
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			text(),
			json(),
		),
	)
}

/*

curl -v http://localhost:8080/echo -H "Accept: text/plain"
curl -v http://localhost:8080/echo -H "Accept: application/json"

*/

func text() µ.Routable {
	return µ.GET(
		µ.Path("echo"),
		headers.Accept.Is(headers.TextPlain),
		func(ctx *µ.Context) error {
			return µ.Status.OK(
				headers.ContentType.Value(headers.TextPlain),
				µ.WithText("hello world."),
			)
		},
	)
}

func json() µ.Routable {
	return µ.GET(
		µ.Path("echo"),
		headers.Accept.Is(headers.ApplicationJSON),
		func(ctx *µ.Context) error {
			return µ.Status.OK(
				headers.ContentType.Value(headers.ApplicationJSON),
				µ.WithText(`{"text": "hellow world."}`),
			)
		},
	)
}
