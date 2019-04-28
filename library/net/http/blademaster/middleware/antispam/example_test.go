package antispam_test

import (
	"time"

	"go-common/library/cache/redis"
	"go-common/library/container/pool"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	xtime "go-common/library/time"
)

// This example create a antispam middleware instance and attach to a blademaster engine,
// it will protect '/ping' API with specified policy.
// If anyone who requests this API more frequently than 1 req/second or 1 req/hour,
// a StatusServiceUnavailable error will be raised.
func Example() {
	anti := antispam.New(&antispam.Config{
		On:     true,
		Second: 1,
		N:      1,
		Hour:   1,
		M:      1,
		Redis: &redis.Config{
			Config: &pool.Config{
				Active:      10,
				Idle:        10,
				IdleTimeout: xtime.Duration(time.Second * 60),
			},
			Name:         "test",
			Proto:        "tcp",
			Addr:         "172.18.33.60:6889",
			DialTimeout:  xtime.Duration(time.Second),
			ReadTimeout:  xtime.Duration(time.Second),
			WriteTimeout: xtime.Duration(time.Second),
		},
	})

	engine := blademaster.Default()
	engine.Use(anti)
	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":18080")
}
