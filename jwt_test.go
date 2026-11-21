package gouldian_test

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	"github.com/fogfish/it/v2"
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

	it.Then(t).
		Should(it.Equal(foo(success), nil)).
		ShouldNot(it.Equal(foo(failure1), nil)).
		ShouldNot(it.Equal(foo(failure2), nil))
}

func TestJWTVar(t *testing.T) {
	type MyT struct{ Sub string }
	sub := µ.Optics1[MyT, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.JWT(µ.Token.Sub, sub),
		),
	)

	t.Run("some", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.Token{"sub": "sub"}))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Sub, "sub"))
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
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

	it.Then(t).
		Should(it.Equal(foo(success), nil)).
		ShouldNot(it.Equal(foo(failure1), nil)).
		ShouldNot(it.Equal(foo(failure2), nil))
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

	it.Then(t).
		Should(it.Equal(foo(success), nil)).
		ShouldNot(it.Equal(foo(failure1), nil)).
		ShouldNot(it.Equal(foo(failure2), nil))
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

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Sub, "sub"))
	})

	t.Run("empty", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.Token{}))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Sub, ""))
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})

}
