package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdRun run project command.
var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "run project",
	Long:  "run project. Example: kratos run",
	Run:   Run,
}

// Run run project.
func Run(cmd *cobra.Command, args []string) {
	var dir string
	if len(args) > 0 {
		dir = args[0]
	}
	base, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
	if dir == "" {
		// find the directory containing the cmd/*
		cmdDir, cmdPath, err := findCMD(base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if len(cmdPath) == 0 {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
			return
		} else if len(cmdPath) == 1 {
			dir = path.Join(cmdDir, cmdPath[0])
		} else {
			prompt := &survey.Select{
				Message: "Which directory do you want to run?",
				Options: cmdPath,
			}
			survey.AskOne(prompt, &dir)
			if dir == "" {
				return
			}
			dir = path.Join(cmdDir, dir)
		}
	}
	fd := exec.Command("go", "run", ".")
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	if err := fd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}
	return
}

func findCMD(base string) (string, []string, error) {
	var root bool
	next := func(dir string) (string, []string, error) {
		var (
			cmdDir  string
			cmdPath []string
		)
		err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			// multi level directory is not allowed under the cmd directory, so it is judged that the path ends with cmd.
			if strings.HasSuffix(walkPath, "cmd") {
				paths, err := ioutil.ReadDir(walkPath)
				if err != nil {
					return err
				}
				for _, fileInfo := range paths {
					if fileInfo.IsDir() {
						cmdPath = append(cmdPath, path.Join("cmd", fileInfo.Name()))
					}
				}
				cmdDir = filepath.Join(walkPath, "..")
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return cmdDir, cmdPath, err
	}
	for i := 0; i < 5; i++ {
		cmdDir, res, err := next(base)
		if err != nil {
			return "", nil, err
		}
		if len(res) > 0 {
			return cmdDir, res, nil
		}
		if root {
			break
		}
		base = filepath.Join(base, "..")
	}
	return "", []string{base}, nil
}
