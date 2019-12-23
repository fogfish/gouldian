# Hello World

Run Hello World example locally

```bash
GOOS=linux GOARCH=amd64 go build main.go

docker run --rm \
  -v `pwd`:/var/task:ro,delegated \
  lambci/lambda:go1.x main \
  '{ "httpMethod": "GET", "path": "/hello" }'
```