package base

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Repo is git repository manager.
type Repo struct {
	url    string
	home   string
	branch string
}

// NewRepo new a repository manager.
func NewRepo(url string, branch string) *Repo {
	var start int
	start = strings.Index(url, "//")
	if start == -1 {
		start = strings.Index(url, ":") + 1
	} else {
		start += 2
	}
	end := strings.LastIndex(url, "/")
	return &Repo{
		url:    url,
		home:   kratosHomeWithDir("repo/" + url[start:end]),
		branch: branch,
	}
}

// Path returns the repository cache path.
func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	branch := ""
	if r.branch == "" {
		branch = "@main"
	} else {
		branch = "@" + r.branch
	}
	return path.Join(r.home, r.url[start+1:end]+branch)
}

// Pull fetch the repository from remote url.
func (r *Repo) Pull(ctx context.Context) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = r.Path()
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if strings.Contains(out.String(), "You are not currently on a branch.") {
		return nil
	}
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx)
	}
	cmd := &exec.Cmd{}
	if r.branch == "" {
		cmd = exec.Command("git", "clone", r.url, r.Path())
	} else {
		cmd = exec.Command("git", "clone", "-b", r.branch, r.url, r.Path())
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
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
