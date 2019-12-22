package main

import (
	"fmt"
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

func verb(method gouldian.Pattern, path string) gouldian.Endpoint {
	h := headers{}

	return method.Path(path).
		HeadString("Accept", &h.Accept).
		HeadString("Host", &h.Host).
		HeadString("Origin", &h.Origin).
		HeadString("Referer", &h.Referer).
		HeadString("User-Agent", &h.UserAgent).
		Then(func() error {
			return gouldian.Ok().
				Json(response{h, fmt.Sprintf("https://%v/%v", h.Host, path)})
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
		HeadString("Authorization", &token).
		Then(func() error {
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
		Then(func() error {
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
		HeadString("Accept", &h.Accept).
		HeadString("Host", &h.Host).
		HeadString("Origin", &h.Origin).
		HeadString("Referer", &h.Referer).
		HeadString("User-Agent", &h.UserAgent).
		Then(func() error {
			return gouldian.Ok().Json(h)
		})
}

func ip() gouldian.Endpoint {
	var ip string
	return gouldian.Get().Path("ip").
		HeadString("X-Forwarded-For", &ip).
		Then(func() error {
			return gouldian.Ok().Json(ip)
		})
}

func ua() gouldian.Endpoint {
	var ua string
	return gouldian.Get().Path("user-agent").
		HeadString("User-Agent", &ua).
		Then(func() error {
			return gouldian.Ok().Json(ip)
		})
}

//-----------------------------------------------------------------------------
//
// Redirects
//
//-----------------------------------------------------------------------------

func redirect1() gouldian.Endpoint {
	return gouldian.Get().Path("redirect").Path("1").
		Then(func() error {
			return gouldian.Success(302).
				With("Location", "https://example.com")
		})
}

func redirectN() gouldian.Endpoint {
	var host string
	var n int
	return gouldian.Get().Path("redirect").Int(&n).
		HeadString("Host", &host).
		Then(func() error {
			return gouldian.Success(302).
				With("Location", fmt.Sprintf("https://%v/redirect/%v", host, n-1))
		})
}

func redirectTo() gouldian.Endpoint {
	var url string
	return gouldian.Get().Path("redirect-to").
		OptString("url", &url).
		Then(func() error {
			return gouldian.Success(302).With("Location", url)
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
