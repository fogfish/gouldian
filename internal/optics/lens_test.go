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

	"github.com/fogfish/golem/hseq"
	lenses "github.com/fogfish/golem/optics"
	"github.com/fogfish/gouldian/v2/internal/optics"
	"github.com/fogfish/it/v2"
)

func structTest[A any](t *testing.T, given string, expect A) {
	t.Helper()

	type T struct {
		A A
		B *A
	}

	a, b := hseq.FMap2(
		hseq.New[T]("A", "B"),
		optics.NewLens(lenses.NewLens[T, A]),
		optics.NewLens(lenses.NewLens[T, *A]),
	)
	x, _ := a.FromString(given)
	y, _ := b.FromString(given)

	var v T

	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
	}
	e := optics.Morph(m, &v)

	it.Then(t).Should(
		it.Nil(e),
		it.Equiv(v.A, expect),
		it.Equiv(*v.B, expect),
	)
}

func TestLensStructString(t *testing.T) {
	type String string

	structTest[string](t, "a", "a")
	structTest[String](t, "a", "a")
}

func TestLensStructInt(t *testing.T) {
	type Int int

	structTest[int](t, "100", 100)
	structTest[Int](t, "100", 100)
}

func TestLensStructIntFail(t *testing.T) {
	//lint:ignore U1000 type is used but not instantiated
	type T struct{ A int }
	a := hseq.FMap1(
		hseq.New[T]("A"),
		optics.NewLens(lenses.NewLens[T, int]),
	)
	_, err := a.FromString("abc")
	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestLensStructFloat(t *testing.T) {
	type Float64 float64

	structTest[float64](t, "100.101", 100.101)
	structTest[Float64](t, "100.101", 100.101)
}

func TestLensStructFloatFail(t *testing.T) {
	//lint:ignore U1000 type is used but not instantiated
	type T struct{ A float64 }
	a := hseq.FMap1(
		hseq.New[T]("A"),
		optics.NewLens(lenses.NewLens[T, float64]),
	)
	_, err := a.FromString("abc")
	it.Then(t).ShouldNot(
		it.Nil(err),
	)
}

func TestLensStructJSON(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}

	structTest[J](t, "{\"x\":\"abc\"}", J{"abc"})
}

func TestLensStructJSONFailt(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct{ A J }
	a := hseq.FMap1(
		hseq.New[T]("A"),
		optics.NewLens(lenses.NewLens[T, J]),
	)
	x, _ := a.FromString("{x:abc}")

	var v T
	m := optics.Morphisms{
		{Lens: a, Value: x},
	}
	e := optics.Morph(m, &v)

	it.Then(t).ShouldNot(
		it.Nil(e),
	)
}

func TestLensStructForm(t *testing.T) {
	type J struct {
		X string `json:"x"`
	}
	type T struct {
		A J `content:"form"`
	}
	a := hseq.FMap1(
		hseq.New[T]("A"),
		optics.NewLens(lenses.NewLens[T, J]),
	)
	x, _ := a.FromString("x=abc")
	m := optics.Morphisms{{Lens: a, Value: x}}

	var v T
	err := optics.Morph(m, &v)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A.X, "abc"),
	)
}

func TestLenses1(t *testing.T) {
	type T struct{ A string }
	a := hseq.FMap1(
		hseq.New[T]("A"),
		optics.NewLens(lenses.NewLens[T, string]),
	)
	x, _ := a.FromString("a")
	m := optics.Morphisms{{Lens: a, Value: x}}

	var v T
	err := optics.Morph(m, &v)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
	)
}

func TestLenses2(t *testing.T) {
	type T struct {
		A string
		B string
	}
	a, b := hseq.FMap2(
		hseq.New[T]("A", "B"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

	x, _ := a.FromString("a")
	y, _ := b.FromString("b")
	m := optics.Morphisms{
		{Lens: a, Value: x},
		{Lens: b, Value: y},
	}

	var v T
	err := optics.Morph(m, &v)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
	)
}

func TestLenses3(t *testing.T) {
	type T struct {
		A string
		B string
		C string
	}
	a, b, c := hseq.FMap3(
		hseq.New[T]("A", "B", "C"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
	)
}

func TestLenses4(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
	}
	a, b, c, d := hseq.FMap4(
		hseq.New[T]("A", "B", "C", "D"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
	)
}

func TestLenses5(t *testing.T) {
	type T struct {
		A string
		B string
		C string
		D string
		E string
	}
	a, b, c, d, e := hseq.FMap5(
		hseq.New[T]("A", "B", "C", "D", "E"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
		it.Equal(v.E, "e"),
	)
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
	a, b, c, d, e, f := hseq.FMap6(
		hseq.New[T]("A", "B", "C", "D", "E", "F"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
		it.Equal(v.E, "e"),
		it.Equal(v.F, "f"),
	)
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
	a, b, c, d, e, f, g := hseq.FMap7(
		hseq.New[T]("A", "B", "C", "D", "E", "F", "G"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
		it.Equal(v.E, "e"),
		it.Equal(v.F, "f"),
		it.Equal(v.G, "g"),
	)
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
	a, b, c, d, e, f, g, h := hseq.FMap8(
		hseq.New[T]("A", "B", "C", "D", "E", "F", "G", "H"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
		it.Equal(v.E, "e"),
		it.Equal(v.F, "f"),
		it.Equal(v.G, "g"),
		it.Equal(v.H, "h"),
	)
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
	a, b, c, d, e, f, g, h, i := hseq.FMap9(
		hseq.New[T]("A", "B", "C", "D", "E", "F", "G", "H", "I"),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
		optics.NewLens(lenses.NewLens[T, string]),
	)

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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(v.A, "a"),
		it.Equal(v.B, "b"),
		it.Equal(v.C, "c"),
		it.Equal(v.D, "d"),
		it.Equal(v.E, "e"),
		it.Equal(v.F, "f"),
		it.Equal(v.G, "g"),
		it.Equal(v.H, "h"),
		it.Equal(v.I, "i"),
	)
}
