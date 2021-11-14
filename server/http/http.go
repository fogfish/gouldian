package http

import (
	"context"
	µ "github.com/fogfish/gouldian"
	"net/http"
)

//
func Serve(endpoints ...µ.Endpoint) {
	routes := &routes{
		endpoint: µ.Or(endpoints...),
	}

	http.ListenAndServe(":8080", routes)
}

type routes struct{ endpoint µ.Endpoint }

func (routes *routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := &µ.Input{
		Context: µ.NewContext(context.Background()),
		Method:  r.Method,
		Params:  µ.Params(r.URL.Query()),
		Headers: µ.Headers(r.Header),
		Payload: r.Body,
	}

	/*output := */
	routes.endpoint(input)

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
