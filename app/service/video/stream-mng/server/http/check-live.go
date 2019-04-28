package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
)

func checkLiveStreamList(c *bm.Context) {
	params := c.Request.URL.Query()
	rooms := params.Get("room_ids")

	c.Set("input_params", params)

	if rooms == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 切割room_id
	roomIDs := strings.Split(rooms, ",")

	if len(roomIDs) <= 0 {
		c.Set("output_data", "room_ids is not right")
		c.JSONMap(map[string]interface{}{"message": "room_ids is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	rids := []int64{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}

		rids = append(rids, roomID)
	}

	c.JSONMap(map[string]interface{}{"data": srv.CheckLiveStreamList(c, rids)}, nil)
}
