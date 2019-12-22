package gouldian

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// TODO: NotImplemented
type NoMatch struct{}

func (err NoMatch) Error() string {
	return fmt.Sprintf("No Match")
}

//
type Matcher struct {
	f Endpoint
}

//
func Delete() Pattern {
	return &Matcher{isVerb("DELETE")}
}

func Get() Pattern {
	return &Matcher{isVerb("GET")}
}

func Patch() Pattern {
	return &Matcher{isVerb("PATCH")}
}

func Post() Pattern {
	return &Matcher{isVerb("POST")}
}

func Put() Pattern {
	return &Matcher{isVerb("PUT")}
}

func isVerb(verb string) Endpoint {
	return func(http *Input) error {
		if http.HTTPMethod == verb {
			return nil
		}
		return NoMatch{}
	}
}

//
func (state *Matcher) Path(segment string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		if len(req.path) > req.segment && req.path[req.segment] == segment {
			req.segment++
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) String(val *string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		if len(req.path) > req.segment {
			*val = req.path[req.segment]
			req.segment++
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) Int(val *int) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		if len(req.path) > req.segment {
			value, err := strconv.Atoi(req.path[req.segment])
			if err != nil {
				return NoMatch{}
			}
			*val = value
			req.segment++
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) Opt(name string, val string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists && opt == val {
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) IsOpt(name string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		_, exists := req.QueryStringParameters[name]
		if exists {
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) OptString(name string, val *string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists {
			*val = opt
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) OptInt(name string, val *int) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists {
			value, err := strconv.Atoi(opt)
			if err != nil {
				return NoMatch{}
			}
			*val = value
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) Head(name string, val string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		head, exists := req.Headers[name]
		if exists && strings.HasPrefix(head, val) {
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) IsHead(name string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		_, exists := req.Headers[name]
		if exists {
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) HeadString(name string, val *string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		head, exists := req.Headers[name]
		if exists {
			*val = head
		}
		return nil
	})
	return state
}

//
func (state *Matcher) Json(val interface{}) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		err := json.Unmarshal([]byte(req.Body), val)
		if err == nil {
			return nil
		}
		return NoMatch{}
	})
	return state
}

//
func (state *Matcher) Text(val *string) Pattern {
	state.f = state.f.Then(func(req *Input) error {
		*val = req.Body
		return nil
	})
	return state
}

//
func (state *Matcher) IsMatch(in *Input) bool {
	return state.f(in) == nil
}

//
func (state *Matcher) Then(f func() error) Endpoint {
	return state.f.Then(func(req *Input) error { return f() })
}
