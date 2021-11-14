package http

import (
	µ "github.com/fogfish/gouldian"
	"net/http"
)

func Serve(endpoints ...µ.Endpoint) {
	http.ListenAndServe(":8080", X)
}

func X(w http.ResponseWriter, r *http.Request) {

}
