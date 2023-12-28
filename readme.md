# Mix

## 特性

* 基于 jsonrpc 的封装
* 支持自定义状态码渗透
* 支持 sdk 生成

## 快速上手

```go
package main

import (
	"context"

	"github.com/gozelle/mix"
)

type IHandler interface {
	Ping(ctx context.Context, message string) (reply string, err error)
}

var _ IHandler = (*Handler)(nil)

type Handler struct {
}

func (h Handler) Ping(ctx context.Context, message string) (reply string, err error) {
	reply = message
	return
}

func main() {
	h := &Handler{}
	s := mix.NewServer()
	
	s.RegisterRPC(s.Group("/rpc/v1"), "foo", h)
	s.RegisterAPI(s.Group("/api/v1"), "foo", h)
	
	s.Run(":1332")
}
```

请求示例

API:

```
curl --location --request POST '127.0.0.1:11111/api/v1/Ping' \
--header 'Content-Type: text/plain' \
--data-raw '"Hello world!"'
```

RPC:

```
curl --location --request POST '127.0.0.1:1332/rpc/v1/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "foo.Ping",
    "params": ["Hello world!"]
}'
```

## 安装命令行工具

```
go install github.com/gozelle/mix/cmd/mix@latest
mix -h
```

## Client 生成

```
# 生成 jsonrpc client
mix generate client --path ./example/api --pkg example_api --outpkg example_api --outfile ./example/api/proxy_gen.go
```

## SDK 生成
> 注意：仅支持项目直接或间接依赖的包引用，建议生成前先使用: go mod tidy
```
# go mod tidy

# 生成 Openapi 文件
mix generate openapi --path ./example/api --interface FullAPI --outfile ./openapi.json 

# 根据 Openapi 文件生成 typescript axios SDK
mix generate sdk --openapi ./example/dist/openapi.json --sdk axios --outdir ./example/dist/sdk
```
