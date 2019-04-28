package dao

import (
	"crypto/md5"
	"encoding/hex"
	"go-common/app/interface/live/web-ucenter/conf"
	rank_api "go-common/app/service/live/rankdb/api/liverpc"
	rc_api "go-common/app/service/live/rc/api/liverpc"
	room_api "go-common/app/service/live/room/api/liverpc"
	"go-common/library/net/rpc/liverpc"
	"net/url"
	"strconv"
	"time"
)

// RoomAPI liverpc room-service api
var RoomAPI *room_api.Client

// RcApi liverpc rc api
var RcApi *rc_api.Client

// RankdbApi liverpc rankdb api
var RankdbApi *rank_api.Client

// InitAPI init all service APIs
func InitAPI() {
	RoomAPI = room_api.New(getConf("room"))
	RcApi = rc_api.New(getConf("rc"))
	RankdbApi = rank_api.New(getConf("rank"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}

// EncodeHttpParams end http params and return encoded string
func EncodeHttpParams(params map[string]string, appKey, appSecret string) string {
	v := url.Values{}
	for key, value := range params {
		v.Set(key, value)
	}
	v.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	v.Set("appkey", appKey)
	tmp := v.Encode() + appSecret
	mh := md5.Sum([]byte(tmp))
	sign := hex.EncodeToString(mh[:])
	v.Set("sign", sign)
	return v.Encode()
}
