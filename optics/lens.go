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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/ajg/form"
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

Codec transforms string ⟼ Value
*/
type Codec interface {
	FromString(string) (Value, error)
}

/*

Lens is composable setter of Value to "some" struct
*/
type Lens interface {
	Codec
	Put(a reflect.Value, s Value) error
}

/*

Setter is product of Lens and Value
*/
type Setter struct {
	Lens
	Value
}

/*

Morphism is collection of lenses and values to be applied for object
*/
type Morphism []Setter

/*

Apply Morphism to "some" struct
*/
func (m Morphism) Apply(a interface{}) error {
	g := reflect.ValueOf(a)
	if g.Kind() != reflect.Ptr {
		return fmt.Errorf("Morphism requires pointer type, %s given", g.Kind().String())
	}

	for _, arrow := range m {
		if err := arrow.Lens.Put(g, arrow.Value); err != nil {
			return err
		}
	}

	return nil
}

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
type lensStructString struct{ lensStruct }

// FromString transforms string ⟼ Value[string]
func (lens lensStructString) FromString(s string) (Value, error) {
	return Value{String: s}, nil
}

// Put Value[string] to struct
func (lens lensStructString) Put(a reflect.Value, s Value) error {
	f := a.Elem().Field(int(lens.field))

	if f.Kind() == reflect.Ptr {
		p := reflect.New(lens.typeof.Elem())
		p.Elem().SetString(s.String)
		f.Set(p)
		return nil
	}

	f.SetString(s.String)
	return nil
}

/*

lensStructInt implements lens for int type
*/
type lensStructInt struct{ lensStruct }

// FromString transforms string ⟼ Value[int]
func (lens lensStructInt) FromString(s string) (Value, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return Value{}, err
	}

	return Value{Number: val}, nil
}

// Put Value[int] to struct
func (lens lensStructInt) Put(a reflect.Value, s Value) error {
	a.Elem().Field(int(lens.field)).SetInt(int64(s.Number))
	return nil
}

/*

lensStructFloat implements lens for float type
*/
type lensStructFloat struct{ lensStruct }

// FromString transforms string ⟼ Value[float64]
func (lens lensStructFloat) FromString(s string) (Value, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{Double: val}, nil
}

// Put Value[float64] to struct
func (lens lensStructFloat) Put(a reflect.Value, s Value) error {
	a.Elem().Field(int(lens.field)).SetFloat(s.Double)
	return nil
}

/*

lensStructJSON implements lens for complex "product" type
*/
type lensStructJSON struct{ lensStruct }

// FromString transforms string ⟼ Value[string]
func (lens lensStructJSON) FromString(s string) (Value, error) {
	return Value{String: s}, nil
}

// Put Value[string] to struct
func (lens lensStructJSON) Put(a reflect.Value, s Value) error {
	c := reflect.New(lens.typeof)
	o := c.Interface()

	if err := json.Unmarshal([]byte(s.String), &o); err != nil {
		return err
	}

	a.Elem().Field(int(lens.field)).Set(c.Elem())
	return nil
}

/*

lensStructForm implements lens for complex "product" type
*/
type lensStructForm struct{ lensStruct }

// FromString transforms string ⟼ Value[string]
func (lens lensStructForm) FromString(s string) (Value, error) {
	return Value{String: s}, nil
}

// Put Value[string] to struct
func (lens lensStructForm) Put(a reflect.Value, s Value) error {
	c := reflect.New(lens.typeof)
	o := c.Interface()

	buf := bytes.NewBuffer([]byte(s.String))
	if err := form.NewDecoder(buf).Decode(&o); err != nil {
		return err
	}

	a.Elem().Field(int(lens.field)).Set(c.Elem())
	return nil
}

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
func newLensStruct(id int, field reflect.StructField) Lens {
	typeof := field.Type.Kind()
	if typeof == reflect.Ptr {
		typeof = field.Type.Elem().Kind()
	}

	switch typeof {
	case reflect.String:
		return &lensStructString{lensStruct{id, field.Type}}
	case reflect.Int:
		return &lensStructInt{lensStruct{id, field.Type}}
	case reflect.Float64:
		return &lensStructFloat{lensStruct{id, field.Type}}
	case reflect.Struct:
		switch field.Tag.Get("content") {
		case "form":
			return &lensStructForm{lensStruct{id, field.Type}}
		case "application/x-www-form-urlencoded":
			return &lensStructForm{lensStruct{id, field.Type}}
		case "json":
			return &lensStructJSON{lensStruct{id, field.Type}}
		case "application/json":
			return &lensStructJSON{lensStruct{id, field.Type}}
		default:
			return &lensStructJSON{lensStruct{id, field.Type}}
		}
	// case reflect.Slice:
	// 	return &lensStructSeq{lensStruct{id, field.Type}}
	default:
		panic(fmt.Errorf("Unknown lens type %v", field.Type))
	}
}

func typeOf(t interface{}) reflect.Type {
	typeof := reflect.TypeOf(t)
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	return typeof
}

/*

ForProduct1 split structure with 1 field to set of lenses
*/
func ForProduct1(t interface{}) Lens {
	tc := typeOf(t)
	if tc.NumField() != 1 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 1", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0))
}

/*

ForProduct2 split structure with 2 fields to set of lenses
*/
func ForProduct2(t interface{}) (Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 2 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 2", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1))
}

/*

ForProduct3 split structure with 3 fields to set of lenses
*/
func ForProduct3(t interface{}) (Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 3 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 3", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2))
}

/*

ForProduct4 split structure with 4 fields to set of lenses
*/
func ForProduct4(t interface{}) (Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 4 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 4", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3))
}

/*

ForProduct5 split structure with 5 fields to set of lenses
*/
func ForProduct5(t interface{}) (Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 5 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 5", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4))
}

/*

ForProduct6 split structure with 6 fields to set of lenses
*/
func ForProduct6(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 6 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 6", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4)),
		newLensStruct(5, tc.Field(5))
}

/*

ForProduct7 split structure with 7 fields to set of lenses
*/
func ForProduct7(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 7 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 7", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4)),
		newLensStruct(5, tc.Field(5)),
		newLensStruct(6, tc.Field(6))
}

/*

ForProduct8 split structure with 8 fields to set of lenses
*/
func ForProduct8(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 8 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 8", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4)),
		newLensStruct(5, tc.Field(5)),
		newLensStruct(6, tc.Field(6)),
		newLensStruct(7, tc.Field(7))
}

/*

ForProduct9 split structure with 9 fields to set of lenses
*/
func ForProduct9(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 9 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 9", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4)),
		newLensStruct(5, tc.Field(5)),
		newLensStruct(6, tc.Field(6)),
		newLensStruct(7, tc.Field(7)),
		newLensStruct(8, tc.Field(8))
}

/*

ForProduct10 split structure with 10 fields to set of lenses
*/
func ForProduct10(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 10 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to 10 lens", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2)),
		newLensStruct(3, tc.Field(3)),
		newLensStruct(4, tc.Field(4)),
		newLensStruct(5, tc.Field(5)),
		newLensStruct(6, tc.Field(6)),
		newLensStruct(7, tc.Field(7)),
		newLensStruct(8, tc.Field(8)),
		newLensStruct(9, tc.Field(9))
}
