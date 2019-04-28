package dao

import (
	"context"
	"fmt"

	"go-common/app/common/live/library/lrucache"

	jsonitor "github.com/json-iterator/go"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
)

const (
	INFO_ROOM = 1 << iota
	INFO_ROOM_EXT
	INFO_TAG
	INFO_ANCHOR
	INFO_SHORT_ID
	INFO_AREA_INFO
)

const (
	INFO_ALL = (INFO_ROOM | INFO_ROOM_EXT | INFO_TAG | INFO_ANCHOR | INFO_SHORT_ID | INFO_AREA_INFO)
)

const (
	FETCH_PAGE_SIZE = 100
)

//消费类型常量 定义
const (
	//弹幕
	//DANMU_NUM 当前弹幕累计数量
	DANMU_NUM = "danmu_num"
	//DANMU_MINUTE_NUM_15 最近15分钟弹幕数量
	DANMU_MINUTE_NUM_15 = "danmu_minute_num_15"
	//DANMU_MINUTE_NUM_30 ...
	DANMU_MINUTE_NUM_30 = "danmu_minute_num_30"
	//DANMU_MINUTE_NUM_45 ...
	DANMU_MINUTE_NUM_45 = "danmu_minute_num_45"
	//DANMU_MINUTE_NUM_60 ...
	DANMU_MINUTE_NUM_60 = "danmu_minute_num_60"
	//人气
	//POPULARITY 当前实时人气
	POPULARITY = "popularity"
	//POPULARITY_MAX_TO_ARG_7 7日峰值人气的均值
	POPULARITY_MAX_TO_ARG_7 = "popularity_max_to_avg_7"
	//POPULARITY_MAX_TO_ARG_30 30日人气峰值的均值
	POPULARITY_MAX_TO_ARG_30 = "popularity_max_to_avg_30"
	//送礼
	//GIFT_NUM 实时送礼数
	GIFT_NUM = "gift_num_current_total"
	//GIFT_GOLD_AMOUNT 实时消费金瓜子数
	GIFT_GOLD_NUM = "gift_gold_num"
	//GIFT_GOLD_AMOUNT 实时消费金瓜子金额
	GIFT_GOLD_AMOUNT = "gift_gold_amount"
	//GIFT_GOLD_AMOUNT_MINUTE_15 最近15分钟金瓜子金额
	GIFT_GOLD_AMOUNT_MINUTE_15 = "gift_gold_num_minute_15"
	//GIFT_GOLD_AMOUNT_MINUTE_30 最近30分钟金瓜子金额
	GIFT_GOLD_AMOUNT_MINUTE_30 = "gift_gold_num_minute_30"
	//GIFT_GOLD_AMOUNT_MINUTE_45 ...
	GIFT_GOLD_AMOUNT_MINUTE_45 = "gift_gold_num_minute_45"
	//GIFT_GOLD_AMOUNT_MINUTE_60 ...
	GIFT_GOLD_AMOUNT_MINUTE_60 = "gift_gold_num_minute_60"
	//有效开播天数
	//VALID_LIVE_DAYS_TYPE_1_DAY_7 7日内有效开播天数；有效开播:一次开播大于5分钟
	VALID_LIVE_DAYS_TYPE_1_DAY_7 = "valid_days_type_1_day_7"
	//VALID_LIVE_DAYS_TYPE_1_DAY_14 14日内有效开播天数；有效开播:一次开播大于5分钟
	VALID_LIVE_DAYS_TYPE_1_DAY_14 = "valid_days_type_1_day_14"
	//VALID_LIVE_DAYS_TYPE_2_DAY_7 7日内有效开播天数；有效开播:大于等于120分钟
	VALID_LIVE_DAYS_TYPE_2_DAY_7 = "valid_days_type_2_day_7"
	//VALID_LIVE_DAYS_TYPE_2_DAY_30 14日内有效开播天数；有效开播:大于等于120分钟
	VALID_LIVE_DAYS_TYPE_2_DAY_30 = "valid_days_type_2_day_30"
	//房间状态
	//ROOM_TAG_CURRENT 房间实时标签
	ROOM_TAG_CURRENT = "room_tag_current"
	//榜单
	//RANK_LIST_CURRENT 排行榜相关数据
	RANK_LIST_CURRENT = "rank_list_current"
	//DAU
	DAU = "dau"
)

const (
	_RoomIdMappingCacheCapacity = 1024
)

// Dao dao
type Dao struct {
	c               *conf.Config
	redis           *redis.Pool
	db              *xsql.DB
	dbLiveApp       *xsql.DB
	shortIDMapping  *lrucache.SyncCache
	areaInfoMapping *lrucache.SyncCache
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:               c,
		redis:           redis.NewPool(c.Redis),
		db:              xsql.NewMySQL(c.MySQL),
		dbLiveApp:       xsql.NewMySQL(c.LiveAppMySQL),
		shortIDMapping:  lrucache.NewSyncCache(c.LRUCache.Bucket, c.LRUCache.Capacity, c.LRUCache.Timeout),
		areaInfoMapping: lrucache.NewSyncCache(c.LRUCache.Bucket, c.LRUCache.Capacity, c.LRUCache.Timeout),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}

// FetchRoomByIDs implementation
// FetchRoomByIDs 查询房间信息
func (d *Dao) FetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	if len(req.RoomIds) > 0 {
		req.RoomIds, err = d.dbNormalizeRoomIDs(ctx, req.RoomIds)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] normalize ids error(%v), req(%v)", err, req)
			return nil, err
		}
	}

	// TODO： 处理部分fields的情况，需要考虑特殊status的依赖问题
	if len(req.RoomIds) > 0 {
		resp = &v1pb.RoomByIDsResp{
			RoomDataSet: make(map[int64]*v1pb.RoomData),
		}

		idsDB := make([]int64, 0, len(req.RoomIds))
		// 从redis获取房间所有信息
		for _, id := range req.RoomIds {
			data, err := d.redisGetRoomInfo(ctx, id, _allRoomInfoFields)
			if err != nil {
				idsDB = append(idsDB, id)
			} else {
				d.dbDealWithStatus(ctx, data)
				resp.RoomDataSet[id] = data
			}
		}

		// 需要回源DB取数据
		if len(idsDB) > 0 {
			// 分段处理
			for start := 0; start < len(idsDB); start += FETCH_PAGE_SIZE {
				end := start + FETCH_PAGE_SIZE
				if end > len(idsDB) {
					end = len(idsDB)
				}

				reqRoom := &v1pb.RoomByIDsReq{
					RoomIds: idsDB[start:end],
					Fields:  _allRoomInfoFields,
				}
				respRoom, err := d.dbFetchRoomByIDs(ctx, reqRoom)
				if err != nil {
					log.Error("[RoomOnlineList] dbFetchRoomByIDs error(%v), reqRoom(%v)", err, reqRoom)
					return nil, err
				}

				// 回写房间信息到redis
				for _, id := range idsDB[start:end] {
					resp.RoomDataSet[id] = respRoom.RoomDataSet[id]
					d.redisSetRoomInfo(ctx, id, _allRoomInfoFields, respRoom.RoomDataSet[id], false)
				}
			}
		}
	} else if len(req.Uids) > 0 {
		// TODO 根据主播ID查询房间号的场景较少，暂不优化，后续先转房间号
		resp, err = d.dbFetchRoomByIDs(ctx, req)
	}
	return
}

// RoomOnlineList implementation
// RoomOnlineList 在线房间列表
func (d *Dao) RoomOnlineList(ctx context.Context, req *v1pb.RoomOnlineListReq) (resp *v1pb.RoomOnlineListResp, err error) {
	log.Info("[dao|RoomOnlineList] req(%v)", err, req)
	ids, err := d.redisGetOnlineList(ctx, _onlineListAllArea)
	if err != nil || len(ids) <= 0 {
		ids, err = d.dbOnlineListByArea(ctx, _onlineListAllArea)
		if err != nil {
			log.Error("[RoomOnlineListByAttrs] dbOnlineListByArea error(%v), req(%v)", err, req)
			return nil, err
		}
		d.redisSetOnlineList(ctx, _onlineListAllArea, ids)
	}

	resp = &v1pb.RoomOnlineListResp{
		RoomDataList: make(map[int64]*v1pb.RoomData),
	}

	// 分页逻辑
	start := int(req.Page * req.PageSize)
	size := len(ids)
	if start >= size {
		return
	}

	end := start + int(req.PageSize)
	if end > size {
		end = size
	}
	ids = ids[start:end]

	idsDB := make([]int64, 0, len(ids))
	// 从redis获取房间信息
	for _, id := range ids {
		data, err := d.redisGetRoomInfo(ctx, id, req.Fields)
		if err != nil {
			idsDB = append(idsDB, id)
		} else {
			d.dbDealWithStatus(ctx, data)
			resp.RoomDataList[id] = data
		}
	}

	// 需要回源DB取数据
	if len(idsDB) > 0 {
		reqRoom := &v1pb.RoomByIDsReq{
			RoomIds: idsDB,
			Fields:  _allRoomInfoFields,
		}
		respRoom, err := d.dbFetchRoomByIDs(ctx, reqRoom)
		if err != nil {
			log.Error("[RoomOnlineList] dbFetchRoomByIDs error(%v), reqRoom(%v)", err, reqRoom)
			return nil, err
		}

		// 回写房间信息到redis
		for _, id := range idsDB {
			resp.RoomDataList[id] = respRoom.RoomDataSet[id]
			d.redisSetRoomInfo(ctx, id, _allRoomInfoFields, respRoom.RoomDataSet[id], false)
		}
	}
	return
}

// RoomOnlineListByArea implementation
// RoomOnlineListByArea 分区在线房间列表
func (d *Dao) RoomOnlineListByArea(ctx context.Context, req *v1pb.RoomOnlineListByAreaReq) (resp *v1pb.RoomOnlineListByAreaResp, err error) {
	idSet := make(map[int64]bool)

	idsDB := make([]int64, 0)
	if len(req.AreaIds) <= 0 {
		req.AreaIds = []int64{0}
	}

	for _, areaID := range req.AreaIds {
		ids, err := d.redisGetOnlineList(ctx, areaID)
		if err != nil {
			idsDB = append(idsDB, areaID)
		} else {
			for _, id := range ids {
				idSet[id] = true
			}
		}
	}

	// 需要回源DB取数据
	if len(idsDB) > 0 {
		for _, areaID := range idsDB {
			roomIds, err := d.dbOnlineListByArea(ctx, areaID)
			if err != nil {
				log.Error("[RoomOnlineListByArea] dbOnlineListByArea error(%v), areaID(%v)", err, areaID)
				return nil, err
			}
			d.redisSetOnlineList(ctx, areaID, roomIds)
			for _, id := range roomIds {
				idSet[id] = true
			}
		}
	}

	resp = &v1pb.RoomOnlineListByAreaResp{
		RoomIds: make([]int64, 0, len(idSet)),
	}

	for id := range idSet {
		resp.RoomIds = append(resp.RoomIds, id)
	}
	return
}

var (
	_fields = []string{"uid", "area_id", "parent_area_id", "popularity_count", "anchor_profile_type"}
)

// RoomOnlineListByAttrs implementation
// RoomOnlineListByAttrs 在线房间维度信息(不传attrs，不查询attr)
func (d *Dao) RoomOnlineListByAttrs(ctx context.Context, req *v1pb.RoomOnlineListByAttrsReq) (resp *v1pb.RoomOnlineListByAttrsResp, err error) {
	ids, err := d.redisGetOnlineList(ctx, _onlineListAllArea)
	if err != nil || len(ids) <= 0 {
		ids, err = d.dbOnlineListByArea(ctx, _onlineListAllArea)
		if err != nil {
			log.Error("[RoomOnlineListByAttrs] dbOnlineListByArea error(%v), req(%v)", err, req)
			return nil, err
		}
		d.redisSetOnlineList(ctx, _onlineListAllArea, ids)
	}

	resp = &v1pb.RoomOnlineListByAttrsResp{
		Attrs: make(map[int64]*v1pb.AttrResp),
	}

	idsDB := make([]int64, 0, len(ids))
	for _, id := range ids {
		// 从redis获取房间基础信息
		data, err := d.redisGetRoomInfo(ctx, id, _fields)
		if err != nil {
			idsDB = append(idsDB, id)
		} else {
			resp.Attrs[id] = &v1pb.AttrResp{
				Uid:               data.Uid,
				RoomId:            id,
				AreaId:            data.AreaId,
				ParentAreaId:      data.ParentAreaId,
				PopularityCount:   data.PopularityCount,
				AnchorProfileType: data.AnchorProfileType,
			}
		}
	}

	// 需要回源DB取数据
	if len(idsDB) > 0 {
		eg := errgroup.Group{}
		// 分段处理
		for start := 0; start < len(idsDB); start += FETCH_PAGE_SIZE {
			end := start + FETCH_PAGE_SIZE
			if end > len(idsDB) {
				end = len(idsDB)
			}
			eg.Go(func(idsDB []int64, start, end int) func() error {
				return func() (err error) {
					reqRoom := &v1pb.RoomByIDsReq{
						RoomIds: idsDB[start:end],
						Fields:  _allRoomInfoFields,
					}
					respRoom, err := d.dbFetchRoomByIDs(ctx, reqRoom)
					if err != nil {
						log.Error("[RoomOnlineList] dbFetchRoomByIDs error(%v), reqRoom(%v)", err, reqRoom)
						return err
					}

					// 回写房间信息到redis
					for _, id := range idsDB[start:end] {
						resp.Attrs[id] = &v1pb.AttrResp{
							Uid:               respRoom.RoomDataSet[id].Uid,
							RoomId:            id,
							AreaId:            respRoom.RoomDataSet[id].AreaId,
							ParentAreaId:      respRoom.RoomDataSet[id].ParentAreaId,
							PopularityCount:   respRoom.RoomDataSet[id].PopularityCount,
							AnchorProfileType: respRoom.RoomDataSet[id].AnchorProfileType,
							TagList:           respRoom.RoomDataSet[id].TagList,
							AttrList:          make([]*v1pb.AttrData, 0, len(req.Attrs)),
						}
						d.redisSetRoomInfo(ctx, id, _allRoomInfoFields, respRoom.RoomDataSet[id], false)
						d.redisSetTagList(ctx, id, respRoom.RoomDataSet[id].TagList)
					}
					return
				}
			}(idsDB, start, end))
		}
		eg.Wait()
	}

	// 重置回源数组
	idsDB = make([]int64, 0, len(ids))
	for _, id := range ids {
		if resp.Attrs[id].TagList == nil {
			// 从redis获取房间Tag信息
			data, err := d.redisGetTagList(ctx, id)
			if err != nil {
				idsDB = append(idsDB, id)
			} else {
				resp.Attrs[id].TagList = data
			}
		}
	}

	// 需要回源DB取数据
	if len(idsDB) > 0 {
		eg := errgroup.Group{}
		// 分段处理
		for start := 0; start < len(idsDB); start += FETCH_PAGE_SIZE {
			end := start + FETCH_PAGE_SIZE
			if end > len(idsDB) {
				end = len(idsDB)
			}
			eg.Go(func(idsDB []int64, start, end int) func() error {
				return func() (err error) {
					respRoom := make(map[int64]*v1pb.RoomData)
					err = d.dbFetchTagInfo(ctx, idsDB[start:end], respRoom)
					if err != nil {
						log.Error("[RoomOnlineList] dbFetchTagInfo error(%v), idsDB[start:end](%v)", err, idsDB[start:end])
						return err
					}

					// 回写房间Tag信息到redis
					for _, id := range idsDB[start:end] {
						if tag, ok := respRoom[id]; ok {
							resp.Attrs[id].TagList = tag.TagList
						} else {
							resp.Attrs[id].TagList = make([]*v1pb.TagData, 0, len(req.Attrs))
						}
						d.redisSetTagList(ctx, id, resp.Attrs[id].TagList)
					}
					return
				}
			}(idsDB, start, end))
		}
		eg.Wait()
	}

	// TODO 从redis获取attr列表
	// TODO 批量从db获取attr列表
	if len(req.Attrs) > 0 {
		eg := errgroup.Group{}
		for _, attr := range req.Attrs {
			// 实时人气值特殊处理
			if attr.AttrId == ATTRID_POPULARITY && attr.AttrSubId == ATTRSUBID_POPULARITY_REALTIME {
				for _, attrResp := range resp.Attrs {
					resp.Attrs[attrResp.RoomId].AttrList = append(resp.Attrs[attrResp.RoomId].AttrList, &v1pb.AttrData{
						RoomId:    attrResp.RoomId,
						AttrId:    attr.AttrId,
						AttrSubId: attr.AttrSubId,
						AttrValue: attrResp.PopularityCount,
					})
				}
				continue
			}
			eg.Go(func(attr *v1pb.AttrReq) func() error {
				return func() (err error) {
					reqAttr := &v1pb.FetchAttrByIDsReq{
						AttrId:    attr.AttrId,
						AttrSubId: attr.AttrSubId,
						RoomIds:   ids,
					}
					respAttr, err := d.FetchAttrByIDs(ctx, reqAttr)
					if err != nil {
						log.Error("[RoomOnlineListByAttrs] FetchAttrByIDs from db error(%v), reqAttr(%v)", err, reqAttr)
						return err
					}
					for _, attr := range respAttr.Attrs {
						resp.Attrs[attr.RoomId].AttrList = append(resp.Attrs[attr.RoomId].AttrList, attr)
					}
					return
				}
			}(attr))
		}
		eg.Wait()
	}

	return
}

// RoomCreate implementation
// RoomCreate 房间创建
func (d *Dao) RoomCreate(ctx context.Context, req *v1pb.RoomCreateReq) (resp *v1pb.RoomCreateResp, err error) {
	return d.roomCreate(ctx, req)
}

// RoomUpdate implementation
// RoomUpdate 房间更新
func (d *Dao) RoomUpdate(ctx context.Context, req *v1pb.RoomUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.roomUpdate(ctx, req)
	if err == nil {
		fields := make([]string, 0, len(req.Fields))
		data := &v1pb.RoomData{
			RoomId:      req.RoomId,
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		for _, f := range req.Fields {
			switch f {
			case "title":
				data.Title = req.Title
			case "cover":
				data.Cover = req.Cover
			case "tags":
				data.Tags = req.Tags
			case "background":
				data.Background = req.Background
			case "description":
				data.Description = req.Description
			case "live_start_time":
				data.LiveStartTime = req.LiveStartTime
				// 更新在播列表
				var areaID int64
				reqRoom := &v1pb.RoomByIDsReq{
					RoomIds: []int64{req.RoomId},
					Fields:  _allRoomInfoFields,
				}
				respRoom, err := d.FetchRoomByIDs(ctx, reqRoom)
				if err != nil {
					log.Error("[RoomOnlineList] dbFetchRoomByIDs error(%v), reqRoom(%v)", err, reqRoom)
				} else {
					if respRoom.RoomDataSet[req.RoomId] != nil {
						areaID = respRoom.RoomDataSet[req.RoomId].AreaId
					}
				}
				if req.LiveStartTime > 0 {
					d.redisAddOnlineList(ctx, _onlineListAllArea, req.RoomId)
					d.redisAddOnlineList(ctx, areaID, req.RoomId)
				} else {
					d.redisDelOnlineList(ctx, _onlineListAllArea, req.RoomId)
					d.redisDelOnlineList(ctx, areaID, req.RoomId)
				}
				// TODO 更新开播状态
			case "live_screen_type":
				data.LiveScreenType = req.LiveScreenType
			case "live_type":
				data.LiveType = req.LiveType
			case "lock_status":
				data.LockStatus = req.LockStatus
			case "lock_time":
				data.LockTime = req.LockTime
			case "hidden_time":
				data.HiddenTime = req.HiddenTime
				// TODO 更新隐藏状态
			case "area_id":
				data.AreaId = req.AreaId
				if req.AreaId > 0 {
					areaInfo, err := d.dbFetchAreaInfo(ctx, req.AreaId)
					if err != nil {
						log.Error("[dao.dao-anchor.mysql|roomUpdate] fetch area info error(%v), req(%v)", err, req)
						err = ecode.InvalidParam
						return nil, err
					}
					data.ParentAreaId = areaInfo.ParentAreaID
					fields = append(fields, "parent_area_id")
				}
			default:
				continue
			}
			fields = append(fields, f)
		}
		d.redisSetRoomInfo(ctx, data.RoomId, fields, data, true)
	}
	return
}

// RoomBatchUpdate implementation
// RoomBatchUpdate 房间更新
func (d *Dao) RoomBatchUpdate(ctx context.Context, req *v1pb.RoomBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	for _, r := range req.Reqs {
		res, err := d.RoomUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|RoomBatchUpdate] update room record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// RoomExtendUpdate implementation
// RoomExtendUpdate 房间更新
func (d *Dao) RoomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.roomExtendUpdate(ctx, req)
	if err == nil {
		fields := make([]string, 0, len(req.Fields))
		data := &v1pb.RoomData{
			RoomId:      req.RoomId,
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		for _, f := range req.Fields {
			switch f {
			case "keyframe":
				data.Keyframe = req.Keyframe
			case "popularity_count":
				data.PopularityCount = req.PopularityCount
			default:
				continue
			}
			fields = append(fields, f)
		}
		d.redisSetRoomInfo(ctx, data.RoomId, fields, data, true)
	}
	return
}

// RoomExtendBatchUpdate implementation
// RoomExtendBatchUpdate 房间更新
func (d *Dao) RoomExtendBatchUpdate(ctx context.Context, req *v1pb.RoomExtendBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	for _, r := range req.Reqs {
		res, err := d.RoomExtendUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|RoomExtendBatchUpdate] update room extend record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// RoomExtendIncre implementation
// RoomExtendIncre 房间增量更新
func (d *Dao) RoomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.roomExtendIncre(ctx, req)
	if err == nil {
		fields := make([]string, 0, len(req.Fields))
		data := &v1pb.RoomData{
			RoomId:      req.RoomId,
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		for _, f := range req.Fields {
			switch f {
			case "popularity_count":
				data.PopularityCount = req.PopularityCount
			default:
				continue
			}
			fields = append(fields, f)
		}
		if len(fields) > 0 {
			d.redisIncreRoomInfo(ctx, data.RoomId, fields, data)
		}
	}
	return
}

// RoomExtendBatchIncre implementation
// RoomExtendBatchIncre 房间增量更新
func (d *Dao) RoomExtendBatchIncre(ctx context.Context, req *v1pb.RoomExtendBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	for _, r := range req.Reqs {
		res, err := d.RoomExtendIncre(ctx, r)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|RoomExtendBatchIncre] update room extend increment record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// RoomTagCreate implementation
// RoomTagCreate 房间Tag创建
func (d *Dao) RoomTagCreate(ctx context.Context, req *v1pb.RoomTagCreateReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.roomTagCreate(ctx, req)
	if err == nil {
		tag := &v1pb.TagData{
			TagId:       req.TagId,
			TagSubId:    req.TagSubId,
			TagValue:    req.TagValue,
			TagExt:      req.TagExt,
			TagExpireAt: req.TagExpireAt,
		}
		d.redisAddTag(ctx, req.RoomId, tag)
	}
	return
}

// RoomAttrCreate implementation
// RoomAttrCreate 房间Attr创建
func (d *Dao) RoomAttrCreate(ctx context.Context, req *v1pb.RoomAttrCreateReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomAttrCreate(ctx, req)
}

// RoomAttrSetEx implementation
// RoomAttrSetEx 房间Attr更新
func (d *Dao) RoomAttrSetEx(ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomAttrSetEx(ctx, req)
}

// AnchorUpdate implementation
// AnchorUpdate 主播更新
func (d *Dao) AnchorUpdate(ctx context.Context, req *v1pb.AnchorUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.anchorUpdate(ctx, req)
	if err == nil {
		roomID := d.dbFetchRoomIDByUID(ctx, req.Uid)
		if roomID == 0 {
			return
		}

		fields := make([]string, 0, len(req.Fields))
		data := &v1pb.RoomData{
			RoomId:      roomID,
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		for _, f := range req.Fields {
			switch f {
			case "profile_type":
				f = "anchor_profile_type"
				data.AnchorProfileType = req.ProfileType
			case "san_score":
				f = "anchor_san"
				data.AnchorSan = req.SanScore
			case "round_status":
				f = "anchor_round_switch"
				data.AnchorRoundSwitch = req.RoundStatus
			case "record_status":
				f = "anchor_record_switch"
				data.AnchorRecordSwitch = req.RecordStatus
			case "exp":
				f = "anchor_exp"
				data.AnchorLevel.Score = req.Exp
			default:
				log.Error("[dao.dao-anchor.mysql|anchorUpdate] unsupported field(%v), req(%s)", f, req)
				err = ecode.InvalidParam
				return
			}
			fields = append(fields, f)
		}
		d.redisSetRoomInfo(ctx, data.RoomId, fields, data, true)
	}
	return
}

// AnchorBatchUpdate implementation
// AnchorBatchUpdate 主播更新
func (d *Dao) AnchorBatchUpdate(ctx context.Context, req *v1pb.AnchorBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	for _, r := range req.Reqs {
		res, err := d.AnchorUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|AnchorBatchUpdate] update anchor record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// AnchorIncre implementation
// AnchorIncre 主播增量更新
func (d *Dao) AnchorIncre(ctx context.Context, req *v1pb.AnchorIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp, err = d.anchorIncre(ctx, req)
	if err == nil {
		roomID := d.dbFetchRoomIDByUID(ctx, req.Uid)
		if roomID == 0 {
			return
		}

		fields := make([]string, 0, len(req.Fields))
		data := &v1pb.RoomData{
			RoomId:      roomID,
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		for _, f := range req.Fields {
			switch f {
			case "san_score":
				f = "anchor_san"
				data.AnchorSan = req.SanScore
			case "exp":
				f = "anchor_exp"
				data.AnchorLevel.Score = req.Exp
			default:
				continue
			}
			fields = append(fields, f)
		}
		if len(fields) > 0 {
			d.redisIncreRoomInfo(ctx, data.RoomId, fields, data)
		}
	}
	return
}

// AnchorBatchIncre implementation
// AnchorBatchIncre 主播增量更新
func (d *Dao) AnchorBatchIncre(ctx context.Context, req *v1pb.AnchorBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	for _, r := range req.Reqs {
		res, err := d.AnchorIncre(ctx, r)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|AnchorBatchIncre] update anchor increment record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// FetchAreas implementation
// FetchAreas 根据父分区号查询子分区
func (d *Dao) FetchAreas(ctx context.Context, req *v1pb.FetchAreasReq) (resp *v1pb.FetchAreasResp, err error) {
	return d.fetchAreas(ctx, req)
}

// FetchAttrByIDs implementation
// FetchAttrByIDs 批量根据房间号查询指标
func (d *Dao) FetchAttrByIDs(ctx context.Context, req *v1pb.FetchAttrByIDsReq) (resp *v1pb.FetchAttrByIDsResp, err error) {
	return d.fetchAttrByIDs(ctx, req)
}

// DeleteAttr implementation
// DeleteAttr 删除一个指标
func (d *Dao) DeleteAttr(ctx context.Context, req *v1pb.DeleteAttrReq) (resp *v1pb.UpdateResp, err error) {
	return d.deleteAttr(ctx, req)
}

type msgVal struct {
	MsgID string `json:"msg_id"`
}

func getConsumedKey(topic string, msgID string) string {
	return fmt.Sprintf("consumed:%s:%s", topic, msgID)
}

// 清除消费过的记录，主要用于测试
func (d *Dao) clearConsumed(ctx context.Context, msg *databus.Message) {
	val := &msgVal{}
	err := jsonitor.Unmarshal(msg.Value, val)
	if err != nil {
		return
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("DEL", getConsumedKey(msg.Topic, val.MsgID))
}

// CanConsume 是否可以消费
func (d *Dao) CanConsume(ctx context.Context, msg *databus.Message) bool {
	val := &msgVal{}
	err := jsonitor.Unmarshal(msg.Value, val)
	if err != nil {
		log.Error("unmarshal msg value error %+v, value: %s", err, string(msg.Value))
		return true
	}
	if val.MsgID == "" {
		log.Warn("msg_id is empty ; value: %s", string(msg.Value))
		return true
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()

	var key = getConsumedKey(msg.Topic, val.MsgID)
	reply, err := conn.Do("SET", key, "1", "NX", "EX", 86400) // 24 hours
	if err == nil {
		if reply == nil {
			log.Info("Already consumed key:%s", key)
			return false
		} else {
			return true
		}
	}

	if err == redis.ErrNil {
		//already consumed
		log.Info("Already consumed key:%s", key)
		return false
	}
	// other redis error happenned, let it pass
	log.Error("redis error when resolve CanConsume %+v", err)
	return true
}
