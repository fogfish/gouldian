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

	"github.com/fogfish/gouldian/optics"
)

/*

Context of HTTP request. The context accumulates matched terms of HTTP and
passes it to destination function.
*/
type Context interface {
	context.Context

	// Free all resource allocated by the context
	Free()

	// Put and Get are optics function store and lifts HTTP terms
	Put(optics.Lens, string) error
	Get(interface{}) error
}

/*

µContext is an internal implementation of the Context interface
*/
type µContext struct {
	context.Context

	morphism optics.Morphism
}

var (
	_ Context = (*µContext)(nil)
)

/*

NewContext create a new context for HTTP request
*/
func NewContext(ctx context.Context) Context {
	return &µContext{
		Context:  ctx,
		morphism: make(optics.Morphism, 0, 10),
	}
}

/*

Free ...
*/
func (ctx *µContext) Free() {
	ctx.morphism = ctx.morphism[:0]
}

/*

Put ...
*/
func (ctx *µContext) Put(lens optics.Lens, str string) error {
	val, err := lens.FromString(str)
	if err != nil {
		return NoMatch{}
	}

	ctx.morphism = append(ctx.morphism, optics.Setter{Lens: lens, Value: val})
	return nil
}

/*

Get ...
*/
func (ctx *µContext) Get(val interface{}) error {
	return ctx.morphism.Apply(val)
}
