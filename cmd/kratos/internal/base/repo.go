package base

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Repo is git repository manager.
type Repo struct {
	url  string
	home string
}

// NewRepo new a repository manager.
func NewRepo(url string) *Repo {
	return &Repo{
		url:  url,
		home: kratosHomeWithDir("repo"),
	}
}

// Path returns the repository cache path.
func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	return path.Join(r.home, r.url[start+1:end])
}

// Pull fetchs the repository from remote url.
func (r *Repo) Pull(ctx context.Context, url string) error {
	repo, err := git.PlainOpen(r.Path())
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	if err = w.PullContext(ctx, &git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	}); errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx, r.url)
	}
	_, err := git.PlainCloneContext(ctx, r.Path(), false, &git.CloneOptions{
		URL:      r.url,
		Progress: os.Stdout,
	})
	return err
}

// CopyTo copies the repository to project path.
func (r *Repo) CopyTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(path.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	return copyDir(r.Path(), to, []string{mod, modPath}, ignores)
}
