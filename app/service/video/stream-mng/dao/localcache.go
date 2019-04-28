package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"net/http"
	"time"
)

const (
	_localLiveStreamList    = "live_stream"
	_localStreamInfo        = "rid:%d"
	_liveExpiredTime        = 600
	_localStreamExpiredTime = 1
)

type OnAirStream struct {
	StreamName string `json:"stream_name, omitempty"`
}

type OnAirStreamList struct {
	List []*OnAirStream `json:"list,omitempty"`
}

type AllOnAirStream struct {
	M map[string]*OnAirStreamList `json:"m,omitempty"`
}

func (d *Dao) getLocalLiveStreamListKey() string {
	return _localLiveStreamList
}

func (d *Dao) getLocalStreamInfoKey(rid int64) string {
	return fmt.Sprintf(_localStreamInfo, rid)
}

// loadLiveStreamList 判断流是否在播
func (d *Dao) LoadLiveStreamList(c context.Context, rids []int64) map[int64]bool {
	list, _ := d.localCache.Get(d.getLocalLiveStreamListKey())

	isLive := map[int64]bool{}
	if res, ok := list.(map[string]int); ok {
		for _, v := range rids {
			// rid => stream
			info, err := d.streamFullInfo(c, v, "")
			if err != nil || info == nil {
				continue
			}

			var sname string
			var origin int64
			for _, v := range info.List {
				if v.Type == 1 {
					sname = v.StreamName
					origin = v.Origin
					break
				}
			}
			// 不在在播列表&orgin 为0===》不在播
			if _, exe := res[sname]; !exe {
				if origin == 0 {
					isLive[v] = false
					continue
				}
			}
			isLive[v] = true
		}
	} else {
		// 当从缓存中获取失败 or 项目刚刚启动
		for _, v := range rids {
			isLive[v] = true
		}
	}
	return isLive
}

// StoreLiveStreamList 刷新在播列表缓存
func (d *Dao) StoreLiveStreamList() {
	ctx := context.Background()

	type liveStream struct {
		Code int             `json:"code,omitempty"`
		Data *AllOnAirStream `json:"data,omitempty"`
	}

	resp := &liveStream{}
	uri := d.getLiveStreamUrl("/api/live/vendor/onairstreamlist?cdn=bvc")

	err := d.NewRequst(ctx, http.MethodGet, uri, nil, nil, nil, resp)
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("http_live_err=%v", err)))
		return
	}

	res := map[string]int{}
	if resp.Code == 0 && resp.Data != nil && resp.Data.M != nil && len(resp.Data.M["BVC"].List) > 0 {
		for _, v := range resp.Data.M["BVC"].List {
			res[v.StreamName] = 1
		}
	} else {
		res, _ := json.Marshal(resp)
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("http_live_err=%v", string(res))))
	}

	// 存10分钟
	d.localCache.SetWithExpire(d.getLocalLiveStreamListKey(), res, _liveExpiredTime*time.Second)
}

// storeStreamInfo 存储单个流信息
func (d *Dao) storeStreamInfo(c context.Context, info *model.StreamFullInfo) {
	if info == nil || info.RoomID < 0 {
		return
	}
	key := d.getLocalStreamInfoKey(info.RoomID)

	d.localCache.SetWithExpire(key, info, _localStreamExpiredTime*time.Second)
}

// loadStreamInfo 读取单个信息
func (d *Dao) loadStreamInfo(c context.Context, rid int64) *model.StreamFullInfo {
	key := d.getLocalStreamInfoKey(rid)
	info, _ := d.localCache.Get(key)
	if res, ok := info.(*model.StreamFullInfo); ok {
		if res != nil && res.RoomID > 0 {
			return res
		}
	}

	return nil
}

// storeMultiStreamInfo 存储多路流
func (d *Dao) storeMultiStreamInfo(c context.Context, infos map[int64]*model.StreamFullInfo) {
	for _, v := range infos {
		key := d.getLocalStreamInfoKey(v.RoomID)

		//log.Warn("key=%v", key)
		d.localCache.SetWithExpire(key, v, _localStreamExpiredTime*time.Second)
	}
}

// loadMultiStreamInfo 读取批量信息
func (d *Dao) loadMultiStreamInfo(c context.Context, rids []int64) (map[int64]*model.StreamFullInfo, []int64) {
	infos := map[int64]*model.StreamFullInfo{}
	missRids := []int64{}
	for _, v := range rids {
		key := d.getLocalStreamInfoKey(v)
		info, _ := d.localCache.Get(key)
		if res, ok := info.(*model.StreamFullInfo); ok {
			if res != nil && res.RoomID > 0 {
				infos[v] = res
			}
		} else {
			missRids = append(missRids, v)
		}
	}

	return infos, missRids
}
