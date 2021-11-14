package headers

import (
	"fmt"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
	"path/filepath"
	"strings"
)

/*

List of supported HTTP header constants, use them instead of explicit definition
https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
*/
const (
	Accept             = µ.Header("Accept")
	AcceptCharset      = µ.Header("Accept-Charset")
	AcceptEncoding     = µ.Header("Accept-Encoding")
	AcceptLanguage     = µ.Header("Accept-Language")
	Authorization      = Authorize("Authorization")
	CacheControl       = µ.Header("Cache-Control")
	Connection         = µ.Header("Connection")
	ContentEncoding    = µ.Header("Content-Encoding")
	ContentLength      = µ.Header("Content-Length")
	ContentType        = Content("Content-Type")
	Cookie             = µ.Header("Cookie")
	Date               = µ.Header("Date")
	Host               = µ.Header("Host")
	IfMatch            = µ.Header("If-Match")
	IfModifiedSince    = µ.Header("If-Modified-Since")
	IfNoneMatch        = µ.Header("If-None-Match")
	IfRange            = µ.Header("If-Range")
	IfUnmodifiedSince  = µ.Header("If-Unmodified-Since")
	Origin             = µ.Header("Origin")
	ProxyAuthorization = µ.Header("Proxy-Authorization")
	Range              = µ.Header("Range")
	TransferEncoding   = µ.Header("Transfer-Encoding")
	UserAgent          = µ.Header("User-Agent")
)

/*

Content is "synonym" to Header type. It defines a few Endpoints that simplify
content negotiation use-cases.

  e := µ.GET( µ.ContentType.JSON )
  e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
  e(mock.Input(mock.Header("Content-Type", "text/plain"))) != nil

*/
type Content µ.Header

func (h Content) eval(req µ.Input, val string) error {
	opt, exists := req.Headers().Get(string(h))
	if exists && strings.HasPrefix(opt, val) {
		return nil
	}
	return µ.NoMatch{}
}

// JSON is a syntax sugar to Header(???).Is("application/json")
func (h Content) JSON(req µ.Input) error {
	return h.eval(req, "application/json")
}

// Form is a syntax sugar to Header(???).Is("application/x-www-form-urlencoded")
func (h Content) Form(req µ.Input) error {
	return h.eval(req, "application/x-www-form-urlencoded")
}

// Text is a syntax sugar to Header(???).Is("text/plain")
func (h Content) Text(req µ.Input) error {
	return h.eval(req, "text/plain")
}

// HTML is a syntax sugar to Header(???).Is("text/html")
func (h Content) HTML(req µ.Input) error {
	return h.eval(req, "text/html")
}

// Is implements matcher for Content type (see Header.Is)
func (h Content) Is(value string) µ.Endpoint {
	return µ.Header(h).Is(value)
}

// Any implements matcher for Content type (see Header.Any)
func (h Content) Any(req µ.Input) error {
	return µ.Header(h).Any(req)
}

// To implements matcher for Content type (see Header.To)
func (h Content) To(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).To(lens)
}

// Maybe implements matcher for Content type (see Header.Maybe)
func (h Content) Maybe(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).Maybe(lens)
}

/*

Authorize is "synonym" to Header type. It defines a few Endpoints that simplify
validation of credentials/tokens supplied within the request

  e := µ.GET( µ.Authorization.With(func(string, string) error { ... }) )
  e(mock.Input(mock.Header("Authorization", "Basic foo"))) == nil
  e(mock.Input(mock.Header("Authorization", "Basic bar"))) != nil

*/
type Authorize µ.Header

// To implements matcher for Content type (see Header.To)
func (h Authorize) To(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).To(lens)
}

// Maybe implements matcher for Content type (see Header.Maybe)
func (h Authorize) Maybe(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).Maybe(lens)
}

// With validates content of HTTP Authorization header
func (h Authorize) With(f func(string, string) error) µ.Endpoint {
	return func(req µ.Input) error {
		auth, exists := req.Headers().Get("Authorization")
		if !exists {
			return µ.Status.Unauthorized(
				fmt.Errorf("Unauthorized %v", filepath.Join(req.Resource()...)),
			)
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			return µ.Status.Unauthorized(
				fmt.Errorf("Unauthorized %v", filepath.Join(req.Resource()...)),
			)
		}

		if err := f(cred[0], cred[1]); err != nil {
			return µ.Status.Unauthorized(err)
		}

		return nil
	}
}
