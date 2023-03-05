package emitter_test

import (
	"bytes"
	"errors"
	µ "github.com/fogfish/gouldian"
	ø "github.com/fogfish/gouldian/emitter"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/guid"
	"github.com/fogfish/it/v2"
	"io"
	"strings"
	"testing"
)

func TestPayloads(t *testing.T) {
	guid.Clock = guid.NewLClock(
		guid.ConfNodeID(0),
		guid.ConfClock(func() uint64 { return 0 }),
	)

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
		{ø.Error(errors.New("test")), `{"instance":"N..............1","type":"https://httpstatuses.com/200","status":200,"title":"OK"}`},
	}

	for _, tt := range spec {
		foo := func(*µ.Context) error { return ø.Status.OK(tt.Result) }
		req := mock.Input()
		out, ok := foo(req).(*µ.Output)

		it.Then(t).Should(
			it.True(ok),
			it.Equal(out.Body, tt.Body),
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
			it.Equal(out.Body, tt.Body),
		)
	}
}
