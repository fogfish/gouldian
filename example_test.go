package gouldian_test

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	µ "github.com/fogfish/gouldian/v3"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/gouldian/v3/server/httpd"
)

// do sends the request through the server, the way a client does
func do(api http.Handler, method, target, body string, headers ...string) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}

	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)

	res := w.Result()
	buf, _ := io.ReadAll(res.Body)
	fmt.Printf("%d %s %s\n", res.StatusCode, res.Header.Get("Content-Type"), buf)
}

func ExampleGET() {
	hello := µ.GET(
		µ.URI(µ.Path("hello")),
		func(ctx *µ.Context) error {
			return ø.Status.OK(
				ø.ContentType.TextPlain,
				ø.Send("Hello World!"),
			)
		},
	)

	api := httpd.Serve(hello)
	do(api, "GET", "/hello", "")
	// Output:
	// 200 text/plain Hello World!
}

type greeting struct {
	Name  string
	Times int
}

func ExampleFMap() {
	name, times := µ.Optics2[greeting, string, int]("Name", "Times")

	hello := µ.GET(
		µ.URI(µ.Path("hello"), µ.Path(name)),
		µ.ParamMaybe("times", times),
		µ.FMap(func(ctx *µ.Context, req *greeting) error {
			return ø.Status.OK(
				ø.ContentType.TextPlain,
				ø.Send(fmt.Sprintf("Hello %s x%d", req.Name, req.Times)),
			)
		}),
	)

	api := httpd.Serve(hello)
	do(api, "GET", "/hello/Joe?times=3", "")
	// a malformed value is not "no match", it is 400 that names the parameter
	do(api, "GET", "/hello/Joe?times=many", "")
	// Output:
	// 200 text/plain Hello Joe x3
	// 400 application/problem+json {"type":"about:blank","title":"Bad Request","status":400,"detail":"query times: strconv.Atoi: parsing \"many\": invalid syntax","instance":"/hello/Joe","errors":[{"source":"query","pointer":"times","detail":"strconv.Atoi: parsing \"many\": invalid syntax"}]}
}

func ExampleNewServer() {
	api := µ.NewServer(
		µ.WithLogger(slog.New(slog.DiscardHandler)),
		µ.WithBodyLimit(1<<20),
		µ.WithTrailingSlash(µ.Ignore),
	).Mount(
		µ.GET(µ.URI(µ.Path("things")), func(*µ.Context) error { return ø.Status.OK(ø.Send("list")) }),
		µ.POST(µ.URI(µ.Path("things")), func(*µ.Context) error { return ø.Status.Created() }),
	)

	h := httpd.Handler(api)
	do(h, "GET", "/things/", "")
	do(h, "DELETE", "/things", "")
	do(h, "GET", "/unknown", "")
	// Output:
	// 200  list
	// 405 application/problem+json {"type":"about:blank","title":"Method Not Allowed","status":405,"instance":"/things"}
	// 404 application/problem+json {"type":"about:blank","title":"Not Found","status":404,"instance":"/unknown"}
}

func ExampleWithIssue() {
	api := µ.NewServer(
		µ.WithIssue(func(issue µ.Issue, cause error) (string, []byte) {
			return "application/json", []byte(fmt.Sprintf(`{"code":"E%d"}`, issue.Status))
		}),
	)

	do(httpd.Handler(api), "GET", "/", "")
	// Output:
	// 404 application/json {"code":"E404"}
}

func ExampleURIf() {
	type file struct {
		User string
		Path string
	}
	user, path := µ.Optics2[file, string, string]("User", "Path")

	get := µ.GET(
		µ.URIf("/users/{user}/files/{path...}", user, path),
		µ.FMap(func(ctx *µ.Context, req *file) error {
			return ø.Status.OK(ø.Send(req.User + ":" + req.Path))
		}),
	)

	do(httpd.Serve(get), "GET", "/users/joe/files/a/b.txt", "")
	// Output:
	// 200  joe:a/b.txt
}

func ExampleAt() {
	type create struct{ Class string }

	post := µ.POST(
		µ.URI(µ.Path("ds"), µ.Path(µ.At[create]("Class"))),
		µ.FMap(func(ctx *µ.Context, req *create) error {
			return ø.Status.Created(ø.Send(req.Class))
		}),
	)

	do(httpd.Serve(post), "POST", "/ds/books", "")
	// Output:
	// 201  books
}

func ExampleBody() {
	type payload struct {
		Name string `json:"name"`
	}
	type request struct {
		Payload payload `content:"json,strict"`
	}

	post := µ.POST(
		µ.URI(µ.Path("users")),
		µ.ContentType.JSON,
		µ.Body(µ.Optics1[request, payload]()),
		µ.FMap(func(ctx *µ.Context, req *request) error {
			return ø.Status.Created(ø.Send(req.Payload.Name))
		}),
	)

	api := httpd.Serve(post)
	json := []string{"Content-Type", "application/json"}
	do(api, "POST", "/users", `{"name":"joe"}`, json...)
	do(api, "POST", "/users", "", json...)
	do(api, "POST", "/users", `{"name":"joe","unknown":1}`, json...)
	do(api, "POST", "/users", `{"name":"joe"}`, "Content-Type", "text/plain")
	// Output:
	// 201  joe
	// 400 application/problem+json {"type":"about:blank","title":"Bad Request","status":400,"detail":"body is required","instance":"/users"}
	// 400 application/problem+json {"type":"about:blank","title":"Bad Request","status":400,"detail":"body /unknown: json: unknown field \"unknown\"","instance":"/users","errors":[{"source":"body","pointer":"/unknown","detail":"json: unknown field \"unknown\""}]}
	// 415 application/problem+json {"type":"about:blank","title":"Unsupported Media Type","status":415,"instance":"/users"}
}

func ExampleFailure() {
	type request struct{ Limit int }
	limit := µ.Optics1[request, int]()

	get := µ.GET(
		µ.URI(µ.Path("items")),
		µ.ParamMaybe("limit", limit),
		func(ctx *µ.Context) error {
			// the application fails with the cause, the server logs the cause
			// and renders the problem, the cause of 5xx is never exposed
			return ø.Status.Forbidden(ø.Error(fmt.Errorf("items are not yours")))
		},
	)

	do(httpd.Serve(get), "GET", "/items", "")
	// Output:
	// 403 application/problem+json {"type":"about:blank","title":"Forbidden","status":403,"detail":"items are not yours","instance":"/items"}
}

func ExampleWithGuard() {
	type principal struct{}

	api := µ.NewServer(
		// the guard admits the request before routing, a rejected client
		// cannot tell known routes from unknown ones
		µ.WithGuard(func(ctx *µ.Context) error {
			key := ctx.Request.Header.Get("X-Api-Key")
			if key == "" {
				return ø.Status.Unauthorized(ø.Error(fmt.Errorf("api key is required")))
			}
			ctx.Set(principal{}, key)
			return nil
		}),
	).Mount(
		µ.GET(µ.URI(µ.Path("me")), func(ctx *µ.Context) error {
			return ø.Status.OK(ø.ContentType.TextPlain, ø.Send(ctx.Value(principal{}).(string)))
		}),
	)

	h := httpd.Handler(api)
	do(h, "GET", "/me", "", "X-Api-Key", "joe")
	do(h, "GET", "/unknown", "")
	// Output:
	// 200 text/plain joe
	// 401 application/problem+json {"type":"about:blank","title":"Unauthorized","status":401,"detail":"api key is required","instance":"/unknown"}
}

// errBusy is a failure of the domain, it knows its own status
type errBusy struct{ retry time.Duration }

func (e errBusy) Error() string { return "busy" }

func ExampleWithIssueHandler() {
	api := µ.NewServer(
		µ.WithLogger(slog.New(slog.DiscardHandler)),
		// the handler owns the response of every failure: the status, the
		// headers and the body
		µ.WithIssueHandler(func(ctx *µ.Context, issue µ.Issue, cause error, out *µ.Output) {
			code := strings.ReplaceAll(strings.ToLower(issue.Title), " ", "_")
			var busy errBusy
			if errors.As(cause, &busy) {
				out.Status, code = http.StatusServiceUnavailable, "busy"
				out.SetHeader("Retry-After", fmt.Sprint(busy.retry.Seconds()))
			}
			out.SetHeader("Content-Type", "application/json")
			out.Body = []byte(fmt.Sprintf(`{"code":"%s"}`, code))
		}),
	).Mount(
		µ.GET(µ.URI(µ.Path("things")), func(*µ.Context) error {
			return errBusy{retry: 5 * time.Second}
		}),
	)

	h := httpd.Handler(api)
	do(h, "GET", "/things", "")
	do(h, "GET", "/unknown", "")
	// Output:
	// 503 application/json {"code":"busy"}
	// 404 application/json {"code":"not_found"}
}

func ExampleWithObserver() {
	api := µ.NewServer(
		// the observer sees the final response and the route that took it
		µ.WithObserver(func(ctx *µ.Context, out *µ.Output, elapsed time.Duration) {
			route := "-"
			if r := ctx.Route(); r != nil {
				route = r.Name
			}
			fmt.Printf("access: %s %d\n", route, out.Status)
		}),
	).Mount(
		µ.Named("things.list",
			µ.GET(µ.URI(µ.Path("things")), func(*µ.Context) error { return ø.Status.OK(ø.Send("list")) }),
		),
	)

	h := httpd.Handler(api)
	do(h, "GET", "/things", "")
	do(h, "GET", "/unknown", "")
	// Output:
	// access: things.list 200
	// 200  list
	// access: - 404
	// 404 application/problem+json {"type":"about:blank","title":"Not Found","status":404,"instance":"/unknown"}
}
