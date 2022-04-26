package repo

import (
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/repo/add"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/repo/new"
	"github.com/spf13/cobra"
)

// CmdRepo represents the repository command.
var CmdRepo = &cobra.Command{
	Use:   "repo",
	Short: "Generate the repository files.",
	Long:  "Generate the repository files.",
	Run:   run,
}

func init() {
	CmdRepo.AddCommand(new.CmdNew)
	CmdRepo.AddCommand(add.CmdAdd)
}

func run(cmd *cobra.Command, args []string) {
}
