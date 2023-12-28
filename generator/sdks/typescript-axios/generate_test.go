package typescript_axios

import (
	"testing"

	"github.com/adminium/fs"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	file, err := fs.LookupPwd("./generator/tests/feature/openapi.json")
	require.NoError(t, err)

	files, err := Generate(file, "")
	require.NoError(t, err)

	for _, v := range files {
		t.Log("生成文件:", v.Name)
		t.Log(v.Content)
	}
}
