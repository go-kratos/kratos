package new

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-kratos/kratos/cmd/kratos/internal/base"
	"golang.org/x/mod/modfile"
)

const (
	serviceLayoutName = "service"
	serviceLayoutURL  = "https://github.com/go-kratos/service-layout.git"
)

// Project is a project template.
type Project struct {
	Name string
}

// Generate .
func (p *Project) Generate(ctx context.Context, dir string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", p.Name)
	}
	fmt.Printf("Creating service %s\n", p.Name)
	repo := base.NewRepo()
	if err := repo.CopyTo(ctx, serviceLayoutName, serviceLayoutURL, to); err != nil {
		return err
	}
	mod := path.Join(to, "go.mod")
	modBytes, err := ioutil.ReadFile(mod)
	if err != nil {
		return err
	}
	modName := modfile.ModulePath(modBytes)
	return ioutil.WriteFile(mod, bytes.Replace(modBytes, []byte(modName), []byte(p.Name), 1), 0644)
}
