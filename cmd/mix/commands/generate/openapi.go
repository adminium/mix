package cmdGenerate

import (
	"os"
	"path/filepath"

	"github.com/adminium/mix/cmd/mix/commands"
	"github.com/adminium/mix/generator/openapi"
	"github.com/gozelle/fs"
	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "文件生成 OpenAPI 文件",
	Long:  ``,
	Run:   generateOpenapi,
}

var (
	openapiTpl       string
	openapiPath      string
	openapiInterface string
	openapiOutfile   string
)

func init() {
	openapiCmd.Flags().StringVar(&openapiPath, "path", "", "[必填]源目录")
	openapiCmd.Flags().StringVar(&openapiInterface, "interface", "", "[必填]源 Interface 名")
	openapiCmd.Flags().StringVar(&openapiOutfile, "outfile", "", "[必填]生成文件路径")
	openapiCmd.Flags().StringVar(&openapiTpl, "template", "", "[可选] OpenAPI 文件模板")
	openapiCmd.MarkFlagsRequiredTogether("path", "interface", "outfile")
}

func generateOpenapi(cmd *cobra.Command, args []string) {

	pwd, err := os.Getwd()
	if err != nil {
		commands.Fatal(err)
	}

	openapiPath = filepath.Join(pwd, openapiPath)
	openapiOutfile = filepath.Join(pwd, openapiOutfile)
	if openapiTpl != "" {
		openapiTpl = filepath.Join(pwd, openapiTpl)
	}

	doc, err := openapi.Parse(openapiTpl, openapiPath, openapiInterface)
	if err != nil {
		commands.Fatal(err)
	}

	c, err := doc.MarshalJSON()
	if err != nil {
		commands.Fatal(err)
	}

	err = fs.Write(openapiOutfile, c)
	if err != nil {
		commands.Fatal(err)
	}

	commands.Info("write file: %s", openapiOutfile)
}
