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

package output

import (
	µ "github.com/fogfish/gouldian/v3"
	"net/http"
)

// StatusCode is a wrapper type over http.StatusCode so that ...
type StatusCode int

// Status is collection of constants for HTTP Status Code
//
//	return µ.Status.Ok()
const Status = StatusCode(0)

func (code StatusCode) build(status int, out []µ.Result) (*µ.Output, error) {
	v := µ.NewOutput(status)
	for _, f := range out {
		if err := f(v); err != nil {
			// the failed result is reported to the caller, server logs it as 500
			v.Free()
			return nil, err
		}
	}
	return v, nil
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
	return code.make(http.StatusOK, out)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created(out ...µ.Result) error {
	return code.make(http.StatusCreated, out)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted(out ...µ.Result) error {
	return code.make(http.StatusAccepted, out)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo(out ...µ.Result) error {
	return code.make(http.StatusNonAuthoritativeInfo, out)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent(out ...µ.Result) error {
	return code.make(http.StatusNoContent, out)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent(out ...µ.Result) error {
	return code.make(http.StatusResetContent, out)
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
	return code.make(http.StatusMultipleChoices, out)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(out ...µ.Result) error {
	return code.make(http.StatusMovedPermanently, out)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(out ...µ.Result) error {
	return code.make(http.StatusFound, out)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(out ...µ.Result) error {
	return code.make(http.StatusSeeOther, out)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(out ...µ.Result) error {
	return code.make(http.StatusNotModified, out)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(out ...µ.Result) error {
	return code.make(http.StatusUseProxy, out)
}

// TemporaryRedirect ⟼ http.StatusTemporaryRedirect
func (code StatusCode) TemporaryRedirect(out ...µ.Result) error {
	return code.make(http.StatusTemporaryRedirect, out)
}

// PermanentRedirect ⟼ http.StatusPermanentRedirect
func (code StatusCode) PermanentRedirect(out ...µ.Result) error {
	return code.make(http.StatusPermanentRedirect, out)
}

//
//
//

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(out ...µ.Result) error {
	return code.make(http.StatusBadRequest, out)
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(out ...µ.Result) error {
	return code.make(http.StatusUnauthorized, out)
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(out ...µ.Result) error {
	return code.make(http.StatusPaymentRequired, out)
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(out ...µ.Result) error {
	return code.make(http.StatusForbidden, out)
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(out ...µ.Result) error {
	return code.make(http.StatusNotFound, out)
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(out ...µ.Result) error {
	return code.make(http.StatusMethodNotAllowed, out)
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(out ...µ.Result) error {
	return code.make(http.StatusNotAcceptable, out)
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(out ...µ.Result) error {
	return code.make(http.StatusProxyAuthRequired, out)
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(out ...µ.Result) error {
	return code.make(http.StatusRequestTimeout, out)
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(out ...µ.Result) error {
	return code.make(http.StatusConflict, out)
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(out ...µ.Result) error {
	return code.make(http.StatusGone, out)
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(out ...µ.Result) error {
	return code.make(http.StatusLengthRequired, out)
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(out ...µ.Result) error {
	return code.make(http.StatusPreconditionFailed, out)
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(out ...µ.Result) error {
	return code.make(http.StatusRequestEntityTooLarge, out)
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(out ...µ.Result) error {
	return code.make(http.StatusRequestURITooLong, out)
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(out ...µ.Result) error {
	return code.make(http.StatusUnsupportedMediaType, out)
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
	return code.make(http.StatusInternalServerError, out)
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(out ...µ.Result) error {
	return code.make(http.StatusNotImplemented, out)
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(out ...µ.Result) error {
	return code.make(http.StatusBadGateway, out)
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(out ...µ.Result) error {
	return code.make(http.StatusServiceUnavailable, out)
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(out ...µ.Result) error {
	return code.make(http.StatusGatewayTimeout, out)
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(out ...µ.Result) error {
	return code.make(http.StatusHTTPVersionNotSupported, out)
}

/*
TODO:
	VariantAlsoNegotiates
	InsufficientStorage
	LoopDetected
	NotExtended
	NetworkAuthenticationRequired
*/

// helper converts the outcome of build to error, a nil *Output must not be
// wrapped to non-nil error interface.
func (code StatusCode) make(status int, out []µ.Result) error {
	v, err := code.build(status, out)
	if err != nil {
		return err
	}
	return v
}

// PartialContent ⟼ http.StatusPartialContent
func (code StatusCode) PartialContent(out ...µ.Result) error {
	return code.make(http.StatusPartialContent, out)
}

// MultiStatus ⟼ http.StatusMultiStatus
func (code StatusCode) MultiStatus(out ...µ.Result) error {
	return code.make(http.StatusMultiStatus, out)
}

// RequestedRangeNotSatisfiable ⟼ http.StatusRequestedRangeNotSatisfiable
func (code StatusCode) RequestedRangeNotSatisfiable(out ...µ.Result) error {
	return code.make(http.StatusRequestedRangeNotSatisfiable, out)
}

// ExpectationFailed ⟼ http.StatusExpectationFailed
func (code StatusCode) ExpectationFailed(out ...µ.Result) error {
	return code.make(http.StatusExpectationFailed, out)
}

// MisdirectedRequest ⟼ http.StatusMisdirectedRequest
func (code StatusCode) MisdirectedRequest(out ...µ.Result) error {
	return code.make(http.StatusMisdirectedRequest, out)
}

// UnprocessableEntity ⟼ http.StatusUnprocessableEntity
func (code StatusCode) UnprocessableEntity(out ...µ.Result) error {
	return code.make(http.StatusUnprocessableEntity, out)
}

// Locked ⟼ http.StatusLocked
func (code StatusCode) Locked(out ...µ.Result) error {
	return code.make(http.StatusLocked, out)
}

// FailedDependency ⟼ http.StatusFailedDependency
func (code StatusCode) FailedDependency(out ...µ.Result) error {
	return code.make(http.StatusFailedDependency, out)
}

// UpgradeRequired ⟼ http.StatusUpgradeRequired
func (code StatusCode) UpgradeRequired(out ...µ.Result) error {
	return code.make(http.StatusUpgradeRequired, out)
}

// PreconditionRequired ⟼ http.StatusPreconditionRequired
func (code StatusCode) PreconditionRequired(out ...µ.Result) error {
	return code.make(http.StatusPreconditionRequired, out)
}

// TooManyRequests ⟼ http.StatusTooManyRequests
func (code StatusCode) TooManyRequests(out ...µ.Result) error {
	return code.make(http.StatusTooManyRequests, out)
}

// RequestHeaderFieldsTooLarge ⟼ http.StatusRequestHeaderFieldsTooLarge
func (code StatusCode) RequestHeaderFieldsTooLarge(out ...µ.Result) error {
	return code.make(http.StatusRequestHeaderFieldsTooLarge, out)
}

// UnavailableForLegalReasons ⟼ http.StatusUnavailableForLegalReasons
func (code StatusCode) UnavailableForLegalReasons(out ...µ.Result) error {
	return code.make(http.StatusUnavailableForLegalReasons, out)
}

// VariantAlsoNegotiates ⟼ http.StatusVariantAlsoNegotiates
func (code StatusCode) VariantAlsoNegotiates(out ...µ.Result) error {
	return code.make(http.StatusVariantAlsoNegotiates, out)
}

// InsufficientStorage ⟼ http.StatusInsufficientStorage
func (code StatusCode) InsufficientStorage(out ...µ.Result) error {
	return code.make(http.StatusInsufficientStorage, out)
}

// NetworkAuthenticationRequired ⟼ http.StatusNetworkAuthenticationRequired
func (code StatusCode) NetworkAuthenticationRequired(out ...µ.Result) error {
	return code.make(http.StatusNetworkAuthenticationRequired, out)
}
