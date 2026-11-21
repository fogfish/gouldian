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
	"errors"
	"net/http"
	"net/textproto"
	"sync"
)

// Global pools
var (
	outputs sync.Pool
)

func init() {
	outputs.New = func() interface{} {
		return &Output{
			pooled: true,
			Headers: make([]struct {
				Header string
				Value  string
			}, 0, 20),
		}
	}
}

// Output is the HTTP response produced by the handler. Outputs created by
// NewOutput are pooled and released by the server, so never keep or share
// them: construct a literal (&µ.Output{...}) if the value is a global constant.
type Output struct {
	Status  int
	Headers []struct{ Header, Value string }
	Body    []byte
	// Failure is the cause of the failed request. The server logs it and,
	// if the Body is empty, renders it as Issue.
	Failure error
	// Title of the Issue
	Title string

	pooled    bool
	render    bool     // server renders Issue
	reason    error    // NoMatch reason, for logs
	stack     string   // stack of recovered panic
	preflight []string // methods of the resource (auto OPTIONS)
}

// Output uses "error" interface
func (out Output) Error() string {
	if out.Failure != nil {
		return out.Failure.Error()
	}
	return string(out.Body)
}

// Unwrap exposes the cause, so errors.Is and errors.As work
func (out *Output) Unwrap() error { return out.Failure }

// NewOutput creates HTTP response with given HTTP Status code
func NewOutput(status int) *Output {
	out := outputs.Get().(*Output)
	out.Status = status
	return out
}

// Free releases output
func (out *Output) Free() {
	if !out.pooled {
		return
	}
	out.Failure = nil
	out.Title = ""
	out.render = false
	out.reason = nil
	out.stack = ""
	out.preflight = nil
	out.Body = nil
	out.Headers = out.Headers[:0]
	outputs.Put(out)
}

// SetHeader sets the header, it replaces the existing value(s).
func (out *Output) SetHeader(header, value string) {
	h := textproto.CanonicalMIMEHeaderKey(header)
	for i := range out.Headers {
		if out.Headers[i].Header == h {
			out.Headers[i].Value = value
			// drop duplicates
			rest := out.Headers[i+1:]
			n := i + 1
			for _, x := range rest {
				if x.Header != h {
					out.Headers[n] = x
					n++
				}
			}
			out.Headers = out.Headers[:n]
			return
		}
	}
	out.AddHeader(h, value)
}

// AddHeader appends the header, existing values are kept (e.g. Set-Cookie).
func (out *Output) AddHeader(header, value string) {
	out.Headers = append(out.Headers,
		struct {
			Header string
			Value  string
		}{textproto.CanonicalMIMEHeaderKey(header), value},
	)
}

// GetHeader returns the first value of the header
func (out *Output) GetHeader(header string) string {
	h := textproto.CanonicalMIMEHeaderKey(header)
	for i := 0; i < len(out.Headers); i++ {
		if out.Headers[i].Header == h {
			return out.Headers[i].Value
		}
	}

	return ""
}

// SetIssue marks the output as failed. It attaches the cause, the server
// logs it and renders the response body as RFC 9457 Problem Details
// (see WithIssue to customise it). The cause of 4xx is shown to the client
// as detail, the cause of 5xx is never exposed.
func (out *Output) SetIssue(failure error, title ...string) {
	out.Failure = failure
	out.render = true
	if len(title) != 0 {
		out.Title = title[0]
	}
}

// Result is a composable function that abstract results of HTTP endpoint.
// The function takes instance of HTTP output and mutates its value
//
//	  return µ.Status.OK(
//			µ.WithHeader(headers.ContentType, headers.ApplicationJson),
//			µ.WithJSON(value),
//		)
type Result func(*Output) error

// Issue implements RFC 9457: Problem Details for HTTP APIs
type Issue struct {
	// Type is URI reference that identifies the problem, about:blank by default
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
	// Instance is the request path
	Instance string       `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// FieldError points to the invalid value
type FieldError struct {
	// Source of the value: path, query, header, body, jwt
	Source string `json:"source,omitempty"`
	// Pointer to the value: RFC 6901 JSON pointer for the body, the name of
	// query param or header otherwise.
	Pointer string `json:"pointer,omitempty"`
	Detail  string `json:"detail"`
}

// NewIssue creates instance of Issue
func NewIssue(status int) Issue {
	return Issue{
		Type:   "about:blank",
		Status: status,
		Title:  http.StatusText(status),
	}
}

// IssueRenderer converts the failure to the body of response
type IssueRenderer func(issue Issue, cause error) (contentType string, body []byte)

// DefaultIssue renders RFC 9457 `application/problem+json`
func DefaultIssue(issue Issue, cause error) (string, []byte) {
	body, err := json.Marshal(issue)
	if err != nil {
		return "text/plain", []byte(http.StatusText(issue.Status))
	}
	return "application/problem+json", body
}

// issue describes the failed output
func (out *Output) issue(instance string) Issue {
	issue := NewIssue(out.Status)
	issue.Instance = instance
	if out.Title != "" {
		issue.Title = out.Title
	}

	if out.Failure == nil || out.Status >= http.StatusInternalServerError {
		return issue
	}

	var failure *Failure
	if errors.As(out.Failure, &failure) {
		issue.Detail = failure.Detail()
		if failure.Key != "" && failure.Err != nil {
			issue.Errors = []FieldError{{
				Source:  failure.Source.String(),
				Pointer: failure.Key,
				Detail:  failure.Err.Error(),
			}}
		}
		return issue
	}

	issue.Detail = out.Failure.Error()
	return issue
}
