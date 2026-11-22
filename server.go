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

package gouldian

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TrailingSlash is the policy of /path/ vs /path
type TrailingSlash uint8

const (
	// Strict (default) treats /path/ and /path as different resources
	Strict TrailingSlash = iota
	// Redirect responds to /path/ with 308 to /path
	Redirect
	// Ignore removes trailing slash before routing
	Ignore
)

// CORS configuration of the server. CORS is disabled by default.
type CORS struct {
	// Origins is allow-list of origins, "*" allows any origin (but never with Credentials)
	Origins []string
	// Methods allowed by preflight. All methods of the resource if empty.
	Methods []string
	// Headers allowed by preflight. It reflects requested headers if empty.
	Headers []string
	// Expose lists headers exposed to the browser
	Expose []string
	// MaxAge of preflight cache
	MaxAge time.Duration
	// Credentials allows cookies and authorization
	Credentials bool
}

// Server is the single place of rules that turns outcomes of endpoints into
// HTTP responses. It is shared by all adapters (httpd, apigateway, ...):
//
//	nil        204 No Content
//	*Output    as is
//	*Failure   400 (or the status of failure) Issue
//	NoMatch    404, 405 (Allow), 406, 415; OPTIONS and HEAD are derived
//	error      500, the cause is logged but never exposed
//	panic      500, recovered, the stack is logged
//
// The lifecycle of a request is
//
//	guards → routes → issue → log, CORS → observers
//
// the guards admit the request before routing (WithGuard), the issue of a
// failure is rendered by WithIssue or WithIssueHandler, the observers see the
// final response (WithObserver).
type Server struct {
	routes    *Routes
	log       *slog.Logger
	issue     IssueRenderer
	issueOut  IssueHandler
	guards    []Guard
	observers []Observer
	cors      *CORS
	bodyLimit int64
	recover   bool
	trailing  TrailingSlash
	logRoutes bool
	pool      sync.Pool
}

// Option of the server
type Option func(*Server)

// WithLogger sets the logger, slog.Default() otherwise. 5xx are logged at
// error level, 4xx at debug level. Only method and path are logged, never
// the query string.
func WithLogger(log *slog.Logger) Option { return func(s *Server) { s.log = log } }

// WithIssue sets the renderer of failures, the application owns the shape
// of problem. The default is RFC 9457 `application/problem+json`.
func WithIssue(f IssueRenderer) Option { return func(s *Server) { s.issue = f } }

// IssueHandler renders the failure into the output. Unlike IssueRenderer, it
// owns the whole response: it may change the status, set headers (e.g.
// Retry-After, WWW-Authenticate) and the body, and it sees the context of
// the request (RequestID, JWT, attributes). The issue is the default
// description of the failure, the cause is the error of the endpoint (nil
// for 404, 405, 406, 415).
type IssueHandler func(ctx *Context, issue Issue, cause error, out *Output)

// WithIssueHandler sets the handler of failures, it takes precedence over
// WithIssue. Every failure is rendered by it: the outcome of endpoints, the
// rejection of a guard and the responses of the router itself.
func WithIssueHandler(f IssueHandler) Option { return func(s *Server) { s.issueOut = f } }

// Guard admits the request before routing. A nil admits the request, any
// other outcome is the response, as if an endpoint returned it (e.g. a
// *Failure, an *Output, an error). The rejection is rendered through the
// issue of the server, the router is never consulted, so a client who is
// not admitted cannot tell known routes from unknown ones.
//
// Guards run for every request, including CORS preflight (OPTIONS): a guard
// that authorizes requests must admit preflight if CORS is enabled. Guards
// may set attributes of the context (Context.Set), they shall not bind lenses.
type Guard func(ctx *Context) error

// WithGuard appends guards to the server, they run in the given order.
func WithGuard(guards ...Guard) Option {
	return func(s *Server) { s.guards = append(s.guards, guards...) }
}

// Observer sees every request once it is served: the final response
// (rendered issue, CORS headers), the route that took it (Context.Route, nil
// if none did), the attributes of the context and the time of dispatch.
// The observer must not keep the context or the output, both are released
// once the response is sent.
type Observer func(ctx *Context, out *Output, elapsed time.Duration)

// WithObserver appends observers to the server (access logs, metrics,
// traces), they run in the given order.
func WithObserver(observers ...Observer) Option {
	return func(s *Server) { s.observers = append(s.observers, observers...) }
}

// WithCORS enables CORS, default is off.
func WithCORS(cors *CORS) Option { return func(s *Server) { s.cors = cors } }

// WithBodyLimit sets the limit of request body (413 if exceeded).
// The default is 6 MiB, 0 disables the limit.
func WithBodyLimit(n int64) Option { return func(s *Server) { s.bodyLimit = n } }

// WithRecover enables recovery from panics, it is on by default.
func WithRecover(on bool) Option { return func(s *Server) { s.recover = on } }

// WithTrailingSlash sets the policy of trailing slash, Strict by default.
func WithTrailingSlash(p TrailingSlash) Option { return func(s *Server) { s.trailing = p } }

// WithLogRoutes enables the log of configured routes (off by default). Mount
// logs, at info level, one record with all mounted routes as "METHOD /path",
// where ":" is a variable, "_" any segment and "*" the rest of the path.
//
//	level=INFO msg="routes mounted" count=2 routes="[GET /users/: POST /users]"
func WithLogRoutes(on bool) Option { return func(s *Server) { s.logRoutes = on } }

// NewServer creates the server
func NewServer(opts ...Option) *Server {
	s := &Server{
		routes:    &Routes{root: &node{}},
		issue:     DefaultIssue,
		bodyLimit: DefaultBodyLimit,
		recover:   true,
	}
	for _, opt := range opts {
		opt(s)
	}

	s.pool.New = func() any { return NewContext(context.Background()) }
	return s
}

// Mount adds routes to the server
func (s *Server) Mount(routes ...Routable) *Server {
	var mounted []string
	if s.logRoutes {
		mounted = make([]string, 0, len(routes))
	}

	for _, r := range routes {
		route := r()
		s.routes.append(route)
		if s.logRoutes {
			mounted = append(mounted, route.String())
		}
	}

	if s.logRoutes && len(mounted) != 0 {
		s.logger().LogAttrs(context.Background(), slog.LevelInfo, "routes mounted",
			slog.Int("count", len(mounted)),
			slog.Any("routes", mounted),
		)
	}
	return s
}

// NewContext creates (from the pool) the context of the request. The
// context of the caller (http.Request, Lambda invocation) is propagated to
// handlers. The context must be released by Release.
func (s *Server) NewContext(ctx context.Context, r *http.Request) *Context {
	c := s.pool.Get().(*Context)
	c.Context = ctx
	c.Request = r
	c.bodyLimit = s.bodyLimit
	c.log = s.log
	return c
}

// Release returns the context to the pool
func (s *Server) Release(c *Context) {
	c.Free()
	c.Context = context.Background()
	s.pool.Put(c)
}

func (s *Server) logger() *slog.Logger {
	if s.log != nil {
		return s.log
	}
	return slog.Default()
}

// Dispatch executes the request. It never returns nil, the output must be
// released with Output.Free once it is sent.
func (s *Server) Dispatch(c *Context) *Output {
	if len(s.observers) == 0 {
		out := s.safeHandle(c)
		s.finalize(c, out)
		return out
	}

	t := time.Now()
	out := s.safeHandle(c)
	s.finalize(c, out)
	elapsed := time.Since(t)
	for _, f := range s.observers {
		f(c, out, elapsed)
	}
	return out
}

func (s *Server) safeHandle(c *Context) (out *Output) {
	if s.recover {
		defer func() {
			if r := recover(); r != nil {
				out = c.output(http.StatusInternalServerError)
				out.Failure = fmt.Errorf("panic: %v", r)
				out.stack = string(debug.Stack())
				out.render = true
			}
		}()
	}

	return s.handle(c)
}

func (s *Server) handle(c *Context) *Output {
	for _, guard := range s.guards {
		if err := guard(c); err != nil {
			return s.output(c, err)
		}
	}

	req := c.Request
	path := req.URL.Path

	if len(path) > 1 && path[len(path)-1] == '/' {
		switch s.trailing {
		case Redirect:
			target := strings.TrimRight(path, "/")
			if target == "" {
				target = "/"
			}
			if req.URL.RawQuery != "" {
				target += "?" + req.URL.RawQuery
			}
			out := c.output(http.StatusPermanentRedirect)
			out.SetHeader("Location", target)
			return out
		case Ignore:
			path = strings.TrimRight(path, "/")
			if path == "" {
				path = "/"
			}
		}
	}

	return s.output(c, s.routes.exec(c, path))
}

func (s *Server) output(c *Context, err error) *Output {
	switch v := err.(type) {
	case nil:
		return c.output(http.StatusNoContent)
	case *Output:
		return v
	case *Failure:
		status := v.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		out := c.output(status)
		out.Failure = v
		out.render = true
		for key, vals := range v.Header {
			for _, val := range vals {
				out.AddHeader(key, val)
			}
		}
		return out
	case NoMatch:
		return s.noMatch(c, v)
	default:
		var o *Output
		if errors.As(err, &o) {
			return o
		}

		out := c.output(http.StatusInternalServerError)
		out.Failure = err
		out.render = true
		return out
	}
}

func (s *Server) noMatch(c *Context, reason NoMatch) *Output {
	var status int
	switch reason {
	case NoMatchMethod:
		allow := allowed(c.allow)
		if c.Request.Method == http.MethodOptions {
			out := c.output(http.StatusNoContent)
			out.SetHeader("Allow", strings.Join(allow, ", "))
			out.preflight = allow
			return out
		}
		out := c.output(http.StatusMethodNotAllowed)
		out.SetHeader("Allow", strings.Join(allow, ", "))
		out.reason = reason
		out.render = true
		return out
	case NoMatchAccept:
		status = http.StatusNotAcceptable
	case NoMatchContentType:
		status = http.StatusUnsupportedMediaType
	default:
		status = http.StatusNotFound
	}

	out := c.output(status)
	out.reason = reason
	out.render = true
	return out
}

// allowed builds sorted list of methods, HEAD is derived from GET and OPTIONS
// is always supported.
func allowed(verbs []string) []string {
	set := make([]string, 0, len(verbs)+2)
	for _, v := range verbs {
		if v == Any {
			continue
		}
		set = append(set, v)
		if v == http.MethodGet {
			set = append(set, http.MethodHead)
		}
	}
	set = append(set, http.MethodOptions)
	slices.Sort(set)
	return slices.Compact(set)
}

// finalize renders the issue, logs the outcome, applies CORS
func (s *Server) finalize(c *Context, out *Output) {
	req := c.Request

	if out.render && len(out.Body) == 0 {
		if s.issueOut != nil {
			s.issueOut(c, out.issue(req.URL.Path), out.Failure, out)
		} else {
			ct, body := s.issue(out.issue(req.URL.Path), out.Failure)
			if out.GetHeader("Content-Type") == "" {
				out.SetHeader("Content-Type", ct)
			}
			out.Body = body
		}
	}

	s.log_(c, out)
	s.applyCORS(c, out)

	if req.Method == http.MethodHead {
		out.Body = nil
	}
}

func (s *Server) log_(c *Context, out *Output) {
	if out.Failure == nil && out.reason == nil && out.Status < http.StatusInternalServerError {
		return
	}

	level := slog.LevelDebug
	if out.Status >= http.StatusInternalServerError {
		level = slog.LevelError
	}

	log := s.logger()
	if !log.Enabled(c, level) {
		return
	}

	attrs := make([]slog.Attr, 0, 7)
	attrs = append(attrs,
		slog.String("method", c.Request.Method),
		// never log the query string, it may carry tokens or personal data
		slog.String("path", c.Request.URL.Path),
		slog.Int("status", out.Status),
	)
	if c.RequestID != "" {
		attrs = append(attrs, slog.String("request_id", c.RequestID))
	}
	if out.reason != nil {
		attrs = append(attrs, slog.String("reason", out.reason.Error()))
	}
	if out.Failure != nil {
		attrs = append(attrs, slog.String("error", out.Failure.Error()))
	}
	if out.stack != "" {
		attrs = append(attrs, slog.String("stack", out.stack))
	}

	log.LogAttrs(c, level, "request failed", attrs...)
}

func (s *Server) applyCORS(c *Context, out *Output) {
	cors := s.cors
	if cors == nil {
		return
	}

	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		return
	}

	allowOrigin := ""
	for _, o := range cors.Origins {
		switch {
		case o == "*" && !cors.Credentials:
			allowOrigin = "*"
		case o == "*" || strings.EqualFold(o, origin):
			allowOrigin = origin
		}
		if allowOrigin != "" {
			break
		}
	}
	if allowOrigin == "" {
		return
	}

	// the headers of the endpoint win over defaults
	set := func(h, v string) {
		if v != "" && out.GetHeader(h) == "" {
			out.SetHeader(h, v)
		}
	}

	set("Access-Control-Allow-Origin", allowOrigin)
	if allowOrigin != "*" {
		set("Vary", "Origin")
	}
	if cors.Credentials {
		set("Access-Control-Allow-Credentials", "true")
	}
	set("Access-Control-Expose-Headers", strings.Join(cors.Expose, ", "))

	// automatic preflight
	if out.preflight != nil && c.Request.Header.Get("Access-Control-Request-Method") != "" {
		methods := cors.Methods
		if len(methods) == 0 {
			methods = out.preflight
		}
		set("Access-Control-Allow-Methods", strings.Join(methods, ", "))

		headers := strings.Join(cors.Headers, ", ")
		if headers == "" {
			headers = c.Request.Header.Get("Access-Control-Request-Headers")
		}
		set("Access-Control-Allow-Headers", headers)

		if cors.MaxAge > 0 {
			set("Access-Control-Max-Age", strconv.Itoa(int(cors.MaxAge.Seconds())))
		}
	}
}
