package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/adminium/mix/generator/openapi"
	"github.com/adminium/mix/generator/parser"
	"github.com/gozelle/spew"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {

	dir, err := os.Getwd()
	require.NoError(t, err)

	t.Log(dir)
	mod, err := parser.PrepareMod(dir)
	if err != nil {
		err = fmt.Errorf("prepare mod error: %s", err)
		return
	}

	pkg, err := parser.Parse(mod, dir)
	require.NoError(t, err)

	inter := pkg.GetInterface("SimpleAPI")
	require.True(t, inter != nil)

	err = inter.Load()
	require.NoError(t, err)

	doc := &openapi.DocumentV3{}

	r := openapi.ConvertAPI(inter)
	openapi.ConvertOpenapi(doc, r)

	spew.Json(doc)
}
