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

// Package httpd adapts the gouldian server to net/http
package httpd

import (
	"net/http"

	µ "github.com/fogfish/gouldian/v3"
)

/*
Serve builds http.Handler for sequence of endpoints

	http.ListenAndServe(":8080", httpd.Serve( ... ))
*/
func Serve(endpoints ...µ.Routable) http.Handler {
	return Handler(µ.NewServer().Mount(endpoints...))
}

/*
ServeAndCommit builds http.Handler for sequence of endpoints.
It executes commit function after each request.

	http.ListenAndServe(":8080", httpd.ServeAndCommit(commit, ... ))
*/
func ServeAndCommit(commit func(), endpoints ...µ.Routable) http.Handler {
	h := Serve(endpoints...)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		commit()
	})
}

/*
Handler adapts the server to http.Handler. The context of the request is
propagated to endpoints.

	api := µ.NewServer(µ.WithCORS(...)).Mount( ... )
	http.ListenAndServe(":8080", httpd.Handler(api))
*/
func Handler(srv *µ.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := srv.NewContext(r.Context(), r)
		// canonical key, direct lookup avoids the cost of Header.Get
		if id := r.Header["X-Request-Id"]; len(id) != 0 {
			ctx.RequestID = id[0]
		}
		out := srv.Dispatch(ctx)

		if len(out.Headers) != 0 {
			header := w.Header()
			for _, h := range out.Headers {
				header.Add(h.Header, h.Value)
			}
		}
		w.WriteHeader(out.Status)
		if len(out.Body) > 0 {
			_, _ = w.Write(out.Body)
		}

		out.Free()
		srv.Release(ctx)
	})
}
