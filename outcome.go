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
	"errors"
	"fmt"
	"net/http"
)

// NoMatch is returned by Endpoint if Context is not matched: "not my route".
// The value is a reason code, it costs nothing to carry (no allocations).
// Coproduct (Or) tries the next endpoint on NoMatch. If no endpoint matches
// then the most specific reason is reported (the one with the largest code).
type NoMatch uint8

// Reasons of NoMatch, ordered from least to most specific.
const (
	NoMatchAny NoMatch = iota + 1
	NoMatchPath
	NoMatchMethod
	NoMatchHeader
	NoMatchParam
	NoMatchJWT
	NoMatchBody
	// content negotiation, a service responds with 406 or 415 if
	// no other endpoint matches the request
	NoMatchAccept
	NoMatchContentType
)

func (err NoMatch) Error() string {
	switch err {
	case NoMatchPath:
		return "No Match: path"
	case NoMatchMethod:
		return "No Match: method"
	case NoMatchHeader:
		return "No Match: header"
	case NoMatchParam:
		return "No Match: param"
	case NoMatchJWT:
		return "No Match: jwt"
	case NoMatchBody:
		return "No Match: body"
	case NoMatchAccept:
		return "No Match: accept"
	case NoMatchContentType:
		return "No Match: content-type"
	default:
		return "No Match"
	}
}

// ErrNoMatch is the generic NoMatch outcome
var ErrNoMatch error = NoMatchAny

// Source of the input that causes a failure
type Source uint8

const (
	SourceUnknown Source = iota
	SourcePath
	SourceParam
	SourceHeader
	SourceBody
	SourceJWT
)

func (s Source) String() string {
	switch s {
	case SourcePath:
		return "path"
	case SourceParam:
		return "query"
	case SourceHeader:
		return "header"
	case SourceBody:
		return "body"
	case SourceJWT:
		return "jwt"
	default:
		return "unknown"
	}
}

// Failure is returned by Endpoint if the request is routed to the endpoint
// but the input is invalid: a value is present but malformed (decode error),
// a body is empty or exceeds the limit. It is the "my route, bad input"
// outcome. Failure is terminal, coproduct (Or) does not try other endpoints.
// The server converts the failure to an Issue (RFC 9457).
type Failure struct {
	// HTTP Status, 400 by default
	Status int
	// Source of the input, the Key identifies the value within the source:
	// name of param, header, JSON pointer of body, etc.
	Source Source
	Key    string
	// Cause of the failure
	Err error
	// Extra headers of the response (e.g. Allow for 405)
	Header http.Header
}

func (f *Failure) Error() string {
	status := f.Status
	if status == 0 {
		status = http.StatusBadRequest
	}

	if f.Key == "" {
		return fmt.Sprintf("%d %s: %v", status, f.Source, f.Err)
	}
	return fmt.Sprintf("%d %s %q: %v", status, f.Source, f.Key, f.Err)
}

func (f *Failure) Unwrap() error { return f.Err }

// Detail is a human-readable explanation, it is safe to show to the client.
func (f *Failure) Detail() string {
	var text string
	if f.Err != nil {
		text = f.Err.Error()
	}

	switch {
	case f.Source == SourceUnknown || f.Source == SourceBody && f.Key == "":
		return text
	case f.Key == "":
		return fmt.Sprintf("%s: %s", f.Source, text)
	default:
		return fmt.Sprintf("%s %s: %s", f.Source, f.Key, text)
	}
}

// Pointer identifies the invalid value: RFC 6901 for the body, the name of
// param or header otherwise.
func (f *Failure) Pointer() string { return f.Key }

// Common failures
var (
	ErrBodyRequired = errors.New("body is required")
	ErrBodyTooLarge = errors.New("body is too large")
	ErrMethod       = errors.New("method is not allowed")
)

func newFailure(src Source, key string, err error) *Failure {
	return &Failure{Status: http.StatusBadRequest, Source: src, Key: key, Err: err}
}

// Lenient turns a Failure of the endpoint into NoMatch. It brings back the
// behaviour of v2: a malformed value is yet another "not my route".
func Lenient(e Endpoint) Endpoint {
	return func(ctx *Context) error {
		err := e(ctx)
		if _, ok := err.(*Failure); ok {
			return NoMatchAny
		}
		return err
	}
}

// Trace logs (debug level) the outcome of the endpoint, use it to understand
// why the request is not matched.
func Trace(name string, e Endpoint) Endpoint {
	return func(ctx *Context) error {
		err := e(ctx)
		if err != nil {
			if _, isOut := err.(*Output); !isOut {
				ctx.logger().DebugContext(ctx, "trace", "endpoint", name, "outcome", err.Error())
			}
		}
		return err
	}
}
