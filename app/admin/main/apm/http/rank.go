package http

import (
	"go-common/app/admin/main/apm/model/ut"
	bm "go-common/library/net/http/blademaster"
)

// @params RankReq
// @router get /x/admin/apm/ut/rank/list
// @response RankResp
func utRank(c *bm.Context) {
	var (
		err        error
		topTens    []*ut.RankResp
		bottomTens []*ut.RankResp
		rankList   = make(map[string]interface{})
		UserList   *ut.RankResp
	)
	v := new(struct {
		UserName string `form:"username"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if topTens, err = apmSvc.RankTen(c, "desc"); err != nil {
		c.JSON(nil, err)
		return
	}
	if bottomTens, err = apmSvc.RankTen(c, "asc"); err != nil {
		c.JSON(nil, err)
		return
	}
	if UserList, err = apmSvc.UserRank(c, v.UserName); err != nil {
		c.JSON(nil, err)
		return
	}
	rankList["topten"] = topTens
	rankList["bottomten"] = bottomTens
	rankList["user"] = UserList
	c.JSON(rankList, nil)
}

// /x/admin/apm/ut/rank/user
func userRank(c *bm.Context) {
	var (
		err      error
		UserList *ut.RankResp
	)

	v := new(struct {
		UserName string `form:"username"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if UserList, err = apmSvc.UserRank(c, v.UserName); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(UserList, nil)
}
