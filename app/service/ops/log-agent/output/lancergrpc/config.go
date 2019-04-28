package lancergrpc

import (
	"errors"
	"time"
	"go-common/app/service/ops/log-agent/output/cache/file"
	streamEvent "go-common/app/service/ops/log-agent/output/lancergrpc/lancergateway"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Local                bool                `tome:"local"`
	Name                 string              `tome:"name"`
	AggrSize             int                 `tome:"aggrSize"`
	SendConcurrency      int                 `tome:"sendConcurrency"`
	CacheConfig          *file.Config        `tome:"cacheConfig"`
	LancerGateway        *streamEvent.Config `tome:"lancerGateway"`
	SendBatchSize        int                 `tome:"sendBatchSize"`
	SendBatchNum         int                 `tome:"sendBatchNum"`
	SendBatchTimeout     xtime.Duration      `tome:"sendBatchTimeout"`
	SendFlushInterval    xtime.Duration      `tome:"sendFlushInterval"`
	InitialRetryDuration xtime.Duration      `tome:"initialRetryDuration"`
	MaxRetryDuration     xtime.Duration      `tome:"maxRetryDuration"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Lancer Output is nil")
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

	if c.SendFlushInterval == 0 {
		c.SendFlushInterval = xtime.Duration(time.Second * 5)
	}

	if c.InitialRetryDuration == 0 {
		c.InitialRetryDuration = xtime.Duration(time.Millisecond * 200)
	}

	if c.MaxRetryDuration == 0 {
		c.MaxRetryDuration = xtime.Duration(time.Second * 2)
	}

	if c.SendBatchNum == 0 {
		c.SendBatchNum = 3000
	}

	if c.SendBatchSize == 0 {
		c.SendBatchSize = 1024 * 1024 * 10
	}

	if c.SendBatchTimeout == 0 {
		c.SendBatchTimeout = xtime.Duration(time.Second * 5)
	}

	if c.LancerGateway == nil {
		c.LancerGateway = &streamEvent.Config{}
	}

	if err := c.LancerGateway.ConfigValidate(); err != nil {
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
