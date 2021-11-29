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
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			echo(),
		),
	)

}

type reqEcho struct {
	Echo string
}

var lensEcho = optics.ForProduct1(reqEcho{})

func echo() µ.Routable {
	return µ.GET(
		µ.Path("echo", lensEcho),
		µ.FMap(func(ctx *µ.Context) error {
			var req reqEcho
			if err := ctx.Get(&req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				headers.ContentType.Value(headers.TextPlain),
				headers.Server.Value("echo"),
				µ.WithText(req.Echo),
			)
		}),
	)
}
