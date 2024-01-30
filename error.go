package mix

import (
	"fmt"

	"github.com/adminium/jsonrpc"
)

type Warning = jsonrpc.Warning

var _ error = (*Warning)(nil)

func Wranf(format string, a ...any) *Warning {
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
