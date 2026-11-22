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
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	µ "github.com/fogfish/gouldian/v3"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/gouldian/v3/server/httpd"
	"github.com/fogfish/it/v2"
)

type principalKey struct{}

func serve(t *testing.T, h http.Handler, method, target string) (int, http.Header, string) {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("X-Request-Id", "rid-1")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	res := w.Result()
	buf, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header, string(buf)
}

func things() µ.Routable {
	return µ.Named("things.list", µ.GET(µ.URI(µ.Path("things")), func(c *µ.Context) error {
		p, _ := c.Value(principalKey{}).(string)
		return ø.Status.OK(ø.ContentType.TextPlain, ø.Send(p))
	}))
}

func TestGuard(t *testing.T) {
	admit := func(c *µ.Context) error {
		switch c.Request.Header.Get("Authorization") {
		case "":
			return &µ.Failure{Status: http.StatusUnauthorized, Err: errors.New("no key")}
		case "deny":
			return ø.Status.Forbidden()
		}
		c.Set(principalKey{}, c.Request.Header.Get("Authorization"))
		return nil
	}

	h := httpd.Handler(µ.NewServer(µ.WithGuard(admit)).Mount(things()))
	do := func(target, auth string) (int, string) {
		req := httptest.NewRequest("GET", target, nil)
		if auth != "" {
			req.Header.Set("Authorization", auth)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		buf, _ := io.ReadAll(w.Result().Body)
		return w.Code, string(buf)
	}

	// the guard answers before the router: an unknown route is not revealed
	code, _ := do("/unknown", "")
	it.Then(t).Should(it.Equal(code, http.StatusUnauthorized))

	code, _ = do("/unknown", "deny")
	it.Then(t).Should(it.Equal(code, http.StatusForbidden))

	code, _ = do("/unknown", "alice")
	it.Then(t).Should(it.Equal(code, http.StatusNotFound))

	// the attribute set by the guard reaches the endpoint
	code, body := do("/things", "alice")
	it.Then(t).
		Should(it.Equal(code, http.StatusOK)).
		Should(it.Equal(body, "alice"))

	// attributes do not leak between pooled contexts
	code, body = do("/things", "bob")
	it.Then(t).Should(it.Equal(code, http.StatusOK)).Should(it.Equal(body, "bob"))
}

func TestGuardsOrder(t *testing.T) {
	var seq []int
	g := func(i int, err error) µ.Guard {
		return func(*µ.Context) error { seq = append(seq, i); return err }
	}

	h := httpd.Handler(µ.NewServer(
		µ.WithGuard(g(1, nil), g(2, ø.Status.Forbidden())),
		µ.WithGuard(g(3, nil)),
	).Mount(things()))

	code, _, _ := serve(t, h, "GET", "/things")
	it.Then(t).
		Should(it.Equal(code, http.StatusForbidden)).
		Should(it.Seq(seq).Equal(1, 2))
}

type problem struct {
	status int
	code   string
}

func (p *problem) Error() string { return p.code }

func TestIssueHandler(t *testing.T) {
	var issues []µ.Issue
	h := httpd.Handler(µ.NewServer(
		µ.WithLogger(slog.New(slog.DiscardHandler)),
		µ.WithIssue(func(µ.Issue, error) (string, []byte) { panic("WithIssueHandler takes precedence") }),
		µ.WithIssueHandler(func(c *µ.Context, issue µ.Issue, cause error, out *µ.Output) {
			issues = append(issues, issue)
			code := "E" + http.StatusText(issue.Status)
			var p *problem
			if errors.As(cause, &p) {
				out.Status, code = p.status, p.code
				out.SetHeader("Retry-After", "1")
			}
			out.SetHeader("Content-Type", "text/plain")
			out.Body = []byte(code + " " + c.RequestID)
		}),
	).Mount(
		µ.GET(µ.URI(µ.Path("busy")), func(*µ.Context) error {
			return &problem{status: http.StatusServiceUnavailable, code: "busy"}
		}),
		µ.GET(µ.URI(µ.Path("ok")), func(*µ.Context) error { return ø.Status.OK() }),
	))

	// a domain error owns the status and the headers of the response
	code, head, body := serve(t, h, "GET", "/busy")
	it.Then(t).
		Should(it.Equal(code, http.StatusServiceUnavailable)).
		Should(it.Equal(head.Get("Retry-After"), "1")).
		Should(it.Equal(body, "busy rid-1"))

	// the responses of the router are rendered by the same handler
	code, _, body = serve(t, h, "GET", "/unknown")
	it.Then(t).
		Should(it.Equal(code, http.StatusNotFound)).
		Should(it.Equal(body, "ENot Found rid-1"))

	code, _, _ = serve(t, h, "GET", "/ok")
	it.Then(t).
		Should(it.Equal(code, http.StatusOK)).
		Should(it.Equal(len(issues), 2))
}

type observed struct {
	route   string
	status  int
	attr    any
	elapsed time.Duration
}

func TestObserver(t *testing.T) {
	var log []observed
	h := httpd.Handler(µ.NewServer(
		µ.WithLogger(slog.New(slog.DiscardHandler)),
		µ.WithRecover(true),
		µ.WithGuard(func(c *µ.Context) error {
			c.Set(principalKey{}, "alice")
			return nil
		}),
		µ.WithObserver(func(c *µ.Context, out *µ.Output, elapsed time.Duration) {
			o := observed{status: out.Status, attr: c.Value(principalKey{}), elapsed: elapsed}
			if r := c.Route(); r != nil {
				o.route = r.Name + " " + r.String()
			}
			log = append(log, o)
		}),
	).Mount(
		things(),
		µ.Named("things.boom", µ.GET(µ.URI(µ.Path("boom")), func(*µ.Context) error { panic("boom") })),
	))

	serve(t, h, "GET", "/things")
	serve(t, h, "GET", "/boom")
	serve(t, h, "GET", "/unknown")
	serve(t, h, "DELETE", "/things")

	it.Then(t).Must(it.Equal(len(log), 4))
	it.Then(t).
		Should(it.Equal(log[0].route, "things.list GET /things")).
		Should(it.Equal(log[0].status, http.StatusOK)).
		Should(it.Equal(log[0].attr, any("alice"))).
		Should(it.True(log[0].elapsed > 0)).
		Should(it.Equal(log[1].route, "things.boom GET /boom")).
		Should(it.Equal(log[1].status, http.StatusInternalServerError)).
		Should(it.Equal(log[2].route, "")).
		Should(it.Equal(log[2].status, http.StatusNotFound)).
		Should(it.Equal(log[3].route, "")).
		Should(it.Equal(log[3].status, http.StatusMethodNotAllowed))
}

func TestRouteOfAlternatives(t *testing.T) {
	var name string
	take := func(c *µ.Context) error { name = c.Route().Name; return nil }

	// two routes share the path and the verb, the one that took the request is reported
	h := httpd.Handler(µ.NewServer(µ.WithLogger(slog.New(slog.DiscardHandler))).Mount(
		µ.Named("json", µ.GET(µ.URI(µ.Path("a")), µ.Accept.JSON, take)),
		µ.Named("text", µ.GET(µ.URI(µ.Path("a")), µ.Accept.TextPlain, take)),
	))

	req := httptest.NewRequest("GET", "/a", nil)
	req.Header.Set("Accept", "text/plain")
	h.ServeHTTP(httptest.NewRecorder(), req)
	it.Then(t).Should(it.Equal(name, "text"))

	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(httptest.NewRecorder(), req)
	it.Then(t).Should(it.Equal(name, "json"))
}

func TestContextValue(t *testing.T) {
	type parentKey struct{}
	type key struct{}

	c := µ.NewContext(context.WithValue(context.Background(), parentKey{}, "parent"))
	c.Set(key{}, 1)
	c.Set(key{}, 2)

	it.Then(t).
		Should(it.Equal(c.Value(key{}), any(2))).
		Should(it.Equal(c.Value(parentKey{}), any("parent"))).
		Should(it.Equal(c.Value("missing"), nil))

	c.Free()
	it.Then(t).Should(it.Equal(c.Value(key{}), nil))
}
