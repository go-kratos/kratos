package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

// clearRoomCacheByRID删除room_id缓存的接口，防止缓存问题出现的bug
func clearRoomCacheByRID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is empty")
		c.JSONMap(map[string]interface{}{"message": "room_id is empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	err = srv.ClearRoomCacheByRID(c, roomID)

	if err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"message": "ok"}, nil)
}
