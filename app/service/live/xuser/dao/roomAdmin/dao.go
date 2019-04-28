package roomAdmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	banned "go-common/app/service/live/banned_service/api/liverpc/v1"
	v12 "go-common/app/service/live/fans_medal/api/liverpc/v2"
	"go-common/app/service/live/room/api/liverpc/v1"
	"go-common/app/service/live/room/api/liverpc/v2"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/dao"
	"go-common/app/service/live/xuser/model"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Dao dao
type Dao struct {
	c               *conf.Config
	mc              *memcache.Pool
	redis           *redis.Pool
	db              *xsql.DB
	orm             *gorm.DB
	RoomAdminExpire int32
	// acc rpc
	acc    *accrpc.Service3
	client *bm.Client
}

const (
	userPrefix   = "up_v1_%d"
	roomPrefix   = "rp_v1_%d"
	mcExpire     = 3600
	maxAdminsNum = 100
)

// KeyUser return the mc key by user mid.
func KeyUser(uid int64) string {
	return fmt.Sprintf(userPrefix, uid)
}

// KeyRoom return the mc key by anchor mid.
func KeyRoom(uid int64) string {
	return fmt.Sprintf(roomPrefix, uid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// 获取主播的房管列表
	// mc: -key=KeyRoom
	CacheRoomAdminRoom(c context.Context, anchor int64) ([]*model.RoomAdmin, error)
	// 获取用户的房管列表
	// mc: -key=KeyUser
	CacheRoomAdminUser(c context.Context, user int64) ([]*model.RoomAdmin, error)

	// mc: -key=KeyRoom -expire=d.RoomAdminExpire -encode=json|gzip
	AddCacheKeyAnchorRoom(c context.Context, anchor int64, value []*model.RoomAdmin) error
	// mc: -key=KeyUser -expire=d.RoomAdminExpire -encode=gob
	AddCacheRoomAdminUser(c context.Context, user int64, value []*model.RoomAdmin) error

	// mc: -key=KeyRoom
	DelCacheKeyAnchorRoom(c context.Context, anchor int64) error
	// mc: -key=KeyUser
	DelCacheRoomAdminUser(c context.Context, user int64) error
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:               c,
		mc:              memcache.NewPool(c.Memcache),
		redis:           redis.NewPool(c.Redis),
		db:              xsql.NewMySQL(c.LiveAppMySQL),
		orm:             orm.NewMySQL(c.LiveAppORM),
		RoomAdminExpire: mcExpire,
		acc:             accrpc.New3(c.AccountRPC),
		client:          bm.NewClient(c.BMClient),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.orm.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}

// HasAnyAdmin whether he has any admin in any room.
func (d *Dao) HasAnyAdmin(c context.Context, uid int64) (int64, error) {
	noAdmin := int64(0)
	hasAdmin := int64(1)

	rst, err := d.GetAllByUid(c, uid)

	if nil == rst {
		return noAdmin, err
	}
	return hasAdmin, err
}

// GetByUidPage get admins by uid and page.
func (d *Dao) GetByUidPage(c context.Context, uid int64, page int64, pageSize int64) (resp *v1pb.RoomAdminGetByUidResp, err error) {
	resp = &v1pb.RoomAdminGetByUidResp{}

	resp.Page = &v1pb.RoomAdminGetByUidResp_Page{
		Page:       page,
		PageSize:   pageSize,
		TotalPage:  1,
		TotalCount: 0,
	}

	rst, err := d.GetAllByUid(c, uid)
	//spew.Dump("GetAllByUid", rst, err)

	if err != nil {
		return
	}

	if rst == nil {
		return
	}

	sort.Sort(sort.Reverse(model.RoomAdmins(rst)))
	//spew.Dump(rst)
	resp.Page.TotalCount = int64(len(rst))
	resp.Page.PageSize = pageSize
	resp.Page.TotalPage = int64(math.Ceil(float64(len(rst)) / float64(resp.Page.PageSize)))

	begin := (page - 1) * pageSize
	end := page * pageSize
	if page*pageSize > int64(len(rst)) {
		end = int64(len(rst))
	}

	roomUid, mids, err := d.getAnchorUidsFromAdmins(c, rst)
	// 没获取到房间信息
	if err != nil {
		return resp, err
	}

	// 可能从room获取主播信息不全
	if int64(len(mids)) < end {
		end = int64(len(mids))
	}

	if begin > end {
		begin = end
	}

	if err != nil {
		return resp, err
	}

	args := &account.ArgMids{Mids: mids[begin:end]}

	accData, err := d.acc.Infos3(c, args)
	//spew.Dump("d.acc.Infos3", accData, err)

	if err != nil {
		log.Error("call account.Infos3(%v) error(%v)", args, err)
		return resp, err
	}

	for _, v := range rst[begin:end] {
		item := &v1pb.RoomAdminGetByUidResp_Data{
			Uid:    v.Uid,
			Roomid: v.Roomid,
			Ctime:  v.Ctime.Time().Format("2006-01-02 15:04:05"),
		}
		if _, ok := roomUid[v.Roomid]; ok {
			item.AnchorId = roomUid[v.Roomid]

			if _, ok := accData[item.AnchorId]; ok {
				item.AnchorCover = accData[item.AnchorId].Face
				item.Uname = accData[item.AnchorId].Name
			} else {
				log.Error("没有这个人的用户信息 uid(%v) data(%v)", item.AnchorId, accData)
			}
		} else {
			log.Error("没有这个人的房间信息 room (%v) data(%v)", v.Roomid, roomUid)
		}

		resp.Data = append(resp.Data, item)
	}

	//spew.Dump("resp.Data", resp.Data)

	return
}

// GetAllByUid get by uid.
func (d *Dao) GetAllByUid(c context.Context, uid int64) ([]*model.RoomAdmin, error) {
	rstMc, err := d.CacheRoomAdminUser(c, uid)

	//spew.Dump("HasAnyAdmin1", rstMc, err)
	//spew.Dump("lenMc", len(rstMc))

	if err != nil {
		return nil, err
	}

	// 空缓存标识
	if rstMc != nil {
		if rstMc[0].Id == -1 {
			return nil, err
		}
		return rstMc, err
	}

	rstDb, err := d.GetByUserMysql(c, uid)

	if err != nil {
		return nil, err
	}

	if len(rstDb) == 0 {
		d.AddCacheNoneUser(c, uid)
		return nil, err
	}
	d.AddCacheRoomAdminUser(c, uid, rstDb)

	return rstDb, err
}

// getAnchorUidsFromAdmins .
// 根据批量房管获取对应主播的房间号和UID
func (d *Dao) getAnchorUidsFromAdmins(c context.Context, admins []*model.RoomAdmin) (roomUid map[int64]int64, uids []int64, err error) {
	var roomIds []int64
	roomUid = make(map[int64]int64)

	for _, a := range admins {
		roomIds = append(roomIds, a.Roomid)
	}

	if len(roomIds) == 0 {
		return
	}

	reply, err := dao.RoomAPI.V2Room.GetByIds(c, &v2.RoomGetByIdsReq{Ids: roomIds})
	if err != nil {
		log.Error("dao.RoomAPI.V2Room.GetByIds (%v) error(%v)", roomIds, err)
		return roomUid, uids, err
	}

	if reply.GetCode() != 0 {
		err = ecode.Int(int(reply.GetCode()))
		log.Error("dao.RoomAPI.V2Room.GetByIds (%v) error code(%v)", roomIds, err)
		return roomUid, uids, err
	}

	for aRoomId, r := range reply.Data {
		roomUid[aRoomId] = r.Uid
		uids = append(uids, r.Uid)
	}

	return
}

// Del delete a roomadmin
func (d *Dao) Del(c context.Context, uid int64, roomId int64) (err error) {

	if err = d.DelAllCache(c, uid, roomId); err != nil {
		log.Error("DelAllCache(%v) uid (%v) roomid (%v) error(%v)", uid, roomId, err)
		return
	}

	admin, err := d.GetByRoomIdUidMysql(c, uid, roomId)
	if len(admin) == 0 {
		log.Error("GetByRoomIdUidMysql empty uid(%v) roomId (%v)error(%v)", uid, roomId, err)
		return
	}

	//spew.Dump("GetByRoomIdUidMysql", admin)

	if err = d.DelDbAdminMysql(c, admin[0].Id); err != nil {
		log.Error("DelCacheRoomAdminUser uid(%v) roomId (%v) error(%v)", uid, roomId, err)
		return err
	}

	return
}

func (d *Dao) getInfoByName(c context.Context, name string) (userInfo *account.Info, err error) {
	userInfo = &account.Info{}

	infosByName, err := d.acc.InfosByName3(c, &account.ArgNames{
		Names: []string{name},
	})

	if err != nil {
		log.Error("d.acc.InfosByName3(%v) error(%v)", name, err)
		return userInfo, err
	}
	log.Info("d.acc.InfosByName3(%v) return (%v)", name, infosByName)
	if len(infosByName) != 0 {
		for _, info := range infosByName {
			return info, err
		}
	}
	return
}

// SearchForAdmin search user list by keyword.
func (d *Dao) SearchForAdmin(c context.Context, keyword string, anchorId int64) (resp []*v1pb.RoomAdminSearchForAdminResp_Data, err error) {
	isUid := 0
	matchUid, _ := strconv.ParseInt(keyword, 10, 64)

	if matchUid != 0 && keyword != "0" {
		isUid = 1
	}

	// get by name
	infoMatchName, err := d.getInfoByName(c, keyword)
	if err != nil {
		return resp, err
	}

	if nil != infoMatchName && infoMatchName.Mid != 0 {
		log.Info("SearchForAdmin infoMatchName keyword (%v) ret (%v)", keyword, infoMatchName)
		isAdminName, _ := d.isAdminByUid(c, infoMatchName.Mid, anchorId)
		medalData, errMedal := d.getMedalInfoByUids(c, []int64{infoMatchName.Mid})

		if errMedal != nil {
			return resp, errMedal
		}

		itemName := &v1pb.RoomAdminSearchForAdminResp_Data{
			Uid:     infoMatchName.Mid,
			IsAdmin: isAdminName,
			Uname:   infoMatchName.Name,
			Face:    infoMatchName.Face,
		}
		if _, ok := medalData[infoMatchName.Mid]; ok {
			itemName.MedalName = medalData[infoMatchName.Mid].MedalName
			itemName.Level = medalData[infoMatchName.Mid].Level
		} else {
			log.Info("没有这个人的勋章信息 uid(%v) data(%v)", infoMatchName.Mid, medalData)
		}

		resp = append(resp, itemName)
	}

	//spew.Dump("searchForadmin2", resp)
	// just name
	if 0 == isUid {
		return resp, nil
	}

	// get by uid
	infoMatchUid, err := d.getInfoByUid(c, matchUid)
	if err != nil {
		return resp, err
	}

	if infoMatchUid == nil {
		return resp, nil
	}

	isAdminUid, _ := d.isAdminByUid(c, matchUid, anchorId)

	medalDataUid, err := d.getMedalInfoByUids(c, []int64{infoMatchUid.Mid})

	if err != nil {
		return resp, err
	}

	itemUid := &v1pb.RoomAdminSearchForAdminResp_Data{
		Uid:     infoMatchUid.Mid,
		IsAdmin: isAdminUid,
		Uname:   infoMatchUid.Name,
		Face:    infoMatchUid.Face,
	}
	if _, ok := medalDataUid[infoMatchUid.Mid]; ok {
		itemUid.MedalName = medalDataUid[infoMatchUid.Mid].MedalName
		itemUid.Level = medalDataUid[infoMatchUid.Mid].Level
	}

	resp = append(resp, itemUid)
	//spew.Dump("searchForadmin2", resp)

	return
}

func (d *Dao) isAdminByUid(c context.Context, uid int64, anchorId int64) (rst int64, err error) {

	roomId, err := d.getRoomIdByUid(c, anchorId)

	if err != nil {
		return rst, err
	}
	if 0 == roomId {
		return rst, nil
	}

	return d.IsAdminByRoomId(c, uid, roomId)
}

func (d *Dao) getInfoByUid(c context.Context, uid int64) (info *account.Info, err error) {
	info, err = d.acc.Info3(c, &account.ArgMid{
		Mid: uid,
	})

	if err != nil {
		log.Error("d.acc.Info3(%v) error(%v)", uid, err)
		return info, err
	}

	log.Info("d.acc.Info3(%v) return (%v)", uid, info)
	return
}

func (d *Dao) getRoomInfoByUid(c context.Context, uid int64) (roomInfo *v1.RoomGetStatusInfoByUidsResp_RoomInfo, err error) {
	roomInfo = &v1.RoomGetStatusInfoByUidsResp_RoomInfo{}

	reply, err := dao.RoomAPI.V1Room.GetStatusInfoByUids(c, &v1.RoomGetStatusInfoByUidsReq{
		Uids:       []int64{uid},
		ShowHidden: 1,
	})
	if err != nil {
		log.Error("dao.RoomAPI.V1Room.GetStatusInfoByUids (%v) error(%v)", uid, err)
		return roomInfo, err
	}

	if reply.GetCode() != 0 {
		err = ecode.Int(int(reply.GetCode()))
		log.Error("dao.RoomAPI.V2Room.GetByIds (%v) error code(%v)", uid, err)
		return roomInfo, err
	}

	if len(reply.Data) == 0 {
		return
	}

	for aUid, aInfo := range reply.Data {
		if aUid == uid {
			return aInfo, nil
		}
	}
	return
}

func (d *Dao) getMedalInfoByUids(c context.Context, uids []int64) (medalInfo map[int64]*v12.AnchorQueryLiveWearingResp_Medal, err error) {
	medalInfo = make(map[int64]*v12.AnchorQueryLiveWearingResp_Medal)

	reply, err := dao.FansMedalAPI.V2Anchor.QueryLiveWearing(c, &v12.AnchorQueryLiveWearingReq{
		UidList: uids,
	})
	log.Info("call dao.FansMedalAPI.V2Anchor.QueryLiveWearing (%v) rst (%v)", uids, reply)

	if err != nil {
		log.Error("dao.FansMedalAPI.V2Anchor.QueryLiveWearing (%v) error(%v)", uids, err)
		return medalInfo, err
	}

	if reply.GetCode() != 0 {
		err = ecode.Int(int(reply.GetCode()))
		log.Error("dao.RoomAPI.V2Room.GetByIds (%v) error code(%v)", uids, err)
		return medalInfo, err
	}

	if len(reply.Data) == 0 {
		return
	}

	return reply.Data, err
}

// IsAdminByRoomId ...
func (d *Dao) IsAdminByRoomId(c context.Context, uid int64, roomId int64) (rst int64, err error) {
	rst = 0
	admins, err := d.GetAllByRoomId(c, roomId)

	if err != nil {
		log.Error("GetAllByRoomId(%v) error(%v)", roomId, err)
		return rst, err
	}

	if len(admins) == 0 {
		return rst, nil
	}

	for _, v := range admins {
		if v.Uid == uid {
			rst = 1
			return rst, nil
		}
	}

	return
}

// GetByAnchorIdPage get by anchor id and page .
func (d *Dao) GetByAnchorIdPage(c context.Context, anchorId int64, page int64, pageSize int64) (resp *v1pb.RoomAdminGetByAnchorResp, err error) {
	resp = &v1pb.RoomAdminGetByAnchorResp{}

	resp.Page = &v1pb.RoomAdminGetByAnchorResp_Page{
		Page:       page,
		PageSize:   pageSize,
		TotalPage:  1,
		TotalCount: 0,
	}

	roomId, err := d.getRoomIdByUid(c, anchorId)

	if err != nil {
		return resp, err
	}
	if 0 == roomId {
		return resp, nil
	}

	allAdmins, err := d.GetAllByRoomId(c, roomId)
	//spew.Dump("GetAllByUid", allAdmins, err)

	if err != nil {
		return resp, err
	}

	if allAdmins == nil {
		return resp, nil
	}

	sort.Sort(sort.Reverse(model.RoomAdmins(allAdmins)))
	//spew.Dump(allAdmins)
	resp.Page.TotalCount = int64(len(allAdmins))
	resp.Page.PageSize = pageSize
	resp.Page.TotalPage = int64(math.Ceil(float64(len(allAdmins)) / float64(resp.Page.PageSize)))

	begin := (page - 1) * pageSize
	end := page * pageSize
	if page*pageSize > int64(len(allAdmins)) {
		end = int64(len(allAdmins))
	}

	uids, _ := d.getUidsFromAdmins(c, allAdmins)
	if int64(len(uids)) < end {
		end = int64(len(uids))
	}
	if begin > end {
		begin = end
	}
	//spew.Dump("getAnchorUidsFromAdmins", uids, err)

	accArgs := &account.ArgMids{Mids: uids[begin:end]}
	accData, err := d.acc.Infos3(c, accArgs)
	if err != nil {
		log.Error("d.acc.Infos3(%v) error(%v)", accArgs, err)
	}

	medalData, err := d.getMedalInfoByUids(c, uids)

	if err != nil {
		log.Error("d.getMedalInfoByUids(%v) error(%v)", uids, err)
		return resp, err
	}

	if err != nil {
		log.Error("call account.Infos3(%v) error(%v)", accArgs, err)
		return resp, err
	}

	for _, v := range allAdmins[begin:end] {
		item := &v1pb.RoomAdminGetByAnchorResp_Data{
			Uid:    v.Uid,
			Ctime:  v.Ctime.Time().Format("2006-01-02 15:04:05"),
			Roomid: v.Roomid,
		}
		if _, ok := accData[v.Uid]; ok {
			item.Uname = accData[v.Uid].Name
			item.Face = accData[v.Uid].Face
		} else {
			log.Error("没有这个人的用户信息 uid(%v) data(%v)", v.Uid, accData)
		}
		if _, ok := medalData[v.Uid]; ok {
			item.Level = medalData[v.Uid].Level
			item.MedalName = medalData[v.Uid].MedalName
		} else {
			log.Info("没有这个人的勋章信息 uid(%v) data(%v)", v.Uid, medalData)
		}
		resp.Data = append(resp.Data, item)
	}

	//spew.Dump("resp.Data", resp.Data)

	return
}

// GetAllByRoomId get by uid.
func (d *Dao) GetAllByRoomId(c context.Context, roomId int64) ([]*model.RoomAdmin, error) {
	rstMc, err := d.CacheRoomAdminRoom(c, roomId)

	//spew.Dump("HasAnyAdmin1", rstMc, err)
	//spew.Dump("lenMc", len(rstMc))

	if err != nil {
		return nil, err
	}

	// 空缓存标识
	if rstMc != nil {
		if rstMc[0].Id == -1 {
			return nil, err
		}
		return rstMc, err
	}

	rstDb, err := d.GetByRoomIdMysql(c, roomId)

	if err != nil {
		return nil, err
	}

	if len(rstDb) == 0 {
		d.AddCacheNoneRoom(c, roomId)
		return nil, err
	}
	d.AddCacheKeyAnchorRoom(c, roomId, rstDb)

	return rstDb, err
}

// getUidsFromAdmins .
// 返回房管列表中的uid
func (d *Dao) getUidsFromAdmins(c context.Context, admins []*model.RoomAdmin) (uids []int64, err error) {

	for _, a := range admins {
		uids = append(uids, a.Uid)
	}

	return
}

// DismissAnchor del a admin
func (d *Dao) DismissAnchor(c context.Context, uid int64, anchorId int64) (resp *v1pb.RoomAdminDismissAdminResp, err error) {
	resp = &v1pb.RoomAdminDismissAdminResp{}
	roomId, err := d.getRoomIdByUid(c, anchorId)

	if err != nil {
		return resp, err
	}
	if 0 == roomId {
		return resp, nil
	}

	isAdmin, err := d.IsAdminByRoomId(c, uid, roomId)
	if err != nil {
		log.Error("IsAdminByRoomId uid(%v) roomid (%v) error(%v)", uid, roomId, err)
		return resp, err
	}

	if 0 == isAdmin {
		err = ecode.Error(ecode.XUserAddRoomAdminNotAdminError, "该用户已经不是房管啦")
		return
	}

	err = d.Del(c, uid, roomId)

	if err != nil {
		log.Error("getRoomInfoByUid uid (%v) roomid (%v) error(%v)", uid, roomId, err)
		return
	}

	return
}

// DelAllCache delete cache .
func (d *Dao) DelAllCache(c context.Context, uid int64, roomId int64) (err error) {
	if err = d.DelCacheKeyAnchorRoom(c, roomId); err != nil {
		log.Error("DelCacheKeyAnchorRoom(%v) error(%v)", roomId, err)
		return err
	}

	if err = d.DelCacheRoomAdminUser(c, uid); err != nil {
		log.Error("DelCacheRoomAdminUser(%v) error(%v)", uid, err)
		return err
	}
	return
}

// Add add a admin
func (d *Dao) Add(c context.Context, uid int64, anchorId int64) (resp *v1pb.RoomAdminAddResp, err error) {
	resp = &v1pb.RoomAdminAddResp{}

	roomId, err := d.getRoomIdByUid(c, anchorId)

	if err != nil {
		return resp, err
	}
	if 0 == roomId {
		return resp, nil
	}

	allRoomAdmin, err := d.GetAllByRoomId(c, roomId)
	//spew.Dump("Add", roomId, allRoomAdmin)
	if err != nil {
		return resp, err
	}

	if len(allRoomAdmin) >= maxAdminsNum {
		err = ecode.Error(ecode.XUserAddRoomAdminOverLimitError, "最多设置100个房间管理员")
		//err = &pb.Error{
		//	ErrCode:    2,
		//	ErrMessage: "最多设置100个房间管理员",
		//}
		return
	}

	isAdmin, err := d.IsAdminByRoomId(c, uid, roomId)
	if err != nil {
		log.Error("IsAdminByRoomId uid(%v) roomid (%v) error(%v)", uid, roomId, err)
		return resp, err
	}

	if 1 == isAdmin {
		err = ecode.Error(ecode.XUserAddRoomAdminIsAdminError, "该用户已经是你的房管啦")
		//err = &pb.Error{
		//	ErrCode:    1,
		//	ErrMessage: "他已经是房管",
		//}
		return
	}

	banArg := &banned.SilentMngIsBlockUserReq{
		Uid:    uid,
		Roomid: roomId,
		Type:   1,
	}

	retBan, err := dao.BannedAPI.V1SilentMng.IsBlockUser(c, banArg)
	if err != nil {
		log.Error("call dao.BannedAPI.V1SilentMng.IsBlockUser(%v) error(%v)", banArg, err)
		return
	}
	if retBan.Code != 0 || nil == retBan.Data {
		log.Error("call dao.BannedAPI.V1SilentMng.IsBlockUser(%v) error return (%v)", banArg, retBan)
	}

	if retBan.Data.GetIsBlockUser() {
		err = ecode.Error(ecode.XUserAddRoomAdminIsSilentError, "他已经被禁言，无法添加房管")
		//err = &pb.Error{
		//	ErrCode:    3,
		//	ErrMessage: "他已经被禁言，无法添加房管",
		//}
		return
	}

	if err = d.DelAllCache(c, uid, roomId); err != nil {
		log.Error("DelCacheKeyAnchorRoom(%v) error(%v)", roomId, err)
		return resp, err
	}

	if err = d.AddAdminMysql(c, uid, roomId); err != nil {
		log.Error("DelCacheKeyAnchorRoom(%v) error(%v)", roomId, err)
		return
	}

	resp.Uid = uid
	resp.Roomid = roomId
	resp.Userinfo = &v1pb.RoomAdminAddResp_UI{}

	resp.Userinfo.Uid = uid

	info, _ := d.getInfoByUid(c, uid)

	if info != nil {
		resp.Userinfo.Uname = info.Name
	}

	d.adminChange(c, uid, roomId)
	return
}

func (d *Dao) getRoomIdByUid(c context.Context, anchorId int64) (roomId int64, err error) {
	roomInfo, err := d.getRoomInfoByUid(c, anchorId)
	if err != nil {
		log.Error("getRoomInfoByUid (%v) error(%v)", anchorId, err)

		return roomId, err
	}
	if roomInfo == nil {
		return roomId, nil
	}

	return roomInfo.RoomId, nil
}

// adminChange send broadcast
func (d *Dao) adminChange(c context.Context, uid int64, roomId int64) (err error) {
	roomInfo, err := dao.RoomAPI.V2Room.GetByIds(c, &v2.RoomGetByIdsReq{
		Ids:       []int64{roomId},
		NeedUinfo: 1,
	})
	if err != nil {
		log.Error("dao.RoomAPI.V2Room.GetByIds(%v) error(%v)", roomId, err)
		return
	}

	if roomInfo.Code != 0 || 0 == len(roomInfo.Data) {
		log.Error("dao.RoomAPI.V1Room.GetInfoById(%v) error code (%v) data (%v)", roomId, roomInfo.Code, roomInfo.Data)
		return
	}

	postJson := make(map[string]interface{})
	postJson["cmd"] = "room_admin_entrance"
	postJson["uid"] = uid
	postJson["msg"] = "系统提示：你已被主播设为房管"

	if err = d.sendBroadcastRoom(roomId, postJson); err != nil {
		return err
	}
	admins, err := d.GetAllByRoomId(c, roomId)
	if err != nil {
		return err
	}

	var adminUids []int64
	//	adminUids := make([]int64, 100)
	for _, v := range admins {
		adminUids = append(adminUids, v.Uid)
	}

	postJson2 := make(map[string]interface{})
	postJson2["cmd"] = "ROOM_ADMINS"
	postJson2["uids"] = adminUids
	if err = d.sendBroadcastRoom(roomId, postJson2); err != nil {
		return err
	}
	return
}

// sendBroadcastRoom .
func (d *Dao) sendBroadcastRoom(roomid int64, postJson map[string]interface{}) (err error) {

	log.Info("send reward broadcast begin:%d", roomid)

	var endPoint = fmt.Sprintf("http://live-dm.bilibili.co/dm/1/push?cid=%d&ensure=1", roomid)

	bytesData, err := json.Marshal(postJson)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", postJson, err)
		return
	}

	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(bytesData))

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Error("http.NewRequest(%v) url(%v) error(%v)", postJson, endPoint, err)
		return
	}

	client := http.Client{
		Timeout: time.Second,
	}

	// use httpClient to send request
	response, err := client.Do(req)

	if err != nil {
		log.Error("sending request to API endpoint(%v) error(%v)", req, err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("parse resp body(%v) error(%v)", body, err)
	}

	log.Info("send reward broadcast end:%d", roomid)

	return
}

// IsAdmin return whether a user is admin.
func (d *Dao) IsAdmin(c context.Context, uid int64, anchorId int64, roomId int64) (resp *v1pb.RoomAdminIsAdminResp, err error) {
	resp = &v1pb.RoomAdminIsAdminResp{}

	isAdmin, err := d.IsAdminByRoomId(c, uid, roomId)
	if err != nil {
		log.Error("IsAdminByRoomId uid(%v) roomid (%v) error(%v)", uid, roomId, err)
		return resp, err
	}

	if 0 == isAdmin {
		if uid != anchorId {
			err = ecode.Error(120014, "Ta不是该主播的房管")
			//err = ecode.Error{
			//	ErrCode:    120014, // 接口迁移, code 为老业务error code
			//	ErrMessage: "Ta不是该主播的房管",
			//}
			return
		}
	}

	userInfo, err := d.getInfoByUid(c, uid)
	if err != nil {
		return resp, err
	}
	if userInfo != nil {
		resp.Userinfo = &v1pb.RoomAdminIsAdminResp_UI{}
		resp.Roomid = roomId
		resp.Uid = uid
		resp.Userinfo.Uid = uid
		resp.Userinfo.Uname = userInfo.Name
		return
	}
	return
}
