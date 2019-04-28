package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
)

// getStreamLastTime 得到流到最后推流时间
func saveMaskByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	roomid := params.Get("room_id")
	mask := params.Get("mask")
	c.Set("input_data", params)
	roomID, err := strconv.ParseInt(roomid, 10, 64)
	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}
	int64mask, err := strconv.ParseInt(mask, 10, 64)
	// 验证传参数
	if err != nil || int64mask < 0 || int64mask > 1 {
		c.Set("output_data", "mask is not right")
		c.JSONMap(map[string]interface{}{"message": "mask is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	//直接修改数据库，更新缓存
	result, err := srv.ChangeMaskStreamByRoomID(c, roomID, "", int64mask)
	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", result)
	c.JSONMap(map[string]interface{}{"data": result}, nil)
}

// getStreamLastTime 得到流到最后推流时间
func saveMaskByStreamName(c *bm.Context) {
	params := c.Request.URL.Query()
	sname := params.Get("stream_name")
	mask := params.Get("mask")
	c.Set("input_data", params)
	if sname == "" {
		c.Set("output_data", "stream_name is not right")
		c.JSONMap(map[string]interface{}{"message": "stream_name is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	int64mask, err := strconv.ParseInt(mask, 10, 64)
	// 验证传参数
	if err != nil || int64mask < 0 || int64mask > 1 {
		c.Set("output_data", "mask is not right")
		c.JSONMap(map[string]interface{}{"message": "mask is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	var newmask int64
	var newsname string
	if len(sname) > 6 && strings.Contains(sname, "_wmask") {
		//设置第三位为1
		if int64mask == 0 {
			newmask = 3
		} else {
			newmask = 2
		}
		newsname = sname[0 : len(sname)-6]
	} else if len(sname) > 6 && strings.Contains(sname, "_mmask") {
		//设置第四位为1
		if int64mask == 0 {
			newmask = 5
		} else {
			newmask = 4
		}
		newsname = sname[0 : len(sname)-6]
	} else {
		c.Set("output_data", "stream_name is not right")
		c.JSONMap(map[string]interface{}{"message": "stream_name is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	//直接修改数据库，更新缓存
	result, err := srv.ChangeMaskStreamByRoomID(c, 0, newsname, newmask)
	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", result)
	c.JSONMap(map[string]interface{}{"data": result}, nil)
}
