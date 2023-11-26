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

// Package optics is an internal, see golem.optics for public api
package optics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/ajg/form"
	"github.com/fogfish/golem/hseq"
	"github.com/fogfish/golem/optics"
)

// Value is co-product types matchable by patterns
// Note: extend it with caution, the structure size is optimized for performance
// See https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/
type Value struct {
	String string
	Number int
	Double float64
}

// Lens is composable setter of Value to "some" struct
type Lens interface {
	FromString(string) (Value, error)
	Put(any, Value) error
}

// Morphism is product of Lens and Value
type Morphism struct {
	Lens
	Value
}

// Morphisms is collection of lenses and values to be applied for object
type Morphisms []Morphism

func Morph[S any](m Morphisms, s *S) error {
	for _, arrow := range m {
		if err := arrow.Lens.Put(s, arrow.Value); err != nil {
			return err
		}
	}

	return nil
}

// Lens to deal with custom string type
type lensString[S, A any] struct{ optics.Reflector[A] }

func (l *lensString[S, A]) Put(s any, a Value) error {
	var t A
	*(*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = a.String

	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensString[S, A]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

// Lens to deal with custom string type
type lensStringPointer[S, A any] struct{ optics.Reflector[A] }

func (l *lensStringPointer[S, A]) Put(s any, a Value) error {
	var t A
	*(**string)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = &a.String
	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensStringPointer[S, A]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

// Lens to deal with custom number type
type lensNumber[S, A any] struct{ optics.Reflector[A] }

func (l *lensNumber[S, A]) Put(s any, a Value) error {
	var t A
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = a.Number
	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensNumber[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return Value{}, err
	}

	return Value{Number: val}, nil
}

// Lens to deal with custom number type
type lensNumberPointer[S, A any] struct{ optics.Reflector[A] }

func (l *lensNumberPointer[S, A]) Put(s any, a Value) error {
	var t A
	*(**int)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = &a.Number
	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensNumberPointer[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return Value{}, err
	}

	return Value{Number: val}, nil
}

// Lens to deal with custom double type
type lensDouble[S, A any] struct{ optics.Reflector[A] }

func (l *lensDouble[S, A]) Put(s any, a Value) error {
	var t A
	*(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = a.Double
	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensDouble[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{Double: val}, nil
}

// Lens to deal with custom double type
type lensDoublePointer[S, A any] struct{ optics.Reflector[A] }

func (l *lensDoublePointer[S, A]) Put(s any, a Value) error {
	var t A
	*(**float64)(unsafe.Pointer(uintptr(unsafe.Pointer(&t)))) = &a.Double
	l.Reflector.Putt(s, t)
	return nil
}

func (l *lensDoublePointer[S, A]) FromString(a string) (Value, error) {
	val, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return Value{}, err
	}

	return Value{Double: val}, nil
}

// Lens to deal with string type
type lensParser[S any] struct{ optics.Reflector[string] }

func (l *lensParser[S]) Put(s any, a Value) error {
	val := l.Reflector.Putt(s, a.String)
	switch err := (val).(type) {
	case error:
		return err
	default:
		return nil
	}
}

func (l *lensParser[S]) FromString(a string) (Value, error) {
	return Value{String: a}, nil
}

// lensStructJSON implements lens for complex "product" type
type lensStructJSON[A any] struct{ optics.Reflector[A] }

func newLensStructJSON[A any](r optics.Reflector[A]) optics.Reflector[string] {
	return &lensStructJSON[A]{r}
}

func (lens *lensStructJSON[A]) Putt(s any, a string) any {
	var o A

	if err := json.Unmarshal([]byte(a), &o); err != nil {
		return err
	}

	return lens.Reflector.Putt(s, o)
}

func (lens *lensStructJSON[A]) Gett(s any) string {
	v, err := json.Marshal(lens.Reflector.Gett(s))
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

func (lens *lensStructForm[A]) Putt(s any, a string) any {
	var o A

	if err := form.DecodeString(&o, a); err != nil {
		return err
	}

	return lens.Reflector.Putt(s, o)
}

func (lens *lensStructForm[A]) Gett(s any) string {
	v, err := form.EncodeToString(lens.Reflector.Gett(s))
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
			if t.StructField.Type.Kind() == reflect.Pointer {
				return &lensStringPointer[S, A]{ln.(optics.Reflector[A])}
			}
			return &lensString[S, A]{ln.(optics.Reflector[A])}
		case reflect.Int:
			if t.StructField.Type.Kind() == reflect.Pointer {
				return &lensNumberPointer[S, A]{ln.(optics.Reflector[A])}
			}
			return &lensNumber[S, A]{ln.(optics.Reflector[A])}
		case reflect.Float64:
			if t.StructField.Type.Kind() == reflect.Pointer {
				return &lensDoublePointer[S, A]{ln.(optics.Reflector[A])}
			}
			return &lensDouble[S, A]{ln.(optics.Reflector[A])}
		case reflect.Struct:
			switch t.Tag.Get("content") {
			case "form":
				return &lensParser[S]{newLensStructForm(ln.(optics.Reflector[A]))}
			case "application/x-www-form-urlencoded":
				return &lensParser[S]{newLensStructForm(ln.(optics.Reflector[A]))}
			case "json":
				return &lensParser[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			case "application/json":
				return &lensParser[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			default:
				return &lensParser[S]{newLensStructJSON(ln.(optics.Reflector[A]))}
			}
		default:
			panic(fmt.Errorf("type %v is not supported", t.Type))
		}
	}
}
