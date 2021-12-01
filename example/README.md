# µ-Gouldian examples 

The folders contains number of ready-to-run examples. Use the following command to execute example. 

```bash
cd name-of-example
go run name-of-example.go
```

Each of the following example becomes available at `http://localhost:8080`, use `curl` to evaluate behavior. 

Example of 
* Using GET, POST, PUT, PATCH and DELETE - [methods.go](methods/methods.go)
* Using path patterns and variable - [paths.go](paths/paths.go)
* Using query string parameters - [params.go](params/params.go) 
* Using headers values - [headers.go](headers/headers.go)
* Receiving and Sending payload - [payloads.go](payloads/payloads.go)

<!--
TODO
* Using serverless / lambdas
* Using authorization
* Using middleware (e.g. logging, open tracing)
* Using custom middleware (Endpoints)
* Using hard deadline (timeouts) for Endpoint
* Using multiple retry
* Using recovery
* Validate request input
* Stream in / out large files
* Sending byte stream from reader

-->

Misc Examples
* [Hello World!](helloworld/helloworld.go)
* [Echo](echo/echo.go)