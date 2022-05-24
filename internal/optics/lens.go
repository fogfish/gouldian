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

/*

Package optics is an internal, see golem.optics for public api

*/
package optics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/ajg/form"
	"github.com/fogfish/golem/optics"
	"github.com/fogfish/golem/pure/hseq"
)

/*

Value is co-product types matchable by patterns
Note: do not extend the structure, optimal size for performance
See https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/
*/
type Value struct {
	String string
	Number int
	Double float64
}

/*

Lens is composable setter of Value to "some" struct
*/
type Lens interface {
	FromString(string) (Value, error)
	Put(reflect.Value, Value) error
}

/*

Morphism is product of Lens and Value
*/
type Morphism struct {
	Lens
	Value
}

/*

Morphisms is collection of lenses and values to be applied for object
*/
type Morphisms []Morphism

func Morph[S any](m Morphisms, s *S) error {
	g := reflect.ValueOf(s)

	for _, arrow := range m {
		if err := arrow.Lens.Put(g, arrow.Value); err != nil {
			return err
		}
	}

	return nil
}

/*

Lens to deal with string type
*/
type lensString[S any] struct{ optics.Reflector[string] }

func (l *lensString[S]) Put(s reflect.Value, a Value) error {
	return l.Reflector.PutValue(s, a.String)
}

func (l *lensString[S]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

/*

Lens to deal with number type
*/
type lensNumber[S any] struct{ optics.Reflector[int] }

func (l *lensNumber[S]) Put(s reflect.Value, a Value) error {
	return l.Reflector.PutValue(s, a.Number)
}

func (l *lensNumber[S]) FromString(a string) (Value, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return Value{}, err
	}

	return Value{Number: val}, nil
}

/*

Lens to deal with double type
*/
type lensDouble[S any] struct{ optics.Reflector[float64] }

func (l *lensDouble[S]) Put(s reflect.Value, a Value) error {
	return l.Reflector.PutValue(s, a.Double)
}

func (l *lensDouble[S]) FromString(a string) (Value, error) {
	val, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{Double: val}, nil
}

/*

lensStructJSON implements lens for complex "product" type
*/
type lensStructJSON[A any] struct{ optics.Reflector[A] }

func newLensStructJSON[A any](r optics.Reflector[A]) optics.Reflector[string] {
	return &lensStructJSON[A]{r}
}

func (lens *lensStructJSON[A]) PutValue(s reflect.Value, a string) error {
	var o A

	if err := json.Unmarshal([]byte(a), &o); err != nil {
		return err
	}

	return lens.Reflector.PutValue(s, o)
}

func (lens *lensStructJSON[A]) GetValue(s reflect.Value) string {
	v, err := json.Marshal(lens.Reflector.GetValue(s))
	if err != nil {
		panic(err)
	}

	return string(v)
}

/*

lensStructForm implements lens for complex "product" type
*/
type lensStructForm[A any] struct{ optics.Reflector[A] }

func newLensStructForm[A any](r optics.Reflector[A]) optics.Reflector[string] {
	return &lensStructForm[A]{r}
}

func (lens *lensStructForm[A]) PutValue(s reflect.Value, a string) error {
	var o A

	if err := form.DecodeString(&o, a); err != nil {
		return err
	}

	return lens.Reflector.PutValue(s, o)
}

func (lens *lensStructForm[A]) GetValue(s reflect.Value) string {
	v, err := form.EncodeToString(lens.Reflector.GetValue(s))
	if err != nil {
		panic(err)
	}

	return string(v)
}

/*

NewLens creates lense instance
*/
func NewLens[S, A any](ln optics.Lens[S, A]) func(t hseq.Type[S]) Lens {
	return func(t hseq.Type[S]) Lens {
		switch t.PureType.Kind() {
		case reflect.String:
			return &lensString[S]{ln.(optics.Reflector[string])}
		case reflect.Int:
			return &lensNumber[S]{ln.(optics.Reflector[int])}
		case reflect.Float64:
			return &lensDouble[S]{ln.(optics.Reflector[float64])}
		case reflect.Struct:
			switch t.Tag.Get("content") {
			case "form":
				return &lensString[S]{newLensStructForm(ln.(optics.Reflector[A]))}
			case "application/x-www-form-urlencoded":
				return &lensString[S]{newLensStructForm(ln.(optics.Reflector[A]))}
			case "json":
				return &lensString[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			case "application/json":
				return &lensString[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			default:
				return &lensString[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			}
		default:
			panic(fmt.Errorf("Type %v is not supported", t.Type))
		}
	}
}

/*

ForProduct1 split structure with 1 field to set of lenses
*/
func ForProduct1[T, A any]() Lens {
	a := optics.ForProduct1[T, A]()
	return hseq.FMap1(
		hseq.Generic[T](),
		NewLens(a),
	)
}

/*

ForProduct2 split structure with 2 fields to set of lenses
*/
func ForProduct2[T, A, B any]() (Lens, Lens) {
	a, b := optics.ForProduct2[T, A, B]()
	return hseq.FMap2(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
	)
}

/*

ForProduct3 split structure with 3 fields to set of lenses
*/
func ForProduct3[T, A, B, C any]() (Lens, Lens, Lens) {
	a, b, c := optics.ForProduct3[T, A, B, C]()
	return hseq.FMap3(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
	)
}

/*

ForProduct4 split structure with 4 fields to set of lenses
*/
func ForProduct4[T, A, B, C, D any]() (Lens, Lens, Lens, Lens) {
	a, b, c, d := optics.ForProduct4[T, A, B, C, D]()
	return hseq.FMap4(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
	)
}

/*

ForProduct5 split structure with 5 fields to set of lenses
*/
func ForProduct5[T, A, B, C, D, E any]() (Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e := optics.ForProduct5[T, A, B, C, D, E]()
	return hseq.FMap5(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
	)
}

/*

ForProduct6 split structure with 6 fields to set of lenses
*/
func ForProduct6[T, A, B, C, D, E, F any]() (Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f := optics.ForProduct6[T, A, B, C, D, E, F]()
	return hseq.FMap6(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
		NewLens(f),
	)
}

/*

ForProduct7 split structure with 7 fields to set of lenses
*/
func ForProduct7[T, A, B, C, D, E, F, G any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g := optics.ForProduct7[T, A, B, C, D, E, F, G]()
	return hseq.FMap7(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
		NewLens(f),
		NewLens(g),
	)
}

/*

ForProduct8 split structure with 8 fields to set of lenses
*/
func ForProduct8[T, A, B, C, D, E, F, G, H any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h := optics.ForProduct8[T, A, B, C, D, E, F, G, H]()
	return hseq.FMap8(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
		NewLens(f),
		NewLens(g),
		NewLens(h),
	)
}

/*

ForProduct9 split structure with 9 fields to set of lenses
*/
func ForProduct9[T, A, B, C, D, E, F, G, H, I any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h, i := optics.ForProduct9[T, A, B, C, D, E, F, G, H, I]()
	return hseq.FMap9(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
		NewLens(f),
		NewLens(g),
		NewLens(h),
		NewLens(i),
	)
}

/*

ForProduct10 split structure with 10 fields to set of lenses
*/
func ForProduct10[T, A, B, C, D, E, F, G, H, I, J any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h, i, j := optics.ForProduct10[T, A, B, C, D, E, F, G, H, I, J]()
	return hseq.FMap10(
		hseq.Generic[T](),
		NewLens(a),
		NewLens(b),
		NewLens(c),
		NewLens(d),
		NewLens(e),
		NewLens(f),
		NewLens(g),
		NewLens(h),
		NewLens(i),
		NewLens(j),
	)
}
