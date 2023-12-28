package mix

import (
	"fmt"

	"github.com/adminium/jsonrpc"
)

type Error = jsonrpc.Error

func NewError(code int, message string) *Error {
	return &Error{Code: jsonrpc.ErrorCode(code), Message: message}
}

func Errorf(format string, a ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func Codef(code int, format string, a ...any) *Error {
	return &Error{
		Code:    jsonrpc.ErrorCode(code),
		Message: fmt.Sprintf(format, a...),
	}
}
