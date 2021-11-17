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

package gouldian

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fogfish/guid"
)

/*

Output HTTP response
*/
type Output interface {
	error

	With(Header, string) Output
	JSON(interface{}) Output
	Bytes([]byte) Output
	Text(string) Output
}

/*

StatusCode is a warpper type over http.StatusCode so that ...
*/
type StatusCode int

/*

Status is collection of constants for HTTP Status Code

  return µ.Status.Ok()
*/
const Status = StatusCode(0)

/*
TODO:
  Continue
	SwitchingProtocols
	Processing
	EarlyHints
*/

// OK ⟼ http.StatusOK
func (code StatusCode) OK() Output {
	return NewSuccess(http.StatusOK)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created() Output {
	return NewSuccess(http.StatusCreated)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted() Output {
	return NewSuccess(http.StatusAccepted)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo() Output {
	return NewSuccess(http.StatusNonAuthoritativeInfo)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent() Output {
	return NewSuccess(http.StatusNoContent)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent() Output {
	return NewSuccess(http.StatusResetContent)
}

/*
TODO:
	PartialContent
	MultiStatus
	AlreadyReported
	IMUsed
*/

// MultipleChoices ⟼ http.StatusMultipleChoices
func (code StatusCode) MultipleChoices() Output {
	return NewSuccess(http.StatusMultipleChoices)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(url string) Output {
	return NewSuccess(http.StatusMovedPermanently).With("Location", url)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(url string) Output {
	return NewSuccess(http.StatusFound).With("Location", url)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(url string) Output {
	return NewSuccess(http.StatusSeeOther).With("Location", url)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(url string) Output {
	return NewSuccess(http.StatusNotModified).With("Location", url)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(url string) Output {
	return NewSuccess(http.StatusUseProxy).With("Location", url)
}

// TemporaryRedirect ⟼ http.StatusTemporaryRedirect
func (code StatusCode) TemporaryRedirect(url string) Output {
	return NewSuccess(http.StatusTemporaryRedirect).With("Location", url)
}

// PermanentRedirect ⟼ http.StatusPermanentRedirect
func (code StatusCode) PermanentRedirect(url string) Output {
	return NewSuccess(http.StatusPermanentRedirect).With("Location", url)
}

//
//
//

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(err error, title ...string) Output {
	return NewFailure(http.StatusBadRequest, err).Bytes([]byte(strings.Join(title, "")))
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(err error, title ...string) Output {
	return NewFailure(http.StatusUnauthorized, err).Bytes([]byte(strings.Join(title, "")))
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(err error, title ...string) Output {
	return NewFailure(http.StatusPaymentRequired, err).Bytes([]byte(strings.Join(title, "")))
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(err error, title ...string) Output {
	return NewFailure(http.StatusForbidden, err).Bytes([]byte(strings.Join(title, "")))
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(err error, title ...string) Output {
	return NewFailure(http.StatusNotFound, err).Bytes([]byte(strings.Join(title, "")))
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(err error, title ...string) Output {
	return NewFailure(http.StatusMethodNotAllowed, err).Bytes([]byte(strings.Join(title, "")))
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(err error, title ...string) Output {
	return NewFailure(http.StatusNotAcceptable, err).Bytes([]byte(strings.Join(title, "")))
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(err error, title ...string) Output {
	return NewFailure(http.StatusProxyAuthRequired, err).Bytes([]byte(strings.Join(title, "")))
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(err error, title ...string) Output {
	return NewFailure(http.StatusRequestTimeout, err).Bytes([]byte(strings.Join(title, "")))
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(err error, title ...string) Output {
	return NewFailure(http.StatusConflict, err).Bytes([]byte(strings.Join(title, "")))
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(err error, title ...string) Output {
	return NewFailure(http.StatusGone, err).Bytes([]byte(strings.Join(title, "")))
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(err error, title ...string) Output {
	return NewFailure(http.StatusLengthRequired, err).Bytes([]byte(strings.Join(title, "")))
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(err error, title ...string) Output {
	return NewFailure(http.StatusPreconditionFailed, err).Bytes([]byte(strings.Join(title, "")))
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(err error, title ...string) Output {
	return NewFailure(http.StatusRequestEntityTooLarge, err).Bytes([]byte(strings.Join(title, "")))
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(err error, title ...string) Output {
	return NewFailure(http.StatusRequestURITooLong, err).Bytes([]byte(strings.Join(title, "")))
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(err error, title ...string) Output {
	return NewFailure(http.StatusUnsupportedMediaType, err).Bytes([]byte(strings.Join(title, "")))
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
func (code StatusCode) InternalServerError(err error, title ...string) Output {
	return NewFailure(http.StatusInternalServerError, err).Bytes([]byte(strings.Join(title, "")))
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(err error, title ...string) Output {
	return NewFailure(http.StatusNotImplemented, err).Bytes([]byte(strings.Join(title, "")))
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(err error, title ...string) Output {
	return NewFailure(http.StatusBadGateway, err).Bytes([]byte(strings.Join(title, "")))
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(err error, title ...string) Output {
	return NewFailure(http.StatusServiceUnavailable, err).Bytes([]byte(strings.Join(title, "")))
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(err error, title ...string) Output {
	return NewFailure(http.StatusGatewayTimeout, err).Bytes([]byte(strings.Join(title, "")))
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(err error, title ...string) Output {
	return NewFailure(http.StatusHTTPVersionNotSupported, err).Bytes([]byte(strings.Join(title, "")))
}

/*
TODO:
	VariantAlsoNegotiates
	InsufficientStorage
	LoopDetected
	NotExtended
	NetworkAuthenticationRequired
*/

//
//
//

// NewSuccess creates HTTP response with given HTTP Status code
func NewSuccess(status StatusCode) Output {
	return &Success{
		Status:  status,
		Headers: map[Header]string{},
		Body:    "",
	}
}

/*

Success defines legitimate HTTP response. It allows to specify
HTTP Headers and Body.

  µ.Ok().With("X-Header", "value").Json(MyStruct{})
*/
type Success struct {
	Status  StatusCode
	Headers map[Header]string
	Body    string
}

//
func (out Success) Error() string {
	return out.Body
}

// JSON appends application/json payload to HTTP response
func (out *Success) JSON(val interface{}) Output {
	body, err := json.Marshal(val)
	if err != nil {
		out.Status = http.StatusInternalServerError
		out.Headers["Content-Type"] = "text/plain"
		out.Body = fmt.Sprintf("JSON serialization is failed for <%T>", val)

		return out
	}

	out.Headers["Content-Type"] = "application/json"
	out.Body = string(body)
	return out
}

// Bytes appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func (out *Success) Bytes(content []byte) Output {
	out.Body = string(content)
	return out
}

// Text appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func (out *Success) Text(content string) Output {
	out.Body = content
	return out
}

// With sets HTTP header to the response
func (out *Success) With(header Header, value string) Output {
	out.Headers[header] = value
	return out
}

//
//
//

// NewFailure creates HTTP issue with given HTTP Status code
func NewFailure(status StatusCode, err error) Output {
	return &Failure{
		ID:      guid.Seq.ID(),
		Type:    fmt.Sprintf("https://httpstatuses.com/%d", status),
		Status:  status,
		Title:   http.StatusText(int(status)),
		Failure: err,
	}
}

// Failure implements RFC 7807: Problem Details for HTTP APIs
type Failure struct {
	ID      string     `json:"instance"`
	Type    string     `json:"type"`
	Status  StatusCode `json:"status"`
	Title   string     `json:"title"`
	Failure error      `json:"-"`
}

func (issue Failure) Error() string {
	return fmt.Sprintf("%d: %s", issue.Status, issue.Title)
}

// JSON appends application/json payload to HTTP response
func (issue *Failure) JSON(val interface{}) Output {
	// Do Nothing
	return issue
}

// Bytes appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func (issue *Failure) Bytes(content []byte) Output {
	if len(content) > 0 {
		issue.Title = string(content)
	}
	return issue
}

// Text appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func (issue *Failure) Text(content string) Output {
	if len(content) > 0 {
		issue.Title = content
	}
	return issue
}

// With sets HTTP header to the response
func (issue *Failure) With(head Header, value string) Output {
	// Do Nothing
	return issue
}
