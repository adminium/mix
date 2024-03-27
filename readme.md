# Mix

## Background


API frameworks are the foundational components of application development. 
In the past, it was common to manually define routes using Gin. 
Now, we simplify the process of defining interfaces by leveraging JSON-RPC.

## Features 

- Some optimizations have been made based on JSONRPC, such as internal method masking, error handling, logging, etc. For more details, please refer to the link.
- Automatically implement corresponding interface definitions based on Interface and Struct public methods.




## Quick Start


API Service:

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
	
	s.Run(":1323")
}
```

API Request:

```
curl --location --request POST '127.0.0.1:11111/api/v1/Ping' \
--header 'Content-Type: text/plain' \
--data-raw '"Hello world!"'
```

RPC Request:
```
curl --location --request POST '127.0.0.1:1332/rpc/v1/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "foo.Ping",
    "params": ["Hello world!"]
}'
```

## API Methods and Parameters
For the registered Struct objects, JSONRPC will scan their exported methods. 
For the exported methods that conform to the interface method format, 
JSONRPC will automatically proxy the API routing to execute this method. The format is as follows:

```
[prefix]/MehtodName
```

Supported function signatures:

```go
type API interface {
    // No Params / Return val
    Func1()
    
    // With Params
    // Note: If param types implement json.[Un]Marshaler, go-jsonrpc will use it
    Func2(param1 int, param2 string, param3 struct{A int})
    
    // Returning errors
    // * For some connection errors, go-jsonrpc will return jsonrpc.RPCConnectionError{}.
    // * RPC-returned errors will be constructed with basic errors.New(__"string message"__)
    // * JSON-RPC error codes can be mapped to typed errors with jsonrpc.Errors - https://pkg.go.dev/github.com/filecoin-project/go-jsonrpc#Errors
    // * For typed errors to work, server needs to be constructed with the `WithServerErrors`
    //   option, and the client needs to be constructed with the `WithErrors` option
    Func3() error
    
    // Returning a value
    // Note: The value must be serializable with encoding/json.
    Func4() int
    
    // Returning a value and an error
    // Note: if the handler returns an error and a non-zero value, the value will not
    //       be returned to the client - the client will see a zero value.
    Func4() (int, error)
    
    // With context
    // * Context isn't passed as JSONRPC param, instead it has a number of different uses
    // * When the context is cancelled on the client side, context cancellation should propagate to the server handler
    //   * In http mode the http request will be aborted
    //   * In websocket mode the client will send a `xrpc.cancel` with a single param containing ID of the cancelled request
    // * If the context contains an opencensus trace span, it will be propagated to the server through a
    //   `"Meta": {"SpanContext": base64.StdEncoding.EncodeToString(propagation.Binary(span.SpanContext()))}` field in
    //   the jsonrpc request
    //   
    Func5(ctx context.Context, param1 string) error
    
    // With non-json-serializable (e.g. interface) params
    // * There are client and server options which make it possible to register transformers for types
    //   to make them json-(de)serializable
    // * Server side: jsonrpc.WithParamDecoder(new(io.Reader), func(ctx context.Context, b []byte) (reflect.Value, error) { ... }
    // * Client side: jsonrpc.WithParamEncoder(new(io.Reader), func(value reflect.Value) (reflect.Value, error) { ... }
    // * For io.Reader specifically there's a simple param encoder/decoder implementation in go-jsonrpc/httpio package
    //   which will pass reader data through separate http streams on a different hanhler.
    // * Note: a similar mechanism for return value transformation isn't supported yet
    Func6(r io.Reader)
    
    // Returning a channel
    // * Only supported in websocket mode
    // * If no error is returned, the return value will be an int channelId
    // * When the server handler writes values into the channel, the client will receive `xrpc.ch.val` notifications
    //   with 2 params: [chanID: int, value: any]
    // * When the channel is closed the client will receive `xrpc.ch.close` notification with a single param: [chanId: int]
    // * The client-side channel will be closed when the websocket connection breaks; Server side will discard writes to
    //   the channel. Handlers should rely on the context to know when to stop writing to the returned channel.
    // NOTE: There is no good backpressure mechanism implemented for channels, returning values faster that the client can
    // receive them may cause memory leaks.
    Func7(ctx context.Context, param1 int, param2 string) (<-chan int, error)
}
```

## Bearer Context Injection
```go
func SetBearer(ctx *gin.Context, bearer string) {
	ctx.Header(X_Bearer, bearer)
}
```

## Client Generator

```shell
$ mix generate client 
```

## OpenAPI Generator
```shell
$ mix generate openapi
```

## SDK Generator
```shell
$ mix generate sdk
```

## Middlewares

- CORS

## Optimization

- Optimized error output:

- Extended warning error types to differentiate between business errors and system errors.




