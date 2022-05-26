package gouldian_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"

	"github.com/fogfish/it"
)

func TestJWTLit(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWT(µ.Token.Sub, "sub"),
		),
	)
	success := mock.Input(mock.JWT(µ.Token{"sub": "sub"}))
	failure1 := mock.Input(mock.JWT(µ.Token{"sub": "foo"}))
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestJWTVar(t *testing.T) {
	type MyT struct{ Sub string }
	sub := µ.Optics1[myT, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWT(µ.Token.Sub, sub),
		),
	)

	t.Run("some", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.Token{"sub": "sub"}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("sub")
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestJWTOneOf(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWTOneOf(µ.Token.Scope, "a", "b", "c"),
		),
	)

	success := mock.Input(mock.JWT(µ.Token{"scope": "x y c"}))
	failure1 := mock.Input(mock.JWT(µ.Token{"scope": "x y"}))
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestJWTAllOf(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWTAllOf(µ.Token.Scope, "a", "b", "c"),
		),
	)

	success := mock.Input(mock.JWT(µ.Token{"scope": "a b c"}))
	failure1 := mock.Input(mock.JWT(µ.Token{"scope": "a b"}))
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestJWTMaybe(t *testing.T) {
	type MyT struct{ Sub string }
	sub := µ.Optics1[MyT, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWTMaybe(µ.Token.Sub, sub),
		),
	)

	t.Run("some", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.Token{"sub": "sub"}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("sub")
	})

	t.Run("empty", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.Token{}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("")
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

}
