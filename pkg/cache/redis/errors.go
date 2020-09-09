package redis

import (
	"strings"

	pkgerr "github.com/pkg/errors"
)

func formatErr(err error, name, addr string) string {
	e := pkgerr.Cause(err)
	switch e {
	case ErrNil, nil:
		if e == ErrNil {
			_metricMisses.Inc(name, addr)
		}
		return ""
	default:
		es := e.Error()
		switch {
		case strings.HasPrefix(es, "read"):
			return "read timeout"
		case strings.HasPrefix(es, "dial"):
			if strings.Contains(es, "connection refused") {
				return "connection refused"
			}
			return "dial timeout"
		case strings.HasPrefix(es, "write"):
			return "write timeout"
		case strings.Contains(es, "EOF"):
			return "eof"
		case strings.Contains(es, "reset"):
			return "reset"
		case strings.Contains(es, "broken"):
			return "broken pipe"
		case strings.Contains(es, "pool exhausted"):
			return "pool exhausted"
		case strings.Contains(es, "pool closed"):
			return "pool closed"
		default:
			return "unexpected err"
		}
	}
}
