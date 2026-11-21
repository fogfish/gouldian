package optics_test

import (
	"errors"
	"net/netip"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/fogfish/golem/hseq"
	"github.com/fogfish/gouldian/v3/internal/optics"
	"github.com/fogfish/it/v2"
)

// bind the first field of T to the string and put it
func bind[T any](t *testing.T, field, given string) (T, error) {
	t.Helper()

	var v T
	l := hseq.FMap1(hseq.New[T](field), optics.NewLens[T])
	val, err := l.FromString(given)
	if err != nil {
		return v, err
	}
	err = optics.Morph(optics.Morphisms{{Lens: l, Value: val}}, &v)
	return v, err
}

func TestScalars(t *testing.T) {
	type String string
	type Int int
	type Float float64

	type T struct {
		S  string
		NS String
		I  int
		NI Int
		I8 int8
		U  uint16
		F  float64
		NF Float
		F3 float32
		B  bool
		D  time.Duration
		T  time.Time
		PS *string
		PI *int
		PB *bool
	}

	for field, given := range map[string]string{
		"S": "abc", "NS": "abc", "I": "-100", "NI": "100", "I8": "127", "U": "65535",
		"F": "100.101", "NF": "100.101", "F3": "1.5", "B": "true", "D": "1m30s",
		"T": "2023-02-01T10:20:30Z", "PS": "abc", "PI": "10", "PB": "false",
	} {
		_, err := bind[T](t, field, given)
		it.Then(t).Should(it.Nil(err))
	}

	v, _ := bind[T](t, "S", "abc")
	it.Then(t).Should(it.Equal(v.S, "abc"))
	v, _ = bind[T](t, "I", "-100")
	it.Then(t).Should(it.Equal(v.I, -100))
	v, _ = bind[T](t, "I8", "127")
	it.Then(t).Should(it.Equal(v.I8, int8(127)))
	v, _ = bind[T](t, "U", "65535")
	it.Then(t).Should(it.Equal(v.U, uint16(65535)))
	v, _ = bind[T](t, "F", "100.101")
	it.Then(t).Should(it.Equal(v.F, 100.101))
	v, _ = bind[T](t, "F3", "1.5")
	it.Then(t).Should(it.Equal(v.F3, float32(1.5)))
	v, _ = bind[T](t, "B", "true")
	it.Then(t).Should(it.True(v.B))
	v, _ = bind[T](t, "D", "1m30s")
	it.Then(t).Should(it.Equal(v.D, 90*time.Second))
	v, _ = bind[T](t, "T", "2023-02-01T10:20:30Z")
	it.Then(t).Should(it.True(v.T.Equal(time.Date(2023, 2, 1, 10, 20, 30, 0, time.UTC))))
	v, _ = bind[T](t, "PS", "abc")
	it.Then(t).Should(it.Equal(*v.PS, "abc"))
	v, _ = bind[T](t, "PB", "false")
	it.Then(t).Should(it.Equal(*v.PB, false))
}

func TestScalarsFail(t *testing.T) {
	type T struct {
		I  int
		I8 int8
		U  uint
		F  float64
		B  bool
		D  time.Duration
		T  time.Time
	}

	for field, given := range map[string]string{
		"I": "abc", "I8": "128", "U": "-1", "F": "abc", "B": "maybe", "D": "1x", "T": "yesterday",
	} {
		_, err := bind[T](t, field, given)
		it.Then(t).ShouldNot(it.Nil(err))
	}
}

type Even int

func (e *Even) UnmarshalText(b []byte) error {
	var n int
	for _, c := range string(b) {
		if c < '0' || c > '9' {
			return errors.New("not a number")
		}
		n = n*10 + int(c-'0')
	}
	if n%2 != 0 {
		return errors.New("must be even")
	}
	*e = Even(n)
	return nil
}

func TestTextUnmarshaler(t *testing.T) {
	type T struct {
		E  Even
		PE *Even
		IP netip.Addr
	}

	v, err := bind[T](t, "E", "10")
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(v.E, Even(10)))

	// the domain type validates itself
	_, err = bind[T](t, "E", "11")
	it.Then(t).ShouldNot(it.Nil(err))
	it.Then(t).Should(it.True(strings.Contains(err.Error(), "must be even")))

	pv, err := bind[T](t, "PE", "4")
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(*pv.PE, Even(4)))

	ip, err := bind[T](t, "IP", "10.0.0.1")
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(ip.IP.String(), "10.0.0.1"))
}

func TestRawBytes(t *testing.T) {
	type T struct {
		B []byte
		R rawMsg
	}

	var v T
	l := hseq.FMap1(hseq.New[T]("B"), optics.NewLens[T])
	val, err := l.FromBytes([]byte{0, 1, 2, 255})
	it.Then(t).Should(it.Nil(err))
	it.Then(t).Should(it.Nil(optics.Morph(optics.Morphisms{{Lens: l, Value: val}}, &v)))
	it.Then(t).Should(it.Equiv(v.B, []byte{0, 1, 2, 255}))

	r := hseq.FMap1(hseq.New[T]("R"), optics.NewLens[T])
	val, _ = r.FromBytes([]byte(`{"a":1}`))
	it.Then(t).Should(it.Nil(optics.Morph(optics.Morphisms{{Lens: r, Value: val}}, &v)))
	it.Then(t).Should(it.Equal(string(v.R), `{"a":1}`))
}

type rawMsg []byte

func TestSlices(t *testing.T) {
	type T struct {
		S []string
		I []int
		C []string `sep:","`
	}

	var v T
	l := hseq.FMap1(hseq.New[T]("S"), optics.NewLens[T])
	it.Then(t).Should(it.True(l.Multi()))
	val, err := l.FromStrings([]string{"a", "b"})
	it.Then(t).Should(it.Nil(err))
	it.Then(t).Should(it.Nil(optics.Morph(optics.Morphisms{{Lens: l, Value: val}}, &v)))
	it.Then(t).Should(it.Equiv(v.S, []string{"a", "b"}))

	li := hseq.FMap1(hseq.New[T]("I"), optics.NewLens[T])
	_, err = li.FromStrings([]string{"1", "x"})
	it.Then(t).ShouldNot(it.Nil(err))

	lc := hseq.FMap1(hseq.New[T]("C"), optics.NewLens[T])
	val, _ = lc.FromString("a, b,c")
	it.Then(t).Should(it.Nil(optics.Morph(optics.Morphisms{{Lens: lc, Value: val}}, &v)))
	it.Then(t).Should(it.Equiv(v.C, []string{"a", "b", "c"}))
}

type J struct {
	X string `json:"x"`
}

func TestJSON(t *testing.T) {
	type T struct {
		A J
		P *J
		S J `content:"json,strict"`
	}

	v, err := bind[T](t, "A", `{"x":"abc"}`)
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(v.A, J{"abc"}))

	pv, err := bind[T](t, "P", `{"x":"abc"}`)
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(pv.P.X, "abc"))

	_, err = bind[T](t, "A", `{"x":`)
	it.Then(t).ShouldNot(it.Nil(err))

	// lenient by default, strict rejects unknown fields
	_, err = bind[T](t, "A", `{"x":"a","z":1}`)
	it.Then(t).Should(it.Nil(err))
	_, err = bind[T](t, "S", `{"x":"a","z":1}`)
	it.Then(t).ShouldNot(it.Nil(err))
}

func TestForm(t *testing.T) {
	type F struct {
		Foo  string   `json:"foo"`
		Bar  int      `form:"bar"`
		Tags []string `form:"tags"`
	}
	type T struct {
		F F `content:"form"`
	}

	v, err := bind[T](t, "F", "foo=a&bar=10&tags=x&tags=y")
	it.Then(t).Should(it.Nil(err)).
		Should(it.Equiv(v.F, F{Foo: "a", Bar: 10, Tags: []string{"x", "y"}}))

	_, err = bind[T](t, "F", "bar=x")
	it.Then(t).ShouldNot(it.Nil(err))
}

func TestRegisteredCodec(t *testing.T) {
	type Color struct{ RGB string }
	type T struct{ C Color }

	optics.Register(func(s string) (Color, error) {
		if !strings.HasPrefix(s, "#") {
			return Color{}, errors.New("no hash")
		}
		return Color{s[1:]}, nil
	})

	v, err := bind[T](t, "C", "#ff00ff")
	it.Then(t).Should(it.Nil(err)).Should(it.Equal(v.C.RGB, "ff00ff"))
	_, err = bind[T](t, "C", "ff00ff")
	it.Then(t).ShouldNot(it.Nil(err))
}

func TestUnsupportedPanicsAtConstruction(t *testing.T) {
	type T struct{ M map[string]string }

	defer func() {
		r := recover()
		it.Then(t).ShouldNot(it.Nil(r))
		err, _ := r.(error)
		it.Then(t).Should(it.True(strings.Contains(err.Error(), "`M`")))
	}()
	hseq.FMap1(hseq.New[T]("M"), optics.NewLens[T])
}

func TestPutWrongType(t *testing.T) {
	type T struct{ A string }
	type X struct{ A string }

	l := hseq.FMap1(hseq.New[T]("A"), optics.NewLens[T])
	val, _ := l.FromString("a")
	err := l.Put(&X{}, val)
	it.Then(t).ShouldNot(it.Nil(err))
	_ = reflect.TypeOf(err)
}

func TestEmbedded(t *testing.T) {
	type Base struct{ ID string }
	type T struct {
		Base
		N int
	}

	var v T
	l := hseq.FMap1(hseq.New[T]("ID"), optics.NewLens[T])
	val, _ := l.FromString("x")
	it.Then(t).Should(it.Nil(l.Put(&v, val))).Should(it.Equal(v.ID, "x"))
}

// Value is copied by value on the hot path, its size is performance critical,
// see the comment of the type before changing it.
func TestValueSize(t *testing.T) {
	it.Then(t).
		Should(it.Equal(unsafe.Sizeof(optics.Value{}), uintptr(32))).
		Should(it.Equal(unsafe.Sizeof(optics.Morphism{}), uintptr(48)))
}
