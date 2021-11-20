package httpd_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/server/httpd"
	"github.com/fogfish/it"
)

func TestServe(t *testing.T) {
	ts := httptest.NewServer(
		httpd.Serve(
			µ.GET(
				µ.Path("echo"),
				µ.FMap(func(ctx µ.Context) error {
					return µ.Status.OK(
						headers.ContentType.Value(headers.TextPlain),
						headers.Server.Value("echo"),
						µ.WithText("echo"),
					)
				}),
			),
		),
	)
	defer ts.Close()

	req, err1 := http.NewRequest("GET", ts.URL+"/echo", nil)
	it.Ok(t).If(err1).Must().Equal(nil)

	out, err2 := http.DefaultClient.Do(req)
	it.Ok(t).If(err2).Must().Equal(nil)

	msg, err3 := ioutil.ReadAll(out.Body)
	it.Ok(t).If(err3).Must().Equal(nil)

	it.Ok(t).
		If(out.Header.Get("Server")).Should().Equal("echo").
		If(out.Header.Get("Content-Type")).Should().Equal("text/plain").
		If(msg).Should().Equal([]byte("echo"))
}
