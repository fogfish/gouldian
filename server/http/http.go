package http

import (
	"context"
	"encoding/json"
	"fmt"
	µ "github.com/fogfish/gouldian"
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
	case *µ.Success:
		routes.sendSuccess(w, v)
	case *µ.Failure:
		routes.sendFailure(w, v)
	case µ.NoMatch:
		err := fmt.Errorf("NoMatch %s", r.URL.Path)
		routes.sendFailure(w, µ.NewFailure(µ.StatusCode(501), err))
	}
}

func (routes *routes) sendSuccess(w http.ResponseWriter, out *µ.Success) {
	for h, v := range out.Headers {
		w.Header().Set(string(h), v)
	}
	w.WriteHeader(int(out.Status))

	if len(out.Body) > 0 {
		w.Write([]byte(out.Body))
	}
}

func (routes *routes) sendFailure(w http.ResponseWriter, out *µ.Failure) {
	w.WriteHeader(int(out.Status))
	if text, err := json.Marshal(out); err == nil {
		w.Write(text)
	}
}
