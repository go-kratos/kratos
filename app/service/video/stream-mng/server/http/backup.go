package http

import (
	"encoding/json"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"strconv"
)

// createBackupStream 创建备用流
func createBackupStream(c *bm.Context) {
	var bs model.BackupStream
	switch c.Request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if len(c.Request.PostForm) == 0 {
			c.Set("output_data", "empty params")
			c.JSONMap(map[string]interface{}{"message": "empty params"}, ecode.RequestErr)
			c.Abort()
			return
		}

		bs.StreamName = c.Request.PostFormValue("stream_name")
		bs.Key = c.Request.PostFormValue("key")

		default_vendor := c.Request.PostFormValue("default_vendor")
		vendor, _ := strconv.ParseInt(default_vendor, 10, 64)
		bs.DefaultVendor = vendor

		id := c.Request.PostFormValue("room_id")
		rid, _ := strconv.ParseInt(id, 10, 64)
		bs.RoomID = rid

		//bs.ExpiresAt = c.Request.PostFormValue("expires_at")
	default:
		defer c.Request.Body.Close()
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": err}, ecode.RequestErr)
			c.Abort()
			return
		}

		if len(b) == 0 {
			c.Set("output_data", "参数不能为空")
			c.JSONMap(map[string]interface{}{"message": "参数不能为空"}, ecode.RequestErr)
			c.Abort()
			return
		}

		err = json.Unmarshal(b, &bs)
		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": "请确认格式是否正常"}, ecode.RequestErr)
			c.Abort()
			return
		}

		if bs.RoomID <= 0 {
			c.Set("output_data", "房间号不正确")
			c.JSONMap(map[string]interface{}{"message": "房间号不正确"}, ecode.RequestErr)
			c.Abort()
			return
		}
	}

	c.Set("input_params", bs)

	_, err := srv.CreateBackupStream(c, &bs)

	c.JSON(bs, err)
}
