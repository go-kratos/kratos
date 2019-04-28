package dao

import (
	"context"
	v1pb "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/service/live/xuser/api/grpc/v1"
	account "go-common/app/service/main/account/model"
	"go-common/library/log"
	"math"
)

// GetByRoomIDPage get roomadmin list by roomid.
func (d *Dao) GetByRoomIDPage(c context.Context, roomID int64, page int64, pageSize int64) (resp *v1pb.RoomAdminGetByRoomResp, err error) {
	resp = &v1pb.RoomAdminGetByRoomResp{}

	resp.Page = &v1pb.RoomAdminGetByRoomResp_Page{
		Page:       page,
		PageSize:   pageSize,
		TotalPage:  1,
		TotalCount: 0,
	}

	ret, err := d.RoomAdminAPI.GetByRoom(c, &v1.RoomAdminGetByRoomReq{
		Roomid: roomID,
	})
	log.Info("call GetByAnchor mid(%v) page (%v) ret(%v)", roomID, page, ret)

	if err != nil {
		log.Error("call GetByAnchor mid(%v) page (%v) error(%v)", roomID, page, err)
		return
	}
	if ret == nil || ret.Data == nil {
		return
	}

	mids, _ := d.getUidsFromAdmins(c, ret.Data)
	dataLen := int64(len(mids))

	d.setPager(dataLen, pageSize, page, resp.Page)
	begin, end := d.getRange(pageSize, page, dataLen)

	args := &account.ArgMids{Mids: mids[begin:end]}

	accData, err := d.acc.Infos3(c, args)
	if err != nil {
		log.Error("call account.Infos3(%v) error(%v)", args, err)
		return resp, err
	}

	for _, v := range ret.Data[begin:end] {
		item := &v1pb.RoomAdminGetByRoomResp_Data{
			Uid:   v.Uid,
			Ctime: v.Ctime,
		}

		if _, ok := accData[item.Uid]; ok {
			item.Face = accData[item.Uid].Face
			item.Uname = accData[item.Uid].Name
		} else {
			log.Error("没有这个人的用户信息 uid(%v) data(%v)", item.Uid, accData)
		}
		resp.Data = append(resp.Data, item)
	}
	return
}

// getAnchorUidsFromAdmins .
// 根据批量房管获取对应主播的房间号和UID
func (d *Dao) getUidsFromAdmins(c context.Context, admins []*v1.RoomAdminGetByRoomResp_Data) (uids []int64, err error) {
	if len(admins) == 0 {
		return
	}

	for _, r := range admins {
		uids = append(uids, r.Uid)
	}

	return
}

// pageSize, page, dataLen
func (d *Dao) getRange(pageSize int64, page int64, dataLen int64) (begin int64, end int64) {
	begin = (page - 1) * pageSize
	end = page * pageSize

	if end > dataLen {
		end = dataLen
	}

	if begin > end {
		begin = end
	}
	return
}

func (d *Dao) setPager(dataLen int64, pageSize int64, page int64, pager *v1pb.RoomAdminGetByRoomResp_Page) {
	pager.Page = page
	pager.TotalCount = dataLen
	pager.TotalPage = int64(math.Ceil(float64(dataLen) / float64(pageSize)))
	pager.PageSize = pageSize
}
