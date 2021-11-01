package optics_test

import (
	"testing"

	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/it"
)

func TestLensStructString(t *testing.T) {
	type T struct {
		A string
		B *string
	}
	a, b := optics.Lenses2(T{})
	x, y := "a", "b"

	t.Run("ByVal", func(t *testing.T) {
		var z T
		m := optics.Morphism{a: x, b: y}
		e := m.Apply(&z)

		it.Ok(t).
			IfNil(e).
			If(z.A).Equal("a").
			If(*z.B).Equal("b")
	})

	t.Run("ByPtr", func(t *testing.T) {
		var z T
		m := optics.Morphism{a: &x, b: &y}
		e := m.Apply(&z)

		it.Ok(t).
			IfNil(e).
			If(z.A).Equal("a").
			If(*z.B).Equal("b").
			If(z.B).Equal(&y)
	})
}

func TestLensStructInt(t *testing.T) {
	type T struct{ A int }
	a := optics.Lenses1(T{})
	m := optics.Morphism{a: 100}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal(100)
}

func TestLensStructFloat(t *testing.T) {
	type T struct{ A float64 }
	a := optics.Lenses1(T{})
	m := optics.Morphism{a: 100.0}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal(100.0)
}

func TestLensStructJSON(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct{ A J }
	a := optics.Lenses1(T{})
	m := optics.Morphism{a: "{\"x\":\"abc\"}"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A.X).Equal("abc")
}

func TestLensStructSeq(t *testing.T) {
	type T struct{ A []string }
	a := optics.Lenses1(T{})
	m := optics.Morphism{a: []string{"a", "b", "c"}}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal([]string{"a", "b", "c"})
}

func TestLenses1(t *testing.T) {
	type T struct {
		A string
	}
	a := optics.Lenses1(T{})
	m := optics.Morphism{a: "a"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a")
}

func TestLenses2(t *testing.T) {
	type T struct {
		A string
		B string
	}
	a, b := optics.Lenses2(T{})
	m := optics.Morphism{a: "a", b: "b"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b")
}

func TestLenses3(t *testing.T) {
	type T struct {
		A string
		B string
		C string
	}
	a, b, c := optics.Lenses3(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c")
}

func TestLenses4(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
	}
	a, b, c, d := optics.Lenses4(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d")
}

func TestLenses5(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
	}
	a, b, c, d, e := optics.Lenses5(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e")
}

func TestLenses6(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
		F string
	}
	a, b, c, d, e, f := optics.Lenses6(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e").
		If(x.F).Equal("f")
}

func TestLenses7(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
	}
	a, b, c, d, e, f, g := optics.Lenses7(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f", g: "g"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e").
		If(x.F).Equal("f").
		If(x.G).Equal("g")
}

func TestLenses8(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
	}
	a, b, c, d, e, f, g, h := optics.Lenses8(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f", g: "g", h: "h"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e").
		If(x.F).Equal("f").
		If(x.G).Equal("g").
		If(x.H).Equal("h")
}

func TestLenses9(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
		I string
	}
	a, b, c, d, e, f, g, h, i := optics.Lenses9(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f", g: "g", h: "h", i: "i"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e").
		If(x.F).Equal("f").
		If(x.G).Equal("g").
		If(x.H).Equal("h").
		If(x.I).Equal("i")
}

func TestLenses10(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
		I string
		K string
	}
	a, b, c, d, e, f, g, h, i, k := optics.Lenses10(T{})
	m := optics.Morphism{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f", g: "g", h: "h", i: "i", k: "k"}

	var x T
	err := m.Apply(&x)

	it.Ok(t).
		IfNil(err).
		If(x.A).Equal("a").
		If(x.B).Equal("b").
		If(x.C).Equal("c").
		If(x.D).Equal("d").
		If(x.E).Equal("e").
		If(x.F).Equal("f").
		If(x.G).Equal("g").
		If(x.H).Equal("h").
		If(x.I).Equal("i").
		If(x.K).Equal("k")
}
