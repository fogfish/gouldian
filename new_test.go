package gouldian_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/path"
)

func TestXGet(t *testing.T) {
	req := gouldian.Mock("/10")
	var str int

	end := gouldian.XGet(
		gouldian.Path(path.Int(&str)),
	)
	val := end.IsMatch(req)
	fmt.Println("====> ")
	fmt.Println(val)
	fmt.Println(str)

}
