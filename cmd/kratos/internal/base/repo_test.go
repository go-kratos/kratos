package base

import (
	"testing"
)

func TestRepo(t *testing.T) {
	if err := Clone("https://github.com/go-kratos/service-layout.git", "template"); err != nil {
		t.Fatal(err)
	}
}
