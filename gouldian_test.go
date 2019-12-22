package gouldian_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/it"
)

func TestVerb(t *testing.T) {
	req := gouldian.NewGet("")

	it.Ok(t).
		If(gouldian.Get().IsMatch(req)).
		Should().Equal(true)
}

func TestPath(t *testing.T) {
	req := gouldian.NewGet("/foo")

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
	req := gouldian.NewGet("/foo/bar")
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
	req := gouldian.NewGet("/foo/10")
	inv := gouldian.NewGet("/foo/bar")
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

func TestOpt(t *testing.T) {
	req := gouldian.NewGet("/foo?bar=foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").Opt("bar", "foo").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").Opt("bar", "bar").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Opt("foo", "").IsMatch(req)).
		Should().Equal(false)
}

func TestIsOpt(t *testing.T) {
	req := gouldian.NewGet("/foo?bar")

	it.Ok(t).
		If(gouldian.Get().Path("foo").IsOpt("bar").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").IsOpt("foo").IsMatch(req)).
		Should().Equal(false)
}

func TestOptString(t *testing.T) {
	req := gouldian.NewGet("/foo?bar=foo")
	bar := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").OptString("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		//
		If(bar).Should().Equal("foo").
		//
		If(gouldian.Get().Path("foo").OptString("foo", &bar).IsMatch(req)).
		Should().Equal(false)
}

func TestOptInt(t *testing.T) {
	req := gouldian.NewGet("/foo?bar=10")
	inv := gouldian.NewGet("/foo?bar=foo")
	bar := 0

	it.Ok(t).
		If(gouldian.Get().Path("foo").OptInt("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		//
		If(bar).Should().Equal(10).
		//
		If(gouldian.Get().Path("foo").OptInt("foo", &bar).IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").OptInt("bar", &bar).IsMatch(inv)).
		Should().Equal(false)
}

func TestHead(t *testing.T) {
	req := gouldian.NewGet("/foo").
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

func TestIsHead(t *testing.T) {
	req := gouldian.NewGet("/foo").
		With("Content-Type", "application/json")

	it.Ok(t).
		If(gouldian.Get().Path("foo").IsHead("Content-Type").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").IsHead("Accept").IsMatch(req)).
		Should().Equal(false)
}

func TestHeadString(t *testing.T) {
	req := gouldian.NewGet("/foo").
		With("Content-Type", "application/json")
	content := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").HeadString("Content-Type", &content).IsMatch(req)).
		Should().Equal(true).
		//
		If(content).Should().Equal("application/json").
		//
		If(gouldian.Get().Path("foo").HeadString("Accept", &content).IsMatch(req)).
		Should().Equal(false)
}

type foobar struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJson(t *testing.T) {
	req := gouldian.NewGet("/foo").
		WithJson(foobar{"foo", 10})
	inv := gouldian.NewGet("/foo").
		WithText("foobar")
	val := foobar{}

	it.Ok(t).
		If(gouldian.Get().Path("foo").Json(&val).IsMatch(req)).
		Should().Equal(true).
		//
		If(val).Should().Equal(foobar{"foo", 10}).
		//
		If(gouldian.Get().Path("foo").Json(&val).IsMatch(inv)).
		Should().Equal(false)
}

func TestThen(t *testing.T) {
	req := gouldian.NewGet("/foo")
	e := gouldian.Get().Path("foo").
		Then(func() error { return gouldian.Ok().Text("bar") })

	v := e(req)
	var o *gouldian.Output

	fmt.Println("==> ", v)
	fmt.Println("==> ", errors.As(v, &o))
}
