package service

import (
	"context"
	"encoding/json"
	"fmt"
	location "go-common/app/service/main/location/api"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/conf/env"
	"net/http"
)

const (
	// 下面三个是老上行推流表中 up_rank 的三个状态
	_originUpRankNothing = 0
	//_originUpRankDefaultSrc       = 1
	_originUpRankForwardStreaming = 2
)

// CDNSalt cdn salt
var CDNSalt = map[string]string{
	"bvc": "bvc_1701101740",
	"js":  "js_1703271720",
	"qn":  "qn_1703271730",
	"txy": "txy_1610171720",
	"ws":  "ws_1608121700",
}

// ValidateParams 验证流合法性
type ValidateParams struct {
	Key        string      `json:"key"`
	Type       json.Number `json:"type,omitempty"`
	StreamName string      `json:"stream_name"`
	Src        string      `json:"src"`
}

// getLiveStreamUrl 对接live-stream.bilibili.co的相关业务
func (s *Service) getLiveStreamUrl(path string) string {
	url := ""
	if env.DeployEnv == env.DeployEnvProd {
		url = fmt.Sprintf("%s%s", "http://prod-live-stream.bilibili.co", path)
	} else {
		url = fmt.Sprintf("%s%s", "http://live-stream.bilibili.co", path)
	}
	return url
}

// getRealRoomID 得到真正的ID
func (s *Service) getRealRoomID(rid int64) int64 {
	type RespStruct struct {
		Data *struct {
			RoomID int64 `json:"room_id"`
		} `json:"data"`
	}

	resp := RespStruct{}
	c := context.Background()
	uri := fmt.Sprintf("http://api.live.bilibili.com/room/v1/Room/room_init?id=%d", rid)
	err := s.dao.NewRequst(c, http.MethodGet, uri, nil, nil, nil, &resp)

	if err != nil || resp.Data == nil || resp.Data.RoomID == 0 {
		return rid
	}

	return resp.Data.RoomID
}

// getLocation grpc 调用
func (s *Service) GetLocation(c context.Context, ip string) *model.UpStreamInfo {
	res := &model.UpStreamInfo{}
	res.IP = ip

	if s.locationCli == nil {
		res.Country = "中国"
		return res
	}

	req := &location.InfoReq{Addr: ip}
	resp, err := s.locationCli.Info(c, req)
	if err != nil {
		res.Country = "中国"
		return res
	}

	res.Country = resp.Country
	res.City = resp.Province
	res.ISP = resp.Isp

	return res
}
