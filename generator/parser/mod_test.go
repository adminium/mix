package parser

import (
	"os"
	"testing"

	"github.com/adminium/fs"
	"github.com/stretchr/testify/require"
)

func TestFindMod(t *testing.T) {

	pwd, _ := os.Getwd()
	pwd = fs.Join(pwd, "../../")
	mod, err := PrepareMod(pwd)
	require.NoError(t, err)

	t.Log(mod.root, mod.file.Module)
	t.Log(mod.Gopath())

	files, err := mod.GetPackageFiles("github.com/adminium/mix/cmd")
	require.NoError(t, err)
	t.Log(files)

	files, err = mod.GetPackageFiles("github.com/adminium/fs")
	require.NoError(t, err)
	t.Log(files)
}
