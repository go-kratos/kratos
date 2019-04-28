package verify_test

import (
	"fmt"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

// This example create a identify middleware instance and attach to several path,
// it will validate request by specified policy and put extra information into context. e.g., `mid`.
// It provides additional handler functions to provide the identification for your business handler.
func Example() {
	idt := verify.New(&verify.Config{
		OpenServiceHost: "http://uat-open.bilibili.co",
		HTTPClient: &bm.ClientConfig{
			App: &bm.App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 10),
		},
	})

	e := bm.Default()
	// mark `/verify` path as Verify policy
	e.GET("/verify", idt.Verify, func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	// mark `/verify` path as VerifyUser policy
	e.GET("/verifyUser", idt.VerifyUser, func(c *bm.Context) {
		mid := metadata.Int64(c, metadata.Mid)
		c.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.Run(":18080")
}
