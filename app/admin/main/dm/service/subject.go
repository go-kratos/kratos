package service

import (
	"context"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	manager "go-common/library/queue/databus/report"
)

// ArchiveList get archive list.
func (s *Service) ArchiveList(c context.Context, req *model.ArchiveListReq) (res *model.ArchiveResult, err error) {
	var (
		oids, aids   []int64
		aidMap       = make(map[int64]struct{})
		arcMap       = make(map[int64]*model.ArcVideo)
		archiveTypes = make(map[int64]*model.ArchiveType)
		subs         map[int64]*model.Subject
		r            = &model.SearchSubjectReq{
			State: req.State,
			Attrs: req.Attrs,
			Pn:    req.Pn,
			Ps:    req.Ps,
			Sort:  req.Sort,
			Order: req.Order,
		}
	)
	res = &model.ArchiveResult{}
	res.Page = &model.Page{}
	if req.ID > 0 {
		switch req.IDType {
		case "oid":
			r.Oids = append(r.Oids, req.ID)
		case "mid":
			r.Mids = append(r.Mids, req.ID)
		case "aid":
			r.Aids = append(r.Aids, req.ID)
		case "ep", "ss":
			if r.Aids, r.Oids, err = s.dao.SeasonInfos(c, req.IDType, req.ID); err != nil {
				err = nil
				return
			}
		}
	}
	data := make([]*model.DMSubject, 0)
	if len(r.Aids) > 0 && req.Page >= 1 {
		var (
			pages []*api.Page
			arg   = archive.ArgAid2{Aid: r.Aids[0]}
		)
		if pages, err = s.arcRPC.Page3(c, &arg); err != nil {
			log.Error("arcRPC.Page3(%v) error(%v)", arg, err)
			return
		}
		if len(pages) < int(req.Page) {
			log.Error("req.Page too big(%d) error(%v)", req.Page, err)
			return
		}
		r.Oids = append(r.Oids, pages[req.Page-1].Cid)
	}
	res = new(model.ArchiveResult)
	if oids, res.Page, err = s.dao.SearchSubject(c, r); err != nil {
		return
	}
	if subs, err = s.dao.Subjects(c, model.SubTypeVideo, oids); err != nil {
		return
	}
	for _, oid := range oids {
		if sub, ok := subs[oid]; ok {
			s := &model.DMSubject{
				OID:    sub.Oid,
				Type:   sub.Type,
				AID:    sub.Pid,
				ACount: sub.ACount,
				Limit:  sub.Maxlimit,
				CTime:  sub.Ctime,
				MTime:  sub.Mtime,
				MID:    sub.Mid,
				State:  sub.State,
			}
			data = append(data, s)
		}
	}
	if len(data) <= 0 {
		return
	}
	for _, idx := range data {
		if _, ok := aidMap[idx.AID]; !ok {
			aidMap[idx.AID] = struct{}{}
			aids = append(aids, idx.AID)
		}
	}
	if arcMap, err = s.dao.ArchiveVideos(c, aids); err != nil {
		return
	}
	if archiveTypes, err = s.dao.TypeInfo(c); err != nil {
		return
	}
	for _, idx := range data {
		info, ok := arcMap[idx.AID] // get archive info
		if !ok {
			continue
		}
		idx.Title = info.Archive.Title
		idx.TID = info.Archive.TID
		if v, ok := archiveTypes[idx.TID]; ok {
			idx.TName = v.Name
		}
		if len(info.Videos) > 0 { // get ep_title name
			for _, video := range info.Videos {
				if video.CID == idx.OID {
					idx.ETitle = video.Title
				}
			}
		}
	}
	res.ArcLists = data
	return
}

// UptSubjectsState change oids subject state and send manager log.
func (s *Service) UptSubjectsState(c context.Context, tp int32, uid int64, uname string, oids []int64, state int32, comment string) (err error) {
	var affect int64
	for _, oid := range oids {
		if affect, err = s.dao.UpSubjectState(c, tp, oid, state); err != nil {
			return
		}
		if affect == 0 {
			log.Info("s.UpSubjectState affect=0 oid(%d)", oid)
			continue
		}
		managerInfo := &manager.ManagerInfo{
			UID:      uid,
			Uname:    uname,
			Business: model.DMLogBizID,
			Type:     int(tp),
			Oid:      oid,
			Ctime:    time.Now(),
			Content: map[string]interface{}{
				"comment": comment,
			},
		}
		if state == model.SubStateOpen {
			managerInfo.Action = "开启弹幕池"
		} else {
			managerInfo.Action = "关闭弹幕池"
		}
		manager.Manager(managerInfo)
		log.Info("s.managerLogSend(%+v)", managerInfo)
	}
	return
}

// UpSubjectMaxLimit update maxlimit in dm subject.
func (s *Service) UpSubjectMaxLimit(c context.Context, tp int32, oid, maxlimit int64) (err error) {
	_, err = s.dao.UpSubjectMaxlimit(c, tp, oid, maxlimit)
	return
}

// SubjectLog get subject log
func (s *Service) SubjectLog(c context.Context, tp int32, oid int64) (data []*model.SubjectLog, err error) {
	data, err = s.dao.SearchSubjectLog(c, tp, oid)
	return
}
