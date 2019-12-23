package main

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fogfish/gouldian"
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

func delete() gouldian.Endpoint {
	return verb(gouldian.Delete(), "delete")
}

func get() gouldian.Endpoint {
	return verb(gouldian.Get(), "get")
}

func patch() gouldian.Endpoint {
	return verb(gouldian.Patch(), "patch")
}

func post() gouldian.Endpoint {
	return verb(gouldian.Patch(), "post")
}

func put() gouldian.Endpoint {
	return verb(gouldian.Patch(), "put")
}

func verb(method gouldian.HTTP, path string) gouldian.Endpoint {
	h := headers{}

	return method.Path(path).
		HString("Accept", &h.Accept).
		HString("Host", &h.Host).
		HString("Origin", &h.Origin).
		HString("Referer", &h.Referer).
		HString("User-Agent", &h.UserAgent).
		FMap(func() error {
			return gouldian.Ok().
				JSON(response{h, fmt.Sprintf("https://%v/%v", h.Host, path)})
		})
}

//-----------------------------------------------------------------------------
//
// Auth
//
//-----------------------------------------------------------------------------

func bearer() gouldian.Endpoint {
	var token string
	return gouldian.Get().Path("bearer").
		HString("Authorization", &token).
		FMap(func() error {
			return gouldian.Unauthorized("Invalid token: " + token)
		})
}

//-----------------------------------------------------------------------------
//
// Status codes
//
//-----------------------------------------------------------------------------

func status() gouldian.Endpoint {
	var code int
	return gouldian.Get().Path("status").Int(&code).
		FMap(func() error {
			return gouldian.Success(code)
		})
}

//-----------------------------------------------------------------------------
//
// Request inpsection
//
//-----------------------------------------------------------------------------

func header() gouldian.Endpoint {
	h := headers{}

	return gouldian.Get().Path("headers").
		HString("Accept", &h.Accept).
		HString("Host", &h.Host).
		HString("Origin", &h.Origin).
		HString("Referer", &h.Referer).
		HString("User-Agent", &h.UserAgent).
		FMap(func() error {
			return gouldian.Ok().JSON(h)
		})
}

func ip() gouldian.Endpoint {
	var ip string
	return gouldian.Get().Path("ip").
		HString("X-Forwarded-For", &ip).
		FMap(func() error {
			return gouldian.Ok().JSON(ip)
		})
}

func ua() gouldian.Endpoint {
	var ua string
	return gouldian.Get().Path("user-agent").
		HString("User-Agent", &ua).
		FMap(func() error {
			return gouldian.Ok().JSON(ua)
		})
}

//-----------------------------------------------------------------------------
//
// Redirects
//
//-----------------------------------------------------------------------------

func redirect1() gouldian.Endpoint {
	return gouldian.Get().Path("redirect").Path("1").
		FMap(func() error {
			return gouldian.Found(
				url.URL{Scheme: "https", Host: "example.com"},
			)
		})
}

func redirectN() gouldian.Endpoint {
	var host string
	var n int
	return gouldian.Get().Path("redirect").Int(&n).
		HString("Host", &host).
		FMap(func() error {
			return gouldian.Found(
				url.URL{Scheme: "https", Host: host, Path: fmt.Sprintf("/api/redirect/%v", n-1)},
			)
		})
}

func redirectTo() gouldian.Endpoint {
	var to string
	return gouldian.Get().Path("redirect-to").
		QString("url", &to).
		FMap(func() error {
			redirect, err := url.Parse(to)
			if err == nil {
				return gouldian.Found(*redirect)
			}
			return gouldian.BadRequest("Invalid url: " + to)
		})
}

func main() {
	lambda.Start(
		gouldian.Serve(
			delete(), get(), patch(), post(), put(),
			bearer(),
			status(),
			header(), ip(), ua(),
			redirect1(), redirectN(), redirectTo(),
		),
	)
}
