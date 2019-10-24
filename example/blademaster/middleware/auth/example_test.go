package auth_test

import (
	"fmt"

	"github.com/bilibili/kratos/example/blademaster/middleware/auth"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/metadata"
)

// This example create a identify middleware instance and attach to several path,
// it will validate request by specified policy and put extra information into context. e.g., `mid`.
// It provides additional handler functions to provide the identification for your business handler.
func Example() {
	myHandler := func(ctx *bm.Context) {
		mid := metadata.Int64(ctx, metadata.Mid)
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	}

	authn := auth.New(&auth.Config{
		DisableCSRF: false,
	})

	e := bm.DefaultServer(nil)

	// mark `/user` path as User policy
	e.GET("/user", authn.User, myHandler)
	// mark `/mobile` path as UserMobile policy
	e.GET("/mobile", authn.UserMobile, myHandler)
	// mark `/web` path as UserWeb policy
	e.GET("/web", authn.UserWeb, myHandler)
	// mark `/guest` path as Guest policy
	e.GET("/guest", authn.Guest, myHandler)

	o := e.Group("/owner", authn.User)
	o.GET("/info", myHandler)
	o.POST("/modify", myHandler)

	e.Start()
}
