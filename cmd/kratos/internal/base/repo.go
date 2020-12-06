package base

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

// Repo is git repository manager.
type Repo struct {
	Dir string
}

// NewRepo new a repository manager.
func NewRepo() *Repo {
	return &Repo{
		Dir: kratosHomeWithDir("repo"),
	}
}

func (r *Repo) Pull(ctx context.Context, name, url string) error {
	dir := path.Join(r.Dir, name)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	if err = w.PullContext(ctx, &git.PullOptions{RemoteName: "origin"}); errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

func (r *Repo) Clone(ctx context.Context, name, url string) error {
	dir := path.Join(r.Dir, name)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return r.Pull(ctx, name, url)
	}
	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func (r *Repo) CopyTo(ctx context.Context, name, url, to string) error {
	dir := path.Join(r.Dir, name)
	if err := r.Clone(ctx, name, url); err != nil {
		return err
	}
	return copyDir(dir, to, []string{".git"})
}
