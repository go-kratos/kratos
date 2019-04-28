package service

import (
	"context"
	"errors"
	"fmt"
	"go-common/app/job/live/push-search/dao"
	"go-common/app/job/live/push-search/model"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	userV3 "go-common/app/service/live/user/api/liverpc/v3"
	accountApi "go-common/app/service/main/account/api"

	"go-common/library/log"
	rpccontext "go-common/library/net/rpc/liverpc/context"
	"strconv"
	"time"
)

var (
	hbaseTable  = "live:PushSearch"
	hbaseFamily = "search"
	fields      = []string{
		"roomid",
		"short_id",
		"uid",
		"uname",
		"area",
		"title",
		"tags",
		"try_time",
		"cover",
		"user_cover",
		"lock_status",
		"hidden_status",
		"attentions",
		"online",
		"live_time",
		"area_v2_id",
		"area_v2_parent_id",
		"virtual",
		"round_status",
		"on_flag",
		"area_v2_name",
		"ctime",
		"mtime",
	}
)

func (s *Service) getBaseRoomInfo(uid int64) (roomInfo *roomV2.RoomGetByIdsResp_RoomInfo, err error) {

	roomIdResp, err := dao.RoomApi.V2Room.RoomIdByUid(rpccontext.WithTimeout(context.TODO(), 50*time.Millisecond), &roomV2.RoomRoomIdByUidReq{
		Uid: uid,
	})
	if err != nil {
		log.Error("[getBaseRoomInfo]RoomIdByUid rpc error, error:%+v", err)
		return
	}
	if roomIdResp.Code != 0 {
		log.Error("[getBaseRoomInfo]RoomIdByUid return error, code:%d, msg:%s", roomIdResp.Code, roomIdResp.Msg)
		err = errors.New("getRoomId return error")
		return
	}
	if roomIdResp.Data == nil {
		log.Error("[getBaseRoomInfo]GetMultiple empty data")
		err = errors.New("getRoomId empty error")
		return
	}

	if roomIdResp.Data.RoomId == 0 {
		log.Error("[getBaseRoomInfo]GetMultiple empty data")
		err = errors.New("roomId not found error")
		return
	}

	roomInfoResp := &roomV2.RoomGetByIdsResp{}
	roomInfoResp, err = dao.RoomApi.V2Room.GetByIds(rpccontext.WithTimeout(context.TODO(), 50*time.Millisecond), &roomV2.RoomGetByIdsReq{
		Ids:    []int64{roomIdResp.Data.RoomId},
		From:   "push-search",
		Fields: fields,
	})
	if err != nil {
		log.Error("[getBaseRoomInfo]GetByIds rpc error, error:%+v", err)
		return
	}
	if roomInfoResp.Code != 0 {
		log.Error("[getBaseRoomInfo]GetByIds return error, code:%d, msg:%s", roomInfoResp.Code, roomInfoResp.Msg)
		err = errors.New("GetByIds return error")
		return
	}
	if roomInfoResp.Data == nil {
		log.Error("[getBaseRoomInfo]GetByIds empty data")
		err = errors.New("GetByIds empty error")
		return
	}

	info, ok := roomInfoResp.Data[roomIdResp.Data.RoomId]
	if !ok {
		log.Error("[getBaseRoomInfo]GetByIds not found")
		err = errors.New("roomId not found error")
		return
	}

	roomInfo = info
	return
}

func (s *Service) getMultiUserInfo(uid int64) (userInfo *userV3.UserGetMultipleResp_Info, err error) {
	userInfo = &userV3.UserGetMultipleResp_Info{}
	pr, err := s.AccountClient.Profile3(context.TODO(), &accountApi.MidReq{Mid: uid})

	if err != nil {
		log.Error("[getMultiUserInfo]Profile3 rpc error, error:%+v", err)
		return
	}

	if pr == nil {
		log.Error("[getMultiUserInfo]Profile3 empty data")
		err = errors.New("user empty error")
		return
	}

	userInfo.Uid = uid
	userInfo.Uname = pr.GetProfile().GetName()
	userInfo.Face = pr.GetProfile().GetFace()

	if userInfo.Uname != "" {
		return
	}

	err = errors.New("user not found")
	log.Error("[getMultiUserInfo]GetMultiple no user, data:%+v", userInfo)

	return
}

func (s *Service) getFc(uid int64) (fc int, err error) {
	fcResp := &relationV1.FeedGetUserFcResp{}
	fcResp, err = dao.RelationApi.V1Feed.GetUserFc(rpccontext.WithTimeout(context.TODO(), 50*time.Millisecond), &relationV1.FeedGetUserFcReq{
		Follow: uid,
	})
	if err != nil {
		log.Error("[getFc]GetFc rpc error, error:%+v", err)
		return
	}
	if fcResp.Code != 0 {
		log.Error("[getFc]GetFc return error, code:%d, msg:%s", fcResp.Code, fcResp.Msg)
		err = errors.New("fc return error")
		return
	}

	if fcResp.Data == nil {
		log.Error("[getFc]GetFc empty data")
		err = errors.New("fc empty error")
		return
	}

	fc = int(fcResp.Data.Fc)
	return
}

func (s *Service) saveHBase(c context.Context, key string, columnInfo map[string][]byte) (err error) {
	var ctx, cancel = context.WithTimeout(c, time.Duration(s.c.SearchHBase.WriteTimeout)*time.Millisecond)

	defer cancel()

	values := map[string]map[string][]byte{hbaseFamily: columnInfo}
	if _, err = s.dao.SearchHBase.PutStr(ctx, hbaseTable, key, values); err != nil {
		log.Error("SearchHBase.PutStr error(%v), table(%s), values(%+v)", err, hbaseTable, values)
	}
	return
}

func (s *Service) getLockStatus(lockStatus string) int {
	status := 0
	if lockStatus != "0000-00-00 00:00:00" {
		status = 1
	}
	return status
}

func (s *Service) getHiddenStatus(HiddenStatus string) int {
	status := 0
	if HiddenStatus != "0000-00-00 00:00:00" {
		status = 1
	}
	return status
}

//hbase key roomID md5
func (s *Service) rowKey(roomId int) string {
	key := fmt.Sprintf("%d_%d", roomId%10, roomId)
	return key
}

func (s *Service) generateSearchInfo(action string, table string, new *model.TableField, old *model.TableField) (ret map[string]interface{}, retByte map[string][]byte) {
	ret = make(map[string]interface{})
	ret["action"] = action
	ret["table"] = table
	//搜索字段转换
	newMap := make(map[string]interface{})
	newMap["id"] = new.RoomId
	newMap["short_id"] = new.ShortId
	newMap["uid"] = new.Uid
	newMap["uname"] = new.UName
	newMap["category"] = new.Area
	newMap["title"] = new.Title
	newMap["tag"] = new.Tag
	newMap["try_time"] = new.TryTime
	newMap["cover"] = new.Cover
	newMap["user_cover"] = new.UserCover
	newMap["lock_status"] = s.getLockStatus(new.LockStatus)
	newMap["hidden_status"] = s.getHiddenStatus(new.HiddenStatus)
	newMap["attentions"] = new.Attentions
	newMap["attention"] = new.Attentions
	newMap["online"] = new.Online
	newMap["live_time"] = new.LiveTime
	newMap["area_v2_id"] = new.AreaV2Id
	newMap["ord"] = new.AreaV2ParentId
	newMap["arcrank"] = new.Virtual
	newMap["lastupdate"] = s.getLastUpdate(new)
	newMap["is_live"] = s.getLiveStatus(new)
	newMap["s_category"] = new.AreaV2Name
	ret["new"] = newMap
	oldMap := make(map[string]interface{})
	if old != nil {
		oldMap["id"] = old.RoomId
		oldMap["short_id"] = old.ShortId
		oldMap["uid"] = old.Uid
		oldMap["uname"] = old.UName
		oldMap["category"] = old.Area
		oldMap["title"] = old.Title
		oldMap["tag"] = old.Tag
		oldMap["try_time"] = old.TryTime
		oldMap["cover"] = old.Cover
		oldMap["user_cover"] = old.UserCover
		oldMap["lock_status"] = s.getLockStatus(old.LockStatus)
		oldMap["hidden_status"] = s.getHiddenStatus(old.HiddenStatus)
		oldMap["attentions"] = old.Attentions
		oldMap["attention"] = old.Attentions
		oldMap["online"] = old.Online
		oldMap["live_time"] = old.LiveTime
		oldMap["area_v2_id"] = old.AreaV2Id
		oldMap["area_v2_name"] = old.AreaV2Name
		oldMap["ord"] = old.AreaV2ParentId
		oldMap["arcrank"] = old.Virtual
		oldMap["lastupdate"] = s.getLastUpdate(old)
		oldMap["is_live"] = s.getLiveStatus(old)
	}
	if action != "insert" && old == nil {
		oldMap["id"] = new.RoomId
		oldMap["short_id"] = new.ShortId
		oldMap["uid"] = new.Uid
		oldMap["uname"] = new.UName
		oldMap["category"] = new.Area
		oldMap["title"] = new.Title
		oldMap["tag"] = new.Tag
		oldMap["try_time"] = new.TryTime
		oldMap["cover"] = new.Cover
		oldMap["user_cover"] = new.UserCover
		oldMap["lock_status"] = s.getLockStatus(new.LockStatus)
		oldMap["hidden_status"] = s.getHiddenStatus(new.HiddenStatus)
		oldMap["attentions"] = new.Attentions
		oldMap["attention"] = new.Attentions
		oldMap["online"] = new.Online
		oldMap["live_time"] = new.LiveTime
		oldMap["area_v2_id"] = new.AreaV2Id
		oldMap["area_v2_name"] = new.AreaV2Name
		oldMap["ord"] = new.AreaV2ParentId
		oldMap["arcrank"] = new.Virtual
		oldMap["lastupdate"] = s.getLastUpdate(new)
		oldMap["is_live"] = s.getLiveStatus(new)
	}
	ret["old"] = oldMap

	newByteMap := make(map[string][]byte)
	newByteMap["id"] = []byte(strconv.Itoa(new.RoomId))
	newByteMap["short_id"] = []byte(strconv.Itoa(new.ShortId))
	newByteMap["uid"] = []byte(strconv.FormatInt(new.Uid, 10))
	newByteMap["uname"] = []byte(new.UName)
	newByteMap["category"] = []byte(strconv.Itoa(new.Area))
	newByteMap["title"] = []byte(new.Title)
	newByteMap["tag"] = []byte(new.Tag)
	newByteMap["try_time"] = []byte(new.TryTime)
	newByteMap["cover"] = []byte(new.Cover)
	newByteMap["user_cover"] = []byte(new.UserCover)
	newByteMap["lock_status"] = []byte(strconv.Itoa(s.getLockStatus(new.LockStatus)))
	newByteMap["hidden_status"] = []byte(strconv.Itoa(s.getHiddenStatus(new.HiddenStatus)))
	newByteMap["attentions"] = []byte(strconv.Itoa(new.Attentions))
	newByteMap["attention"] = []byte(strconv.Itoa(new.Attentions))
	newByteMap["online"] = []byte(strconv.Itoa(new.Online))
	newByteMap["live_time"] = []byte(new.LiveTime)
	newByteMap["area_v2_id"] = []byte(strconv.Itoa(new.AreaV2Id))
	newByteMap["ord"] = []byte(strconv.Itoa(new.AreaV2ParentId))
	newByteMap["arcrank"] = []byte(strconv.Itoa(new.Virtual))
	newByteMap["lastupdate"] = []byte(s.getLastUpdate(new))
	newByteMap["is_live"] = []byte(strconv.Itoa(s.getLiveStatus(new)))
	newByteMap["s_category"] = []byte(new.AreaV2Name)
	return ret, newByteMap
}

//获取直播状态
func (s *Service) getLiveStatus(roomInfo *model.TableField) int {
	if roomInfo.LiveTime != "0000-00-00 00:00:00" {
		return 1
	}

	if roomInfo.RoundStatus == 1 && roomInfo.OnFlag == 1 {
		return 2
	}

	return 0
}

//获取房间最后更新时间
func (s *Service) getLastUpdate(roomInfo *model.TableField) string {
	if roomInfo.MTime != "0000-00-00 00:00:00" {
		return roomInfo.MTime
	}
	return roomInfo.CTime
}

func (s *Service) getAreaV2Detail(areaV2Id int) (areaInfo *roomV1.AreaGetDetailResp_AreaInfo, err error) {
	areaResp, err := dao.RoomApi.V1Area.GetDetail(rpccontext.WithTimeout(context.TODO(), 50*time.Millisecond), &roomV1.AreaGetDetailReq{
		Id: int64(areaV2Id),
	})
	if err != nil {
		log.Error("[getAreaV2Detail]GetMultiple rpc error, error:%+v", err)
		return
	}
	if areaResp.Code != 0 {
		log.Error("[getAreaV2Detail]GetMultiple return error, code:%d, msg:%s", areaResp.Code, areaResp.Msg)
		err = errors.New("user return error")
		return
	}
	if areaResp.Data == nil {
		log.Error("[getAreaV2Detail]GetMultiple empty data")
		err = errors.New("area detail empty error")
		return
	}

	return areaResp.Data, err
}
