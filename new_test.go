package gouldian_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/path"
)

func TestXGet(t *testing.T) {
	req := mock.Input(mock.URL("/10"))
	var i int
	var s string

	x := path.OR(path.Int(&i), path.String(&s))

	end := gouldian.GET(
		gouldian.Path(x),
	)
	val := end.IsMatch(req)
	fmt.Println("====> ")
	fmt.Println(val)
	fmt.Println(i)
	fmt.Println(s)
}
