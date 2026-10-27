package run

import (
	"path/filepath"
	"testing"
)

func TestAbsPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantAbs bool
	}{
		{
			name:    "absolute path",
			input:   "/absolute/path",
			wantAbs: true,
		},
		{
			name:    "relative path",
			input:   "relative/path",
			wantAbs: false,
		},
		{
			name:    "relative path with dot",
			input:   "./cmd/foo",
			wantAbs: false,
		},
		{
			name:    "relative path with parent",
			input:   "../cmd/foo",
			wantAbs: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAbs := filepath.IsAbs(tt.input)
			if isAbs != tt.wantAbs {
				t.Errorf("filepath.IsAbs(%q) = %v, want %v", tt.input, isAbs, tt.wantAbs)
			}

			if !isAbs {
				abs, err := filepath.Abs(tt.input)
				if err != nil {
					t.Fatalf("filepath.Abs(%q) error = %v", tt.input, err)
				}
				if !filepath.IsAbs(abs) {
					t.Errorf("filepath.Abs(%q) = %q, not absolute", tt.input, abs)
				}
			}
		})
	}
}
