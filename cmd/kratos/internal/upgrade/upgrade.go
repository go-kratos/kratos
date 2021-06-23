package upgrade

import (
	"fmt"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"

	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the kratos tools",
	Long:  "Upgrade the kratos tools. Example: kratos upgrade",
	Run:   Run,
}

// Run upgrade the kratos tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoGet(
		"github.com/go-kratos/kratos/cmd/kratos/v2",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/envoyproxy/protoc-gen-validate",
	)
	if err != nil {
		fmt.Println(err)
	}
}
