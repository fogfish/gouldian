package gouldian_test

import (
	"net/http"
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/it/v2"
)

// N3: scope is matched by whole tokens, not by substring
func TestJWTScopeIsNotSubstring(t *testing.T) {
	one := mock.Endpoint(µ.GET(µ.URI(), µ.JWTOneOf(µ.Token.Scope, "ro")))
	all := mock.Endpoint(µ.GET(µ.URI(), µ.JWTAllOf(µ.Token.Scope, "ro", "rw")))

	for scope, ok := range map[string]bool{
		"prod:write":    false,
		"xro":           false,
		"ro:x":          false,
		"":              false,
		"ro":            true,
		"x ro y":        true,
		"x\tro\ny":      true,
		"  ro  ":        true,
		"prod:write ro": true,
	} {
		err := one(mock.Input(mock.JWT(µ.Token{"scope": scope})))
		it.Then(t).
			Should(it.Equal(err == nil, ok))
	}

	for scope, ok := range map[string]bool{
		"ro":         false,
		"prod:write": false,
		"ro rw":      true,
		"rw  ro x":   true,
		"pro prw":    false,
	} {
		err := all(mock.Input(mock.JWT(µ.Token{"scope": scope})))
		it.Then(t).
			Should(it.Equal(err == nil, ok))
	}
}

// N4: numeric claims do not panic
func TestNewTokenNonStringClaims(t *testing.T) {
	var tkn µ.Token
	it.Then(t).
		Should(it.Equal(func() error {
			tkn = µ.NewToken(map[string]any{
				"sub":   "sub",
				"exp":   1234.0,
				"iat":   float64(1700000000),
				"flag":  true,
				"obj":   map[string]any{"a": 1},
				"list":  []any{"a"},
				"empty": nil,
			})
			return nil
		}(), nil))

	it.Then(t).
		Should(it.Equal(tkn.Sub(), "sub")).
		Should(it.Equal(tkn.Exp(), "1234")).
		Should(it.Equal(tkn.String("iat"), "1700000000")).
		Should(it.Equal(tkn.String("flag"), "true"))

	exp, hasExp := tkn.Number("exp")
	flag, hasFlag := tkn.Bool("flag")
	_, hasObj := tkn["obj"]
	_, hasNil := tkn["empty"]
	it.Then(t).
		Should(it.Equal(exp, 1234.0)).
		Should(it.Equal(hasExp, true)).
		Should(it.Equal(flag, true)).
		Should(it.Equal(hasFlag, true)).
		Should(it.Equal(tkn.String("list"), "a")).
		Should(it.Equal(hasObj, true)).
		Should(it.Equal(tkn.String("obj"), "")).
		Should(it.Equal(hasNil, false))
}

type sec1 struct{ X string }
type sec2 struct{ Y string }
type sec3 struct{ X, Y string }

// N1: bindings of a rejected branch are not visible to the next branch.
// Before the fix the value of `x` leaked into the struct of the next branch.
func TestOrDoesNotLeakBindings(t *testing.T) {
	x := µ.Optics1[sec3, string]("X")
	y := µ.Optics1[sec3, string]("Y")

	t.Run("Endpoint.Or", func(t *testing.T) {
		e := µ.Param("x", x).Then(µ.Header("H", "1")).Or(µ.Param("y", y))

		req := mock.Input(mock.URL("/a?x=leak&y=v"))
		var val sec3
		it.Then(t).
			Should(it.Equal(e(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Y, "v")).
			Should(it.Equal(val.X, ""))
	})

	t.Run("Or", func(t *testing.T) {
		e := µ.Or(
			µ.Param("x", x).Then(µ.Header("H", "1")),
			µ.Param("y", y),
		)

		req := mock.Input(mock.URL("/a?x=leak&y=v"))
		var val sec3
		it.Then(t).
			Should(it.Equal(e(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Y, "v")).
			Should(it.Equal(val.X, ""))
	})

	t.Run("accepted branch keeps bindings", func(t *testing.T) {
		e := µ.Or(
			µ.Param("x", x).Then(µ.Header("H", "1")),
			µ.Param("y", y),
		)

		req := mock.Input(mock.URL("/a?x=a&y=b"), mock.Header("H", "1"))
		var val sec3
		it.Then(t).
			Should(it.Equal(e(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.X, "a")).
			Should(it.Equal(val.Y, ""))
	})
}

// N1 through the router: the second route used to panic with
// `invalid type passed to Reflector` on a value from the first route.
func TestRoutesDoNotLeakBindings(t *testing.T) {
	lx := µ.Optics1[sec1, string]()
	ly := µ.Optics1[sec2, string]()

	var got sec2
	api := µ.NewRoutes(
		µ.GET(µ.URI(µ.Path("a")),
			µ.Param("x", lx),
			µ.Header("H", "1"),
			µ.FMap(func(*µ.Context, *sec1) error { return ø.Status.OK() }),
		),
		µ.GET(µ.URI(µ.Path("a")),
			µ.Param("x", ly),
			µ.FMap(func(_ *µ.Context, a *sec2) error {
				got = *a
				return ø.Status.Accepted()
			}),
		),
	).Endpoint()

	req := mock.Input(mock.URL("/a?x=v"))

	var err error
	it.Then(t).
		Should(it.Equal(func() error {
			err = api(req)
			return nil
		}(), nil)) // no panic

	out, ok := err.(*µ.Output)
	it.Then(t).
		Should(it.Equal(ok, true)).
		Should(it.Equal(out.Status, http.StatusAccepted)).
		Should(it.Equal(got.Y, "v"))
}
