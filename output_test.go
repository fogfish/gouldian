package gouldian_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/it"
)

func TestSuccess(t *testing.T) {
	output(t, gouldian.Success(200), gouldian.Success(200))
	output(t, gouldian.Ok(), gouldian.Success(200))
	output(t, gouldian.Created(), gouldian.Success(201))
	output(t, gouldian.Accepted(), gouldian.Success(202))
	output(t, gouldian.NoContent(), gouldian.Success(204))
}

func TestIssue(t *testing.T) {
	issue(t, gouldian.Failure(500), gouldian.Failure(500))
	issue(t, gouldian.BadRequest("issue"), gouldian.Failure(400).Reason("issue"))
	issue(t, gouldian.Unauthorized("issue"), gouldian.Failure(401).Reason("issue"))
	issue(t, gouldian.Forbidden("issue"), gouldian.Failure(403).Reason("issue"))
	issue(t, gouldian.NotFound("issue"), gouldian.Failure(404).Reason("issue"))
	issue(t, gouldian.InternalServerError("issue"), gouldian.Failure(500).Reason("issue"))
	issue(t, gouldian.NotImplemented("issue"), gouldian.Failure(501).Reason("issue"))
	issue(t, gouldian.ServiceUnavailable("issue"), gouldian.Failure(503).Reason("issue"))
}

type myT struct {
	A string
}

func TestJSON(t *testing.T) {
	output(t,
		gouldian.Ok().JSON(myT{"Hello"}),
		gouldian.Ok().JSON(myT{"Hello"}),
	)
}

func output(t *testing.T, a, b *gouldian.Output) {
	t.Helper()
	handle := func() error { return a }

	it.Ok(t).
		If(gouldian.Get().Then(handle)(gouldian.Mock(""))).Should().
		Assert(
			func(be interface{}) bool {
				var out *gouldian.Output
				if errors.As(be.(error), &out) {
					return reflect.DeepEqual(b, out)
				}
				return false
			},
		)
}

func issue(t *testing.T, a, b *gouldian.Issue) {
	t.Helper()
	handle := func() error { return a }

	it.Ok(t).
		If(gouldian.Get().Then(handle)(gouldian.Mock(""))).Should().
		Assert(
			func(be interface{}) bool {
				var out *gouldian.Issue
				if errors.As(be.(error), &out) {
					return reflect.DeepEqual(b, out)
				}
				return false
			},
		)
}
