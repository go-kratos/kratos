package proto

import (
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/proto/add"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/proto/client"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/proto/server"

	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:       "proto",
	Short:     "Generate the proto files",
	Long:      "Generate the proto files.",
	ValidArgs: []string{"add", "client", "server"},
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(client.CmdClient)
	CmdProto.AddCommand(server.CmdServer)
}
