package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"strconv"
)

// 单独文件，切换cdn

// changeSrc 切换cdn
func changeSrc(c *bm.Context) {
	// 这里传递的src是新的src
	type changStruct struct {
		RoomID      int64  `json:"room_id,omitempty"`
		ToOrigin    int8   `json:"src,omitempty"`
		Source      string `json:"source,omitempty"`
		OperateName string `json:"operate_name,omitempty"`
		Reason      string `json:"reason,omitempty"`
	}

	var cs changStruct
	switch c.Request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if len(c.Request.PostForm) == 0 {
			c.Set("output_data", "参数为空")
			c.JSONMap(map[string]interface{}{"message": "参数为空"}, ecode.RequestErr)
			c.Abort()
			return
		}
		toOrigin := c.Request.PostFormValue("src")
		or, _ := strconv.ParseInt(toOrigin, 10, 64)
		cs.ToOrigin = int8(or)

		rid, _ := strconv.ParseInt(c.Request.PostFormValue("room_id"), 10, 64)
		cs.RoomID = rid

		cs.Source = c.Request.PostFormValue("source")
		cs.OperateName = c.Request.PostFormValue("operate_name")
	default:
		// 验证传参数
		defer c.Request.Body.Close()
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSONMap(map[string]interface{}{"message": err}, ecode.RequestErr)
			c.Abort()
			return
		}

		if len(b) == 0 {
			c.JSONMap(map[string]interface{}{"message": "参数为空"}, ecode.RequestErr)
			c.Abort()
			return
		}

		err = json.Unmarshal(b, &cs)
		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
			c.Abort()
			return
		}
	}

	c.Set("input_params", cs)

	// 校验：房间号+src+平台来源+操作人 都是必须的， 操作理由可以不填
	if cs.RoomID <= 0 || cs.ToOrigin == 0 || cs.Source == "" || cs.OperateName == "" {
		c.Set("output_data", "some fields are not right")
		c.JSONMap(map[string]interface{}{"message": "部分参数为空"}, ecode.RequestErr)
		c.Abort()
		return
	}

	// todo  先使用老的src, 后续改为新的src
	src := int8(cs.ToOrigin)

	if _, ok := common.SrcMapBitwise[src]; !ok {
		c.Set("output_data", "src is not right")
		c.JSONMap(map[string]interface{}{"message": "src is not right"}, ecode.RequestErr)
		c.Abort()
		return
	}

	err := srv.ChangeSrc(c, cs.RoomID, common.SrcMapBitwise[src], cs.Source, cs.OperateName, cs.Reason)

	if err == nil {
		c.Set("output_data", fmt.Sprintf("room_id = %d, change src success", cs.RoomID))
		c.JSONMap(map[string]interface{}{"message": "ok"}, nil)
		c.Abort()
		return
	}

	c.Set("output_data", fmt.Sprintf("room_id = %d, change src faild = %v", cs.RoomID, err))
	c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
}
