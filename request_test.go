package gouldian_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	µ "github.com/fogfish/gouldian/v2"
	"github.com/fogfish/gouldian/v2/mock"
	"github.com/fogfish/it/v2"
)

func TestHTTP(t *testing.T) {
	foo := mock.Endpoint(
		µ.HTTP(
			http.MethodGet,
			µ.URI(µ.Path("foo")),
		),
	)

	req := mock.Input(mock.URL("/foo"))
	err := foo(req)
	it.Then(t).Should(
		it.Nil(err),
	)
}

func TestMethod(t *testing.T) {
	spec := []struct {
		Verb func(µ.Routable, ...µ.Endpoint) µ.Routable
		Mock mock.Mock
	}{
		{µ.GET, mock.Method("GET")},
		{µ.PUT, mock.Method("PUT")},
		{µ.POST, mock.Method("POST")},
		{µ.DELETE, mock.Method("DELETE")},
		{µ.PATCH, mock.Method("PATCH")},
		{µ.ANY, mock.Method("GET")},
		{µ.ANY, mock.Method("PUT")},
	}

	for _, tt := range spec {
		foo := mock.Endpoint(
			tt.Verb(
				µ.URI(µ.Path("foo")),
			),
		)

		req := mock.Input(tt.Mock, mock.URL("/foo"))
		err := foo(req)
		it.Then(t).Should(
			it.Nil(err),
		)
	}
}

func TestMethodNoMatch(t *testing.T) {
	spec := []struct {
		Verb func(µ.Routable, ...µ.Endpoint) µ.Routable
		Mock mock.Mock
	}{
		{µ.GET, mock.Method("OTHER")},
		{µ.PUT, mock.Method("OTHER")},
		{µ.POST, mock.Method("OTHER")},
		{µ.DELETE, mock.Method("OTHER")},
		{µ.PATCH, mock.Method("OTHER")},
	}

	for _, tt := range spec {
		foo := mock.Endpoint(
			tt.Verb(
				µ.URI(µ.Path("foo")),
			),
		)

		req := mock.Input(tt.Mock, mock.URL("/foo"))
		err := foo(req)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	}
}

func TestBodyJson(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock   *µ.Context
		Expect foobar
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.JSON(foobar{"foo1", 10}),
			),
			foobar{"foo1", 10},
		},
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.Text(`{"foo":"foo1","bar":10}`),
			),
			foobar{"foo1", 10},
		},
	}

	type request struct {
		FooBar foobar `content:"json"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}

func TestBodyJsonNoMatch(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock *µ.Context
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.Text(`{"foo:"foo1,"bar":10}`),
			),
		},
	}

	type request struct {
		FooBar foobar `content:"json"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).
			Should(it.Nil(err)).
			ShouldNot(it.Nil(µ.FromContext(tt.Mock, &req)))
	}
}

func TestBodyForm(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock   *µ.Context
		Expect foobar
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/x-www-form-urlencoded"),
				mock.Text("foo=foo1&bar=10"),
			),
			foobar{"foo1", 10},
		},
	}

	type request struct {
		FooBar foobar `content:"form"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}

func TestBodyFormNoMatch(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock *µ.Context
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/x-www-form-urlencoded"),
				mock.Text("foobar"),
			),
		},
	}

	type request struct {
		FooBar foobar `content:"form"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).
			Should(it.Nil(err)).
			ShouldNot(it.Nil(µ.FromContext(tt.Mock, &req)))
	}
}

func TestBodyText(t *testing.T) {
	spec := []struct {
		Mock   *µ.Context
		Expect string
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "text/plain"),
				mock.Text("foobar"),
			),
			"foobar",
		},
	}

	type request struct {
		FooBar string
	}
	var lens = µ.Optics1[request, string]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}

func TestBodyOctetStream(t *testing.T) {
	spec := []struct {
		Mock   *µ.Context
		Expect string
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/octet-stream"),
				mock.Text("foobar"),
			),
			"foobar",
		},
	}

	type request struct {
		FooBar string
	}
	var lens = µ.Optics1[request, string]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}

}

func TestFMapSuccess(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.FMap(func(ctx *µ.Context, t *T) error {
				out := µ.NewOutput(http.StatusOK)
				out.Body = t.A
				return out
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckOutput(err, "bar"),
	)
}

func TestFMapFailure(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.FMap(func(*µ.Context, *T) error {
				out := µ.NewOutput(http.StatusUnauthorized)
				out.SetIssue(errors.New(""))
				return out
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckStatusCode(err, http.StatusUnauthorized),
	)
}

func TestMapSuccess(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.Map(func(ctx *µ.Context, t *T) (*T, error) { return t, nil }),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckOutput(err, `{"A":"bar"}`),
	)
}

func TestContextFree(t *testing.T) {
	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("foo"))))
	req := mock.Input(mock.URL("/foo"))
	err := foo(req)

	it.Then(t).Should(it.Nil(err))

	req.Free()
	err = foo(req)

	it.Then(t).ShouldNot(it.Nil(err))
}

func TestOutputFree(t *testing.T) {
	out := µ.NewOutput(200)
	out.SetHeader("X-Foo", "bar")
	out.Body = "test"

	it.Then(t).Should(
		it.Equal(out.Status, 200),
		it.Equal(out.GetHeader("X-Foo"), "bar"),
		it.Equal(out.Body, "test"),
	)

	out.Free()

	it.Then(t).Should(
		it.Equal(out.GetHeader("X-Foo"), ""),
		it.Equal(out.Body, ""),
	)
}

func TestHandlerSuccess(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			mock.Output(http.StatusOK, "bar"),
		),
	)
	req := mock.Input(mock.URL("/foo"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckOutput(err, "bar"),
	)
}

func TestHandler2Success(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			mock.Output(http.StatusOK, "bar"),
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("bar")),
			mock.Output(http.StatusOK, "foo"),
		),
	)
	req := mock.Input(mock.URL("/foo"))
	err := µ.Endpoints{foo, bar}.Or(req)

	it.Then(t).Should(
		mock.CheckOutput(err, "bar"),
	)
}

func TestHandlerFailure(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			func(*µ.Context) error {
				out := µ.NewOutput(http.StatusUnauthorized)
				out.SetIssue(errors.New(""))
				return out
			},
		),
	)
	req := mock.Input(mock.URL("/foo"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckStatusCode(err, http.StatusUnauthorized),
	)
}

func TestMapFailure(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.Map(func(*µ.Context, *T) (*T, error) {
				out := µ.NewOutput(http.StatusUnauthorized)
				out.SetIssue(errors.New(""))
				return nil, out
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))
	err := foo(req)

	it.Then(t).Should(
		mock.CheckStatusCode(err, http.StatusUnauthorized),
	)
}

func TestBodyLeak(t *testing.T) {
	type Pair struct {
		Key int    `json:"key,omitempty"`
		Val string `json:"val,omitempty"`
	}
	type Item struct {
		Seq []Pair `json:"seq,omitempty"`
	}
	type request struct {
		Item Item
	}
	lens := µ.Optics1[request, Item]()

	endpoint := func() µ.Routable {
		return µ.GET(
			µ.URI(),
			µ.Body(lens),
			func(ctx *µ.Context) error {
				var req request
				if err := µ.FromContext(ctx, &req); err != nil {
					return err
				}

				seq := []Pair{}
				for key, val := range req.Item.Seq {
					if val.Key == 0 {
						seq = append(seq, Pair{Key: key + 1, Val: val.Val})
					}
				}
				req.Item = Item{Seq: seq}
				out := µ.NewOutput(http.StatusOK)

				val, _ := json.Marshal(req.Item)
				out.Body = string(val)
				return out
			},
		)
	}

	foo := mock.Endpoint(endpoint())
	for val, expect := range map[string]string{
		"{\"seq\":[{\"val\":\"a\"},{\"val\":\"b\"}]}":                 "{\"seq\":[{\"key\":1,\"val\":\"a\"},{\"key\":2,\"val\":\"b\"}]}",
		"{\"seq\":[{\"val\":\"c\"}]}":                                 "{\"seq\":[{\"key\":1,\"val\":\"c\"}]}",
		"{\"seq\":[{\"val\":\"d\"},{\"val\":\"e\"},{\"val\":\"f\"}]}": "{\"seq\":[{\"key\":1,\"val\":\"d\"},{\"key\":2,\"val\":\"e\"},{\"key\":3,\"val\":\"f\"}]}",
	} {
		req := mock.Input(
			mock.Method("GET"),
			mock.Header("Content-Type", "application/json"),
			mock.Text(val),
		)
		out := foo(req)

		it.Then(t).Should(
			it.Equal(out.Error(), expect),
		)
	}
}
