package add

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/repo/project"
	"github.com/spf13/cobra"
)

var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Create a repository's service in the base repository directory",
	Long:  "Create a repository's service in the base repository directory: Example: kratos repo add app/user",
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
	CmdAdd.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdAdd.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdAdd.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, e := os.Stat(path.Join(wd, "go.mod")); os.IsNotExist(e) {
		fmt.Printf("🚫 go.mod don't exists in %s\n", wd)
		return
	}

	mod, err := base.ModulePath(path.Join(wd, "go.mod"))
	if err != nil {
		panic(err)
	}

	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	addPath := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project path ?",
			Help:    "Created project path, example: /app/project",
		}
		err = survey.AskOne(prompt, &addPath)
		if err != nil || addPath == "" {
			return
		}
	} else {
		addPath = args[0]
	}

	p := &project.Project{Name: path.Base(addPath), Path: addPath}
	done := make(chan error, 1)
	go func() {
		done <- p.Add(ctx, wd, repoURL, branch, mod)
	}()
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: add service timed out\033[m\n")
		} else {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to add service(%s)\033[m\n", ctx.Err().Error())
		}
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to add service(%s)\033[m\n", err.Error())
		}
	}
}
