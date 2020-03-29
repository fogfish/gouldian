//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package main

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fogfish/gouldian"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/param"
	"github.com/fogfish/gouldian/path"
)

type headers struct {
	Accept    string `json:"Accept,omitempty"`
	Host      string `json:"Host,omitempty"`
	Origin    string `json:"Origin,omitempty"`
	Referer   string `json:"Referer,omitempty"`
	UserAgent string `json:"User-Agent,omitempty"`
}

type response struct {
	Headers headers `json:"headers,omitempty"`
	URL     string  `json:"url,omitempty"`
}

//-----------------------------------------------------------------------------
//
// Http Methods
//
//-----------------------------------------------------------------------------

func anyMethod(subpath string) µ.Endpoint {
	var h headers

	return µ.Join(
		µ.Path(path.Is(subpath)),
		µ.Header(
			header.MaybeString("Accept", &h.Accept),
			header.MaybeString("Host", &h.Host),
			header.MaybeString("Origin", &h.Origin),
			header.MaybeString("Referer", &h.Referer),
			header.MaybeString("User-Agent", &h.UserAgent),
		),
		µ.FMap(
			func() error {
				return µ.Ok().
					JSON(response{h, fmt.Sprintf("https://%v/%v", h.Host, subpath)})
			},
		),
	)
}

func delete() µ.Endpoint {
	return µ.DELETE(anyMethod("delete"))
}

func get() µ.Endpoint {
	return µ.GET(anyMethod("get"))
}

func patch() µ.Endpoint {
	return µ.PATCH(anyMethod("patch"))
}

func post() µ.Endpoint {
	return µ.POST(anyMethod("post"))
}

func put() µ.Endpoint {
	return µ.PUT(anyMethod("put"))
}

//-----------------------------------------------------------------------------
//
// Auth
//
//-----------------------------------------------------------------------------

func bearer() µ.Endpoint {
	var token string
	return µ.GET(
		µ.Path(path.Is("bearer")),
		µ.Header(header.String("Authorization", &token)),
		µ.FMap(
			func() error {
				return gouldian.Unauthorized("Invalid token: " + token)
			},
		),
	)
}

//-----------------------------------------------------------------------------
//
// Status codes
//
//-----------------------------------------------------------------------------

func status() µ.Endpoint {
	var code int
	return µ.GET(
		µ.Path(path.Is("status"), path.Int(&code)),
		µ.FMap(
			func() error {
				return gouldian.Success(code)
			},
		),
	)
}

//-----------------------------------------------------------------------------
//
// Request inpsection
//
//-----------------------------------------------------------------------------

func head() µ.Endpoint {
	var h headers

	return µ.GET(
		µ.Path(path.Is("headers")),
		µ.Header(
			header.MaybeString("Accept", &h.Accept),
			header.MaybeString("Host", &h.Host),
			header.MaybeString("Origin", &h.Origin),
			header.MaybeString("Referer", &h.Referer),
			header.MaybeString("User-Agent", &h.UserAgent),
		),
		µ.FMap(
			func() error {
				return µ.Ok().JSON(h)
			},
		),
	)
}

func ip() µ.Endpoint {
	var ip string
	return µ.GET(
		µ.Path(path.Is("ip")),
		µ.Header(header.String("X-Forwarded-For", &ip)),
		µ.FMap(
			func() error {
				return gouldian.Ok().JSON(ip)
			},
		),
	)
}

func ua() µ.Endpoint {
	var ua string
	return µ.GET(
		µ.Path(path.Is("user-agent")),
		µ.Header(header.String("User-Agent", &ua)),
		µ.FMap(
			func() error {
				return gouldian.Ok().JSON(ip)
			},
		),
	)
}

//-----------------------------------------------------------------------------
//
// Redirects
//
//-----------------------------------------------------------------------------

func redirect1() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("redirect"), path.Is("1")),
		µ.FMap(
			func() error {
				return µ.Found(url.URL{Scheme: "https", Host: "example.com"})
			},
		),
	)
}

func redirectN() µ.Endpoint {
	var host string
	var n int
	return µ.GET(
		µ.Path(path.Is("redirect"), path.Int(&n)),
		µ.Header(header.String("Host", &host)),
		µ.FMap(
			func() error {
				return µ.Found(url.URL{Scheme: "https", Host: host, Path: fmt.Sprintf("/api/redirect/%v", n-1)})
			},
		),
	)
}

func redirectTo() µ.Endpoint {
	var to string
	return µ.GET(
		µ.Path(path.Is("redirect-to")),
		µ.Param(param.String("url", &to)),
		µ.FMap(
			func() error {
				redirect, err := url.Parse(to)
				if err == nil {
					return gouldian.Found(*redirect)
				}
				return gouldian.BadRequest("Invalid url: " + to)
			},
		),
	)
}

func main() {
	lambda.Start(
		gouldian.Serve(
			delete(), get(), patch(), post(), put(),
			bearer(),
			status(),
			head(), ip(), ua(),
			redirect1(), redirectN(), redirectTo(),
		),
	)
}
