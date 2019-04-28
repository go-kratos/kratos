package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
)

// 简单的流信息的处理和返回

// getStreamLastTime 得到流到最后推流时间
func getStreamLastTime(c *bm.Context) {
	// 获取url中的room_id
	params := c.Request.URL.Query()
	room := params.Get("room_id")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	t, err := srv.GetStreamLastTime(c, roomID)
	if err != nil {
		log.Warn("%v", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": map[string]int64{"last_time": t}}, nil)
}

// getStreamNameByRoomID 根据房间号获取流名
func getStreamNameByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	back := params.Get("back")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	flag := false
	if back == "1" {
		flag = true
	}

	info, err := srv.GetStreamNameByRoomID(c, roomID, flag)
	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	if !flag && len(info) > 0 {
		c.JSONMap(map[string]interface{}{"data": info[0]}, nil)
		return
	}

	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getRoomIdByStreamName 得到房间号，传递流名,可传入备用流名
func getRoomIDByStreamName(c *bm.Context) {
	params := c.Request.URL.Query()
	sname := params.Get("stream_name")

	sname = strings.TrimSpace(sname)
	if len(sname) == 0 {
		c.Set("output_data", "stream name is empty")
		c.JSONMap(map[string]interface{}{"message": "stream name is empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	rid, err := srv.GetRoomIDByStreamName(c, sname)

	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": map[string]int64{"room_id": rid}}, nil)
}

// getAdapterStreamByStreamName 得到适配的流信息，迁移PHP接口
func getAdapterStreamByStreamName(c *bm.Context) {
	params := c.Request.URL.Query()
	snames := params.Get("stream_names")

	snames = strings.TrimSpace(snames)
	if len(snames) == 0 {
		c.Set("output_data", "stream names is empty")
		c.JSONMap(map[string]interface{}{"message": "stream names is empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 最多查询500个数据
	nameSlice := strings.Split(snames, ",")
	if len(nameSlice) > 500 {
		c.Set("output_data", "too many names")
		c.JSONMap(map[string]interface{}{"message": "too many names"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info := srv.GetAdapterStreamByStreamName(c, nameSlice)

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getSrcByRoom 获取线路接口, 适配原PHP代码； 线路名称+线路编码src+是否当前选择的线路
func getSrcByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetSrcByRoomID(c, roomID)

	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getLineListByRoomID 得下线路信息， 和getSrcByRoomID 只有返回的格式不一样
func getLineListByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetLineListByRoomID(c, roomID)

	if err != nil {
		c.Set("output_data", err)
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}
