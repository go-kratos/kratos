package base

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	r := NewRepo()
	if err := r.Clone(context.Background(), "test", "https://github.com/golang-standards/project-layout.git"); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "test", "https://github.com/golang-standards/project-layout.git", "/tmp/test_repo"); err != nil {
		t.Fatal(err)
	}
}
