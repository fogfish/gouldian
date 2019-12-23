package gouldian_test

import (
	"errors"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/it"
)

func TestEndpointThen(t *testing.T) {
	var ok = errors.New("b")
	var a gouldian.Endpoint = func(x *gouldian.Input) error { return nil }
	var b gouldian.Endpoint = func(x *gouldian.Input) error { return ok }

	it.Ok(t).
		If(a.Then(b)(gouldian.Mock(""))).Should().Equal(ok)
}

func TestEndpointOr(t *testing.T) {
	var ok = errors.New("a")
	var a gouldian.Endpoint = func(x *gouldian.Input) error { return ok }
	var b gouldian.Endpoint = func(x *gouldian.Input) error { return nil }

	it.Ok(t).
		If(a.Or(b)(gouldian.Mock(""))).Should().Equal(ok)
}
