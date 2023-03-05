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

package emitter_test

import (
	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/gouldian/v2/mock"
	ø "github.com/fogfish/gouldian/v2/output"
	"github.com/fogfish/it/v2"
	"testing"
	"time"
)

func TestHeaders(t *testing.T) {
	spec := []struct {
		Result µ.Result
		Header string
		Value  string
	}{
		{ø.Age.Set(1024), string(ø.Age), "1024"},
		{ø.CacheControl.Set("nocache"), string(ø.CacheControl), "nocache"},
		{ø.Connection.Set("keep-alive"), string(ø.Connection), "keep-alive"},
		{ø.Connection.KeepAlive, string(ø.Connection), "keep-alive"},
		{ø.Connection.Close, string(ø.Connection), "close"},
		{ø.ContentEncoding.Set("foo"), string(ø.ContentEncoding), "foo"},
		{ø.ContentLanguage.Set("foo"), string(ø.ContentLanguage), "foo"},
		{ø.ContentLength.Set(1024), string(ø.ContentLength), "1024"},
		{ø.Header("Content-Length", 1024), string(ø.ContentLength), "1024"},
		{ø.ContentLocation.Set("foo"), string(ø.ContentLocation), "foo"},
		{ø.ContentMD5.Set("foo"), string(ø.ContentMD5), "foo"},
		{ø.ContentRange.Set("foo"), string(ø.ContentRange), "foo"},
		{ø.ContentType.Set("foo"), string(ø.ContentType), "foo"},
		{ø.ContentType.ApplicationJSON, string(ø.ContentType), "application/json"},
		{ø.ContentType.JSON, string(ø.ContentType), "application/json"},
		{ø.ContentType.Form, string(ø.ContentType), "application/x-www-form-urlencoded"},
		{ø.ContentType.TextPlain, string(ø.ContentType), "text/plain"},
		{ø.ContentType.Text, string(ø.ContentType), "text/plain"},
		{ø.ContentType.TextHTML, string(ø.ContentType), "text/html"},
		{ø.ContentType.HTML, string(ø.ContentType), "text/html"},
		{ø.Date.Set(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(ø.Date), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{ø.Header("Date", time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(ø.Date), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{ø.ETag.Set("foo"), string(ø.ETag), "foo"},
		{ø.Expires.Set(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(ø.Expires), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{ø.LastModified.Set(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(ø.LastModified), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{ø.Link.Set("foo"), string(ø.Link), "foo"},
		{ø.Location.Set("foo"), string(ø.Location), "foo"},
		{ø.RetryAfter.Set(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(ø.RetryAfter), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{ø.Server.Set("foo"), string(ø.Server), "foo"},
		{ø.SetCookie.Set("foo"), string(ø.SetCookie), "foo"},
		{ø.TransferEncoding.Set("chunked"), string(ø.TransferEncoding), "chunked"},
		{ø.TransferEncoding.Chunked, string(ø.TransferEncoding), "chunked"},
		{ø.TransferEncoding.Identity, string(ø.TransferEncoding), "identity"},
		{ø.Via.Set("foo"), string(ø.Via), "foo"},
		{ø.Header("X-Value", "foo"), "X-Value", "foo"},
		{ø.Header("x-value", "foo"), "X-Value", "foo"},
	}

	for _, tt := range spec {
		foo := func(*µ.Context) error { return ø.Status.OK(tt.Result) }
		req := mock.Input()
		out, ok := foo(req).(*µ.Output)

		it.Then(t).Should(
			it.True(ok),
			it.Equal(out.GetHeader(tt.Header), tt.Value),
		)
	}
}
