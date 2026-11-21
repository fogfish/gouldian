package gouldian_test

import (
	"errors"

	µ "github.com/fogfish/gouldian/v3"
)

func isFailure(err error) bool {
	var f *µ.Failure
	return errors.As(err, &f)
}
