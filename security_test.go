package gouldian_test

import (
	"net/http"
	"testing"

	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/gouldian/v2/mock"
	ø "github.com/fogfish/gouldian/v2/output"
	"github.com/fogfish/it"
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
		it.Ok(t).If(err == nil).Should().Equal(ok)
	}

	for scope, ok := range map[string]bool{
		"ro":         false,
		"prod:write": false,
		"ro rw":      true,
		"rw  ro x":   true,
		"pro prw":    false,
	} {
		err := all(mock.Input(mock.JWT(µ.Token{"scope": scope})))
		it.Ok(t).If(err == nil).Should().Equal(ok)
	}
}

// N4: numeric claims do not panic
func TestNewTokenNonStringClaims(t *testing.T) {
	var tkn µ.Token
	it.Ok(t).If(func() error {
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
	}()).Should().Equal(nil)

	it.Ok(t).
		If(tkn.Sub()).Should().Equal("sub").
		If(tkn.Exp()).Should().Equal("1234").
		If(tkn["iat"]).Should().Equal("1700000000").
		If(tkn["flag"]).Should().Equal("true")

	_, hasObj := tkn["obj"]
	_, hasList := tkn["list"]
	_, hasNil := tkn["empty"]
	it.Ok(t).
		If(hasObj).Should().Equal(false).
		If(hasList).Should().Equal(false).
		If(hasNil).Should().Equal(false)
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
		it.Ok(t).
			If(e(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Y).Should().Equal("v").
			If(val.X).Should().Equal("")
	})

	t.Run("Or", func(t *testing.T) {
		e := µ.Or(
			µ.Param("x", x).Then(µ.Header("H", "1")),
			µ.Param("y", y),
		)

		req := mock.Input(mock.URL("/a?x=leak&y=v"))
		var val sec3
		it.Ok(t).
			If(e(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Y).Should().Equal("v").
			If(val.X).Should().Equal("")
	})

	t.Run("accepted branch keeps bindings", func(t *testing.T) {
		e := µ.Or(
			µ.Param("x", x).Then(µ.Header("H", "1")),
			µ.Param("y", y),
		)

		req := mock.Input(mock.URL("/a?x=a&y=b"), mock.Header("H", "1"))
		var val sec3
		it.Ok(t).
			If(e(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.X).Should().Equal("a").
			If(val.Y).Should().Equal("")
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
	it.Ok(t).If(func() error {
		err = api(req)
		return nil
	}()).Should().Equal(nil) // no panic

	out, ok := err.(*µ.Output)
	it.Ok(t).
		If(ok).Should().Equal(true).
		If(out.Status).Should().Equal(http.StatusAccepted).
		If(got.Y).Should().Equal("v")
}
