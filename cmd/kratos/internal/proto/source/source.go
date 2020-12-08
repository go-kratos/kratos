package source

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// CmdSource represents the source command.
var CmdSource = &cobra.Command{
	Use:   "source",
	Short: "Generate the proto source code",
	Long:  "Generate the proto source code. Example: kratos proto source ./**/*.proto",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	input := []string{"--go_out=paths=source_relative:.", "--go-grpc_out=paths=source_relative:."}
	input = append(input, args...)
	do := exec.Command("protoc", input...)
	out, err := do.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to execute: %s\n", err)
	}
	fmt.Println(string(out))
}
