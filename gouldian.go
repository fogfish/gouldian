//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

// Serve HTTP service
func Serve(seq ...Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := JoinOr(seq...)

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var output *Output
		var issue *Issue

		http := NewRequest(req)
		err := api(http)
		if errors.As(err, &output) {
			return events.APIGatewayProxyResponse{
				Body:       output.body,
				StatusCode: output.status,
				Headers:    output.headers,
			}, nil
		} else if errors.As(err, &issue) {
			text, _ := json.Marshal(issue)
			return events.APIGatewayProxyResponse{Body: string(text), StatusCode: issue.Status}, nil
		} else if errors.Is(err, NoMatch{}) {
			return events.APIGatewayProxyResponse{StatusCode: 501}, nil
		}
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
}
