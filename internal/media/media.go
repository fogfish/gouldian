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

// Package media helps adapters to serialize responses
package media

import (
	"mime"
	"strings"
	"unicode/utf8"
)

// IsBinary checks if the response body must be base64 encoded: the media type
// is not textual or the payload is not valid UTF-8.
func IsBinary(contentType string, body []byte) bool {
	if len(body) == 0 {
		return false
	}

	if contentType == "" {
		return !utf8.Valid(body)
	}

	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return !utf8.Valid(body)
	}

	if isText(strings.ToLower(mt)) {
		return !utf8.Valid(body)
	}

	return true
}

func isText(mt string) bool {
	switch {
	case strings.HasPrefix(mt, "text/"):
		return true
	case strings.HasSuffix(mt, "+json"), strings.HasSuffix(mt, "+xml"):
		return true
	}

	switch mt {
	case "application/json", "application/xml", "application/javascript",
		"application/x-www-form-urlencoded", "application/x-ndjson",
		"application/yaml", "application/x-yaml", "application/graphql":
		return true
	}

	return false
}
