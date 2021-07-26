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
	paths := []string{
		"github.com/go-kratos/kratos/cmd/kratos/v2",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/envoyproxy/protoc-gen-validate",
	}
	goVersion, err := base.GetGoVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	if goVersion > "1.17" {
		err = base.GoInstall(paths...)
	} else {
		err = base.GoGet(paths...)
	}

	if err != nil {
		fmt.Println(err)
	}
}
