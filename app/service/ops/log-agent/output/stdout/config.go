package stdout

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name string `tome:"name"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Stdout Output is nil")
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
