//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Output defines legitimate HTTP response. It allows to specify
// HTTP Headers and Body. The structure allows to use any HTTP
// status code.
//   gouldian.Ok().With("X-Header", "value").Json(MyStruct{})
type Output struct {
	Status  int
	Headers map[string]string
	Body    string
}

func (out Output) Error() string {
	return out.Body
}

// Success creates HTTP response with given HTTP Status code
func Success(status int) *Output {
	return &Output{status, map[string]string{}, ""}
}

// Ok is an alias of "200 Ok" output
func Ok() *Output { return Success(http.StatusOK) }

// Created is an alias of "201 Created" output
func Created() *Output { return Success(http.StatusCreated) }

// Accepted is an alias of "202 Accepted" output
func Accepted() *Output { return Success(http.StatusAccepted) }

// NoContent is an alias of "204 No Content" output
func NoContent() *Output { return Success(http.StatusNoContent) }

// MovedPermanently is an alias of "301 Moved Permanently" output
func MovedPermanently(uri url.URL) *Output {
	return Success(http.StatusMovedPermanently).With("Location", uri.String())
}

// Found is an alias of "302 Found" output
func Found(uri url.URL) *Output {
	return Success(http.StatusFound).With("Location", uri.String())
}

// SeeOther is an alias of "303 See Other" output
func SeeOther(uri url.URL) *Output {
	return Success(http.StatusSeeOther).With("Location", uri.String())
}

// NotModified is an alias of "304 Not Modified" output
func NotModified(uri url.URL) *Output {
	return Success(http.StatusNotModified).With("Location", uri.String())
}

// TemporaryRedirect is an alias of "307 Temporary Redirect" output
func TemporaryRedirect(uri url.URL) *Output {
	return Success(http.StatusTemporaryRedirect).With("Location", uri.String())
}

// PermanentRedirect is an alias of "308 Permanent Redirect" output
func PermanentRedirect(uri url.URL) *Output {
	return Success(http.StatusPermanentRedirect).With("Location", uri.String())
}

// JSON appends application/json payload to HTTP response
func (out *Output) JSON(val interface{}) *Output {
	body, _ := json.Marshal(val)
	out.Headers["Content-Type"] = "application/json"
	out.Body = string(body)
	return out
}

// Text appends text/plain payload to HTTP response
func (out *Output) Text(text string) *Output {
	out.Body = text
	out.Headers["Content-Type"] = "text/plain"
	return out
}

// With sets HTTP header to the response
func (out *Output) With(head string, value string) *Output {
	out.Headers[head] = value
	return out
}

// Issue implements RFC 7807: Problem Details for HTTP APIs
type Issue struct {
	Type    string      `json:"type"`
	Status  int         `json:"status"`
	Title   string      `json:"title"`
	Details interface{} `json:"details"`
}

func (err Issue) Error() string {
	return fmt.Sprintf(strconv.Itoa(err.Status) + ": " + err.Title)
}

// Reason defines details of the issue
func (err *Issue) Reason(reason interface{}) *Issue {
	err.Details = reason
	return err
}

// Failure creates HTTP issue with given HTTP Status code
func Failure(status int) *Issue {
	return &Issue{typeOf(status), status, http.StatusText(status), ""}
}

// BadRequest is an alias of "400 Bad Request" issue
func BadRequest(reason interface{}) *Issue {
	return Failure(http.StatusBadRequest).Reason(reason)
}

// Unauthorized is an alias of "401 Unauthorized" issue
func Unauthorized(reason interface{}) *Issue {
	return Failure(http.StatusUnauthorized).Reason(reason)
}

// Forbidden is an alias of "403 Forbidden" issue
func Forbidden(reason interface{}) *Issue {
	return Failure(http.StatusForbidden).Reason(reason)
}

// NotFound is an alias of "404 Not Found" issue
func NotFound(reason interface{}) *Issue {
	return Failure(http.StatusNotFound).Reason(reason)
}

// InternalServerError is an alias of "500 Internal Server Error" issue
func InternalServerError(reason interface{}) *Issue {
	return Failure(http.StatusInternalServerError).Reason(reason)
}

// NotImplemented is an alias of "501 Not Implemented" issue
func NotImplemented(reason interface{}) *Issue {
	return Failure(http.StatusNotImplemented).Reason(reason)
}

// ServiceUnavailable is an alias of "503 Service Unavailable" issue
func ServiceUnavailable(reason interface{}) *Issue {
	return Failure(http.StatusServiceUnavailable).Reason(reason)
}

func typeOf(status int) string {
	return fmt.Sprintf("https://httpstatuses.com/%v", status)
}
