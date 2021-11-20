package gouldian

import (
	"path/filepath"

	"github.com/fogfish/gouldian/optics"
)

//
type pathArrow func(Context, string) error

/*

Path is an endpoint to match URL of HTTP request. The function takes a path
pattern as arguments. The pattern is sequence of either literals or lenses,
where each term corresponds to the path segment. The function do not match
if length of path is not equal to the length of pattern or segment do not
match to pattern

  e := µ.GET( µ.Path("foo") )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Path(arrows ...interface{}) Endpoint {
	return mkPathEndpoint(mkPathMatcher(arrows))
}

func mkPathEndpoint(farrows []pathArrow) Endpoint {
	return func(req *Input) error {
		if len(farrows) != len(req.Resource) {
			return NoMatch{}
		}

		for i, f := range farrows {
			if err := f(req.Context, req.Resource[i]); err != nil {
				return err
			}
		}
		return nil
	}
}

/*

PathSeq is like Path but last element in the pattern must be lens that lifts
the tail of path.

  e := µ.GET( µ.PathSeq("foo", suffix) )
  e(mock.Input(mock.URL("/foo/bar"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func PathSeq(arrows ...interface{}) Endpoint {
	return mkPathSeqEndpoint(mkPathMatcher(arrows))
}

func mkPathSeqEndpoint(farrows []pathArrow) Endpoint {
	return func(req *Input) error {
		if len(farrows) > len(req.Resource) || len(farrows) < 1 {
			return NoMatch{}
		}

		last := len(farrows) - 1
		for i := 0; i < last; i++ {
			if err := farrows[i](req.Context, req.Resource[i]); err != nil {
				return err
			}
		}

		if err := farrows[last](req.Context, filepath.Join(req.Resource[last:]...)); err != nil {
			return err
		}

		return nil
	}
}

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
  e := µ.GET( µ.Path("foo") )
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

None matches nothing
*/
func pathNone() pathArrow {
	return func(Context, string) error {
		return NoMatch{}
	}
}

/*

Any is a wildcard matcher of path segment
  e := µ.GET( µ.Path(path.Any) )
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
