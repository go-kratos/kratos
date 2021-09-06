package httputil

import (
	"testing"
)

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

func TestContentType(t *testing.T) {
	tests := []struct {
		name    string
		subtype string
		want    string
	}{
		{"kratos", "kratos", "application/kratos"},
		{"json", "json", "application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentType(tt.subtype); got != tt.want {
				t.Errorf("ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
