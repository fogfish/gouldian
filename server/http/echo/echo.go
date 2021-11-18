package main

import (
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/http"
)

func main() {
	http.Serve(echo())
}

type reqEcho struct {
	Echo string
}

var lensEcho = optics.Lenses1(reqEcho{})

func echo() µ.Endpoint {
	return µ.GET(
		µ.Path("echo", lensEcho),
		func(r *µ.Input) error {
			var req reqEcho
			if err := r.Context.Get(&req); err != nil {
				return err
			}

			return µ.Status.
				OK().
				With("Server", "echo").
				With("Content-Type", "text/plain").
				Bytes([]byte(req.Echo))
		},
	)
}
