package flagvar

import (
	"strings"
)

// StringVars []string implement flag.Value
type StringVars []string

func (s StringVars) String() string {
	return strings.Join(s, ",")
}

// Set implement flag.Value
func (s *StringVars) Set(val string) error {
	*s = append(*s, val)
	return nil
}
