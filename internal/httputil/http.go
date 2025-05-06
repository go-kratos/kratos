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

// ContentSubtype returns the content-subtype for the given content-type.
// The contentType is assumed to be lowercase, as per RFC7231.
// The function extracts the subtype from the content-type in the format: "<type>/<subtype>",
// and returns an empty string if the content-type is not well-formed or lacks a subtype.
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
