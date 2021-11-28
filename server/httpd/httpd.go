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

package httpd

import (
	"context"
	"fmt"
	µ "github.com/fogfish/gouldian"
	"net/http"
	"sync"
)

/*

Serve builds http.Handler for sequence of endpoints

  http.ListenAndServe(":8080", httpd.Server( ... ))
*/
func Serve(endpoints ...µ.Route) http.Handler {
	root := µ.NewRoutes()
	for _, endpoint := range endpoints {
		endpoint(root)
	}

	routes := &routes{
		endpoint: root.Endpoint(),
	}

	// root.Println()
	routes.pool.New = func() interface{} {
		return µ.NewContext(context.Background())
	}

	return routes
}

type routes struct {
	endpoint µ.Endpoint
	pool     sync.Pool
}

func (routes *routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := routes.pool.Get().(*µ.Context)
	req.Free()
	req.Request = r

	switch v := routes.endpoint(req).(type) {
	case *µ.Output:
		routes.output(w, v)
	case µ.NoMatch:
		failure := µ.Status.NotImplemented(
			µ.WithIssue(fmt.Errorf("NoMatch %s", r.URL.Path)),
		).(*µ.Output)
		routes.output(w, failure)
	}

	routes.pool.Put(req)
}

func (routes *routes) output(w http.ResponseWriter, out *µ.Output) {
	for _, h := range out.Headers {
		w.Header().Set(h.Header, h.Value)
	}
	w.WriteHeader(int(out.Status))

	if len(out.Body) > 0 {
		w.Write([]byte(out.Body))
	}
	out.Free()
}
