package dao

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// 获取流完整信息
	// cache: -singleflight=true -ignores=||id,sname -nullcache=&model.StreamFullInfo{RoomID:-1} -check_null_code=$!=nil&&$.RoomID<=0
	streamFullInfo(c context.Context, id int64, sname string) (res *model.StreamFullInfo, err error)

	// 获取rid
	// cache: -singleflight=true -ignores=||sname -nullcache=&model.StreamFullInfo{RoomID:-1} -check_null_code=$!=nil&&$.RoomID<=0
	streamRIDByName(c context.Context, sname string) (res *model.StreamFullInfo, err error)

	// 批量获取接口
	// cache: -nullcache=&model.StreamFullInfo{RoomID:-1} -check_null_code=$!=nil&&$.RoomID<=0
	multiStreamInfo(c context.Context, rid []int64) (res map[int64]*model.StreamFullInfo, err error)
}

func (d *Dao) cacheSFstreamFullInfo(id int64, sname string) string {
	if sname != "" {
		return fmt.Sprintf("sf_sname_%s", sname)
	}
	return fmt.Sprintf("sf_rid_%d", id)
}

func (d *Dao) cacheSFstreamRIDByName(sname string) string {
	return fmt.Sprintf("sf_rid_map_name_%s", sname)
}

// StreamFullInfo 传入rid或者sname 获取房间流信息
func (d *Dao) StreamFullInfo(c context.Context, rid int64, sname string) (res *model.StreamFullInfo, err error) {
	info, err := d.streamFullInfo(c, rid, sname)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("can not find by room_id=%d", rid)
	}

	if len(info.List) == 1 && info.List[0].StreamName == "miss" {
		return nil, fmt.Errorf("can not find any info by room_id=%d", rid)
	}

	return info, nil
}

// OriginUpStreamInfo 原始上行流名和src
func (d *Dao) OriginUpStreamInfo(c context.Context, rid int64) (sname string, origin int64, err error) {
	info, err := d.streamFullInfo(c, rid, "")
	if err != nil {
		return "", 0, err
	}
	if info == nil {
		return "", 0, fmt.Errorf("can not find by room_id=%d", rid)
	}

	for _, v := range info.List {
		if v.Type == 1 {
			//  优先级高
			if v.Origin != 0 {
				return v.StreamName, v.Origin, nil
			}
			return v.StreamName, v.DefaultUpStream, nil
		}
	}
	return "", 0, fmt.Errorf("can not find by room_id=%d", rid)
}

func (d *Dao) DefaultUpStreamInfo(c context.Context, rid int64) (sname string, origin int64, err error) {
	info, err := d.streamFullInfo(c, rid, "")
	if err != nil {
		return "", 0, err
	}
	if info == nil {
		return "", 0, fmt.Errorf("can not find by room_id=%d", rid)
	}

	for _, v := range info.List {
		if v.Type == 1 {
			return v.StreamName, v.DefaultUpStream, nil
		}
	}
	return "", 0, fmt.Errorf("can not find by room_id=%d", rid)
}

// OriginUpStreamInfoBySName 查询流的上行，正在推流上行和默认上行; 包含备用流
func (d *Dao) OriginUpStreamInfoBySName(c context.Context, sname string) (rid int64, origin int64, err error) {
	info, err := d.streamFullInfo(c, 0, sname)
	if err != nil {
		return 0, 0, err
	}
	if info == nil {
		return 0, 0, fmt.Errorf("can not find by sname=%s", sname)
	}

	for _, v := range info.List {
		if v.StreamName == sname {
			//  优先级高
			if v.Origin != 0 {
				return info.RoomID, v.Origin, nil
			}
			return info.RoomID, v.DefaultUpStream, nil
		}
	}
	return 0, 0, fmt.Errorf("can not find by sname=%s", sname)
}

// StreamRIDByName 获取rid
func (d *Dao) StreamRIDByName(c context.Context, sname string) (int64, error) {
	info, err := d.streamRIDByName(c, sname)
	if err != nil {
		return -1, err
	}

	if info != nil && info.RoomID > 0 {
		return info.RoomID, nil
	}
	return -1, fmt.Errorf("can not find by sname=%s", sname)
}

// MultiStreamInfo 批量接口
func (d *Dao) MultiStreamInfo(c context.Context, rids []int64) (res map[int64]*model.StreamFullInfo, err error) {
	infos, err := d.multiStreamInfo(c, rids)
	if err != nil {
		return res, err
	}

	resp := map[int64]*model.StreamFullInfo{}
	for k, v := range infos {
		if len(v.List) == 1 && v.List[0].StreamName == "miss" {
			continue
		}
		resp[k] = v
	}

	return resp, nil
}
