package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
)

var projects = map[string]string{
	"service": "https://github.com/go-kratos/kratos-layout.git",
	"admin":   "https://github.com/go-kratos/kratos-admin.git",
}

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: kratos new helloworld",
	Run:   run,
}

var (
	nomod   bool
	branch  string
	timeout = "60s"
)

func init() {
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdNew.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
	CmdNew.Flags().BoolVarP(&nomod, "nomod", "", nomod, "retain go mod")
}

func run(_ *cobra.Command, args []string) {
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
	projectName, workingDir := processProjectParams(name, wd)
	p := &Project{Name: projectName}
	done := make(chan error, 1)
	repoURL, err := selectRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to select repo(%s)\033[m\n", err.Error())
		return
	}
	go func() {
		if !nomod {
			done <- p.New(ctx, workingDir, repoURL, branch)
			return
		}
		projectRoot := getgomodProjectRoot(workingDir)
		if gomodIsNotExistIn(projectRoot) {
			done <- fmt.Errorf("🚫 go.mod don't exists in %s", projectRoot)
			return
		}

		packagePath, e := filepath.Rel(projectRoot, filepath.Join(workingDir, projectName))
		if e != nil {
			done <- fmt.Errorf("🚫 failed to get relative path: %v", err)
			return
		}
		packagePath = strings.ReplaceAll(packagePath, "\\", "/")

		mod, e := base.ModulePath(filepath.Join(projectRoot, "go.mod"))
		if e != nil {
			done <- fmt.Errorf("🚫 failed to parse `go.mod`: %v", e)
			return
		}
		// Get the relative path for adding a project based on Go modules
		p.Path = filepath.Join(strings.TrimPrefix(workingDir, projectRoot+"/"), p.Name)
		done <- p.Add(ctx, workingDir, repoURL, branch, mod, packagePath)
	}()
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
			return
		}
		fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
		}
	}
}

func processProjectParams(projectName string, workingDir string) (projectNameResult, workingDirResult string) {
	_projectDir := projectName
	_workingDir := workingDir
	// Process ProjectName with system variable
	if strings.HasPrefix(projectName, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// cannot get user home return fallback place dir
			return _projectDir, _workingDir
		}
		_projectDir = filepath.Join(homeDir, projectName[2:])
	}

	// check path is relative
	if !filepath.IsAbs(projectName) {
		absPath, err := filepath.Abs(projectName)
		if err != nil {
			return _projectDir, _workingDir
		}
		_projectDir = absPath
	}

	return filepath.Base(_projectDir), filepath.Dir(_projectDir)
}

func getgomodProjectRoot(dir string) string {
	if dir == filepath.Dir(dir) {
		return dir
	}
	if gomodIsNotExistIn(dir) {
		return getgomodProjectRoot(filepath.Dir(dir))
	}
	return dir
}

func gomodIsNotExistIn(dir string) bool {
	_, e := os.Stat(filepath.Join(dir, "go.mod"))
	return os.IsNotExist(e)
}

func selectRepo() (string, error) {
	var (
		choice    string
		customURL string
	)
	form := huh.NewForm(
		// 1) Select group (always visible)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a template").
				Options(
					huh.NewOption("Service", "service"),
					huh.NewOption("Admin", "admin"),
					huh.NewOption("Custom (enter repo URL)", "custom"),
				).
				Value(&choice),
		),
		// 2) Input group (only visible when choice == "custom")
		huh.NewGroup(
			huh.NewInput().
				Title("Enter custom repository URL").
				Placeholder("https://github.com/owner/repo.git").
				Value(&customURL).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return fmt.Errorf("repo URL cannot be empty")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return choice != "custom"
		}),
	)
	if err := form.Run(); err != nil {
		panic(err)
	}
	if choice == "custom" {
		return strings.TrimSpace(customURL), nil
	}
	return projects[choice], nil
}
