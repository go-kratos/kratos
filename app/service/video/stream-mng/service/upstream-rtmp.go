package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/conf"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"net/http"
	"net/url"
	"strings"
)

// 此处写具体的实现方式

// UpStreamRtmp 上行推流地址格式
type UpStreamRtmp struct {
	Addr    string `json:"addr,omitempty"`
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	NewLink string `json:"new_link,omitempty"`
}

// JudgeChangeRtmp 判断是否需要切换cdn参数
type JudgeChangeRtmp struct {
	RId        int64  `json:"rid,omitempty"`
	FreeFlow   string `json:"freeflow,omitempty"`
	IP         string `json:"ip,omitempty"`
	AreaID     int64  `json:"area_id,omitempty"`
	Attentions int    `json:"attentions,omitempty"`
	Uid        int64  `json:"uid,omitempty"`
	Platform   string `json:"platform,omitempty"`
}

// GetUpStreamRtmp UpStream
func (s *Service) GetUpStreamRtmp(c context.Context, rid int64, freeFlow string, ip string, areaID int64, attentions int, uid int64, platform string) (up *UpStreamRtmp, err error) {
	// 判断是否是该房间用户
	if uid != 0 {
		if uid != s.getRoomUserID(c, rid) {
			return nil, fmt.Errorf("无权限,不是该房间用户")
		}
	}

	judge := &JudgeChangeRtmp{
		RId:        rid,
		FreeFlow:   freeFlow,
		IP:         ip,
		AreaID:     areaID,
		Attentions: attentions,
		Platform:   platform,
	}

	log.Warn("%v", judge)

	sname, key, endUpStream, err := s.getEndUpStream(c, judge, false)

	if err != nil {
		return nil, err
	}

	// 获取url
	addr, code := s.generateUpStreamUrl(c, sname, endUpStream, key)

	newLink := s.formatNewLink(c, endUpStream, addr, code, false)
	if freeFlow == "unicom" {
		newLink = fmt.Sprintf("%s&unicom_free=1", newLink)
	}

	up = &UpStreamRtmp{
		Addr:    addr,
		Code:    code,
		NewLink: newLink,
	}

	return up, nil
}

// GetWebRtmp web端调用
func (s *Service) GetWebRtmp(c context.Context, rid int64, uid int64, ip string, platform string) (interface{}, error) {
	// 判断是否是该房间用户
	if uid != s.getRoomUserID(c, rid) {
		return nil, fmt.Errorf("无权限,不是该房间用户")
	}

	judge := &JudgeChangeRtmp{
		RId:        rid,
		FreeFlow:   "",
		IP:         ip,
		AreaID:     0,
		Attentions: 0,
		Platform:   platform,
	}
	sname, key, endUpStream, err := s.getEndUpStream(c, judge, true)

	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{}

	// 获取url
	addr, code := s.generateUpStreamUrl(c, sname, endUpStream, key)

	resp["rtmp"] = map[string]string{
		"addr": addr,
		"code": code,
	}

	streamLine := []map[string]interface{}{}

	streamLine = append(streamLine, map[string]interface{}{
		"name":     common.LineName[0],
		"src":      common.BitwiseMapSrc[endUpStream],
		"cdn_name": common.BitwiseMapName[endUpStream],
		"checked":  1,
	})

	// 获取是否是签约主播
	sign := s.isSignRoom(c, uid)

	// 测试房间
	testRoom := map[int64]int64{
		537499:   537499,
		11891462: 11891462,
	}

	ok := false
	_, ok = testRoom[rid]
	if sign || ok {
		for or, n := range common.BitwiseMapName {
			if or != endUpStream && or != common.BitWiseWS {
				streamLine = append(streamLine, map[string]interface{}{
					"name":     common.LineName[or],
					"src":      common.BitwiseMapSrc[or],
					"cdn_name": n,
					"checked":  0,
				})
			}
		}
	}

	resp["stream_line"] = streamLine

	return resp, nil
}

// GetRoomRtmp 拜年祭推流接口
func (s *Service) GetRoomRtmp(c context.Context, rid int64) (interface{}, error) {
	info, err := s.dao.StreamFullInfo(c, rid, "")
	if err != nil || info.RoomID <= 0 {
		return nil, err
	}

	type rtmp struct {
		Type       int             `json:"type"`
		StreamName string          `json:"stream_name"`
		Key        string          `json:"key,omitempty"`
		List       []*UpStreamRtmp `json:"list"`
	}

	resp := []rtmp{}

	// 只返回主流的一个推流码， 备用流的所有cdn推流码
	for _, v := range info.List {
		// 主流
		if v.Type == 1 {
			var or int64
			if v.Origin != 0 {
				or = v.Origin
			} else {
				or = v.DefaultUpStream
			}
			if or == 0 {
				return nil, fmt.Errorf("can not find upstream by room_id=%d", rid)
			}

			addr, code := s.generateUpStreamUrl(c, v.StreamName, or, v.Key)

			item := rtmp{
				Type:       1,
				StreamName: v.StreamName,
				Key:        v.Key,
			}

			item.List = append(item.List, &UpStreamRtmp{
				Addr: addr,
				Code: code,
				Name: common.NameMapBitwise[or],
			})

			resp = append(resp, item)
		} else {
			item := rtmp{
				Type:       2,
				StreamName: v.StreamName,
				Key:        v.Key,
			}

			for k, b := range common.ChinaNameMapBitwise {
				if b == common.BitWiseWS {
					continue
				}

				addr, code := s.generateUpStreamUrl(c, v.StreamName, b, v.Key)

				item.List = append(item.List, &UpStreamRtmp{
					Name: k,
					Addr: addr,
					Code: code,
				})
			}

			resp = append(resp, item)
		}
	}

	if len(resp) == 0 {
		return nil, nil
	}
	return resp, nil
}

func (s *Service) getEndUpStream(c context.Context, judge *JudgeChangeRtmp, web bool) (sname string, key string, end int64, err error) {
	info := &model.StreamFullInfo{}

	info, err = s.dao.StreamFullInfo(c, judge.RId, "")
	if err != nil {
		return
	}

	if info == nil || info.RoomID <= 0 {
		err = fmt.Errorf("获取房间信息失败")
		return
	}

	// 线路给默认上行的数据，而不是origin
	sname = ""
	var origin int64
	var endUpStream int64
	key = ""
	for _, v := range info.List {
		if v.Type == 1 {
			sname = v.StreamName
			if v.DefaultUpStream != 0 {
				origin = v.DefaultUpStream
				endUpStream = v.DefaultUpStream
			} else {
				origin = v.Origin
				endUpStream = v.Origin
			}
			key = v.Key
			break
		}
	}

	if sname == "" || key == "" || origin == 0 {
		errStr, _ := json.Marshal(info)
		log.Errorv(c, log.KV("log", string(errStr)))
		err = fmt.Errorf("获取房间信息失败, room_id=%d", judge.RId)
		return
	}

	// web 端调用， 不需要考虑免流和ismust
	fromSrc, toSrc, isMust, isChange, reason, dispatch := s.getChangeSrcRule(c, origin, judge.RId, judge.FreeFlow, judge.IP, judge.AreaID, judge.Attentions, web)
	if isChange {
		// 重新设置上行, 免流的必须切成功，其他无所谓
		err := s.dao.UpdateOfficialStreamStatus(c, judge.RId, common.BitwiseMapSrc[toSrc])
		if err == nil {
			endUpStream = toSrc

			go func(ctx context.Context, rid, toSrc, fromSrc int64, reason string, sname string) {
				s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
					RoomID:          rid,
					StreamName:      sname,
					DefaultUpStream: toSrc,
					DefaultChange:   true,
				})

				// 更新main-stream
				if err := s.dao.ChangeDefaultVendor(ctx, rid, toSrc); err != nil {
					log.Infov(ctx, log.KV("change_main_stream_default_err", err.Error()))
				}

				s.RecordChangeLog(ctx, &model.StreamChangeLog{
					RoomID:      rid,
					FromOrigin:  fromSrc,
					ToOrigin:    toSrc,
					Reason:      reason,
					OperateName: "auto_change",
					Source:      "background",
				})

				log.Infov(ctx, log.KV("log", fmt.Sprintf("auto_change_room:room_id=%d,from src=%d, tosrc=%d", rid, fromSrc, toSrc)))
			}(metadata.WithContext(c), judge.RId, toSrc, fromSrc, reason, sname)
		}
	} else {
		// 非web端必须切
		if isMust && !web {
			err = fmt.Errorf("auto_change_room failed:room_id=%d", judge.RId)
			return
		}
	}

	dispatch.PlatForm = judge.Platform
	dispatch.RoomID = judge.RId
	dispatch.CDN = endUpStream

	log.Warn("%v", dispatch)

	s.liveAside.Do(c, func(ctx context.Context) {
		err := s.dao.CreateUpStreamDispatch(ctx, dispatch)
		if err != nil {
			log.Errorv(ctx, log.KV("log", fmt.Sprintf("tidb_err:=%v", err)))
		}
	})

	return sname, key, endUpStream, nil
}

// getChangeSrcRule
func (s *Service) getChangeSrcRule(c context.Context, iLastCDN int64, rid int64, freeFlow string, ip string, areaID int64, attentions int, web bool) (fromOrigin int64, toOrigin int64, isMust bool, isChange bool, reason string, local *model.UpStreamInfo) {
	// 1.免流卡切至视频云(优先),下次如果不免流则切回上一次的CDN， 免流必须切成功
	// 2.签约主播不自动切 ----》 废弃
	// 3.是否被手动切过，是则不再做以下逻辑
	// 4.国内放映厅切至视频云
	// 5.海外切至腾讯
	// 6.按尾号切 --->废弃
	// 7.吃鸡分区,且关注数大于

	fromOrigin = iLastCDN

	// 判断用户所在区域
	local = s.GetLocation(c, ip)
	country := local.Country

	log.Infov(c, log.KV("log", fmt.Sprintf("ip=%s;country=%s", ip, country)))

	// 免流切到bvc
	if freeFlow == "unicom" {
		if iLastCDN != common.BitWiseBVC {
			// 设置上次的redis
			err := s.dao.UpdateLastCDNCache(c, rid, iLastCDN)
			if err != nil {
				log.Errorv(c, log.KV("log", fmt.Sprintf("set_last_cdn_error:%v", err)))
			}

			toOrigin = common.BitWiseBVC
			isMust = true
			isChange = true
		}

		return fromOrigin, toOrigin, isMust, isChange, "bilibili unicom card", local
	}

	thisconf := *conf.Conf
	if drop, ok := thisconf.DropCDN["dropCDN"].([]interface{}); ok {
		for _, v := range drop {
			if v.(int64) == fromOrigin {
				toOrigin = common.BitWiseBVC
				isChange = true
				isMust = true
				s.dao.DeleteLastCDNFromCache(c, rid)
				return fromOrigin, toOrigin, isMust, isChange, "drop cdn", local
			}
		}
	}

	// 切空某家上行到另一家
	// 此次是bvc, 上次不是bvc,切回原来的，并删除上次切cdn的记录
	if iLastCDN == common.BitWiseBVC {
		origin, err := s.dao.GetLastCDNFromCache(c, rid)
		if err != nil {
			return 0, 0, isMust, isChange, "", local
		}
		if origin > 0 {
			toOrigin = origin
			isChange = true
			s.dao.DeleteLastCDNFromCache(c, rid)
			return fromOrigin, toOrigin, isMust, isChange, "cdn is bvc", local
		}
	}

	// 被手动切过，是则不再做以下逻辑
	or, err := s.dao.GetChangeSrcFromCache(c, rid)
	if or != 0 && err == nil {
		s.dao.DeleteLastCDNFromCache(c, rid)
		return fromOrigin, toOrigin, isMust, isChange, "", local
	}

	// 国内吃鸡数大于1000给bvc
	preAttention := common.ChickenAttention
	if country == "中国" && areaID == common.AREAIDCHICKEN && attentions >= preAttention {
		if fromOrigin != common.BitWiseBVC {
			toOrigin = common.BitWiseBVC
			isChange = true
			s.dao.DeleteLastCDNFromCache(c, rid)
			return fromOrigin, toOrigin, isMust, isChange, "country is china and area is chicken", local
		}
	}

	// 海外节点给腾讯，无论web端还是移动端，均切换
	if country != "中国" && country != "局域网" && country != "本机地址" {
		if fromOrigin != common.BitWiseTC {
			toOrigin = common.BitWiseTC
			isChange = true
			s.dao.DeleteLastCDNFromCache(c, rid)
			return fromOrigin, toOrigin, isMust, isChange, "country is foreign country", local
		}
	}

	// bvc 移动端 切给qn，把尾号为2和5的的切给qn,先切20%
	if !web {
		if fromOrigin == common.BitWiseBVC && (rid%10 == 2 || rid%10 == 5) {
			toOrigin = common.BitWiseQN
			isChange = true
			s.dao.DeleteLastCDNFromCache(c, rid)
			return fromOrigin, toOrigin, isMust, isChange, "change bvc to qn", local
		}
	}

	// 全量切网宿, 3:2=tc:bvc
	if fromOrigin == common.BitWiseWS {
		if rid%10 < 5 {
			toOrigin = common.BitWiseTC
		} else {
			toOrigin = common.BitWiseBVC
		}
		isChange = true
		isMust = true
		s.dao.DeleteLastCDNFromCache(c, rid)
		return fromOrigin, toOrigin, isMust, isChange, fmt.Sprintf("change ws to %s", common.BitwiseMapName[toOrigin]), local
	}

	return 0, 0, isMust, isChange, "", local
}

// generateUpStreamUrl 生成推流url
func (s *Service) generateUpStreamUrl(c context.Context, sname string, origin int64, key string) (rtmp string, code string) {
	srcKw := common.BitwiseMapName[origin]
	if srcKw == common.WSName {
		rtmp = "rtmp://live-send.acg.tv/live"
		code = fmt.Sprintf("%s?streamname=%s&key=%s", sname, sname, key)
	} else {
		rtmp = fmt.Sprintf("rtmp://%s.live-send.acg.tv/live-%s/", srcKw, srcKw)
		code = fmt.Sprintf("?streamname=%s&key=%s", sname, key)
	}

	return
}

// formatNewLink
func (s *Service) formatNewLink(c context.Context, origin int64, addr string, code string, isHttps bool) string {
	upRtmp := ""
	if common.BitwiseMapName[origin] == common.WSName {
		upRtmp = fmt.Sprintf("%s/%s", addr, code)
	} else {
		upRtmp = fmt.Sprintf("%s%s", addr, code)
	}

	val := ""
	if isHttps {
		val = common.NewLinkMap[origin]["newLinkHttps"]
	} else {
		val = common.NewLinkMap[origin]["newLink"]
	}

	if val == "" {
		return val
	}

	// val = http://tcdns.myqcloud.com:8086/bilibili_redirect?up_rtmp=
	// upRtmp = rtmp://live-send.acg.tv/live/live_19148701_6447624?streamname=live_19148701_6447624&key=f2b23e1a938b6f979fa080f2bddc9225

	upRtmpEncode := url.QueryEscape(strings.Replace(upRtmp, "rtmp://", "", -1))

	newlink := fmt.Sprintf("%s%s", val, upRtmpEncode)

	return newlink
}

// getRoomUserID 获取一个房间的uid
func (s *Service) getRoomUserID(c context.Context, rid int64) int64 {
	bd := map[string]interface{}{
		"ids": []int64{rid},
	}

	bdj, _ := json.Marshal(bd)

	url := ""
	if env.DeployEnv == env.DeployEnvProd {
		url = "http://api.live.bilibili.co/room/v2/Room/get_by_ids"
	} else {
		url = "http://api.live.bilibili.com/room/v2/Room/get_by_ids"
	}

	type liveTimeStruct struct {
		Data map[string]map[string]interface{} `json:"data"`
	}

	liveResp := &liveTimeStruct{}

	head := map[string]string{
		"Content-Type": "application/json",
	}

	err := s.dao.NewRequst(c, http.MethodPost, url, nil, bdj, head, liveResp)
	if err != nil {
		log.Warn("%v", err)
		return 0
	}

	log.Infov(c, log.KV("log", fmt.Sprintf("url=%s;resp=%v", url, liveResp)))

	if id, ok := liveResp.Data[fmt.Sprintf("%d", rid)]["uid"].(float64); ok {
		return int64(id)
	}
	return 0
}

// 判断是否是签约的房间
func (s *Service) isSignRoom(c context.Context, rid int64) bool {
	if rid <= 0 {
		return false
	}
	url := "http://api.live.bilibili.co/live_user/v1/Anchor/get_badge"

	type signStatus struct {
		Data map[string]map[string]bool `json:"data,omitempty"`
	}

	var bd bytes.Buffer
	fmt.Fprint(&bd, "uids[]=", strings.Join([]string{fmt.Sprint(rid)}, ","))

	log.Warn("%s", bd.String())
	//bdj, _ := json.Marshal(bd)

	resp := &signStatus{}

	head := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err := s.dao.NewRequst(c, http.MethodPost, url, nil, bd.Bytes(), head, resp)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("%v=%v", err, resp)))
		return false
	}

	log.Infov(c, log.KV("log", fmt.Sprintf("%v", resp)))
	return resp.Data[fmt.Sprintf("%d", rid)]["is_signed_anchor"]
}
