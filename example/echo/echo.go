package main

import (
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080",
		httpd.Serve(
			echo(),
		),
	)

}

type reqEcho struct {
	Echo string
}

var lensEcho = optics.ForProduct1(reqEcho{})

func echo() µ.Endpoint {
	return µ.GET(
		µ.Path("echo", lensEcho),
		µ.FMap(func(ctx µ.Context) error {
			var req reqEcho
			if err := ctx.Get(&req); err != nil {
				return µ.Status.BadRequest(µ.WithIssue(err))
			}

			return µ.Status.OK(
				headers.ContentType.Value(headers.TextPlain),
				headers.Server.Value("echo"),
				µ.WithText(req.Echo),
			)
		}),
	)
}
