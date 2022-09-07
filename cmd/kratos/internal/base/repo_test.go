package base

import (
	"context"
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	urls := []string{
		// ssh://[user@]host.xz[:port]/path/to/repo.git/
		"ssh://git@github.com:7875/go-kratos/kratos.git",
		// git://host.xz[:port]/path/to/repo.git/
		"git://github.com:7875/go-kratos/kratos.git",
		// http[s]://host.xz[:port]/path/to/repo.git/
		"https://github.com:7875/go-kratos/kratos.git",
		// ftp[s]://host.xz[:port]/path/to/repo.git/
		"ftps://github.com:7875/go-kratos/kratos.git",
		//[user@]host.xz:path/to/repo.git/
		"git@github.com:go-kratos/kratos.git",
		// ssh://[user@]host.xz[:port]/~[user]/path/to/repo.git/
		"ssh://git@github.com:7875/go-kratos/kratos.git",
		// git://host.xz[:port]/~[user]/path/to/repo.git/
		"git://github.com:7875/go-kratos/kratos.git",
		//[user@]host.xz:/~[user]/path/to/repo.git/
		"git@github.com:go-kratos/kratos.git",
		///path/to/repo.git/
		"//github.com/go-kratos/kratos.git",
		// file:///path/to/repo.git/
		"file://./github.com/go-kratos/kratos.git",
	}
	for _, url := range urls {
		dir := repoDir(url)
		if dir != "github.com/go-kratos" && dir != "/go-kratos" {
			t.Fatal(url, "repoDir test failed", dir)
		}
	}
}

func TestRepoClone(t *testing.T) {
	r := NewRepo("https://github.com/go-kratos/service-layout.git", "")
	if err := r.Clone(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "/tmp/test_repo", "github.com/go-kratos/kratos-layout", nil); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll("/tmp/test_repo")
	})
}
