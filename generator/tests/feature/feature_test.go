package tests_feature

import (
	"encoding/json"
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/openapi"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/spew"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestFeature(t *testing.T) {
	mod, err := parser.PrepareMod()
	require.NoError(t, err)
	
	dir, err := fs.Lookup("./generator/tests/feature")
	require.NoError(t, err)
	
	pkg, err := parser.Parse(mod, dir)
	require.NoError(t, err)
	
	defJson, err := json.Marshal(pkg.GetInterface("FeatureAPI"))
	require.NoError(t, err)
	
	t.Log(string(defJson))
	
	//c, err := fs.Read("./parser_basic_def.json")
	//require.NoError(t, err)
	//
	//// 比较 Parser 类型定义
	//err = fastjson.EqualsBytes(c, defJson)
	//require.NoError(t, err)
	//
	parserInterface := pkg.GetInterface("FeatureAPI")
	require.True(t, parserInterface != nil)
	
	renderInterface := openapi.ConvertAPI(parserInterface)
	d, _ := json.Marshal(renderInterface)
	t.Log(string(d))
	
	doc := &openapi.DocumentV3{}
	
	openapi.ConvertOpenapi(doc, renderInterface)
	
	d, err = doc.MarshalJSON()
	require.NoError(t, err)
	
	//
	t.Log(string(d))
	//err = fs.Write("./render_interface.json", d)
	//require.NoError(t, err)
	
	//t.Log(string(d))
}

type A struct {
	Name  string
	Child *A
}

func TestJson(t *testing.T) {
	a := A{
		Name: "a",
		Child: &A{
			Name: "b",
		},
	}
	spew.Json(a)
}
