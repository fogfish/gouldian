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
