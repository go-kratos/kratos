package http

import (
	"fmt"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"net"
	"strconv"
)

// 获取上行推流地址， 一共三个方法调用

// UpStream
func getUpStreamRtmp(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	// 来源：pc：表示PC端；ios：表示ios端；android：安卓端；ios_link：表示ios端；android_link： 安卓端；live_mng：表示live后台;vc_mng：表示vc后台;
	platform := params.Get("platform")
	// client_ip
	ip := params.Get("ip")
	// 分区id
	area := params.Get("area_id")
	// 免流标志
	freeFlow := params.Get("free_flow")
	attentions := params.Get("attentions")

	c.Set("input_params", params)

	if room == "" || platform == "" || area == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "房间号错误"}, ecode.RequestErr)
		c.Abort()
		return
	}

	var attentionsInt int
	if attentions == "" {
		attentionsInt = 0
	} else {
		att, _ := strconv.ParseInt(attentions, 10, 64)
		attentionsInt = int(att)
	}

	// ip映射
	realIP := ip
	if ip == "" {
		remoteAddr := c.Request.RemoteAddr
		// 使用header: X-REAL-IP + X_FORWARED_FOR + reamoteadd
		if add := c.Request.Header.Get("X-REAL-IP"); add != "" {
			remoteAddr = add
		} else if add = c.Request.Header.Get("X_FORWARED_FOR"); add != "" {
			remoteAddr = add
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}

		if remoteAddr == "::1" {
			remoteAddr = "127.0.0.1"
		}

		realIP = remoteAddr
	}

	areaID, _ := strconv.ParseInt(area, 10, 64)

	info, err := srv.GetUpStreamRtmp(c, roomID, freeFlow, realIP, areaID, attentionsInt, 0, platform)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	if info == nil {
		c.Set("output_data", fmt.Sprintf("can find any info by room_id=%d", roomID))
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": map[string]interface{}{"up_stream": info}}, nil)
}

// getWebRtmp web端调用
func getWebRtmp(c *bm.Context) {
	// 获取room_id
	params := c.Request.URL.Query()
	room := params.Get("room_id")

	c.Set("input_params", params)

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "房间号不正确"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// 获取uid
	uid, ok := metadata.Value(c, metadata.Mid).(int64)
	//uid = 19148701
	//ok = true

	//log.Infov(c, log.KV("log", fmt.Sprintf("uid=%v", uid)))
	if !ok {
		log.Warn("%v=%v", uid, ok)
		c.Set("output_data", "未登陆")
		c.JSONMap(map[string]interface{}{"message": fmt.Sprintf("未登陆")}, ecode.RequestErr)
		c.Abort()
		return
	}

	remoteAddr := c.Request.RemoteAddr
	// 使用header: X-REAL-IP + X_FORWARED_FOR + reamoteadd
	if add := c.Request.Header.Get("X-REAL-IP"); add != "" {
		remoteAddr = add
	} else if add = c.Request.Header.Get("X_FORWARED_FOR"); add != "" {
		remoteAddr = add
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	realIP := remoteAddr

	info, err := srv.GetWebRtmp(c, roomID, uid, realIP, "web")
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	if info == nil {
		c.Set("output_data", fmt.Sprintf("can find any info by room_id=%d", roomID))
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}

// getMobileRtmp 移动端调用
func getMobileRtmp(c *bm.Context) {
	params := c.Request.URL.Query()
	room := params.Get("room_id")
	// 来源：pc：表示PC端；ios：表示ios端；android：安卓端；ios_link：表示ios端；android_link： 安卓端；live_mng：表示live后台;vc_mng：表示vc后台;
	platform := params.Get("platform")
	// client_ip
	ip := params.Get("ip")
	// 分区id
	area := params.Get("area_id")
	// 免流标志
	freeFlow := params.Get("free_flow")

	c.Set("input_params", params)

	if room == "" || platform == "" || area == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "房间号错误")
		c.JSONMap(map[string]interface{}{"message": "房间号错误"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// ip映射
	realIP := ip
	if ip == "" {
		remoteAddr := c.Request.RemoteAddr
		// 使用header: X-REAL-IP + X_FORWARED_FOR + reamoteadd
		if add := c.Request.Header.Get("X-REAL-IP"); add != "" {
			remoteAddr = add
		} else if add = c.Request.Header.Get("X_FORWARED_FOR"); add != "" {
			remoteAddr = add
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}

		if remoteAddr == "::1" {
			remoteAddr = "127.0.0.1"
		}

		realIP = remoteAddr
	}

	areaID, _ := strconv.ParseInt(area, 10, 64)

	// 获取uid
	uid, ok := metadata.Value(c, metadata.Mid).(int64)
	//uid = 19148701
	//ok = true
	if !ok {
		c.Set("output_data", "未登陆")
		c.JSONMap(map[string]interface{}{"message": fmt.Sprintf("未登陆")}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetUpStreamRtmp(c, roomID, freeFlow, realIP, areaID, 0, uid, platform)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	if info == nil {
		c.Set("output_data", fmt.Sprintf("can find any info by room_id=%d", roomID))
		c.JSONMap(map[string]interface{}{"message": "获取线路信息失败，刷新页面或稍后重试"}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": map[string]interface{}{"up_stream": info}}, nil)
}

// getRoomRtmp 拜年祭房间推流码接口
func getRoomRtmp(c *bm.Context) {
	params := c.Request.URL.Query()

	c.Set("input_params", params)

	room := params.Get("room_id")
	roomID, err := strconv.ParseInt(room, 10, 64)

	if err != nil || roomID <= 0 {
		c.Set("output_data", "房间号不正确")
		c.JSONMap(map[string]interface{}{"message": "房间号不正确"}, ecode.RequestErr)
		c.Abort()
		return
	}

	info, err := srv.GetRoomRtmp(c, roomID)

	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": "获取房间信息失败"}, ecode.RequestErr)
		c.Abort()
		return
	}

	if info == nil {
		c.Set("output_data", fmt.Sprintf("can find any info by room_id=%d", roomID))
		c.JSONMap(map[string]interface{}{"message": "获取房间信息失败，请确认是否房间存在"}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", info)
	c.JSONMap(map[string]interface{}{"data": info}, nil)
}
