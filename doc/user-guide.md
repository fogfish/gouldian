# User Guide

The guide explains how to declare HTTP services with µ-Gouldian v3. Every Go snippet below is copied from [example_test.go](../example_test.go) and executed by `go test`, so the guide never drifts from the library. See [MIGRATION-v3.md](../MIGRATION-v3.md) if you come from v2.

* [Endpoint and its outcomes](#endpoint-and-its-outcomes)
* [Routes](#routes)
* [Binding request to a struct](#binding-request-to-a-struct)
* [Failures and problems](#failures-and-problems)
* [Server options](#server-options)
* [Request lifecycle: guards, issues, observers](#request-lifecycle-guards-issues-observers)
* [Body](#body)
* [Testing through the server](#testing-through-the-server)

## Endpoint and its outcomes

An endpoint is a function `func(*µ.Context) error`. It is composable with `Then` (and-then, product) and `Or` (or-else, coproduct). The returned value tells what happened:

| Outcome        | Meaning                                                               | Server response |
| -------------- | --------------------------------------------------------------------- | --------------- |
| `nil`          | matched, continue the product                                         | `204` if it is the end of the route |
| `µ.NoMatch`    | not my route, `Or` tries the next endpoint (reason code, no allocation) | `404`, `405`, `406`, `415` by the reason |
| `*µ.Failure`   | my route, bad input: malformed value, empty body, body is too large   | `400`, `413`, ... a problem |
| `*µ.Output`    | the response of a handler, never retried                              | as is |
| any `error`    | an unexpected error                                                   | `500`, the cause is logged but never exposed |

`Or` is a Finch-like coproduct: it tries the next branch on `NoMatch` only. The first `Failure` is terminal. A branch that does not match leaves nothing in the context.

A request that is routed but has a malformed value (`?limit=abc`) is therefore not "no match": the client gets `400` that names the parameter. Use `µ.Lenient(endpoint)` to turn a `Failure` into `NoMatch` where the v2 behaviour is required, and `µ.Trace(name, endpoint)` to log (at debug level) why an endpoint did not match.

## Routes

A route is an HTTP verb, a path and endpoints. The verb is routing metadata, so the server distinguishes `404` (no such resource) from `405` (no such method, with the `Allow` header). `HEAD` is derived from `GET`, `OPTIONS` is answered automatically.

```go
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
```

`µ.URI` takes segments: `µ.Path("literal")`, `µ.Path(lens)` captures a segment, `µ.PathAny()` matches any segment, `µ.PathAll(lens)` captures the rest of the path. A path segment that cannot be decoded by its lens does not match the route, so `/user/:id` can be overloaded by an `int` lens and a `string` lens.

`µ.URIf` declares the same route as a string, the syntax is the pattern of Go 1.22 `http.ServeMux` plus the RFC 6570 query form `{?limit,q}`. The pattern is checked when the route is constructed: a mismatch between placeholders and lenses panics at start-up.

```go
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
```

## Binding request to a struct

A handler receives a typed struct and never touches `http.Request`. A **lens** binds a decoded value to a field. `µ.OpticsN[T, A, B, ...](fields...)` builds lenses for fields of a struct. Without names the fields are found by type: the n-th occurrence of a type binds the n-th field of the type (`Optics3[T, User, string, User]` binds the first and the second `User` field). An ambiguous or missing type panics at start-up with a fix-it message. `µ.At[T]("Field")` builds one lens inline.

```go
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
```

```go
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
```

The codec of a lens is chosen once, when the lens is constructed:

1. the `content` tag of the field: `json`, `json,strict`, `form`;
2. a codec registered by the application, `µ.RegisterCodec[T](func(string) (T, error))`;
3. `encoding.TextUnmarshaler`, the domain type validates itself and its message reaches the client;
4. built-in: `string`, `bool`, `int*`, `uint*`, `float*`, `time.Duration`, `[]byte` (raw), slices of scalars (repeated `?a=1&a=2` or a `sep:","` tag), pointers;
5. a struct without a tag is JSON.

An unsupported type panics when the lens is constructed, naming the struct and the field.

A literal header is matched exactly (case-insensitive). `µ.ContentType.JSON` matches the media type (parameters such as charset are ignored), `µ.Accept.JSON` negotiates the `Accept` header (lists, `q` values, `*/*`, `type/*`). A mismatch is `NoMatch` with a reason, the server responds `415` or `406` if no other route matches.

## Failures and problems

Failures are rendered as [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) `application/problem+json`: `type` (`about:blank`), `title`, `status`, `detail`, `instance` (the request path) and `errors` (source, pointer and detail of invalid values). A handler fails with a cause, the server logs the cause. The detail of `4xx` is shown to the client, the cause of `5xx` never is.

```go
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
```

The application owns the shape of the problem: `µ.WithIssue` replaces the renderer for all adapters.

```go
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
```

## Server options

`µ.NewServer(options...).Mount(routes...)` is the single place of the rules above. The adapters (`httpd`, `apigateway`, `apigatewayv2`) only convert their requests and responses.

```go
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
```

| Option                  | Default            |
| ----------------------- | ------------------ |
| `WithLogger(*slog.Logger)` | `slog.Default()`; `5xx` at error level, `4xx` at debug level, method and path only (never the query string) |
| `WithIssue(renderer)`   | RFC 9457           |
| `WithIssueHandler(h)`   | off; owns the whole response of a failure (status, headers, body), takes precedence over `WithIssue` |
| `WithGuard(guards...)`  | none; admit the request before routing |
| `WithObserver(observers...)` | none; see every final response, the route and the time of dispatch |
| `WithCORS(*µ.CORS)`     | off; origin allow-list, automatic preflight, endpoint headers win |
| `WithBodyLimit(n)`      | 6 MiB, over the limit `413` |
| `WithRecover(bool)`     | on; a panic is `500`, the stack is logged |
| `WithTrailingSlash(p)`  | `µ.Strict`; `µ.Redirect` (308), `µ.Ignore` |
| `WithLogRoutes(bool)`   | off; `Mount` logs one info record with the routes: `routes="[GET /users/: POST /users]"` (`:` variable, `_` any segment, `*` rest of the path) |

```text
http.ListenAndServe(":8080", httpd.Handler(api))     // net/http
lambda.Start(apigateway.Handler(api))                // REST API, payload 1.0
lambda.Start(apigatewayv2.Handler(api))              // HTTP API and Function URL, payload 2.0
```

The API Gateway adapters take options: `apigateway.WithToken(f)` (and `apigatewayv2.WithToken(f)`) builds `ctx.JWT` from the event, e.g. from the flat context of a Lambda authorizer. A body that is not valid base64 fails the endpoint that reads it (`400`, rendered by the server), after the guards admitted the request.

The context of the caller (`http.Request`, the Lambda invocation) is propagated: `ctx.Done()`, the deadline and the trace identity are available to handlers. Binary responses are base64 encoded by the API Gateway adapters, multi-value headers and `Set-Cookie` are preserved.

## Request lifecycle: guards, issues, observers

The server runs every request through one pipeline, the application plugs into it without wrapping the adapter:

```text
guards → routes → issue (WithIssue | WithIssueHandler) → log, CORS → observers
```

A guard admits the request before routing. Any outcome other than `nil` is the response, rendered as any other failure; the router is never consulted, so a client who is not admitted cannot probe the routes. A guard passes facts to endpoints and observers with `ctx.Set(key, value)`; `ctx.Value(key)` reads them (and the values of the caller's context). Guards run for CORS preflight too: admit `OPTIONS` if CORS is on.

```go
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
```

`WithIssueHandler` owns the whole response of a failure: the outcomes of endpoints, the rejections of guards and the responses of the router itself (`404`, `405`, ...). Unlike `WithIssue`, it sets the status and headers, so a domain error is simply returned by the handler and mapped in one place.

```go
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
```

An observer sees every request once it is served: the final response, the route that took it (`ctx.Route()`, `nil` if none did) and the time of dispatch. `µ.Named` names the route, so logs and metrics identify the operation without parsing the path.

```go
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
```

## Body

The body is read once and shared by every `µ.Body` lens of the request, so a route may take the raw bytes and the typed struct. An empty body is `400` (`body is required`), a body over the limit is `413`, a malformed body is `400` with the JSON pointer of the invalid value. Use `µ.BodyMaybe` if the body is optional.

```go
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
```

## Testing through the server

Bare endpoints cannot tell routing bugs. Test an API the way a client uses it: through the server and an adapter. `httptest` is enough for `httpd`; the repository [conformance](../conformance) suite runs the same request through all adapters.

```text
api := httpd.Serve(routes...)
w := httptest.NewRecorder()
api.ServeHTTP(w, httptest.NewRequest("GET", "/hello", nil))
```

Each [example](../example) has a test that drives it through `httpd` and `apigatewayv2`.
