package generate

import (
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成 RPC Client、API SDK",
	Long:  ``,
}

func init() {
	GenerateCmd.AddCommand(
		clientCmd,
		openapiCmd,
		sdkCmd,
	)
}
