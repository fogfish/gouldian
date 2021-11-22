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

	"github.com/fogfish/guid"
)

/*

Output is HTTP response
*/
type Output struct {
	Status  int
	Headers map[string]string
	Body    string
	Failure error
}

// Output uses "error" interface
func (out Output) Error() string {
	return out.Body
}

// NewOutput creates HTTP response with given HTTP Status code
func NewOutput(status int) *Output {
	return &Output{
		Status:  status,
		Headers: map[string]string{},
	}
}

/*

Result is a composable function that abstract results of HTTP endpoint.
The function takes instance of HTTP output and mutates its value

  return µ.Status.OK(
		headers.ContentType.Value("application/json"),
		µ.WithJSON(value),
	)
*/
type Result func(*Output) error

// Issue implements RFC 7807: Problem Details for HTTP APIs
type Issue struct {
	ID     string `json:"instance"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Title  string `json:"title"`
}

// NewIssue creates instance of Issue
func NewIssue(status int) *Issue {
	return &Issue{
		ID:     guid.G.K(guid.Clock).String(),
		Type:   fmt.Sprintf("https://httpstatuses.com/%d", status),
		Status: status,
		Title:  http.StatusText(status),
	}
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

func (code StatusCode) output(status int, out []Result) *Output {
	v := NewOutput(status)
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
func (code StatusCode) OK(out ...Result) error {
	return code.output(http.StatusOK, out)
}

// Created ⟼ http.StatusCreated
func (code StatusCode) Created(out ...Result) error {
	return code.output(http.StatusCreated, out)
}

// Accepted ⟼ http.StatusAccepted
func (code StatusCode) Accepted(out ...Result) error {
	return code.output(http.StatusAccepted, out)
}

// NonAuthoritativeInfo ⟼ http.StatusNonAuthoritativeInfo
func (code StatusCode) NonAuthoritativeInfo(out ...Result) error {
	return code.output(http.StatusNonAuthoritativeInfo, out)
}

// NoContent ⟼ http.StatusNoContent
func (code StatusCode) NoContent(out ...Result) error {
	return code.output(http.StatusNoContent, out)
}

// ResetContent ⟼ http.StatusResetContent
func (code StatusCode) ResetContent(out ...Result) error {
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
func (code StatusCode) MultipleChoices(out ...Result) error {
	return code.output(http.StatusMultipleChoices, out)
}

// MovedPermanently ⟼ http.StatusMovedPermanently
func (code StatusCode) MovedPermanently(out ...Result) error {
	return code.output(http.StatusMovedPermanently, out)
}

// Found ⟼ http.StatusFound
func (code StatusCode) Found(out ...Result) error {
	return code.output(http.StatusFound, out)
}

// SeeOther ⟼ http.StatusSeeOther
func (code StatusCode) SeeOther(out ...Result) error {
	return code.output(http.StatusSeeOther, out)
}

// NotModified ⟼ http.StatusNotModified
func (code StatusCode) NotModified(out ...Result) error {
	return code.output(http.StatusNotModified, out)
}

// UseProxy ⟼ http.StatusUseProxy
func (code StatusCode) UseProxy(out ...Result) error {
	return code.output(http.StatusUseProxy, out)
}

// TemporaryRedirect ⟼ http.StatusTemporaryRedirect
func (code StatusCode) TemporaryRedirect(out ...Result) error {
	return code.output(http.StatusTemporaryRedirect, out)
}

// PermanentRedirect ⟼ http.StatusPermanentRedirect
func (code StatusCode) PermanentRedirect(out ...Result) error {
	return code.output(http.StatusPermanentRedirect, out)
}

//
//
//

// BadRequest ⟼ http.StatusBadRequest
func (code StatusCode) BadRequest(out ...Result) error {
	return code.output(http.StatusBadRequest, out)
}

// Unauthorized ⟼ http.StatusUnauthorized
func (code StatusCode) Unauthorized(out ...Result) error {
	return code.output(http.StatusUnauthorized, out)
}

// PaymentRequired ⟼ http.StatusPaymentRequired
func (code StatusCode) PaymentRequired(out ...Result) error {
	return code.output(http.StatusPaymentRequired, out)
}

// Forbidden ⟼ http.StatusForbidden
func (code StatusCode) Forbidden(out ...Result) error {
	return code.output(http.StatusForbidden, out)
}

// NotFound ⟼ http.StatusNotFound
func (code StatusCode) NotFound(out ...Result) error {
	return code.output(http.StatusNotFound, out)
}

// MethodNotAllowed ⟼ http.StatusMethodNotAllowed
func (code StatusCode) MethodNotAllowed(out ...Result) error {
	return code.output(http.StatusMethodNotAllowed, out)
}

// NotAcceptable ⟼ http.StatusNotAcceptable
func (code StatusCode) NotAcceptable(out ...Result) error {
	return code.output(http.StatusNotAcceptable, out)
}

// ProxyAuthRequired ⟼ http.StatusProxyAuthRequired
func (code StatusCode) ProxyAuthRequired(out ...Result) error {
	return code.output(http.StatusProxyAuthRequired, out)
}

// RequestTimeout ⟼ http.StatusRequestTimeout
func (code StatusCode) RequestTimeout(out ...Result) error {
	return code.output(http.StatusRequestTimeout, out)
}

// Conflict ⟼ http.StatusConflict
func (code StatusCode) Conflict(out ...Result) error {
	return code.output(http.StatusConflict, out)
}

// Gone ⟼ http.StatusGone
func (code StatusCode) Gone(out ...Result) error {
	return code.output(http.StatusGone, out)
}

// LengthRequired ⟼ http.StatusLengthRequired
func (code StatusCode) LengthRequired(out ...Result) error {
	return code.output(http.StatusLengthRequired, out)
}

// PreconditionFailed ⟼ http.StatusPreconditionFailed
func (code StatusCode) PreconditionFailed(out ...Result) error {
	return code.output(http.StatusPreconditionFailed, out)
}

// RequestEntityTooLarge ⟼ http.StatusRequestEntityTooLarge
func (code StatusCode) RequestEntityTooLarge(out ...Result) error {
	return code.output(http.StatusRequestEntityTooLarge, out)
}

// RequestURITooLong ⟼ http.StatusRequestURITooLong
func (code StatusCode) RequestURITooLong(out ...Result) error {
	return code.output(http.StatusRequestURITooLong, out)
}

// UnsupportedMediaType ⟼ http.StatusUnsupportedMediaType
func (code StatusCode) UnsupportedMediaType(out ...Result) error {
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
func (code StatusCode) InternalServerError(out ...Result) error {
	return code.output(http.StatusInternalServerError, out)
}

// NotImplemented ⟼ http.StatusNotImplemented
func (code StatusCode) NotImplemented(out ...Result) error {
	return code.output(http.StatusNotImplemented, out)
}

// BadGateway ⟼ http.StatusBadGateway
func (code StatusCode) BadGateway(out ...Result) error {
	return code.output(http.StatusBadGateway, out)
}

// ServiceUnavailable ⟼ http.StatusServiceUnavailable
func (code StatusCode) ServiceUnavailable(out ...Result) error {
	return code.output(http.StatusServiceUnavailable, out)
}

// GatewayTimeout ⟼ http.StatusGatewayTimeout
func (code StatusCode) GatewayTimeout(out ...Result) error {
	return code.output(http.StatusGatewayTimeout, out)
}

// HTTPVersionNotSupported ⟼ http.StatusHTTPVersionNotSupported
func (code StatusCode) HTTPVersionNotSupported(out ...Result) error {
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

// WithHeader appends header to HTTP response
func WithHeader(header, value string) Result {
	return func(out *Output) error {
		out.Headers[header] = value
		return nil
	}
}

// WithJSON appends application/json payload to HTTP response
func WithJSON(val interface{}) Result {
	return func(out *Output) error {
		body, err := json.Marshal(val)
		if err != nil {
			out.Status = http.StatusInternalServerError
			out.Headers["Content-Type"] = "text/plain"
			out.Body = fmt.Sprintf("JSON serialization is failed for <%T>", val)

			return nil
		}

		out.Headers["Content-Type"] = "application/json"
		out.Body = string(body)
		return nil
	}
}

// WithBytes appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func WithBytes(content []byte) Result {
	return func(out *Output) error {
		out.Body = string(content)
		return nil
	}
}

// WithText appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func WithText(content string) Result {
	return func(out *Output) error {
		out.Body = content
		return nil
	}
}

// WithIssue appends Issue, RFC 7807: Problem Details for HTTP APIs
func WithIssue(err error, title ...string) Result {
	return func(out *Output) error {
		issue := NewIssue(out.Status)
		if len(title) != 0 {
			issue.Title = title[0]
		}

		body, err := json.Marshal(issue)
		if err != nil {
			out.Status = http.StatusInternalServerError
			out.Headers["Content-Type"] = "text/plain"
			out.Body = fmt.Sprintf("JSON serialization is failed for <Issue>")

			return nil
		}

		out.Headers["Content-Type"] = "application/json"
		out.Body = string(body)
		out.Failure = err
		return nil
	}
}
