package base

import (
	"context"
	"fmt"
	stdurl "net/url"
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

func repoDir(url string) string {
	if !strings.Contains(url, "//") {
		url = "//" + url
	}
	if strings.HasPrefix(url, "//git@") {
		url = "ssh:" + url
	} else if strings.HasPrefix(url, "//") {
		url = "https:" + url
	}
	u, err := stdurl.Parse(url)
	if err == nil {
		url = fmt.Sprintf("%s://%s%s", u.Scheme, u.Hostname(), u.Path)
	}
	var start int
	start = strings.Index(url, "//")
	if start == -1 {
		start = strings.Index(url, ":") + 1
	} else {
		start += 2
	}
	end := strings.LastIndex(url, "/")
	return url[start:end]
}

// NewRepo new a repository manager.
func NewRepo(url string, branch string) *Repo {
	return &Repo{
		url:    url,
		home:   kratosHomeWithDir("repo/" + repoDir(url)),
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
	var branch string
	if r.branch == "" {
		branch = "@main"
	} else {
		branch = "@" + r.branch
	}
	return path.Join(r.home, r.url[start+1:end]+branch)
}

// Pull fetch the repository from remote url.
func (r *Repo) Pull(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "HEAD")
	cmd.Dir = r.Path()
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	cmd = exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.Path()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx)
	}
	var cmd *exec.Cmd
	if r.branch == "" {
		cmd = exec.CommandContext(ctx, "git", "clone", r.url, r.Path())
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", "-b", r.branch, r.url, r.Path())
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
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
