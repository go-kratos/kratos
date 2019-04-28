package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

// webAdGameList fn
func webAdGameList(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn <= 0 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 0 || ps > 20 {
		ps = 20
	}
	keywordStr := params.Get("keyword")
	letterStr := params.Get("letter")
	data, err := adSvc.GameList(c, keywordStr, letterStr, pn, ps, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
