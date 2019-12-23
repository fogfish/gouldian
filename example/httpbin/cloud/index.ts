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

import * as api from '@aws-cdk/aws-apigateway'
import * as lambda from '@aws-cdk/aws-lambda'
import * as logs from '@aws-cdk/aws-logs'
import * as cdk from '@aws-cdk/core'
import * as pure from 'aws-cdk-pure'

const app = new cdk.App()
const stack = new cdk.Stack(app, 'httpbin', { })

const HttpBin = (): api.RestApiProps => ({
  deploy: true,
  deployOptions: {
    stageName: 'api',
  },
  endpointTypes: [api.EndpointType.REGIONAL],
  failOnWarnings: true,
})

const Lambda = (): lambda.FunctionProps => ({
  code: new lambda.AssetCode('../bin'),
  handler: 'main',
  logRetention: logs.RetentionDays.FIVE_DAYS,
  memorySize: 256,
  reservedConcurrentExecutions: 5,
  runtime: lambda.Runtime.GO_1_X,
  timeout: cdk.Duration.seconds(10),
})

pure.join(stack,
  pure.use({
    api: pure.iaac(api.RestApi)(HttpBin),
    httpbin: pure.wrap(api.LambdaIntegration)(pure.iaac(lambda.Function)(Lambda)),
  })
  .effect(x => {
    const proxy = x.api.root.addResource('{any+}')
    const methods = ['DELETE', 'GET', 'PATCH', 'POST', 'PUT']
    methods.map(m => proxy.addMethod(m, x.httpbin))
  })
)

app.synth()
