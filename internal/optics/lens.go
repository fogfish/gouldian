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

// Package optics is an internal, it implements lenses that bind strings
// received by HTTP request to fields of Go struct.
//
// A codec is chosen once, when the lens is constructed. The order is
//  1. content tag of the field (`content:"json"`, `json,strict`, `form`);
//  2. codec registered by application (see Register);
//  3. encoding.TextUnmarshaler;
//  4. built-in codecs: string, bool, ints, uints, floats, time.Duration,
//     []byte (raw), slices of scalars, pointers to any of above;
//  5. struct without a tag is decoded as JSON.
//
// An unsupported type panics at construction (start-up), naming the field.
package optics

import (
	"encoding"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/fogfish/golem/hseq"
)

// Value is the result of decoding. Most common scalars use dedicated fields
// to avoid allocations, any other type is carried by Ptr.
//
// WARNING: THE SIZE OF THIS STRUCTURE IS PERFORMANCE CRITICAL, EDIT WITH CARE.
//
// Value is copied on the hot path (Lens.FromString returns it, Morphism
// embeds it, Put receives it). It is 32 bytes (4 words) on purpose, the same
// size as in v2, and the whole request path is tuned around it.
//
// Go has no pass-by-reference, every argument is a copy. What differs is where
// the copy lives. The compiler keeps a value in registers (SSA-able, see
// CanSSA in cmd/compile/internal/ssa/value.go) only if its size is at most
// 4*PtrSize = 32 bytes and a struct has at most 4 fields. A bigger struct
// lives in memory and every pass/return of it is a memory copy. When the
// previous revision grew the structure to 40 bytes (String, Number, Any `any`)
// the routing benchmarks became 10-20% slower without a single extra
// allocation, and shrinking it back to 32 bytes recovered it:
//
//	            PathParam5  ServerParam5  GitHub
//	32 bytes      55 ns        87 ns      150 ns
//	40 bytes      66 ns       101 ns      160 ns
//
// Rules:
//   - do not add a field (4 is the limit), do not widen a field (e.g. float64
//     next to Number, an `any`, a slice, a second string);
//   - floats are stored in Number as math.Float64bits (see Double);
//   - any other decoded value lives in the memory the codec allocated, Ptr
//     points to it, the codec that created it knows the type;
//   - Morphism (Lens + Value) is 48 bytes, it is stored in a pooled slice;
//   - TestValueSize guards the size. If you must change it, run
//     `go test -run '^$' -bench . -benchmem -count 6 .` before and after
//     and compare the medians with `python3 bench/compare.py`
//     (see bench/README.md), PathParam*, ServerParam*, LensForProduct* and
//     GitHub must not regress.
//
// See also https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/
type Value struct {
	String string
	// Number carries integers, booleans (0/1) and floats (math.Float64bits)
	Number int64
	// Ptr points to the decoded value of any other type (*T allocated by
	// the codec), the codec that created it knows T. It keeps Value in 32 bytes.
	Ptr unsafe.Pointer
}

// Double returns the float stored by float codecs
func (v Value) Double() float64 { return math.Float64frombits(uint64(v.Number)) }

func fromDouble(f float64) Value { return Value{Number: int64(math.Float64bits(f))} }

// Lens is composable setter of Value to "some" struct
type Lens interface {
	// FromString decodes a single string (path segment, header, param)
	FromString(string) (Value, error)
	// FromBytes decodes a binary payload (body)
	FromBytes([]byte) (Value, error)
	// FromStrings decodes all values of repeated param, `?a=1&a=2`
	FromStrings([]string) (Value, error)
	// Multi is true if the lens consumes all values of repeated param
	Multi() bool
	// Put updates a field of struct `s` (must be a pointer to the struct)
	Put(s any, v Value) error
}

// Morphism is product of Lens and Value
type Morphism struct {
	Lens
	Value
}

// Morphisms is collection of lenses and values to be applied for object
type Morphisms []Morphism

// Morph applies all morphisms to struct
func Morph[S any](m Morphisms, s *S) error {
	for _, arrow := range m {
		if err := arrow.Lens.Put(s, arrow.Value); err != nil {
			return err
		}
	}

	return nil
}

//------------------------------------------------------------------------------

var registry sync.Map // reflect.Type -> *codec

// Register custom decoder for the type
func Register[T any](f func(string) (T, error)) {
	ft := reflect.TypeFor[T]()
	registry.Store(ft, &codec{
		str: func(s string) (Value, error) {
			v, err := f(s)
			if err != nil {
				return Value{}, err
			}
			p := new(T)
			*p = v
			return Value{Ptr: unsafe.Pointer(p)}, nil
		},
		set: func(dst reflect.Value, v Value) { dst.Set(reflect.NewAt(ft, v.Ptr).Elem()) },
	})
}

type codec struct {
	str   func(string) (Value, error)
	raw   func([]byte) (Value, error) // optional, str is used otherwise
	multi func([]string) (Value, error)
	set   func(dst reflect.Value, v Value)
}

type fastKind uint8

const (
	fastNone fastKind = iota
	fastString
	fastInt
	fastFloat
)

type lens struct {
	root reflect.Type // *S
	off  uintptr
	ft   reflect.Type
	fast fastKind
	c    *codec
}

func (l *lens) FromString(s string) (Value, error) {
	switch l.fast {
	case fastString:
		return Value{String: s}, nil
	case fastInt:
		v, err := strconv.Atoi(s)
		return Value{Number: int64(v)}, err
	}

	if l.c.multi != nil {
		return l.c.multi([]string{s})
	}
	return l.c.str(s)
}

func (l *lens) FromBytes(b []byte) (Value, error) {
	if l.c.raw != nil {
		return l.c.raw(b)
	}
	return l.FromString(string(b))
}

func (l *lens) FromStrings(s []string) (Value, error) {
	if l.c.multi != nil {
		return l.c.multi(s)
	}
	return l.c.str(s[0])
}

func (l *lens) Multi() bool { return l.c.multi != nil }

func (l *lens) Put(s any, v Value) error {
	if reflect.TypeOf(s) != l.root {
		return fmt.Errorf("invalid type %T passed to lens of %v", s, l.root)
	}

	ptr := unsafe.Add(reflect.ValueOf(s).UnsafePointer(), l.off)
	switch l.fast {
	case fastString:
		*(*string)(ptr) = v.String
	case fastInt:
		*(*int)(ptr) = int(v.Number)
	case fastFloat:
		*(*float64)(ptr) = v.Double()
	default:
		l.c.set(reflect.NewAt(l.ft, ptr).Elem(), v)
	}
	return nil
}

// NewLens creates lens for the field of struct S
func NewLens[S any](t hseq.Type[S]) Lens {
	ft := t.StructField.Type
	c, err := newCodec(ft, t.StructField.Tag)
	if err != nil {
		root := reflect.TypeFor[S]()
		panic(fmt.Errorf("\n==> Critical Error:\n❌ field `%s` of `%v` cannot be bound to HTTP request: %w\n\n🛠️ How To Fix:\n- use string, bool, int*, uint*, float*, time.Duration, []byte, a slice of them;\n- implement encoding.TextUnmarshaler for the type;\n- register decoder with µ.RegisterCodec;\n- add tag `content:\"json\"` or `content:\"form\"` to the field",
			t.StructField.Name, root, err))
	}

	l := &lens{
		root: reflect.PointerTo(reflect.TypeFor[S]()),
		off:  t.StructField.Offset + t.RootOffs,
		ft:   ft,
		c:    c,
	}

	if _, custom := registry.Load(ft); !custom && !implementsText(ft) && c.multi == nil {
		switch ft.Kind() {
		case reflect.String:
			l.fast = fastString
		case reflect.Int:
			if ft != durationType {
				l.fast = fastInt
			}
		case reflect.Float64:
			l.fast = fastFloat
		}
		if tagContent(t.StructField.Tag) != "" {
			l.fast = fastNone
		}
	}

	return l
}

var (
	textType     = reflect.TypeFor[encoding.TextUnmarshaler]()
	durationType = reflect.TypeFor[time.Duration]()
)

func implementsText(ft reflect.Type) bool {
	return reflect.PointerTo(ft).Implements(textType)
}
