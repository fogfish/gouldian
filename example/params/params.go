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

var lensString = µ.Optics1[paramString, string]()

func qString() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("echo")),
		µ.Param("q", lensString),
		µ.FMap(func(ctx *µ.Context, req *paramString) error {
			return µ.Status.OK(
				µ.WithText(fmt.Sprintf("query string (%s)", req.Value)),
			)
		}),
	)
}

/*

matches integer query parameter /echo?v=number
*/
type paramInt struct {
	Value int
}

var lensInt = µ.Optics1[paramInt, int]()

func qInt() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("echo")),
		µ.Param("v", lensInt),
		µ.FMap(func(ctx *µ.Context, req *paramInt) error {
			return µ.Status.OK(
				µ.WithText(fmt.Sprintf("query number (%d)", req.Value)),
			)
		}),
	)
}
