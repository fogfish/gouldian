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
	Put(optics.Lens, ...string) error
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
		morphism: make(optics.Morphism),
	}
}

/*

Free ...
*/
func (ctx *µContext) Free() {
	ctx.morphism = make(optics.Morphism)
}

/*

Put ...
*/
func (ctx *µContext) Put(lens optics.Lens, raw ...string) error {
	val, err := lens.FromSeq(raw)
	if err != nil {
		return NoMatch{}
	}

	ctx.morphism[lens] = val
	return nil
}

/*

Get ...
*/
func (ctx *µContext) Get(val interface{}) error {
	return ctx.morphism.Apply(val)
}
