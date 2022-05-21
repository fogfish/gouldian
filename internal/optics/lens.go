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

Package optics ...

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

Apply Morphism to "some" struct
@deprecated
*/
func (m Morphisms) Apply(s interface{}) error {
	g := reflect.ValueOf(s)
	// if g.Kind() != reflect.Ptr {
	// 	return fmt.Errorf("Morphism requires pointer type, %s given", g.Kind().String())
	// }

	// p := unsafe.Pointer(&a)
	for _, arrow := range m {
		if err := arrow.Lens.Put(g, arrow.Value); err != nil {
			return err
		}
	}

	return nil
}

/*

...
*/
type lensString[S any] struct{ optics.Reflector[string] }

func (l *lensString[S]) Put(s reflect.Value, a Value) error {
	return l.Reflector.PutValue(s, a.String)
}

func (l *lensString[S]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

/*

...
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

...
*/
// func newCodec[T, A any](t hseq.Type[T]) Codec[A] {
// 	switch t.Type.Kind() {
// 	case reflect.String:
// 		return codecForString().(Codec[A])
// 	case reflect.Int:
// 		return codecForInt().(Codec[A])
// 	case reflect.Float64:
// 		return codecForFloat64().(Codec[A])
// 	}
// 	return nil
// }

// //
// type codecString string

// func codecForString() Codec[string] { return codecString("codec.string") }

// func (codecString) FromString(a string) (string, error) {
// 	return a, nil
// }

// //
// type codecInt string

// func codecForInt() Codec[int] { return codecInt("codec.int") }

// func (codecInt) FromString(a string) (int, error) {
// 	return strconv.Atoi(a)
// }

// //
// type codecFloat64 string

// func codecForFloat64() Codec[float64] { return codecFloat64("codec.float64") }

// func (codecFloat64) FromString(a string) (float64, error) {
// 	return strconv.ParseFloat(a, 64)
// }

/*

lensStruct is a type for any lens
*/
type lensStruct struct {
	field  int
	typeof reflect.Type
}

/*

lensStructString implements lens for string type
*/
// type lensStructString struct{ lensStruct }

// // FromString transforms string ⟼ Value[string]
// func (lens lensStructString) FromString(s string) (Value, error) {
// 	return Value{String: s}, nil
// }

// // Put Value[string] to struct
// func (lens lensStructString) Put(a reflect.Value, s Value) error {
// 	f := a.Elem().Field(int(lens.field))

// 	if f.Kind() == reflect.Ptr {
// 		p := reflect.New(lens.typeof.Elem())
// 		p.Elem().SetString(s.String)
// 		f.Set(p)
// 		return nil
// 	}

// 	f.SetString(s.String)
// 	return nil
// }

/*

lensStructInt implements lens for int type
*/
// type lensStructInt struct{ lensStruct }

// // FromString transforms string ⟼ Value[int]
// func (lens lensStructInt) FromString(s string) (Value, error) {
// 	val, err := strconv.Atoi(s)
// 	if err != nil {
// 		return Value{}, err
// 	}

// 	return Value{Number: val}, nil
// }

// // Put Value[int] to struct
// func (lens lensStructInt) Put(a reflect.Value, s Value) error {
// 	a.Elem().Field(int(lens.field)).SetInt(int64(s.Number))
// 	return nil
// }

/*

lensStructFloat implements lens for float type
*/
// type lensStructFloat struct{ lensStruct }

// // FromString transforms string ⟼ Value[float64]
// func (lens lensStructFloat) FromString(s string) (Value, error) {
// 	val, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		return Value{}, err
// 	}

// 	return Value{Double: val}, nil
// }

// // Put Value[float64] to struct
// func (lens lensStructFloat) Put(a reflect.Value, s Value) error {
// 	a.Elem().Field(int(lens.field)).SetFloat(s.Double)
// 	return nil
// }

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

// // FromString transforms string ⟼ Value[string]
// func (lens lensStructJSON) FromString(s string) (Value, error) {
// 	return Value{String: s}, nil
// }

// // Put Value[string] to struct
// func (lens lensStructJSON) Put(a reflect.Value, s Value) error {
// 	c := reflect.New(lens.typeof)
// 	o := c.Interface()

// 	if err := json.Unmarshal([]byte(s.String), &o); err != nil {
// 		return err
// 	}

// 	a.Elem().Field(int(lens.field)).Set(c.Elem())
// 	return nil
// }

/*

lensStructForm implements lens for complex "product" type
*/
type lensStructForm[S, A any] struct{ optics.Lens[S, A] }

func newLensStructForm[S, A any](l optics.Lens[S, A]) optics.Lens[S, string] {
	return lensStructForm[S, A]{l}
}

func (lens lensStructForm[S, A]) Put(s *S, a string) error {
	var o A

	if err := form.DecodeString(&o, a); err != nil {
		return err
	}
	// buf := bytes.NewBuffer([]byte(a))
	// if err := form.NewDecoder(buf).Decode(&o); err != nil {
	// 	return err
	// }

	return lens.Lens.Put(s, o)
}

func (lens lensStructForm[S, A]) Get(s *S) string {
	v, err := form.EncodeToString(lens.Lens.Get(s))
	if err != nil {
		panic(err)
	}

	return v
}

// type lensStructForm struct{ lensStruct }

// FromString transforms string ⟼ Value[string]
// func (lens lensStructForm) FromString(s string) (Value, error) {
// 	return Value{String: s}, nil
// }

// // Put Value[string] to struct
// func (lens lensStructForm) Put(a reflect.Value, s Value) error {
// 	c := reflect.New(lens.typeof)
// 	o := c.Interface()

// 	buf := bytes.NewBuffer([]byte(s.String))
// 	if err := form.NewDecoder(buf).Decode(&o); err != nil {
// 		return err
// 	}

// 	a.Elem().Field(int(lens.field)).Set(c.Elem())
// 	return nil
// }

/*

lensStructSeq ...
*/
/*
type lensStructSeq struct{ lensStruct }

func (lensStructSeq) FromString(s []string) (Value, error) {
	return Value{String: s[0]}, nil
}

func (lens lensStructSeq) Put(a reflect.Value, s Value) error {
	v := reflect.ValueOf(s.String)
	switch v.Type().Kind() {
	case reflect.Slice:
		a.Elem().Field(int(lens.field)).Set(v)
	default:
		return fmt.Errorf("Cannot cast %T to Seq", s)
	}

	return nil
}
*/

/*

newLensStruct creates lens
*/
// func newLensStruct(id int, field reflect.StructField) Lens {
// 	typeof := field.Type.Kind()
// 	if typeof == reflect.Ptr {
// 		typeof = field.Type.Elem().Kind()
// 	}

// 	switch typeof {
// 	case reflect.String:
// 		return &lensStructString{lensStruct{id, field.Type}}
// 	case reflect.Int:
// 		return &lensStructInt{lensStruct{id, field.Type}}
// 	case reflect.Float64:
// 		return &lensStructFloat{lensStruct{id, field.Type}}
// 	case reflect.Struct:
// 		switch field.Tag.Get("content") {
// 		case "form":
// 			return &lensStructForm{lensStruct{id, field.Type}}
// 		case "application/x-www-form-urlencoded":
// 			return &lensStructForm{lensStruct{id, field.Type}}
// 		case "json":
// 			return &lensStructJSON{lensStruct{id, field.Type}}
// 		case "application/json":
// 			return &lensStructJSON{lensStruct{id, field.Type}}
// 		default:
// 			return &lensStructJSON{lensStruct{id, field.Type}}
// 		}
// 	// case reflect.Slice:
// 	// 	return &lensStructSeq{lensStruct{id, field.Type}}
// 	default:
// 		panic(fmt.Errorf("Unknown lens type %v", field.Type))
// 	}
// }

// func typeOf(t interface{}) reflect.Type {
// 	typeof := reflect.TypeOf(t)
// 	if typeof.Kind() == reflect.Ptr {
// 		typeof = typeof.Elem()
// 	}

// 	return typeof
// }

func mkLens[S, A any](ln optics.Lens[S, A]) func(t hseq.Type[S]) Lens {
	return func(t hseq.Type[S]) Lens {
		switch t.Type.Kind() {
		case reflect.String:
			return &lensString[S]{ln.(optics.Reflector[string])}
		case reflect.Int:
			return &lensNumber[S]{ln.(optics.Reflector[int])}
		default:
			panic(fmt.Errorf("Type %v is not supported", t.Type))
		}

		// if t.Type.Kind() == reflect.Struct {
		// 	switch t.Tag.Get("content") {
		// 	// 	case "form":
		// 	// 		ln.lens = newLensStructForm(l) //.(optics.Lens[S, A])
		// 	// 	case "application/x-www-form-urlencoded":
		// 	// 		ln.lens = newLensStructForm(l) //.(optics.Lens[S, A])
		// 	// 	case "json":
		// 	// 		ln.lens = newLensStructJSON(l) //.(optics.Lens[S, A])
		// 	// 	case "application/json":
		// 	// 		ln.lens = newLensStructJSON(l) //.(optics.Lens[S, A])
		// 	default:
		// 		ln.lens = newLensStructJSON(refl) //.(optics.Reflector[A])
		// 	}
		// }
		// return ln
	}
}

/*

ForProduct1 split structure with 1 field to set of lenses
*/
func ForProduct1[T, A any]() Lens {
	a := optics.ForProduct1[T, A]()
	return hseq.FMap1(
		hseq.Generic[T](),
		mkLens(a),
	)
}

// func ForProduct1(t interface{}) Lens {
// 	tc := typeOf(t)
// 	if tc.NumField() != 1 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 1", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0))
// }

/*

ForProduct2 split structure with 2 fields to set of lenses
*/
// func ForProduct2(t interface{}) (Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 2 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 2", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1))
// }

// /*

// ForProduct3 split structure with 3 fields to set of lenses
// */
// func ForProduct3(t interface{}) (Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 3 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 3", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2))
// }

// /*

// ForProduct4 split structure with 4 fields to set of lenses
// */
// func ForProduct4(t interface{}) (Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 4 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 4", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3))
// }

// /*

// ForProduct5 split structure with 5 fields to set of lenses
// */
// func ForProduct5(t interface{}) (Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 5 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 5", tc.Name(), tc.NumField()))
// 	}
func ForProduct5[T, A, B, C, D, E any]() (Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e := optics.ForProduct5[T, A, B, C, D, E]()
	return hseq.FMap5(
		hseq.Generic[T](),
		mkLens(a),
		mkLens(b),
		mkLens(c),
		mkLens(d),
		mkLens(e),
	)
}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4))
// }

// /*

// ForProduct6 split structure with 6 fields to set of lenses
// */
// func ForProduct6(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 6 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 6", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4)),
// 		newLensStruct(5, tc.Field(5))
// }

// /*

// ForProduct7 split structure with 7 fields to set of lenses
// */
// func ForProduct7(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 7 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 7", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4)),
// 		newLensStruct(5, tc.Field(5)),
// 		newLensStruct(6, tc.Field(6))
// }

// /*

// ForProduct8 split structure with 8 fields to set of lenses
// */
// func ForProduct8(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 8 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 8", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4)),
// 		newLensStruct(5, tc.Field(5)),
// 		newLensStruct(6, tc.Field(6)),
// 		newLensStruct(7, tc.Field(7))
// }

// /*

// ForProduct9 split structure with 9 fields to set of lenses
// */
// func ForProduct9(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 9 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 9", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4)),
// 		newLensStruct(5, tc.Field(5)),
// 		newLensStruct(6, tc.Field(6)),
// 		newLensStruct(7, tc.Field(7)),
// 		newLensStruct(8, tc.Field(8))
// }

// /*

// ForProduct10 split structure with 10 fields to set of lenses
// */
// func ForProduct10(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
// 	tc := typeOf(t)
// 	if tc.NumField() != 10 {
// 		panic(fmt.Errorf("Unable to unapply type |%s| = %d to 10 lens", tc.Name(), tc.NumField()))
// 	}

// 	return newLensStruct(0, tc.Field(0)),
// 		newLensStruct(1, tc.Field(1)),
// 		newLensStruct(2, tc.Field(2)),
// 		newLensStruct(3, tc.Field(3)),
// 		newLensStruct(4, tc.Field(4)),
// 		newLensStruct(5, tc.Field(5)),
// 		newLensStruct(6, tc.Field(6)),
// 		newLensStruct(7, tc.Field(7)),
// 		newLensStruct(8, tc.Field(8)),
// 		newLensStruct(9, tc.Field(9))
// }
func ForProduct10[T, A, B, C, D, E, F, G, H, I, J any]() (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	a, b, c, d, e, f, g, h, i, j := optics.ForProduct10[T, A, B, C, D, E, F, G, H, I, J]()
	return hseq.FMap10(
		hseq.Generic[T](),
		mkLens(a),
		mkLens(b),
		mkLens(c),
		mkLens(d),
		mkLens(e),
		mkLens(f),
		mkLens(g),
		mkLens(h),
		mkLens(i),
		mkLens(j),
	)
}
