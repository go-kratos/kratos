package change

import "testing"

func TestParseGithubURL(t *testing.T) {
	urls := []struct {
		url   string
		owner string
		repo  string
	}{
		{"https://github.com/go-kratos/kratos.git", "go-kratos", "kratos"},
		{"https://github.com/go-kratos/kratos", "go-kratos", "kratos"},
		{"git@github.com:go-kratos/kratos.git", "go-kratos", "kratos"},
		{"https://github.com/go-kratos/go-kratos.dev.git", "go-kratos", "go-kratos.dev"},
	}
	for _, url := range urls {
		owner, repo := ParseGithubURL(url.url)
		if owner != url.owner {
			t.Fatalf("owner want: %s, got: %s", owner, url.owner)
		}
		if repo != url.repo {
			t.Fatalf("repo want: %s, got: %s", repo, url.repo)
		}
	}
}
