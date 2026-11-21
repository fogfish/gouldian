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

package gouldian

import (
	"fmt"
	"reflect"

	"github.com/fogfish/golem/hseq"
	"github.com/fogfish/gouldian/v3/internal/optics"
)

// Lens is a setter of decoded HTTP value to a field of struct.
// The decoder is chosen once when the lens is constructed:
//
//  1. field tag `content:"json"`, `content:"json,strict"`, `content:"form"`;
//  2. decoder registered by RegisterCodec;
//  3. encoding.TextUnmarshaler (the type validates itself);
//  4. string, bool, int*, uint*, float*, time.Duration, []byte (raw),
//     slices of them (repeated param or tag `sep:","`), pointers;
//  5. struct is JSON by default.
//
// Unsupported type panics when the lens is constructed (at start-up).
type Lens struct{ optics.Lens }

func newLens[S any](t hseq.Type[S]) Lens {
	return Lens{optics.NewLens(t)}
}

// RegisterCodec registers a decoder for the type T. The decoder is used by
// all lenses created after the registration.
func RegisterCodec[T any](decode func(string) (T, error)) {
	optics.Register(decode)
}

/*
At builds a lens for the field of struct T. It is a type-safe-enough
alternative to OpticsN if the lens is required inline. The name of the field
is checked when the lens is constructed, an unknown field panics at start-up.

	µ.POST(µ.URI(µ.Path("ds"), µ.Path(µ.At[reqCreate]("Class"))))
*/
func At[T any](field string) Lens {
	return hseq.FMap1(hseq.New[T](field), newLens[T])
}

// unfold finds the fields of struct T. Without names the fields are found by
// types (n-th occurrence of the type binds n-th field of the type). With names
// the type arguments are asserted against the types of the fields, a mismatch
// panics at start-up.
func unfold[T any](n int, attr []string, byType func() hseq.Seq[T], types ...reflect.Type) hseq.Seq[T] {
	if len(attr) == 0 {
		return byType()
	}

	seq := hseq.New[T](attr[0:n]...)
	for i, f := range seq {
		if want := types[i]; f.Type.String() != want.String() || !f.Type.AssignableTo(want) {
			cat := reflect.TypeFor[T]()
			if cat.Kind() == reflect.Pointer {
				cat = cat.Elem()
			}
			panic(fmt.Errorf("\n==> Critical Error:\n❌ field `%s` of struct `%s` has type `%s`, not `%s`\n\n🛠️ How To Fix:\n- use the type of the field as the type argument of Optics.\n- check for typos in the name of the field",
				f.Name, cat.Name(), f.Type, want))
		}
	}
	return seq
}

// Optics1 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type (hseq v1.3.2).
func Optics1[T, A any](attr ...string) Lens {
	seq := unfold(1, attr, hseq.New1[T, A], reflect.TypeFor[A]())

	return hseq.FMap1(seq,
		newLens[T],
	)
}

// Optics2 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics2[T, A, B any](attr ...string) (Lens, Lens) {
	seq := unfold(2, attr, hseq.New2[T, A, B], reflect.TypeFor[A](), reflect.TypeFor[B]())

	return hseq.FMap2(seq,
		newLens[T],
		newLens[T],
	)
}

// Optics3 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics3[T, A, B, C any](attr ...string) (Lens, Lens, Lens) {
	seq := unfold(3, attr, hseq.New3[T, A, B, C], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C]())

	return hseq.FMap3(seq,
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics4 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics4[T, A, B, C, D any](attr ...string) (Lens, Lens, Lens, Lens) {
	seq := unfold(4, attr, hseq.New4[T, A, B, C, D], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D]())

	return hseq.FMap4(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics5 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics5[T, A, B, C, D, E any](attr ...string) (Lens, Lens, Lens, Lens, Lens) {
	seq := unfold(5, attr, hseq.New5[T, A, B, C, D, E], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D](), reflect.TypeFor[E]())

	return hseq.FMap5(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics6 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics6[T, A, B, C, D, E, F any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens) {
	seq := unfold(6, attr, hseq.New6[T, A, B, C, D, E, F], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D](), reflect.TypeFor[E](), reflect.TypeFor[F]())

	return hseq.FMap6(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics7 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics7[T, A, B, C, D, E, F, G any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	seq := unfold(7, attr, hseq.New7[T, A, B, C, D, E, F, G], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D](), reflect.TypeFor[E](), reflect.TypeFor[F](), reflect.TypeFor[G]())

	return hseq.FMap7(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics8 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics8[T, A, B, C, D, E, F, G, H any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	seq := unfold(8, attr, hseq.New8[T, A, B, C, D, E, F, G, H], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D](), reflect.TypeFor[E](), reflect.TypeFor[F](), reflect.TypeFor[G](), reflect.TypeFor[H]())

	return hseq.FMap8(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}

// Optics9 unfold attribute(s) of type T. If the names of attributes are not
// defined, the fields are discovered by types. The n-th occurrence of the
// type binds the n-th field of the type.
func Optics9[T, A, B, C, D, E, F, G, H, I any](attr ...string) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	seq := unfold(9, attr, hseq.New9[T, A, B, C, D, E, F, G, H, I], reflect.TypeFor[A](), reflect.TypeFor[B](), reflect.TypeFor[C](), reflect.TypeFor[D](), reflect.TypeFor[E](), reflect.TypeFor[F](), reflect.TypeFor[G](), reflect.TypeFor[H](), reflect.TypeFor[I]())

	return hseq.FMap9(seq,
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
		newLens[T],
	)
}
