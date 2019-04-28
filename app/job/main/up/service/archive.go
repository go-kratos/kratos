package service

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"go-common/app/admin/main/up/util/mathutil"
	"go-common/app/job/main/up/model"
	"go-common/app/job/main/up/model/upcrmmodel"
	upGRPCv1 "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"go-common/app/job/main/up/model/archivemodel"
)

// action
const (
	ActionUpdate = "update"
	ActionInsert = "insert"
)

// table name
const (
	TableArchiveStaff = "archive_staff"
)

//ArchiveUpInfo .
type ArchiveUpInfo struct {
	Table  string                     `json:"table"`
	Action string                     `json:"action"`
	New    *archivemodel.ArchiveCanal `json:"new"`
	Old    *archivemodel.ArchiveCanal `json:"old"`
}

// CanalMsg canal message struct
type CanalMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// handle notify t message
func (s *Service) handleArchiveNotifyT(msg *databus.Message) (err error) {
	m := &ArchiveUpInfo{}
	if err = json.Unmarshal(msg.Value, m); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
		return
	}

	switch m.Table {
	case "archive":
		err = s.checkArchiveNotify(m)
	}
	if err != nil {
		log.Error("handle msg fail, err=%s", err)
	}
	return
}
func (s *Service) checkArchiveNotify(msg *ArchiveUpInfo) (err error) {
	var arch *archivemodel.ArchiveCanal
	switch {
	default:
		if msg.Action == ActionInsert {
			if msg.New != nil && msg.New.State >= 0 {
				arch = msg.New
			}
			break
		}
		if msg.Action == ActionUpdate {
			if archiveStateChange(msg.New, msg.Old) {
				arch = msg.New
			}
			break
		}
	}
	if arch == nil {
		log.Warn("no need to update up cache, msg value=%v", msg)
		return
	}
	s.worker.Add(func() {
		s.upRPC.UpCount(context.Background(), &upGRPCv1.UpCountReq{
			Mid: arch.Mid,
		})
		var upCacheReq = &upGRPCv1.UpCacheReq{
			Mid: arch.Mid,
			Aid: arch.AID,
		}
		if arch.State >= 0 {
			s.upRPC.AddUpPassedCache(context.Background(), upCacheReq)
			log.Info("rpc add up cache, mid=%d, aid=%d", upCacheReq.Mid, upCacheReq.Aid)
		} else {
			s.upRPC.DelUpPassedCache(context.Background(), upCacheReq)
			log.Info("rpc delete up cache, mid=%d, aid=%d", upCacheReq.Mid, upCacheReq.Aid)
		}
	})
	return
}

func (s *Service) handleArchiveT(msg *databus.Message) (err error) {
	var m = &CanalMsg{}
	if err = json.Unmarshal(msg.Value, m); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
		return
	}

	switch m.Table {
	case TableArchiveStaff:
		var new, old archivemodel.ArchiveStaff
		switch m.Action {
		case ActionInsert:
			if err = json.Unmarshal(m.New, &new); err != nil {
				log.Error("m.New -> json.Unmarshal(%v) error(%v)", m.New, err)
				return
			}
		case ActionUpdate:
			if err = json.Unmarshal(m.New, &new); err != nil {
				log.Error("m.New -> json.Unmarshal(%v) error(%v)", m.New, err)
				return
			}
			if err = json.Unmarshal(m.Old, &old); err != nil {
				log.Error("m.Old -> json.Unmarshal(%v) error(%v)", m.New, err)
				return
			}
			if new.State == old.State {
				log.Warn("new staff state(%d) eq old staff state(%d)", new.State, old.State)
				return
			}
		}
		// state是正常说明是新增，否则是删除
		needInsert := new.State == archivemodel.StaffStateNormal
		var req = &upGRPCv1.UpCacheReq{Mid: new.StaffMid, Aid: new.Aid}
		if needInsert {
			_, err = s.upRPC.AddUpPassedCacheByStaff(context.Background(), req)
			if err != nil {
				log.Error("rpc call add up staff, new=%v, err=%v", new, err)
			} else {
				log.Info("rpc call add up staff, new=%v", new)
			}
		} else {
			_, err = s.upRPC.DelUpPassedCacheByStaff(context.Background(), req)
			if err != nil {
				log.Error("rpc call del up staff, new=%v, err=%v", new, err)
			} else {
				log.Info("rpc call del up staff, new=%v", new)
			}
		}
	}

	return
}

func archiveStateChange(a, b *archivemodel.ArchiveCanal) bool {
	if a == b {
		return false
	} else if a == nil || b == nil {
		return true
	}

	if a.State == b.State {
		return false
	}

	var min, max int
	if a.State > b.State {
		min, max = b.State, a.State
	} else {
		min, max = a.State, b.State
	}
	if min < 0 && max >= 0 {
		return true
	}
	return false
}

//WarmUp warm up
func (s *Service) WarmUp(c context.Context, req *model.WarmUpReq) (res *model.WarmUpReply, err error) {
	var (
		d                = s.crmdb.GetDb()
		lastID           = req.LastID
		limit, thisCount = 100, 100
	)
	var count = 0
	if req.Size == -1 {
		req.Size = math.MaxInt32
	}
	for ; count < req.Size && thisCount == limit; count += thisCount {
		time.Sleep(time.Millisecond * 1000)
		var end = count + limit
		if end > req.Size {
			limit = req.Size - count
		}
		var upList []*upcrmmodel.UpBaseInfo
		err = d.Select("mid, id").Where("id>?", lastID).Limit(limit).Find(&upList).Error
		if err != nil {
			log.Error("fail to query db, err=%v", err)
			return
		}

		thisCount = len(upList)
		var mids []int64
		for _, v := range upList {
			lastID = mathutil.Max(lastID, int(v.ID))
			mids = append(mids, v.Mid)
		}

		var _, e = s.upRPC.UpsArcs(context.Background(), &upGRPCv1.UpsArcsReq{
			Mids: mids,
			Pn:   1,
			Ps:   1,
		})

		if e != nil {
			log.Warn("up rpc UpsArcs return err=%v", e)
		}
		_, e = s.upRPC.UpsCount(context.Background(), &upGRPCv1.UpsCountReq{
			Mids: mids,
		})
		if e != nil {
			log.Warn("up rpc UpsCount return err=%v", e)
		}
		log.Info("warm ups, handled last id=%d, count=%d", lastID, count)
	}
	log.Info("warm ups, begin id=%d, expect up size=%d, end id=%d, real count=%d", req.LastID, req.Size, lastID, count)
	res = &model.WarmUpReply{
		LastID: lastID,
	}
	return
}

//WarmUpMid warm up by mid
func (s *Service) WarmUpMid(c context.Context, req *model.WarmUpReq) (res *model.WarmUpReply, err error) {
	if _, err = s.upRPC.UpArcs(context.Background(), &upGRPCv1.UpArcsReq{
		Mid: req.Mid,
		Pn:  1,
		Ps:  1,
	}); err != nil {
		log.Error("up rpc UpsArc(%d) return err=%v", req.Mid, err)
		return
	}
	var (
		count    int64
		cntReply *upGRPCv1.UpCountReply
	)
	if cntReply, err = s.upRPC.UpCount(context.Background(), &upGRPCv1.UpCountReq{
		Mid: req.Mid,
	}); err != nil {
		log.Error("up rpc UpsCount(%d) return err=%v", req.Mid, err)
		return
	}
	if cntReply != nil {
		count = cntReply.Count
	}
	log.Info("warm up(%d) real count=%d", req.Mid, count)
	return
}

//AddStaff .
func (s *Service) AddStaff(c context.Context, req *model.AddStaffReq) (res *upGRPCv1.NoReply, err error) {
	return s.upRPC.AddUpPassedCacheByStaff(c, &upGRPCv1.UpCacheReq{Aid: req.Aid, Mid: req.StaffMid})
}

//DeleteStaff .
func (s *Service) DeleteStaff(c context.Context, req *model.AddStaffReq) (res *upGRPCv1.NoReply, err error) {
	return s.upRPC.DelUpPassedCacheByStaff(c, &upGRPCv1.UpCacheReq{Aid: req.Aid, Mid: req.StaffMid})
}
