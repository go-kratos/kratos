package args

import (
	"strings"
)

// StringSliceVar is
type StringSliceVar []string

func (v StringSliceVar) String() string {
	return strings.Join(v, ",")
}

// Set is
func (v *StringSliceVar) Set(in string) error {
	for _, a := range strings.Split(in, ",") {
		*v = append(*v, a)
	}
	return nil
}
