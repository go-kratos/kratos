package supervisor_test

import (
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/supervisor"
	"time"
)

// This example create a supervisor middleware instance and attach to a blademaster engine,
// it will allow '/ping' API can be requested with specified policy.
// This example will block all http method except `GET` on '/ping' API in next hour,
// and allow in further.
func Example() {
	now := time.Now()
	end := now.Add(time.Hour * 1)
	spv := supervisor.New(&supervisor.Config{
		On:    true,
		Begin: now,
		End:   end,
	})

	engine := blademaster.Default()
	engine.Use(spv)
	engine.GET("/ping", func(c *blademaster.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":18080")
}
