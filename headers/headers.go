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

package headers

/*

List of supported HTTP header constants, use them instead of explicit definition
*/
const (
	// Common HTTP headers
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
	CacheControl     = "Cache-Control"
	Connection       = "Connection"
	ContentEncoding  = "Content-Encoding"
	ContentLength    = "Content-Length"
	ContentType      = "Content-Type"
	Date             = "Date"
	TransferEncoding = "Transfer-Encoding"

	//
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
	Accept             = "Accept"
	AcceptCharset      = "Accept-Charset"
	AcceptEncoding     = "Accept-Encoding"
	AcceptLanguage     = "Accept-Language"
	Cookie             = "Cookie"
	Host               = "Host"
	IfMatch            = "If-Match"
	IfModifiedSince    = "If-Modified-Since"
	IfNoneMatch        = "If-None-Match"
	IfRange            = "If-Range"
	IfUnmodifiedSince  = "If-Unmodified-Since"
	Origin             = "Origin"
	ProxyAuthorization = "Proxy-Authorization"
	Range              = "Range"
	UserAgent          = "User-Agent"
	// Authorization      = "Authorization"

	//
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
	ContentLanguage = "Content-Language"
	ETag            = "ETag"
	Expires         = "Expires"
	LastModified    = "Last-Modified"
	Link            = "Link"
	Location        = "Location"
	Server          = "Server"
	SetCookie       = "Set-Cookie"
)

//
const (
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
	TextPlain       = "text/plain"
	TextHTML        = "text/html"
)
