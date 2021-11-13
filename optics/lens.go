/*

Package optics ...

*/
package optics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

/*

Codec ...
*/
type Codec interface {
	FromString(string) (interface{}, error)
	FromSeq([]string) (interface{}, error)
}

/*

Lens ...
*/
type Lens interface {
	Codec
	Put(a reflect.Value, s interface{}) error
}

/*

Morphism ... stack of lenses and values to apply
*/
type Morphism map[Lens]interface{}

/*

Apply ...
*/
func (m Morphism) Apply(a interface{}) error {
	g := reflect.ValueOf(a)
	if g.Kind() != reflect.Ptr {
		return fmt.Errorf("Morphism requires pointer type, %s given", g.Kind().String())
	}

	for lens, s := range m {
		if err := lens.Put(g, s); err != nil {
			return err
		}
	}

	return nil
}

/*

lensStruct ...
*/
type lensStruct struct {
	field  int
	typeof reflect.Type
}

/*

lensStructString ...
*/
type lensStructString struct{ lensStruct }

func (lensStructString) FromString(s string) (interface{}, error) {
	return s, nil
}

func (lens lensStructString) FromSeq(s []string) (interface{}, error) {
	if len(s) == 0 {
		return "", nil
	}

	return lens.FromString(s[0])
}

func (lens lensStructString) Put(a reflect.Value, s interface{}) error {
	f := a.Elem().Field(int(lens.field))

	if f.Kind() == reflect.Ptr {
		return lens.putToPtr(f, s)
	}

	return lens.putToVal(f, s)
}

func (lens lensStructString) putToVal(a reflect.Value, s interface{}) error {
	switch v := s.(type) {
	case string:
		a.SetString(v)
	case *string:
		a.SetString(*v)
	}
	return nil
}

func (lens lensStructString) putToPtr(a reflect.Value, s interface{}) error {
	switch v := s.(type) {
	case string:
		p := reflect.New(lens.typeof.Elem())
		p.Elem().SetString(s.(string))
		a.Set(p)
	case *string:
		a.Set(reflect.ValueOf(v))
	}
	return nil
}

/*

lensStructInt ...
*/
type lensStructInt struct{ lensStruct }

func (lensStructInt) FromString(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

func (lens lensStructInt) FromSeq(s []string) (interface{}, error) {
	if len(s) == 0 {
		return 0, nil
	}

	return lens.FromString(s[0])
}

func (lens lensStructInt) Put(a reflect.Value, s interface{}) error {
	a.Elem().Field(int(lens.field)).SetInt(int64(s.(int)))
	return nil
}

/*

lensStructFloat ...
*/
type lensStructFloat struct{ lensStruct }

func (lensStructFloat) FromString(s string) (interface{}, error) {
	return strconv.ParseFloat(s, 64)
}

func (lens lensStructFloat) FromSeq(s []string) (interface{}, error) {
	if len(s) == 0 {
		return 0.0, nil
	}

	return lens.FromString(s[0])
}

func (lens lensStructFloat) Put(a reflect.Value, s interface{}) error {
	a.Elem().Field(int(lens.field)).SetFloat(s.(float64))
	return nil
}

/*

lensStructJSON ...
*/
type lensStructJSON struct{ lensStruct }

func (lensStructJSON) FromString(s string) (interface{}, error) {
	return s, nil
}

func (lens lensStructJSON) FromSeq(s []string) (interface{}, error) {
	if len(s) == 0 {
		return "", nil
	}

	return lens.FromString(s[0])
}

func (lens lensStructJSON) Put(a reflect.Value, s interface{}) error {
	c := reflect.New(lens.typeof)
	o := c.Interface()
	switch v := s.(type) {
	case []byte:
		if err := json.Unmarshal(v, &o); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), &o); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Cannot cast %T to JSON", s)
	}

	a.Elem().Field(int(lens.field)).Set(c.Elem())
	return nil
}

/*

lensStructSeq ...
*/
type lensStructSeq struct{ lensStruct }

func (lensStructSeq) FromString(s string) (interface{}, error) {
	return s, nil
}

func (lensStructSeq) FromSeq(s []string) (interface{}, error) {
	return s, nil
}

func (lens lensStructSeq) Put(a reflect.Value, s interface{}) error {
	v := reflect.ValueOf(s)
	switch v.Type().Kind() {
	case reflect.Slice:
		a.Elem().Field(int(lens.field)).Set(v)
	default:
		return fmt.Errorf("Cannot cast %T to Seq", s)
	}

	return nil
}

/*

newLensStruct ...
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
		return &lensStructJSON{lensStruct{id, field.Type}}
	case reflect.Slice:
		return &lensStructSeq{lensStruct{id, field.Type}}
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

Lenses1 ...
*/
func Lenses1(t interface{}) Lens {
	tc := typeOf(t)
	if tc.NumField() != 1 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 1", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0))
}

/*

Lenses2 ...
*/
func Lenses2(t interface{}) (Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 2 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 2", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1))
}

/*

Lenses3 ...
*/
func Lenses3(t interface{}) (Lens, Lens, Lens) {
	tc := typeOf(t)
	if tc.NumField() != 3 {
		panic(fmt.Errorf("Unable to unapply type |%s| = %d to lens of 3", tc.Name(), tc.NumField()))
	}

	return newLensStruct(0, tc.Field(0)),
		newLensStruct(1, tc.Field(1)),
		newLensStruct(2, tc.Field(2))
}

/*

Lenses4 ...
*/
func Lenses4(t interface{}) (Lens, Lens, Lens, Lens) {
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

Lenses5 ...
*/
func Lenses5(t interface{}) (Lens, Lens, Lens, Lens, Lens) {
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

Lenses6 ...
*/
func Lenses6(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens) {
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

Lenses7 ...
*/
func Lenses7(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
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

Lenses8 ...
*/
func Lenses8(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
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

Lenses9 ...
*/
func Lenses9(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
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

Lenses10 ...
*/
func Lenses10(t interface{}) (Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens, Lens) {
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
