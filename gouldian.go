package gouldian

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

func Serve(seq ...Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := JoinOr(seq...)

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var output *Output
		var issue Issue

		// TODO: handle headers
		http := NewRequest(req)
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
