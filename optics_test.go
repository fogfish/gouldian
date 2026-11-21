package gouldian_test

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	"github.com/fogfish/it/v2"
)

func TestLenses1(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]("A")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a))))
	req := mock.Input(mock.URL("/t/a"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a"))
}

func TestLenses2(t *testing.T) {
	type T struct{ A, B string }
	a, b := µ.Optics2[T, string, string]("A", "B")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b))))
	req := mock.Input(mock.URL("/t/a/b"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b"))
}

func TestLenses3(t *testing.T) {
	type T struct{ A, B, C string }
	a, b, c := µ.Optics3[T, string, string, string]("A", "B", "C")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c))))
	req := mock.Input(mock.URL("/t/a/b/c"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c"))
}

func TestLenses4(t *testing.T) {
	type T struct{ A, B, C, D string }
	a, b, c, d := µ.Optics4[T, string, string, string, string]("A", "B", "C", "D")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d))))
	req := mock.Input(mock.URL("/t/a/b/c/d"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d"))
}

func TestLenses5(t *testing.T) {
	type T struct{ A, B, C, D, E string }
	a, b, c, d, e := µ.Optics5[T, string, string, string, string, string]("A", "B", "C", "D", "E")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d")).
		Should(it.Equal(v.E, "e"))
}

func TestLenses6(t *testing.T) {
	type T struct{ A, B, C, D, E, F string }
	a, b, c, d, e, f := µ.Optics6[T, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d")).
		Should(it.Equal(v.E, "e")).
		Should(it.Equal(v.F, "f"))
}

func TestLenses7(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G string }
	a, b, c, d, e, f, g := µ.Optics7[T, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d")).
		Should(it.Equal(v.E, "e")).
		Should(it.Equal(v.F, "f")).
		Should(it.Equal(v.G, "g"))
}

func TestLenses8(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H string }
	a, b, c, d, e, f, g, h := µ.Optics8[T, string, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G", "H")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g), µ.Path(h))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g/h"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d")).
		Should(it.Equal(v.E, "e")).
		Should(it.Equal(v.F, "f")).
		Should(it.Equal(v.G, "g")).
		Should(it.Equal(v.H, "h"))
}

func TestLenses9(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I string }
	a, b, c, d, e, f, g, h, i := µ.Optics9[T, string, string, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G", "H", "I")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g), µ.Path(h), µ.Path(i))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g/h/i"))

	var v T
	it.Then(t).
		Should(it.Equal(foo(req), nil)).
		Should(it.Equal(µ.FromContext(req, &v), nil)).
		Should(it.Equal(v.A, "a")).
		Should(it.Equal(v.B, "b")).
		Should(it.Equal(v.C, "c")).
		Should(it.Equal(v.D, "d")).
		Should(it.Equal(v.E, "e")).
		Should(it.Equal(v.F, "f")).
		Should(it.Equal(v.G, "g")).
		Should(it.Equal(v.H, "h")).
		Should(it.Equal(v.I, "i"))
}
