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

// Value is co-product types matchable by patterns
// Note: do not extend the structure, optimal size for performance
// See https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/
type Value struct {
	String string
	Number int
	Double float64
}

// Lens is composable setter of Value to "some" struct
type Lens interface {
	FromString(string) (Value, error)
	Put(reflect.Value, Value) error
}

// Morphism is product of Lens and Value
type Morphism struct {
	Lens
	Value
}

// Morphisms is collection of lenses and values to be applied for object
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

// Lens to deal with string type
type lensString[S any] struct{ optics.Reflector[string] }

func (l *lensString[S]) Put(s reflect.Value, a Value) error {
	return l.Reflector.PutValue(s, a.String)
}

func (l *lensString[S]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

// Lens to deal with custom string type
type lensStringTyped[S, A any] struct{ optics.Reflector[A] }

func (l *lensStringTyped[S, A]) Put(s reflect.Value, a Value) error {
	var t A
	r := reflect.ValueOf(&t).Elem()
	r.SetString(a.String)
	return l.Reflector.PutValue(s, t)
}

func (l *lensStringTyped[S, A]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

// Lens to deal with number type
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

// Lens to deal with custom number type
type lensNumberTyped[S, A any] struct{ optics.Reflector[A] }

func (l *lensNumberTyped[S, A]) Put(s reflect.Value, a Value) error {
	var t A
	r := reflect.ValueOf(&t).Elem()
	r.SetInt(int64(a.Number))
	return l.Reflector.PutValue(s, t)
}

func (l *lensNumberTyped[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return Value{}, err
	}

	return Value{Number: val}, nil
}

// Lens to deal with double type
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

// Lens to deal with custom double type
type lensDoubleTyped[S, A any] struct{ optics.Reflector[A] }

func (l *lensDoubleTyped[S, A]) Put(s reflect.Value, a Value) error {
	var t A
	r := reflect.ValueOf(&t).Elem()
	r.SetFloat(a.Double)
	return l.Reflector.PutValue(s, t)
}

func (l *lensDoubleTyped[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{Double: val}, nil
}

// lensStructJSON implements lens for complex "product" type
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

// lensStructForm implements lens for complex "product" type
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

// NewLens creates lense instance
func NewLens[S, A any](fln func(t hseq.Type[S]) optics.Lens[S, A]) func(t hseq.Type[S]) Lens {
	return func(t hseq.Type[S]) Lens {
		ln := fln(t)
		switch t.PureType.Kind() {
		case reflect.String:
			if t.PureType.Name() == "string" {
				return &lensString[S]{ln.(optics.Reflector[string])}
			}
			return &lensStringTyped[S, A]{ln.(optics.Reflector[A])}
		case reflect.Int:
			if t.PureType.Name() == "int" {
				return &lensNumber[S]{ln.(optics.Reflector[int])}
			}
			return &lensNumberTyped[S, A]{ln.(optics.Reflector[A])}
		case reflect.Float64:
			if t.PureType.Name() == "float64" {
				return &lensDouble[S]{ln.(optics.Reflector[float64])}
			}
			return &lensDoubleTyped[S, A]{ln.(optics.Reflector[A])}
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
			panic(fmt.Errorf("type %v is not supported", t.Type))
		}
	}
}
