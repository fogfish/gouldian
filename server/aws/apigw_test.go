package aws_test

/*
func TestServeSuccess(t *testing.T) {
	fun := µ.Serve(hello())
	req := mock.Input(mock.URL("/hello"))
	rsp, _ := fun(req.APIGatewayProxyRequest)

	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "text/plain",
	}

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(200).
		If(rsp.Headers).Should().Equal(head).
		If(rsp.Body).Should().Equal("Hello World!")
}

func hello() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("hello")),
		µ.FMap(
			func() error { return µ.Ok().Text("Hello World!") },
		),
	)
}

func TestServeFailure(t *testing.T) {
	fun := µ.Serve(unauthorized())
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/issue",
	}
	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}
	rsp, _ := fun(req)
	var issue µ.Issue

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(401).
		If(rsp.Headers).Should().Equal(head).
		If(json.Unmarshal([]byte(rsp.Body), &issue)).Should().Equal(nil).
		If(issue.Type).Should().Equal("https://httpstatuses.com/401").
		If(issue.Status).Should().Equal(401).
		If(issue.Title).Should().Equal("Unauthorized").
		If(issue.ID).ShouldNot().Equal("")
}

func unauthorized() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("issue")),
		µ.FMap(func() error { return µ.Unauthorized(errors.New("some reason")) }),
	)
}

func TestServeNoMatch(t *testing.T) {
	fun := µ.Serve(hello())
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/issue",
	}
	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}
	rsp, _ := fun(req)

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(501).
		If(rsp.Headers).Should().Equal(head)
}

func TestServeNoMatchLogger(t *testing.T) {
	fun := µ.Serve(µ.NoMatchLogger())
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/issue",
	}
	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}
	rsp, _ := fun(req)

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(501).
		If(rsp.Headers).Should().Equal(head)
}

func TestServeUnescapedPath(t *testing.T) {
	fun := µ.Serve(unescaped())
	req := mock.Input(mock.URL("/"))
	req.APIGatewayProxyRequest.Path = "/h%rt"
	req.Path = []string{"h%rt"}
	rsp, _ := fun(req.APIGatewayProxyRequest)

	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "text/plain",
	}

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(200).
		If(rsp.Headers).Should().Equal(head).
		If(rsp.Body).Should().Equal("Hello World!")
}

func unescaped() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("h%rt")),
		µ.FMap(
			func() error { return µ.Ok().Text("Hello World!") },
		),
	)
}

func TestServeUnknownError(t *testing.T) {
	fun := µ.Serve(unknown())
	req := mock.Input(mock.URL("/"))
	req.APIGatewayProxyRequest.Path = "/h%rt"
	req.Path = []string{"h%rt"}
	rsp, _ := fun(req.APIGatewayProxyRequest)

	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(500).
		If(rsp.Headers).Should().Equal(head)
}

func unknown() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("h%rt")),
		µ.FMap(
			func() error { return fmt.Errorf("Unknown error") },
		),
	)
}
*/
