package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	credit "go-common/app/job/main/credit/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	xtime "go-common/library/time"
)

func (s *Service) afterSetGrpResult(grp *param.GroupResParam, g *model.Group, tinyChalls map[int64]*model.TinyChall) {
	logging := func() {
		// write group log
		l := &model.WLog{
			AdminID:  grp.AdminID,
			Admin:    grp.AdminName,
			Oid:      g.Oid,
			Business: g.Business,
			Target:   g.ID,
			Module:   model.WLogModuleGroup,
			Remark:   fmt.Sprintf(`“工单编号 %d” %s 理由：%s`, g.ID, s.StateDescr(grp.Business, 0, grp.State), grp.Reason),
			Note:     grp.Reason,
		}
		s.writeAuditLog(l)

		// write challenge log
		for _, c := range tinyChalls {
			l := &model.WLog{
				AdminID:  grp.AdminID,
				Admin:    grp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   c.Cid,
				Module:   model.WLogModuleChallenge,
				Remark:   fmt.Sprintf(`“工单详情编号 %d” %s 理由：%s`, c.Cid, s.StateDescr(grp.Business, 0, grp.State), grp.Reason),
				Note:     grp.Reason,
			}
			s.writeAuditLog(l)
		}
	}

	callback := func() {
		metas := s.dao.AllMetas(context.Background())
		meta, ok := metas[grp.Business]
		if !ok || meta.ItemType != "group" {
			return
		}

		// handle callback
		cb, ok := s.callbackCache[grp.Business]
		if !ok || cb == nil {
			return
		}

		mids := make([]int64, 0, len(tinyChalls))
		for _, c := range tinyChalls {
			mids = append(mids, c.Mid)
		}

		// origin callback
		if cb.URL != "" && cb.State == model.CallbackEnable {
			payload := &model.Payload{
				Verb: model.GroupSetResult,
				Actor: model.Actor{
					AdminID: grp.AdminID,
				},
				CTime:  xtime.Time(time.Now().Unix()),
				Object: grp,
				Target: g,
				Influence: map[string]interface{}{
					"mids": mids,
				},
			}
			if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
				log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
				// continue
			}
		}
	}

	extra := func() {
		// extra condition
		switch grp.Business {
		case model.ArchiveComplain:
			s.notifyUsers(grp.Business, grp.Oid, grp.AdminID, tinyChalls, grp.IsMessage)
		default:
		}
	}

	logging()
	callback()
	extra()
}

func (s *Service) afterBatchSetGrpResult(bgrp *param.BatchGroupResParam, groups map[int64]*model.Group, tinyChalls map[int64]*model.TinyChall) {
	// record log
	logging := func() {
		for _, g := range groups {
			l := &model.WLog{
				AdminID:  bgrp.AdminID,
				Admin:    bgrp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   g.ID,
				Module:   model.WLogModuleGroup,
				Remark:   fmt.Sprintf(`“工单编号 %d” %s 理由：%s`, g.ID, s.StateDescr(bgrp.Business, 0, bgrp.State), bgrp.Reason),
				Note:     bgrp.Reason,
			}
			s.writeAuditLog(l)
		}
		for _, c := range tinyChalls {
			g, ok := groups[c.Gid]
			if !ok {
				log.Warn("Failed to retrive group by gid %d", c.Gid)
				continue
			}

			l := &model.WLog{
				AdminID:  bgrp.AdminID,
				Admin:    bgrp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   c.Cid,
				Module:   model.WLogModuleChallenge,
				Remark:   fmt.Sprintf(`“工单详情编号 %d” 状态从 %s 修改为 %s 理由：%s`, c.Cid, s.StateDescr(bgrp.Business, 0, c.State), s.StateDescr(bgrp.Business, 0, bgrp.State), bgrp.Reason),
				Note:     bgrp.Reason,
			}
			s.writeAuditLog(l)
		}
	}

	// origin callback
	callback := func() {
		metas := s.dao.AllMetas(context.Background())
		meta, ok := metas[bgrp.Business]
		if !ok || meta.ItemType != "group" {
			return
		}
		cb, ok := s.callbackCache[bgrp.Business]
		if !ok || cb == nil {
			return
		}

		toids := make([]int64, 0, len(groups))
		for _, g := range groups {
			toids = append(toids, g.Oid)
		}
		bgrp.Oids = toids

		gmids := make(map[int64][]int64, len(groups))
		for _, c := range tinyChalls {
			if _, ok := gmids[c.Gid]; !ok {
				gmids[c.Gid] = make([]int64, 0)
			}
			gmids[c.Gid] = append(gmids[c.Gid], c.Mid)
		}

		if cb.URL != "" && cb.State == model.CallbackEnable {
			for _, g := range groups {
				payload := &model.Payload{
					Verb: model.BatchGroupSetResult,
					Actor: model.Actor{
						AdminID: bgrp.AdminID,
					},
					CTime:  xtime.Time(time.Now().Unix()),
					Object: bgrp,
					Target: g,
					Influence: map[string]interface{}{
						"mids": gmids[g.ID],
					},
				}
				if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
					log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
					// continue
				}
			}
		}
	}

	// extra condition
	extra := func() {
		switch bgrp.Business {
		case model.ArchiveComplain:
			oidToChalls := make(map[int64]map[int64]*model.TinyChall, len(groups))
			for _, c := range tinyChalls {
				g, ok := groups[c.Gid]
				if !ok {
					continue
				}
				if _, ok := oidToChalls[g.Oid]; !ok {
					oidToChalls[g.Oid] = make(map[int64]*model.TinyChall)
				}
				oidToChalls[g.Oid][c.Cid] = c
			}
			for oid, cs := range oidToChalls {
				s.notifyUsers(bgrp.Business, oid, bgrp.AdminID, cs, bgrp.IsMessage)
			}
		default:
		}
	}

	logging()
	callback()
	extra()
}

func (s *Service) afterSetChallResult(crp *param.ChallResParam, c *model.Chall) {
	logging := func() {
		// write challenge log
		l := &model.WLog{
			AdminID:  crp.AdminID,
			Admin:    crp.AdminName,
			Oid:      c.Oid,
			Business: c.Business,
			Target:   c.Cid,
			Module:   model.WLogModuleChallenge,
			Remark:   fmt.Sprintf(`“工单详情编号 %d” 处理状态从 %s 修改为 %s 理由：%s`, c.Cid, s.StateDescr(c.Business, 0, c.State), s.StateDescr(c.Business, 0, crp.State), crp.Reason),
			Note:     crp.Reason,
		}
		s.writeAuditLog(l)
	}

	callback := func() {
		metas := s.dao.AllMetas(context.Background())
		meta, ok := metas[c.Business]
		if !ok || meta.ItemType != "challenge" {
			log.Warn("item type not exist or is not challenge, meta(%+v)", meta)
			return
		}

		// handle callback
		cb, ok := s.callbackCache[c.Business]
		if !ok || cb == nil {
			log.Error("can not find callback func, business(%d)", c.Business)
			return
		}

		mids := []int64{c.Mid}
		if cb.URL != "" && cb.State == model.CallbackEnable {
			payload := &model.Payload{
				Verb: model.ChallSetResult,
				Actor: model.Actor{
					AdminID: crp.AdminID,
				},
				CTime:  xtime.Time(time.Now().Unix()),
				Object: crp,
				Target: c,
				Influence: map[string]interface{}{
					"mids": mids,
				},
			}
			if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
				log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
				// continue
			}
		}
	}

	extra := func() {
		// extra condition
	}

	logging()
	callback()
	extra()
}

func (s *Service) afterBatchSetChallResult(bcrp *param.BatchChallResParam, challs map[int64]*model.Chall) {
	// Record Log
	logging := func() {
		for _, c := range challs {
			l := &model.WLog{
				AdminID:  bcrp.AdminID,
				Admin:    bcrp.AdminName,
				Oid:      c.Oid,
				Business: c.Business,
				Target:   c.Cid,
				Module:   model.WLogModuleChallenge,
				Remark:   fmt.Sprintf(`“工单详情编号 %d” 状态从 %s 修改为 %s 理由：%s`, c.Cid, s.StateDescr(c.Business, 0, c.State), s.StateDescr(c.Business, 0, bcrp.State), bcrp.Reason),
				Note:     bcrp.Reason,
			}
			s.writeAuditLog(l)
		}
	}

	// origin callback
	callback := func() {
		metas := s.dao.AllMetas(context.Background())
		for _, c := range challs {
			meta, ok := metas[c.Business]
			if !ok || meta.ItemType != "challenge" {
				return
			}

			cb, ok := s.callbackCache[c.Business]
			if !ok || cb == nil {
				return
			}

			mids := []int64{c.Mid}
			// origin callback
			if cb.URL != "" && cb.State == model.CallbackEnable {
				payload := &model.Payload{
					Verb: model.BatchChallSetResult,
					Actor: model.Actor{
						AdminID: bcrp.AdminID,
					},
					CTime:  xtime.Time(time.Now().Unix()),
					Object: bcrp,
					Target: c,
					Influence: map[string]interface{}{
						"mids": mids,
					},
				}
				if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
					log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
					// continue
				}
			}
		}
	}

	extra := func() {
		// extra condition
	}
	logging()
	callback()
	extra()
}

func (s *Service) afterSetGroupRole(grsp *param.GroupRoleSetParam, groups map[int64]*model.Group) {
	// Record Log
	logging := func() {
		logs, err := s.LastLogStat(context.Background(), grsp.GID, []int{model.WLogModuleRoleShift}, []string{"int_1", "ctime", "int_2", "uid"})
		if err != nil {
			log.Error("s.LastLogStat() failed error:%v", err)
			return
		}
		for _, g := range groups {
			var (
				preTag *model.TagMeta
				newTag *model.TagMeta
				err    error
			)
			if preTag, err = s.tag(g.Business, g.Tid); err != nil {
				log.Warn("s.tag(%d,%d) error(%v)", g.Business, g.Tid, err)
				err = nil
				preTag.Bid = g.Business
				preTag.Tid = g.Tid
				preTag.RID = g.Rid
			}
			if newTag, err = s.tag(g.Business, grsp.TID); err != nil {
				log.Warn("s.tag(%d,%d) error(%v)", g.Business, g.Tid, err)
				err = nil
				newTag.Bid = g.Business
				newTag.Tid = g.Tid
			}

			l := &model.WLog{
				AdminID:  grsp.AdminID,
				Admin:    grsp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   g.ID,
				Module:   model.WLogModuleRoleShift,
				Remark: fmt.Sprintf(`“工单编号 %d” 业务角色流转从 %d:%s 到 %d:%s, 管理 tag 从 %d:%s 变更到 %d:%s, 备注: %s`,
					g.ID, g.Rid, s.RidDesc(g.Business, g.Rid), grsp.RID, s.RidDesc(g.Business, grsp.RID), g.Tid, preTag.Name, grsp.TID, newTag.Name, grsp.Note),
				Note:   grsp.Note,
				OpType: strconv.Itoa(int(model.RoleShift)),
				PreRid: strconv.Itoa(int(preTag.RID)),
			}
			if roleLog, ok := logs[g.ID]; !ok {
				l.TimeConsume = int64(time.Since(g.LastTime.Time()) / time.Second)
			} else {
				lastRoleTime, _ := time.ParseInLocation("2006-01-02 15:04:05", roleLog.CTime, time.Local)
				if lastRoleTime.After(g.LastTime.Time()) {
					l.TimeConsume = int64(time.Since(lastRoleTime) / time.Second)
				} else {
					l.TimeConsume = int64(time.Since(g.LastTime.Time()) / time.Second)
				}
			}
			s.writeAuditLog(l)
		}
	}
	logging()
}

// set state v3
func (s *Service) afterSetGroupState(gssp *param.GroupStateSetParam, groups map[int64]*model.Group, tinyChalls map[int64]*model.TinyChall) {
	// record log
	logging := func() {
		logs, err := s.LastLogStat(context.Background(), gssp.ID, []int{model.WLogModuleRoleShift}, []string{"int_1", "ctime", "int_2", "uid"})
		if err != nil {
			log.Error("s.LastLogStat() failed error:%v", err)
			return
		}
		for _, g := range groups {
			l := &model.WLog{
				AdminID:  gssp.AdminID,
				Admin:    gssp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   g.ID,
				Module:   model.WLogModuleGroup,
				Remark:   fmt.Sprintf(`“工单编号 %d” %s 理由：%s`, g.ID, s.StateDescr(g.Business, 0, gssp.State), gssp.Reason),
				Note:     gssp.Reason,
				OpType:   strconv.Itoa(int(gssp.State)),
				PreRid:   strconv.Itoa(int(gssp.Rid)),
				Param:    gssp,
			}
			if roleLog, ok := logs[g.ID]; !ok {
				l.TimeConsume = int64(time.Since(g.LastTime.Time()) / time.Second)
			} else {
				lastRoleTime, _ := time.ParseInLocation("2006-01-02 15:04:05", roleLog.CTime, time.Local)
				if lastRoleTime.After(g.LastTime.Time()) {
					l.TimeConsume = int64(time.Since(lastRoleTime) / time.Second)
				} else {
					l.TimeConsume = int64(time.Since(g.LastTime.Time()) / time.Second)
				}
			}
			if g.BusinessObject != nil {
				l.Mid = g.BusinessObject.Mid
			}
			s.writeAuditLog(l)
		}
		for _, c := range tinyChalls {
			g, ok := groups[c.Gid]
			if !ok {
				log.Warn("Failed to retrive group by gid %d", c.Gid)
				continue
			}

			l := &model.WLog{
				AdminID:  gssp.AdminID,
				Admin:    gssp.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   c.Cid,
				Module:   model.WLogModuleChallenge,
				Remark:   fmt.Sprintf(`“工单详情编号 %d” 状态从 %s 修改为 %s 理由：%s`, c.Cid, s.StateDescr(g.Business, 0, c.State), s.StateDescr(g.Business, 0, gssp.State), gssp.Reason),
				Note:     gssp.Reason,
				Mid:      c.Mid,
			}
			s.writeAuditLog(l)
		}
	}

	// origin callback
	callback := func() {
		cb, ok := s.callbackCache[gssp.Business]
		if !ok || cb == nil {
			return
		}

		gmids := make(map[int64][]int64, len(groups))
		for _, c := range tinyChalls {
			if _, ok := gmids[c.Gid]; !ok {
				gmids[c.Gid] = make([]int64, 0)
			}
			gmids[c.Gid] = append(gmids[c.Gid], c.Mid)
		}

		if cb.URL != "" && cb.State == model.CallbackEnable {
			//common
			payload := &model.Payload{
				Bid:  int(gssp.Business),
				Verb: model.GroupSetState,
				Actor: model.Actor{
					AdminID:   gssp.AdminID,
					AdminName: gssp.AdminName,
				},
				CTime:  xtime.Time(time.Now().Unix()),
				Object: gssp,
			}

			// influence
			// 长短评
			if gssp.Business == model.ReviewShortComplain || gssp.Business == model.ReviewLongComplain {
				var reportMids = make([]int64, len(tinyChalls))
				for _, tc := range tinyChalls {
					reportMids = append(reportMids, tc.Mid)
				}
				payload.Influence = map[string]interface{}{
					"mids": reportMids,
				}
			}

			// target
			for _, g := range groups {
				payload.Targets = append(payload.Targets, g)
			}
			//todo: extra

			if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
				log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
				// continue
			}
		}
	}

	// notify user
	notify := func() {
		s.notifyUsersV4(gssp, groups, tinyChalls)
	}

	// 扣节操
	decreaseMoral := func() {
		if gssp.DecreaseMoral == 0 {
			return
		}
		var (
			gids []int64
			mids []int64
		)

		for _, g := range groups {
			gids = append(gids, g.ID)
		}
		// 被举报人 mid
		var (
			bus map[int64]*model.Business
			err error
		)
		if bus, err = s.dao.BusObjectByGids(context.Background(), gids); err != nil {
			log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
			return
		}
		for _, b := range bus {
			mids = append(mids, b.Mid)
		}

		if err := s.dao.AddMoral(context.Background(), mids, gssp); err != nil {
			log.Error("s.dao.AddMoral(%v,%d) error(%v)", mids, gssp.DecreaseMoral, err)
			return
		}
		l := &model.WLog{
			AdminID: gssp.AdminID,
			Admin:   gssp.AdminName,
			Module:  model.WLogModuleAddMoral,
			Mids:    mids,
			Param:   gssp,
		}
		s.writeAuditLog(l)
	}

	// 封禁账号
	block := func() {
		if gssp.BlockDay == 0 {
			return
		}
		var (
			gids []int64
			mids []int64
		)
		for _, g := range groups {
			gids = append(gids, g.ID)
		}
		// 被举报人 mid
		var (
			bus map[int64]*model.Business
			err error
		)
		if bus, err = s.dao.BusObjectByGids(context.Background(), gids); err != nil {
			log.Error("s.dao.BusObjectByGids(%v) error(%v)", gids, err)
			return
		}
		for _, b := range bus {
			mids = append(mids, b.Mid)
		}

		if err := s.dao.AddBlock(context.Background(), mids, gssp); err != nil {
			log.Error("s.dao.AddBlock(%d,%v) error(%v)", mids, gssp, err)
			return
		}

		if err := s.dao.AddCreditBlockInfo(context.Background(), bus, gssp); err != nil {
			log.Error("s.dao.AddCreditBlockInfo(%d,%v) error(%v)", mids, gssp, err)
			return
		}
		l := &model.WLog{
			AdminID: gssp.AdminID,
			Admin:   gssp.AdminName,
			Module:  model.WLogModuleBlock,
			Mids:    mids,
			Param:   gssp,
		}
		s.writeAuditLog(l)
	}

	logging()
	callback()
	notify()
	decreaseMoral()
	block()
}

// only set group state
func (s *Service) afterSimpleSetState(gspr *param.GroupStatePublicReferee, groups map[int64]*model.Group) {
	// record log
	logging := func() {
		for _, g := range groups {
			l := &model.WLog{
				AdminID:  gspr.AdminID,
				Admin:    gspr.AdminName,
				Oid:      g.Oid,
				Business: g.Business,
				Target:   g.ID,
				Module:   model.WLogModulePublicReferee,
				Remark:   fmt.Sprintf(`“工单编号 %d” %s 理由：%s`, g.ID, s.StateDescr(g.Business, 0, gspr.State), ""),
				OpType:   strconv.Itoa(int(gspr.State)),
			}
			s.writeAuditLog(l)
		}
	}

	// callback
	callback := func() {
		cb, ok := s.callbackCache[gspr.Business]
		if !ok || cb == nil {
			return
		}
		gssp := &param.GroupStateSetParam{
			ID:        gspr.ID,
			Business:  gspr.Business,
			AdminID:   gspr.AdminID,
			AdminName: gspr.AdminName,
			State:     gspr.State,
		}
		if cb.URL != "" && cb.State == model.CallbackEnable {
			//common
			payload := &model.Payload{
				Bid:  int(gssp.Business),
				Verb: model.GroupSetPublicReferee,
				Actor: model.Actor{
					AdminID:   gssp.AdminID,
					AdminName: gssp.AdminName,
				},
				CTime:  xtime.Time(time.Now().Unix()),
				Object: gssp,
			}
			// target
			for _, g := range groups {
				payload.Targets = append(payload.Targets, g)
			}

			if err := s.SendCallbackRetry(context.Background(), cb, payload); err != nil {
				log.Error("Failed to s.SendCallbackRetry(%+v, %+v): %v", cb, payload, err)
				// continue
			}
		}
	}

	logging()
	callback()
}

func (s *Service) afterSetBusinessState(challs []*model.Chall) {
	// todo: after set business state
	logging := func() {
		for _, c := range challs {
			l := &model.WLog{
				AdminID:  c.AssigneeAdminID,
				Admin:    c.AssigneeAdminName,
				Oid:      c.Oid,
				Business: c.Business,
				Target:   c.Cid,
				Module:   model.WLogModuleReply,
				Mid:      c.Mid,
				TypeID:   c.TypeID,
			}
			s.writeAuditLog(l)
		}
	}
	logging()
}

func (s *Service) afterAddReply(ep *param.EventParam, c *model.Chall) {
	logging := func() {
		l := &model.WLog{
			AdminID:  ep.AdminID,
			Admin:    ep.AdminName,
			Oid:      c.Oid,
			Business: c.Business,
			Target:   c.Cid,
			Module:   model.FeedBackTypeReply,
			Mid:      c.Mid,
			Remark:   ep.Content,
			Note:     strconv.Itoa(int(ep.Event)),
			Meta: map[string]interface{}{
				"ep":        ep,
				"challenge": c,
			},
		}
		s.writeReplyLog(l)
	}
	logging()
}

func (s *Service) afterAddMultiReply(bep *param.BatchEventParam, cs map[int64]*model.Chall) {
	loggin := func() {
		for _, cid := range bep.Cids {
			var (
				c  *model.Chall
				ok bool
			)
			if c, ok = cs[cid]; !ok {
				continue
			}
			l := &model.WLog{
				AdminID:  bep.AdminID,
				Admin:    bep.AdminName,
				Oid:      c.Oid,
				Business: c.Business,
				Target:   c.Cid,
				Module:   model.FeedBackTypeReply,
				Mid:      c.Mid,
				Remark:   bep.Content,
				Note:     strconv.Itoa(int(bep.Event)),
				Meta: map[string]interface{}{
					"bep":       bep,
					"challenge": cs,
				},
			}
			s.writeReplyLog(l)
		}
	}

	loggin()
}

func (s *Service) notifyUsers(business int8, oid int64, adminid int64, challs map[int64]*model.TinyChall, isDisposeMsg bool) (err error) {
	var makeDealMsgParam = model.DealArcComplainMsg
	var makeReceivedMsgParam = model.ReceivedArcComplainMsg

	mids := make([]int64, 0, len(challs))
	for _, c := range challs {
		mids = append(mids, int64(c.Mid))
	}

	if business == int8(model.ArchiveComplain) {
		// if challenge ctime < 10min, report appeal received
		var rChall = make([]*model.TinyChall, 0)
		var rMids = make([]int64, 0)
		for _, c := range challs {
			if time.Since(c.CTime.Time()) < time.Minute*10 {
				//todo: rmids send message archive complain received
				rChall = append(rChall, c)
				rMids = append(rMids, c.Mid)
			}
		}

		if len(rChall) > 0 {
			rmp := makeReceivedMsgParam(oid, rMids)
			if err = s.dao.SendMessage(context.Background(), rmp); err != nil {
				log.Error("Failed to s.dao.SendMessage(%v, %v, %+v): %v", oid, mids, rmp, err)
				err = nil
			}
			log.Info("send archive complain received message business(%d) oid(%d) mids(%v) mc(%s) message(%+v)", business, oid, mids, rmp.MC, rmp)

			// report archive complain received msg
			for _, tc := range rChall {
				report.Manager(&report.ManagerInfo{
					Uname:    "",
					UID:      0,
					Business: 11,
					Type:     2,
					Oid:      oid,
					Action:   "notify_users_received",
					Ctime:    time.Now(),
					Index:    []interface{}{tc.Cid, tc.Gid, tc.Mid},
					Content: map[string]interface{}{
						"mid":     tc.Mid,
						"message": rmp,
					},
				})
			}
		}
	}
	if !isDisposeMsg {
		return
	}
	time.Sleep(500 * time.Millisecond)
	mdmp := makeDealMsgParam(oid, mids)
	if err = s.dao.SendMessage(context.Background(), mdmp); err != nil {
		log.Error("Failed to s.dao.SendMessage(%v, %v, %+v): %v", oid, mids, mdmp, err)
		return
	}

	// group dispose report
	report.Manager(&report.ManagerInfo{
		Uname:    "",
		UID:      adminid,
		Business: 11,
		Type:     3,
		Oid:      oid,
		Action:   "notify_users_dispose",
		Ctime:    time.Now(),
		Index:    []interface{}{business},
		Content: map[string]interface{}{
			"mid":     mids,
			"message": mdmp,
		},
	})

	log.Info("send message business(%d) oid(%d) mids(%v) mc(%s) message(%+v)", business, oid, mids, mdmp.MC, mdmp)
	return
}

func (s *Service) notifyUsersV4(gssp *param.GroupStateSetParam, groups map[int64]*model.Group, tinyChalls map[int64]*model.TinyChall) (err error) {
	// 稿件投诉补发已收到消息
	if gssp.Business == model.ArchiveComplain {
		for _, g := range groups {
			// if challenge ctime < 10min, report appeal received
			var rMids = make([]int64, 0)
			for _, c := range tinyChalls {
				if time.Since(c.CTime.Time()) < time.Minute*10 && c.Gid == g.ID {
					//todo: rmids send message archive complain received
					rMids = append(rMids, c.Mid)
				}
			}
			if len(rMids) > 0 {
				rmp := &param.MessageParam{
					Type:     "json",
					Source:   1,
					DataType: 4,
					MC:       model.ArcComplainRevMC,
					Title:    "您的投诉已收到",
					Context:  fmt.Sprintf("您对稿件（av%d）的举报我们已经收到。感谢您对 bilibili 社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~", g.Oid),
					MidList:  rMids,
				}
				if err = s.dao.SendMessage(context.Background(), rmp); err != nil {
					log.Error("Failed to s.dao.SendMessage(%+v) mid(%v): %v", rmp, rMids, err)
					err = nil
				}
				log.Info("send archive complain received message gid(%d) oid(%d) mids(%v) message(%+v)", g.ID, g.Oid, rMids, rmp)
				// report archive complain received msg
				report.Manager(&report.ManagerInfo{
					Uname:    "",
					UID:      0,
					Business: 11,
					Type:     model.FeedBackTypeNotifyUserReceived,
					Oid:      g.Oid,
					Action:   "notify_users_received_v4",
					Ctime:    time.Now(),
					Index:    []interface{}{g.ID},
					Content: map[string]interface{}{
						"message": rmp,
					},
				})
			}
		}
	}
	time.Sleep(500 * time.Millisecond)

	// 通知举报人
	if gssp.IsMessage {
		mps := make(map[int64]*param.MessageParam) // map[oid]*mp or map[eid]*mp
		// all mids need to report
		cmids := make(map[int64][]int64) //map[gid][]mids
		for _, c := range tinyChalls {
			if _, ok := cmids[c.Gid]; !ok {
				cmids[c.Gid] = make([]int64, 0)
			}
			cmids[c.Gid] = append(cmids[c.Gid], c.Mid)
		}

		switch gssp.Business {
		case model.ArchiveComplain:
			for _, g := range groups {
				reportMids, ok := cmids[g.ID]
				if !ok {
					log.Warn("report mid not found gid(%d)", g.ID)
					continue
				}
				mp := &param.MessageParam{
					Type:     "json",
					Source:   1,
					DataType: 4,
					MC:       model.ArcComplainDealMC,
					Title:    "您的投诉已被受理",
					Context:  fmt.Sprintf("您对稿件（av%d）的投诉已被受理。感谢您对 bilibili 社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~ ", g.Oid),
					MidList:  reportMids,
				}
				mps[g.ID] = mp
			}
		case model.ChannelComplain:
			for _, g := range groups {
				reportMids, ok := cmids[g.ID]
				if !ok {
					log.Warn("report mid not found gid(%d)", g.ID)
					continue
				}
				mp := &param.MessageParam{
					Type:     "json",
					Source:   1,
					DataType: 4,
					MC:       model.WkfNotifyMC,
					Title:    "你的举报已成功处理",
					Context: fmt.Sprintf(`您在稿件【#{av%d}{"https://www.bilibili.com/video/av%d"}】举报的频道【%s】已处理，感谢反馈，点击进去查看`,
						g.Oid, g.Oid, g.BusinessObject.Title),
					MidList: reportMids,
				}
				mps[g.ID] = mp
			}
		case model.CommentComplain: // 评论举报
			var blockStr, isDel string
			if gssp.BlockDay == -1 {
				blockStr = "该用户已被永久封禁。"
			} else if gssp.BlockDay > 0 {
				blockStr = fmt.Sprintf("并被封禁%d天。", gssp.BlockDay)
			}
			if gssp.DisposeMode == 1 {
				isDel = "已被移除"
			} else {
				isDel = "已被处罚"
			}
			tMeta, _ := s.tag(gssp.Business, gssp.Tid)
			tName := tMeta.Name // 举报理由
			for _, g := range groups {
				if g.Fid == model.ReplyFidManga {
					continue
				}
				var extMap map[string]interface{}
				reportMids, ok := cmids[g.ID]
				if !ok {
					log.Warn("report mid not found gid(%d)", g.ID)
					continue
				}
				json.Unmarshal([]byte(g.BusinessObject.Extra), &extMap)
				extTitle, _ := extMap["title"].(string)
				extTitle = subString(extTitle, 0, 40)
				extLink, _ := extMap["link"].(string)
				content := fmt.Sprintf(`您好，您在#{%s}{%s}下举报的评论『%s』%s，%s 理由：%s。`+model.NotifyComRulesReport,
					extTitle, extLink, pretreat(g.BusinessObject.Title), isDel, blockStr, tName)
				mp := &param.MessageParam{
					Type:     "json",
					Source:   1,
					DataType: 4,
					MC:       model.WkfNotifyMC,
					Title:    "举报处理结果通知",
					Context:  content,
					MidList:  reportMids,
				}
				mps[g.ID] = mp
			}
		}
		for gid, mp := range mps {
			if err = s.dao.SendMessage(context.Background(), mp); err != nil {
				log.Error("Failed to s.dao.SendMessage(%+v) mids(%v): %v", mp, err)
				continue
			}
			// group dispose report
			report.Manager(&report.ManagerInfo{
				Uname:    gssp.AdminName,
				UID:      gssp.AdminID,
				Business: 11,
				Type:     model.FeedBackTypeNotifyUserDisposed,
				Action:   "notify_users_v4_reporters",
				Ctime:    time.Now(),
				Index:    []interface{}{gssp.Business},
				Content: map[string]interface{}{
					"message": mp,
				},
			})
			log.Info("notifyUsersV4 success send to reporters message(%+v) gid(%d)", mp, gid)
		}
	}

	// 通知被举报人(up主)
	if gssp.IsMessageUper {
		mps := make(map[int64]*param.MessageParam) // map[oid]*mp or map[eid]*mp
		switch gssp.Business {
		case model.CommentComplain: //评论举报
			var (
				blockStr   string
				disposeStr string
			)
			if gssp.BlockDay == -1 {
				blockStr = "本帐号已被永久封禁。"
			} else if gssp.BlockDay > 0 {
				blockStr = fmt.Sprintf("并被封禁%d天。", gssp.BlockDay)
			}
			if gssp.DisposeMode == 1 {
				disposeStr = "已被举报并移除"
			} else {
				disposeStr = "已被举报并处罚"
			}
			tMeta, _ := s.tag(gssp.Business, gssp.Tid)
			tName := tMeta.Name // 举报理由
			for _, g := range groups {
				if g.Fid == model.ReplyFidManga {
					continue
				}
				var extMap map[string]interface{}
				json.Unmarshal([]byte(g.BusinessObject.Extra), &extMap)
				extTitle, _ := extMap["title"].(string)
				extTitle = subString(extTitle, 0, 40)
				extLink, _ := extMap["link"].(string)
				content := fmt.Sprintf(`您好，您在#{%s}{%s}下发布的评论『%s』%s，%s 理由：%s。`, extTitle, extLink, pretreat(g.BusinessObject.Title), disposeStr, blockStr, tName)
				content = suffixContent(content, gssp)
				mp := &param.MessageParam{
					Type:     "json",
					Source:   1,
					DataType: 4,
					MC:       model.WkfNotifyMC,
					Title:    "评论违规处理通知",
					Context:  content,
					MidList:  []int64{g.BusinessObject.Mid},
				}
				mps[g.ID] = mp
			}
		}
		for gid, mp := range mps {
			if err = s.dao.SendMessage(context.Background(), mp); err != nil {
				log.Error("Failed to s.dao.SendMessage(%+v) gid(%d): %v", mp, gid, err)
				continue
			}
			// group dispose report
			report.Manager(&report.ManagerInfo{
				Uname:    gssp.AdminName,
				UID:      gssp.AdminID,
				Business: 11,
				Type:     3,
				Action:   "notify_users_v4_upers",
				Ctime:    time.Now(),
				Index:    []interface{}{gssp.Business},
				Content: map[string]interface{}{
					"mid":     mp.MidList,
					"message": *mp,
				},
			})
			log.Info("notifyUsersV4 success send to uper message(%+v) gid(%d)", mp, gid)
		}
	}
	return
}

// 站内信评论内容预处理
func pretreat(str string) string {
	str = subString(str, 0, 40)
	str = filterViolationMsg(str)
	return str
}

// 字符串截断
func subString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// 字符串加*
func filterViolationMsg(msg string) string {
	s := []rune(msg)
	for i := 0; i < len(s); i++ {
		if i%3 != 0 {
			s[i] = '*'
		}
	}
	return string(s)
}

// 评论站内信后缀(被举报人)
func suffixContent(content string, gssp *param.GroupStateSetParam) string {
	switch gssp.BlockReason {
	// 剧透, 广告, 抢楼, 刷屏
	case credit.ReasonSpoiler, credit.ReasonGarbageAds, credit.ReasonGrabFloor, credit.ReasonBrushScreen:
		content += model.NotifyComRules
	// 引战, 人身攻击
	case credit.ReasonLeadBattle, credit.ReasonPersonalAttacks:
		content += model.NotifyComProvoke
	default:
		content += model.NofityComProhibited
	}
	return content
}
