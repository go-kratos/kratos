package dsn

import (
	"errors"
	"strings"
)

var (
	errInvalidDSN = errors.New("invalid dsn params")
)

// DSN is a configuration parsed from a DSN string
// key:secret@group/topic=?&role=?
type DSN struct {
	Key    string // app key
	Secret string // app secret
	Group  string // kafka group
	Topic  string // kafka topic
	Role   string // pub or sub
	Color  string // env color
}

// ParseDSN parse databus info.
func ParseDSN(s string) (*DSN, error) {
	if strings.Count(s, "@") != 1 || strings.Count(s, "/") != 1 || strings.Count(s, ":") != 1 {
		return nil, errInvalidDSN
	}
	var (
		c      = &DSN{}
		params string
	)
	i := len(s) - 1
	var j, k int
	for j = i; j >= 0; j-- {
		// found key:passwd
		if s[j] == '@' {
			for k = 0; k < j; k++ {
				if s[k] == ':' {
					c.Secret = s[k+1 : j]
					break
				}
			}
			c.Key = s[:k]
			break
		}
	}
	// group
	for k = j + 1; k < i; k++ {
		if s[k] == '/' {
			break
		}
	}
	c.Group = s[j+1 : k]
	params = s[k+1:]
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}
		switch value := param[1]; strings.ToLower(param[0]) {
		case "topic":
			c.Topic = value
		case "role":
			c.Role = value
		case "color":
			c.Color = value
		}
	}
	return c, nil
}
