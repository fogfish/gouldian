package gouldian_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/it/v2"
)

func TestEndpointOrReasons(t *testing.T) {
	req := mock.Input(mock.Header("Content-Type", "text/plain"))

	e := µ.ContentType.JSON
	err := µ.Endpoint(e).Or(µ.Header("X-Missing", "1"))(req)
	it.Then(t).Should(it.Equal(err, error(µ.NoMatchContentType)))

	err = µ.Or(µ.Header("X-Missing", "1"), µ.Endpoint(µ.ContentType.JSON))(req)
	it.Then(t).Should(it.Equal(err, error(µ.NoMatchContentType)))
}

func TestFailureIsTerminalInOr(t *testing.T) {
	type q struct{ N int }
	n := µ.Optics1[q, int]()
	req := mock.Input(mock.URL("/?n=x"))

	called := false
	err := µ.Or(
		µ.Param("n", n),
		func(*µ.Context) error { called = true; return nil },
	)(req)

	it.Then(t).
		Should(it.True(isFailure(err))).
		Should(it.Equal(called, false))

	err = µ.Endpoint(µ.Param("n", n)).Or(func(*µ.Context) error { called = true; return nil })(req)
	it.Then(t).
		Should(it.True(isFailure(err))).
		Should(it.Equal(called, false))
}

func TestLenient(t *testing.T) {
	type q struct{ N int }
	n := µ.Optics1[q, int]()
	req := mock.Input(mock.URL("/?n=x"))

	err := µ.Or(
		µ.Lenient(µ.Param("n", n)),
		func(*µ.Context) error { return nil },
	)(req)
	it.Then(t).Should(it.Equal(err, nil))
}

func TestTrace(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	srv := µ.NewServer(µ.WithLogger(log)).Mount(
		µ.GET(µ.URI(µ.Path("a")), µ.Trace("guard", µ.Header("X-A", "1"))),
	)

	c := srv.NewContext(t.Context(), mock.Input(mock.URL("/a")).Request)
	out := srv.Dispatch(c)
	out.Free()
	srv.Release(c)

	it.Then(t).Should(it.True(strings.Contains(buf.String(), "endpoint=guard")))
}

func TestRouterPriorityAndBacktracking(t *testing.T) {
	type s struct{ ID string }
	id := µ.Optics1[s, string]()
	handler := func(s string) µ.Endpoint {
		return func(*µ.Context) error { return ø.Status.OK(ø.Send(s)) }
	}
	srv := µ.NewServer().Mount(
		µ.GET(µ.URI(µ.Path("a"), µ.Path("b")), µ.Header("X-Guard", "1"), handler("literal-guarded")),
		µ.GET(µ.URI(µ.Path("a"), µ.Path(id)), handler("variable")),
		µ.GET(µ.URI(µ.Path("a"), µ.PathAny()), handler("any")),
		µ.GET(µ.URI(µ.Path("a"), µ.PathAll(id)), handler("rest")),
	)

	do := func(path string, headers ...string) string {
		r := mock.Input(mock.URL(path))
		for i := 0; i+1 < len(headers); i += 2 {
			r.Request.Header.Set(headers[i], headers[i+1])
		}
		c := srv.NewContext(t.Context(), r.Request)
		defer srv.Release(c)
		out := srv.Dispatch(c)
		defer out.Free()
		return string(out.Body)
	}

	it.Then(t).
		Should(it.Equal(do("/a/b", "X-Guard", "1"), "literal-guarded")).
		// the guard of the literal route fails, the next candidate is tried
		Should(it.Equal(do("/a/b"), "variable")).
		Should(it.Equal(do("/a/c"), "variable")).
		Should(it.Equal(do("/a/c/d"), "rest"))
}

func TestOpticsOrderedAmbiguity(t *testing.T) {
	type acc string
	type thr struct {
		A acc
		B string
		C acc
	}

	mustPanic := func(name string, f func(), contains string) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: must panic", name)
			}
			if !strings.Contains(r.(error).Error(), contains) {
				t.Fatalf("%s: %v", name, r)
			}
		}()
		f()
	}

	// acc is not unique, but requested once
	mustPanic("not unique", func() { µ.Optics2[thr, acc, string]() }, "is not unique")
	// too many occurrences
	mustPanic("too many", func() { µ.Optics3[thr, string, string, acc]() }, "requested 2 times")
	// type not found
	mustPanic("not found", func() { µ.Optics1[thr, int]() }, "is not field")

	a, b, c := µ.Optics3[thr, acc, string, acc]()
	_, _, _ = a, b, c
}

func TestFailureDetail(t *testing.T) {
	f := &µ.Failure{Source: µ.SourceParam, Key: "limit", Err: errors.New("bad")}
	it.Then(t).
		Should(it.Equal(f.Detail(), "query limit: bad")).
		Should(it.True(errors.Is(f, f.Err)))
}

func TestOutputHeaders(t *testing.T) {
	out := µ.NewOutput(http.StatusOK)
	defer out.Free()

	out.AddHeader("set-cookie", "a=1")
	out.AddHeader("Set-Cookie", "b=2")
	out.SetHeader("Content-Type", "a")
	out.SetHeader("content-type", "b")
	it.Then(t).
		Should(it.Equal(len(out.Headers), 3)).
		Should(it.Equal(out.GetHeader("Content-Type"), "b"))

	out.SetHeader("Set-Cookie", "c=3")
	it.Then(t).Should(it.Equal(len(out.Headers), 2))
}

func TestLiteralOutputIsNotPooled(t *testing.T) {
	out := &µ.Output{Status: 200, Body: []byte("x")}
	out.Free()
	it.Then(t).Should(it.Equal(string(out.Body), "x"))
}

func TestQueryGetAll(t *testing.T) {
	q := µ.Query{"a": {"1", "2"}}
	v, ok := q.GetAll("a")
	it.Then(t).Should(it.Equal(len(v), 2)).Should(it.True(ok))
	first, _ := q.Get("a")
	it.Then(t).Should(it.Equal(first, "1"))
}

func TestTokenAccessors(t *testing.T) {
	tkn := µ.NewToken(map[string]any{
		"s": "x", "f": 1.5, "i": 7, "i64": int64(8), "n": json.Number("9"),
		"b": true, "bs": "true", "arr": []string{"a", "b"}, "anyarr": []any{"a", 1, "b"},
		"num": "12.5",
	})

	n, ok := tkn.Number("i")
	n2, ok2 := tkn.Number("n")
	n3, ok3 := tkn.Number("num")
	_, ok4 := tkn.Number("s")
	b, okb := tkn.Bool("bs")
	_, okc := tkn.Bool("i")
	v, found := tkn.Get("s")
	it.Then(t).
		Should(it.Equal(tkn.String("f"), "1.5")).
		Should(it.Equal(tkn.String("i"), "7")).
		Should(it.Equal(tkn.String("i64"), "8")).
		Should(it.Equal(tkn.String("n"), "9")).
		Should(it.Equal(tkn.String("arr"), "a b")).
		Should(it.Equal(tkn.String("anyarr"), "a b")).
		Should(it.Equal(tkn.String("missing"), "")).
		Should(it.Equal(len(tkn.Strings("arr")), 2)).
		Should(it.Equal(n, 7.0)).Should(it.True(ok)).
		Should(it.Equal(n2, 9.0)).Should(it.True(ok2)).
		Should(it.Equal(n3, 12.5)).Should(it.True(ok3)).
		Should(it.Equal(ok4, false)).
		Should(it.True(b)).Should(it.True(okb)).
		Should(it.Equal(okc, false)).
		Should(it.Equal(v, any("x"))).Should(it.True(found))
}

func TestRegisterCodec(t *testing.T) {
	type Color struct{ V string }
	type q struct{ C Color }
	µ.RegisterCodec(func(s string) (Color, error) { return Color{s}, nil })

	var v q
	req := mock.Input(mock.URL("/?c=red"))
	it.Then(t).
		Should(it.Equal(µ.Param("c", µ.Optics1[q, Color]())(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.C.V, "red"))
}

func TestBodyMaybeAndMethod(t *testing.T) {
	type r struct{ Raw []byte }
	lens := µ.Optics1[r, []byte]()

	it.Then(t).
		Should(it.Equal(µ.BodyMaybe(lens)(mock.Input()), nil)).
		Should(it.Equal(µ.BodyMaybe(lens)(mock.Input(mock.Text("x"))), nil)).
		Should(it.Equal(µ.Method("GET")(mock.Input()), nil)).
		Should(it.Equal(µ.Method("_")(mock.Input()), nil)).
		Should(it.Equal(µ.Method("POST")(mock.Input()), error(µ.NoMatchMethod)))
}

func TestStatusHelpers(t *testing.T) {
	for status, f := range map[int]func(...µ.Result) error{
		206: ø.Status.PartialContent, 207: ø.Status.MultiStatus,
		416: ø.Status.RequestedRangeNotSatisfiable, 417: ø.Status.ExpectationFailed,
		421: ø.Status.MisdirectedRequest, 422: ø.Status.UnprocessableEntity,
		423: ø.Status.Locked, 424: ø.Status.FailedDependency, 426: ø.Status.UpgradeRequired,
		428: ø.Status.PreconditionRequired, 429: ø.Status.TooManyRequests,
		431: ø.Status.RequestHeaderFieldsTooLarge, 451: ø.Status.UnavailableForLegalReasons,
		506: ø.Status.VariantAlsoNegotiates, 507: ø.Status.InsufficientStorage,
		511: ø.Status.NetworkAuthenticationRequired,
	} {
		out := f().(*µ.Output)
		it.Then(t).Should(it.Equal(out.Status, status))
		out.Free()
	}
}

func TestNoMatchStrings(t *testing.T) {
	it.Then(t).Should(it.Equal(µ.ErrNoMatch.Error(), "No Match"))
	for _, e := range []µ.NoMatch{µ.NoMatchPath, µ.NoMatchMethod, µ.NoMatchHeader, µ.NoMatchParam,
		µ.NoMatchJWT, µ.NoMatchBody, µ.NoMatchAccept, µ.NoMatchContentType} {
		it.Then(t).ShouldNot(it.Equal(e.Error(), "No Match"))
	}

	f := &µ.Failure{Source: µ.SourceBody, Err: errors.New("x")}
	it.Then(t).
		Should(it.Equal(f.Error(), "400 body: x")).
		Should(it.Equal(f.Pointer(), ""))
}

func TestLogRoutes(t *testing.T) {
	routes := []µ.Routable{
		µ.GET(µ.URI(µ.Path("users"), µ.Path(µ.Optics1[greeting, string]("Name"))), µ.FMap(func(*µ.Context, *greeting) error { return nil })),
		µ.POST(µ.URI(µ.Path("users")), func(*µ.Context) error { return nil }),
		µ.ANY(µ.URI(µ.Path("files"), µ.PathAll(µ.Optics1[greeting, string]("Name"))), func(*µ.Context) error { return nil }),
		µ.GET(µ.URI(), func(*µ.Context) error { return nil }),
	}

	// off by default
	var off bytes.Buffer
	µ.NewServer(µ.WithLogger(slog.New(slog.NewTextHandler(&off, nil)))).Mount(routes...)
	it.Then(t).Should(it.Equal(off.Len(), 0))

	var on bytes.Buffer
	µ.NewServer(
		µ.WithLogger(slog.New(slog.NewTextHandler(&on, nil))),
		µ.WithLogRoutes(true),
	).Mount(routes...)

	it.Then(t).
		Should(it.True(strings.Contains(on.String(), "level=INFO"))).
		Should(it.True(strings.Contains(on.String(), `msg="routes mounted"`))).
		Should(it.True(strings.Contains(on.String(), "count=4"))).
		Should(it.True(strings.Contains(on.String(), "GET /users/:"))).
		Should(it.True(strings.Contains(on.String(), "POST /users"))).
		Should(it.True(strings.Contains(on.String(), "ANY /files/*"))).
		Should(it.True(strings.Contains(on.String(), "GET /")))
	it.Then(t).Should(it.Equal(strings.Count(strings.TrimSpace(on.String()), "\n"), 0))
}

func TestOpticsNamedFieldTypeIsAsserted(t *testing.T) {
	type T struct {
		Name string
		Age  int
	}

	mustPanic := func(f func(), contains ...string) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("must panic: the type argument does not match the field")
			}
			for _, c := range contains {
				if !strings.Contains(r.(error).Error(), c) {
					t.Fatalf("%v does not contain %q", r, c)
				}
			}
		}()
		f()
	}

	mustPanic(func() { µ.Optics1[T, int]("Name") }, "field `Name` of struct `T` has type `string`, not `int`")
	mustPanic(func() { µ.Optics2[T, string, string]("Name", "Age") }, "field `Age`")

	// matching types and the inferred At are fine
	µ.Optics2[T, string, int]("Name", "Age")
	µ.At[T]("Age")
}
