package roomAdmin

import (
	"context"
	"go-common/app/service/live/xuser/model"
	"go-common/library/log"
)

// GetByUserMysql get admins by user .
func (dao *Dao) GetByUserMysql(c context.Context, uid int64) (admins []*model.RoomAdmin, err error) {

	err = dao.orm.Model(&model.RoomAdmin{}).Find(&admins, "uid=?", uid).Order("ctime DESC").Error

	if err != nil {
		log.Error("GetByUserMysql (%v) error(%v)", uid, err)
		return nil, err
	}
	return
}

// GetByRoomIdMysql get admins by roomId.
func (dao *Dao) GetByRoomIdMysql(c context.Context, roomId int64) (admins []*model.RoomAdmin, err error) {

	err = dao.orm.Table("ap_room_admin").Model(&model.RoomAdmin{}).Find(&admins, "roomid=?", roomId).Order("ctime DESC").Error

	if err != nil {
		log.Error("GetByUserMysql (%v) error(%v)", roomId, err)
		return nil, err
	}
	return
}

// DelDbAdminMysql delete a admin .
func (dao *Dao) DelDbAdminMysql(c context.Context, id int64) (err error) {
	param := &model.RoomAdmin{
		Id: id,
	}
	//spew.Dump("deldbAdmin", param);
	err = dao.orm.Table("ap_room_admin").Model(&model.RoomAdmin{}).Delete(param).Error

	if err != nil {
		log.Error("DelDbAdminMysql (%v) error(%v)", id, err)
		return err
	}

	return
}

// AddAdminMysql delete a admin .
func (dao *Dao) AddAdminMysql(c context.Context, uid int64, roomId int64) (err error) {
	param := &model.RoomAdmin{
		Uid:    uid,
		Roomid: roomId,
	}
	//spew.Dump("AddAdminMysql", param)
	err = dao.orm.Model(&model.RoomAdmin{}).Save(param).Error

	if err != nil {
		log.Error("AddAdminMysql uid (%v) roomid (%v)error(%v)", uid, roomId, err)
		return err
	}

	return
}

// GetByRoomIdUidMysql delete a admin .
func (dao *Dao) GetByRoomIdUidMysql(c context.Context, uid int64, roomId int64) (resp []*model.RoomAdmin, err error) {
	//param := &model.RoomAdmin{
	//	Uid:    uid,
	//	Roomid: roomId,
	//}
	err = dao.orm.Model(&model.RoomAdmin{}).Where("uid=? AND roomid =?", uid, roomId).Find(&resp).Error
	//spew.Dump("GetByRoomIdUidMysql", resp)

	if err != nil {
		log.Error("AddAdminMysql uid (%v) roomid (%v)error(%v)", uid, roomId, err)
		return resp, err
	}

	return
}
