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
	"io/ioutil"
	"net/http"

	"github.com/fogfish/gouldian/optics"
)

/*

Context of HTTP request. The context accumulates matched terms of HTTP and
passes it to destination function.
*/
type Context struct {
	context.Context

	Request *http.Request
	values  []string
	params  Params
	payload []byte

	JWT JWT

	morphism optics.Morphisms
}

/*

NewContext create a new context for HTTP request
*/
func NewContext(ctx context.Context) *Context {
	return &Context{
		Context:  ctx,
		values:   make([]string, 0, 20),
		morphism: make(optics.Morphisms, 0, 20),
	}
}

/*

Free the context
*/
func (ctx *Context) free() {
	ctx.values = ctx.values[:0]
	ctx.morphism = ctx.morphism[:0]
}

/*

Free the context
*/
func (ctx *Context) Free() {
	ctx.JWT = nil
	ctx.params = nil
	ctx.payload = nil
	ctx.Request = nil
	ctx.values = ctx.values[:0]
	ctx.morphism = ctx.morphism[:0]
}

/*

Put injects value to the context
*/
func (ctx *Context) Put(lens optics.Lens, str string) error {
	val, err := lens.FromString(str)
	if err != nil {
		return ErrNoMatch
	}

	ctx.morphism = append(ctx.morphism, optics.Morphism{Lens: lens, Value: val})
	// optics.Setter{Lens: lens, Value: val})
	return nil
}

/*

Get decodes context into structure
*/
func (ctx *Context) Get(val interface{}) error {
	if err := ctx.morphism.Apply(val); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) cacheBody() error {
	if ctx.Request.Body != nil {
		buf, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return err
		}
		// This is copied from runtime. It relies on the string
		// header being a prefix of the slice header!
		// ctx.payload = *(*string)(unsafe.Pointer(&buf))
		ctx.payload = buf
		return nil
	}

	return nil
}
