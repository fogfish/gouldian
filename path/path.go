package path

import (
	"strconv"

	"github.com/fogfish/gouldian/core"
)

//
type Arrow func(val string) error

//
func Is(val string) Arrow {
	return func(segment string) error {
		if segment == val {
			return nil
		}
		return core.NoMatch{}
	}
}

//
func Any() Arrow {
	return func(string) error {
		return nil
	}
}

//
func String(val *string) Arrow {
	return func(segment string) error {
		*val = segment
		return nil
	}
}

//
func Int(val *int) Arrow {
	return func(segment string) error {
		if value, err := strconv.Atoi(segment); err == nil {
			*val = value
			return nil
		}
		return core.NoMatch{}
	}
}
