package proto

import (
	"github.com/go-kratos/kratos/cmd/kratos/internal/proto/add"
	"github.com/go-kratos/kratos/cmd/kratos/internal/proto/service"
	"github.com/go-kratos/kratos/cmd/kratos/internal/proto/source"

	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
	Run:   run,
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(source.CmdSource)
	CmdProto.AddCommand(service.CmdService)
}

func run(cmd *cobra.Command, args []string) {

}
