package openapi

import (
	"encoding/json"
	"testing"

	"github.com/adminium/fs"
	"github.com/adminium/mix/generator/parser"
	"github.com/stretchr/testify/require"
)

const basicRender = `
{
	"Type": "struct",
	"StructFields": [
		{
			"Field": "Int",
			"Type": "int"
		},
		{
			"Type": "uint",
			"Pointer": true
		},
		{
			"Field": "Bool",
			"Type": "bool"
		},
		{
			"Type": "float32",
			"Pointer": true
		},
		{
			"Field": "String",
			"Type": "string"
		},
		{
			"Field": "Map",
			"Type": "map"
		},
		{
			"Field": "IntArray",
			"Json": "int_array",
			"Type": "slice",
			"Elem": {
				"Type": "int"
			},
			"Tags": "json:\"int_array\""
		},
		{
			"Field": "UintArray",
			"Json": "uint_array",
			"Type": "slice",
			"Elem": {
				"Type": "uint",
				"Pointer": true
			},
			"Tags": "json:\"uint_array\""
		},
		{
			"Field": "BoolArray",
			"Json": "bool_array",
			"Type": "slice",
			"Elem": {
				"Type": "bool"
			},
			"Tags": "json:\"bool_array\""
		},
		{
			"Field": "FloatArray",
			"Json": "float_array",
			"Type": "slice",
			"Elem": {
				"Type": "float32",
				"Pointer": true
			},
			"Tags": "json:\"float_array\""
		},
		{
			"Field": "StringArray",
			"Json": "string_array",
			"Type": "slice",
			"Elem": {
				"Type": "string"
			},
			"Tags": "json:\"string_array\""
		},
		{
			"Field": "StringArrayArray",
			"Json": "string_array_array",
			"Type": "slice",
			"Elem": {
				"Type": "slice",
				"Elem": {
					"Type": "string"
				}
			},
			"Tags": "json:\"string_array_array\""
		}
	]
}
`

func TestConvertType(t *testing.T) {
	path, err := fs.LookupPwd("./generator/tests/basic/parser_basic_def.json")
	require.NoError(t, err)

	d, err := fs.Read(path)
	require.NoError(t, err)

	pDef := &parser.Def{}
	err = json.Unmarshal(d, pDef)
	require.NoError(t, err)

	def := convertRenderType(pDef.Type)

	d, err = json.Marshal(def)
	require.NoError(t, err)

	//err = fastjson.Equals(basicRender, string(d))
	//require.NoError(t, err)
}
