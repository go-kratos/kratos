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
		cmdDir, err := findCMD(base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if len(cmdDir) == 0 {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
			return
		} else if len(cmdDir) == 1 {
			dir = cmdDir[0]
		} else {
			prompt := &survey.Select{
				Message: "Which directory do you want to run?",
				Options: cmdDir,
			}
			survey.AskOne(prompt, &dir)
		}
	}
	fd := exec.Command("go", "run", ".")
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	if fd.Dir == "" {
		return
	}
	if err := fd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}
	return
}

func findCMD(dir string) ([]string, error) {
	hasCMDPaths, err := walk(dir)
	if err != nil {
		return []string{}, err
	}
	if len(hasCMDPaths) == 0 {
		// find go.mod from the parent directory.
		p := dir
		d := ""
		for {
			d = filepath.Join(p, "..")
			if d == p {
				return []string{}, nil
			}
			paths, err := ioutil.ReadDir(p)
			if err != nil {
				return []string{}, err
			}
			for _, fileInfo := range paths {
				if fileInfo.Name() == "go.mod" {
					hasCMDPaths, err = walk(p)
					if err != nil {
						return []string{}, err
					}
					return hasCMDPaths, nil
				}
			}
			p = d
		}
	}
	return hasCMDPaths, err
}

func walk(dir string) ([]string, error) {
	var hasCMDPaths []string
	err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
		// multi level directory is not allowed under the cmd directory, so it is judged that the path ends with cmd.
		if strings.HasSuffix(walkPath, "cmd") {
			paths, err := ioutil.ReadDir(walkPath)
			if err != nil {
				return err
			}
			for _, fileInfo := range paths {
				if fileInfo.IsDir() {
					hasCMDPaths = append(hasCMDPaths, path.Join(walkPath, fileInfo.Name()))
				}
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return hasCMDPaths, nil
}