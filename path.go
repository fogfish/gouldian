package gouldian

import (
	"github.com/fogfish/gouldian/optics"
)

const (
	//
	Any = "_"
)

//
type pathArrow func(Context, string) error

/*

Path is an endpoint to match URL of HTTP request. The function takes
url matching primitives, which are defined by the package `path`.

  import "github.com/fogfish/gouldian/path"

  e := µ.GET( µ.Path(path.Is("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Path(arrows ...interface{}) Endpoint {
	return mkPathEndpoint(mkPathMatcher(arrows))
}

func mkPathEndpoint(farrows []pathArrow) Endpoint {
	return func(req *Input) error {
		if len(req.Resource) != len(farrows) {
			return NoMatch{}
		}

		ctx := req.Context
		for i, f := range farrows {
			if err := f(ctx, req.Resource[i]); err != nil {
				return err
			}
		}

		return nil
	}
}

/*

Prefix is an endpoint to match URL of HTTP request. The function takes
url matching primitives, which are defined by the package `path`.

  import "github.com/fogfish/gouldian/path"

  e := µ.GET( µ.Path(path.Is("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
// func Prefix(arrows ...interface{}) Endpoint {
// 	farrows := mkPathMatcher(arrows)

// 	return func(req *Input) error {
// 		if len(req.Resource) < len(farrows) {
// 			return NoMatch{}
// 		}

// 		ctx := req.Context
// 		for i, f := range farrows {
// 			if err := f(ctx, req.Resource[i]); err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	}
// }

func mkPathMatcher(arrows []interface{}) []pathArrow {
	seq := make([]pathArrow, len(arrows))

	for i, arrow := range arrows {
		switch v := arrow.(type) {
		case string:
			switch v {
			case Any:
				seq[i] = pathAny()
			default:
				seq[i] = pathIs(v)
			}
		case optics.Lens:
			seq[i] = pathTo(v)
		default:
			seq[i] = pathNone()
		}
	}

	return seq
}

/*

Is matches a path segment to defined literal
  e := µ.GET( µ.Path(path.Is("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func pathIs(val string) pathArrow {
	return func(ctx Context, segment string) error {
		if segment == val {
			return nil
		}
		return NoMatch{}
	}
}

/*

Any is a wildcard matcher of path segment
  e := µ.GET( µ.Path(path.Any()) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) == nil
*/
func pathNone() pathArrow {
	return func(Context, string) error {
		return NoMatch{}
	}
}

/*

Any is a wildcard matcher of path segment
  e := µ.GET( µ.Path(path.Any()) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) == nil
*/
func pathAny() pathArrow {
	return func(Context, string) error {
		return nil
	}
}

/*

Lifts the path segment to lens
*/
func pathTo(l optics.Lens) pathArrow {
	return func(ctx Context, segment string) error {
		return ctx.Put(l, segment)
	}
}
