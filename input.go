package gouldian

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"net/url"
)

// Input wraps HTTP request
type Input struct {
	events.APIGatewayProxyRequest
	segment int
	path    []string
	body    string
}

// Mock creates new Input - HTTP GET request
func Mock(httpURL string) *Input {
	return MockVerb("GET", httpURL)
}

// MockVerb creates new Input with any verb
func MockVerb(verb string, httpURL string) *Input {
	uri, _ := url.Parse(httpURL)
	query := map[string]string{}
	for key, val := range uri.Query() {
		query[key] = strings.Join(val, "")
	}

	return NewRequest(
		events.APIGatewayProxyRequest{
			HTTPMethod:            verb,
			Path:                  uri.Path,
			Headers:               map[string]string{},
			QueryStringParameters: query,
		},
	)
}

// NewRequest creates new Input from API Gateway request
func NewRequest(req events.APIGatewayProxyRequest) *Input {
	return &Input{req, 1, strings.Split(req.Path, "/"), ""}
}

// With adds HTTP header to mocked request
func (input *Input) With(head string, value string) *Input {
	input.Headers[head] = value
	return input
}

// WithJSON adds Json payload to mocked request
func (input *Input) WithJSON(val interface{}) *Input {
	body, _ := json.Marshal(val)
	input.Body = string(body)
	return input
}

// WithText adds Text payload to mocked requets
func (input *Input) WithText(val string) *Input {
	input.Body = val
	return input
}
