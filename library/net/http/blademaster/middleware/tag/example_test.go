package tag_test

import (
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/tag"
)

// This example create a tag middleware instance and attach to a global.
// It will put a tag into Keys field of context by specified policy custom defines
// You can define a custom policy through any request param or http header.
// Register several tag middlewares to put several tags.
func Example() {
	var pf tag.PolicyFunc
	// create your tag policy
	pf = func(ctx *blademaster.Context) string {
		if ctx.Request.Form.Get("group") == "a" {
			return "a"
		}
		return "b"
	}
	t := tag.New("abtest", pf)

	engine := blademaster.Default()
	engine.Use(t)
	engine.GET("/abtest", HandlerMap)

	engine.Run(":18080")
}

func HandlerMap(ctx *blademaster.Context) {
	value, ok := tag.Value(ctx, "abtest")
	if !ok {
		ctx.String(-400, "failed to parse group")
		ctx.Abort()
		return
	}

	if value == "a" {
		HandlerA(ctx)
	}
	if value == "b" {
		HandlerB(ctx)
	}
}

func HandlerA(ctx *blademaster.Context) {
	// your business
	ctx.String(200, "group a")
	return
}

func HandlerB(ctx *blademaster.Context) {
	// your business
	ctx.String(200, "group b")
	return
}
