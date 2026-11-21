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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/fogfish/gouldian/v3/internal/optics"
)

// DefaultBodyLimit is the default limit of the request body (6 MiB,
// the payload limit of AWS Lambda).
const DefaultBodyLimit int64 = 6 << 20

/*
Context of HTTP request. The context accumulates matched terms of HTTP and
passes it to destination function. It embeds context.Context of the caller
(http.Request, AWS Lambda invocation), so the cancellation, the deadline and
the trace identity are available to the handlers.
*/
type Context struct {
	context.Context

	Request *http.Request
	// RequestID is the identity of the request assigned by the platform
	RequestID string
	JWT       Token

	values    []string
	allow     []string
	params    Query
	payload   []byte
	bodyRead  bool
	bodyLimit int64
	log       *slog.Logger

	morphism optics.Morphisms
}

/*
NewContext create a new context for HTTP request
*/
func NewContext(ctx context.Context) *Context {
	return &Context{
		Context:   ctx,
		values:    make([]string, 0, 20),
		morphism:  make(optics.Morphisms, 0, 20),
		bodyLimit: DefaultBodyLimit,
	}
}

// reset the context before routing
func (ctx *Context) free() {
	ctx.values = ctx.values[:0]
	ctx.morphism = ctx.morphism[:0]
}

/*
Free the context
*/
func (ctx *Context) Free() {
	ctx.JWT = nil
	ctx.RequestID = ""
	ctx.params = nil
	ctx.payload = nil
	ctx.bodyRead = false
	ctx.allow = ctx.allow[:0]
	ctx.Request = nil
	ctx.values = ctx.values[:0]
	ctx.morphism = ctx.morphism[:0]
}

// rollback removes bindings captured by a rejected branch
func (ctx *Context) rollback(mark int) {
	if len(ctx.morphism) > mark {
		clear(ctx.morphism[mark:])
		ctx.morphism = ctx.morphism[:mark]
	}
}

func (ctx *Context) logger() *slog.Logger {
	if ctx.log != nil {
		return ctx.log
	}
	return slog.Default()
}

// Query returns the parsed query string
func (ctx *Context) Query() Query {
	if ctx.params == nil {
		ctx.params = Query(ctx.Request.URL.Query())
	}
	return ctx.params
}

// Put injects value to the context. A value that cannot be decoded is
// a Failure.
func (ctx *Context) Put(lens Lens, str string) error {
	return ctx.put(SourceUnknown, "", lens, str)
}

func (ctx *Context) put(src Source, key string, lens Lens, str string) error {
	val, err := lens.FromString(str)
	if err != nil {
		return newFailure(src, key, err)
	}

	ctx.morphism = append(ctx.morphism, optics.Morphism{Lens: lens.Lens, Value: val})
	return nil
}

// putMulti injects all values of the repeated param
func (ctx *Context) putMulti(src Source, key string, lens Lens, vals []string) error {
	val, err := lens.FromStrings(vals)
	if err != nil {
		return newFailure(src, key, err)
	}

	ctx.morphism = append(ctx.morphism, optics.Morphism{Lens: lens.Lens, Value: val})
	return nil
}

// putPath injects value of path segment. The segment that cannot be decoded
// does not match the route (e.g. /user/:id can be overloaded by int and string)
func (ctx *Context) putPath(lens Lens, str string) error {
	val, err := lens.FromString(str)
	if err != nil {
		return NoMatchPath
	}

	ctx.morphism = append(ctx.morphism, optics.Morphism{Lens: lens.Lens, Value: val})
	return nil
}

func (ctx *Context) putBody(lens Lens, buf []byte) error {
	val, err := lens.FromBytes(buf)
	if err != nil {
		key := bodyPointer(err)
		var syntax *json.SyntaxError
		if errors.As(err, &syntax) {
			err = fmt.Errorf("%w (offset %d)", err, syntax.Offset)
		}
		return newFailure(SourceBody, key, err)
	}

	ctx.morphism = append(ctx.morphism, optics.Morphism{Lens: lens.Lens, Value: val})
	return nil
}

// FromContext decodes the context into structure
func FromContext[S any](ctx *Context, val *S) error {
	return optics.Morph(ctx.morphism, val)
}

// body reads the request body once, all endpoints share the cached bytes.
func (ctx *Context) body() ([]byte, error) {
	if ctx.bodyRead {
		return ctx.payload, nil
	}

	ctx.bodyRead = true
	req := ctx.Request
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}

	limit := ctx.bodyLimit
	if limit > 0 && req.ContentLength > limit {
		return nil, &Failure{Status: http.StatusRequestEntityTooLarge, Source: SourceBody, Err: ErrBodyTooLarge}
	}

	var r io.Reader = req.Body
	if limit > 0 {
		r = io.LimitReader(r, limit+1)
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, newFailure(SourceBody, "", err)
	}

	if limit > 0 && int64(len(buf)) > limit {
		return nil, &Failure{Status: http.StatusRequestEntityTooLarge, Source: SourceBody, Err: ErrBodyTooLarge}
	}

	ctx.payload = buf
	return buf, nil
}
