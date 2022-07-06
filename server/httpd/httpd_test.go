/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package httpd_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/server/httpd"
	"github.com/fogfish/it"
)

func TestServeMatch(t *testing.T) {
	ts := httptest.NewServer(httpd.Serve(mock()))
	defer ts.Close()

	req, err1 := http.NewRequest("GET", ts.URL+"/echo", nil)
	it.Ok(t).If(err1).Must().Equal(nil)

	out, err2 := http.DefaultClient.Do(req)
	it.Ok(t).If(err2).Must().Equal(nil)

	msg, err3 := ioutil.ReadAll(out.Body)
	it.Ok(t).If(err3).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusOK).
		If(out.Header.Get("Server")).Should().Equal("echo").
		If(out.Header.Get("Content-Type")).Should().Equal("text/plain").
		If(msg).Should().Equal([]byte("echo"))
}

func TestServeNoMatch(t *testing.T) {
	ts := httptest.NewServer(httpd.Serve(mock()))
	defer ts.Close()

	req, err1 := http.NewRequest("GET", ts.URL+"/foo", nil)
	it.Ok(t).If(err1).Must().Equal(nil)

	out, err2 := http.DefaultClient.Do(req)
	it.Ok(t).If(err2).Must().Equal(nil)

	msg, err3 := ioutil.ReadAll(out.Body)
	it.Ok(t).If(err3).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusNotImplemented).
		If(out.Header.Get("Content-Type")).Should().Equal("application/json").
		If(msg).ShouldNot().Equal([]byte{})
}

func TestServeUnknownError(t *testing.T) {
	ts := httptest.NewServer(
		httpd.Serve(
			µ.GET(µ.URI(), func(c *µ.Context) error { return fmt.Errorf("Fail") }),
		),
	)
	defer ts.Close()

	req, err1 := http.NewRequest("GET", ts.URL+"/", nil)
	it.Ok(t).If(err1).Must().Equal(nil)

	out, err2 := http.DefaultClient.Do(req)
	it.Ok(t).If(err2).Must().Equal(nil)

	msg, err3 := ioutil.ReadAll(out.Body)
	it.Ok(t).If(err3).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusInternalServerError).
		If(out.Header.Get("Content-Type")).Should().Equal("application/json").
		If(msg).ShouldNot().Equal([]byte{})
}

func TestServeAndCommit(t *testing.T) {
	cnt := 0
	ts := httptest.NewServer(
		httpd.ServeAndCommit(
			func() { cnt = cnt + 1 },
			mock(),
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
		If(out.StatusCode).Should().Equal(http.StatusOK).
		If(out.Header.Get("Server")).Should().Equal("echo").
		If(out.Header.Get("Content-Type")).Should().Equal("text/plain").
		If(msg).Should().Equal([]byte("echo")).
		If(cnt).Should().Equal(1)
}

func mock() µ.Routable {
	return µ.GET(
		µ.URI(µ.Path("echo")),
		func(ctx *µ.Context) error {
			return µ.Status.OK(
				µ.WithHeader(headers.ContentType, headers.TextPlain),
				µ.WithHeader(headers.Server, "echo"),
				µ.WithText("echo"),
			)
		},
	)

}
