package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/base"
)

const (
	goModFileName  = "go.mod"
	goSumFileName  = "go.sum"
	readmeFileName = "README.md"
)

var repoAddIgnores = []string{
	".git", ".github", "api", readmeFileName, "LICENSE", goModFileName, goSumFileName, "third_party", "openapi.yaml", ".gitignore",
}

func (p *Project) Add(ctx context.Context, dir string, layout string, branch string, mod string, pkgPath string) error {
	to := filepath.Join(dir, p.Name)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Printf("🚀 Add service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	pkgPath = fmt.Sprintf("%s/%s", mod, pkgPath)
	repo := base.NewRepo(layout, branch)
	err := repo.CopyToV2(ctx, to, pkgPath, repoAddIgnores, []string{filepath.Join(p.Path, "api"), "api"})
	if err != nil {
		return err
	}

	e := os.Rename(
		filepath.Join(to, "cmd", "server"),
		filepath.Join(to, "cmd", p.Name),
	)
	if e != nil {
		if !os.IsNotExist(e) {
			return e
		}
	}

	base.Tree(to, dir)

	fmt.Printf("\n🍺 Repository creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to add a project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go generate ./..."))
	fmt.Println(color.WhiteString("$ go build -o ./bin/ ./... "))
	fmt.Println(color.WhiteString("$ ./bin/%s -conf ./configs\n", p.Name))
	fmt.Println("			🤝 Thanks for using Kratos")
	fmt.Println("	📚 Tutorial: https://go-kratos.dev/docs/getting-started/start")
	return nil
}
