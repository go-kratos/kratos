package http

import (
	"fmt"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

// cutStream 切流， 内部调用
func cutStream(c *bm.Context) {
	// roomid 必须
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	cutTime := params.Get("cut_time")

	c.Set("input_params", params)

	// 验证传参数
	roomID, err := strconv.ParseInt(room, 10, 64)
	if err != nil || roomID <= 0 {
		c.Set("output_data", "roomid is not right")
		c.JSONMap(map[string]interface{}{"message": "roomid is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 默认切流时间为1s,可以传入-1
	ct, err := strconv.ParseInt(cutTime, 10, 64)
	if err != nil || ct == 0 {
		ct = 1
	}

	srv.StreamCut(c, roomID, ct, 0)

	c.Set("output_data", "ok")
	c.JSONMap(map[string]interface{}{"data": map[string]int{}}, nil)
}

// cutStreamByExt 外部调用
func cutStreamByMobile(c *bm.Context) {
	// roomid 必须
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	cutTime := params.Get("cut_time")

	c.Set("input_params", params)

	// 验证传参数
	roomID, err := strconv.ParseInt(room, 10, 64)
	if err != nil || roomID <= 0 {
		c.Set("output_data", "roomid is not right")
		c.JSONMap(map[string]interface{}{"message": "roomid is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 默认切流时间为1s,可以传入-1
	ct, err := strconv.ParseInt(cutTime, 10, 64)
	if err != nil || ct == 0 {
		ct = 1
	}

	uid, ok := metadata.Value(c, metadata.Mid).(int64)
	//uid = 19148701
	//ok = true
	if !ok {
		c.Set("output_data", "未登陆")
		c.JSONMap(map[string]interface{}{"message": fmt.Sprintf("未登陆")}, ecode.RequestErr)
		c.Abort()
		return
	}

	err = srv.StreamCut(c, roomID, ct, uid)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", "ok")
	c.JSONMap(map[string]interface{}{"data": map[string]int{}}, nil)
}
