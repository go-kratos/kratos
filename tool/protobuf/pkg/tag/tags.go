package tag

import (
	"reflect"
	"strings"
)

// GetCommentWithoutTag strip tags in comment
func GetCommentWithoutTag(comment string) []string {
	var lines []string
	if comment == "" {
		return lines
	}
	split := strings.Split(strings.TrimRight(comment, "\n\r"), "\n")
	for _, line := range split {
		tag, _, _ := GetLineTag(line)
		if tag == "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func GetTagsInComment(comment string) []reflect.StructTag {
	split := strings.Split(comment, "\n")
	var tagsInComment []reflect.StructTag
	for _, line := range split {
		tag, _, _ := GetLineTag(line)
		if tag != "" {
			tagsInComment = append(tagsInComment, tag)
		}
	}
	return tagsInComment
}

func GetTagValue(key string, tags []reflect.StructTag) string {
	for _, t := range tags {
		val := t.Get(key)
		if val != "" {
			return val
		}
	}
	return ""
}

// find tag between backtick, start & end is the position of backtick
func GetLineTag(line string) (tag reflect.StructTag, start int, end int) {
	start = strings.Index(line, "`")
	end = strings.LastIndex(line, "`")
	if end <= start {
		return
	}
	tag = reflect.StructTag(line[start+1 : end])
	return
}
