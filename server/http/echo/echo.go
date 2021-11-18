package main

import (
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
	ƒ "github.com/fogfish/gouldian/output"
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
				return µ.Status.BadRequest(ƒ.Issue(err))
			}

			return µ.Status.OK(
				ƒ.ContentType.Text,
				ƒ.Server.Is("echo"),
				ƒ.Text(req.Echo),
			)
		},
	)
}
