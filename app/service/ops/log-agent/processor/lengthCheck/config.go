package lengthCheck

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	MaxLength int `toml:"maxLength"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return fmt.Errorf("Error can't be nil")
	}
	if c.MaxLength == 0 || c.MaxLength > 1024*32 {
		c.MaxLength = 1024 * 32 //32K by default
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
