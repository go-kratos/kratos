package redis

import (
	"strings"

	pkgerr "github.com/pkg/errors"
)

func formatErr(err error) string {
	e := pkgerr.Cause(err)
	switch e {
	case ErrNil, nil:
		return ""
	default:
		es := e.Error()
		switch {
		case strings.HasPrefix(es, "read"):
			return "read timeout"
		case strings.HasPrefix(es, "dial"):
			return "dial timeout"
		case strings.HasPrefix(es, "write"):
			return "write timeout"
		case strings.Contains(es, "EOF"):
			return "eof"
		case strings.Contains(es, "reset"):
			return "reset"
		case strings.Contains(es, "broken"):
			return "broken pipe"
		default:
			return "unexpected err"
		}
	}
}
