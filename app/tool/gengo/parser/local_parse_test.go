package parser

import (
	"testing"
)

func TestImportBuildPackage(t *testing.T) {
	b := New()
	if _, err := b.importBuildPackage("go-common/app/tool/gengo/testdata/fake/dep"); err != nil {
		t.Fatal(err)
	}
	if _, ok := b.buildPackages["go-common/app/tool/gengo/testdata/fake/dep"]; !ok {
		t.Errorf("missing expected, but got %v", b.buildPackages)
	}

	if len(b.buildPackages) > 1 {
		// this would happen if the canonicalization failed to normalize the path
		// you'd get a go-common/app/tool/gengo/testdata/fake/dep key too
		t.Errorf("missing one, but got %v", b.buildPackages)
	}
}

func TestCanonicalizeImportPath(t *testing.T) {
	tcs := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "passthrough",
			input:  "github.com/foo/bar",
			output: "github.com/foo/bar",
		},
		{
			name:   "simple",
			input:  "github.com/foo/vendor/k8s.io/kubernetes/pkg/api",
			output: "k8s.io/kubernetes/pkg/api",
		},
		{
			name:   "deeper",
			input:  "github.com/foo/bar/vendor/k8s.io/kubernetes/pkg/api",
			output: "k8s.io/kubernetes/pkg/api",
		},
	}

	for _, tc := range tcs {
		actual := canonicalizeImportPath(tc.input)
		if string(actual) != tc.output {
			t.Errorf("%v: expected %q got %q", tc.name, tc.output, actual)
		}
	}
}
