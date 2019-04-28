package http

import (
	"encoding/json"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

// clearStreamStatus 清理互推标志
func clearStreamStatus(c *bm.Context) {
	type room struct {
		RoomID json.Number `json:"room_id"`
	}

	vp := &room{}

	switch c.Request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if len(c.Request.PostForm) == 0 {
			c.Set("output_data", "clearStreamStatus = empty post body")
			c.JSONMap(map[string]interface{}{"message": "empty post body"}, ecode.RequestErr)
			c.Abort()
			return
		}
		vp.RoomID = json.Number(c.Request.PostFormValue("room_id"))
	default:
		defer c.Request.Body.Close()

		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
			c.Abort()
			return
		}

		if len(b) == 0 {
			c.Set("output_data", "clearStreamStatus  empty params")
			c.JSONMap(map[string]interface{}{"message": "empty params"}, ecode.RequestErr)
			c.Abort()
			return
		}

		err = json.Unmarshal(b, &vp)
		if err != nil {
			c.Set("output_data", "room_id is not right")
			c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
			c.Abort()
			return
		}
	}

	roomID, err := vp.RoomID.Int64()

	if roomID <= 0 || err != nil {
		c.Set("output_data", "room_id is not right")
		c.JSONMap(map[string]interface{}{"message": "room_id is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("input_params", map[string]int64{"room_id": roomID})

	err = srv.ClearStreamStatus(c, roomID)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.Set("output_data", "clear status success")
	c.JSONMap(map[string]interface{}{"message": "ok"}, nil)
}
