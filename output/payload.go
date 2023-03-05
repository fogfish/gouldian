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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	µ "github.com/fogfish/gouldian/v2"
)

func Send(data any) µ.Result {
	return func(out *µ.Output) error {
		chunked := out.GetHeader(string(TransferEncoding)) == "chunked"
		content := out.GetHeader(string(ContentType))

		switch stream := data.(type) {
		case string:
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(len(stream)))
			}
			out.Body = stream
			return nil
		case *strings.Reader:
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(stream.Len()))
			}
			val, err := io.ReadAll(stream)
			if err != nil {
				return err
			}
			out.Body = string(val)
		case []byte:
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(len(stream)))
			}
			out.Body = string(stream)
		case *bytes.Buffer:
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(stream.Len()))
			}
			val, err := io.ReadAll(stream)
			if err != nil {
				return err
			}
			out.Body = string(val)
		case *bytes.Reader:
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(stream.Len()))
			}
			val, err := io.ReadAll(stream)
			if err != nil {
				return err
			}
			out.Body = string(val)
		case io.Reader:
			val, err := io.ReadAll(stream)
			if err != nil {
				return err
			}
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(len(val)))
			}
			out.Body = string(val)
		default:
			val, err := encode(content, data)
			if err != nil {
				out.Status = http.StatusInternalServerError
				out.SetIssue(fmt.Errorf("serialization is failed for <%T>", val))
				return nil
			}
			if !chunked {
				out.SetHeader(string(ContentLength), strconv.Itoa(len(val)))
			}
			out.Body = string(val)
		}
		return nil
	}
}

func encode(content string, data any) (buf []byte, err error) {
	switch {
	// "application/json" and other variants
	case strings.Contains(content, "json"):
		buf, err = encodeJSON(data)
	// "application/x-www-form-urlencoded"
	case strings.Contains(content, "www-form"):
		buf, err = encodeForm(data)
	default:
		buf, err = encodeJSON(data)
	}

	return
}

func encodeJSON(data interface{}) ([]byte, error) {
	json, err := json.Marshal(data)
	return json, err
}

func encodeForm(data interface{}) ([]byte, error) {
	bin, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var req map[string]string
	err = json.Unmarshal(bin, &req)
	if err != nil {
		return nil, fmt.Errorf("encode application/x-www-form-urlencoded: %w", err)
	}

	var payload url.Values = make(map[string][]string)
	for key, val := range req {
		payload[key] = []string{val}
	}
	return []byte(payload.Encode()), nil
}

// Error appends Issue, RFC 7807: Problem Details for HTTP APIs
func Error(failure error, title ...string) µ.Result {
	return func(out *µ.Output) error {
		out.SetIssue(failure, title...)
		return nil
	}
}
