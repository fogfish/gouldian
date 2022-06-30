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
	"github.com/fogfish/logger"
	"net/http"
	"sync"
)

/*

Serve builds http.Handler for sequence of endpoints

  http.ListenAndServe(":8080", httpd.Server( ... ))
*/
func Serve(endpoints ...µ.Routable) http.Handler {
	routes := &routes{
		endpoint: µ.NewRoutes(endpoints...).Endpoint(),
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
	case nil:
	case *µ.Output:
		routes.output(w, r, v)
	case µ.NoMatch:
		failure := µ.Status.NotImplemented(
			µ.WithIssue(fmt.Errorf("NoMatch %s", r.URL.Path)),
		).(*µ.Output)
		routes.output(w, r, failure)
	default:
		failure := µ.Status.InternalServerError(
			µ.WithIssue(fmt.Errorf("unknown response %s", r.URL.Path)),
		).(*µ.Output)
		routes.output(w, r, failure)
	}

	routes.pool.Put(req)
}

func (routes *routes) output(w http.ResponseWriter, r *http.Request, out *µ.Output) {
	if out.Failure != nil {
		logger.Error("%s %v", r.RequestURI, out.Failure)
	}

	for _, h := range out.Headers {
		w.Header().Set(h.Header, h.Value)
	}
	w.WriteHeader(int(out.Status))

	if len(out.Body) > 0 {
		w.Write([]byte(out.Body))
	}
	out.Free()
}
