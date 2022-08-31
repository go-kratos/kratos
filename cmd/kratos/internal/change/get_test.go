package change

import "testing"

func TestParseGithubURL(t *testing.T) {
	urls := [][]string{
		{"https://github.com/go-kratos/kratos.git", "go-kratos", "kratos"},
		{"https://github.com/go-kratos/kratos", "go-kratos", "kratos"},
		{"git@github.com:go-kratos/kratos.git", "go-kratos", "kratos"},
		{"https://github.com/go-kratos/go-kratos.dev.git", "go-kratos", "go-kratos.dev"},
	}
	for _, url := range urls {
		owner, repo := ParseGithubURL(url[0])
		if owner != url[1] {
			t.Fatalf("owner want: %s, got: %s", owner, url[1])
		}
		if repo != url[2] {
			t.Fatalf("repo want: %s, got: %s", repo, url[2])
		}
	}
}
