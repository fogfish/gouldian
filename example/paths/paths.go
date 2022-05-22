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
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			root(),
			path(),
			file(),
		),
	)
}

/*

curl http://localhost:8080/
curl http://localhost:8080/path/value
curl http://localhost:8080/file/so/me/va/lue

*/

/*

matches root (/)
*/
func root() µ.Routable {
	return µ.GET(
		µ.URI(),
		func(ctx *µ.Context) error {
			return µ.Status.OK(
				µ.WithText("matches root (/)"),
			)
		},
	)
}

/*

matches path with param (/echo/:text)
*/
type param struct {
	Text string
}

var lensEcho = µ.Optics1[param, string]()

func path() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("path"), µ.Path(lensEcho)),
		func(ctx *µ.Context) error {
			var req param
			if err := µ.FromContext(ctx, &req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				µ.WithText("matches /path/" + req.Text),
			)
		},
	)
}

/*

matches root (/file/file+)
*/
func file() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("file"), µ.PathAll(lensEcho)),
		func(ctx *µ.Context) error {
			var req param
			if err := µ.FromContext(ctx, &req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				µ.WithText("matches /file/" + req.Text),
			)
		},
	)
}
