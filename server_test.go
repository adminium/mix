package mix

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type ITestHandler interface {
	Ping(ctx context.Context, msg string) (reply string, err error)
	Query(ctx context.Context, page, limit int) (reply string, err error)
	Download(ctx context.Context, file string) io.Reader
	Error(ctx context.Context) error
	Code(ctx context.Context) error
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
	return fmt.Errorf("some error")
}

func (t TestHandler) Code(ctx context.Context) error {
	return &Error{
		Code:    1000,
		Message: "自定义错误",
	}
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

	server := NewServer()
	group := server.Group("/api/v1")

	server.RegisterRPC(server.Group("/rpc/v1"), "", h)
	server.RegisterAPI(group, "", h)

	group.GET("/download", WrapHandler(func(ctx *gin.Context) (data any, err error) {
		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		return h.Download(ctx, "ok"), nil
	}))

	require.NoError(t, server.Run(":11111"))
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
