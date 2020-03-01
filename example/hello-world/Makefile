all: main 
	@docker run --rm \
  	  -v `pwd`:/var/task:ro,delegated \
  	  lambci/lambda:go1.x main \
  	  '{ "httpMethod": "GET", "path": "/hello" }'

main: main.go
	@GOOS=linux GOARCH=amd64 go build main.go

