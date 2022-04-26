package new

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/repo/project"
	"github.com/spf13/cobra"
)

var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a base repository template",
	Long:  "Create a base repository template: Example: kratos reqo new my-project",
	Run:   run,
}

var (
	repoURL string
	branch  string
	timeout string
)

func init() {
	if repoURL = os.Getenv("KRATOS_LAYOUT_REPO"); repoURL == "" {
		repoURL = "https://github.com/go-kratos/kratos-layout.git"
	}
	timeout = "60s"
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdNew.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}

	p := &project.Project{Name: path.Base(name), Path: name}
	done := make(chan error, 1)
	go func() {
		done <- p.New(ctx, wd, repoURL, branch)
	}()
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: repository creation timed out\033[m\n")
		} else {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create repository(%s)\033[m\n", ctx.Err().Error())
		}
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create repository(%s)\033[m\n", err.Error())
		}
	}
}
