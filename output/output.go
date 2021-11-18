package output

import (
	"encoding/json"
	"fmt"
	µ "github.com/fogfish/gouldian"
	"net/http"
)

// JSON appends application/json payload to HTTP response
func JSON(val interface{}) µ.Result {
	return func(out *µ.Output) error {
		body, err := json.Marshal(val)
		if err != nil {
			out.Status = http.StatusInternalServerError
			out.Headers["Content-Type"] = "text/plain"
			out.Body = fmt.Sprintf("JSON serialization is failed for <%T>", val)

			return nil
		}

		out.Headers["Content-Type"] = "application/json"
		out.Body = string(body)
		return nil
	}
}

// Bytes appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func Bytes(content []byte) µ.Result {
	return func(out *µ.Output) error {
		out.Body = string(content)
		return nil
	}
}

// Text appends arbitrary octet/stream payload to HTTP response
// content type shall be specified using With method
func Text(content string) µ.Result {
	return func(out *µ.Output) error {
		out.Body = content
		return nil
	}
}

// Issue ...
func Issue(err error, title ...string) µ.Result {
	return func(out *µ.Output) error {
		issue := µ.NewIssue(out.Status)
		if len(title) != 0 {
			issue.Title = title[0]
		}

		body, err := json.Marshal(issue)
		if err != nil {
			out.Status = http.StatusInternalServerError
			out.Headers["Content-Type"] = "text/plain"
			out.Body = fmt.Sprintf("JSON serialization is failed for <Issue>")

			return nil
		}

		out.Headers["Content-Type"] = "application/json"
		out.Body = string(body)
		out.Failure = err
		return nil
	}
}

//
type Header string

//
func (header Header) Is(value string) µ.Result {
	return func(out *µ.Output) error {
		out.Headers[string(header)] = value
		return nil
	}
}

// Content defines headers for content negotiation
type Content Header

// JSON matches header `???: application/json`
func (h Content) JSON(out *µ.Output) error {
	out.Headers[string(h)] = "application/json"
	return nil
}

// Form matches Header `???: application/x-www-form-urlencoded`
func (h Content) Form(out *µ.Output) error {
	out.Headers[string(h)] = "application/x-www-form-urlencoded"
	return nil
}

// Text matches Header `???: text/plain`
func (h Content) Text(out *µ.Output) error {
	out.Headers[string(h)] = "text/plain"
	return nil
}

// HTML matches Header `???: text/html`
func (h Content) HTML(out *µ.Output) error {
	out.Headers[string(h)] = "text/html"
	return nil
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (h Content) Is(value string) µ.Result {
	return Header(h).Is(value)
}

/*

List of supported HTTP header constants
https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
*/
const (
	CacheControl     = Header("Cache-Control")
	Connection       = Header("Connection")
	ContentEncoding  = Header("Content-Encoding")
	ContentLanguage  = Header("Content-Language")
	ContentLength    = Header("Content-Length")
	ContentType      = Content("Content-Type")
	Date             = Header("Date")
	ETag             = Header("ETag")
	Expires          = Header("Expires")
	LastModified     = Header("Last-Modified")
	Link             = Header("Link")
	Location         = Header("Location")
	Server           = Header("Server")
	SetCookie        = Header("Set-Cookie")
	TransferEncoding = Header("Transfer-Encoding")
)
