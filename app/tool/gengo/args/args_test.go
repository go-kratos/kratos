package args

import (
	"testing"

	"go-common/app/tool/gengo/types"
)

func TestInputIncludes(t *testing.T) {
	a := &GeneratorArgs{
		InputDirs: []string{"a/b/..."},
	}
	if !a.InputIncludes(&types.Package{Path: "a/b/c"}) {
		t.Errorf("Expected /... syntax to work")
	}
	if a.InputIncludes(&types.Package{Path: "a/c/b"}) {
		t.Errorf("Expected correctness")
	}
}
