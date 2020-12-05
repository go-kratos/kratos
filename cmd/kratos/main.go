package main

import (
	"log"

	"github.com/go-kratos/kratos/v2/cmd/kratos/internal/gen"
	"github.com/go-kratos/kratos/v2/cmd/kratos/internal/new"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kratos",
	Short: "Kratos: An elegant toolkit for Go microservices.",
	Long:  `Kratos: An elegant toolkit for Go microservices.`,
}

func init() {
	rootCmd.AddCommand(new.CmdNew)
	rootCmd.AddCommand(gen.CmdGen)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
