package httputil

import (
	"strings"
)

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return baseContentType + "/" + subtype
}

// ContentSubtype extracts and returns the content subtype from a given Content-Type string.
// The input is expected to be lowercase, following the conventions of RFC 7231.
// It handles formats like "type/subtype" or "type/subtype; key=value".
// If the input is not well-formed, an empty string is returned.
func ContentSubtype(contentType string) string {
	switch contentType {
	case "":
		return ""
	case "application/json":
		return "json"
	}

	left := strings.IndexByte(contentType, '/')
	if left == -1 {
		return ""
	}
	right := strings.IndexByte(contentType, ';')
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
