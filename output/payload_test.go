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

package output_test

import (
	"bytes"
	"errors"
	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/it/v2"
	"io"
	"strings"
	"testing"
)

func TestPayloads(t *testing.T) {
	spec := []struct {
		Result µ.Result
		Body   string
	}{
		{ø.Send("test"), "test"},
		{ø.Send(strings.NewReader("test")), "test"},
		{ø.Send([]byte("test")), "test"},
		{ø.Send(bytes.NewBuffer([]byte("test"))), "test"},
		{ø.Send(bytes.NewReader([]byte("test"))), "test"},
		{ø.Send(io.LimitReader(strings.NewReader("test"), 4)), "test"},
		{ø.Send(struct {
			T string `json:"t"`
		}{"test"}), `{"t":"test"}`},
	}

	for _, tt := range spec {
		foo := func(*µ.Context) error { return ø.Status.OK(tt.Result) }
		req := mock.Input()
		out, ok := foo(req).(*µ.Output)

		it.Then(t).Should(
			it.True(ok),
			it.Equal(string(out.Body), tt.Body),
		)
	}
}

func TestPayloadCodecs(t *testing.T) {
	spec := []struct {
		Result error
		Body   string
	}{
		{ø.Status.OK(
			ø.ContentType.JSON,
			ø.Send(struct {
				A string `json:"a"`
				B string `json:"b"`
			}{"a", "b"}),
		), `{"a":"a","b":"b"}`},
		{ø.Status.OK(
			ø.ContentType.Form,
			ø.Send(struct {
				A string `json:"a"`
				B string `json:"b"`
			}{"a", "b"}),
		), `a=a&b=b`},
	}

	for _, tt := range spec {
		foo := func(*µ.Context) error { return tt.Result }
		req := mock.Input()
		out, ok := foo(req).(*µ.Output)

		it.Then(t).Should(
			it.True(ok),
			it.Equal(string(out.Body), tt.Body),
		)
	}
}

func TestErrorAttachesCause(t *testing.T) {
	cause := errors.New("test")
	out := ø.Status.BadRequest(ø.Error(cause)).(*µ.Output)

	it.Then(t).
		Should(it.Equal(out.Status, 400)).
		Should(it.True(errors.Is(out, cause))).
		Should(it.Equal(len(out.Body), 0))
}

func TestSendBinary(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0, 0xff, 0xfe}
	out := ø.Status.OK(ø.ContentType.Set("image/png"), ø.Send(png)).(*µ.Output)

	it.Then(t).
		Should(it.Equiv(out.Body, png)).
		Should(it.Equal(out.GetHeader("Content-Length"), "7"))
}

func TestResultErrorPropagates(t *testing.T) {
	boom := errors.New("boom")
	err := ø.Status.OK(func(*µ.Output) error { return boom })

	it.Then(t).Should(it.Equal(err, boom))
}

func TestSerializationFailureIs500(t *testing.T) {
	out := ø.Status.OK(ø.Send(make(chan int))).(*µ.Output)

	it.Then(t).
		Should(it.Equal(out.Status, 500)).
		ShouldNot(it.Nil(out.Failure))
}
