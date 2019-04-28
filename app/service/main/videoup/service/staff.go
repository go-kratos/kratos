package service

import (
	"context"

	"encoding/json"
	"fmt"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"strings"
)

//Apply .
func (s *Service) Apply(c context.Context, ID int64) (data *archive.StaffApply, err error) {
	if data, err = s.arc.Apply(c, ID); err != nil {
		log.Error(" Apply(%d) is ,err(%+v)", ID, err)
		return
	}
	return
}

//MidCount count .
func (s *Service) MidCount(c context.Context, ID int64) (data int64, err error) {
	if data, err = s.arc.MidCount(c, ID); err != nil {
		log.Error(" MidCount(%d) is ,err(%+v)", ID, err)
		return
	}
	return
}

//Applys .
func (s *Service) Applys(c context.Context, ids []int64) (data []*archive.StaffApply, err error) {
	if data, err = s.arc.Applys(c, ids); err != nil {
		log.Error(" Applys(%d) is ,err(%+v)", ids, err)
		return
	}
	return
}

//FilterApplys .
func (s *Service) FilterApplys(c context.Context, aids []int64, mid int64) (data []*archive.StaffApply, err error) {
	if data, err = s.arc.FilterApplys(c, aids, mid); err != nil {
		log.Error(" FilterApplys(%v,%d) is ,err(%+v)", aids, mid, err)
		return
	}
	return
}

//ApplysByAID .
func (s *Service) ApplysByAID(c context.Context, aid int64) (ret []*archive.StaffApply, err error) {
	var data []*archive.StaffApply
	if data, err = s.arc.ApplysByAID(c, aid); err != nil {
		log.Error(" ApplysByAID(%d) is ,err(%+v)", aid, err)
		return
	}
	ret = make([]*archive.StaffApply, 0)
	for _, v := range data {
		if !s.HiddenApply(v) {
			ret = append(ret, v)
		}
	}
	return
}

//Staffs .
func (s *Service) Staffs(c context.Context, AID int64) (data []*archive.Staff, err error) {
	if data, err = s.arc.Staffs(c, AID); err != nil {
		log.Error(" Staffs(%d) is ,err(%+v)", AID, err)
		return
	}
	if len(data) == 0 {
		return
	}
	var staffMids []int64
	for _, staff := range data {
		staffMids = append(staffMids, staff.StaffMID)
	}
	users, err := s.Infos(c, staffMids)
	if err != nil {
		log.Error("s.Infos() err(%v)", err)
		return
	}
	for _, staff := range data {
		if _, ok := users[staff.StaffMID]; ok {
			staff.StaffName = users[staff.StaffMID].Name
		}
	}
	return
}

// CheckStaff check
func (s *Service) CheckStaff(vps []*archive.StaffParam) bool {
	if len(vps) == 0 {
		return true
	}
	for _, vp := range vps {
		if vp.MID == 0 || vp.Title == "" {
			return false
		}
	}
	return true
}

//HandleArchiveApplys   up 在稿件编辑页批量操作
func (s *Service) HandleArchiveApplys(c context.Context, aid int64, params []*archive.StaffParam, source string, syncAttr bool) (err error) {
	defer func() {
		//提供管理员触发属性位计算逻辑  该逻辑影响 属性位，合作角标,动态推送，逸飞同步
		if syncAttr {
			s.SyncStaffAttr(aid)
		}
	}()

	var g = &errgroup.Group{}
	switch source {
	case "add":
		//新稿件
		for _, v := range params {
			apply := &archive.ApplyParam{ApplyAID: aid}
			apply.ApplyStaffMID = v.MID
			apply.ApplyTitle = v.Title
			apply.ApplyTitleID = v.TitleID
			apply.Type = archive.TYPEUPADD
			apply.State = archive.APPLYSTATEOPEN
			//新增保持顺序
			if _, err = s.AddApply(c, apply, "HandleArchiveApplys/add"); err != nil {
				log.Error(" HandleArchiveApplys s.AddApply(%d) err(%+v)", aid, err)
				err = nil
			}
		}
		return nil
	case "edit":
		//稿件编辑
		var applys []*archive.StaffApply
		if applys, err = s.ApplysByAID(c, aid); err != nil || applys == nil {
			log.Error(" s.ApplysByAID(%d) err(%+v)", aid, err)
			return ecode.RequestErr
		}
		oldMap := make(map[int64]*archive.StaffApply)
		for _, k := range applys {
			oldMap[k.ApplyStaffMID] = k
		}
		var nvs, evs, dvs []*archive.StaffApply
		var change bool
		for _, v := range params {
			var (
				ov, ok = oldMap[v.MID]
				nv     = &archive.StaffApply{ApplyAID: aid}
				ovChg  bool
			)
			if !ok {
				//add staff
				nv.ApplyStaffMID = v.MID
				nv.Type = archive.TYPEUPADD
				nv.ApplyTitle = v.Title
				nv.ApplyTitleID = v.TitleID
				nv.State = archive.APPLYSTATEOPEN
				nvs = append(nvs, nv)
				change = true
			} else {
				// NOTE: edit staff
				*nv = *ov
				//注意up编辑时 应对staff 结束工单的逻辑  up操作直接覆盖staff申请单
				if nv.ApplyTitle != v.Title {
					nv.ApplyTitle = v.Title
					if nv.StaffState == archive.STATEON {
						//edit 逻辑 diff两次 解决aba 反复修改恢复问题
						nv.Type = archive.TYPEUPMODIFY
						if v.Title == nv.StaffTitle {
							nv.State = archive.APPLYSTATEDEL
						} else {
							nv.State = archive.APPLYSTATEOPEN
						}
					} else {
						//直接修改 不更新type
						nv.State = archive.APPLYSTATEOPEN
					}
					ovChg = true
					change = true
				}
				if ovChg {
					evs = append(evs, nv)
				}
				delete(oldMap, nv.ApplyStaffMID)
			}
		}
		//del staff
		if len(oldMap) > 0 {
			for _, v := range oldMap {
				//1.up 删除了staff 解除请求的工单   staff 发起解除  up 编辑页 删除操作 这算 up 同意了  staff工单
				if v.Type == archive.TYPESTAFFDEL && v.State == archive.APPLYSTATEOPEN {
					//up编辑页同意了staff解除申请 工单状态设置为已删除 （已解除）工单性质不变
					v.State = archive.APPLYSTATEDEL
				} else {
					//2.up 申请解除合作
					v.Type = archive.TYPEUPDEL
					if v.StaffState == archive.STATEON {
						v.State = archive.APPLYSTATEOPEN
					} else {
						v.State = archive.APPLYSTATEDEL
					}
				}
				change = true
				dvs = append(dvs, v)
			}
		}
		if !change {
			return
		}
		if len(nvs) > 0 {
			//新增保持顺序
			for _, v := range nvs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				if _, err = s.AddApply(c, apply, "HandleArchiveApplys/edit/addItem"); err != nil {
					log.Error(" HandleArchiveApplys  Edit s.AddApply(%d) err(%+v)", aid, err)
					err = nil
				}
			}
		}
		if len(evs) > 0 {
			for _, v := range evs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				g.Go(func() error {
					_, err = s.AddApply(c, apply, "HandleArchiveApplys/edit/editItem")
					return err
				})
			}
		}
		if len(dvs) > 0 {
			for _, v := range dvs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				g.Go(func() error {
					_, err = s.AddApply(c, apply, "HandleArchiveApplys/edit/delItem")
					return err
				})
			}
		}
		log.Info(" s.HandleArchiveApplys success aid(%d) params(%+v) nvs(%+v) evs(%+v) dvs(%+v) err(%+v)", aid, params, nvs, evs, dvs, err)
		return g.Wait()
	case "admin_edit":
		//管理员修改是直接生效  数据同步从staff->apply 自动cancel掉未处理掉的applys
		var applys []*archive.StaffApply
		if applys, err = s.arc.ApplysByAID(c, aid); err != nil || applys == nil {
			log.Error(" s.ApplysByAID(%d) err(%+v)", aid, err)
			return ecode.RequestErr
		}
		oldMap := make(map[int64]*archive.StaffApply)
		for _, k := range applys {
			oldMap[k.ApplyStaffMID] = k
		}
		var nvs, evs, dvs []*archive.StaffApply
		var change bool
		for _, v := range params {
			var (
				ov, ok = oldMap[v.MID]
				nv     = &archive.StaffApply{ApplyAID: aid}
				ovChg  bool
			)
			if !ok {
				//add staff
				nv.ApplyStaffMID = v.MID
				nv.Type = archive.TYPEUPADD
				nv.ApplyTitle = v.Title
				nv.ApplyTitleID = v.TitleID
				nv.State = archive.APPLYSTATEOPEN
				nvs = append(nvs, nv)
				change = true
			} else {
				// NOTE: edit staff
				*nv = *ov
				//注意up编辑时 应对staff 结束工单的逻辑  up操作直接覆盖staff申请单
				if nv.StaffTitle != v.Title {
					nv.StaffTitle = v.Title
					ovChg = true
					change = true
				}
				if ovChg {
					evs = append(evs, nv)
				}
				delete(oldMap, nv.ApplyStaffMID)
			}
		}
		//del staff
		if len(oldMap) > 0 {
			for _, v := range oldMap {
				//基于staffs
				if v.StaffState == archive.STATEOFF {
					continue
				}
				v.Type = archive.TYPEADMINDEL
				v.State = archive.APPLYSTATEDEL
				change = true
				dvs = append(dvs, v)
			}
		}
		if !change {
			return
		}
		//admin暂时不做修改
		nvs = (nvs)[0:0]
		evs = (evs)[0:0]
		if len(nvs) > 0 {
			for _, v := range nvs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				g.Go(func() error {
					_, err = s.AddApply(c, apply, "HandleArchiveApplys/edit/addItem")
					return err
				})
			}
		}
		if len(evs) > 0 {
			for _, v := range evs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				g.Go(func() error {
					_, err = s.AddApply(c, apply, "HandleArchiveApplys/edit/editItem")
					return err
				})
			}
		}
		//第一阶段只做admin 删除staff逻辑
		if len(dvs) > 0 {
			adminNotifyUp := false
			var StaffsName []string
			for _, v := range dvs {
				if pfl, _ := s.profile(c, v.ApplyStaffMID); pfl != nil {
					StaffsName = append(StaffsName, pfl.Profile.Name)
				}
			}
			for _, v := range dvs {
				apply := &archive.ApplyParam{ApplyAID: aid}
				apply.Copy(v)
				if !adminNotifyUp {
					adminNotifyUp = true
					apply.NotifyUp = true
					apply.StaffsName = strings.Join(StaffsName, ",")
				}
				if _, err = s.AddApply(c, apply, "HandleArchiveApplys/admin/edit/delItem"); err != nil {
					log.Error(" HandleArchiveApplys/admin/edit/delItem  s.AddApply(%d) err(%+v)", aid, err)
					err = nil
				}
			}
		}
		log.Info(" s.HandleArchiveApplys success aid(%d) params(%+v) nvs(%+v) evs(%+v) dvs(%+v) err(%+v)", aid, params, nvs, evs, dvs, err)
		return g.Wait()
	default:
		err = ecode.RequestErr
	}
	return
}

//HandleUpApplys  staff批量拒绝up主
func (s *Service) HandleUpApplys(c context.Context, upMid, staffMid int64, source string) (err error) {
	//批量拒绝：STAFF批量拒绝某UP主时， STAFF针对此UP主的“待处理申请“，全部拒绝，无论“新增”、“修改”、”解绑"
	var applys []*archive.StaffApply
	if applys, err = s.arc.ApplysByMIDAndStaff(c, upMid, staffMid); err != nil {
		log.Error("s.arc.ApplysByMIDAndStaff(%d ,%d) error(%v)", upMid, staffMid, err)
		return
	}
	if len(applys) == 0 {
		return
	}
	var g = &errgroup.Group{}
	for _, v := range applys {
		if v.State != archive.APPLYSTATEOPEN {
			continue
		}
		apply := &archive.ApplyParam{}
		apply.Copy(v)
		apply.State = archive.APPLYSTATEREFUSE
		g.Go(func() error {
			_, err = s.AddApply(c, apply, "HandleUpApplys/batch/refuse")
			return err
		})
	}
	err = g.Wait()
	if err != nil {
		log.Error(" HandleUpApplys(%d,%d) aid(%d) err(%+v)", upMid, staffMid, err)
	}
	return
}

//FillApplyParam  参数补全
func (s *Service) FillApplyParam(c context.Context, staffParam *archive.ApplyParam) (err error) {
	//申请单提交时补充title  申请单不能修改title
	if staffParam.ApplyTitle == "" && staffParam.ID > 0 {
		var apply *archive.StaffApply
		if apply, err = s.arc.Apply(c, staffParam.ID); err != nil {
			log.Error("s.arc.Apply(%+v) error(%v)", staffParam, err)
			return
		}
		if apply == nil {
			err = ecode.NothingFound
			return
		}
		staffParam.ApplyUpMID = apply.ApplyUpMID
		staffParam.ApplyAID = apply.ApplyAID
		staffParam.ApplyStaffMID = apply.ApplyStaffMID
		staffParam.ApplyTitle = apply.ApplyTitle
		staffParam.ApplyTitleID = apply.ApplyTitleID
		staffParam.StaffState = apply.StaffState
		staffParam.StaffTitle = apply.StaffTitle
		//case 1 申请单staff 交互
		if staffParam.Type == 0 {
			staffParam.Type = apply.Type
		}
		if apply.ASID > 0 {
			staffParam.ASID = apply.ASID
		}
		//case 1 staff发起解除申请 额
	}
	//带上archive信息
	if staffParam.ApplyAID > 0 {
		var a *archive.Archive
		if a, err = s.arc.Archive(c, staffParam.ApplyAID); err != nil {
			log.Error("s.arc.Archive(%d) error(%v)", staffParam.ApplyAID, err)
			return
		}
		if a == nil {
			err = ecode.NothingFound
			return
		}
		staffParam.ApplyUpMID = a.Mid
		staffParam.Archive = a
	}
	if staffParam.ApplyAID == 0 || staffParam.ApplyStaffMID == 0 {
		return ecode.RequestErr
	}
	//参数验证
	if staffParam.Type > archive.TYPEADMINDEL || staffParam.Type < archive.TYPEUPADD {
		err = ecode.RequestErr
		return
	}
	//带上staff线上数据
	if staffParam.ApplyStaffMID > 0 && staffParam.ApplyAID > 0 {
		var staff *archive.Staff
		if staff, err = s.arc.StaffByAidAndMid(c, staffParam.ApplyAID, staffParam.ApplyStaffMID); err != nil {
			log.Error("s.arc.StaffByAidAndMid(%+v) error(%v)", staffParam, err)
			return
		}
		if staff != nil {
			staffParam.StaffTitle = staff.StaffTitle
			staffParam.StaffState = staff.State
			staffParam.ASID = staff.ID
		}
		//staff申请解除时补充职位
		if staffParam.ApplyTitle == "" {
			staffParam.ApplyTitle = staff.StaffTitle
		}
	}
	return
}

//SyncStaffAttr .
func (s *Service) SyncStaffAttr(aid int64) (err error) {
	var staffs []*archive.Staff
	if staffs, err = s.arc.Staffs(context.TODO(), aid); err != nil {
		log.Error("SyncStaffAttr aid(%d) s.arc.Staffs(%+v) error(%v)", aid, staffs, err)
		return
	}
	var isStaff bool
	if len(staffs) == 0 {
		isStaff = false
	} else {
		isStaff = true
	}
	var a *archive.Archive
	if a, err = s.arc.Archive(context.TODO(), aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	//写属性位
	if s.isStaff(a) != isStaff {
		var attVal int32
		var tx *sql.Tx
		if tx, err = s.arc.BeginTran(context.TODO()); err != nil {
			log.Error("s.arc.BeginTran() error(%v)", err)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Error("wocao jingran recover le error(%v)", r)
			}
		}()
		if isStaff {
			attVal = 1
		} else {
			attVal = 0
		}
		if _, err = s.arc.TxUpArcAttr(tx, a.Aid, archive.AttrBitSTAFF, attVal); err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
			return
		}
	}
	//无脑同步动态
	s.worker.Add(func() {
		s.busSyncArchive(aid)
	})
	log.Info("SyncStaffAttr aid(%d) isStaff(%v)  success", aid, isStaff)
	return
}

//HiddenApply .无效的申请单
func (s *Service) HiddenApply(staffParam *archive.StaffApply) bool {
	return (staffParam.Type == archive.TYPEUPDEL && (staffParam.State == archive.APPLYSTATEDEL || staffParam.State == archive.APPLYSTATEACCEPT)) ||
		(staffParam.Type == archive.TYPESTAFFDEL && (staffParam.State == archive.APPLYSTATEDEL || staffParam.State == archive.APPLYSTATEACCEPT)) ||
		(staffParam.Type == archive.TYPEADMINDEL && (staffParam.State == archive.APPLYSTATEDEL || staffParam.State == archive.APPLYSTATEACCEPT)) ||
		(staffParam.Type == archive.TYPEUPADD && (staffParam.State == archive.APPLYSTATEREFUSE || staffParam.State == archive.APPLYSTATEDEL))
}

//IsMidSilence .
func (s *Service) IsMidSilence(c context.Context, mid int64) bool {
	if pfl, _ := s.profile(c, mid); pfl != nil {
		return pfl.Profile.Silence == 1
	}
	return false

}

//HandleMsg .
func (s *Service) HandleMsg(c context.Context, staffParam *archive.ApplyParam, source string) {
	if pfl, _ := s.profile(c, staffParam.ApplyUpMID); pfl != nil {
		staffParam.UpName = pfl.Profile.Name
	}
	if pfl, _ := s.profile(c, staffParam.ApplyStaffMID); pfl != nil {
		staffParam.StaffName = pfl.Profile.Name
	}
	//管理员删除通知 需要通知双方
	if staffParam.MsgId == archive.MSG_12 {
		if staffParam.NotifyUp {
			//admin批量删除时 通知up仅一次 其他staff每人一次
			var msg2Up = archive.ArgMsg{MSGID: archive.MSG_12, Apply: staffParam}
			s.sendMsg(&msg2Up)
		}
		var msg2Staff = archive.ArgMsg{MSGID: archive.MSG_13, Apply: staffParam}
		s.sendMsg(&msg2Staff)
	} else {
		var msg = archive.ArgMsg{MSGID: staffParam.MsgId, Apply: staffParam}
		s.sendMsg(&msg)
	}

}

//DispatchEvent  集中处理事件  up侧  staff侧  admin侧
func (s *Service) DispatchEvent(c context.Context, staffParam *archive.ApplyParam, source string) {
	staffParam.DealState = 0
	switch staffParam.Type {
	//五类申请单
	case archive.TYPEUPADD:
		switch staffParam.State {
		case archive.APPLYSTATEOPEN:
			staffParam.MsgId = archive.MSG_1
			staffParam.DealState = archive.DEALSTATEOPEN
		case archive.APPLYSTATEREFUSE:
			staffParam.MsgId = archive.MSG_3
		case archive.APPLYSTATEACCEPT:
			staffParam.SyncStaff = true
			staffParam.MsgId = archive.MSG_2
			staffParam.DealState = archive.DEALSTATEDONE
			staffParam.CleanCache = true
		case archive.APPLYSTATEIGNORE:
			staffParam.NoNotify = true
			staffParam.DealState = archive.DEALSTATEIGNORE
		case archive.APPLYSTATEDEL:
			staffParam.NoNotify = true
		}
	case archive.TYPEUPMODIFY:
		switch staffParam.State {
		case archive.APPLYSTATEOPEN:
			staffParam.MsgId = archive.MSG_4
			staffParam.DealState = archive.DEALSTATEOPEN
		case archive.APPLYSTATEREFUSE:
			if staffParam.StaffState == archive.STATEON {
				staffParam.OldTitle = staffParam.ApplyTitle
				staffParam.ApplyTitle = staffParam.StaffTitle
				staffParam.MsgId = archive.MSG_6
			}
			staffParam.DealState = archive.DEALSTATEDONE
		case archive.APPLYSTATEACCEPT:
			staffParam.MsgId = archive.MSG_5
			staffParam.SyncStaff = true
			staffParam.CleanCache = true
			staffParam.DealState = archive.DEALSTATEDONE
		case archive.APPLYSTATEIGNORE:
			staffParam.NoNotify = true
			staffParam.DealState = archive.DEALSTATEIGNORE
		case archive.APPLYSTATEDEL:
			//处理up反复修复恢复case  自动关闭上一个（撤销上一个修改操作）
			staffParam.NoNotify = true
		}
	case archive.TYPEUPDEL:
		switch staffParam.State {
		case archive.APPLYSTATEOPEN:
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_9
			}
			staffParam.DealState = archive.DEALSTATEOPEN
		case archive.APPLYSTATEREFUSE:
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_11
			}
			staffParam.DealState = archive.DEALSTATEDONE
		case archive.APPLYSTATEACCEPT:
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_10
				staffParam.SyncStaff = true
				staffParam.CleanCache = true
			}
		case archive.APPLYSTATEIGNORE:
			staffParam.DealState = archive.DEALSTATEIGNORE
			staffParam.NoNotify = true
		case archive.APPLYSTATEDEL:
			staffParam.NoNotify = true
		}
	case archive.TYPESTAFFDEL:
		//staff申请解除 up只能在编辑页去删除对应行 up edit submit type应该是 TYPESTAFFDEL   且不提交删除行
		//up 主处理staff工单 只能在编辑页的话  怎么体现 同意或者拒绝呢  up主又在staff申请解除基础上无视申请 再继续修改title
		switch staffParam.State {
		case archive.APPLYSTATEOPEN:
			staffParam.MsgId = archive.MSG_7
			staffParam.DealState = archive.DEALSTATEDONE
		case archive.APPLYSTATEREFUSE:
			if staffParam.StaffState == archive.STATEON {
				//拒绝 不发消息？？？
				staffParam.MsgId = 0
			}
		case archive.APPLYSTATEACCEPT:
			//目前并没有在 staff申请解除时  up可以操作到同意的转换
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_8
				staffParam.CleanCache = true
				staffParam.SyncStaff = true
			}
		case archive.APPLYSTATEIGNORE:
			staffParam.NoNotify = true
		case archive.APPLYSTATEDEL:
			//staff申请删除  up手动删视为同意
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_8
				staffParam.CleanCache = true
				staffParam.SyncStaff = true
				staffParam.NoNotify = false
			}
		}
	case archive.TYPEADMINDEL:
		switch staffParam.State {
		case archive.APPLYSTATEDEL:
			if staffParam.StaffState == archive.STATEON {
				staffParam.MsgId = archive.MSG_12
				staffParam.CleanCache = true
				staffParam.SyncStaff = true
			}
		}
	}
}

//DoApply .
func (s *Service) DoApply(c context.Context, staffParam *archive.ApplyParam, source string) (id int64, err error) {
	//同意，拒绝，忽略 前置验证
	if staffParam.ID > 0 {
		switch staffParam.State {
		case archive.APPLYSTATEACCEPT, archive.APPLYSTATEREFUSE, archive.APPLYSTATEIGNORE:
			var oldApply *archive.StaffApply
			if oldApply, err = s.arc.Apply(c, staffParam.ID); err != nil {
				log.Error("s.arc.Apply(%+v) error(%v)", staffParam, err)
				return
			}
			if oldApply == nil {
				err = ecode.VideoupStaffApply404
				return
			}
			//前置状态
			if oldApply.State != archive.APPLYSTATEOPEN && oldApply.State != archive.APPLYSTATEIGNORE {
				err = ecode.VideoupStaffApplyStateNotMatch
				return
			}
			//工单性质改变(up/admin修改了)
			if oldApply.Type != staffParam.Type {
				err = ecode.VideoupStaffApplyTypeChange
				return
			}
			if staffParam.Type == archive.TYPEUPADD && oldApply.StaffState == archive.STATEON {
				err = ecode.RequestErr
				return
			}
			if (staffParam.Type == archive.TYPEUPMODIFY || staffParam.Type == archive.TYPEUPDEL) && oldApply.StaffState == archive.STATEOFF {
				err = ecode.RequestErr
				return
			}
			staffParam.ApplyUpMID = oldApply.ApplyUpMID
		default:
			err = ecode.RequestErr
			return
		}
	} else {
		//申请解除参数错误
		if staffParam.ApplyAID == 0 {
			err = ecode.RequestErr
			return
		}
		//staff申请解除
		var oldApplys []*archive.StaffApply
		if oldApplys, err = s.arc.ApplysByAID(c, staffParam.ApplyAID); err != nil {
			log.Error("s.ApplysByAID.Apply(%+v) error(%v)", staffParam, err)
			return
		}
		if len(oldApplys) == 0 {
			err = ecode.NothingFound
			return
		}
		in := false
		checkState := false
		for _, v := range oldApplys {
			if v.ApplyStaffMID == staffParam.ApplyStaffMID && !in {
				in = true
				staffParam.ApplyUpMID = v.ApplyUpMID
				if v.StaffState == archive.STATEON && !checkState {
					checkState = true
				}
				break
			}
		}
		//是否是staff
		if !in {
			err = ecode.VideoupStaffApplyMidNotIn
			return
		}
		//是否可以发起解除操作
		if !checkState {
			err = ecode.VideoupStaffApplyStateNotMatch
			return
		}
	}
	//拉黑关系 拦截
	if staffParam.State != archive.APPLYSTATEIGNORE {
		var relation int64
		if relation, err = s.relation.Relation(context.TODO(), staffParam.ApplyUpMID, staffParam.ApplyStaffMID, staffParam.ApplyAID); err != nil {
			log.Error("s.relation.Relation(%+v) error(%v)", staffParam, err)
		}
		if relation >= archive.UPRELATIONBLACK {
			err = ecode.VideoupStaffApplyUpMidBlack
			return
		}
	}
	//staff是否被封禁
	if s.IsMidSilence(c, staffParam.ApplyStaffMID) {
		log.Error("s.IsMidSilence(%+v) error(%v)", staffParam, err)
		err = ecode.VideoupStaffMidSilence
		return
	}
	_, err = s.AddApply(c, staffParam, source)
	s.SearchUpdate(c, staffParam, source)
	return
}

//SearchUpdate .
func (s *Service) SearchUpdate(c context.Context, staffParam *archive.ApplyParam, source string) (id int64, err error) {
	var indexItems = make([]*archive.IndexItem, 0)
	var applys []*archive.StaffApply
	if applys, err = s.arc.ApplysByAID(c, staffParam.ApplyAID); err != nil {
		log.Error("s.ApplysByAID.Apply(%+v) error(%v)", staffParam, err)
		return
	}
	if len(applys) == 0 {
		err = ecode.NothingFound
		return
	}

	for _, v := range applys {
		indexItem := &archive.IndexItem{DealState: v.DealState, ApplyStaffMID: v.ApplyStaffMID}
		indexItems = append(indexItems, indexItem)
	}
	if len(indexItems) > 0 {
		var dataByte []byte
		indexs := make([]*archive.Index, 0)
		indexs = append(indexs, &archive.Index{ID: staffParam.ApplyAID, Item: indexItems})
		indexData := &archive.SearchApplyIndex{Indexs: indexs}
		if dataByte, err = json.Marshal(indexData); err != nil {
			log.Error("SearchUpdate  json.Marshal(%+v) error(%v)", indexData, err)
			return
		}
		s.mng.SearchUpdate(c, "creative_archive_staff", string(dataByte), staffParam.ApplyAID)
	}
	return
}

//AddApply .
//a.【UP】稿件添加或者编辑批量修改applys   []mid title title_id  diff后自动生成对应type 批量生产 add/edit/del 对应的type
//b.【STAFF】申请单交互流转 (同意、拒绝、忽略、staff申请解除)   id  staff_mid  flag_add_black  flag_refuse_mid
//c.【ADMIN】审核后台提交的批量删除
func (s *Service) AddApply(c context.Context, staffParam *archive.ApplyParam, source string) (id int64, err error) {
	var hits = make([]string, 0)
	hits = append(hits, fmt.Sprintf("hitSource(%s)", source))
	log.Info("aid (%d) AddApply params(%+v) source(%s) start", staffParam.ApplyAID, staffParam, source)
	if err = s.FillApplyParam(c, staffParam); err != nil {
		log.Error("AddApply(%d) s.FillApplyParam(%+v) error(%v)", staffParam.ApplyAID, staffParam, err)
		return
	}
	s.DispatchEvent(c, staffParam, source)
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("AddApply(%d) s.arc.BeginTran() error(%v)", staffParam.ApplyAID, err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("AddApply(%d) wocao jingran recover le error(%v)", staffParam.ApplyAID, r)
		}
	}()
	if id, err = s.arc.TxAddApply(tx, staffParam); err != nil {
		log.Error("AddApply(%d) TxAddApply(%d) err(%+v)", staffParam.ApplyAID, staffParam, err)
		tx.Rollback()
		return
	}
	//1.传导至staff
	if staffParam.SyncStaff {
		staff := &archive.Staff{}
		staff.Copy(staffParam)
		//for log trace
		staffParam.StaffState = staff.State
		var staffID int64
		if staffID, err = s.arc.TxAddStaff(tx, staff); err != nil {
			log.Error("AddApply(%d) TxAddStaff(%d) err(%+v)", staffParam.ApplyAID, staffParam, err)
			tx.Rollback()
			return
		}
		if staffParam.Archive.State >= archive.StateOpen {
			staffParam.SyncDynamic = true
		}
		hits = append(hits, fmt.Sprintf("addStaff AID(%d)", staffParam.ApplyAID))
		//2.回写as_id建立绑定关系
		if staffParam.ASID == 0 {
			staffParam.ASID = staffID
			if _, err = s.arc.TxAddApply(tx, staffParam); err != nil {
				log.Error("AddApply(%d) TxAddApply(%d) renew as_id(%d) err(%+v)", staffParam.ApplyAID, staffParam, staffParam.ASID, err)
				tx.Rollback()
				return
			}
			hits = append(hits, fmt.Sprintf(" bindAsid(%d) AID(%d)", staffParam.ASID, staffParam.ApplyAID))
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("AddApply(%d) tx.Commit() error(%v)", staffParam.ApplyAID, err)
		return
	}
	//3.通知
	if !staffParam.NoNotify && staffParam.MsgId > 0 {
		hits = append(hits, fmt.Sprintf("sendMsg (%d)", staffParam.MsgId))
		s.worker.Add(func() {
			s.HandleMsg(context.TODO(), staffParam, source)
		})
	}
	//4.清理cache
	if staffParam.CleanCache {
		hits = append(hits, fmt.Sprintf("cleanCache AID(%d)", staffParam.ApplyAID))
		s.worker.Add(func() {
			log.Info("aid(%d) apply(%+v) cleanCache", staffParam.ApplyAID, staffParam)
			s.arc.DelCacheStaffData(context.TODO(), staffParam.ApplyAID)
		})
	}
	//6 拉黑
	if staffParam.State == archive.APPLYSTATEREFUSE && staffParam.FlagAddBlack {
		hits = append(hits, fmt.Sprintf("staff AddBlack MID(%d) AID(%d)", staffParam.ApplyUpMID, staffParam.ApplyAID))
		s.worker.Add(func() {
			if err = s.relation.AddBalck(context.TODO(), staffParam.ApplyStaffMID, staffParam.ApplyUpMID, staffParam.ApplyAID); err != nil {
				log.Error("AddApply(%d) s.relation.AddBalck(%+v) error(%v)", staffParam.ApplyAID, staffParam, err)
			}
		})
	}
	//7 staff拒绝Up
	if staffParam.State == archive.APPLYSTATEREFUSE && staffParam.FlagRefuse {
		hits = append(hits, fmt.Sprintf("staff RefuseMid A(%d) rejuse B(%d)", staffParam.ApplyStaffMID, staffParam.ApplyUpMID))
		s.worker.Add(func() {
			if err = s.HandleUpApplys(context.TODO(), staffParam.ApplyUpMID, staffParam.ApplyStaffMID, "RefuseMid"); err != nil {
				log.Error("AddApply(%d) s.HandleUpApplys(%+v) error(%v)", staffParam.ApplyAID, staffParam, err)
			}
		})
	}
	//8.staff属性位聚合
	if staffParam.SyncStaff {
		hits = append(hits, fmt.Sprintf("syncStaffAttr AID(%d)", staffParam.ApplyAID))
		s.worker.Add(func() {
			if err = s.SyncStaffAttr(staffParam.ApplyAID); err != nil {
				log.Error("AddApply(%d) s.SyncStaffAttr(%+v) error(%v)", staffParam.ApplyAID, staffParam, err)
			}
		})
	}
	//9.行为日志
	hits = append(hits, fmt.Sprintf("AddStaffLog AID(%d)", staffParam.ApplyAID))
	s.worker.Add(func() {
		index := []interface{}{staffParam.ApplyAID, staffParam.ApplyStaffMID}
		content := map[string]interface{}{
			"aid":            staffParam.ApplyAID,
			"as_id":          staffParam.ASID,
			"mid":            staffParam.ApplyUpMID,
			"staff_mid":      staffParam.ApplyStaffMID,
			"apply_title":    staffParam.ApplyTitle,
			"apply_title_id": staffParam.ApplyTitleID,
			"state":          staffParam.State,
			"staff_state":    staffParam.StaffState,
			"staff_title":    staffParam.StaffTitle,
			"source":         source,
			"event":          hits,
		}
		s.AddAuditLog(context.Background(), archive.STAFFLogBizID, staffParam.Type, source, staffParam.ApplyStaffMID, staffParam.ApplyTitle, []int64{staffParam.ApplyAID}, index, content)
	})
	log.Info(" aid (%d) AddApply  success  params(%+v) source(%s) and hits event(%v)", staffParam.ApplyAID, staffParam, source, hits)
	return
}
