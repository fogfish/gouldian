package http

import (
	"context"
	"fmt"
	µ "github.com/fogfish/gouldian"
	ƒ "github.com/fogfish/gouldian/output"
	"net/http"
	"strings"
)

//
func Serve(endpoints ...µ.Endpoint) error {
	routes := &routes{
		endpoint: µ.Or(endpoints...),
	}

	return http.ListenAndServe(":8080", routes)
}

type routes struct{ endpoint µ.Endpoint }

func (routes *routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := &µ.Input{
		Context:  µ.NewContext(context.Background()),
		Method:   r.Method,
		Resource: strings.Split(r.URL.Path, "/")[1:],
		Params:   µ.Params(r.URL.Query()),
		Headers:  µ.Headers(r.Header),
		Stream:   r.Body,
	}

	switch v := routes.endpoint(input).(type) {
	case *µ.Output:
		routes.output(w, v)
	case µ.NoMatch:
		failure := µ.Status.NotImplemented(
			ƒ.Issue(fmt.Errorf("NoMatch %s", r.URL.Path)),
		).(*µ.Output)
		routes.output(w, failure)
	}
}

func (routes *routes) output(w http.ResponseWriter, out *µ.Output) {
	for h, v := range out.Headers {
		w.Header().Set(string(h), v)
	}
	w.WriteHeader(int(out.Status))

	if len(out.Body) > 0 {
		w.Write([]byte(out.Body))
	}
}
