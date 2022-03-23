package base

import (
	"context"
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	os.RemoveAll("/tmp/test_repo")
	r := NewRepo("https://github.com/go-kratos/service-layout.git", "")
	if err := r.Clone(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "/tmp/test_repo", "github.com/go-kratos/kratos-layout", nil); err != nil {
		t.Fatal(err)
	}
	urls := []string{
		"ssh://git@gitlab.xxx.com:1234/foo/bar.git",
		"ssh://gitlab.xxx.com:1234/foo/bar.git",
		"//git@gitlab.xxx.com:1234/foo/bar.git",
		"git@gitlab.xxx.com:1234/foo/bar.git",
		"gitlab.xxx.com:1234/foo/bar.git",
		"gitlab.xxx.com/foo/bar.git",
		"gitlab.xxx.com/foo/bar",
	}
	for _, url := range urls {
		dir := repoDir(url)
		if dir != "gitlab.xxx.com/foo" {
			t.Fatal("repoDir test failed", dir)
		}
	}
}
