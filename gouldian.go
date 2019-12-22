package gouldian

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

//
type Pattern interface {
	//
	Path(segment string) Pattern
	String(val *string) Pattern
	Int(val *int) Pattern

	//
	Opt(name string, val string) Pattern
	IsOpt(name string) Pattern
	// TODO: MayBe
	OptString(name string, val *string) Pattern
	OptInt(name string, val *int) Pattern

	//
	Head(name string, val string) Pattern
	IsHead(name string) Pattern
	// MayBe
	HeadString(name string, val *string) Pattern

	//
	Json(val interface{}) Pattern
	Text(val *string) Pattern

	//
	Then(f func() error) Endpoint
	IsMatch(in *Input) bool
}

//
/*
type Output interface {
	Json(json interface{})
	Text(text string)
}
*/

func Serve(endpoints ...Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := endpoints[0]
	for _, x := range endpoints[1:] {
		api = api.Or(x)
	}

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var output *Output
		var issue Issue

		// TODO: handle headers
		http := New(req)
		err := api(http)
		if errors.As(err, &output) {
			return events.APIGatewayProxyResponse{Body: output.body, StatusCode: output.status}, nil
		} else if errors.As(err, &issue) {
			text, _ := json.Marshal(issue)
			return events.APIGatewayProxyResponse{Body: string(text), StatusCode: issue.Status}, nil
		} else if errors.Is(err, NoMatch{}) {
			nm, _ := json.Marshal(issue)
			return events.APIGatewayProxyResponse{Body: string(nm), StatusCode: 501}, nil
		}
		fatal, _ := json.Marshal(issue)
		return events.APIGatewayProxyResponse{Body: string(fatal), StatusCode: 500}, nil
	}
}
