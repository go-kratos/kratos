package main

import (
	"log"

	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/change"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/project"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/proto"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/run"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/upgrade"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kratos",
	Short:   "Kratos: An elegant toolkit for Musicverse.",
	Long:    `Kratos: An elegant toolkit for Musicverse.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
