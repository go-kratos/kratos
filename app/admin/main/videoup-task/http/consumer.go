package http

import (
	"net/http"
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

// 登录管理
func on(c *bm.Context) {
	uid, uname := getUIDName(c)
	err := srv.HandsUp(c, uid, uname)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

// 踢出
func forceoff(c *bm.Context) {
	uidS := c.Request.Form.Get("uid")
	uid, _ := strconv.ParseInt(uidS, 10, 64)
	if uid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	adminuid, _ := getUIDName(c)
	err := srv.HandsOff(c, adminuid, uid)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

func off(c *bm.Context) {
	adminuid, _ := getUIDName(c)
	err := srv.HandsOff(c, adminuid, 0)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

func online(c *bm.Context) {
	c.JSON(srv.Online(c))
}

func inoutlist(c *bm.Context) {
	v := new(struct {
		Unames string `form:"unames" default:""`
		Bt     string `form:"bt" default:""`
		Et     string `form:"et" default:""`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.InOutList(c, v.Unames, v.Bt, v.Et))
}
