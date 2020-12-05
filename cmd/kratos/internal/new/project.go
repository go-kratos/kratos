package new

// Project is a project template.
type Project struct {
	Name string
	Repo *Repository
}

// Generate .
func (p *Project) Generate(path string) error {

	return nil
}
