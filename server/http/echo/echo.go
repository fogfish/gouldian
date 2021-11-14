package main

import (
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/server/http"
)

func main() {
	http.Serve(echo())
}

func echo() µ.Endpoint {
	return µ.GET(
		µ.Path("echo"),
		µ.FMap(func(c µ.Context) error {
			return µ.Status.OK().Bytes([]byte("Hello"))
		}),
	)
}
