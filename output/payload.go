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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	µ "github.com/fogfish/gouldian/v3"
)

func Send(data any) µ.Result {
	return func(out *µ.Output) error {
		chunked := out.GetHeader(string(TransferEncoding)) == "chunked"
		content := out.GetHeader(string(ContentType))

		var body []byte
		switch stream := data.(type) {
		case string:
			body = []byte(stream)
		case []byte:
			body = stream
		case io.Reader:
			val, err := io.ReadAll(stream)
			if err != nil {
				return err
			}
			body = val
		default:
			val, err := encode(content, data)
			if err != nil {
				out.Status = http.StatusInternalServerError
				out.SetIssue(fmt.Errorf("serialization is failed for <%T>: %w", data, err))
				return nil
			}
			body = val
		}

		if !chunked {
			out.SetHeader(string(ContentLength), strconv.Itoa(len(body)))
		}
		out.Body = body
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

// Error attaches the cause of failure to the output. The server logs the cause
// and renders RFC 9457 Problem Details.
func Error(failure error, title ...string) µ.Result {
	return func(out *µ.Output) error {
		out.SetIssue(failure, title...)
		return nil
	}
}
