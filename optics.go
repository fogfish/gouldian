package gouldian

import (
	lenses "github.com/fogfish/golem/optics"
	"github.com/fogfish/golem/pure/hseq"
	"github.com/fogfish/gouldian/internal/optics"
)

type Lens struct{ optics.Lens }

func newLens[S, A any](ln lenses.Lens[S, A]) func(hseq.Type[S]) Lens {
	return func(t hseq.Type[S]) Lens {
		return Lens{optics.NewLens(ln)(t)}
	}
}

/*

ForProduct1 split structure with 1 field to set of lenses
*/
func Optics1[T, A any]() Lens {
	a := lenses.ForProduct1[T, A]()
	return hseq.FMap1(
		hseq.Generic[T](),
		newLens(a),
	)
}

/*

Optics2 split structure with 2 field to set of lenses
*/
func Optics2[T, A, B any]() (Lens, Lens) {
	a, b := lenses.ForProduct2[T, A, B]()
	return hseq.FMap2(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
	)
}

/*

Optics3 split structure with 3 field to set of lenses
*/
func Optics3[T, A, B, C, D any]() (Lens, Lens, Lens) {
	a, b, c := lenses.ForProduct3[T, A, B, C]()
	return hseq.FMap3(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
	)
}

/*

Optics4 split structure with 4 field to set of lenses
*/
func Optics4[T, A, B, C, D, E any]() (Lens, Lens, Lens, Lens) {
	a, b, c, d := lenses.ForProduct4[T, A, B, C, D]()
	return hseq.FMap4(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
	)
}

/*

Optics5 split structure with 5 field to set of lenses
*/
func Optics5[T, A, B, C, D, E any]() (Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e := lenses.ForProduct5[T, A, B, C, D, E]()
	return hseq.FMap5(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
	)
}

/*

Optics6 split structure with 6 field to set of lenses
*/
func Optics6[T, A, B, C, D, E, F any]() (Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f := lenses.ForProduct6[T, A, B, C, D, E, F]()
	return hseq.FMap6(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
		newLens(f),
	)
}

/*

Optics7 split structure with 7 field to set of lenses
*/
func Optics7[T, A, B, C, D, E, F, G any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g := lenses.ForProduct7[T, A, B, C, D, E, F, G]()
	return hseq.FMap7(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
		newLens(f),
		newLens(g),
	)
}

/*

Optics8 split structure with 8 field to set of lenses
*/
func Optics8[T, A, B, C, D, E, F, G, H any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h := lenses.ForProduct8[T, A, B, C, D, E, F, G, H]()
	return hseq.FMap8(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
		newLens(f),
		newLens(g),
		newLens(h),
	)
}

/*

Optics9 split structure with 9 field to set of lenses
*/
func Optics9[T, A, B, C, D, E, F, G, H, I any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h, i := lenses.ForProduct9[T, A, B, C, D, E, F, G, H, I]()
	return hseq.FMap9(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
		newLens(f),
		newLens(g),
		newLens(h),
		newLens(i),
	)
}

/*

Optics10 split structure with 10 field to set of lenses
*/
func Optics10[T, A, B, C, D, E, F, G, H, I, J any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h, i, j := lenses.ForProduct10[T, A, B, C, D, E, F, G, H, I, J]()
	return hseq.FMap10(
		hseq.Generic[T](),
		newLens(a),
		newLens(b),
		newLens(c),
		newLens(d),
		newLens(e),
		newLens(f),
		newLens(g),
		newLens(h),
		newLens(i),
		newLens(j),
	)
}
