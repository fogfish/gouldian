/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package optics_test

import (
	"testing"

	"github.com/fogfish/gouldian/internal/optics"
	"github.com/fogfish/it"
)

func TestLensStructString(t *testing.T) {
	type String string
	type T struct {
		A string
		B *string
		C String
	}
	a, b, c := optics.ForProduct3[T, string, string, String]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")

	t.Run("ByVal", func(t *testing.T) {
		var v T

		m := optics.Morphisms{
			{Lens: a, Value: x},
			{Lens: b, Value: y},
			{Lens: c, Value: z},
		}
		e := optics.Morph(m, &v)

		it.Ok(t).
			IfNil(e).
			If(v.A).Equal("a").
			If(*v.B).Equal("b").
			If(v.C).Equal(String("c"))
	})
}

func TestLensStructInt(t *testing.T) {
	type Int int
	type T struct {
		A int
		B *int
		C Int
	}
	a, b, c := optics.ForProduct3[T, int, int, Int]()
	x, err := a.FromString("100")
	it.Ok(t).If(err).Must().Equal(nil)

	y, err := b.FromString("100")
	it.Ok(t).If(err).Must().Equal(nil)

	z, err := c.FromString("100")
	it.Ok(t).If(err).Must().Equal(nil)

	var v T

	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
	}
	err = optics.Morph(m, &v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal(100).
		If(*v.B).Equal(100).
		If(v.C).Equal(Int(100))
}

func TestLensStructIntFail(t *testing.T) {
	type T struct{ A int }
	a := optics.ForProduct1[T, int]()
	_, err := a.FromString("abc")
	it.Ok(t).If(err).MustNot().Equal(nil)
}

func TestLensStructFloat(t *testing.T) {
	type Float64 float64
	type T struct {
		A float64
		B *float64
		C Float64
	}
	a, b, c := optics.ForProduct3[T, float64, float64, Float64]()

	x, err := a.FromString("100.0")
	it.Ok(t).If(err).Must().Equal(nil)

	y, err := b.FromString("100.0")
	it.Ok(t).If(err).Must().Equal(nil)

	z, err := c.FromString("100.0")
	it.Ok(t).If(err).Must().Equal(nil)

	var v T

	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
	}
	err = optics.Morph(m, &v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal(100.0).
		If(*v.B).Equal(100.0).
		If(v.C).Equal(Float64(100.0))
}

func TestLensStructFloatFail(t *testing.T) {
	type T struct{ A float64 }
	a := optics.ForProduct1[T, float64]()
	_, err := a.FromString("abc")
	it.Ok(t).If(err).MustNot().Equal(nil)
}

func TestLensStructJSON(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct{ A J }
	a := optics.ForProduct1[T, J]()
	x, _ := a.FromString("{\"x\":\"abc\"}")
	m := optics.Morphisms{{Lens: a, Value: x}}

	var v T
	err := optics.Morph(m, &v)

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
	a := optics.ForProduct1[T, J]()
	x, _ := a.FromString("x=abc")
	m := optics.Morphisms{{Lens: a, Value: x}}

	var v T
	err := optics.Morph(m, &v)

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
	a := optics.ForProduct1[T, string]()
	x, _ := a.FromString("a")
	m := optics.Morphisms{{Lens: a, Value: x}}

	var v T
	err := optics.Morph(m, &v)

	it.Ok(t).
		IfNil(err).
		If(v.A).Equal("a")
}

func TestLenses2(t *testing.T) {
	type T struct {
		A string
		B string
	}
	a, b := optics.ForProduct2[T, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c := optics.ForProduct3[T, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c, d := optics.ForProduct4[T, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c, d, e := optics.ForProduct5[T, string, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c, d, e, f := optics.ForProduct6[T, string, string, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c, d, e, f, g := optics.ForProduct7[T, string, string, string, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
		{Lens: c, Value: z},
		{Lens: d, Value: k},
		{Lens: e, Value: q},
		{Lens: f, Value: w},
		{Lens: g, Value: s},
	}

	var v T
	err := optics.Morph(m, &v)

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
	a, b, c, d, e, f, g, h := optics.ForProduct8[T, string, string, string, string, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	r, _ := h.FromString("h")
	m := optics.Morphisms{
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
	err := optics.Morph(m, &v)

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
	a, b, c, d, e, f, g, h, i := optics.ForProduct9[T, string, string, string, string, string, string, string, string, string]()
	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	z, _ := c.FromString("c")
	k, _ := d.FromString("d")
	q, _ := e.FromString("e")
	w, _ := f.FromString("f")
	s, _ := g.FromString("g")
	r, _ := h.FromString("h")
	u, _ := i.FromString("i")
	m := optics.Morphisms{
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
	err := optics.Morph(m, &v)

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
	a, b, c, d, e, f, g, h, i, k := optics.ForProduct10[T, string, string, string, string, string, string, string, string, string, string]()
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
	m := optics.Morphisms{
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
	err := optics.Morph(m, &v)

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
