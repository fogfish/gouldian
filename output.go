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
	"net/textproto"
	"sync"

	"github.com/fogfish/guid"
)

// Global pools
var (
	outputs sync.Pool
)

func init() {
	outputs.New = func() interface{} {
		return &Output{
			Headers: make([]struct {
				Header string
				Value  string
			}, 0, 20),
		}
	}
}

// Output
type Output struct {
	Status  int
	Headers []struct{ Header, Value string }
	Body    string
	Failure error
}

// Output uses "error" interface
func (out Output) Error() string {
	return out.Body
}

// NewOutput creates HTTP response with given HTTP Status code
func NewOutput(status int) *Output {
	out := outputs.Get().(*Output)
	out.Status = status
	return out
}

// Free releases output
func (out *Output) Free() {
	out.Failure = nil
	out.Body = ""
	out.Headers = out.Headers[:0]
	outputs.Put(out)
}

func (out *Output) SetHeader(header, value string) {
	out.Headers = append(out.Headers,
		struct {
			Header string
			Value  string
		}{textproto.CanonicalMIMEHeaderKey(header), value},
	)
}

func (out *Output) GetHeader(header string) string {
	h := textproto.CanonicalMIMEHeaderKey(header)
	for i := 0; i < len(out.Headers); i++ {
		if out.Headers[i].Header == h {
			return out.Headers[i].Value
		}
	}

	return ""
}

// WithIssue appends Issue, RFC 7807: Problem Details for HTTP APIs
func (out *Output) SetIssue(failure error, title ...string) {
	issue := NewIssue(out.Status)
	if len(title) != 0 {
		issue.Title = title[0]
	}

	body, err := json.Marshal(issue)
	if err != nil {
		out.Status = http.StatusInternalServerError
		out.Headers = append(out.Headers,
			struct {
				Header string
				Value  string
			}{"Content-Type", "text/plain"},
		)
		out.Body = "JSON serialization is failed for <Issue>"

		return
	}

	out.Headers = append(out.Headers,
		struct {
			Header string
			Value  string
		}{"Content-Type", "application/json"},
	)
	out.Body = string(body)
	out.Failure = fmt.Errorf("%s: %d %s - %w", issue.ID, out.Status, issue.Title, failure)
}

// Result is a composable function that abstract results of HTTP endpoint.
// The function takes instance of HTTP output and mutates its value
//
//	  return µ.Status.OK(
//			µ.WithHeader(headers.ContentType, headers.ApplicationJson),
//			µ.WithJSON(value),
//		)
type Result func(*Output) error

// Issue implements RFC 7807: Problem Details for HTTP APIs
type Issue struct {
	ID     string `json:"instance"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Title  string `json:"title"`
}

// NewIssue creates instance of Issue
func NewIssue(status int) Issue {
	return Issue{
		ID:     guid.G.K(guid.Clock).String(),
		Type:   fmt.Sprintf("https://httpstatuses.com/%d", status),
		Status: status,
		Title:  http.StatusText(status),
	}
}
