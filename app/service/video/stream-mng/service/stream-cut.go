package service

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/library/log"
	"net/http"
)

// StreamCut 切流的房间和时间, 内部调用
func (s *Service) StreamCut(c context.Context, rid int64, ct int64, uid int64) error {
	if uid != 0 {
		if uid != s.getRoomUserID(c, rid) {
			return fmt.Errorf("无权限,不是该房间用户")
		}
	}

	// 根据房间号查询流名和直推&转推的cdn
	info, err := s.dao.StreamFullInfo(c, rid, "")

	if err != nil || info == nil || len(info.List) == 0 {
		log.Errorv(c, log.KV("log", fmt.Sprintf("stream cut err = %v", err)))
		return err
	}

	for _, v := range info.List {
		if v.Type == 1 {
			// 发送切流请求
			uri := s.getLiveStreamUrl("/api/live/vendor/cutstream")

			uri = fmt.Sprintf("%s?stream_name=%s&cut_time=%d", uri, v.StreamName, ct)

			if v.Origin != 0 {
				uri = fmt.Sprintf("%s&src=%s", uri, common.BitwiseMapName[v.Origin])
			}
			for _, i := range v.Forward {
				if i != 0 {
					uri = fmt.Sprintf("%s&src=%s", uri, common.BitwiseMapName[i])
				}
			}

			log.Infov(c, log.KV("log", fmt.Sprintf("url=%s", uri)))
			err = s.dao.NewRequst(c, http.MethodGet, uri, nil, nil, nil, nil)
			if err != nil {
				log.Errorv(c, log.KV("log", fmt.Sprintf("stream cut err = %v", err)))
			}
			break
		}
	}
	return nil
}
