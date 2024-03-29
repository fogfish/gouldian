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

	"net/http"

	"github.com/fogfish/gouldian/v2/server/httpd"
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

var lensEcho = µ.Optics1[reqEcho, string]()

func echo() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("echo"), µ.Path(lensEcho)),
		µ.Accept.Text,
		µ.FMap(func(ctx *µ.Context, req *reqEcho) error {
			return ø.Status.OK(
				ø.Server.Set("echo"),
				ø.ContentType.TextPlain,
				ø.Send(req.Echo),
			)
		}),
	)
}
