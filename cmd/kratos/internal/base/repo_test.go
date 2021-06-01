package base

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	r := NewRepo("https://github.com/go-kratos/service-layout.git", "")
	if err := r.Clone(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "/tmp/test_repo", "github.com/go-kratos/kratos-layout", nil); err != nil {
		t.Fatal(err)
	}
}
