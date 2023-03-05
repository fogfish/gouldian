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

package emitter_test

import (
	"net/http"
	"reflect"
	"testing"

	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/gouldian/v2/mock"
	ø "github.com/fogfish/gouldian/v2/output"
	"github.com/fogfish/it"
)

func TestSuccess(t *testing.T) {
	//
	output(t, ø.Status.OK(), µ.NewOutput(http.StatusOK))
	output(t, ø.Status.Created(), µ.NewOutput(http.StatusCreated))
	output(t, ø.Status.Accepted(), µ.NewOutput(http.StatusAccepted))
	output(t, ø.Status.NonAuthoritativeInfo(), µ.NewOutput(http.StatusNonAuthoritativeInfo))
	output(t, ø.Status.NoContent(), µ.NewOutput(http.StatusNoContent))
	output(t, ø.Status.ResetContent(), µ.NewOutput(http.StatusResetContent))

	//
	output(t, ø.Status.MultipleChoices(), µ.NewOutput(http.StatusMultipleChoices))
	output(t, ø.Status.MovedPermanently(), µ.NewOutput(http.StatusMovedPermanently))
	output(t, ø.Status.Found(), µ.NewOutput(http.StatusFound))
	output(t, ø.Status.SeeOther(), µ.NewOutput(http.StatusSeeOther))
	output(t, ø.Status.NotModified(), µ.NewOutput(http.StatusNotModified))
	output(t, ø.Status.UseProxy(), µ.NewOutput(http.StatusUseProxy))
	output(t, ø.Status.TemporaryRedirect(), µ.NewOutput(http.StatusTemporaryRedirect))
	output(t, ø.Status.PermanentRedirect(), µ.NewOutput(http.StatusPermanentRedirect))

	//
	output(t, ø.Status.BadRequest(), µ.NewOutput(http.StatusBadRequest))
	output(t, ø.Status.Unauthorized(), µ.NewOutput(http.StatusUnauthorized))
	output(t, ø.Status.PaymentRequired(), µ.NewOutput(http.StatusPaymentRequired))
	output(t, ø.Status.Forbidden(), µ.NewOutput(http.StatusForbidden))
	output(t, ø.Status.NotFound(), µ.NewOutput(http.StatusNotFound))
	output(t, ø.Status.MethodNotAllowed(), µ.NewOutput(http.StatusMethodNotAllowed))
	output(t, ø.Status.NotAcceptable(), µ.NewOutput(http.StatusNotAcceptable))
	output(t, ø.Status.ProxyAuthRequired(), µ.NewOutput(http.StatusProxyAuthRequired))
	output(t, ø.Status.RequestTimeout(), µ.NewOutput(http.StatusRequestTimeout))
	output(t, ø.Status.Conflict(), µ.NewOutput(http.StatusConflict))
	output(t, ø.Status.Gone(), µ.NewOutput(http.StatusGone))
	output(t, ø.Status.LengthRequired(), µ.NewOutput(http.StatusLengthRequired))
	output(t, ø.Status.PreconditionFailed(), µ.NewOutput(http.StatusPreconditionFailed))
	output(t, ø.Status.RequestEntityTooLarge(), µ.NewOutput(http.StatusRequestEntityTooLarge))
	output(t, ø.Status.RequestURITooLong(), µ.NewOutput(http.StatusRequestURITooLong))
	output(t, ø.Status.UnsupportedMediaType(), µ.NewOutput(http.StatusUnsupportedMediaType))

	//
	output(t, ø.Status.InternalServerError(), µ.NewOutput(http.StatusInternalServerError))
	output(t, ø.Status.NotImplemented(), µ.NewOutput(http.StatusNotImplemented))
	output(t, ø.Status.BadGateway(), µ.NewOutput(http.StatusBadGateway))
	output(t, ø.Status.ServiceUnavailable(), µ.NewOutput(http.StatusServiceUnavailable))
	output(t, ø.Status.GatewayTimeout(), µ.NewOutput(http.StatusGatewayTimeout))
	output(t, ø.Status.HTTPVersionNotSupported(), µ.NewOutput(http.StatusHTTPVersionNotSupported))
}

type myT struct {
	A string
}

func TestWithJSON(t *testing.T) {
	output(t,
		ø.Status.OK(ø.Send(myT{"Hello"})),
		ø.Status.OK(ø.Send(myT{"Hello"})),
	)
}

func TestWithText(t *testing.T) {
	output(t,
		ø.Status.OK(ø.Send("Hello")),
		ø.Status.OK(ø.Send("Hello")),
	)
}

func TestWithBytes(t *testing.T) {
	output(t,
		ø.Status.OK(ø.Send([]byte("Hello"))),
		ø.Status.OK(ø.Send([]byte("Hello"))),
	)
}

func TestErrorOnJSON(t *testing.T) {
	output := ø.Status.OK(ø.Send(make(chan int))).(*µ.Output)

	it.Ok(t).
		If(output.Status).Should().Equal(http.StatusInternalServerError)
}

func output(t *testing.T, a, b error) {
	t.Helper()
	foo := func(*µ.Context) error { return a }
	req := mock.Input()

	it.Ok(t).
		If(foo(req)).Should().
		Assert(
			func(be interface{}) bool {
				switch out := be.(type) {
				case error:
					eq := reflect.DeepEqual(b, out)
					a.(*µ.Output).Free()
					b.(*µ.Output).Free()
					return eq
				}

				return false
			},
		)
}

func TestHeaderOutput(t *testing.T) {
	out := ø.Status.OK(
		ø.Header("foo", "bar"),
	).(*µ.Output)

	it.Ok(t).
		If(out.Status).Should().Equal(200) //.
	// If(out.Headers["foo"]).Should().Equal("bar")
}
