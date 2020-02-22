package gouldian

import (
	"github.com/fogfish/gouldian/path"
)

//
type Arrow func(*APIGateway) *APIGateway

//
func Path(arrows ...path.Arrow) Endpoint {
	return func(req *Input) error {
		for i, f := range arrows {
			if i+1 > len(req.path) {
				return NoMatch{}
			}
			return f(req.path[i+1])
		}
		return nil
	}
}

//
func XGet(arrows ...Endpoint) *APIGateway {
	state := &APIGateway{isVerb("GET")}
	for _, f := range arrows {
		state.f = state.f.Then(f)
	}
	return state
}
