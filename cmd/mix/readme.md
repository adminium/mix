# Mix Cli

## 安装

```shell
$ go install github.com/gozelle/mix/cmd/mix
```

## 项目初始化

```shell
$ mix new hello
$ cd hello 
$ make run
```


## Client 生成
```
# 生成 jsonrpc client
mix generate client --path ./example/api --pkg example_api --outpkg example_api --outfile ./example/api/proxy_gen.go
```

## OpenAPI 文件生成
> 注意：仅支持项目直接或间接依赖的包引用，建议生成前先使用: go mod tidy
```
# 生成 Openapi 文件
mix generate openapi --path ./example/api --interface FullAPI --outfile ./openapi.json 
```

## 根据 OpenAPI 文件生成 SDK

```shell
# 根据 Openapi 文件生成 typescript axios SDK
mix generate sdk --openapi ./example/dist/openapi.json --sdk axios --outdir ./example/dist/sdk
```