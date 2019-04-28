package util

import (
	"flag"
	"fmt"
	"strings"
)

// ParseParamSetFlag parse param from a=b,c=d
func ParseParamSetFlag(param string, fset *flag.FlagSet) (err error) {
	if param == "" {
		return nil
	}
	args := strings.Split(param, ",")
	for _, arg := range args {
		spec := strings.SplitN(arg, "=", 2)
		if len(spec) == 2 {
			err = fset.Set(spec[0], spec[1])
		} else {
			err = fset.Set(spec[0], "")
		}
		if err != nil {
			return fmt.Errorf("set flag error: %s", err)
		}
	}
	return nil
}
