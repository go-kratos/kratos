package file

import (
	"errors"
	"regexp"
	"fmt"
)

type MultilineConf struct {
	Pattern  string `toml:"pattern"`
	MaxLines int    `toml:"maxLines"`
}

func (c *MultilineConf) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Multiline  is nil")
	}

	if c.Pattern == "" {
		return errors.New("Pattern in Multiline can't be nil")
	}

	if _, err := regexp.Compile(c.Pattern); err != nil {
		return fmt.Errorf("Multiline pattern compile error: %s", err)
	}

	if c.MaxLines == 0 {
		c.MaxLines = 200
	}
	return nil
}
