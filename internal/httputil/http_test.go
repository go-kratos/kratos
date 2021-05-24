package httputil

import "testing"

func TestContentSubtype(t *testing.T) {
	tests := []struct {
		contentType string
		want        string
	}{
		{"text/html; charset=utf-8", "html"},
		{"multipart/form-data; boundary=something", "form-data"},
		{"application/json; charset=utf-8", "json"},
		{"application/json", "json"},
		{"application/xml", "xml"},
		{"text/xml", "xml"},
		{";text/xml", ""},
		{"application", ""},
	}
	for _, test := range tests {
		t.Run(test.contentType, func(t *testing.T) {
			got := ContentSubtype(test.contentType)
			if got != test.want {
				t.Fatalf("want %v got %v", test.want, got)
			}
		})
	}
}
