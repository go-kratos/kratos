package run

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "absolute path is unchanged",
			input: "/absolute/path",
			want:  "/absolute/path",
		},
		{
			name:  "relative path becomes absolute",
			input: "relative/path",
			want:  filepath.Join(cwd, "relative/path"),
		},
		{
			name:  "dot-relative path becomes absolute",
			input: "./cmd/foo",
			want:  filepath.Join(cwd, "cmd/foo"),
		},
		{
			name:  "parent-relative path becomes absolute",
			input: "../cmd/foo",
			want:  filepath.Join(cwd, "../cmd/foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveDir(tt.input)
			if err != nil {
				t.Fatalf("resolveDir(%q) error = %v", tt.input, err)
			}

			// Clean both paths so the comparison is canonical.
			want := filepath.Clean(tt.want)
			got = filepath.Clean(got)

			if got != want {
				t.Errorf("resolveDir(%q) = %q, want %q", tt.input, got, want)
			}

			if !filepath.IsAbs(got) {
				t.Errorf("resolveDir(%q) = %q, expected absolute path", tt.input, got)
			}
		})
	}
}
