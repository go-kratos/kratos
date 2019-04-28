package service

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
)

// GetStreamInfo 获取单个流信息
func (s *Service) GetStreamInfo(c context.Context, roomID int64, sname string) (*model.StreamFullInfo, error) {
	result, err := s.dao.StreamFullInfo(c, roomID, sname)
	return result, err
}

// AddHotStreamInfo 增加热流信息到REDIS
func (s *Service) AddHotStreamInfo(c context.Context, streamName string) (string, error) {
	roomid, _, err := s.dao.OriginUpStreamInfoBySName(c, streamName)
	if err != nil {
		return "error", err
	} else {
		if roomid <= 0 {
			return "error", fmt.Errorf("room id < 0")
		}
		err := s.dao.UpdateRoomHotStatusCache(c, roomid, 1)
		//err := s.dao.AddRoomHotCache(c, roomid)
		if err == nil {
			return "success", nil
		}
		return "error", err
	}
}

// GetMultiStreamInfo 批量获取流信息
func (s *Service) GetMultiStreamInfo(c context.Context, rids []int64) (map[int64]*model.StreamFullInfo, error) {
	result, err := s.dao.MultiStreamInfo(c, rids)

	return result, err
}

// GetStreamInfoByRIDMapSrcFromDB 传入房间号，返回房间流信息，包含备用流+使用原始src
func (s *Service) GetStreamInfoByRIDMapSrcFromDB(c context.Context, roomID int64) (*model.StreamFullInfo, error) {
	info, err := s.GetStreamInfo(c, roomID, "")
	if err != nil {
		return nil, err
	}

	return s.translateInfoBit2Src(info), nil
}

// GetStreamInfoBySNameMapSrcFromDB 传入流名，返回房间流信息，包含备用流+使用原始src
func (s *Service) GetStreamInfoBySNameMapSrcFromDB(c context.Context, sname string) (*model.StreamFullInfo, error) {
	info, err := s.GetStreamInfo(c, 0, sname)
	if err != nil {
		return nil, err
	}

	return s.translateInfoBit2Src(info), nil
}

// translatebit2Src 将cdn对应的bit 转为原始的src
func (s *Service) translatebit2Src(b int64) int64 {
	for k, v := range common.CdnBitwiseMap {
		if v == b {
			return int64(common.CdnMapSrc[k])
		}
	}
	return 0
}

// translateInfoBit2Src 将原始结构转为src对应的list
func (s *Service) translateInfoBit2Src(info *model.StreamFullInfo) *model.StreamFullInfo {
	resp := &model.StreamFullInfo{}
	// 适配src
	if info != nil {
		resp.RoomID = info.RoomID

		for _, v := range info.List {
			fw := []int64{}
			for _, v2 := range v.Forward {
				fw = append(fw, s.translatebit2Src(v2))
			}

			resp.List = append(resp.List, &model.StreamBase{
				StreamName:      v.StreamName,
				DefaultUpStream: s.translatebit2Src(v.DefaultUpStream),
				Origin:          s.translatebit2Src(v.Origin),
				Forward:         fw,
				Type:            v.Type,
			})
		}
	}
	return resp
}
