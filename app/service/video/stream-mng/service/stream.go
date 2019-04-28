package service

import (
	"context"
	"go-common/app/service/video/stream-mng/common"
	"go-common/library/log"
	"net/http"
	"net/url"
)

type adpterStream struct {
	Src     int64  `json:"src"`
	RoomID  int64  `json:"room_id"`
	UpRank  int64  `json:"up_rank"`
	SrcName string `json:"src_name"`
}

type roomSrc struct {
	Src     int64  `json:"src"`
	Checked int    `json:"checked"`
	Desc    string `json:"desc"`
}

type lineList struct {
	Src  int64  `json:"src"`
	Use  bool   `json:"use"`
	Desc string `json:"desc"`
}

// GetStreamLastTime 得到流到最后推流时间;主流的推流时间up_rank = 1
func (s *Service) GetStreamLastTime(c context.Context, rid int64) (t int64, err error) {
	// 读取上行
	streamName, origin, err := s.dao.OriginUpStreamInfo(c, rid)

	if err != nil {
		return 0, err
	}

	cdn := ""
	for k, v := range common.BitwiseMapName {
		if k == origin {
			cdn = v
			break
		}
	}

	// 发送请求，三次重试机制
	uri := s.getLiveStreamUrl("/api/live/vendor/checkstream")
	params := make(url.Values)
	params.Set("cdn", cdn)
	params.Set("stream_name", streamName)

	type RespStruct struct {
		Code int                 `json:"code"`
		Data []*map[string]int64 `json:"data"`
	}

	// 重试机制
	for i := 0; i < 3; i++ {
		resp := RespStruct{}
		err := s.dao.NewRequst(c, http.MethodGet, uri, params, nil, nil, &resp)
		if err != nil {
			log.Warn("http request err = %v", err)
			continue
		}

		if resp.Code == 0 && len(resp.Data) > 0 {
			d := *resp.Data[0]
			return d[streamName], nil
		}
	}

	return 0, nil
}

// GetStreamNameByRoomID 需要考虑备用流 + 考虑短号
func (s *Service) GetStreamNameByRoomID(c context.Context, rid int64, back bool) ([]string, error) {
	resp := []string{}

	// 考虑短号
	realRoomID := s.getRealRoomID(rid)

	if !back {
		name, _, err := s.dao.OriginUpStreamInfo(c, realRoomID)

		if err != nil {
			return resp, err
		}

		resp = append(resp, name)

		return resp, nil
	}

	infos, err := s.dao.StreamFullInfo(c, realRoomID, "")

	if err != nil {
		return resp, err
	}

	if infos != nil && infos.List != nil {
		for _, v := range infos.List {
			resp = append(resp, v.StreamName)
		}
	}

	return resp, nil
}

// GetRoomIDByStreamName 查询房间号
func (s *Service) GetRoomIDByStreamName(c context.Context, sname string) (int64, error) {
	return s.dao.StreamRIDByName(c, sname)
}

// GetAdapterStreamByStreamName 适配结果输出, 此处也可以输入备用流， 该结果只输出直推上行
func (s *Service) GetAdapterStreamByStreamName(c context.Context, names []string) map[string]*adpterStream {
	resp := map[string]*adpterStream{}

	for _, n := range names {
		if n != "" {
			rid, origin, err := s.dao.OriginUpStreamInfoBySName(c, n)

			if err != nil {
				return resp
			}

			item := &adpterStream{}
			item.RoomID = rid

			item.UpRank = 1
			// todo changesrc 完全上线后才使用
			//item.Src = origin
			item.Src = int64(common.BitwiseMapSrc[origin])
			item.SrcName = common.NameMapBitwise[origin]

			resp[n] = item
		}
	}
	return resp
}

// GetSrcByRoomID 得到src
func (s *Service) GetSrcByRoomID(c context.Context, rid int64) ([]*roomSrc, error) {
	_, origin, err := s.dao.DefaultUpStreamInfo(c, rid)

	if err != nil {
		return nil, err
	}

	resp := []*roomSrc{}

	for k, v := range common.ChinaNameMapBitwise {
		checked := 0
		if v == origin {
			checked = 1
		}

		// todo 等changesrc上线
		resp = append(resp, &roomSrc{
			//Src:     v,
			Src:     int64(common.BitwiseMapSrc[v]),
			Checked: checked,
			Desc:    k,
		})
	}

	return resp, nil
}

// GetLineListByRoomID 获取线路信息
func (s *Service) GetLineListByRoomID(c context.Context, rid int64) ([]*lineList, error) {
	_, origin, err := s.dao.DefaultUpStreamInfo(c, rid)
	if err != nil {
		return nil, err
	}

	resp := []*lineList{}
	//aUse := lineList{}

	//empty := false
	for k, v := range common.ChinaNameMapBitwise {
		// 网宿暂不开发到后台
		if v == common.BitWiseWS {
			continue
		}

		use := false
		if v == origin {
			use = true
		}

		//if use {
		//	//aUse.Src = v
		//	aUse.Src = int64(common.BitwiseMapSrc[v])
		//	aUse.Use = use
		//	aUse.Desc = k
		//} else {
		//	resp = append(resp, lineList{
		//		Src:  int64(common.BitwiseMapSrc[v]),
		//		Use:  use,
		//		Desc: k,
		//	})
		//	//empty = true
		//}

		resp = append(resp, &lineList{
			Src:  int64(common.BitwiseMapSrc[v]),
			Use:  use,
			Desc: k,
		})
	}

	// 十九大cdn情况兼容
	//setEmpty := &lineList{
	//	Src:  int64(999),
	//	Use:  empty,
	//	Desc: "空",
	//}
	//
	//if aUse.Src != 0 {
	//	resp = append(resp, &aUse, setEmpty)
	//} else {
	//	// 上行在网宿，就会执行
	//	resp = append(resp, empty)
	//}
	return resp, nil
}
