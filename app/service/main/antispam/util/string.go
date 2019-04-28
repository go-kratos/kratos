package util

import (
	"strconv"
	"strings"
)

// StrSliToSQLVarchars convert string slice to varchar in sql syntax
// eg: ["default", "deleted", "modified"] -> " 'default', 'deleted', 'modified' "
// so that we can used it in 'SELECT * ... WHERE xxx IN ('default', 'deleted', 'modified')'
func StrSliToSQLVarchars(s []string) string {
	sli := make([]string, len(s))
	for i, ss := range s {
		sli[i] = "'" + ss + "'"
	}
	return strings.Join(sli, ",")
}

// StrToIntSli convert string to int slice, eg: "1,2,3" -> [1,2,3]
func StrToIntSli(s string, delimiter string) ([]int64, error) {
	var result []int64
	sli := strings.Split(s, delimiter)
	for _, intStr := range sli {
		i, err := strconv.ParseInt(intStr, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

// IntSliToSQLVarchars convert int slice to string, eg: [1,2,3] -> "1,2,3"
func IntSliToSQLVarchars(ints []int64) string {
	return intSliToStr(ints, ",")
}

func intSliToStr(ints []int64, delimiter string) string {
	sli := make([]string, len(ints))
	for i, ii := range ints {
		sli[i] = strconv.FormatInt(ii, 10)
	}
	return strings.Join(sli, delimiter)
}

// SameChar check if string consists of same characters
func SameChar(content string) bool {
	content = strings.ToLower(content)
	first := content[0]
	for _, s := range content {
		if s != rune(first) {
			return false
		}
	}
	return true
}

// StripPrefix remove prefix from string if exists
func StripPrefix(s string, prefix string, suffix string) string {
	if strings.HasPrefix(s, prefix) {
		i := strings.Index(s, suffix)
		return s[i+1:]
	}
	return s
}
