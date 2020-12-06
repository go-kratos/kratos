package new

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// CmdNew represents the init command
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  `Create a service project using the repository template. Example: kratos new helloworld`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p := &Project{Name: args[0]}
	if err := p.Generate(ctx, wd); err != nil {
		log.Fatal(err)
	}
}
