package v0

import (
	"context"
	"encoding/json"
	"time"

	"go-common/library/cache/redis"

	"go-common/app/service/live/dao-anchor/api/grpc/v0"
	"go-common/app/service/live/dao-anchor/conf"

	"go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/library/log"
)

// CreateDataService struct
type CreateDataService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewCreateDataService init
func NewCreateDataService(c *conf.Config) (s *CreateDataService) {
	s = &CreateDataService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

//job定制接口,内部调用
//需要定时入库的数据常量定义
var contentMap = map[string]int64{dao.DANMU_NUM: 1, dao.GIFT_GOLD_AMOUNT: 1}

//GetContentMap 需要定时入库的数据常量定义，与job
//@todo 待优化逻辑
func (s *CreateDataService) GetContentMap(ctx context.Context, req *v0.GetContentMapReq) (resp *v0.GetContentMapResp, err error) {
	resp = &v0.GetContentMapResp{}
	resp.List = contentMap
	return
}

//CreateCacheList 生成历史记录list
func (s *CreateDataService) CreateCacheList(ctx context.Context, req *v0.CreateCacheListReq) (resp *v0.CreateCacheListResp, err error) {
	resp = &v0.CreateCacheListResp{}
	content := req.Content
	log.Info("createCacheList_start")
	if contentMap[content] <= 0 {
		log.Error("createCacheList_check_content_error:content=%s", content)
		return
	}
	log.Info("createCacheList_end")
	return
}

//CreateLiveCacheList 生成开播历史记录list
func (s *CreateDataService) CreateLiveCacheList(ctx context.Context, req *v0.CreateLiveCacheListReq) (resp *v0.CreateLiveCacheListResp, err error) {
	resp = &v0.CreateLiveCacheListResp{}
	content := req.Content
	roomIds := req.RoomIds
	log.Info("createLiveCacheList_start")
	if contentMap[content] <= 0 {
		log.Error("createLiveCacheList_check_content_error:content=%s", content)
		return
	}

	roomReq := &v1.RoomByIDsReq{
		RoomIds: roomIds,
		Fields:  []string{"live_status"},
	}

	reply, err := s.dao.FetchRoomByIDs(ctx, roomReq)
	if err != nil || len(reply.RoomDataSet) <= 0 {
		log.Error("createLiveCacheList_get_room_error:reply=%v;err=%v;", reply, err)
		return
	}

	// Order of data in contentList is determined by `realRoomIds`.
	contentList, err := s.dao.GetRoomRecordsCurrent(ctx, content, roomIds)
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("createLiveCacheList_get_room_content_records_error:err=%v", err)
		return
	}

	toSetValues := make(map[int64]string)
	// dao.redisKeyResp is private in dao package.
	toSetKeys := make(map[int64]interface{})
	toSetRoomIds := make([]int64, 0)
	for i, roomId := range roomIds {
		if reply.RoomDataSet[roomId] == nil {
			log.Error("createLiveCacheList_find_room_error:%d", roomId)
			continue
		}
		if reply.RoomDataSet[roomId].LiveStatus != dao.LIVE_OPEN {
			log.Info("createLiveCacheList_room_is_close;room_id=%d", roomId)
			continue
		}
		setValue := make(map[string]interface{})
		setValue["time"] = time.Now().Unix()
		setValue["value"] = contentList[i]

		valueBytes, err := json.Marshal(setValue)
		if err != nil {
			panic(err)
		}
		toSetValues[roomId] = string(valueBytes)

		liveTime := reply.RoomDataSet[roomId].LiveStartTime
		setKey := s.dao.LRoomLiveRecordList(content, roomId, liveTime)
		toSetKeys[roomId] = setKey
		toSetRoomIds = append(toSetRoomIds, roomId)
	}
	err = s.dao.SetRoomRecordsList(ctx, toSetRoomIds, toSetKeys, toSetValues)
	s.dao.DelRoomRecordsCurrent(ctx, content, toSetRoomIds)
	if err != nil {
		log.Error("createLiveCacheList_set_room_records_list:err=%v", err)
		return
	}

	log.Info("createLiveCacheList_end")
	return
}

func (s *CreateDataService) CreateDBData(ctx context.Context, req *v0.CreateDBDataReq) (resp *v0.CreateDBDataResp, err error) {
	resp = &v0.CreateDBDataResp{}
	content := req.Content
	roomIds := req.RoomIds
	log.Info("CreateDB_start")
	if contentMap[content] <= 0 {
		log.Error("CreateDB_check_content_error:content=%s", content)
		return
	}
	attrId := int64(0)
	switch content {
	case dao.DAU:
		//@todo
		break
	case dao.GIFT_NUM:
		//@todo
		break
	case dao.GIFT_GOLD_NUM:
		//@todo
		break
	case dao.GIFT_GOLD_AMOUNT:
		attrId = dao.ATTRID_REVENUE
		break
	case dao.DANMU_NUM:
		attrId = dao.ATTRID_DANMU
		break
	default:
		break
	}

	recordsRange, err := s.dao.GetRoomLiveRecordsRange(ctx, content, roomIds, 0, 0, 59)
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("CreateDBData_redis_GetRoomLiveRecordsRange failed completely error(%v)", err)
		return nil, err
	}

	if attrId < 0 {
		return
	}
	for roomId, valueInfos := range recordsRange {
		attrReq := &v1.RoomAttrCreateReq{
			AttrId: attrId,
			RoomId: roomId,
		}

		values := getAttrData(attrId, valueInfos)
		for attrSubId, attrValue := range values {
			attrReq.AttrSubId = attrSubId
			attrReq.AttrValue = attrValue
			s.dao.RoomAttrCreate(ctx, attrReq)
		}
	}

	return
}
func (s *CreateDataService) CreateRoomExtendDBData(ctx context.Context, content string, roomIds []int64) (err error) {
	log.Info("CreateDB_start")
	if contentMap[content] <= 0 {
		log.Error("CreateDB_check_content_error:content=%s", content)
		return
	}
	RoomExtendReq, err := s.GetRoomExtendDataReq(content)
	if err != nil {
		log.Error("CreateDB_RoomExtendReq_error:content=%s;roomIds=%v", content, roomIds)
	}
	if err == nil && RoomExtendReq != nil {
		for _, roomId := range roomIds {
			RoomExtendReq.RoomId = roomId
		}
		s.dao.RoomExtendUpdate(ctx, RoomExtendReq)
	}
	return
}

func (s *CreateDataService) GetRoomExtendDataReq(content string) (resp *v1.RoomExtendUpdateReq, err error) {

	resp = &v1.RoomExtendUpdateReq{}
	switch content {
	case dao.DAU:
		resp.AudienceCount = 0
		break
	case dao.GIFT_NUM:
		resp.GiftCount = 0
		break
	case dao.GIFT_GOLD_NUM:
		resp.GiftGoldCount = 0
		break
	case dao.GIFT_GOLD_AMOUNT:
		resp.GiftGoldAmount = 0
		break
	default:
		break
	}
	return
}
func getAttrData(AttrId int64, listInfos []*dao.ListIntValueInfo) (value map[int64]int64) {

	value = make(map[int64]int64)
	switch AttrId {
	case dao.ATTRID_DANMU:
		for i, listInfo := range listInfos {
			if i < 15 {
				value[dao.ATTRSUBID_DANMU_MINUTE_NUM_15] += listInfo.Value
			}
			if i < 30 {
				value[dao.ATTRSUBID_DANMU_MINUTE_NUM_30] += listInfo.Value
			}
			if i < 45 {
				value[dao.ATTRSUBID_DANMU_MINUTE_NUM_45] += listInfo.Value
			}
			if i < 60 {
				value[dao.ATTRSUBID_DANMU_MINUTE_NUM_60] += listInfo.Value
			}
		}
		break
	case dao.ATTRID_REVENUE:
		for i, listInfo := range listInfos {
			if i < 15 {
				value[dao.ATTRSUBID_REVENUE_MINUTE_NUM_15] += listInfo.Value
			}
			if i < 30 {
				value[dao.ATTRSUBID_REVENUE_MINUTE_NUM_30] += listInfo.Value
			}
			if i < 45 {
				value[dao.ATTRSUBID_REVENUE_MINUTE_NUM_45] += listInfo.Value
			}
			if i < 60 {
				value[dao.ATTRSUBID_REVENUE_MINUTE_NUM_60] += listInfo.Value
			}
		}
		break
	default:
		break
	}
	return
}
