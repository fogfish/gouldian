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

package gouldian_test

import (
	"context"
	"fmt"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

//
// Microbenchmark
//

/*

Endpoints Or

*/

var (
	pathA = µ.Path("resu", name)
	pathB = µ.Path("user", name)
	apiOr = µ.Or(
		µ.GET(pathA),
		µ.GET(pathA),
		µ.GET(pathA),
		µ.GET(pathA),
		µ.GET(pathB),
		// func(c *µ.Context) error { return µ.NoMatch{} },
		// func(c *µ.Context) error { return µ.NoMatch{} },
		// func(c *µ.Context) error { return µ.NoMatch{} },
		// func(c *µ.Context) error { return µ.NoMatch{} },
		// func(c *µ.Context) error { return nil },
	)
	reqOr = mock.Input(mock.URL("/user/123456"))
)

//go:noinline
func aa() bool {
	return false
}

//go:noinline
func call() {
	for i := 0; i < 3000; i++ {
		aa()
	}
}

func BenchmarkCall(mb *testing.B) {
	for i := 0; i < mb.N; i++ {
		call()
	}
}

/*

TODO
 - botleneck is O(n) function calls
 - convert issue to O(log(N)) using trie of method + path segments

*/

func BenchmarkOr(mb *testing.B) {
	// root := &µ.Node{}

	// µ.JoinN(
	// 	µ.Method2("GET"),
	// 	µ.Path2("a"),
	// )(root)

	// µ.JoinN(
	// 	µ.Method2("GET"),
	// 	µ.Path2("b"),
	// )(root)

	// µ.JoinN(
	// 	µ.Method2("GET"),
	// 	µ.Path2("c", "1"),
	// )(root)

	// µ.JoinN(
	// 	µ.Method2("GET"),
	// 	µ.Path2("c", "2"),
	// )(root)

	// µ.JoinN(
	// 	µ.Method2("GET"),
	// 	µ.Path2("c", "3", "d"),
	// )(root)

	// router := httpd.Serve(root.EvalRoot) // seq...,
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathA),
	// µ.GET(pathB),

	// w := new(mockResponseWriter)
	// r, _ := http.NewRequest("GET", "/c/3/d", nil)

	// mb.ReportAllocs()
	// mb.ResetTimer()

	// for i := 0; i < mb.N; i++ {
	// 	router.ServeHTTP(w, r)
	// }

	// mb.ReportAllocs()
	// mb.ResetTimer()

	// for i := 0; i < mb.N; i++ {
	// 	apiOr(reqOr)
	// }
}

/*

Path Pattern with 1 param

*/

type MyT1 struct {
	Name string
}

var (
	name  = optics.ForProduct1(MyT1{})
	path1 = µ.Path("user", name)
	foo1  = µ.GET(path1)
	req1  = mock.Input(mock.URL("/user/123456"))
)

func BenchmarkPathParam1(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo1(req1)
	}
}

func BenchmarkServerParam1(mb *testing.B) {
	w := new(mockResponseWriter)
	router := httpd.Serve(
		µ.Join(
			foo1,
			func(c *µ.Context) error { return nil },
		),
	)
	r, _ := http.NewRequest("GET", "/user/123456", nil)

	mb.ReportAllocs()
	mb.ResetTimer()
	for i := 0; i < mb.N; i++ {
		router.ServeHTTP(w, r)
	}
}

/*

Path Pattern with 5 param

*/

type MyT5 struct{ A, B, C, D, E string }

var (
	a, b, c, d, e = optics.ForProduct5(MyT5{})
	path5         = µ.Path(a, b, c, d, e)
	foo5          = µ.GET(path5)
	req5          = mock.Input(mock.URL("/a/b/c/d/e"))
)

func BenchmarkPathParam5(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo5(req5)
	}
}

func BenchmarkServerParam5(mb *testing.B) {
	w := new(mockResponseWriter)
	router := httpd.Serve(
		µ.Join(
			foo5,
			func(c *µ.Context) error { return nil },
		),
	)
	r, _ := http.NewRequest("GET", "/a/b/c/d/e", nil)

	mb.ReportAllocs()
	mb.ResetTimer()
	for i := 0; i < mb.N; i++ {
		router.ServeHTTP(w, r)
	}
}

/*

Lens decode with 1 param

*/

func BenchmarkLensForProduct1(mb *testing.B) {
	ctx := µ.NewContext(context.Background())
	ctx.Put(name, "123456")

	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		ctx.Get(&val)
	}
}

/*

Lens decode with 1 param

*/

func BenchmarkLensForProduct5(mb *testing.B) {
	ctx := µ.NewContext(context.Background())
	ctx.Put(a, "a")
	ctx.Put(b, "b")
	ctx.Put(c, "c")
	ctx.Put(d, "d")
	ctx.Put(e, "e")

	var val MyT5

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		ctx.Get(&val)
	}
}

/*

Endpoint decode with 1 param

*/

var endpoint1 = µ.GET(
	path1,
	µ.FMap(func(ctx *µ.Context) error {
		var req MyT1
		if err := ctx.Get(&req); err != nil {
			return µ.Status.BadRequest(µ.WithIssue(err))
		}

		return µ.Status.OK(
			headers.ContentType.Value(headers.TextPlain),
			headers.Server.Value("echo"),
			µ.WithText(req.Name),
		)
	}),
)

func BenchmarkEndpoint1(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		endpoint1(req1)
	}
}

/*

Endpoint decode with 5 param

*/

var endpoint5 = µ.GET(
	path5,
	µ.FMap(func(ctx *µ.Context) error {
		var req MyT5
		if err := ctx.Get(&req); err != nil {
			return µ.Status.BadRequest(µ.WithIssue(err))
		}

		return µ.Status.OK(
			headers.ContentType.Value(headers.TextPlain),
			headers.Server.Value("echo"),
			µ.WithText(filepath.Join(req.A, req.B, req.C, req.D, req.E)),
		)
	}),
)

func BenchmarkEndpoint5(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		endpoint5(req5)
	}
}

/*

Simulation of GitHub API, the test is borrowed from
https://github.com/julienschmidt/go-http-routing-benchmark

*/
// http://developer.github.com/v3/
var githubAPI = []struct{ method, path string }{
	// OAuth Authorizations
	{"GET", "/authorizations"},
	{"GET", "/authorizations/:id"},
	{"POST", "/authorizations"},
	//{"PUT", "/authorizations/clients/:client_id"},
	//{"PATCH", "/authorizations/:id"},
	{"DELETE", "/authorizations/:id"},
	{"GET", "/applications/:client_id/tokens/:access_token"},
	{"DELETE", "/applications/:client_id/tokens"},
	{"DELETE", "/applications/:client_id/tokens/:access_token"},

	// Activity
	{"GET", "/events"},
	{"GET", "/repos/:owner/:repo/events"},
	{"GET", "/networks/:owner/:repo/events"},
	{"GET", "/orgs/:org/events"},
	{"GET", "/users/:user/received_events"},
	{"GET", "/users/:user/received_events/public"},
	{"GET", "/users/:user/events"},
	{"GET", "/users/:user/events/public"},
	{"GET", "/users/:user/events/orgs/:org"},
	{"GET", "/feeds"},
	{"GET", "/notifications"},
	{"GET", "/repos/:owner/:repo/notifications"},
	{"PUT", "/notifications"},
	{"PUT", "/repos/:owner/:repo/notifications"},
	{"GET", "/notifications/threads/:id"},
	//{"PATCH", "/notifications/threads/:id"},
	{"GET", "/notifications/threads/:id/subscription"},
	{"PUT", "/notifications/threads/:id/subscription"},
	{"DELETE", "/notifications/threads/:id/subscription"},
	{"GET", "/repos/:owner/:repo/stargazers"},
	{"GET", "/users/:user/starred"},
	{"GET", "/user/starred"},
	{"GET", "/user/starred/:owner/:repo"},
	{"PUT", "/user/starred/:owner/:repo"},
	{"DELETE", "/user/starred/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/subscribers"},
	{"GET", "/users/:user/subscriptions"},
	{"GET", "/user/subscriptions"},
	{"GET", "/repos/:owner/:repo/subscription"},
	{"PUT", "/repos/:owner/:repo/subscription"},
	{"DELETE", "/repos/:owner/:repo/subscription"},
	{"GET", "/user/subscriptions/:owner/:repo"},
	{"PUT", "/user/subscriptions/:owner/:repo"},
	{"DELETE", "/user/subscriptions/:owner/:repo"},

	// Gists
	{"GET", "/users/:user/gists"},
	{"GET", "/gists"},
	//{"GET", "/gists/public"},
	//{"GET", "/gists/starred"},
	{"GET", "/gists/:id"},
	{"POST", "/gists"},
	//{"PATCH", "/gists/:id"},
	{"PUT", "/gists/:id/star"},
	{"DELETE", "/gists/:id/star"},
	{"GET", "/gists/:id/star"},
	{"POST", "/gists/:id/forks"},
	{"DELETE", "/gists/:id"},

	// Git Data
	{"GET", "/repos/:owner/:repo/git/blobs/:sha"},
	{"POST", "/repos/:owner/:repo/git/blobs"},
	{"GET", "/repos/:owner/:repo/git/commits/:sha"},
	{"POST", "/repos/:owner/:repo/git/commits"},
	//{"GET", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/refs"},
	{"POST", "/repos/:owner/:repo/git/refs"},
	//{"PATCH", "/repos/:owner/:repo/git/refs/*ref"},
	//{"DELETE", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/tags/:sha"},
	{"POST", "/repos/:owner/:repo/git/tags"},
	{"GET", "/repos/:owner/:repo/git/trees/:sha"},
	{"POST", "/repos/:owner/:repo/git/trees"},

	// Issues
	{"GET", "/issues"},
	{"GET", "/user/issues"},
	{"GET", "/orgs/:org/issues"},
	{"GET", "/repos/:owner/:repo/issues"},
	{"GET", "/repos/:owner/:repo/issues/:number"},
	{"POST", "/repos/:owner/:repo/issues"},
	//{"PATCH", "/repos/:owner/:repo/issues/:number"},
	{"GET", "/repos/:owner/:repo/assignees"},
	{"GET", "/repos/:owner/:repo/assignees/:assignee"},
	{"GET", "/repos/:owner/:repo/issues/:number/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments/:id"},
	{"POST", "/repos/:owner/:repo/issues/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/issues/comments/:id"},
	//{"DELETE", "/repos/:owner/:repo/issues/comments/:id"},
	{"GET", "/repos/:owner/:repo/issues/:number/events"},
	//{"GET", "/repos/:owner/:repo/issues/events"},
	//{"GET", "/repos/:owner/:repo/issues/events/:id"},
	{"GET", "/repos/:owner/:repo/labels"},
	{"GET", "/repos/:owner/:repo/labels/:name"},
	{"POST", "/repos/:owner/:repo/labels"},
	//{"PATCH", "/repos/:owner/:repo/labels/:name"},
	{"DELETE", "/repos/:owner/:repo/labels/:name"},
	{"GET", "/repos/:owner/:repo/issues/:number/labels"},
	{"POST", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name"},
	{"PUT", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones"},
	{"GET", "/repos/:owner/:repo/milestones/:number"},
	{"POST", "/repos/:owner/:repo/milestones"},
	//{"PATCH", "/repos/:owner/:repo/milestones/:number"},
	{"DELETE", "/repos/:owner/:repo/milestones/:number"},

	// // Miscellaneous
	{"GET", "/emojis"},
	{"GET", "/gitignore/templates"},
	{"GET", "/gitignore/templates/:name"},
	{"POST", "/markdown"},
	{"POST", "/markdown/raw"},
	{"GET", "/meta"},
	{"GET", "/rate_limit"},

	// Organizations
	{"GET", "/users/:user/orgs"},
	{"GET", "/user/orgs"},
	{"GET", "/orgs/:org"},
	//{"PATCH", "/orgs/:org"},
	{"GET", "/orgs/:org/members"},
	{"GET", "/orgs/:org/members/:user"},
	{"DELETE", "/orgs/:org/members/:user"},
	{"GET", "/orgs/:org/public_members"},
	{"GET", "/orgs/:org/public_members/:user"},
	{"PUT", "/orgs/:org/public_members/:user"},
	{"DELETE", "/orgs/:org/public_members/:user"},
	{"GET", "/orgs/:org/teams"},
	{"GET", "/teams/:id"},
	{"POST", "/orgs/:org/teams"},
	//{"PATCH", "/teams/:id"},
	{"DELETE", "/teams/:id"},
	{"GET", "/teams/:id/members"},
	{"GET", "/teams/:id/members/:user"},
	{"PUT", "/teams/:id/members/:user"},
	{"DELETE", "/teams/:id/members/:user"},
	{"GET", "/teams/:id/repos"},
	{"GET", "/teams/:id/repos/:owner/:repo"},
	{"PUT", "/teams/:id/repos/:owner/:repo"},
	{"DELETE", "/teams/:id/repos/:owner/:repo"},
	{"GET", "/user/teams"},

	// Pull Requests
	{"GET", "/repos/:owner/:repo/pulls"},
	{"GET", "/repos/:owner/:repo/pulls/:number"},
	{"POST", "/repos/:owner/:repo/pulls"},
	//{"PATCH", "/repos/:owner/:repo/pulls/:number"},
	{"GET", "/repos/:owner/:repo/pulls/:number/commits"},
	{"GET", "/repos/:owner/:repo/pulls/:number/files"},
	{"GET", "/repos/:owner/:repo/pulls/:number/merge"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/merge"},
	{"GET", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments/:number"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/pulls/comments/:number"},
	//{"DELETE", "/repos/:owner/:repo/pulls/comments/:number"},

	// Repositories
	{"GET", "/user/repos"},
	{"GET", "/users/:user/repos"},
	{"GET", "/orgs/:org/repos"},
	{"GET", "/repositories"},
	{"POST", "/user/repos"},
	{"POST", "/orgs/:org/repos"},
	{"GET", "/repos/:owner/:repo"},
	//{"PATCH", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/contributors"},
	{"GET", "/repos/:owner/:repo/languages"},
	{"GET", "/repos/:owner/:repo/teams"},
	{"GET", "/repos/:owner/:repo/tags"},
	{"GET", "/repos/:owner/:repo/branches"},
	{"GET", "/repos/:owner/:repo/branches/:branch"},
	{"DELETE", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/collaborators"},
	{"GET", "/repos/:owner/:repo/collaborators/:user"},
	{"PUT", "/repos/:owner/:repo/collaborators/:user"},
	{"DELETE", "/repos/:owner/:repo/collaborators/:user"},
	{"GET", "/repos/:owner/:repo/comments"},
	{"GET", "/repos/:owner/:repo/commits/:sha/comments"},
	{"POST", "/repos/:owner/:repo/commits/:sha/comments"},
	{"GET", "/repos/:owner/:repo/comments/:id"},
	//{"PATCH", "/repos/:owner/:repo/comments/:id"},
	{"DELETE", "/repos/:owner/:repo/comments/:id"},
	{"GET", "/repos/:owner/:repo/commits"},
	{"GET", "/repos/:owner/:repo/commits/:sha"},
	{"GET", "/repos/:owner/:repo/readme"},
	//{"GET", "/repos/:owner/:repo/contents/*path"},
	//{"PUT", "/repos/:owner/:repo/contents/*path"},
	//{"DELETE", "/repos/:owner/:repo/contents/*path"},
	//{"GET", "/repos/:owner/:repo/:archive_format/:ref"},
	{"GET", "/repos/:owner/:repo/keys"},
	{"GET", "/repos/:owner/:repo/keys/:id"},
	{"POST", "/repos/:owner/:repo/keys"},
	//{"PATCH", "/repos/:owner/:repo/keys/:id"},
	{"DELETE", "/repos/:owner/:repo/keys/:id"},
	{"GET", "/repos/:owner/:repo/downloads"},
	{"GET", "/repos/:owner/:repo/downloads/:id"},
	{"DELETE", "/repos/:owner/:repo/downloads/:id"},
	{"GET", "/repos/:owner/:repo/forks"},
	{"POST", "/repos/:owner/:repo/forks"},
	{"GET", "/repos/:owner/:repo/hooks"},
	{"GET", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks"},
	//{"PATCH", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks/:id/tests"},
	{"DELETE", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/merges"},
	{"GET", "/repos/:owner/:repo/releases"},
	{"GET", "/repos/:owner/:repo/releases/:id"},
	{"POST", "/repos/:owner/:repo/releases"},
	//{"PATCH", "/repos/:owner/:repo/releases/:id"},
	{"DELETE", "/repos/:owner/:repo/releases/:id"},
	{"GET", "/repos/:owner/:repo/releases/:id/assets"},
	{"GET", "/repos/:owner/:repo/stats/contributors"},
	{"GET", "/repos/:owner/:repo/stats/commit_activity"},
	{"GET", "/repos/:owner/:repo/stats/code_frequency"},
	{"GET", "/repos/:owner/:repo/stats/participation"},
	{"GET", "/repos/:owner/:repo/stats/punch_card"},
	{"GET", "/repos/:owner/:repo/statuses/:ref"},
	{"POST", "/repos/:owner/:repo/statuses/:ref"},

	// Search
	{"GET", "/search/repositories"},
	{"GET", "/search/code"},
	{"GET", "/search/issues"},
	{"GET", "/search/users"},
	{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword"},
	{"GET", "/legacy/repos/search/:keyword"},
	{"GET", "/legacy/user/search/:keyword"},
	{"GET", "/legacy/user/email/:email"},

	// Users
	{"GET", "/users/:user"},
	{"GET", "/user"},
	//{"PATCH", "/user"},
	{"GET", "/users"},
	{"GET", "/user/emails"},
	{"POST", "/user/emails"},
	{"DELETE", "/user/emails"},
	{"GET", "/users/:user/followers"},
	{"GET", "/user/followers"},
	{"GET", "/users/:user/following"},
	{"GET", "/user/following"},
	{"GET", "/user/following/:user"},
	{"GET", "/users/:user/following/:target_user"},
	{"PUT", "/user/following/:user"},
	{"DELETE", "/user/following/:user"},
	{"GET", "/users/:user/keys"},
	{"GET", "/user/keys"},
	{"GET", "/user/keys/:id"},
	{"POST", "/user/keys"},
	//{"PATCH", "/user/keys/:id"},
	{"DELETE", "/user/keys/:id"},
}

var github http.Handler = loadRouter(githubAPI)

type githubReq struct {
	V0, V1, V2, V3, V4, V5, V6, V7, V8, V9 string
}

var v0, v1, v2, v3, v4, v5, v6, v7, v8, v9 = optics.ForProduct10(githubReq{})

func githubHandle(c *µ.Context) error { return nil }

func loadRouter(routes []struct{ method, path string }) http.Handler {
	root := &µ.Node{}
	// seq := make([]µ.Builder, 0, len(routes))
	for _, ep := range routes {
		if ep.method != "GET" {
			continue
		}

		// lens := []interface{}{v0, v1, v2, v3, v4, v5, v6, v7, v8, v9}
		segs := []string{}
		path := strings.Split(ep.path, "/")[1:]
		for _, seg := range path {
			switch {
			case len(seg) == 0:
				break
			case seg[0] == ':':
				segs = append(segs, ":")
			default:
				segs = append(segs, seg)
				// case seg[0] != ':':
				// 	segs = append(segs, seg)
				// 	continue
				// default:
				// 	segs = append(segs, lens[0])
				// 	lens = lens[1:]
			}
		}
		// seq = append(seq,
		µ.JoinN(
			// func(c *µ.Context) error { return µ.ErrNoMatch },
			// func(c *µ.Context) error { return µ.ErrNoMatch },
			// func(c *µ.Context) error { return nil },
			//
			// µ.Method2(ep.method),
			// µ.Method(ep.method),
			µ.Path2(segs...),
			// githubHandle,
		)(root)
		// )
	}

	µ.Walk(root,
		func(i int, n *µ.Node) {
			fmt.Println(strings.Repeat(" ", i), n.Path)
		},
	)

	// root.Print(0)
	// panic('x')

	return httpd.Serve(root.EvalRoot)
}

func benchRoutes(b *testing.B, router http.Handler, routes []struct{ method, path string }) {
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, route := range routes {
			if route.method == "GET" {
				r.Method = route.method
				r.RequestURI = route.path
				u.Path = route.path
				u.RawQuery = rq
				router.ServeHTTP(w, r)
			}
		}
	}
}

func BenchmarkGitHub(b *testing.B) {
	benchRoutes(b, github, githubAPI)
}
