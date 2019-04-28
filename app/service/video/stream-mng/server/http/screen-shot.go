package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
	"time"
)

//  截图相关业务

// getScreeShotByRoomID 得到一个房间某个时间段的截图
func getSingleScreenShot(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	// 2018-10-24 14:27:07
	start := params.Get("start_time")
	end := params.Get("end_time")
	channel := params.Get("channel")

	c.Set("input_params", params)

	if room == "" || start == "" || end == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		c.Set("output_data", "Start time format is incorrect")
		c.JSONMap(map[string]interface{}{"message": "Start time format is incorrect"}, ecode.RequestErr)
		c.Abort()
		return
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err != nil {
		c.Set("output_data", "End time format is incorrect")
		c.JSONMap(map[string]interface{}{"message": "End time format is incorrect"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetSingleScreeShot(c, roomID, startTime.Unix(), endTime.Unix(), channel)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": map[string][]string{"list": info}}, nil)
}

// getMultiScreenShot 得到多个房间一个时间点截图
func getMultiScreenShot(c *bm.Context) {
	params := c.Request.URL.Query()
	rooms := params.Get("room_ids")
	ts := params.Get("ts")
	channel := params.Get("channel")

	c.Set("input_params", params)

	if rooms == "" || ts == "" {
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

	tsInt, _ := strconv.ParseInt(ts, 10, 64)
	rids := []int64{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warn("room id is not right")
			continue
		}

		rids = append(rids, roomID)
	}

	resp, err := srv.GetMultiScreenShot(c, rids, tsInt, channel)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", resp)
	c.JSONMap(map[string]interface{}{"data": map[string]interface{}{"list": resp}}, nil)
}

// getOriginScreenShotPic 获取原始图片地址
func getOriginScreenShotPic(c *bm.Context) {
	params := c.Request.URL.Query()
	rooms := params.Get("room_ids")
	ts := params.Get("ts")
	//tp := params.Get("type")

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

	tsInt, _ := strconv.ParseInt(ts, 10, 64)

	rids := []int64{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warn("room id is not right")
			continue
		}

		rids = append(rids, roomID)
	}

	resp, err := srv.GetOriginScreenShotPic(c, rids, tsInt)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", resp)
	c.JSONMap(map[string]interface{}{"data": map[string]interface{}{"list": resp}}, nil)
}

// getTimePeriodScreenShot 获取多个房间一个时间段内的截图
func getTimePeriodScreenShot(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_ids")
	// 2018-10-24 14:27:07
	start := params.Get("start_time")
	end := params.Get("end_time")
	channel := params.Get("channel")

	c.Set("input_params", params)

	if room == "" || start == "" || end == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		c.Set("output_data", "Start time format is incorrect")
		c.JSONMap(map[string]interface{}{"message": "Start time format is incorrect"}, ecode.RequestErr)
		c.Abort()
		return
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err != nil {
		c.Set("output_data", "End time format is incorrect")
		c.JSONMap(map[string]interface{}{"message": "End time format is incorrect"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 切割room_id
	roomIDs := strings.Split(room, ",")

	if len(roomIDs) <= 0 {
		c.Set("output_data", "room_ids is not right")
		c.JSONMap(map[string]interface{}{"message": "room_ids is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	resp := map[int64][]string{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Warn("room id is not right")
			continue
		}

		urls, err := srv.GetSingleScreeShot(c, roomID, startTime.Unix(), endTime.Unix(), channel)
		if err != nil {
			log.Warn("%v", err)
			continue
		}

		resp[roomID] = urls
	}

	c.Set("output_data", resp)
	c.JSONMap(map[string]interface{}{"data": resp}, nil)
}
