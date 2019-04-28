package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// ping check server ok.
func ping(c *bm.Context) {
	var (
		err error
	)
	if err = confSvc.Ping(c); err != nil {
		log.Error("config service ping error(%v)", err)
		c.JSON(nil, err)
		http.Error(c.Writer, "", http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{}, nil)
}
