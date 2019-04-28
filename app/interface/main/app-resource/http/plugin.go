package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func plugin(c *bm.Context) {
	var (
		params   = c.Request.Form
		build    int
		baseCode int
		seed     int
		err      error
	)
	name := params.Get("name")
	buildStr := params.Get("build")
	baseCodeStr := params.Get("base_code")
	seedStr := params.Get("seed")
	if build, err = strconv.Atoi(buildStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if seed, err = strconv.Atoi(seedStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if baseCode, err = strconv.Atoi(baseCodeStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pg := pgSvr.Plugin(build, baseCode, seed, name); pg != nil {
		c.JSON(pg, nil)
	} else {
		c.JSON(struct{}{}, nil)
	}

}
