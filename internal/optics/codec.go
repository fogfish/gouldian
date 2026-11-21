package optics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// content tag, e.g. `content:"json"`, `content:"json,strict"`
func tagContent(tag reflect.StructTag) string { return tag.Get("content") }

func newCodec(ft reflect.Type, tag reflect.StructTag) (*codec, error) {
	// 1. explicit content tag
	if content := tagContent(tag); content != "" {
		return newContentCodec(ft, content)
	}

	// 2. application registered
	if c, ok := registry.Load(ft); ok {
		return c.(*codec), nil
	}

	// 3. pointers
	if ft.Kind() == reflect.Pointer && !implementsText(ft) {
		return newPointerCodec(ft, tag)
	}

	// 4. encoding.TextUnmarshaler
	if implementsText(ft) {
		return newTextCodec(ft), nil
	}

	if ft == durationType {
		return &codec{
			str: func(s string) (Value, error) {
				v, err := time.ParseDuration(s)
				return Value{Number: int64(v)}, err
			},
			set: func(dst reflect.Value, v Value) { dst.SetInt(v.Number) },
		}, nil
	}

	switch ft.Kind() {
	case reflect.String:
		return &codec{
			str: func(s string) (Value, error) { return Value{String: s}, nil },
			set: func(dst reflect.Value, v Value) { dst.SetString(v.String) },
		}, nil
	case reflect.Bool:
		return &codec{
			str: func(s string) (Value, error) {
				v, err := strconv.ParseBool(s)
				if v {
					return Value{Number: 1}, err
				}
				return Value{}, err
			},
			set: func(dst reflect.Value, v Value) { dst.SetBool(v.Number != 0) },
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits := ft.Bits()
		return &codec{
			str: func(s string) (Value, error) {
				v, err := strconv.ParseInt(s, 10, bits)
				return Value{Number: v}, err
			},
			set: func(dst reflect.Value, v Value) { dst.SetInt(v.Number) },
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits := ft.Bits()
		return &codec{
			str: func(s string) (Value, error) {
				v, err := strconv.ParseUint(s, 10, bits)
				return Value{Number: int64(v)}, err
			},
			set: func(dst reflect.Value, v Value) { dst.SetUint(uint64(v.Number)) },
		}, nil
	case reflect.Float32, reflect.Float64:
		bits := ft.Bits()
		return &codec{
			str: func(s string) (Value, error) {
				v, err := strconv.ParseFloat(s, bits)
				return fromDouble(v), err
			},
			set: func(dst reflect.Value, v Value) { dst.SetFloat(v.Double()) },
		}, nil
	case reflect.Slice:
		if ft.Elem().Kind() == reflect.Uint8 {
			return newBytesCodec(), nil
		}
		return newSliceCodec(ft, tag)
	case reflect.Struct:
		// default for struct is JSON (compatibility)
		return newContentCodec(ft, "json")
	}

	return nil, fmt.Errorf("type %v is not supported", ft)
}

func newTextCodec(ft reflect.Type) *codec {
	return &codec{
		str: func(s string) (Value, error) {
			p := reflect.New(ft)
			if err := p.Interface().(interface{ UnmarshalText([]byte) error }).UnmarshalText([]byte(s)); err != nil {
				return Value{}, err
			}
			return Value{Ptr: p.UnsafePointer()}, nil
		},
		set: func(dst reflect.Value, v Value) { dst.Set(reflect.NewAt(ft, v.Ptr).Elem()) },
	}
}

func newBytesCodec() *codec {
	return &codec{
		str: func(s string) (Value, error) { b := []byte(s); return Value{Ptr: unsafe.Pointer(&b)}, nil },
		raw: func(b []byte) (Value, error) { return Value{Ptr: unsafe.Pointer(&b)}, nil },
		set: func(dst reflect.Value, v Value) { dst.SetBytes(*(*[]byte)(v.Ptr)) },
	}
}

func newPointerCodec(ft reflect.Type, tag reflect.StructTag) (*codec, error) {
	ec, err := newCodec(ft.Elem(), tag)
	if err != nil {
		return nil, err
	}

	return &codec{
		str:   ec.str,
		raw:   ec.raw,
		multi: ec.multi,
		set: func(dst reflect.Value, v Value) {
			p := reflect.New(ft.Elem())
			ec.set(p.Elem(), v)
			dst.Set(p)
		},
	}, nil
}

// slice of scalars, repeated `?a=1&a=2`, or tagged with `sep:","`
func newSliceCodec(ft reflect.Type, tag reflect.StructTag) (*codec, error) {
	ec, err := newCodec(ft.Elem(), "")
	if err != nil {
		return nil, fmt.Errorf("slice of %w", err)
	}
	sep := tag.Get("sep")

	decode := func(vals []string) (Value, error) {
		if sep != "" {
			var split []string
			for _, v := range vals {
				split = append(split, strings.Split(v, sep)...)
			}
			vals = split
		}

		holder := reflect.New(ft)
		slice := reflect.MakeSlice(ft, len(vals), len(vals))
		for i, s := range vals {
			ev, err := ec.str(strings.TrimSpace(s))
			if err != nil {
				return Value{}, fmt.Errorf("element %d: %w", i, err)
			}
			ec.set(slice.Index(i), ev)
		}
		holder.Elem().Set(slice)
		return Value{Ptr: holder.UnsafePointer()}, nil
	}

	return &codec{
		str:   func(s string) (Value, error) { return decode([]string{s}) },
		multi: decode,
		set:   func(dst reflect.Value, v Value) { dst.Set(reflect.NewAt(ft, v.Ptr).Elem()) },
	}, nil
}

// content codec: json, json strict, form
func newContentCodec(ft reflect.Type, content string) (*codec, error) {
	kind, opt, _ := strings.Cut(content, ",")
	switch strings.TrimSpace(kind) {
	case "json", "application/json":
		return newJSONCodec(ft, strings.TrimSpace(opt) == "strict"), nil
	case "form", "application/x-www-form-urlencoded":
		return newFormCodec(ft)
	default:
		return nil, fmt.Errorf("unknown content type `%s`", content)
	}
}

func newJSONCodec(ft reflect.Type, strict bool) *codec {
	decode := func(b []byte) (Value, error) {
		p := reflect.New(ft)
		if strict {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.DisallowUnknownFields()
			if err := dec.Decode(p.Interface()); err != nil {
				return Value{}, err
			}
			if dec.More() {
				return Value{}, errors.New("unexpected data after top-level JSON value")
			}
		} else if err := json.Unmarshal(b, p.Interface()); err != nil {
			return Value{}, err
		}
		return Value{Ptr: p.UnsafePointer()}, nil
	}

	return &codec{
		str: func(s string) (Value, error) { return decode([]byte(s)) },
		raw: decode,
		set: func(dst reflect.Value, v Value) { dst.Set(reflect.NewAt(ft, v.Ptr).Elem()) },
	}
}

// form codec decodes `application/x-www-form-urlencoded` into struct.
// Field name is taken from tag `form`, then `json`, then name of the field.
func newFormCodec(ft reflect.Type) (*codec, error) {
	if ft.Kind() == reflect.Pointer {
		ft = ft.Elem()
	}
	if ft.Kind() != reflect.Struct {
		return nil, fmt.Errorf("form content requires struct, got %v", ft)
	}

	type field struct {
		name  string
		index []int
		c     *codec
	}
	var fields []field
	for i := 0; i < ft.NumField(); i++ {
		f := ft.Field(i)
		if !f.IsExported() {
			continue
		}
		name := f.Name
		for _, key := range []string{"json", "form"} {
			if v, _, _ := strings.Cut(f.Tag.Get(key), ","); v == "-" {
				name = ""
			} else if v != "" {
				name = v
			}
		}
		if name == "" {
			continue
		}

		c, err := newCodec(f.Type, f.Tag)
		if err != nil {
			return nil, fmt.Errorf("form field %s: %w", f.Name, err)
		}
		fields = append(fields, field{name: name, index: f.Index, c: c})
	}

	decode := func(s string) (Value, error) {
		q, err := url.ParseQuery(s)
		if err != nil {
			return Value{}, err
		}

		p := reflect.New(ft)
		for _, f := range fields {
			vals, ok := q[f.name]
			if !ok || len(vals) == 0 {
				continue
			}

			var v Value
			if f.c.multi != nil {
				v, err = f.c.multi(vals)
			} else {
				v, err = f.c.str(vals[0])
			}
			if err != nil {
				return Value{}, fmt.Errorf("%s: %w", f.name, err)
			}
			f.c.set(p.Elem().FieldByIndex(f.index), v)
		}
		return Value{Ptr: p.UnsafePointer()}, nil
	}

	return &codec{
		str: decode,
		raw: func(b []byte) (Value, error) { return decode(string(b)) },
		set: func(dst reflect.Value, v Value) {
			if dst.Kind() == reflect.Pointer {
				dst.Set(reflect.NewAt(ft, v.Ptr))
				return
			}
			dst.Set(reflect.NewAt(ft, v.Ptr).Elem())
		},
	}, nil
}
