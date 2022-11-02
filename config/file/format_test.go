package file

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "",
			expect: "",
		},
		{
			input:  " ",
			expect: "",
		},
		{
			input:  ".",
			expect: "",
		},
		{
			input:  "a.",
			expect: "",
		},
		{
			input:  ".b",
			expect: "b",
		},
		{
			input:  "a.b",
			expect: "b",
		},
	}
	for _, v := range tests {
		content := format(v.input)
		if got, want := content, v.expect; got != want {
			t.Errorf("expect %v,got %v", want, got)
		}
	}
}
