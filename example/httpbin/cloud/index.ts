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
