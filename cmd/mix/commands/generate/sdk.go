package cmdGenerate

import (
	"fmt"
	"os"

	"github.com/adminium/fs"
	"github.com/adminium/mix/cmd/mix/commands"
	typescript_axios "github.com/adminium/mix/generator/sdks/typescript-axios"
	"github.com/adminium/mix/generator/writter"
	"github.com/spf13/cobra"
)

var sdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: "生成 API SDK",
	Long:  ``,
	Example: `
mix generate sdk --openapi ./example/openapi.json --sdk axios --outdir ./sdk # 生成 Axios SDK 
`,
	Run: generateSDK,
}

var (
	sdkOpenapi string
	sdkType    string
	sdkOutdir  string
	sdkOptions string
)

func init() {
	sdkCmd.Flags().StringVar(&sdkOpenapi, "openapi", "", "[必填] OpenAPI 文件")
	sdkCmd.Flags().StringVar(&sdkType, "sdk", "", "[必填] SDK 类型，如：axios")
	sdkCmd.Flags().StringVar(&sdkOutdir, "outdir", "", "[必填] SDK 存放目录")
	sdkCmd.Flags().StringVar(&sdkOptions, "option", "", "[可选]配置选项，请查看不同 SDK 配置选项")
	sdkCmd.MarkFlagsRequiredTogether("openapi", "sdk", "outdir")
}

func generateSDK(cmd *cobra.Command, args []string) {
	pwd, err := os.Getwd()
	if err != nil {
		commands.Fatal(err)
	}
	sdkOpenapi = fs.Join(pwd, sdkOpenapi)
	sdkOutdir = fs.Join(pwd, sdkOutdir)

	if fs.Exist(sdkOutdir) {
		if !fs.IsDir(sdkOutdir) {
			commands.Fatal(fmt.Errorf("outdir: %s is not dir", sdkOutdir))
		}
	} else {
		err = fs.MakeDir(sdkOutdir)
		if err != nil {
			commands.Fatal(fmt.Errorf("make outdir: %s error: %s", sdkOutdir, err))
		}
		commands.Info("make dir: %s", sdkOutdir)
	}

	var files []*writter.File
	switch sdkType {
	case "axios":
		files, err = typescript_axios.Generate(sdkOpenapi, sdkOptions)
	default:
		commands.Fatal(fmt.Errorf("sdk type: %s unsupported", sdkType))
	}
	if err != nil {
		commands.Fatal(fmt.Errorf("generate error: %s", err))
	}

	paths, err := writter.WriteFiles(sdkOutdir, files)
	if err != nil {
		commands.Fatal(fmt.Errorf("write file error: %s", err))
	}

	for _, v := range paths {
		commands.Info("write file: %s", v)
	}
}
