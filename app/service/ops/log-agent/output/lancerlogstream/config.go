package lancerlogstream

import (
	"errors"
	"go-common/app/service/ops/log-agent/output/cache/file"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Local           bool            `tome:"local"`
	Name            string          `tome:"name"`
	AggrSize        int             `tome:"aggrSize"`
	SendConcurrency int             `tome:"sendConcurrency"`
	CacheConfig     *file.Config    `tome:"cacheConfig"`
	PoolConfig      *ConnPoolConfig `tome:"poolConfig"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Sock Input is nil")
	}

	if c.Name == "" {
		return errors.New("output Name can't be nil")
	}

	if c.AggrSize == 0 {
		c.AggrSize = 819200
	}

	if c.SendConcurrency == 0 {
		c.SendConcurrency = 5
	}

	if err := c.CacheConfig.ConfigValidate(); err != nil {
		return err
	}

	if err := c.PoolConfig.ConfigValidate(); err != nil {
		return err
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
