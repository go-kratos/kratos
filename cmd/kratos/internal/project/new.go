package project

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
)

// Project is a project template.
type Project struct {
	Name       string
	Path       string
	ModuleName string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := filepath.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		var override bool
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	modName, notifyFunc := p.getModName()
	repo := base.NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, modName, []string{".git", ".github"}); err != nil {
		return err
	}
	e := os.Rename(
		filepath.Join(to, "cmd", "server"),
		filepath.Join(to, "cmd", p.Name),
	)
	if e != nil {
		return e
	}
	e = base.ModuleName(
		path.Join(to, "go.mod"),
		p.ModuleName,
		p.Name,
	)
	if e != nil {
		return e
	}
	base.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	notifyFunc()
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go generate ./..."))
	fmt.Println(color.WhiteString("$ go build -o ./bin/ ./... "))
	fmt.Println(color.WhiteString("$ ./bin/%s -conf ./configs\n", p.Name))
	fmt.Println("			🤝 Thanks for using Kratos")
	fmt.Println("	📚 Tutorial: https://go-kratos.dev/docs/getting-started/start")
	return nil
}

func (p *Project) getModName() (string, func()) {
	if p.ModuleName != "" {
		return p.ModuleName, func() {
			fmt.Printf("👉 Module name replace succeeded %s\n", color.BlueString(p.ModuleName))
		}
	}

	return p.Name, func() {}
}
