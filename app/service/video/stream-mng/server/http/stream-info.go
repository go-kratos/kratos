package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"strconv"
	"strings"
)

// addHotStream 增加房间热流标记
func addHotStream(c *bm.Context) {
	req := c.Request
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.Set("output_data", "ioutil.ReadAll() error")
		c.JSONMap(map[string]interface{}{"message": "outil.ReadAll() error"}, err)
		c.Abort()
		return
	}
	req.Body.Close()
	var hrbody []string
	if err := json.Unmarshal(bs, &hrbody); err != nil {
		c.Set("output_data", "json.Unmarshal() error")
		c.JSONMap(map[string]interface{}{"message": "json.Unmarshal() error"}, err)
		c.Abort()
		return
	}
	if len(hrbody) <= 0 {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}
	for _, streamName := range hrbody {
		srv.AddHotStreamInfo(c, streamName)
	}
	c.Set("output_data", "success")
	c.JSONMap(map[string]interface{}{"data": "success"}, nil)
}

// getStream 获取单个流信息
func getStream(c *bm.Context) {
	params := c.Request.URL.Query()
	rid := params.Get("room_id")
	sname := params.Get("stream_name")

	if rid == "" && sname == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	var roomID int64
	var err error
	var info *model.StreamFullInfo
	if sname == "" {
		roomID, err = strconv.ParseInt(rid, 10, 64)

		// 验证传参数
		if err != nil || roomID <= 0 {
			c.Set("output_data", "roomid is not right")
			c.JSONMap(map[string]interface{}{"message": "roomid is not right"}, ecode.RequestErr)
			c.Abort()
			return
		}
		info, err = srv.GetStreamInfo(c, roomID, "")
	} else {
		info, err = srv.GetStreamInfo(c, 0, sname)
	}

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getMulitiStreams 批量查询流接口
func getMultiStreams(c *bm.Context) {
	params := c.Request.URL.Query()
	roomID := params.Get("room_ids")

	if roomID == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	roomIDs := []int64{}

	rids := strings.Split(roomID, ",")
	for _, v := range rids {
		rid, err := strconv.ParseInt(v, 10, 64)

		// 验证传参数
		if err == nil && rid > 0 {
			roomIDs = append(roomIDs, rid)
		}
	}

	if len(roomIDs) > 30 {
		c.Set("output_data", "The number of rooms must be less than 30")
		c.JSONMap(map[string]interface{}{"message": "The number of rooms must be less than 30"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetMultiStreamInfo(c, roomIDs)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	if info == nil || len(info) == 0 {
		c.Set("output_data", fmt.Sprintf("can not find any info by room_ids=%s", roomID))
		c.JSONMap(map[string]interface{}{"message": fmt.Sprintf("can not find any info by room_ids=%s", roomID)}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getOldStreamInfoByRoomID map 到原始src数据
func getOldStreamInfoByRoomID(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("roomid")
	room2 := params.Get("room_id")

	rid := ""
	if room == "" {
		rid = room2
	} else {
		rid = room
	}

	roomID, err := strconv.ParseInt(rid, 10, 64)

	// 验证传参数
	if err != nil || roomID <= 0 {
		c.Set("output_data", "roomid is not right")
		c.JSONMap(map[string]interface{}{"message": "roomid is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetStreamInfoByRIDMapSrcFromDB(c, roomID)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getOldStreamInfoByStreamName map到原始src数据
func getOldStreamInfoByStreamName(c *bm.Context) {
	params := c.Request.URL.Query()
	sname := params.Get("stream_name")

	sname = strings.TrimSpace(sname)

	if sname == "" {
		c.JSONMap(map[string]interface{}{"message": "stream name is empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetStreamInfoBySNameMapSrcFromDB(c, sname)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}
