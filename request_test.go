package gouldian_test

import (
	"errors"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/it"
)

func TestVerbDelete(t *testing.T) {
	it.Ok(t).
		If(
			gouldian.Delete().IsMatch(gouldian.MockVerb("DELETE", "")),
		).Should().Equal(true).
		If(
			gouldian.Delete().IsMatch(gouldian.MockVerb("XXX", "")),
		).Should().Equal(false)
}

func TestVerbGet(t *testing.T) {
	it.Ok(t).
		If(
			gouldian.Get().IsMatch(gouldian.MockVerb("GET", "")),
		).Should().Equal(true).
		If(
			gouldian.Get().IsMatch(gouldian.MockVerb("XXX", "")),
		).Should().Equal(false)
}

func TestVerbPatch(t *testing.T) {
	it.Ok(t).
		If(
			gouldian.Patch().IsMatch(gouldian.MockVerb("PATCH", "")),
		).Should().Equal(true).
		If(
			gouldian.Patch().IsMatch(gouldian.MockVerb("XXX", "")),
		).Should().Equal(false)
}

func TestVerbPost(t *testing.T) {
	it.Ok(t).
		If(
			gouldian.Post().IsMatch(gouldian.MockVerb("POST", "")),
		).Should().Equal(true).
		If(
			gouldian.Post().IsMatch(gouldian.MockVerb("XXX", "")),
		).Should().Equal(false)
}

func TestVerbPut(t *testing.T) {
	it.Ok(t).
		If(
			gouldian.Put().IsMatch(gouldian.MockVerb("PUT", "")),
		).Should().Equal(true).
		If(
			gouldian.Put().IsMatch(gouldian.MockVerb("XXX", "")),
		).Should().Equal(false)
}

func TestPath(t *testing.T) {
	req := gouldian.Mock("/foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("bar").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Path("bar").IsMatch(req)).
		Should().Equal(false)
}

func TestString(t *testing.T) {
	req := gouldian.Mock("/foo/bar")
	foo := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").String(&foo).IsMatch(req)).
		Should().Equal(true).
		//
		If(foo).Should().Equal("bar").
		//
		If(gouldian.Get().Path("foo").Path("bar").String(&foo).IsMatch(req)).
		Should().Equal(false)
}

func TestInt(t *testing.T) {
	req := gouldian.Mock("/foo/10")
	inv := gouldian.Mock("/foo/bar")
	foo := 0

	it.Ok(t).
		If(gouldian.Get().Path("foo").Int(&foo).IsMatch(req)).
		Should().Equal(true).
		//
		If(foo).Should().Equal(10).
		//
		If(gouldian.Get().Path("foo").Path("bar").Int(&foo).IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Int(&foo).IsMatch(inv)).
		Should().Equal(false)
}

func TestParam(t *testing.T) {
	req := gouldian.Mock("/foo?bar=foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").Param("bar", "foo").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").Param("bar", "bar").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Param("foo", "").IsMatch(req)).
		Should().Equal(false)
}

func TestHasParam(t *testing.T) {
	req := gouldian.Mock("/foo?bar")
	foo := gouldian.Mock("/foo?bar=foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").HasParam("bar").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasParam("bar").IsMatch(foo)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasParam("foo").IsMatch(req)).
		Should().Equal(false)
}

func TestQString(t *testing.T) {
	req := gouldian.Mock("/foo?bar=foo")
	bar := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").QString("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal("foo").
		//
		If(gouldian.Get().Path("foo").QString("foo", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal("")
}

func TestQInt(t *testing.T) {
	req := gouldian.Mock("/foo?bar=10")
	inv := gouldian.Mock("/foo?bar=foo")
	bar := 0

	it.Ok(t).
		If(gouldian.Get().Path("foo").QInt("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal(10).
		//
		If(gouldian.Get().Path("foo").QInt("foo", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal(0).
		//
		If(gouldian.Get().Path("foo").QInt("bar", &bar).IsMatch(inv)).
		Should().Equal(false)
}

func TestHead(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")

	it.Ok(t).
		If(gouldian.Get().Path("foo").Head("Content-Type", "application/json").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").Head("Content-Type", "text/plain").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Head("Accept", "application/json").IsMatch(req)).
		Should().Equal(false)
}

func TestHasHead(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")

	it.Ok(t).
		If(gouldian.Get().Path("foo").HasHead("Content-Type").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasHead("Accept").IsMatch(req)).
		Should().Equal(false)
}

func TestHeadString(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")
	content := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").HString("Content-Type", &content).IsMatch(req)).
		Should().Equal(true).
		If(content).Should().Equal("application/json").
		//
		If(gouldian.Get().Path("foo").HString("Accept", &content).IsMatch(req)).
		Should().Equal(true).
		If(content).Should().Equal("")
}

type foobar struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJson(t *testing.T) {
	req := gouldian.Mock("/foo").
		WithJSON(foobar{"foo", 10})
	inv := gouldian.Mock("/foo").
		WithText("foobar")
	val := foobar{}

	it.Ok(t).
		If(gouldian.Get().Path("foo").JSON(&val).IsMatch(req)).
		Should().Equal(true).
		//
		If(val).Should().Equal(foobar{"foo", 10}).
		//
		If(gouldian.Get().Path("foo").JSON(&val).IsMatch(inv)).
		Should().Equal(false)
}

func TestThenSuccess(t *testing.T) {
	req := gouldian.Mock("/foo")
	handle := func() error { return gouldian.Ok().Text("bar") }

	it.Ok(t).
		If(gouldian.Get().Path("foo").FMap(handle)(req)).Should().
		Assert(
			func(be interface{}) bool {
				var rsp *gouldian.Output
				return errors.As(be.(error), &rsp)
			},
		)
}

func TestThenFailure(t *testing.T) {
	req := gouldian.Mock("/foo")
	handle := func() error { return gouldian.Unauthorized("") }

	it.Ok(t).
		If(gouldian.Get().Path("foo").FMap(handle)(req)).Should().
		Assert(
			func(be interface{}) bool {
				var rsp *gouldian.Output
				return !errors.As(be.(error), &rsp)
			},
		)
}
