package gouldian_test

import (
	"testing"

	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/gouldian/v2/mock"
	"github.com/fogfish/it"
)

func TestLenses1(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]("A")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a))))
	req := mock.Input(mock.URL("/t/a"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a")
}

func TestLenses2(t *testing.T) {
	type T struct{ A, B string }
	a, b := µ.Optics2[T, string, string]("A", "B")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b))))
	req := mock.Input(mock.URL("/t/a/b"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b")
}

func TestLenses3(t *testing.T) {
	type T struct{ A, B, C string }
	a, b, c := µ.Optics3[T, string, string, string]("A", "B", "C")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c))))
	req := mock.Input(mock.URL("/t/a/b/c"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c")
}

func TestLenses4(t *testing.T) {
	type T struct{ A, B, C, D string }
	a, b, c, d := µ.Optics4[T, string, string, string, string]("A", "B", "C", "D")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d))))
	req := mock.Input(mock.URL("/t/a/b/c/d"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d")
}

func TestLenses5(t *testing.T) {
	type T struct{ A, B, C, D, E string }
	a, b, c, d, e := µ.Optics5[T, string, string, string, string, string]("A", "B", "C", "D", "E")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e")
}

func TestLenses6(t *testing.T) {
	type T struct{ A, B, C, D, E, F string }
	a, b, c, d, e, f := µ.Optics6[T, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f")
}

func TestLenses7(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G string }
	a, b, c, d, e, f, g := µ.Optics7[T, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f").
		If(v.G).Equal("g")
}

func TestLenses8(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H string }
	a, b, c, d, e, f, g, h := µ.Optics8[T, string, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G", "H")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g), µ.Path(h))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g/h"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f").
		If(v.G).Equal("g").
		If(v.H).Equal("h")
}

func TestLenses9(t *testing.T) {
	type T struct{ A, B, C, D, E, F, G, H, I string }
	a, b, c, d, e, f, g, h, i := µ.Optics9[T, string, string, string, string, string, string, string, string, string]("A", "B", "C", "D", "E", "F", "G", "H", "I")

	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("t"), µ.Path(a), µ.Path(b), µ.Path(c), µ.Path(d), µ.Path(e), µ.Path(f), µ.Path(g), µ.Path(h), µ.Path(i))))
	req := mock.Input(mock.URL("/t/a/b/c/d/e/f/g/h/i"))

	var v T
	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(µ.FromContext(req, &v)).Should().Equal(nil).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f").
		If(v.G).Equal("g").
		If(v.H).Equal("h").
		If(v.I).Equal("i")
}
