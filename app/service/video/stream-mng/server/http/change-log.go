package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

// getChangeLogByRoomID 查询cdn切换记录
func getChangeLogByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	limit := params.Get("limit")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 默认查询1最近一条记录
	limitInt, _ := strconv.ParseInt(limit, 10, 64)

	if limitInt <= 0 {
		limitInt = 1
	}

	infos, err := srv.GetChangeLogByRoomID(c, roomID, limitInt)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": infos}, nil)
}
