package gouldian

import (
	"context"

	"github.com/fogfish/gouldian/optics"
)

/*

Context ...
*/
type Context interface {
	context.Context

	// Free all resource allocated by the context
	Free()

	Put(optics.Lens, interface{})
	Get(interface{}) error
}

/*

HList

µContext ...
*/
type µContext struct {
	context.Context

	morphism optics.Morphism
}

var (
	_ Context = (*µContext)(nil)
)

/*

NewContext ...
*/
func NewContext(ctx context.Context) Context {
	return &µContext{
		Context:  ctx,
		morphism: make(optics.Morphism),
	}
}

//
func (ctx *µContext) Free() {
	ctx.morphism = make(optics.Morphism)
}

//
func (ctx *µContext) Put(lens optics.Lens, val interface{}) {
	ctx.morphism[lens] = val
}

//
func (ctx *µContext) Get(val interface{}) error {
	return ctx.morphism.Apply(val)
}
