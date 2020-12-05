package new

import "os"

// Repository is a repository template.
type Repository struct {
	Name       string
	Path       string
	References string
}

// Clone .
func (r *Repository) Clone() error {
	return nil
}

// ListFiles .
func (r *Repository) ListFiles() ([]*os.File, error) {
	return nil, nil
}

// CopyTo .
func (r *Repository) CopyTo(path string) error {
	return nil
}
