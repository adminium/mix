package generate

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate RPC client、API SDK",
	Long:  ``,
}

func init() {
	Cmd.AddCommand(
		clientCmd,
		openapiCmd,
		sdkCmd,
	)
}
