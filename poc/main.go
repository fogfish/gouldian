package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Body struct {
	A string `json:"a"`
}

//
type Req struct {
	Foo string
	Bar Body
}

// type XStruct interface {
// 	XX()
// }

// type AX struct{ A int }

// func (AX) XX() {}

// type BX struct{ A int }

// func (BX) XX() {}

type XStruct struct{ A int }
type AX XStruct
type BX XStruct

/*
type Lens interface {
	TypeOf() reflect.Type
}

type LensStruct struct {
	typeof reflect.Type
	field int
}

type LensStructString struct { LensStruct }

type (lens LensStructString) Put(x reflect.Value, s string)

...

*/
type Lens int

func (l Lens) String(x reflect.Value, s string) {
	x.Elem().Field(int(l)).SetString(s)
}

func (l Lens) JSON(x reflect.Value, raw []byte) {
	typeof := reflect.TypeOf(Req{})
	objType := typeof.Field(1).Type
	// reflect.TypeOf(Body{})

	obj := reflect.New(objType)
	o := obj.Interface()

	// o := x.Elem().Field(int(l)).Interface()
	if err := json.Unmarshal(raw, &o); err != nil {
		fmt.Println(err)
	}
	x.Elem().Field(int(l)).Set(obj.Elem())
	fmt.Println(o)
}

func Lens1(t interface{}) Lens {
	// 	typeof := reflect.TypeOf(t)
	// 	if typeof.Kind() == reflect.Ptr {
	// 		typeof = typeof.Elem()
	// 	}

	// v := reflect.ValueOf(t)

	// f := typeof.Field(0)
	// f.Type

	return Lens(0)

	// func(x reflect.Value) {
	// 	x.Elem().Field(0).SetString("xxx")
	// }
}

func Lens2(t interface{}) (Lens, Lens) {
	return Lens(0), Lens(1)
}

func main() {
	a, b := Lens2(Req{})
	var z Req
	v := reflect.ValueOf(&z)
	a.String(v, "abc")
	b.JSON(v, []byte("{\"a\":\"zzz\"}"))

	// example of JSON
	// var z Req
	// v := reflect.ValueOf(&z)
	// o := v.Interface()
	// json.Unmarshal([]byte("{\"a\":\"zzz\"}"), &o)

	// f(v)

	fmt.Println(z)

	a1 := &AX{1}
	b1 := &BX{2}

	seq := []*XStruct{(*XStruct)(a1), (*XStruct)(b1)}

	for _, x := range seq {
		fmt.Printf("==>>> %T %v \n", x, x)
	}

}

/*
With Context -- this is brilliant -- super brilliant!!!

type Req struct {
	Foo string
	Bar int
}
var FOO, BAR := Lenses(Req{})

// internally values
// seq of
//   {Lens, typed value}

return µ.GET(
	µ.Path(path.Is("status"), path.Int(FOO)),
  µ.Param(param.String("foo", BAR))

	µ.FMap(
		func(ctx µ.Context) error {
			var req Req
			ctx.Values(&req)

			return gouldian.Success(code)
		},
	),

)

Split the arch for local host and serverless

*/

/*

- How to define a service ?

type XYZ struct {
	*Service
	code string
	foox string
}

func (x XYZ) MyEndpoint() Endpoint {
	return µ.GET(

	µ.Path(path.Is("status"), path.Int(&x.code)),
  µ.Param(param.String("foo", &x.foox))

	µ.FMap(
		func() error {
			return gouldian.Success(code)
		},
	),
)

}


func Endpoints() {
	return []Endpoints{ &x{}, &b{}, ... }
}


func (Service) XYZ() {
	return &x{*Service}
}

*/

/*
type XYZ struct {
	code string
	foox string
}

func (XYZ) Int() Endpoint { ... }

return µ.GET(
	µ.Context(XYZ{})

	µ.Path(path.Is("status"), path.Int(ID)),
  µ.Param(param.String("foo", FOO))

	µ.FMap(
		func(ctx µ.Context) error {
			code := ctx.Int(ID)
			foox := ctx.String(FOO)

			return gouldian.Success(code)
		},
	),

)




*/
