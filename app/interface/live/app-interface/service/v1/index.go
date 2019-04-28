package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	v1indexpb "go-common/app/interface/live/app-interface/api/http/v1"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	liveuserV1 "go-common/app/service/live/live_user/api/liverpc/v1"
	relationV2 "go-common/app/service/live/relation/api/liverpc/v2"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	bannerV1 "go-common/app/service/live/room_ex/api/liverpc/v1"
	"go-common/app/service/live/third_api/bvc"
	userextV1 "go-common/app/service/live/userext/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/liverpc"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"math/rand"

	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"

	"go-common/library/net/http/blademaster"
)

const (
	_bannerType           = 1
	_navigatorType        = 2
	_yunyingRecFormType   = 3
	_yunyingRecSquareType = 4
	_recFormType          = 6
	_recSquareType        = 7
	_feedType             = 8
	_parentAreaFormType   = 9
	_parentAreaSquareType = 10
	_myAreaTagType        = 12
	_seaPatrolType        = 14
	_myAreaTagListType    = 13
	_activityType         = 11

	// _recTypeOnline   = 1
	// _recTypeIncome   = 2
	_recTypeForce    = 3
	_recTypeSkyHorse = 4

	_defaultRecNum = 24

	_skyHorseRecTimeOut = 100

	_areaModuleLink = "https://live.bilibili.com/app/area?parent_area_id=%d&parent_area_name=%s&area_id=%d&area_name=%s"

	_activityGo     = 0
	_activityBook   = 1
	_activityUnbook = 2

	_mobileIndexBadgeColorDefault = "#FB9E60"
)

// Service struct
type Service struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao         *dao.Dao
	allListInfo atomic.Value
}

type roomItem struct {
	RoomId           int64   `json:"roomid"`
	Title            string  `json:"title"`
	Uname            string  `json:"uname"`
	Online           int64   `json:"online"`
	Cover            string  `json:"cover"`
	Link             string  `json:"link"`
	Face             string  `json:"face"`
	AreaV2ParentId   int64   `json:"area_v2_parent_id"`
	AreaV2ParentName string  `json:"area_v2_parent_name"`
	AreaV2Id         int64   `json:"area_v2_id"`
	AreaV2Name       string  `json:"area_v2_name"`
	PlayUrl          string  `json:"play_url"`
	PlayUrlH265      string  `json:"play_url_h265"`
	CurrentQuality   int64   `json:"current_quality"`
	BroadcastType    int64   `json:"broadcast_type"`
	PendentRu        string  `json:"pendent_ru"`
	PendentRuPic     string  `json:"pendent_ru_pic"`
	PendentRuColor   string  `json:"pendent_ru_color"`
	RecType          int64   `json:"rec_type"`
	PkId             int64   `json:"pk_id"`
	AcceptQuality    []int64 `json:"accept_quality"`
}

type offlineItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type userTagItem struct {
	AreaV2Id         int    `json:"area_v2_id"`
	AreaV2Name       string `json:"area_v2_name"`
	AreaV2ParentId   int    `json:"area_v2_parent_id"`
	AreaV2ParentName string `json:"area_v2_parent_name"`
	Pic              string `json:"pic"`
	Link             string `json:"link"`
	IsAdvice         int    `json:"is_advice"`
}

type commonResp struct {
	ModuleInfo map[string]interface{}
	ExtraInfo  map[string]interface{}
	List       interface{}
}

type ModuleResp struct {
	Interval   int                      `json:"interval"`
	ModuleList []map[string]interface{} `json:"module_list"`
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  dao.New(c),
	}
	go s.tickCacheAllList(context.TODO())
	return s
}

// GetAllList implementation
// 首页大接口
func (s *Service) GetAllList(ctx context.Context, req *v1indexpb.GetAllListReq) (ret interface{}, err error) {
	resp := &ModuleResp{
		Interval: 10,
	}
	build := req.Build
	relationTimeout := conf.GetTimeout("relation", 200)
	// dao.LiveUserApi.V1UserSetting.GetTag(ctx, &liveuserV1.UserSettingGetTagReq{})
	midInterface, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64) // 大多使用header里的mid解析, 框架已封装请求的header
	isSkyHorseGray := false
	mid := int64(0)
	if isUIDSet {
		mid = midInterface
		// 天马灰度
		isSkyHorseGray = s.isSkyHorseRec(mid)
	}

	buvid := ""
	// 主站封好的，可从device里获取到sid、buvid、buvid3、build、channel、device、mobi_app、platform
	device, ok := metadata.Value(ctx, metadata.Device).(*blademaster.Device)
	if ok {
		buvid = device.Buvid
	}
	// deviceInterface := req.Device
	// device := deviceInterface.(*bm.Device)
	if req.Platform == "" || req.Device == "" || req.Scale == "" || req.RelationPage == 0 {
		err = errors.WithMessage(ecode.InvalidParam, "INVALID PARAM")
		return
	}

	allListTimeout := time.Duration(conf.GetTimeout("allList", 50)) * time.Millisecond
	rawModuleList := s.getAllListFromCache(rpcCtx.WithTimeout(ctx, allListTimeout))

	if rawModuleList == nil {
		err = errors.WithMessage(ecode.GetAllListReturnError, "")
		return
	}

	// 大分区常量定义
	parentName := map[int64]string{
		1: "娱乐",
		2: "游戏",
		3: "手游",
		4: "绘画",
		5: "电台",
	}

	// 天马灰度/保底
	defaultRecSlice := make([]map[string]interface{}, 0)
	loginRecRoomIDSlice := make([]int64, 0)
	// [{module_info:{},list:{}}...]
	resp.ModuleList = make([]map[string]interface{}, len(rawModuleList))
	for _, m := range rawModuleList {
		module := m.(map[string]interface{})
		if module["module_info"] == nil {
			log.Error("empty_module_info:%+v \n", m)
			fmt.Printf("empty_module_info:raw_all:%+v \n", rawModuleList)
		}
		moduleInfo := module["module_info"].(map[string]interface{})
		moduleType := jsonMustInt(moduleInfo["type"], 0)
		if moduleType == 0 {
			continue
		}
		list := module["list"].([]interface{})
		if moduleType == _recFormType || moduleType == _recSquareType {
			for _, r := range list {
				recItem := r.(map[string]interface{})
				defaultRecSlice = append(defaultRecSlice, recItem)
				roomID := jsonMustInt(recItem["roomid"], 0)
				if roomID == 0 {
					continue
				}
				// 登录了也有可能请求不到数据，登录的保底用
				loginRecRoomIDSlice = append(loginRecRoomIDSlice, roomID)
			}
		}
	}
	// 常用标签 roomListMap
	myTagRoomListMap := make(map[int64][]*roomV2.AppIndexGetMultiRoomListResp_RoomList)
	myTagAreaIds := make([]int64, 0, 4)
	myTagAreaInfoMap := make(map[int64]*liveuserV1.UserSettingGetTagResp_Tags)

	myTagResp := commonResp{
		List: make([]interface{}, 0),
	}
	attentionResp := commonResp{
		List: make([]interface{}, 0),
	}
	loginRecResp := commonResp{}
	seaResp := commonResp{}
	bannerResp := commonResp{
		List: make([]interface{}, 0),
	}
	// playurl定义
	attentionRoomListPlayURLMap := make(map[int64]*bvc.PlayUrlItem)
	loginRecRoomListPlayURLMap := make(map[int64]*bvc.PlayUrlItem)
	myTagRoomListPlayURLMap := make(map[int64]*bvc.PlayUrlItem)
	otherRoomListPlayURLMap := make(map[int64]*bvc.PlayUrlItem)
	otherRoomIDSlice := make([]int64, 0)

	isSkyHorseGrayOk := 0

	// 此group包含首页的一些任务
	// 但是任务之间不能同时cancel
	// 不然一个接口出错所有任务都cancel首页就空了
	// 所以return固定为nil(一个wg的任务使用的是从bm context继承来的ctx，cancel后一起推出)
	// 只有上层ctx（http的ctx）出问题（超时等）才会退出后续任务1
	wg, _ := errgroup.WithContext(ctx)
	for _, m := range rawModuleList {
		// {module_info:xx}
		module := m.(map[string]interface{})
		// {id:xx,type:xx,pic:xx,title:xx,link:xx,...}
		moduleInfo := module["module_info"].(map[string]interface{})
		moduleList := module["list"].([]interface{})
		moduleType := jsonMustInt(moduleInfo["type"], 0)
		if moduleType == 0 {
			continue
		}
		// banner分支 分端分版本
		if moduleType == _bannerType {
			bannerTimeout := time.Duration(conf.GetTimeout("banner", 100)) * time.Millisecond
			wg.Go(func() error {
				bannerList, bannerErr := dao.RoomExtApi.V1Banner.GetNewBanner(rpcCtx.WithTimeout(ctx, time.Duration(bannerTimeout)*time.Millisecond), &bannerV1.BannerGetNewBannerReq{UserPlatform: req.Platform, Build: build, UserDevice: req.Device})
				if bannerErr != nil {
					log.Error("[GetAllList]get banner rpc error, roomex.v1.Banner.GetNewBanner, error:%+v,rpctimeout:%d", bannerErr, bannerTimeout)
					return nil
				}
				if bannerList.Code != 0 || bannerList.Data == nil {
					log.Error("[GetAllList]get banner response error, code, %d, msg: %s, error:%+v", bannerList.Code, bannerList.Msg, bannerErr)
					return nil
				}
				if len(bannerList.Data) > 0 {
					for _, bannerInfo := range bannerList.Data {
						bannerResp.List = append(bannerResp.List.([]interface{}), map[string]interface{}{
							"id":    bannerInfo.Id,
							"pic":   bannerInfo.Pic,
							"link":  bannerInfo.Link,
							"title": bannerInfo.Title,
						})
					}
				}
				return nil
			})
		}
		// 关注分支
		if moduleType == _feedType {
			if !isUIDSet {
				continue
			}
			wg.Go(func() error {
				currentAttentionRoomMap := make(map[int64]bool)
				currentAttentionRoomSlice := make([]int64, 0)
				attentionResp.ModuleInfo = moduleInfo
				currentAttention, attentionErr := dao.RelationApi.V2App.LiveHomePage(rpcCtx.WithTimeout(ctx, time.Duration(relationTimeout)*time.Millisecond), &relationV2.AppLiveHomePageReq{RelationPage: req.RelationPage})
				if attentionErr != nil {
					log.Error("[GetAllList]get user attention rpc error, relation.v2.App.liveHomePage, error:%+v,rpctimeout:%d", attentionErr, relationTimeout)
				} else if currentAttention.Code != 0 || currentAttention.Data == nil {
					log.Error("[GetAllList]get user attention response error, code, %d, msg: %s, error:%+v", currentAttention.Code, currentAttention.Msg, attentionErr)
				} else {
					attentionResp.ExtraInfo = map[string]interface{}{
						"total_count":   currentAttention.Data.TotalCount,
						"time_desc":     currentAttention.Data.TimeDesc,
						"uname_desc":    currentAttention.Data.UnameDesc,
						"tags_desc":     currentAttention.Data.TagsDesc,
						"relation_page": currentAttention.Data.RelationPage,
					}
					// 存关注map
					for _, attentionCard := range currentAttention.Data.Rooms {
						currentAttentionRoomMap[attentionCard.Roomid] = true
						currentAttentionRoomSlice = append(currentAttentionRoomSlice, attentionCard.Roomid)
						attentionResp.List = append(attentionResp.List.([]interface{}), attentionCard)
					}
					// playurl
					attentionRoomListPlayURLMap = dao.BvcApi.GetPlayUrlMulti(ctx, currentAttentionRoomSlice, 0, 4, build, req.Platform)
				}
				useDefaultRec := false
				isOpen := conf.Conf.SkyHorseStatus
				if isOpen && isSkyHorseGray {
					// 取天马数据，传入关注当前刷列表roomid+强推，天马会对传入的roomid去重
					// duplicates := append(currentAttentionRoomSlice, forceRecSlice...)
					// getSkyHorseRoomList已经对强推去重
					duplicates := currentAttentionRoomSlice
					skyRecResult, skyHorseErr := getSkyHorseRoomList(ctx, mid, buvid, req.Build, req.Platform, duplicates, 1)
					if skyHorseErr != nil {
						log.Warn("[GetAllList]get data from skyHorse err: %v", skyHorseErr)
						useDefaultRec = true
					}
					if skyRecResult == nil {
						log.Warn("[GetAllList]get data from skyHorse empty: %v", skyHorseErr)
						useDefaultRec = true
					}
					if len(skyRecResult) < 6 {
						log.Warn("[GetAllList]get data from skyHorse not enough: %v", skyRecResult)
						useDefaultRec = true
					}
					skyRecResultInterface := make([]map[string]interface{}, 0)
					for _, item := range skyRecResult {
						loginRecRoomIDSlice = append(loginRecRoomIDSlice, item.RoomId)
						skyRecResultInterface = append(skyRecResultInterface, map[string]interface{}{
							"roomid":              item.RoomId,
							"title":               item.Title,
							"uname":               item.Uname,
							"online":              item.Online,
							"cover":               item.Cover,
							"link":                item.Link,
							"face":                item.Face,
							"area_v2_parent_id":   item.AreaV2ParentId,
							"area_v2_parent_name": item.AreaV2ParentName,
							"area_v2_id":          item.AreaV2Id,
							"area_v2_name":        item.AreaV2Name,
							"play_url":            item.PlayUrl,
							"current_quality":     item.CurrentQuality,
							"broadcast_type":      item.BroadcastType,
							"pendent_ru":          item.PendentRu,
							"pendent_ru_pic":      item.PendentRuPic,
							"pendent_ru_color":    item.PendentRuColor,
							"rec_type":            item.RecType,
							"pk_id":               item.PkId,
							"accept_quality":      item.AcceptQuality,
						})
					}
					loginRecResp.List = skyRecResultInterface
					isSkyHorseGrayOk = 1
				}
				// 保底逻辑
				if !isOpen || !isSkyHorseGray || useDefaultRec {
					newDefaultRecSlice := make([]map[string]interface{}, 0)
					for _, defaultRecRoom := range defaultRecSlice {
						roomID := jsonMustInt(defaultRecRoom["roomid"], 0)
						if roomID == 0 {
							continue
						}
						if _, exist := currentAttentionRoomMap[roomID]; !exist {
							// 天马没取到保底：默认推荐要对关注当前刷去重
							newDefaultRecSlice = append(newDefaultRecSlice, defaultRecRoom)
						}
					}
					// 只返24个
					if len(newDefaultRecSlice) > _defaultRecNum {
						loginRecResp.List = newDefaultRecSlice[:_defaultRecNum]
					} else {
						loginRecResp.List = newDefaultRecSlice
					}
					isSkyHorseGrayOk = 0
				}
				loginRecRoomListPlayURLMap = dao.BvcApi.GetPlayUrlMulti(ctx, loginRecRoomIDSlice, 0, 4, build, req.Platform)
				return nil
			})
		}

		// 常用标签分支
		if moduleType == _myAreaTagType {
			wg.Go(func() error {
				myTagResp.ModuleInfo = moduleInfo
				getMyTagTimeout := time.Duration(conf.GetTimeout("getMyTag", 100)) * time.Millisecond
				myTagListResp, userTagError := dao.LiveUserApi.V1UserSetting.GetTag(rpcCtx.WithTimeout(ctx, getMyTagTimeout), &liveuserV1.UserSettingGetTagReq{})
				if userTagError != nil {
					log.Error("[GetAllList]get user tag rpc error, live_user.v1.usersetting.get_tag, error:%+v", userTagError)
					return nil // 如果return err 则所有当前group的任务都会cancel
				}

				if myTagListResp.Code != 0 || myTagListResp.Data == nil {
					log.Error("[GetAllList]get user tag return error, code, %d, msg: %s, error:%+v", myTagListResp.Code, myTagListResp.Msg, userTagError)
					return nil
				}
				if myTagListResp.Data != nil {
					myTagResp.ExtraInfo = make(map[string]interface{})
					myTagResp.ExtraInfo["is_gray"] = myTagListResp.Data.IsGray
					myTagResp.ExtraInfo["offline"] = make([]interface{}, 0)
					for _, offlineInfo := range myTagListResp.Data.Offline {
						myTagResp.ExtraInfo["offline"] = append(myTagResp.ExtraInfo["offline"].([]interface{}), &offlineItem{Id: int(offlineInfo.Id), Name: offlineInfo.Name})
					}
				}

				for _, tagInfo := range myTagListResp.Data.Tags {
					myTagAreaIds = append(myTagAreaIds, tagInfo.Id)
					myTagAreaInfoMap[tagInfo.Id] = tagInfo
					link := fmt.Sprintf("http://live.bilibili.com/app/area?parent_area_id=%d&parent_area_name=%s&area_id=%d&area_name=%s", tagInfo.ParentId, parentName[tagInfo.ParentId], tagInfo.Id, tagInfo.Name)
					myTagResp.List = append(myTagResp.List.([]interface{}), &userTagItem{AreaV2Id: int(tagInfo.Id), AreaV2Name: tagInfo.Name, AreaV2ParentId: int(tagInfo.ParentId), AreaV2ParentName: parentName[tagInfo.ParentId], Link: link, Pic: tagInfo.Pic, IsAdvice: int(tagInfo.IsAdvice)})
				}
				if (req.Platform == "ios" && build > 8220) || (req.Platform == "android" && build > 5333002) {
					myTagResp.List = append(myTagResp.List.([]interface{}), &userTagItem{AreaV2Id: int(0), AreaV2Name: "全部标签", AreaV2ParentId: int(0), AreaV2ParentName: "", Pic: "http://i0.hdslb.com/bfs/vc/ff03528785fc8c91491d79e440398484811d6d87.png", Link: "http://live.bilibili.com/app/mytag/", IsAdvice: 1})
				}
				if len(myTagAreaIds) <= 0 {
					log.Error("[GetAllList]get user tag return empty!uid:%d", mid)
					return nil
				}
				// 常用标签房间列表 先生成最后wait替换就好了
				getMultiRoomListTimeout := time.Duration(conf.GetTimeout("getMultiRoomList", 100)) * time.Millisecond
				myTagRoomListResp, multiRoomListErr := dao.RoomApi.V2AppIndex.GetMultiRoomList(rpcCtx.WithTimeout(ctx, getMultiRoomListTimeout), &roomV2.AppIndexGetMultiRoomListReq{AreaIds: xstr.JoinInts(myTagAreaIds), Platform: req.Platform})

				if multiRoomListErr != nil {
					log.Error("[GetAllList]get multi list rpc error, room.v2.AppIndex.GetMultiRoomList, error:%+v", multiRoomListErr)
					return nil
				}
				if myTagRoomListResp.Code != 0 || myTagRoomListResp.Data == nil {
					log.Error("[GetAllList]get multi list response error, code, %d, msg: %s, error:%+v", myTagRoomListResp.Code, myTagRoomListResp.Msg, multiRoomListErr)
					return nil
				}
				// 保存roomListMap，wait 聚合数据
				myTagRoomIDSlice := make([]int64, 0)
				for _, myTagRoomItem := range myTagRoomListResp.Data {
					myTagRoomListMap[myTagRoomItem.Id] = myTagRoomItem.List
					for _, item := range myTagRoomItem.List {
						myTagRoomIDSlice = append(myTagRoomIDSlice, item.Roomid)
					}
				}
				myTagRoomListPlayURLMap = dao.BvcApi.GetPlayUrlMulti(ctx, myTagRoomIDSlice, 0, 4, build, req.Platform)

				return nil
			})
		}

		if moduleType == _seaPatrolType {
			seaPatrolList := make([]interface{}, 0)
			if !isUIDSet {
				continue
			}
			// 大航海分支
			wg.Go(func() error {
				seaPatrolTimeout := time.Duration(conf.GetTimeout("seaPatrol", 100)) * time.Millisecond
				seaPatrol, seaPatrolError := dao.LiveUserApi.V1Note.Get(rpcCtx.WithTimeout(ctx, seaPatrolTimeout), &liveuserV1.NoteGetReq{})
				if seaPatrolError != nil {
					log.Error("[GetAllList]get sea patrol rpc error, liveuser.v1.Note.Get, error:%+v", seaPatrolError)
					return nil
				}
				if seaPatrol.Code != 0 || seaPatrol.Data == nil {
					log.Error("[GetAllList]get sea patrol note from liveuser response error, code, %d, msg: %s, error:%+v", seaPatrol.Code, seaPatrol.Msg, seaPatrolError)
					return nil
				}

				if seaPatrol.Data.Title != "" {
					seaPatrolList = append(seaPatrolList, map[string]interface{}{
						"pic":     seaPatrol.Data.Logo,
						"title":   seaPatrol.Data.Title,
						"link":    seaPatrol.Data.Link,
						"content": seaPatrol.Data.Content,
					})
				}
				seaResp.List = seaPatrolList
				seaResp.ModuleInfo = moduleInfo
				seaResp.ExtraInfo = make(map[string]interface{})
				return nil
			})
		}

		if moduleType == _activityType {
			cardList := module["list"].([]interface{})
			actyInfo := cardList[0].(map[string]interface{})
			bookStatus := jsonMustInt(actyInfo["status"], 0)
			// status=0(非预约类型的活动)
			if bookStatus == _activityGo {
				actyInfo["button_text"] = "去围观"
				actyInfo["status"] = _activityGo
				continue
			}
			// 未登入 显示预约
			if !isUIDSet {
				actyInfo["button_text"] = "预约"
				actyInfo["status"] = _activityBook
				continue
			}
			// 登入状态 设置保底数据
			actyInfo["button_text"] = "去围观"
			actyInfo["status"] = _activityGo
			// 获取活动id
			materialID := jsonMustInt(moduleInfo["material_id"], 0)
			if materialID == 0 {
				continue
			}
			log.Info("[GetAllList]materialID is %v", materialID)
			wg.Go(func() error {
				activityQueryTimeout := time.Duration(conf.GetTimeout("activityQuery", 100)) * time.Millisecond
				bookInfo, userExtError := dao.UserExtApi.V1Remind.Query(rpcCtx.WithTimeout(ctx, activityQueryTimeout), &userextV1.RemindQueryReq{Aid: materialID})
				if userExtError != nil {
					log.Error("[GetAllList]get activity book info rpc error, userext.v1.Remind.Query, error:%+v", userExtError)
					return nil
				}
				if bookInfo.Code != 0 {
					log.Error("[GetAllList]get activity book info response error, code, %d, msg: %s, error:%+v", bookInfo.Code, bookInfo.Msg, userExtError)
					return nil
				}
				log.Info("[GetAllList]materialID is %v and bookInfo.Data.Status is %v", materialID, bookInfo.Data.Status)
				switch bookInfo.Data.Status {
				case _activityBook:
					actyInfo["button_text"] = "已预约"
					actyInfo["status"] = _activityUnbook
				case _activityUnbook:
					actyInfo["button_text"] = "预约"
					actyInfo["status"] = _activityBook
				default:
					actyInfo["button_text"] = "去围观"
					actyInfo["status"] = _activityGo
				}
				return nil
			})
		}

		// 其他playurl，注意这里取的推荐是未登录下的推荐play_url
		if moduleType == _yunyingRecFormType || moduleType == _yunyingRecSquareType ||
			moduleType == _parentAreaFormType || moduleType == _parentAreaSquareType ||
			moduleType == _recFormType || moduleType == _recSquareType {
			// append Roomid
			for _, item := range moduleList {
				itemV := item.(map[string]interface{})
				roomID := jsonMustInt(itemV["roomid"], 0)
				if roomID == 0 {
					continue
				}
				otherRoomIDSlice = append(otherRoomIDSlice, roomID)
			}
		}
	}
	// +其他模块playurl
	wg.Go(func() error {
		otherRoomListPlayURLMap = dao.BvcApi.GetPlayUrlMulti(ctx, otherRoomIDSlice, 0, 4, build, req.Platform)
		return nil
	})

	waitErr := wg.Wait()
	if waitErr != nil {
		log.Error("[GetAllList]wait error: %s", waitErr)
		return
	}

	// 封装
	tagIndex := 0
	for index, m := range rawModuleList {
		module := m.(map[string]interface{})
		moduleInfo := module["module_info"].(map[string]interface{})
		moduleList := module["list"].([]interface{})
		moduleType := jsonMustInt(moduleInfo["type"], 0)
		if moduleType == 0 {
			continue
		}
		// 初始化
		resp.ModuleList[index] = make(map[string]interface{})
		resp.ModuleList[index]["list"] = moduleList
		resp.ModuleList[index]["module_info"] = moduleInfo

		if moduleType == _bannerType {
			resp.ModuleList[index]["list"] = bannerResp.List
		}
		if moduleType == _navigatorType && req.Platform == "android" && build <= 5333002 {
			// 分区入口5.33版本还返回4个(前3个+全部)，5.34透传后台的5个
			if len(moduleList) > 3 {
				resp.ModuleList[index]["list"] = append(moduleList[:3], map[string]interface{}{
					"id":    12,
					"pic":   "https://i0.hdslb.com/bfs/vc/ff03528785fc8c91491d79e440398484811d6d87.png",
					"link":  "https://live.bilibili.com/app/mytag/",
					"title": "全部标签",
				})
			}
		}
		if moduleType == _seaPatrolType {
			if seaResp.List != nil {
				resp.ModuleList[index]["list"] = seaResp.List
			}
			if seaResp.ModuleInfo != nil {
				resp.ModuleList[index]["module_info"] = seaResp.ModuleInfo
			}
			if seaResp.ExtraInfo != nil {
				resp.ModuleList[index]["extra_info"] = seaResp.ExtraInfo
			}
		}
		if moduleType == _myAreaTagType {
			if myTagResp.List != nil {
				resp.ModuleList[index]["list"] = myTagResp.List
			}
			if myTagResp.ModuleInfo != nil {
				resp.ModuleList[index]["module_info"] = myTagResp.ModuleInfo
			}
			if myTagResp.ExtraInfo != nil {
				resp.ModuleList[index]["extra_info"] = myTagResp.ExtraInfo
			}
		}
		var isTagGray int
		iTmp, _ := myTagResp.ExtraInfo["is_gray"].(int64)
		isTagGray = int(iTmp)
		// 常用分区房间列表填充
		if moduleType == _myAreaTagListType {
			if isTagGray == 0 {
				continue
			}
			if len(myTagAreaIds) == 0 || tagIndex >= len(myTagAreaIds) {
				continue
			}
			mTagAreaID := myTagAreaIds[tagIndex]
			if _, ok := myTagRoomListMap[mTagAreaID]; ok {
				for _, v := range myTagRoomListMap[mTagAreaID] {
					if myTagRoomListPlayURLMap[v.Roomid] != nil {
						v.AcceptQuality = myTagRoomListPlayURLMap[v.Roomid].AcceptQuality
						v.CurrentQuality = myTagRoomListPlayURLMap[v.Roomid].CurrentQuality
						v.PlayUrl = myTagRoomListPlayURLMap[v.Roomid].Url["h264"]
						v.PlayUrlH265 = myTagRoomListPlayURLMap[v.Roomid].Url["h265"]
					}
				}
				resp.ModuleList[index]["list"] = myTagRoomListMap[mTagAreaID]
				if _, ok := myTagAreaInfoMap[mTagAreaID]; ok {
					areaInfo := myTagAreaInfoMap[mTagAreaID]
					moduleInfo["title"] = areaInfo.Name
					moduleInfo["link"] = fmt.Sprintf(_areaModuleLink, areaInfo.ParentId, parentName[areaInfo.ParentId], areaInfo.Id, areaInfo.Name)
					resp.ModuleList[index]["module_info"] = moduleInfo
				}
			}
			tagIndex++
		}

		// 运营推荐分区对常用分区去重
		if moduleType == _yunyingRecFormType || moduleType == _yunyingRecSquareType {
			link := moduleInfo["link"]
			u, err := url.Parse(link.(string))
			if err != nil {
				log.Warn("[GetAllList]url.Parse (%s) error: %v", link, err)
				continue
			}
			m, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				log.Warn("[GetAllList]url.ParseQuery (%s) error: %v", link, err)
				continue
			}
			area, ok := m["area_id"]
			if !ok {
				log.Warn("[GetAllList]url ((%s) area_id lost: %v", link, ok)
				continue
			}

			trueArea, err := strconv.Atoi(area[0])
			if err != nil {
				log.Warn("[GetAllList]get trueAreaId error: %v", link, ok)
				continue
			}
			if _, ok := myTagRoomListMap[int64(trueArea)]; ok && isTagGray == 1 {
				resp.ModuleList[index]["list"] = nil
			} else {
				for _, v := range moduleList {
					if v == nil {
						continue
					}
					vv := v.(map[string]interface{})
					roomID := jsonMustInt(vv["roomid"], 0)
					if roomID == 0 {
						continue
					}
					if otherRoomListPlayURLMap[roomID] != nil {
						vv["accept_quality"] = otherRoomListPlayURLMap[roomID].AcceptQuality
						vv["current_quality"] = otherRoomListPlayURLMap[roomID].CurrentQuality
						vv["play_url"] = otherRoomListPlayURLMap[roomID].Url["h264"]
						vv["play_url_h265"] = otherRoomListPlayURLMap[roomID].Url["h265"]
					}
				}

				resp.ModuleList[index]["list"] = moduleList
				resp.ModuleList[index]["module_info"] = moduleInfo
			}
		}

		if moduleType == _feedType {
			if attentionResp.ModuleInfo != nil {
				resp.ModuleList[index]["module_info"] = attentionResp.ModuleInfo
			}
			if attentionResp.ExtraInfo != nil {
				resp.ModuleList[index]["extra_info"] = attentionResp.ExtraInfo
			}
			if attentionResp.List != nil {
				for _, v := range attentionResp.List.([]interface{}) {
					vv := v.(*relationV2.AppLiveHomePageResp_Rooms)
					if attentionRoomListPlayURLMap[vv.Roomid] != nil {
						vv.AcceptQuality = attentionRoomListPlayURLMap[vv.Roomid].AcceptQuality
						vv.CurrentQuality = attentionRoomListPlayURLMap[vv.Roomid].CurrentQuality
						vv.PlayUrl = attentionRoomListPlayURLMap[vv.Roomid].Url["h264"]
						vv.PlayUrlH265 = attentionRoomListPlayURLMap[vv.Roomid].Url["h265"]
					}
				}
				resp.ModuleList[index]["list"] = attentionResp.List
			}
		}

		if moduleType == _recFormType || moduleType == _recSquareType {
			moduleInfo["is_sky_horse_gray"] = isSkyHorseGrayOk
			resp.ModuleList[index]["module_info"] = moduleInfo
			if loginRecResp.List != nil {
				// is uid set
				for _, v := range loginRecResp.List.([]map[string]interface{}) {
					roomID := int64(0)
					r, ok := v["roomid"].(json.Number)
					if ok {
						rr, intErr := r.Int64()
						if intErr != nil {
							continue
						}
						roomID = rr
					} else {
						roomID = v["roomid"].(int64)
					}
					if roomID == 0 {
						continue
					}

					if loginRecRoomListPlayURLMap[roomID] != nil {
						v["accept_quality"] = loginRecRoomListPlayURLMap[roomID].AcceptQuality
						v["current_quality"] = loginRecRoomListPlayURLMap[roomID].CurrentQuality
						v["play_url"] = loginRecRoomListPlayURLMap[roomID].Url["h264"]
						v["play_url_h265"] = loginRecRoomListPlayURLMap[roomID].Url["h265"]
					}
				}
				resp.ModuleList[index]["list"] = loginRecResp.List
			} else {
				for _, v := range moduleList {
					if v == nil {
						continue
					}
					vv := v.(map[string]interface{})
					roomID := jsonMustInt(vv["roomid"], 0)
					if roomID == 0 {
						continue
					}
					if otherRoomListPlayURLMap[roomID] != nil {
						vv["accept_quality"] = otherRoomListPlayURLMap[roomID].AcceptQuality
						vv["current_quality"] = otherRoomListPlayURLMap[roomID].CurrentQuality
						vv["play_url"] = otherRoomListPlayURLMap[roomID].Url["h264"]
						vv["play_url_h265"] = otherRoomListPlayURLMap[roomID].Url["h265"]
					}
				}
				// 只返24个，新推荐已在上面做处理
				if len(moduleList) > _defaultRecNum {
					resp.ModuleList[index]["list"] = moduleList[:_defaultRecNum]
				} else {
					resp.ModuleList[index]["list"] = moduleList
				}

			}
		}

		if moduleType == _parentAreaFormType || moduleType == _parentAreaSquareType {
			for _, v := range moduleList {
				if v == nil {
					continue
				}
				vv := v.(map[string]interface{})
				roomID := jsonMustInt(vv["roomid"], 0)
				if roomID == 0 {
					continue
				}
				if otherRoomListPlayURLMap[roomID] != nil {
					vv["accept_quality"] = otherRoomListPlayURLMap[roomID].AcceptQuality
					vv["current_quality"] = otherRoomListPlayURLMap[roomID].CurrentQuality
					vv["play_url"] = otherRoomListPlayURLMap[roomID].Url["h264"]
					vv["play_url_h265"] = otherRoomListPlayURLMap[roomID].Url["h265"]
				}
			}
			resp.ModuleList[index]["list"] = moduleList
		}
	}
	ret = resp
	return
}

func jsonMustInt(arg interface{}, def int64) int64 {
	if arg == nil {
		log.Warn("jsonMustInt arg(%v) nil!", arg)
		return def
	}
	r, ok := arg.(json.Number)
	if !ok {
		log.Warn("jsonMustInt arg(%v) is not json.Number but %v", arg, reflect.TypeOf(arg))
		return def
	}
	rr, err := r.Int64()
	if err != nil {
		log.Warn("jsonMustInt arg(%v) transfer error: %v", arg, err)
		return def
	}
	return rr
}

// Change implementation
// 首页换一换接口 for 天马
func (s *Service) Change(ctx context.Context, req *v1indexpb.ChangeReq) (resp *v1indexpb.ChangeResp, err error) {
	resp = &v1indexpb.ChangeResp{
		ModuleList: make([]*v1indexpb.ChangeResp_ModuleList, 0),
	}
	mid, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64)

	var uid int64

	if isUIDSet {
		uid = mid
	}
	// deviceInterface, _ := ctx.Get("device")
	// device := req.Device

	duplicates, _ := xstr.SplitInts(req.AttentionRoomId)
	duplicatesMap := make(map[int64]bool)

	for _, roomID := range duplicates {
		duplicatesMap[roomID] = true
	}
	build := req.Build
	buvid := ""
	// 主站封好的，可从device里获取到sid、buvid、buvid3、build、channel、device、mobi_app、platform
	device, ok := metadata.Value(ctx, metadata.Device).(*blademaster.Device)
	if ok {
		buvid = device.Buvid
	}
	var recModuleInfo *v1indexpb.ChangeResp_ModuleInfo
	allListOut, callErr := dao.RoomApi.V2AppIndex.GetAllList(rpcCtx.WithTimeout(ctx, 100*time.Millisecond), &roomV2.AppIndexGetAllListReq{
		Platform: req.Platform,
		Device:   req.Device,
		Scale:    req.Scale,
		Build:    int64(build),
		ModuleId: req.ModuleId,
	})

	if callErr != nil {
		log.Error("[Change]get all list rpc error, room.v2.AppIndex.setAllList, error:%+v", callErr)
		err = errors.WithMessage(ecode.ChangeGetAllListRPCError, "CHANGE GET ALL LIST FAIL#1")
		return
	}

	if allListOut.Code != 0 || allListOut.Data == nil {
		log.Error("[Change]get all list return data error, code, %d, msg: %s, error:%+v", allListOut.Code, allListOut.Msg, err)
		err = errors.WithMessage(ecode.ChangeGetAllListReturnError, "CHANGE GET ALL LIST FAIL#2")
		return
	}

	if len(allListOut.Data.ModuleList) == 0 {
		log.Error("[Change]get all list return empty, code, %d, msg: %s, error:%+v", allListOut.Code, allListOut.Msg, err)
		err = errors.WithMessage(ecode.ChangeGetAllListEmptyError, "CHANGE GET ALL LIST FAIL#3")
		return
	}

	m := allListOut.Data.ModuleList[0]
	duplicateList := make([]*roomV2.AppIndexGetAllListResp_RoomList, 0)
	for _, itemInfo := range m.List {
		if _, ok := duplicatesMap[itemInfo.Roomid]; !ok {
			duplicateList = append(duplicateList, itemInfo)
		}
	}

	m.List = duplicateList

	recModuleInfo = &v1indexpb.ChangeResp_ModuleInfo{
		Id:             m.ModuleInfo.Id,
		Title:          m.ModuleInfo.Title,
		Pic:            m.ModuleInfo.Pic,
		Type:           m.ModuleInfo.Type,
		Link:           m.ModuleInfo.Link,
		Count:          m.ModuleInfo.Count,
		IsSkyHorseGray: 0,
	}
	// 目前只有推荐-天马有换一换，都必须有登陆态

	roomIds := make([]int64, 0)
	list := make([]*v1indexpb.ChangeResp_List, 0)
	for i, itemInfo := range m.List {
		if i >= 24 {
			break
		}
		roomIds = append(roomIds, itemInfo.Roomid)
		list = append(list, &v1indexpb.ChangeResp_List{
			Roomid:           itemInfo.Roomid,
			Title:            itemInfo.Title,
			Uname:            itemInfo.Uname,
			Online:           itemInfo.Online,
			Cover:            itemInfo.Cover,
			Link:             "/" + strconv.Itoa(int(itemInfo.Roomid)),
			Face:             itemInfo.Face,
			AreaV2ParentId:   itemInfo.AreaV2ParentId,
			AreaV2ParentName: itemInfo.AreaV2ParentName,
			AreaV2Id:         itemInfo.AreaV2Id,
			AreaV2Name:       itemInfo.AreaV2Name,
			BroadcastType:    itemInfo.BroadcastType,
			PendentRu:        itemInfo.PendentRu,
			PendentRuPic:     itemInfo.PendentRuPic,
			PendentRuColor:   itemInfo.PendentRuColor,
			RecType:          itemInfo.RecType,
			CurrentQuality:   itemInfo.CurrentQuality,
			AcceptQuality:    itemInfo.AcceptQuality,
			PlayUrl:          itemInfo.PlayUrl,
		})
	}

	resp.ModuleList = append(resp.ModuleList, &v1indexpb.ChangeResp_ModuleList{
		List:       list,
		ModuleInfo: recModuleInfo,
	})

	isOpen := conf.Conf.SkyHorseStatus
	if isOpen && isUIDSet && s.isSkyHorseRec(uid) {
		recPage := rand.Intn(4)
		if recPage == 1 || recPage == 0 {
			recPage = 2
		}
		recList, skyHorseErr := getSkyHorseRoomList(ctx, uid, buvid, req.Build, req.Platform, duplicates, int64(recPage))
		if skyHorseErr != nil {
			log.Error("[Change]getSkyHorseRoomList error:%+v", skyHorseErr)
			// err = errors.WithMessage(ecode.SkyHorseError, "")
		} else if len(recList) <= 0 {
			log.Error("[Change]getSkyHorseRoomList empty:%+v", recList)
			// err = errors.WithMessage(ecode.ChangeSkyHorseEmptyError, "")
		} else {
			list := make([]*v1indexpb.ChangeResp_List, 0)
			for i, recInfo := range recList {
				if i >= 6 {
					continue
				}
				roomIds = append(roomIds, recInfo.RoomId)
				list = append(list, &v1indexpb.ChangeResp_List{
					Roomid:           recInfo.RoomId,
					Title:            recInfo.Title,
					Uname:            recInfo.Uname,
					Online:           recInfo.Online,
					Cover:            recInfo.Cover,
					Link:             "/" + strconv.Itoa(int(recInfo.RoomId)),
					Face:             recInfo.Face,
					AreaV2ParentId:   recInfo.AreaV2ParentId,
					AreaV2ParentName: recInfo.AreaV2ParentName,
					AreaV2Id:         recInfo.AreaV2Id,
					AreaV2Name:       recInfo.AreaV2Name,
					BroadcastType:    recInfo.BroadcastType,
					PendentRu:        recInfo.PendentRu,
					PendentRuPic:     recInfo.PendentRuPic,
					PendentRuColor:   recInfo.PendentRuColor,
					RecType:          _recTypeSkyHorse,
				})
			}
			skyHorseList := make([]*v1indexpb.ChangeResp_ModuleList, 0)
			recModuleInfo.IsSkyHorseGray = 1
			skyHorseList = append(skyHorseList, &v1indexpb.ChangeResp_ModuleList{
				List:       list,
				ModuleInfo: recModuleInfo,
			})
			resp.ModuleList = skyHorseList
		}
	}
	changeRoomListPlayURLMap := dao.BvcApi.GetPlayUrlMulti(ctx, roomIds, 0, 4, build, req.Platform)

	for _, v := range resp.ModuleList[0].List {
		if changeRoomListPlayURLMap[v.Roomid] != nil {
			v.AcceptQuality = changeRoomListPlayURLMap[v.Roomid].AcceptQuality
			v.CurrentQuality = changeRoomListPlayURLMap[v.Roomid].CurrentQuality
			v.PlayUrl = changeRoomListPlayURLMap[v.Roomid].Url["h264"]
			v.PlayUrlH265 = changeRoomListPlayURLMap[v.Roomid].Url["h265"]
		}
	}
	// 赋值
	return
}

// 获取天马房间信息列表
// 已将强推roomids传给天马，其他的可通过传duplicates来merge进去
func getSkyHorseRoomList(ctx context.Context, uid int64, buvid string, build int64, platform string, duplicates []int64, recPage int64) (resp []*roomItem, err error) {
	clientRecStrongTimeout := time.Duration(conf.GetTimeout("clientRecStrong", 100)) * time.Millisecond
	strongRecList, strongRecErr := dao.RoomApi.V1RoomRecommend.ClientRecStrong(rpcCtx.WithTimeout(ctx, clientRecStrongTimeout), &roomV1.RoomRecommendClientRecStrongReq{RecPage: recPage})
	// liverpc.NewClient().CallRaw()
	recDuplicate := make([]int64, 0)
	if strongRecErr != nil {
		log.Error("[getSkyHorseRoomList]room.v1.ClientRecStrong rpc error:%+v", strongRecErr)
	} else if strongRecList.Code != 0 {
		log.Error("[getSkyHorseRoomList]room.v1.ClientRecStrong response error:%+v,code:%d,msg:%s", strongRecErr, strongRecList.Code, strongRecList.Msg)
	} else {
		for _, strongInfo := range strongRecList.Data.Result {
			if strongInfo.Roomid == 0 {
				continue
			}
			recDuplicate = append(recDuplicate, strongInfo.Roomid)
		}
	}

	strongLen := len(recDuplicate)

	duplicates = append(duplicates, recDuplicate...)
	skyHorseRec, skyHorseErr := dao.SkyHorseApi.GetSkyHorseRec(ctx, uid, buvid, build, platform, duplicates, strongLen, _skyHorseRecTimeOut)
	if skyHorseErr != nil {
		err = errors.WithMessage(ecode.SkyHorseError, "")
		return
	}

	roomIds := make([]int64, 0)

	for _, skyHorseInfo := range skyHorseRec.Data {
		roomIds = append(roomIds, int64(skyHorseInfo.Id))
	}

	indexRoomListFields := []string{
		"roomid",
		"title",
		"uname",
		"online",
		"cover",
		"user_cover",
		"link",
		"face",
		"area_v2_parent_id",
		"area_v2_parent_name",
		"area_v2_id",
		"area_v2_name",
		"broadcast_type",
		"uid",
	}

	wg, _ := errgroup.WithContext(ctx)

	// 房间基础信息（是map，但是天马返回是无序的）
	var multiRoomListResp *roomV2.RoomGetByIdsResp
	wg.Go(func() error {
		getByIdsTimeout := time.Duration(conf.GetTimeout("getByIds", 50)) * time.Millisecond
		multiRoomList, getByIdsError := dao.RoomApi.V2Room.GetByIds(rpcCtx.WithTimeout(ctx, getByIdsTimeout), &roomV2.RoomGetByIdsReq{
			Ids:               roomIds,
			NeedBroadcastType: 1,
			NeedUinfo:         1,
			Fields:            indexRoomListFields,
			From:              "app-interface.gateway",
		})
		if getByIdsError != nil {
			log.Error("[getSkyHorseRoomList]room.v2.getByIds rpc error:%+v", getByIdsError)
			// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
			return errors.WithMessage(ecode.GetRoomError, "room.v2.getByIds rpc error")
		}
		if multiRoomList.Code != 0 {
			log.Error("[getSkyHorseRoomList]room.v2.getByIds response error:%+v,code:%d,msg:%s", getByIdsError, multiRoomList.Code, multiRoomList.Msg)
			// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
			return errors.WithMessage(ecode.GetRoomError, "room.v2.getByIds response error")
		}
		multiRoomListResp = multiRoomList

		return nil
	})

	// 房间角标信息
	pendantRoomListResp := &roomV1.RoomPendantGetPendantByIdsResp{}
	wg.Go(func() error {
		getPendantByIdsTimeout := time.Duration(conf.GetTimeout("getPendantByIds", 50)) * time.Millisecond
		pendantRoomList, getPendantError := dao.RoomApi.V1RoomPendant.GetPendantByIds(rpcCtx.WithTimeout(ctx, getPendantByIdsTimeout), &roomV1.RoomPendantGetPendantByIdsReq{
			Ids:      roomIds,
			Type:     "mobile_index_badge",
			Position: 2, // 历史原因，取右上，但客户端展示在左上
		})

		if getPendantError != nil {
			log.Error("[getSkyHorseRoomList]room.v1.getPendantByIds rpc error:%+v", getPendantError)
			return nil
		}
		if pendantRoomList.Code != 0 {
			log.Error("[getSkyHorseRoomList]room.v1.getPendantByIds response error:%+v,code:%d,msg:%s", getPendantError, pendantRoomList.Code, pendantRoomList.Msg)
			return nil
		}
		pendantRoomListResp = pendantRoomList
		return nil
	})

	waitErr := wg.Wait()
	if waitErr != nil {
		log.Error("[getSkyHorseRoomList]wait error(%+v)", waitErr)
		return
	}

	pendantResult := make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	// 天马返回是无序的
	if multiRoomListResp == nil {
		err = errors.WithMessage(ecode.GetRoomEmptyError, "")
		return
	}

	respSlice := make([]*roomV2.RoomGetByIdsResp_RoomInfo, 0)
	for _, roomBaseInfo := range multiRoomListResp.Data {
		respSlice = append(respSlice, roomBaseInfo)
	}
	for i := 0; i < 6; i++ {
		if strongRecList != nil && strongRecList.Data != nil && strongRecList.Data.Result != nil {
			if recInfo, ok := strongRecList.Data.Result[int64(i)]; ok {
				resp = append(resp, &roomItem{
					RoomId:           recInfo.Roomid,
					Title:            recInfo.Title,
					Uname:            recInfo.Uname,
					Online:           recInfo.Online,
					Cover:            recInfo.Cover,
					Link:             "/" + strconv.Itoa(int(recInfo.Roomid)),
					Face:             recInfo.Face,
					AreaV2ParentId:   recInfo.AreaV2ParentId,
					AreaV2ParentName: recInfo.AreaV2ParentName,
					AreaV2Id:         recInfo.AreaV2Id,
					AreaV2Name:       recInfo.AreaV2Name,
					BroadcastType:    recInfo.BroadcastType,
					PendentRu:        recInfo.PendentRu,
					PendentRuPic:     recInfo.PendentRuPic,
					PendentRuColor:   recInfo.PendentRuColor,
					RecType:          _recTypeForce,
					CurrentQuality:   recInfo.CurrentQuality,
					AcceptQuality:    recInfo.AcceptQuality,
				})
				continue
			}
		}

		if len(respSlice) <= 0 {
			continue
		}
		tmpItem := respSlice[0:1][0]
		respSlice = respSlice[1:]
		pendantValue := ""
		pendantBgPic := ""
		pendantBgColor := ""
		if pendantRoomListResp != nil && pendantRoomListResp.Data != nil {
			pendantResult = pendantRoomListResp.Data.Result
			if pendantResult[tmpItem.Roomid] != nil {
				// 移动端取value, web取name
				pendantValue = pendantResult[tmpItem.Roomid].Value
				pendantBgPic = pendantResult[tmpItem.Roomid].BgPic
				if pendantResult[tmpItem.Roomid].BgColor != "" {
					pendantBgColor = pendantResult[tmpItem.Roomid].BgColor
				} else {
					pendantBgColor = _mobileIndexBadgeColorDefault
				}

			}
		}

		cover := ""
		if tmpItem.UserCover != "" {
			cover = tmpItem.UserCover
		} else {
			cover = tmpItem.Cover
		}

		resp = append(resp, &roomItem{
			RoomId:           tmpItem.Roomid,
			Title:            tmpItem.Title,
			Uname:            tmpItem.Uname,
			Online:           tmpItem.Online,
			Cover:            cover,
			Link:             "/" + strconv.Itoa(int(tmpItem.Roomid)),
			Face:             tmpItem.Face,
			AreaV2ParentId:   tmpItem.AreaV2ParentId,
			AreaV2ParentName: tmpItem.AreaV2ParentName,
			AreaV2Id:         tmpItem.AreaV2Id,
			AreaV2Name:       tmpItem.AreaV2Name,
			BroadcastType:    tmpItem.BroadcastType,
			PendentRu:        pendantValue,
			PendentRuPic:     pendantBgPic,
			PendentRuColor:   pendantBgColor,
			RecType:          _recTypeSkyHorse,
		})
	}

	return
}

func (s *Service) isSkyHorseRec(mid int64) bool {
	lastMid := strconv.Itoa(int(mid % 100))
	if len(lastMid) < 2 {
		lastMid = "0" + lastMid
	}
	_, isSkyHorseGray := s.conf.SkyHorseGray[lastMid]

	return isSkyHorseGray
}

func (s *Service) getAllList(ctx context.Context) (json.RawMessage, error) {
	allListOut, err := dao.RoomRawApi.CallRaw(ctx, 2,
		"AppIndex.getAllRawList", &liverpc.Args{})
	if err != nil {
		log.Error("[getAllList]get all list rpc error, room.v2.AppIndex.getAllRawList, error:%+v", err)
		err = errors.WithMessage(ecode.GetAllListRPCError, "GET ALL LIST FAIL#1")
		return json.RawMessage{}, err
	}
	if allListOut == nil {
		log.Error("[getAllList]get all list raw data nil, room.v2.AppIndex.getAllRawList")
		err = errors.WithMessage(ecode.GetAllListRPCError, "GET ALL LIST FAIL#2")
		return json.RawMessage{}, err
	}
	if allListOut.Code != 0 || allListOut.Data == nil {
		log.Error("[getAllList]get all list return data error, code, %d, msg: %s, error:%+v", allListOut.Code, allListOut.Message, err)
		err = errors.WithMessage(ecode.GetAllListReturnError, "GET ALL LIST FAIL#3")
		return json.RawMessage{}, err
	}
	allListJSONObj, jsonErr := simplejson.NewJson(allListOut.Data)
	if jsonErr != nil {
		log.Error("[getAllList]get all list simplejson error, error:%+v", err)
		err = errors.WithMessage(ecode.GetAllListReturnError, "GET ALL LIST FAIL#4")
		return json.RawMessage{}, err
	}

	allListCache := allListJSONObj.Get("module_list").MustArray()
	if len(allListCache) > 1 && allListCache[1] == nil {
		log.Error("[getAllList]abnormal module, allListCache:%+v", allListCache)
		err = errors.WithMessage(ecode.GetAllListReturnError, "GET ALL LIST FAIL#5")
		return json.RawMessage{}, err
	}

	return allListOut.Data, nil
}

func (s *Service) tickCacheAllList(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			allListData, err := s.getAllList(ctx)
			if err != nil {
				log.Error("[tickCacheAllList] setAllList error(%+v)", err)
				continue
			}
			if len(allListData) <= 0 {
				log.Error("[tickCacheAllList] setAllList empty data(%+v)", allListData)
				continue
			}
			log.Info("[tickCacheAllList] setAllList success!")
			s.allListInfo.Store(allListData)
		}
	}
}

func (s *Service) getAllListFromCache(ctx context.Context) (allListCache []interface{}) {
	allListRawCache := s.allListInfo.Load()

	if allListRawCache == nil {
		log.Warn("[getAllListFromCache] cache miss!")
		allList, err := s.getAllList(ctx)
		if err != nil {
			log.Error("[getAllListFromCache] pass through error(%+v)", err)
		}
		allListRawCache = allList
	}

	allListJSONObj, err := simplejson.NewJson(allListRawCache.(json.RawMessage))
	if err != nil {
		log.Error("[getAllListFromCache]get all list simplejson error, error:%+v", err)
		return
	}

	allListCache = allListJSONObj.Get("module_list").MustArray()
	if allListCache[1] == nil {
		fmt.Printf("abnormal module, allListRawCache: %+v, module_list: %+v", allListRawCache, allListJSONObj.Get("module_list"))
	}
	log.Info("[getAllListFromCache] cache hit! len: %d", len(allListCache))
	return
}
