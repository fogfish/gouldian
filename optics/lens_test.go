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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")

	t.Run("ByVal", func(t *testing.T) {
		var z T

		m := optics.Morphism{
			{Lens: a, Value: x},
			{Lens: b, Value: y},
		}
		e := m.Apply(&z)

		it.Ok(t).
			IfNil(e).
			If(z.A).Equal("a").
			If(*z.B).Equal("b")
	})
}

func TestLensStructInt(t *testing.T) {
	type T struct{ A int }
	a := optics.Lenses1(T{})
	x, _ := a.FromString("100")
	m := optics.Morphism{{Lens: a, Value: x}}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal(100)
}

func TestLensStructFloat(t *testing.T) {
	type T struct{ A float64 }
	a := optics.Lenses1(T{})
	x, _ := a.FromString("100.0")
	m := optics.Morphism{{Lens: a, Value: x}}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal(100.0)
}

func TestLensStructJSON(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct{ A J }
	a := optics.Lenses1(T{})
	x, _ := a.FromString("{\"x\":\"abc\"}")
	m := optics.Morphism{{Lens: a, Value: x}}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A.X).Equal("abc")
}

func TestLensStructForm(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct {
		A J `content:"form"`
	}
	a := optics.Lenses1(T{})
	x, _ := a.FromString("x=abc")
	m := optics.Morphism{{Lens: a, Value: x}}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A.X).Equal("abc")
}

// func TestLensStructSeq(t *testing.T) {
// 	type T struct{ A []string }
// 	a := optics.Lenses1(T{})
// 	m := optics.Morphism{a: []string{"a", "b", "c"}}

// 	var x T
// 	err := m.Apply(&x)

// 	it.Ok(t).
// 		IfNil(err).
// 		If(x.A).Equal([]string{"a", "b", "c"})
// }

func TestLenses1(t *testing.T) {
	type T struct {
		A string
	}
	a := optics.Lenses1(T{})
	x, _ := a.FromString("a")
	m := optics.Morphism{{Lens: a, Value: x}}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a")
}

func TestLenses2(t *testing.T) {
	type T struct {
		A string
		B string
	}
	a, b := optics.Lenses2(T{})
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b")
}

func TestLenses3(t *testing.T) {
	type T struct {
		A string
		B string
		C string
	}
	a, b, c := optics.Lenses3(T{})
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c")
}

func TestLenses4(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
	}
	a, b, c, d := optics.Lenses4(T{})
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d")
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e")
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f")
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
		{Lens: g, Value: s},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f").
		If(v.G).Equal("g")
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	r, _ := h.FromString("h")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
		{Lens: g, Value: s},
		{Lens: h, Value: r},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	r, _ := h.FromString("h")
	u, _ := i.FromString("i")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
		{Lens: g, Value: s},
		{Lens: h, Value: r},
		{Lens: i, Value: u},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
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
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	p, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	r, _ := h.FromString("h")
	u, _ := i.FromString("i")
	n, _ := k.FromString("k")
	m := optics.Morphism{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: p},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
		{Lens: g, Value: s},
		{Lens: h, Value: r},
		{Lens: i, Value: u},
		{Lens: k, Value: n},
	}

	var v T
	err := m.Apply(&v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a").
		If(v.B).Equal("b").
		If(v.C).Equal("c").
		If(v.D).Equal("d").
		If(v.E).Equal("e").
		If(v.F).Equal("f").
		If(v.G).Equal("g").
		If(v.H).Equal("h").
		If(v.I).Equal("i").
		If(v.K).Equal("k")
}
