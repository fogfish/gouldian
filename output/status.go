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

package emitter

import (
	µ "github.com/fogfish/gouldian/v2"
	"net/http"
)

// StatusCode is a wrapper type over http.StatusCode so that ...
type StatusCode int

// Status is collection of constants for HTTP Status Code
//
//	return µ.Status.Ok()
const Status = StatusCode(0)

func (code StatusCode) output(status int, out []µ.Result) *µ.Output {
	v := µ.NewOutput(status)
	for _, f := range out {
		f(v)
	}
	return v
}

/*
TODO:
  Continue
	SwitchingProtocols
	Processing
	EarlyHints
*/

// OK ⟼ http.StatusOK
func (code StatusCode) OK(out ...µ.Result) error {
	return code.output(http.StatusOK, out)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created(out ...µ.Result) error {
	return code.output(http.StatusCreated, out)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted(out ...µ.Result) error {
	return code.output(http.StatusAccepted, out)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo(out ...µ.Result) error {
	return code.output(http.StatusNonAuthoritativeInfo, out)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent(out ...µ.Result) error {
	return code.output(http.StatusNoContent, out)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent(out ...µ.Result) error {
	return code.output(http.StatusResetContent, out)
}

/*
TODO:
	PartialContent
	MultiStatus
	AlreadyReported
	IMUsed
*/

// MultipleChoices ⟼ http.StatusMultipleChoices
func (code StatusCode) MultipleChoices(out ...µ.Result) error {
	return code.output(http.StatusMultipleChoices, out)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(out ...µ.Result) error {
	return code.output(http.StatusMovedPermanently, out)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(out ...µ.Result) error {
	return code.output(http.StatusFound, out)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(out ...µ.Result) error {
	return code.output(http.StatusSeeOther, out)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(out ...µ.Result) error {
	return code.output(http.StatusNotModified, out)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(out ...µ.Result) error {
	return code.output(http.StatusUseProxy, out)
}

// TemporaryRedirect ⟼ http.StatusTemporaryRedirect
func (code StatusCode) TemporaryRedirect(out ...µ.Result) error {
	return code.output(http.StatusTemporaryRedirect, out)
}

// PermanentRedirect ⟼ http.StatusPermanentRedirect
func (code StatusCode) PermanentRedirect(out ...µ.Result) error {
	return code.output(http.StatusPermanentRedirect, out)
}

//
//
//

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(out ...µ.Result) error {
	return code.output(http.StatusBadRequest, out)
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(out ...µ.Result) error {
	return code.output(http.StatusUnauthorized, out)
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(out ...µ.Result) error {
	return code.output(http.StatusPaymentRequired, out)
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(out ...µ.Result) error {
	return code.output(http.StatusForbidden, out)
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(out ...µ.Result) error {
	return code.output(http.StatusNotFound, out)
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(out ...µ.Result) error {
	return code.output(http.StatusMethodNotAllowed, out)
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(out ...µ.Result) error {
	return code.output(http.StatusNotAcceptable, out)
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(out ...µ.Result) error {
	return code.output(http.StatusProxyAuthRequired, out)
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(out ...µ.Result) error {
	return code.output(http.StatusRequestTimeout, out)
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(out ...µ.Result) error {
	return code.output(http.StatusConflict, out)
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(out ...µ.Result) error {
	return code.output(http.StatusGone, out)
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(out ...µ.Result) error {
	return code.output(http.StatusLengthRequired, out)
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(out ...µ.Result) error {
	return code.output(http.StatusPreconditionFailed, out)
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(out ...µ.Result) error {
	return code.output(http.StatusRequestEntityTooLarge, out)
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(out ...µ.Result) error {
	return code.output(http.StatusRequestURITooLong, out)
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(out ...µ.Result) error {
	return code.output(http.StatusUnsupportedMediaType, out)
}

/*
TODO:
	RequestedRangeNotSatisfiable
	ExpectationFailed
	Teapot
	MisdirectedRequest
	UnprocessableEntity
	Locked
	FailedDependency
	TooEarly
	UpgradeRequired
	PreconditionRequired
	TooManyRequests
	RequestHeaderFieldsTooLarge
	UnavailableForLegalReasons
*/

// InternalServerError ⟼ http.StatusInternalServerError
func (code StatusCode) InternalServerError(out ...µ.Result) error {
	return code.output(http.StatusInternalServerError, out)
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(out ...µ.Result) error {
	return code.output(http.StatusNotImplemented, out)
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(out ...µ.Result) error {
	return code.output(http.StatusBadGateway, out)
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(out ...µ.Result) error {
	return code.output(http.StatusServiceUnavailable, out)
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(out ...µ.Result) error {
	return code.output(http.StatusGatewayTimeout, out)
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(out ...µ.Result) error {
	return code.output(http.StatusHTTPVersionNotSupported, out)
}

/*
TODO:
	VariantAlsoNegotiates
	InsufficientStorage
	LoopDetected
	NotExtended
	NetworkAuthenticationRequired
*/
