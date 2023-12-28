package mix

import (
	"fmt"
	"strings"

	"github.com/adminium/jsonrpc"
	"github.com/adminium/mix/middlewares/cors"
	"github.com/gin-gonic/gin"
)

func NewServer() *Server {
	s := &Server{}
	s.Engine = gin.New()
	s.Engine.Use(Logger(), gin.Recovery(), cors.Cors())
	return s
}

type Server struct {
	*gin.Engine
}

func (s *Server) RegisterRPC(router gin.IRouter, namespace string, handler any, middlewares ...gin.HandlerFunc) {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)
	router.POST("", append([]gin.HandlerFunc{gin.WrapH(rpcServer)}, middlewares...)...)
}

func (s *Server) RegisterAPI(router gin.IRouter, namespace string, handler any, middlewares ...gin.HandlerFunc) {

	var path string
	if strings.TrimSpace(namespace) == "" {
		path = fmt.Sprintf("/:%s", method)
	} else {
		path = fmt.Sprintf("/:%s/:%s", module, method)
	}

	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)

	router.POST(path, append([]gin.HandlerFunc{wrapAPI(namespace, rpcServer)}, middlewares...)...)
}
