package httpstream

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return fmt.Errorf("Error can't be nil")
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