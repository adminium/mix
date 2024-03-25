package mix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

type ITestHandler interface {
	Ping(ctx context.Context, msg string) (reply string, err error)
	Query(ctx context.Context, page, limit int) (reply string, err error)
	Download(ctx context.Context, file string) io.Reader
	Error(ctx context.Context) error
	Code(ctx context.Context) error
	Warn(ctx context.Context) error
	Upload(ctx context.Context, file string, size int64, data []byte) (err error)
}

var _ ITestHandler = (*TestHandler)(nil)

type TestHandler struct {
}

func (t TestHandler) Upload(ctx context.Context, file string, size int64, data []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t TestHandler) Error(ctx context.Context) error {
	return fmt.Errorf("default error")
}

func (t TestHandler) Code(ctx context.Context) error {
	return Codef(1000, "error with code")
}

func (t TestHandler) Warn(ctx context.Context) error {
	return Warnf("warning")
}

func (t TestHandler) Query(ctx context.Context, page, limit int) (reply string, err error) {
	reply = fmt.Sprintf("page:%d limit:%d", page, limit)
	return
}

func (t TestHandler) Ping(ctx context.Context, msg string) (reply string, err error) {
	reply = fmt.Sprintf("received: %s", msg)
	return
}

func (t TestHandler) Download(ctx context.Context, file string) io.Reader {
	buf := &bytes.Buffer{}
	buf.WriteString("<h1>Hello world</h1>")
	return buf
}

func TestServer(t *testing.T) {

	h := &TestHandler{}

	s := NewServer()
	group := s.Group("/api/v1")

	s.RegisterRPC(s.Group("/rpc/v1"), "test", h)
	s.RegisterAPI(group, "", h)

	group.GET("/download", WrapHandler(func(ctx *gin.Context) (data any, err error) {
		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		return h.Download(ctx, "ok"), nil
	}))

	const port = ":10000"

	go func() {
		require.NoError(t, s.Run(port))
	}()

	c := &TestClient{
		addr:   fmt.Sprintf("http://127.0.0.1%s", port),
		client: resty.New(),
	}
	time.Sleep(50 * time.Millisecond)
	c.TestPing(t)
	c.TestError(t)
	c.TestCode(t)
	c.TestWarn(t)
	c.TestQuery(t)
}

type TestCase struct {
	requests []*resty.Request
	handle   func(t *testing.T, index int, resp *resty.Response)
}

type TestClient struct {
	addr   string
	client *resty.Client
}

func (c *TestClient) apiUrl(path ...string) string {
	return fmt.Sprintf("%s/api/v1%s", c.addr, strings.Join(path, "/"))
}

func (c *TestClient) rpcUrl() string {
	return fmt.Sprintf("%s/rpc/v1", c.addr)
}

func (c *TestClient) newRequest(url string) *resty.Request {
	r := c.client.R()
	r.URL = url
	r.Method = resty.MethodPost
	return r
}

func (c *TestClient) executeTestCases(t *testing.T, cases []*TestCase) {
	for _, v := range cases {
		for index, vv := range v.requests {
			r, err := vv.Send()
			require.NoError(t, err, vv.URL)
			v.handle(t, index, r)
		}
	}
}

func (c *TestClient) wrapRPCBody(method string, params ...any) any {
	return map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      time.Now().Nanosecond(),
	}
}

func (c *TestClient) TestError(t *testing.T) {
	c.executeTestCases(t, []*TestCase{
		{
			requests: []*resty.Request{
				c.newRequest(c.apiUrl("/Error")),
				c.newRequest(c.rpcUrl()).SetBody(c.wrapRPCBody("test.Error")),
			},
			handle: func(t *testing.T, index int, resp *resty.Response) {
				require.True(t, resp.IsError())
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
				status := &Error{}
				require.NoError(t, json.Unmarshal(resp.Body(), status))
				require.Equal(t, 1, int(status.Code))
				require.Equal(t, "default error", status.Message)
			},
		},
	})
}

func (c *TestClient) TestCode(t *testing.T) {
	c.executeTestCases(t, []*TestCase{
		{
			requests: []*resty.Request{
				c.newRequest(c.apiUrl("/Code")),
				c.newRequest(c.rpcUrl()).SetBody(c.wrapRPCBody("test.Code")),
			},
			handle: func(t *testing.T, index int, resp *resty.Response) {
				require.True(t, resp.IsError())
				require.Equal(t, http.StatusBadRequest, resp.StatusCode())
				status := &Error{}
				require.NoError(t, json.Unmarshal(resp.Body(), status))
				require.Equal(t, 1000, int(status.Code))
				require.Equal(t, "error with code", status.Message)
			},
		},
	})
}

func (c *TestClient) TestWarn(t *testing.T) {
	c.executeTestCases(t, []*TestCase{
		{
			requests: []*resty.Request{
				c.newRequest(c.apiUrl("/Warn")),
				c.newRequest(c.rpcUrl()).SetBody(c.wrapRPCBody("test.Warn")),
			},
			handle: func(t *testing.T, index int, resp *resty.Response) {
				require.True(t, resp.IsError())
				require.Equal(t, http.StatusBadRequest, resp.StatusCode())
				status := &Error{}
				require.NoError(t, json.Unmarshal(resp.Body(), status))
				require.Equal(t, 2, int(status.Code))
				require.Equal(t, "warning", status.Message)
			},
		},
	})
}

func (c *TestClient) TestPing(t *testing.T) {
	c.executeTestCases(t, []*TestCase{
		{
			requests: []*resty.Request{
				c.newRequest(c.apiUrl("/Ping")).SetBody([]string{"hello"}),
				c.newRequest(c.rpcUrl()).SetBody(c.wrapRPCBody("test.Ping", "hello")),
			},
			handle: func(t *testing.T, index int, resp *resty.Response) {
				require.True(t, !resp.IsError(), index)
				type Data struct {
					Result string
				}
				d := &Data{}
				require.NoError(t, json.Unmarshal(resp.Body(), d), index)
				require.Equal(t, "received: hello", d.Result, index)
			},
		},
	})
}

func (c *TestClient) TestQuery(t *testing.T) {
	c.executeTestCases(t, []*TestCase{
		{
			requests: []*resty.Request{
				c.newRequest(c.apiUrl("/Query")).SetBody([]int{1, 10}),
				c.newRequest(c.rpcUrl()).SetBody(c.wrapRPCBody("test.Query", 1, 10)),
			},
			handle: func(t *testing.T, index int, resp *resty.Response) {
				require.True(t, !resp.IsError(), index)
				type Data struct {
					Result string
				}
				d := &Data{}
				require.NoError(t, json.Unmarshal(resp.Body(), d), index)
				require.Equal(t, "page:1 limit:10", d.Result, index)
			},
		},
	})
}

func TestOpenAPI(t *testing.T) {
	//doc := &openapi3.T{}
	//doc.Paths = map[string]*openapi3.PathItem{
	//	"/Ping": {
	//		Post: &openapi3.Operation{
	//			//Extensions:  map[string]interface{}{},
	//			Tags:        []string{"Tag1"},
	//			Summary:     "一点简介",
	//			Description: "一点描述",
	//			OperationID: "TestPing",
	//			Parameters: openapi3.Parameters{
	//				{
	//					Ref:   "",
	//					Value: &openapi3.Parameter{}, // TODO 不能为空
	//				},
	//			},
	//
	//			RequestBody: &openapi3.RequestBodyRef{
	//				Value: &openapi3.RequestBody{
	//					Required: true,
	//					Content: map[string]*openapi3.MediaType{
	//						"application/json": &openapi3.MediaType{
	//							Schema: &openapi3.SchemaRef{
	//								Ref: "#/components/schemas/centaurus.v1.LoginByLdapRequest",
	//							},
	//						},
	//					},
	//				},
	//			},
	//			Responses: map[string]*openapi3.ResponseRef{},
	//		},
	//	},
	//}
	//d, err := doc.MarshalJSON()
	//require.NoError(t, err)
	//t.Log(string(d))
}
