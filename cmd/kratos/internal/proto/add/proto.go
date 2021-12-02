package add

import (
	"fmt"
	"os"
	"path"
)

// Proto is a proto generator.
type Proto struct {
	Name        string
	Path        string
	Service     string
	Package     string
	GoPackage   string
	JavaPackage string
}

// Generate generate a proto template.
func (p *Proto) Generate() error {
	body, err := p.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, p.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0o700); err != nil {
			return err
		}
	}
	name := path.Join(to, p.Name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", p.Name)
	}
	return os.WriteFile(name, body, 0o644)
}
