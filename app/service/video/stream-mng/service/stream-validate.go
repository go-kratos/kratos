package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/library/net/metadata"
	"net/http"
	"strings"
)

// CheckStreamValidate 根据流名查询正式流
func (s *Service) CheckStreamValidate(c context.Context, vp *ValidateParams, notify bool) (int, error) {
	// 首先验证参数不能为空
	vp.StreamName = strings.TrimSpace(vp.StreamName)
	vp.Key = strings.TrimSpace(vp.Key)
	vp.Src = strings.TrimSpace(vp.Src)
	if vp.StreamName == "" || vp.Src == "" {
		return 0, errors.New("some fields are empty")
	}

	checked := false
	rid, _ := s.dao.StreamRIDByName(c, vp.StreamName)
	if rid > 0 {
		if !strings.Contains(vp.StreamName, "_bs_") {
			checked = true
			// 检验房间是否开播状态
			bd := map[string]interface{}{
				"ids":    []int64{rid},
				"fields": []string{"live_time"},
			}
			bdj, _ := json.Marshal(bd)

			url := "http://api.live.bilibili.co/room/v2/Room/get_by_ids"

			type liveTimeStruct struct {
				Data map[string]map[string]interface{} `json:"data"`
			}

			liveResp := &liveTimeStruct{}

			head := map[string]string{
				"Content-Type": "application/json",
			}

			err := s.dao.NewRequst(c, http.MethodPost, url, nil, bdj, head, liveResp)
			if err != nil {
				return 0, err
			}

			liveTime := liveResp.Data[fmt.Sprintf("%d", rid)]["live_time"]
			if liveTime == "0000-00-00 00:00:00" {
				return 0, errors.New("room is closed")
			}
		}
	}

	// 验证src是否合法
	iSc := common.CdnMapSrc[vp.Src]
	if iSc == 0 {
		return 0, errors.New("src is not right")
	}

	// 转换type
	vpType, err := vp.Type.Int64()
	if err != nil {
		return 0, errors.New("type is not right")
	}

	// 返回信息
	roomInfo, err := s.dao.StreamFullInfo(c, 0, vp.StreamName)

	// 可以查询到
	if err == nil && roomInfo != nil && len(roomInfo.List) > 0 {
		realKey := ""
		var origin int64
		var rid int64
		isBack := false
		for _, v := range roomInfo.List {
			if v.StreamName == vp.StreamName {
				if v.Type == 2 {
					isBack = true
				}
				realKey = v.Key
				origin = v.DefaultUpStream
				rid = roomInfo.RoomID
			}
		}

		if !notify {
			// 验证key是否匹配
			if !strings.Contains(vp.Key, realKey) {
				return 0, errors.New("key is not right")
			}
		}

		// 主流需要验证，备用流不需要验证
		if !isBack {
			// 直/互推参数检验，type 1:互推; 0直推
			// src相等，但是type为互推1 或者 src不相等， type为直推0
			upSrc := common.BitwiseMapSrc[origin]

			if ((upSrc == iSc) && vpType == 1) || (upSrc != iSc && vpType == 0) {
				return 0, errors.New("key type is not right")
			}

			if !checked {
				// 检验房间是否开播状态
				bd := map[string]interface{}{
					"ids":    []int64{rid},
					"fields": []string{"live_time"},
				}
				bdj, _ := json.Marshal(bd)

				url := "http://api.live.bilibili.co/room/v2/Room/get_by_ids"

				type liveTimeStruct struct {
					Data map[string]map[string]interface{} `json:"data"`
				}

				liveResp := &liveTimeStruct{}

				head := map[string]string{
					"Content-Type": "application/json",
				}

				err = s.dao.NewRequst(c, http.MethodPost, url, nil, bdj, head, liveResp)
				if err != nil {
					return 0, err
				}

				liveTime := liveResp.Data[fmt.Sprintf("%d", rid)]["live_time"]
				if liveTime == "0000-00-00 00:00:00" {
					return 0, errors.New("room is closed")
				}
			}
		}

		// 同步数据
		go func(ctx context.Context, roomID int64) {
			s.syncMainStream(ctx, roomID, "")
		}(metadata.WithContext(c), rid)
		return 1, nil
	}

	return 0, nil
}
