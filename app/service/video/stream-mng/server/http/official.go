package http

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"math/rand"
	"time"
)

// createOfficalStream 创建正式流
// optional string debug;  1表示线下测试
// required int uid; uid线下测试必传
func createOfficalStream(c *bm.Context) {
	defer c.Request.Body.Close()
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSONMap(map[string]interface{}{"message": err}, ecode.RequestErr)
		c.Abort()
		return
	}

	if len(b) == 0 {
		c.JSONMap(map[string]interface{}{"message": "empty params"}, ecode.RequestErr)
		c.Abort()
		return
	}

	type officialParams struct {
		RoomID     int64  `json:"room_id,omitempty"`
		StreamName string `json:"stream_name,omitempty"`
		Key        string `json:"key,omitempty"`
		Debug      string `json:"debug,omitempty"`
		Uid        int    `json:"uid,omitempty"`
	}

	var off officialParams

	err = json.Unmarshal(b, &off)
	if err != nil {
		c.JSONMap(map[string]interface{}{"message": err}, ecode.RequestErr)
		c.Abort()
		return
	}

	streamName := off.StreamName
	key := off.Key
	uid := off.Uid
	roomID := off.RoomID

	// 线下测试， 1表示线下测试,uid线下测试必传
	if uid != 0 {
		id := fmt.Sprintf("%d", uid)
		key = mockStreamKey(id)
		streamName = mockStreamName(id)
	}

	// 检查参数
	if streamName == "" || key == "" || roomID <= 0 {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	flag := srv.CreateOfficalStream(c, streamName, key, roomID)

	c.Set("output_data", fmt.Sprintf("create stream success = %v, room_id = %d", flag, roomID))
	c.JSONMap(map[string]interface{}{"data": map[string]bool{"succ": flag}}, nil)
}

// mockStream 模拟生成的流名
func mockStreamName(uid string) string {
	num := rand.Int63n(88888888)
	return fmt.Sprintf("live_%s_%d", uid, num+1111111)
}

// mockStreamKey 模拟生成的key
func mockStreamKey(uid string) string {
	str := fmt.Sprintf("nvijqwopW1%s%d", uid, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	return md5Str
}
