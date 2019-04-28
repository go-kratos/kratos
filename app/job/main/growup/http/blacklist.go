package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func initBlacklistMID(c *bm.Context) {
	log.Info("begin init blacklist mid")
	err := svr.InitBlacklistMID(c)
	if err != nil {
		log.Error("job growup svr.InitBlacklistMID err(%v)", err)
	}
	c.JSON(nil, err)
}

func updateBlacklist(c *bm.Context) {
	log.Info("begin update blacklist")
	err := svr.UpdateBlacklist(c)
	if err != nil {
		log.Error("(job growup svr.UpdateBlacklist error(%v)", err)
	}
	c.JSON(nil, err)
}
