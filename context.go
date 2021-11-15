package gouldian

import (
	"context"
	"io"

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
	PutStream(optics.Lens, io.Reader) error
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
		morphism: make(optics.Morphism, 0, 30),
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

	ctx.morphism = append(ctx.morphism, optics.Arrow{Lens: lens, Value: val})
	return nil
}

func (ctx *µContext) PutStream(lens optics.Lens, r io.Reader) error {
	// TODO
	return nil
}

/*

Get ...
*/
func (ctx *µContext) Get(val interface{}) error {
	return ctx.morphism.Apply(val)
}
