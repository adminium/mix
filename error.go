package mix

import (
	"fmt"

	"github.com/adminium/jsonrpc"
)

type Error = jsonrpc.Error

type Warning = jsonrpc.Warning

var _ error = (*Warning)(nil)

func Warnf(format string, a ...any) *Warning {
	return &Warning{
		Message: fmt.Sprintf(format, a...),
	}
}

func Codef(code int, format string, a ...any) *Warning {
	return &Warning{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
