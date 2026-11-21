package conformance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	µ "github.com/fogfish/gouldian/v3"
	ø "github.com/fogfish/gouldian/v3/output"
)

func server(opts ...µ.Option) func(...µ.Routable) func() *µ.Server {
	return func(routes ...µ.Routable) func() *µ.Server {
		return func() *µ.Server { return µ.NewServer(opts...).Mount(routes...) }
	}
}

var plain = server()

func text(s string) µ.Endpoint {
	return func(*µ.Context) error {
		return ø.Status.OK(ø.ContentType.TextPlain, ø.Send(s))
	}
}

func want(t *testing.T, r Resp, status int) {
	t.Helper()
	if r.Status != status {
		t.Fatalf("status %d, want %d, body %q", r.Status, status, r.Body)
	}
}

func problem(t *testing.T, r Resp) µ.Issue {
	t.Helper()
	if ct := r.Header.Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("content-type %q, want problem+json (body %q)", ct, r.Body)
	}
	var issue µ.Issue
	if err := json.Unmarshal(r.Body, &issue); err != nil {
		t.Fatal(err)
	}
	if issue.Status != r.Status {
		t.Fatalf("issue status %d, response %d", issue.Status, r.Status)
	}
	return issue
}

// F1, F2, F7, P1-1: malformed value is 400 naming the param, RFC 9457
func TestMalformedParam(t *testing.T) {
	type q struct{ Limit int }
	limit := µ.Optics1[q, int]()
	srv := plain(µ.GET(µ.URI(µ.Path("items")), µ.ParamMaybe("limit", limit),
		µ.FMap(func(_ *µ.Context, a *q) error { return ø.Status.OK(ø.Send(fmt.Sprint(a.Limit))) })))

	each(t, srv, Req{Method: "GET", Path: "/items", Query: url.Values{"limit": {"abc"}}}, func(t *testing.T, r Resp) {
		want(t, r, 400)
		issue := problem(t, r)
		if issue.Type != "about:blank" || issue.Instance != "/items" {
			t.Fatalf("unexpected issue %+v", issue)
		}
		if len(issue.Errors) != 1 || issue.Errors[0].Pointer != "limit" || issue.Errors[0].Source != "query" {
			t.Fatalf("issue must name the param: %+v", issue)
		}
	})

	each(t, srv, Req{Method: "GET", Path: "/items", Query: url.Values{"limit": {"10"}}}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "10" {
			t.Fatalf("body %q", r.Body)
		}
	})
}

// §6.2: the first decode failure is terminal, a later route does not hide it
func TestFailureIsTerminal(t *testing.T) {
	type q struct{ N int }
	type s struct{ N string }
	n := µ.Optics1[q, int]()
	str := µ.Optics1[s, string]()
	srv := plain(
		µ.GET(µ.URI(µ.Path("x")), µ.Param("n", n), text("int")),
		µ.GET(µ.URI(µ.Path("x")), µ.Param("n", str), text("string")),
	)

	each(t, srv, Req{Method: "GET", Path: "/x", Query: url.Values{"n": {"abc"}}}, func(t *testing.T, r Resp) {
		want(t, r, 400)
	})
}

// path segments are routing, they overload instead of failing
func TestPathOverload(t *testing.T) {
	type a struct{ ID int }
	type b struct{ ID string }
	srv := plain(
		µ.GET(µ.URI(µ.Path("user"), µ.Path(µ.Optics1[a, int]())), text("int")),
		µ.GET(µ.URI(µ.Path("user"), µ.Path(µ.Optics1[b, string]())), text("string")),
	)

	each(t, srv, Req{Method: "GET", Path: "/user/10"}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "int" {
			t.Fatal(string(r.Body))
		}
	})
	each(t, srv, Req{Method: "GET", Path: "/user/alice"}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "string" {
			t.Fatal(string(r.Body))
		}
	})
}

// F2, P4-2: application owns the problem shape
func TestCustomIssue(t *testing.T) {
	srv := server(µ.WithIssue(func(i µ.Issue, cause error) (string, []byte) {
		return "application/vnd.app.error+json", []byte(fmt.Sprintf(`{"code":%d}`, i.Status))
	}))(µ.GET(µ.URI(µ.Path("x")), text("ok")))

	each(t, srv, Req{Method: "GET", Path: "/nope"}, func(t *testing.T, r Resp) {
		want(t, r, 404)
		if r.Header.Get("Content-Type") != "application/vnd.app.error+json" || string(r.Body) != `{"code":404}` {
			t.Fatalf("%v %q", r.Header, r.Body)
		}
	})
}

// F3, P2-2: 404, 405 with Allow, HEAD, OPTIONS
func TestRouting404405(t *testing.T) {
	srv := plain(
		µ.GET(µ.URI(µ.Path("things")), text("list")),
		µ.POST(µ.URI(µ.Path("things")), text("create")),
		µ.GET(µ.URI(µ.Path("things"), µ.Path("one")), text("one")),
	)

	each(t, srv, Req{Method: "GET", Path: "/unknown"}, func(t *testing.T, r Resp) { want(t, r, 404); problem(t, r) })
	each(t, srv, Req{Method: "DELETE", Path: "/things"}, func(t *testing.T, r Resp) {
		want(t, r, 405)
		if got := r.Header.Get("Allow"); got != "GET, HEAD, OPTIONS, POST" {
			t.Fatalf("Allow: %q", got)
		}
	})
	each(t, srv, Req{Method: "HEAD", Path: "/things"}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if len(r.Body) != 0 {
			t.Fatal("HEAD has a body")
		}
	})
	each(t, srv, Req{Method: "OPTIONS", Path: "/things"}, func(t *testing.T, r Resp) {
		want(t, r, 204)
		if got := r.Header.Get("Allow"); got != "GET, HEAD, OPTIONS, POST" {
			t.Fatalf("Allow: %q", got)
		}
	})
	each(t, srv, Req{Method: "OPTIONS", Path: "/missing"}, func(t *testing.T, r Resp) { want(t, r, 404) })
}

// F5, P2-3: PathAny and PathAll through the trie
func TestPathAnyWildcards(t *testing.T) {
	type all struct{ Rest string }
	srv := plain(
		µ.POST(µ.URI(µ.PathAny(), µ.Path("things")), text("post-things")),
		µ.GET(µ.URI(µ.PathAny(), µ.Path("things"), µ.PathAll(µ.Optics1[all, string]())),
			func(c *µ.Context) error {
				var a all
				_ = µ.FromContext(c, &a)
				return ø.Status.OK(ø.Send(a.Rest))
			}),
		µ.GET(µ.URI(µ.PathAny(), µ.Path("schema")), text("schema")),
		µ.POST(µ.URI(µ.PathAny(), µ.Path("query"), µ.Path("find")), text("find")),
	)

	each(t, srv, Req{Method: "POST", Path: "/demo/things"}, func(t *testing.T, r Resp) { want(t, r, 200) })
	each(t, srv, Req{Method: "GET", Path: "/demo/things/a/b"}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "a/b" {
			t.Fatal(string(r.Body))
		}
	})
	each(t, srv, Req{Method: "GET", Path: "/demo/schema"}, func(t *testing.T, r Resp) { want(t, r, 200) })
	each(t, srv, Req{Method: "POST", Path: "/demo/query/find"}, func(t *testing.T, r Resp) { want(t, r, 200) })
}

// shared segments between routes (/a/b, /a/:/c)
func TestSharedSegments(t *testing.T) {
	type s struct{ ID string }
	id := µ.Optics1[s, string]()
	srv := plain(
		µ.GET(µ.URI(µ.Path("a"), µ.Path("b")), text("literal")),
		µ.GET(µ.URI(µ.Path("a"), µ.Path(id), µ.Path("c")), text("var")),
	)

	each(t, srv, Req{Method: "GET", Path: "/a/b"}, func(t *testing.T, r Resp) { want(t, r, 200); _ = r })
	each(t, srv, Req{Method: "GET", Path: "/a/b/c"}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "var" {
			t.Fatal(string(r.Body))
		}
	})
}

// N11, P2-4: trailing slash policies
func TestTrailingSlash(t *testing.T) {
	routes := µ.GET(µ.URI(µ.Path("nil")), text("nil"))

	each(t, server()(routes), Req{Method: "GET", Path: "/nil/"}, func(t *testing.T, r Resp) { want(t, r, 404) })
	each(t, server(µ.WithTrailingSlash(µ.Ignore))(routes), Req{Method: "GET", Path: "/nil/"}, func(t *testing.T, r Resp) { want(t, r, 200) })
	each(t, server(µ.WithTrailingSlash(µ.Redirect))(routes),
		Req{Method: "GET", Path: "/nil/", Query: url.Values{"a": {"b"}}}, func(t *testing.T, r Resp) {
			want(t, r, 308)
			if got := r.Header.Get("Location"); got != "/nil?a=b" {
				t.Fatalf("Location: %q", got)
			}
		})
}

// N1, P1-2: bindings of a rejected route do not leak into the next one
func TestRoutesDoNotLeakBindings(t *testing.T) {
	type s1 struct{ X string }
	type s3 struct{ Y string }
	x := µ.Optics1[s1, string]()
	srv := plain(
		µ.GET(µ.URI(µ.Path("a")), µ.Param("x", x), µ.Header("H", "1"), text("r1")),
		µ.GET(µ.URI(µ.Path("a")), µ.FMap(func(_ *µ.Context, a *s3) error {
			return ø.Status.OK(ø.Send("r2:" + a.Y))
		})),
	)

	each(t, srv, Req{Method: "GET", Path: "/a", Query: url.Values{"x": {"leaked"}}}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "r2:" {
			t.Fatal(string(r.Body))
		}
	})
}

// N2, F4, P3-3: body handling
func TestBody(t *testing.T) {
	type payload struct {
		N int `json:"n"`
	}
	type req struct {
		Raw  []byte
		Data payload `content:"json"`
		Str  payload `content:"json,strict"`
	}
	raw := µ.Optics1[req, []byte]()
	data := µ.Optics1[req, payload]("Data")
	handler := µ.FMap(func(_ *µ.Context, a *req) error {
		return ø.Status.OK(ø.Send(fmt.Sprintf("%s|%d", a.Raw, a.Data.N)))
	})

	// two lenses read the same body, the first route fails on a header
	srv := plain(
		µ.POST(µ.URI(µ.Path("a")), µ.Header("X-Mode", "other"), µ.Body(data), text("first")),
		µ.POST(µ.URI(µ.Path("a")), µ.Body(raw), µ.Body(data), handler),
	)
	each(t, srv, Req{Method: "POST", Path: "/a", Body: []byte(`{"n":7}`)}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != `{"n":7}|7` {
			t.Fatalf("%q", r.Body)
		}
	})

	// empty body
	each(t, srv, Req{Method: "POST", Path: "/a"}, func(t *testing.T, r Resp) {
		want(t, r, 400)
		if !strings.Contains(problem(t, r).Detail, "body is required") {
			t.Fatalf("%s", r.Body)
		}
	})

	// malformed JSON, pointer to the field
	each(t, srv, Req{Method: "POST", Path: "/a", Body: []byte(`{"n":"x"}`)}, func(t *testing.T, r Resp) {
		want(t, r, 400)
		issue := problem(t, r)
		if len(issue.Errors) != 1 || issue.Errors[0].Pointer != "/n" {
			t.Fatalf("%+v", issue)
		}
	})

	// strict
	strict := µ.Optics1[req, payload]("Str")
	srv = plain(µ.POST(µ.URI(µ.Path("s")), µ.Body(strict), text("ok")))
	each(t, srv, Req{Method: "POST", Path: "/s", Body: []byte(`{"n":1,"zzz":2}`)}, func(t *testing.T, r Resp) {
		want(t, r, 400)
		if issue := problem(t, r); len(issue.Errors) != 1 || issue.Errors[0].Pointer != "/zzz" {
			t.Fatalf("%+v", issue)
		}
	})
}

// F4, P3-3: body limit is 413
func TestBodyLimit(t *testing.T) {
	type req struct{ Raw []byte }
	srv := server(µ.WithBodyLimit(8))(µ.POST(µ.URI(µ.Path("a")), µ.Body(µ.Optics1[req, []byte]()), text("ok")))

	each(t, srv, Req{Method: "POST", Path: "/a", Body: []byte("123456789")}, func(t *testing.T, r Resp) { want(t, r, 413) })
	each(t, srv, Req{Method: "POST", Path: "/a", Body: []byte("12345678")}, func(t *testing.T, r Resp) { want(t, r, 200) })
}

// #66, N8, P4-6: binary payload round trip
func TestBinaryRoundTrip(t *testing.T) {
	type req struct{ Raw []byte }
	srv := plain(µ.POST(µ.URI(µ.Path("png")), µ.Body(µ.Optics1[req, []byte]()),
		µ.FMap(func(_ *µ.Context, a *req) error {
			return ø.Status.OK(ø.ContentType.Set("image/png"), ø.Send(a.Raw))
		})))

	png := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0xff, 0xfe, 0x80}
	each(t, srv, Req{Method: "POST", Path: "/png", Body: png}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if !bytes.Equal(r.Body, png) {
			t.Fatalf("% x", r.Body)
		}
	})
}

// F6, P3-4: headers
func TestContentNegotiation(t *testing.T) {
	srv := plain(
		µ.GET(µ.URI(µ.Path("n")), µ.Accept.Is("application/json"), text("json")),
	)
	for accept, status := range map[string]int{
		"application/json":                   200,
		"APPLICATION/JSON":                   200,
		"*/*":                                200,
		"application/*":                      200,
		"text/html, application/json":        200,
		"application/json, text/plain;q=0.9": 200,
		"application/jsonx":                  406,
		"text/html":                          406,
		"application/json;q=0, text/html":    406,
	} {
		each(t, srv, Req{Method: "GET", Path: "/n", Header: http.Header{"Accept": {accept}}}, func(t *testing.T, r Resp) {
			if r.Status != status {
				t.Fatalf("Accept %q: status %d want %d", accept, r.Status, status)
			}
		})
	}

	srv = plain(µ.POST(µ.URI(µ.Path("c")), µ.ContentType.JSON, text("json")))
	for ct, status := range map[string]int{
		"application/json":                200,
		"application/json; charset=utf-8": 200,
		"APPLICATION/JSON":                200,
		"application/jsonx":               415,
		"text/plain":                      415,
	} {
		each(t, srv, Req{Method: "POST", Path: "/c", Header: http.Header{"Content-Type": {ct}}}, func(t *testing.T, r Resp) {
			if r.Status != status {
				t.Fatalf("Content-Type %q: status %d want %d", ct, r.Status, status)
			}
		})
	}

	// both content types on the same path are alternatives
	srv = plain(
		µ.POST(µ.URI(µ.Path("c")), µ.ContentType.JSON, text("json")),
		µ.POST(µ.URI(µ.Path("c")), µ.ContentType.Form, text("form")),
	)
	each(t, srv, Req{Method: "POST", Path: "/c", Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if string(r.Body) != "form" {
			t.Fatal(string(r.Body))
		}
	})
}

// F6: literal header is an exact match
func TestHeaderExactMatch(t *testing.T) {
	srv := plain(µ.GET(µ.URI(µ.Path("h")), µ.Header("X-Mode", "fast"), text("ok")))
	each(t, srv, Req{Method: "GET", Path: "/h", Header: http.Header{"X-Mode": {"FAST"}}}, func(t *testing.T, r Resp) { want(t, r, 200) })
	each(t, srv, Req{Method: "GET", Path: "/h", Header: http.Header{"X-Mode": {"faster"}}}, func(t *testing.T, r Resp) { want(t, r, 404) })
}

// F8, P5-4: CORS is off by default, allow-list, preflight
func TestCORS(t *testing.T) {
	routes := µ.GET(µ.URI(µ.Path("c")), text("ok"))
	origin := http.Header{"Origin": {"https://evil.example"}}

	each(t, server()(routes), Req{Method: "GET", Path: "/c", Header: origin}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if r.Header.Get("Access-Control-Allow-Origin") != "" {
			t.Fatal("CORS must be off by default")
		}
	})

	cors := &µ.CORS{Origins: []string{"https://app.example"}, MaxAge: 10 * time.Minute, Headers: []string{"Authorization"}}
	srv := server(µ.WithCORS(cors))(routes)

	each(t, srv, Req{Method: "GET", Path: "/c", Header: origin}, func(t *testing.T, r Resp) {
		want(t, r, 200)
		if r.Header.Get("Access-Control-Allow-Origin") != "" {
			t.Fatal("disallowed origin got ACAO")
		}
	})
	each(t, srv, Req{Method: "GET", Path: "/c", Header: http.Header{"Origin": {"https://app.example"}}}, func(t *testing.T, r Resp) {
		if r.Header.Get("Access-Control-Allow-Origin") != "https://app.example" || r.Header.Get("Vary") != "Origin" {
			t.Fatalf("%v", r.Header)
		}
	})
	each(t, srv, Req{Method: "OPTIONS", Path: "/c", Header: http.Header{
		"Origin":                         {"https://app.example"},
		"Access-Control-Request-Method":  {"GET"},
		"Access-Control-Request-Headers": {"authorization"},
	}}, func(t *testing.T, r Resp) {
		want(t, r, 204)
		h := r.Header
		if h.Get("Access-Control-Allow-Origin") != "https://app.example" ||
			h.Get("Access-Control-Allow-Methods") != "GET, HEAD, OPTIONS" ||
			h.Get("Access-Control-Allow-Headers") != "Authorization" ||
			h.Get("Access-Control-Max-Age") != "600" {
			t.Fatalf("%v", h)
		}
	})
}

// F8: the headers of the endpoint win over defaults
func TestCORSEndpointHeadersWin(t *testing.T) {
	srv := server(µ.WithCORS(&µ.CORS{Origins: []string{"*"}}))(µ.GET(µ.URI(µ.Path("c")),
		func(*µ.Context) error {
			return ø.Status.OK(ø.Header("Access-Control-Allow-Origin", "https://mine.example"))
		}))
	each(t, srv, Req{Method: "GET", Path: "/c", Header: http.Header{"Origin": {"https://x.example"}}}, func(t *testing.T, r Resp) {
		if got := r.Header.Values("Access-Control-Allow-Origin"); len(got) != 1 || got[0] != "https://mine.example" {
			t.Fatalf("%v", got)
		}
	})
}

// F9, P3-7: multi-value query
func TestMultiValueQuery(t *testing.T) {
	type q struct {
		Rels []string
		Csv  []string `sep:","`
		Flag bool
	}
	rels, csv, flag := µ.Optics3[q, []string, []string, bool]("Rels", "Csv", "Flag")
	srv := plain(µ.GET(µ.URI(µ.Path("q")), µ.Param("rels", rels), µ.ParamMaybe("csv", csv), µ.ParamMaybe("flag", flag),
		µ.FMap(func(_ *µ.Context, a *q) error {
			return ø.Status.OK(ø.Send(fmt.Sprintf("%v|%v|%v", a.Rels, a.Csv, a.Flag)))
		})))

	each(t, srv, Req{Method: "GET", Path: "/q", Query: url.Values{"rels": {"a", "b"}, "csv": {"x,y"}, "flag": {"true"}}},
		func(t *testing.T, r Resp) {
			want(t, r, 200)
			if string(r.Body) != "[a b]|[x y]|true" {
				t.Fatalf("%q", r.Body)
			}
		})
	// P3-1: bool
	each(t, srv, Req{Method: "GET", Path: "/q", Query: url.Values{"rels": {"a"}, "flag": {"maybe"}}},
		func(t *testing.T, r Resp) { want(t, r, 400) })
}

// F11, P3-1: TextUnmarshaler validates itself
type Even int

func (e *Even) UnmarshalText(b []byte) error {
	var n int
	if _, err := fmt.Sscanf(string(b), "%d", &n); err != nil || n%2 != 0 {
		return errors.New("must be an even number")
	}
	*e = Even(n)
	return nil
}

func TestTextUnmarshaler(t *testing.T) {
	type q struct{ N Even }
	srv := plain(µ.GET(µ.URI(µ.Path("e")), µ.Param("n", µ.Optics1[q, Even]()), text("ok")))
	each(t, srv, Req{Method: "GET", Path: "/e", Query: url.Values{"n": {"3"}}}, func(t *testing.T, r Resp) {
		want(t, r, 400)
		if !strings.Contains(problem(t, r).Detail, "must be an even number") {
			t.Fatalf("%s", r.Body)
		}
	})
	each(t, srv, Req{Method: "GET", Path: "/e", Query: url.Values{"n": {"4"}}}, func(t *testing.T, r Resp) { want(t, r, 200) })
}

// #36, N5, P1-5: panic recovery; N6: error and nil results
func TestPanicErrorNil(t *testing.T) {
	srv := plain(
		µ.GET(µ.URI(µ.Path("panic")), func(*µ.Context) error { panic("boom") }),
		µ.GET(µ.URI(µ.Path("error")), func(*µ.Context) error { return errors.New("secret detail") }),
		µ.GET(µ.URI(µ.Path("nil")), func(*µ.Context) error { return nil }),
		µ.GET(µ.URI(µ.Path("failure")), func(*µ.Context) error {
			return ø.Status.Forbidden(ø.Error(errors.New("not yours")))
		}),
	)

	for _, path := range []string{"/panic", "/error"} {
		each(t, srv, Req{Method: "GET", Path: path}, func(t *testing.T, r Resp) {
			want(t, r, 500)
			problem(t, r)
			if strings.Contains(string(r.Body), "boom") || strings.Contains(string(r.Body), "secret") {
				t.Fatalf("cause is exposed: %s", r.Body)
			}
		})
	}
	each(t, srv, Req{Method: "GET", Path: "/nil"}, func(t *testing.T, r Resp) { want(t, r, 204) })
	each(t, srv, Req{Method: "GET", Path: "/failure"}, func(t *testing.T, r Resp) {
		want(t, r, 403)
		if problem(t, r).Detail != "not yours" {
			t.Fatalf("%s", r.Body)
		}
	})

	// recovery can be disabled
	srv = server(µ.WithRecover(false))(µ.GET(µ.URI(µ.Path("panic")), func(*µ.Context) error { panic("boom") }))
	defer func() {
		if recover() == nil {
			t.Fatal("panic must escape if recovery is disabled")
		}
	}()
	runHTTPD(t, srv(), Req{Method: "GET", Path: "/panic"})
}

// N7, P1-4: the context of the caller reaches handlers
func TestContextPropagation(t *testing.T) {
	type key struct{}
	srv := plain(µ.GET(µ.URI(µ.Path("ctx")), func(c *µ.Context) error {
		return ø.Status.OK(ø.Send(fmt.Sprint(c.Value(key{}))))
	}))

	each(t, srv, Req{Method: "GET", Path: "/ctx", Ctx: context.WithValue(context.Background(), key{}, "trace-1")},
		func(t *testing.T, r Resp) {
			if string(r.Body) != "trace-1" {
				t.Fatalf("%q", r.Body)
			}
		})
}

// N9, P4-3: Set-Cookie survives, no duplicate Content-Type
func TestMultiValueResponseHeaders(t *testing.T) {
	srv := plain(µ.GET(µ.URI(µ.Path("c")), func(*µ.Context) error {
		out := µ.NewOutput(200)
		out.AddHeader("Set-Cookie", "a=1")
		out.AddHeader("Set-Cookie", "b=2")
		out.SetHeader("Content-Type", "text/plain")
		out.SetHeader("Content-Type", "text/html")
		return out
	}))

	each(t, srv, Req{Method: "GET", Path: "/c"}, func(t *testing.T, r Resp) {
		if got := r.Header.Values("Set-Cookie"); len(got) != 2 {
			t.Fatalf("Set-Cookie: %v", got)
		}
		if got := r.Header.Values("Content-Type"); len(got) != 1 || got[0] != "text/html" {
			t.Fatalf("Content-Type: %v", got)
		}
	})
}

// N10, P4-5: Map serialization failure is 500
func TestMapSerializationFailure(t *testing.T) {
	type in struct{}
	srv := plain(µ.GET(µ.URI(µ.Path("m")), µ.Map(func(*µ.Context, *in) (*chan int, error) {
		c := make(chan int)
		return &c, nil
	})))
	each(t, srv, Req{Method: "GET", Path: "/m"}, func(t *testing.T, r Resp) { want(t, r, 500) })
}

// N12, #63, P3-6: repeated types bind in order
func TestOrderedOptics(t *testing.T) {
	type acc string
	type thr struct {
		Author acc
		Story  string
		User   acc
	}
	author, story, user := µ.Optics3[thr, acc, string, acc]()
	srv := plain(µ.GET(µ.URI(µ.Path(author), µ.Path(story), µ.Path(user)),
		µ.FMap(func(_ *µ.Context, a *thr) error { return ø.Status.OK(ø.Send(fmt.Sprintf("%+v", *a))) })))

	each(t, srv, Req{Method: "GET", Path: "/alice/s1/bob"}, func(t *testing.T, r Resp) {
		if string(r.Body) != "{Author:alice Story:s1 User:bob}" {
			t.Fatalf("%q", r.Body)
		}
	})
}

// #58, P6-1: inline lens
func TestAt(t *testing.T) {
	type reqCreate struct{ Class string }
	srv := plain(µ.POST(µ.URI(µ.Path("ds"), µ.Path(µ.At[reqCreate]("Class"))),
		µ.FMap(func(_ *µ.Context, a *reqCreate) error { return ø.Status.OK(ø.Send(a.Class)) })))
	each(t, srv, Req{Method: "POST", Path: "/ds/books"}, func(t *testing.T, r Resp) {
		if string(r.Body) != "books" {
			t.Fatalf("%q", r.Body)
		}
	})

	defer func() {
		if recover() == nil {
			t.Fatal("unknown field must panic at start-up")
		}
	}()
	µ.At[reqCreate]("Unknown")
}

// #61, P2-5: URI templates are equivalent to the segment form
func TestURIf(t *testing.T) {
	type req struct {
		ID    string
		Path  string
		Limit int
	}
	id, path, limit := µ.Optics3[req, string, string, int]("ID", "Path", "Limit")
	handler := µ.FMap(func(_ *µ.Context, a *req) error { return ø.Status.OK(ø.Send(fmt.Sprintf("%+v", *a))) })

	tpl := plain(µ.GET(µ.URIf("/users/{id}/files/{path...}{?limit}", id, path, limit), handler))
	seg := plain(µ.GET(µ.URI(µ.Path("users"), µ.Path(id), µ.Path("files"), µ.PathAll(path)), µ.ParamMaybe("limit", limit), handler))

	for _, r := range []Req{
		{Method: "GET", Path: "/users/u1/files/a/b.txt", Query: url.Values{"limit": {"5"}}},
		{Method: "GET", Path: "/users/u1/files/c"},
		{Method: "GET", Path: "/users/u1"},
		{Method: "GET", Path: "/users/u1/files/c", Query: url.Values{"limit": {"x"}}},
	} {
		a := runHTTPD(t, tpl(), r)
		b := runHTTPD(t, seg(), r)
		if a.Status != b.Status || !bytes.Equal(a.Body, b.Body) {
			t.Fatalf("template %d %q != segments %d %q", a.Status, a.Body, b.Status, b.Body)
		}
	}

	for _, bad := range []func(){
		func() { µ.URIf("/users/{id}", id, path) },
		func() { µ.URIf("/users/{id}/{p}", id) },
		func() { µ.URIf("users/{id}", id) },
		func() { µ.URIf("/a/{rest...}/b", id) },
		func() { µ.URIf("/a/x{id}", id) },
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(fmt.Sprint(r), "invalid URIf pattern") {
					t.Fatalf("must panic with the template in the message: %v", r)
				}
			}()
			bad()
		}()
	}
}

// #62, P4-1: instance is the request path
func TestIssueInstance(t *testing.T) {
	srv := plain(µ.GET(µ.URI(µ.Path("a")), text("ok")))
	each(t, srv, Req{Method: "GET", Path: "/missing/thing"}, func(t *testing.T, r Resp) {
		if got := problem(t, r).Instance; got != "/missing/thing" {
			t.Fatalf("instance %q", got)
		}
	})
}

// #60, F9, P4-4: failures are logged, the query string is not
type logCapture struct{ bytes.Buffer }

func TestLogging(t *testing.T) {
	var buf logCapture
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	srv := server(µ.WithLogger(log))(
		µ.GET(µ.URI(µ.Path("error")), func(*µ.Context) error { return errors.New("db is down") }),
	)

	each(t, srv, Req{Method: "GET", Path: "/error", Query: url.Values{"token": {"s3cr3t"}}}, func(t *testing.T, r Resp) {
		want(t, r, 500)
	})
	each(t, srv, Req{Method: "GET", Path: "/nope", Query: url.Values{"token": {"s3cr3t"}}}, func(t *testing.T, r Resp) {
		want(t, r, 404)
	})

	out := buf.String()
	for _, must := range []string{"level=ERROR", "db is down", "path=/error", "status=500", "level=DEBUG", "status=404", "method=GET"} {
		if !strings.Contains(out, must) {
			t.Errorf("log misses %q:\n%s", must, out)
		}
	}
	if strings.Contains(out, "s3cr3t") || strings.Contains(out, "token=") {
		t.Errorf("query string leaked to the log:\n%s", out)
	}
}

// N3, N4: JWT scopes are whole tokens, claims typed
func TestJWTScopes(t *testing.T) {
	srv := plain(µ.GET(µ.URI(µ.Path("s")), µ.JWTOneOf(µ.Token.Scope, "ro"), text("ok")))
	eachAWS(t, srv, Req{Method: "GET", Path: "/s", JWT: map[string]string{"scope": "prod:write"}}, func(t *testing.T, r Resp) { want(t, r, 404) })
	eachAWS(t, srv, Req{Method: "GET", Path: "/s", JWT: map[string]string{"scope": "rw ro"}}, func(t *testing.T, r Resp) { want(t, r, 200) })
}

// F10, P5-3: payload 2.0 specifics, cookie request header
func TestRequestCookies(t *testing.T) {
	srv := plain(µ.GET(µ.URI(µ.Path("c")), µ.Cookie.Any, text("ok")))
	each(t, srv, Req{Method: "GET", Path: "/c", Header: http.Header{"Cookie": {"a=1"}}}, func(t *testing.T, r Resp) { want(t, r, 200) })
}
