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
			input:  "a",
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
		{
			input:  "a.b.c",
			expect: "c",
		},
	}
	for _, v := range tests {
		content := format(v.input)
		if got, want := content, v.expect; got != want {
			t.Errorf("expect %v,got %v", want, got)
		}
	}
}

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		format("abc.txt")
	}
}
