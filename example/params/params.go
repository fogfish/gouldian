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
	"fmt"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			qString(),
			qInt(),
		),
	)
}

/*

curl http://localhost:8080/echo?q=value
curl http://localhost:8080/echo?v=12345

*/

/*

matches string query parameter /echo?q=text
*/
type paramString struct {
	Value string
}

var lensString = optics.ForProduct1(paramString{})

func qString() µ.Routable {
	return µ.GET(
		µ.Path("echo"),
		µ.Param("q").To(lensString),
		func(ctx *µ.Context) error {
			var req paramString
			if err := ctx.Get(&req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				µ.WithText(fmt.Sprintf("query string (%s)", req.Value)),
			)
		},
	)
}

/*

matches integer query parameter /echo?v=number
*/
type paramInt struct {
	Value int
}

var lensInt = optics.ForProduct1(paramInt{})

func qInt() µ.Routable {
	return µ.GET(
		µ.Path("echo"),
		µ.Param("v").To(lensInt),
		func(ctx *µ.Context) error {
			var req paramInt
			if err := ctx.Get(&req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				µ.WithText(fmt.Sprintf("query number (%d)", req.Value)),
			)
		},
	)
}
