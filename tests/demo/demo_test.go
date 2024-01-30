package demo

import (
	"fmt"
	"os"
	"testing"

	"github.com/adminium/fs"
	"github.com/adminium/mix/generator/openapi"
	"github.com/adminium/mix/generator/parser"
	"github.com/stretchr/testify/require"
)

func TestMixServer(t *testing.T) {

}

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

	d, err := doc.MarshalJSON()
	require.NoError(t, err)

	file := "./simple_api.json"
	err = fs.Write(file, d)
	require.NoError(t, err)

	t.Logf("write file: %s", file)
}

func TestParseSDT(t *testing.T) {

	dir := "/Users/koyeo/shadowtokens/shadowtokens"
	t.Log(dir)

	mod, err := parser.PrepareMod(dir)
	if err != nil {
		err = fmt.Errorf("prepare mod error: %s", err)
		return
	}

	pkg, err := parser.Parse(mod, fs.Join(dir, "module/client/api"))
	require.NoError(t, err)

	inter := pkg.GetInterface("FullAPI")
	require.True(t, inter != nil)

	err = inter.Load()
	require.NoError(t, err)

	doc := &openapi.DocumentV3{}

	r := openapi.ConvertAPI(inter)
	openapi.ConvertOpenapi(doc, r)

	d, err := doc.MarshalJSON()
	require.NoError(t, err)

	file := "./sdt.json"
	err = fs.Write(file, d)
	require.NoError(t, err)

	t.Logf("write file: %s", file)
}
