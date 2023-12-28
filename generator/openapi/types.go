package openapi

import "github.com/gozelle/openapi/openapi3"

type DocumentV3 struct {
	openapi3.T
}

const (
	ApplicationJson = "application/json"
)

const (
	Array       = "array"
	ArrayParams = "array_params"
	String      = "string"
	Integer     = "integer"
	Number      = "number"
	Boolean     = "boolean"
	Object      = "object"
)
