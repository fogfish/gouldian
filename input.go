package gouldian

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"net/url"
)

type Input struct {
	events.APIGatewayProxyRequest
	segment int
	path    []string
	body    string
}

/*
func (x *Input) Json(val interface{}) {
	body, _ := json.Marshal(val)
	x.body = string(body)
}

func (x *Input) Text(text string) {
	x.body = text
}
*/

func New(req events.APIGatewayProxyRequest) *Input {
	return &Input{
		req,
		1,
		strings.Split(req.Path, "/"),
		"",
	}
}

func NewGet(spec string) *Input {
	uri, _ := url.Parse(spec)
	query := map[string]string{}
	for key, val := range uri.Query() {
		query[key] = strings.Join(val, "")
	}

	return &Input{
		events.APIGatewayProxyRequest{
			HTTPMethod:            "GET",
			Path:                  uri.Path,
			Headers:               map[string]string{},
			QueryStringParameters: query,
		},
		1,
		strings.Split(uri.Path, "/"),
		"",
	}
}

func (input *Input) With(head string, value string) *Input {
	input.Headers[head] = value
	return input
}

func (input *Input) WithJson(val interface{}) *Input {
	body, _ := json.Marshal(val)
	input.Body = string(body)
	return input
}

func (input *Input) WithText(val string) *Input {
	input.Body = val
	return input
}
