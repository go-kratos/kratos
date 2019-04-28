package middleware

import (
	"encoding/json"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func Logger() bm.HandlerFunc {
	return func(c *bm.Context) {
		// panic("sss")
		c.Next()
		i, _ := c.Get("input_params")
		ji, _ := json.Marshal(i)

		o, _ := c.Get("output_data")
		jo, _ := json.Marshal(o)

		log.Infov(c,
			log.KV("path", c.Request.URL.Path),
			log.KV("method", c.Request.Method),
			log.KV("input_params", string(ji)),
			log.KV("output_data", string(jo)),
		)
	}
}
