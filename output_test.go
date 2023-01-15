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

package gouldian_test

import (
	"net/http"
	"reflect"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestSuccess(t *testing.T) {
	//
	output(t, µ.Status.OK(), µ.NewOutput(http.StatusOK))
	output(t, µ.Status.Created(), µ.NewOutput(http.StatusCreated))
	output(t, µ.Status.Accepted(), µ.NewOutput(http.StatusAccepted))
	output(t, µ.Status.NonAuthoritativeInfo(), µ.NewOutput(http.StatusNonAuthoritativeInfo))
	output(t, µ.Status.NoContent(), µ.NewOutput(http.StatusNoContent))
	output(t, µ.Status.ResetContent(), µ.NewOutput(http.StatusResetContent))

	//
	output(t, µ.Status.MultipleChoices(), µ.NewOutput(http.StatusMultipleChoices))
	output(t, µ.Status.MovedPermanently(), µ.NewOutput(http.StatusMovedPermanently))
	output(t, µ.Status.Found(), µ.NewOutput(http.StatusFound))
	output(t, µ.Status.SeeOther(), µ.NewOutput(http.StatusSeeOther))
	output(t, µ.Status.NotModified(), µ.NewOutput(http.StatusNotModified))
	output(t, µ.Status.UseProxy(), µ.NewOutput(http.StatusUseProxy))
	output(t, µ.Status.TemporaryRedirect(), µ.NewOutput(http.StatusTemporaryRedirect))
	output(t, µ.Status.PermanentRedirect(), µ.NewOutput(http.StatusPermanentRedirect))

	//
	output(t, µ.Status.BadRequest(), µ.NewOutput(http.StatusBadRequest))
	output(t, µ.Status.Unauthorized(), µ.NewOutput(http.StatusUnauthorized))
	output(t, µ.Status.PaymentRequired(), µ.NewOutput(http.StatusPaymentRequired))
	output(t, µ.Status.Forbidden(), µ.NewOutput(http.StatusForbidden))
	output(t, µ.Status.NotFound(), µ.NewOutput(http.StatusNotFound))
	output(t, µ.Status.MethodNotAllowed(), µ.NewOutput(http.StatusMethodNotAllowed))
	output(t, µ.Status.NotAcceptable(), µ.NewOutput(http.StatusNotAcceptable))
	output(t, µ.Status.ProxyAuthRequired(), µ.NewOutput(http.StatusProxyAuthRequired))
	output(t, µ.Status.RequestTimeout(), µ.NewOutput(http.StatusRequestTimeout))
	output(t, µ.Status.Conflict(), µ.NewOutput(http.StatusConflict))
	output(t, µ.Status.Gone(), µ.NewOutput(http.StatusGone))
	output(t, µ.Status.LengthRequired(), µ.NewOutput(http.StatusLengthRequired))
	output(t, µ.Status.PreconditionFailed(), µ.NewOutput(http.StatusPreconditionFailed))
	output(t, µ.Status.RequestEntityTooLarge(), µ.NewOutput(http.StatusRequestEntityTooLarge))
	output(t, µ.Status.RequestURITooLong(), µ.NewOutput(http.StatusRequestURITooLong))
	output(t, µ.Status.UnsupportedMediaType(), µ.NewOutput(http.StatusUnsupportedMediaType))

	//
	output(t, µ.Status.InternalServerError(), µ.NewOutput(http.StatusInternalServerError))
	output(t, µ.Status.NotImplemented(), µ.NewOutput(http.StatusNotImplemented))
	output(t, µ.Status.BadGateway(), µ.NewOutput(http.StatusBadGateway))
	output(t, µ.Status.ServiceUnavailable(), µ.NewOutput(http.StatusServiceUnavailable))
	output(t, µ.Status.GatewayTimeout(), µ.NewOutput(http.StatusGatewayTimeout))
	output(t, µ.Status.HTTPVersionNotSupported(), µ.NewOutput(http.StatusHTTPVersionNotSupported))
}

type myT struct {
	A string
}

func TestWithHeader(t *testing.T) {
	output(t,
		µ.Status.OK(µ.WithHeader("foo", "bar")),
		µ.Status.OK(µ.WithHeader("foo", "bar")),
	)
}

func TestWithJSON(t *testing.T) {
	output(t,
		µ.Status.OK(µ.WithJSON(myT{"Hello"})),
		µ.Status.OK(µ.WithJSON(myT{"Hello"})),
	)
}

func TestWithText(t *testing.T) {
	output(t,
		µ.Status.OK(µ.WithText("Hello")),
		µ.Status.OK(µ.WithText("Hello")),
	)
}

func TestWithBytes(t *testing.T) {
	output(t,
		µ.Status.OK(µ.WithBytes([]byte("Hello"))),
		µ.Status.OK(µ.WithBytes([]byte("Hello"))),
	)
}

func TestErrorOnJSON(t *testing.T) {
	output := µ.Status.OK(µ.WithJSON(make(chan int))).(*µ.Output)

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
