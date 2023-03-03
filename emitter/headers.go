package emitter

import (
	"strconv"
	"time"

	µ "github.com/fogfish/gouldian"
)

// Header defines HTTP headers to the request
//
//	ø.Header("User-Agent", "gurl")
func Header[T µ.ReadableHeaderValues](header string, value T) µ.Result {
	return HeaderOf[T](header).Set(value)
}

// Type of HTTP Header
//
//	const Host = HeaderOf[string]("Host")
//	ø.Host.Set("example.com")
type HeaderOf[T µ.ReadableHeaderValues] string

// Sets value of HTTP header
func (h HeaderOf[T]) Set(value T) µ.Result {
	switch v := any(value).(type) {
	case string:
		return func(out *µ.Output) error {
			out.SetHeader(string(h), v)
			return nil
		}
	case int:
		return func(out *µ.Output) error {
			out.SetHeader(string(h), strconv.Itoa(v))
			return nil
		}
	case time.Time:
		return func(out *µ.Output) error {
			out.SetHeader(string(h), v.UTC().Format(time.RFC1123))
			return nil
		}
	default:
		panic("invalid type")
	}
}

// Type of HTTP Header, Content-Type enumeration
//
//	const ContentType = HeaderEnumContent("Content-Type")
//	ø.ContentType.JSON
type HeaderEnumContent string

// Sets value of HTTP header
func (h HeaderEnumContent) Set(value string) µ.Result {
	return func(out *µ.Output) error {
		out.SetHeader(string(h), value)
		return nil
	}
}

// ApplicationJSON defines header `???: application/json`
func (h HeaderEnumContent) ApplicationJSON(out *µ.Output) error {
	out.SetHeader(string(h), "application/json")
	return nil
}

// JSON defines header `???: application/json`
func (h HeaderEnumContent) JSON(out *µ.Output) error {
	out.SetHeader(string(h), "application/json")
	return nil
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderEnumContent) Form(out *µ.Output) error {
	out.SetHeader(string(h), "application/x-www-form-urlencoded")
	return nil
}

// TextPlain defined Header `???: text/plain`
func (h HeaderEnumContent) TextPlain(out *µ.Output) error {
	out.SetHeader(string(h), "text/plain")
	return nil
}

// Text defined Header `???: text/plain`
func (h HeaderEnumContent) Text(out *µ.Output) error {
	out.SetHeader(string(h), "text/plain")
	return nil
}

// TextHTML defined Header `???: text/html`
func (h HeaderEnumContent) TextHTML(out *µ.Output) error {
	out.SetHeader(string(h), "text/html")
	return nil
}

// HTML defined Header `???: text/html`
func (h HeaderEnumContent) HTML(out *µ.Output) error {
	out.SetHeader(string(h), "text/html")
	return nil
}

// Type of HTTP Header, Connection enumeration
//
//	const Connection = HeaderEnumConnection("Connection")
//	ø.Connection.KeepAlive
type HeaderEnumConnection string

// Sets value of HTTP header
func (h HeaderEnumConnection) Set(value string) µ.Result {
	return func(out *µ.Output) error {
		out.SetHeader(string(h), value)
		return nil
	}
}

// KeepAlive defines header `???: keep-alive`
func (h HeaderEnumConnection) KeepAlive(out *µ.Output) error {
	out.SetHeader(string(h), "keep-alive")
	return nil
}

// Close defines header `???: close`
func (h HeaderEnumConnection) Close(out *µ.Output) error {
	out.SetHeader(string(h), "close")
	return nil
}

// Type of HTTP Header, Transfer-Encoding enumeration
//
//	const TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
//	ø.TransferEncoding.Chunked
type HeaderEnumTransferEncoding string

// Sets value of HTTP header
func (h HeaderEnumTransferEncoding) Set(value string) µ.Result {
	return func(out *µ.Output) error {
		out.SetHeader(string(h), value)
		return nil
	}
}

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderEnumTransferEncoding) Chunked(out *µ.Output) error {
	out.SetHeader(string(h), "chunked")
	return nil
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderEnumTransferEncoding) Identity(out *µ.Output) error {
	out.SetHeader(string(h), "identity")
	return nil
}

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
const (
	Age              = HeaderOf[int]("Age")
	CacheControl     = HeaderOf[string]("Cache-Control")
	Connection       = HeaderEnumConnection("Connection")
	ContentEncoding  = HeaderOf[string]("Content-Encoding")
	ContentLanguage  = HeaderOf[string]("Content-Language")
	ContentLength    = HeaderOf[int]("Content-Length")
	ContentLocation  = HeaderOf[string]("Content-Location")
	ContentMD5       = HeaderOf[string]("Content-MD5")
	ContentRange     = HeaderOf[string]("Content-Range")
	ContentType      = HeaderEnumContent("Content-Type")
	Date             = HeaderOf[time.Time]("Date")
	ETag             = HeaderOf[string]("ETag")
	Expires          = HeaderOf[time.Time]("Expires")
	LastModified     = HeaderOf[time.Time]("Last-Modified")
	Link             = HeaderOf[string]("Link")
	Location         = HeaderOf[string]("Location")
	RetryAfter       = HeaderOf[time.Time]("Retry-After")
	Server           = HeaderOf[string]("Server")
	SetCookie        = HeaderOf[string]("Set-Cookie")
	TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
	Via              = HeaderOf[string]("Via")
)
