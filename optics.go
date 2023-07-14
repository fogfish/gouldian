package gouldian

import (
	"github.com/fogfish/golem/hseq"
	lenses "github.com/fogfish/golem/optics"
	"github.com/fogfish/gouldian/v2/internal/optics"
)

// Lens type
type Lens struct{ optics.Lens }

func newLens[S, A any](ln func(t hseq.Type[S]) lenses.Lens[S, A]) func(hseq.Type[S]) Lens {
	return func(t hseq.Type[S]) Lens {
		return Lens{optics.NewLens(ln)(t)}
	}
}

// Optics1 unfold attribute(s) of type T
func Optics1[T, A any](attr ...string) Lens {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New1[T, A]()
	} else {
		seq = hseq.New[T](attr[0])
	}

	return hseq.FMap1(seq,
		newLens(lenses.NewLens[T, A]),
	)
}

// Optics2 unfold attribute(s) of type T
func Optics2[T, A, B any](attr ...string) (Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New2[T, A, B]()
	} else {
		seq = hseq.New[T](attr[0:2]...)
	}

	return hseq.FMap2(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
	)
}

// Optics3 unfold attribute(s) of type T
func Optics3[T, A, B, C any](attr ...string) (Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New3[T, A, B, C]()
	} else {
		seq = hseq.New[T](attr[0:3]...)
	}

	return hseq.FMap3(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
	)
}

// Optics4 unfold attribute(s) of type T
func Optics4[T, A, B, C, D any](attr ...string) (Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New4[T, A, B, C, D]()
	} else {
		seq = hseq.New[T](attr[0:4]...)
	}

	return hseq.FMap4(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
	)
}

// Optics5 unfold attribute(s) of type T
func Optics5[T, A, B, C, D, E any](attr ...string) (Lens, Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New5[T, A, B, C, D, E]()
	} else {
		seq = hseq.New[T](attr[0:5]...)
	}

	return hseq.FMap5(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
		newLens(lenses.NewLens[T, E]),
	)
}

// Optics6 unfold attribute(s) of type T
func Optics6[T, A, B, C, D, E, F any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New6[T, A, B, C, D, E, F]()
	} else {
		seq = hseq.New[T](attr[0:6]...)
	}

	return hseq.FMap6(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
		newLens(lenses.NewLens[T, E]),
		newLens(lenses.NewLens[T, F]),
	)
}

// Optics7 unfold attribute(s) of type T
func Optics7[T, A, B, C, D, E, F, G any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New7[T, A, B, C, D, E, F, G]()
	} else {
		seq = hseq.New[T](attr[0:7]...)
	}

	return hseq.FMap7(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
		newLens(lenses.NewLens[T, E]),
		newLens(lenses.NewLens[T, F]),
		newLens(lenses.NewLens[T, G]),
	)
}

// Optics8 unfold attribute(s) of type T
func Optics8[T, A, B, C, D, E, F, G, H any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New8[T, A, B, C, D, E, F, G, H]()
	} else {
		seq = hseq.New[T](attr[0:8]...)
	}

	return hseq.FMap8(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
		newLens(lenses.NewLens[T, E]),
		newLens(lenses.NewLens[T, F]),
		newLens(lenses.NewLens[T, G]),
		newLens(lenses.NewLens[T, H]),
	)
}

// Optics9 unfold attribute(s) of type T
func Optics9[T, A, B, C, D, E, F, G, H, I any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	var seq hseq.Seq[T]

	if len(attr) == 0 {
		seq = hseq.New9[T, A, B, C, D, E, F, G, H, I]()
	} else {
		seq = hseq.New[T](attr[0:9]...)
	}

	return hseq.FMap9(seq,
		newLens(lenses.NewLens[T, A]),
		newLens(lenses.NewLens[T, B]),
		newLens(lenses.NewLens[T, C]),
		newLens(lenses.NewLens[T, D]),
		newLens(lenses.NewLens[T, E]),
		newLens(lenses.NewLens[T, F]),
		newLens(lenses.NewLens[T, G]),
		newLens(lenses.NewLens[T, H]),
		newLens(lenses.NewLens[T, I]),
	)
}
