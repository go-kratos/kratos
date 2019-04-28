package http

import (
	"strconv"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func atoi(param string) int {
	if param == "" {
		return 0
	}
	res, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return res
}

func mDetailParam(c *bm.Context, param *model.MediaParam) (err error) {
	var msg string
	if err = c.Bind(param); err != nil {
		return
	}
	if param.EpID == 0 && param.SeasonID == 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if param.EpID > 0 {
		if param.SeasonID, msg, err = pgcSvc.EpControl(c, param.EpID); err != nil {
			c.JSONMap(map[string]interface{}{"message": msg}, err)
			return
		}
	}
	return
}

// get ep/season detail
func mediaDetail(c *bm.Context) {
	var param = new(model.MediaParam)
	if err := mDetailParam(c, param); err != nil {
		return
	}
	detail, msg, err := pgcSvc.SnDetail(c, param)
	if err != nil {
		c.JSONMap(map[string]interface{}{"message": msg}, err)
		return
	}
	c.JSONMap(map[string]interface{}{"result": detail, "message": "success"}, nil)
}

func mDetailV2(c *bm.Context) {
	var param = new(model.MediaParam)
	if err := mDetailParam(c, param); err != nil {
		return
	}
	detail, msg, err := pgcSvc.SnDetailV2(c, param)
	if err != nil {
		c.JSONMap(map[string]interface{}{"message": msg}, err)
		return
	}
	c.JSONMap(map[string]interface{}{"result": detail, "message": "success"}, nil)
}
