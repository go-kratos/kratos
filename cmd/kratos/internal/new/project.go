package new

import (
	"context"
	"path"

	"github.com/go-kratos/kratos/cmd/kratos/internal/base"
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
	repo := base.NewRepo()
	return repo.CopyTo(ctx, serviceLayoutName, serviceLayoutURL, to)
}
