package grok

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Pattern string `toml:"pattern"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return fmt.Errorf("config of grok can't be nil")
	}

	if c.Pattern == "" {
		return fmt.Errorf("pattern of grok can't be nil")
	}

	return nil
}

func DecodeConfig(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	c = new(Config)
	if err = md.PrimitiveDecode(primValue, c); err != nil {
		return nil, err
	}
	return c, nil
}
