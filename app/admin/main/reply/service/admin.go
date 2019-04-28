package service

import (
	"go-common/app/admin/main/reply/model"
	"go-common/library/log"

	"context"
	"time"
)

func (s *Service) addAdminLog(c context.Context, oid, rpID, adminID int64, typ, isNew, isReport, state int32, result, remark string, now time.Time) (err error) {
	rpMap := map[int64]*model.Reply{rpID: &model.Reply{
		Oid: oid,
		ID:  rpID,
	}}
	return s.addAdminLogs(c, rpMap, adminID, typ, isNew, isReport, state, result, remark, now)
}

func (s *Service) addAdminIDLogs(c context.Context, oids []int64, rpIDs []int64, adminID int64, typ, isNew, isReport, state int32, result, remark string, now time.Time) (err error) {
	if _, err = s.dao.UpAdminNotNew(c, rpIDs, now); err != nil {
		log.Error("s.dao.UpAdminNotNew(%v) error(%v)", rpIDs, err)
	}
	if _, err = s.dao.AddAdminLog(c, oids, rpIDs, adminID, typ, isNew, isReport, state, result, remark, now); err != nil {
		log.Error("s.dao.AddAdminLog(admin:%d, oid:%v rpID:%v type:%d result:%s remark:%s isReport:%d state:%d) error(%v)", adminID, oids, rpIDs, typ, result, remark, isReport, state, err)
	}
	return
}

func (s *Service) addAdminLogs(c context.Context, rps map[int64]*model.Reply, adminID int64, typ, isNew, isReport, state int32, result, remark string, now time.Time) (err error) {
	if len(rps) == 0 {
		return
	}
	rpIDs := make([]int64, 0, len(rps))
	oids := make([]int64, 0, len(rps))
	for _, rp := range rps {
		rpIDs = append(rpIDs, rp.ID)
		oids = append(oids, rp.Oid)
	}
	if _, err = s.dao.UpAdminNotNew(c, rpIDs, now); err != nil {
		log.Error("s.dao.UpAdminNotNew(%v) error(%v)", rpIDs, err)
	}
	if _, err = s.dao.AddAdminLog(c, oids, rpIDs, adminID, typ, isNew, isReport, state, result, remark, now); err != nil {
		log.Error("s.dao.AddAdminLog(admin:%d, oid:%v rpID:%v type:%d result:%s remark:%s isReport:%d state:%d) error(%v)", adminID, oids, rpIDs, typ, result, remark, isReport, state, err)
	}
	return
}

// ReplyAdminLog ReplyAdminLog
type ReplyAdminLog struct {
	model.AdminLog
	AdminName string `json:"admin_name"`
}

// LogsByRpID get log by reply id.
func (s *Service) LogsByRpID(c context.Context, rpID int64) (res []ReplyAdminLog, err error) {
	var resDao []*model.AdminLog
	res = []ReplyAdminLog{}
	if resDao, err = s.dao.AdminLogsByRpID(c, rpID); err != nil {
		log.Error("s.dao.AdminLogsByRpID(%d) error(%v)", rpID, err)
		return
	}
	admins := map[int64]string{}
	for _, data := range resDao {
		admins[data.AdminID] = ""
		res = append(res, ReplyAdminLog{*data, ""})
	}
	s.dao.AdminName(c, admins)
	for i := range res {
		res[i].AdminName = admins[res[i].AdminID]
	}
	return
}
