package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// VideoList get ugc video list
func VideoList(c *bm.Context) {
	var (
		res   = map[string]interface{}{}
		err   error
		pager *model.VideoListPager
	)
	param := new(model.VideoListParam)
	if err = c.Bind(param); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pager, err = tvSrv.VideoList(c, param); err != nil {
		res["message"] = "获取数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

// VideoOnline online ugc video
func VideoOnline(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if err := tvSrv.VideoOnline(param.IDs); err != nil {
		res["message"] = "更新数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON("成功", nil)
}

// VideoHidden hidden ugc video
func VideoHidden(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if err := tvSrv.VideoHidden(param.IDs); err != nil {
		res["message"] = "更新数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON("成功", nil)
}

// VideoPreview ugc video preview
func VideoPreview(c *bm.Context) {
	var (
		playurl string
		err     error
		res     = map[string]interface{}{}
	)
	param := new(struct {
		ID int `form:"id" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if playurl, err = tvSrv.UPlayurl(param.ID); err != nil {
		res["message"] = "获取playurl失败" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(playurl, nil)
}

// videoUpdate ugc video update
func videoUpdate(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(struct {
		ID    int    `form:"id" validate:"required"`
		Title string `form:"title" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if err = tvSrv.VideoUpdate(param.ID, param.Title); err != nil {
		res["message"] = "更新失败" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func abnorExport(c *bm.Context) {
	buf, fileName := tvSrv.AbnormExport()
	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fileName)
	c.Writer.Write(buf.Bytes())
}

func abnorDebug(c *bm.Context) {
	c.JSON(tvSrv.AbnDebug(), nil)
}
