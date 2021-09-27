package aws

import (
	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian"
)

type µInput struct {
	events.APIGatewayProxyRequest
	Path []string
}

var (
	_ µ.Input = (*µInput)(nil)
)

//
// req.APIGatewayProxyRequest.QueryStringParameters)
// req.APIGatewayProxyRequest.Headers
