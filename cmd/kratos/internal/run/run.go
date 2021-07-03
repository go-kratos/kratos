package run

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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
		// find the directory containing the cmd/*/main.go
		mainDir, err := searchMain()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		dir = path.Join(base, mainDir)
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

func searchMain() (string, error) {
	var mainPath string
	err := filepath.Walk(".", func(walkPath string, info os.FileInfo, err error) error {
		if strings.HasPrefix(walkPath, "cmd") && info.Name() == "main.go" {
			mainPath = walkPath
			_ = filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(mainPath, "main.go")
	if index != -1 {
		mainPath = mainPath[:index]
	}
	return mainPath, err
}
