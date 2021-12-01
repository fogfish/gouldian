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
			get(),
			post(),
			put(),
			patch(),
			delete(),
		),
	)
}

/*

curl -XGET http://localhost:8080/echo
curl -XPOST http://localhost:8080/echo
curl -XPUT http://localhost:8080/echo
curl -XPATCH http://localhost:8080/echo
curl -XDELETE http://localhost:8080/echo

*/

func get() µ.Routable {
	return µ.GET(µ.Path("echo"), serve)
}

func post() µ.Routable {
	return µ.POST(µ.Path("echo"), serve)
}

func put() µ.Routable {
	return µ.PUT(µ.Path("echo"), serve)
}

func patch() µ.Routable {
	return µ.PATCH(µ.Path("echo"), serve)
}

func delete() µ.Routable {
	return µ.DELETE(µ.Path("echo"), serve)
}

func serve(ctx *µ.Context) error {
	return µ.Status.OK(
		µ.WithText(ctx.Request.Method + " echo"),
	)
}
