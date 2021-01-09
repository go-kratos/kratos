package main

import (
	"log"

	"github.com/go-kratos/kratos/cmd/kratos/internal/new"
	"github.com/go-kratos/kratos/cmd/kratos/internal/proto"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kratos",
	Short:   "Kratos: An elegant toolkit for Go microservices.",
	Long:    `Kratos: An elegant toolkit for Go microservices.`,
	Version: Version,
}

func init() {
	rootCmd.AddCommand(new.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
