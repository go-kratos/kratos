package http

import "testing"

func TestSubtype(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"application/json", "json"},
		{"application/json;", "json"},
		{"application/json; charset=utf-8", "json"},
		{"application/", ""},
		{"application", ""},
		{"foo", ""},
		{"", ""},
	}
	for _, test := range tests {
		if contentSubtype(test.input) != test.expected {
			t.Errorf("expected %s got %s", test.expected, test.input)
		}
	}
}
