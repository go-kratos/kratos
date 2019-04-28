package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/usersuit/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// PendantInfoList pendant list by group rank.
func (s *Service) PendantInfoList(c context.Context, arg *model.ArgPendantGroupList) (pis []*model.PendantInfo, pager *model.Pager, err error) {
	var (
		total int64
		pids  []int64
		ppm   map[int64][]*model.PendantPrice
	)
	pager = &model.Pager{
		PN: arg.PN,
		PS: arg.PS,
	}
	// all group pendants
	if arg.GID == 0 {
		if total, err = s.d.PendantGroupInfoTotal(c); err != nil {
			err = errors.Wrap(err, "s.d.PendantGroupInfoTotal()")
			return
		}
		if total <= 0 {
			return
		}
		pager.Total = total
		if pis, pids, err = s.d.PendantInfoAll(c, arg.PN, arg.PS); err != nil {
			err = errors.Wrapf(err, "s.d.PendantInfoAll(%d,%d)", arg.PN, arg.PS)
			return
		}
		if len(pis) == 0 || len(pids) == 0 {
			log.Warn("no pendant list")
			return
		}
		if ppm, err = s.d.PendantPriceIDs(c, pids); err != nil {
			err = errors.Wrapf(err, "s.d.PendantPriceIDs(%s)", xstr.JoinInts(pids))
			return
		}
		for _, pi := range pis {
			if pp, ok := ppm[pi.ID]; ok {
				pi.Prices = pp
			}
		}
		return
	}
	// one group pendants
	if total, err = s.d.PendantGroupRefsGidTotal(c, arg.GID); err != nil {
		err = errors.Wrap(err, "s.d.PendantGroupRefsGidTotal()")
		return
	}
	if total <= 0 {
		return
	}
	pager.Total = total
	var pg *model.PendantGroup
	if pids, err = s.d.PendantGroupPIDs(c, arg.GID, arg.PN, arg.PS); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroupPIDs(%d,%d,%d)", arg.GID, arg.PN, arg.PS)
		return
	}
	if pg, err = s.d.PendantGroupID(c, arg.GID); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroupID(%d)", arg.GID)
		return
	}
	if len(pids) == 0 {
		log.Warn("no pendant group relation")
		return
	}
	if pis, _, err = s.d.PendantInfoIDs(c, pids); err != nil {
		err = errors.Wrapf(err, "s.d.PendantInfoIDs(%s)", xstr.JoinInts(pids))
		return
	}
	if len(pis) == 0 {
		log.Warn("no pendant list")
		return
	}
	if ppm, err = s.d.PendantPriceIDs(c, pids); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPriceIDs(%s)", xstr.JoinInts(pids))
		return
	}
	for _, pi := range pis {
		if pp, ok := ppm[pi.ID]; ok {
			pi.Prices = pp
		}
		pi.GID = arg.GID
		pi.GroupName = pg.Name
		pi.GroupRank = pg.Rank
	}
	return
}

// PendantInfoID pendant info by pid and gid.
func (s *Service) PendantInfoID(c context.Context, pid, gid int64) (pi *model.PendantInfo, err error) {
	if pi, err = s.d.PendantInfoID(c, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantInfoID(%d)", pid)
		return
	}
	if pi == nil {
		err = ecode.PendantNotFound
		return
	}
	var (
		pg  *model.PendantGroup
		ppm map[int64][]*model.PendantPrice
	)
	if pg, err = s.d.PendantGroupID(c, gid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroupID(%d)", gid)
		return
	}
	if ppm, err = s.d.PendantPriceIDs(c, []int64{pid}); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPriceIDs(%d)", pid)
		return
	}
	pi.GID = gid
	pi.GroupName = pg.Name
	pi.GroupRank = pg.Rank
	if pp, ok := ppm[pid]; ok {
		pi.Prices = pp
	}
	return
}

// PendantGroupID pendant group by ID.
func (s *Service) PendantGroupID(c context.Context, gid int64) (pg *model.PendantGroup, err error) {
	if pg, err = s.d.PendantGroupID(c, gid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroupID(%d)", gid)
	}
	return
}

// PendantGroupList group page.
func (s *Service) PendantGroupList(c context.Context, arg *model.ArgPendantGroupList) (pgs []*model.PendantGroup, pager *model.Pager, err error) {
	var total int64
	pager = &model.Pager{
		PN: arg.PN,
		PS: arg.PS,
	}
	if total, err = s.d.PendantGroupsTotal(c); err != nil {
		err = errors.Wrap(err, "s.d.PendantGroupsTotal()")
		return
	}
	if total <= 0 {
		return
	}
	pager.Total = total
	if pgs, err = s.d.PendantGroups(c, arg.PN, arg.PS); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroups(%d,%d)", arg.PN, arg.PS)
	}
	return
}

// PendantGroupAll all groups.
func (s *Service) PendantGroupAll(c context.Context) (pgs []*model.PendantGroup, err error) {
	if pgs, err = s.d.PendantGroupAll(c); err != nil {
		err = errors.Wrap(err, "s.d.PendantGroupAll()")
		return
	}
	return
}

// PendantInfoAllNoPage all info on no page.
func (s *Service) PendantInfoAllNoPage(c context.Context) (pis []*model.PendantInfo, err error) {
	if pis, err = s.d.PendantInfoAllNoPage(c); err != nil {
		err = errors.Wrap(err, "s.d.PendantInfoAllNoPage()")
		return
	}
	return
}

// AddPendantInfo add pendantInfo .
func (s *Service) AddPendantInfo(c context.Context, arg *model.ArgPendantInfo) (err error) {
	var pg *model.PendantGroup
	if pg, err = s.d.PendantGroupID(c, arg.GID); err != nil {
		err = errors.Wrapf(err, "s.d.PendantGroupID(%d)", arg.GID)
		return
	}
	if pg == nil {
		err = errors.New("group no exist")
		return
	}
	// begin tran
	var tx *sql.Tx
	if tx, err = s.d.BeginTran(c); err != nil {
		err = errors.Wrap(err, "s.arc.BeginTran()")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	pi := &model.PendantInfo{
		Name:       arg.Name,
		Image:      arg.Image,
		ImageModel: arg.ImageModel,
		Status:     arg.Status,
		Rank:       arg.Rank,
	}
	var pid int64
	if pid, err = s.d.TxAddPendantInfo(tx, pi); err != nil {
		err = errors.Wrapf(err, "s.d.TxAddPendantInfo(%+v)", pi)
		return
	}
	arg.PID = pid
	pr := &model.PendantGroupRef{GID: arg.GID, PID: arg.PID}
	if _, err = s.d.TxAddPendantGroupRef(tx, pr); err != nil {
		err = errors.Wrapf(err, "s.d.TxAddPendantGroupRef(%+v)", pr)
		return
	}
	pp := &model.PendantPrice{}
	for _, tp := range model.PriceTypes {
		pp.BulidPendantPrice(arg, tp)
		if pp.Price != 0 {
			if _, err = s.d.TxAddPendantPrices(tx, pp); err != nil {
				err = errors.Wrapf(err, "s.d.TxAddPendantPrices(%+v)", pp)
				return
			}
		}
	}
	return
}

// UpPendantInfo  update pendant info .
func (s *Service) UpPendantInfo(c context.Context, arg *model.ArgPendantInfo) (err error) {
	// begin tran
	var tx *sql.Tx
	if tx, err = s.d.BeginTran(c); err != nil {
		err = errors.Wrap(err, "s.arc.BeginTran()")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	pi := &model.PendantInfo{
		ID:         arg.PID,
		Name:       arg.Name,
		Image:      arg.Image,
		ImageModel: arg.ImageModel,
		Status:     arg.Status,
		Rank:       arg.Rank,
	}
	if _, err = s.d.TxUpPendantInfo(tx, pi); err != nil {
		err = errors.Wrapf(err, "s.d.TxUpPendantInfo(%+v)", pi)
		return
	}
	if _, err = s.d.TxUpPendantGroupRef(tx, arg.GID, arg.PID); err != nil {
		err = errors.Wrapf(err, "s.d.TxUpPendantGroupRef(%d,%d)", arg.GID, arg.PID)
		return
	}
	pp := &model.PendantPrice{}
	for _, tp := range model.PriceTypes {
		pp.BulidPendantPrice(arg, tp)
		if pp.Price != 0 {
			if _, err = s.d.TxAddPendantPrices(tx, pp); err != nil {
				err = errors.Wrapf(err, "s.d.TxAddPendantPrices(%+v)", pp)
				return
			}
		}
	}
	return
}

// UpPendantGroupStatus update pendant group status
func (s *Service) UpPendantGroupStatus(c context.Context, gid int64, status int8) (err error) {
	if _, err = s.d.UpPendantGroupStatus(c, gid, status); err != nil {
		err = errors.Wrapf(err, "s.d.UpPendantGroupStatus(%d,%d)", gid, status)
	}
	return
}

// UpPendantInfoStatus update pendant info status
func (s *Service) UpPendantInfoStatus(c context.Context, pid int64, status int8) (err error) {
	if _, err = s.d.UpPendantInfoStatus(c, pid, status); err != nil {
		err = errors.Wrapf(err, "s.d.UpPendantInfoStatus(%d,%d)", pid, status)
	}
	return
}

// AddPendantGroup update pendant group.
func (s *Service) AddPendantGroup(c context.Context, arg *model.ArgPendantGroup) (err error) {
	pg := &model.PendantGroup{
		Name:   arg.Name,
		Rank:   arg.Rank,
		Status: arg.Status,
	}
	if _, err = s.d.AddPendantGroup(c, pg); err != nil {
		err = errors.Wrapf(err, "s.d.AddPendantGroup(%+v)", pg)
	}
	return
}

// UpPendantGroup update pendant group.
func (s *Service) UpPendantGroup(c context.Context, arg *model.ArgPendantGroup) (err error) {
	if arg.GID == 0 {
		err = ecode.PendantNotFound
		return
	}
	pg := &model.PendantGroup{
		ID:     arg.GID,
		Name:   arg.Name,
		Rank:   arg.Rank,
		Status: arg.Status,
	}
	if _, err = s.d.UpPendantGroup(c, pg); err != nil {
		err = errors.Wrapf(err, "s.d.UpPendantGroup(%+v)", pg)
	}
	return
}

// PendantOrders get pendant order historys.
func (s *Service) PendantOrders(c context.Context, arg *model.ArgPendantOrder) (pos []*model.PendantOrder, pager *model.Pager, err error) {
	var total int64
	pager = &model.Pager{
		PN: arg.PN,
		PS: arg.PS,
	}
	if total, err = s.d.MaxOrderHistory(c); err != nil {
		err = errors.Wrapf(err, "s.d.MaxOrderHistory(%+v)", arg)
		return
	}
	if total <= 0 {
		return
	}
	pager.Total = total
	var pids []int64
	if pos, pids, err = s.d.OrderHistorys(c, arg); err != nil {
		err = errors.Wrapf(err, "s.d.OrderHistorys(%+v)", arg)
		return
	}
	if len(pids) == 0 {
		return
	}
	var pim map[int64]*model.PendantInfo
	if _, pim, err = s.d.PendantInfoIDs(c, pids); err != nil {
		err = errors.Wrapf(err, "s.d.PendantInfoIDs(%v)", xstr.JoinInts(pids))
		return
	}
	for _, v := range pos {
		if pi, ok := pim[v.PID]; ok {
			v.PName = pi.Name
		}
		v.CoverToPlatform()
	}
	return
}

// PendantPKG  get pendant in pkg.
func (s *Service) PendantPKG(c context.Context, uid int64) (pkgs []*model.PendantPKG, equip *model.PendantPKG, err error) {
	if equip, err = s.d.PendantEquipUID(c, uid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantEquipUID(%d)", uid)
		return
	}
	if pkgs, err = s.d.PendantPKGs(c, uid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPKGs(%d)", uid)
		return
	}
	if len(pkgs) == 0 {
		return
	}
	var time = time.Now().Unix()
	for _, pkg := range pkgs {
		if equip != nil && pkg.PID == equip.PID {
			pkg.Status = model.PendantPKGOnEquip
		}
		if pkg.Status == model.PendantPKGValid && pkg.Expires < time {
			pkg.Status = model.PendantPKGInvalid
		}
	}
	return
}

// UserPKGDetails  get user pkg 's pendant details.
func (s *Service) UserPKGDetails(c context.Context, uid, pid int64) (pkg *model.PendantPKG, err error) {
	if pkg, err = s.d.PendantPKG(c, uid, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPKG(%d,%d)", uid, pid)
	}
	return
}

// EquipPendant equip pendant.
func (s *Service) EquipPendant(c context.Context, uid, pid int64) (err error) {
	var pkg *model.PendantPKG
	if pkg, err = s.d.PendantPKG(c, uid, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPKG(%d,%d)", uid, pid)
		return
	}
	if pkg == nil || pkg.Expires < time.Now().Unix() {
		log.Warn("pid(%d) not exist or expires(%d) is failed", pid, time.Now().Unix())
		return
	}
	if _, err = s.d.AddPendantEquip(c, pkg); err != nil {
		err = errors.Wrapf(err, "s.d.AddPendantEquip(%+v)", pkg)
		return
	}
	s.addAsyn(func() {
		if err = s.d.DelEquipsCache(context.Background(), []int64{uid}); err != nil {
			log.Error("s.d.DelEquipsCache(%d) error(%+v)", uid, err)
			return
		}
	})
	s.addAsyn(func() {
		if err = s.accNotify(context.Background(), uid, model.AccountNotifyUpdatePendant); err != nil {
			log.Error("s.accNotify(%d) error(%+v)", uid, err)
			return
		}
	})
	return
}

// UpPendantPKG update user pkg.
func (s *Service) UpPendantPKG(c context.Context, uid, pid int64, day int64, msg *model.SysMsg, oid int64) (err error) {
	var (
		pi  *model.PendantInfo
		pkg *model.PendantPKG
	)
	if pi, err = s.d.PendantInfoID(c, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantInfoID(%d)", pid)
		return
	}
	if pi == nil {
		err = ecode.PendantNotFound
		return
	}
	if pkg, err = s.d.PendantPKG(c, uid, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPKG(%d,%d)", uid, pid)
		return
	}
	if pkg == nil {
		pkg = &model.PendantPKG{}
	}
	var operAction string
	switch msg.Type {
	case model.PendantAddStyleDay:
		if pkg.Expires < time.Now().Unix() {
			pkg.Expires = time.Now().Unix() + day*86400
			pkg.PID = pid
			pkg.UID = uid
			pkg.TP = model.PendantAddStyleDay
		} else {
			pkg.Expires = pkg.Expires + day*86400
		}
		operAction = fmt.Sprintf("新增%s挂件 %d天", pi.Name, day)
	case model.PendantAddStyleDate:
		if pkg.ID == 0 {
			err = errors.New("no pkg")
			return
		}
		oldExpires := pkg.Expires
		pkg.Expires = day
		operAction = fmt.Sprintf("修改%s挂件 到期时间从%s至%s", pi.Name,
			xtime.Time(oldExpires).Time().Format("2006-01-02 15:04:05"),
			xtime.Time(pkg.Expires).Time().Format("2006-01-02 15:04:05"))
	}
	if _, err = s.d.AddPendantPKG(c, pkg); err != nil {
		err = errors.Wrapf(err, "s.d.AddPendantPKG(%+v)", pkg)
		return
	}
	if msg.IsMsg {
		s.addAsyn(func() {
			msg.Type = model.MsgTypeCustom
			title, content, ip := model.MsgInfo(msg)
			if err = s.d.SendSysMsg(context.Background(), []int64{uid}, title, content, ip); err != nil {
				log.Error("s.d.MutliSendSysMsg(%d,%s,%s,%s) error(%+v)", uid, title, content, ip, err)
			}
		})
	}
	s.addAsyn(func() {
		if err = s.d.DelPKGCache(context.Background(), []int64{uid}); err != nil {
			log.Error("s.d.DelPKGCache(%d) error(%+v)", uid, err)
		}
		if err = s.d.SetPendantPointCache(context.Background(), uid, pid); err != nil {
			log.Error("s.d.SetPendantPointCache(%d,%d) error(%+v)", uid, pid, err)
		}
		if _, err = s.d.AddPendantOperLog(context.Background(), oid, []int64{uid}, pid, operAction); err != nil {
			log.Error("s.d.AddPendantOperLog(%d,%s) error(%+v)", oid, operAction, err)
		}
	})
	s.addAsyn(func() {
		if err = s.accNotify(context.Background(), uid, model.AccountNotifyUpdatePendant); err != nil {
			log.Error("s.accNotify(%d) error(%+v)", uid, err)
			return
		}
	})
	return
}

// MutliSendPendant mutli send pendant.
func (s *Service) MutliSendPendant(c context.Context, uids []int64, pid int64, day int64, msg *model.SysMsg, oid int64) (err error) {
	var (
		pi           *model.PendantInfo
		opkgs, npkgs []*model.PendantPKG
	)
	if pi, err = s.d.PendantInfoID(c, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantInfoID(%d)", pid)
		return
	}
	if pi == nil {
		err = ecode.PendantNotFound
		return
	}
	if opkgs, err = s.d.PendantPKGUIDs(c, uids, pid); err != nil {
		err = errors.Wrapf(err, "s.d.PendantPKGUIDs(%+v,%d)", uids, pid)
		return
	}
	var ouids, nuids []int64
	for _, pkg := range opkgs {
		if pkg.Expires < time.Now().Unix() {
			pkg.Expires = time.Now().Unix() + day*86400
		} else {
			pkg.Expires = pkg.Expires + day*86400
		}
		ouids = append(ouids, pkg.UID)
	}
	nuids = s.diffSlice(uids, ouids)
	for _, nuid := range nuids {
		npkgs = append(npkgs, &model.PendantPKG{PID: pid, UID: nuid, Expires: time.Now().Unix() + day*86400})
	}
	// begin tran
	var tx *sql.Tx
	if tx, err = s.d.BeginTran(c); err != nil {
		err = errors.Wrap(err, "s.arc.BeginTran()")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if len(npkgs) != 0 {
		if _, err = s.d.TxAddPendantPKGs(tx, npkgs); err != nil {
			err = errors.Wrapf(err, "s.d.TxAddPendantPKGs(%d,%s)", pid, xstr.JoinInts(nuids))
			return
		}
	}
	if len(opkgs) != 0 {
		if _, err = s.d.TxUpPendantPKGs(tx, opkgs); err != nil {
			err = errors.Wrapf(err, "s.d.TxUpPendantPKGs(%d,%s)", pid, xstr.JoinInts(ouids))
			return
		}
	}
	if msg.IsMsg {
		s.addAsyn(func() {
			title, content, ip := model.MsgInfo(msg)
			if err = s.d.MutliSendSysMsg(context.Background(), uids, title, content, ip); err != nil {
				log.Error("s.d.MutliSendSysMsg(%s,%s,%s,%s) error(%+v)", xstr.JoinInts(uids), title, content, ip, err)
			}
		})
	}
	s.addAsyn(func() {
		if err = s.d.DelPKGCache(context.Background(), uids); err != nil {
			log.Error("s.d.DelPKGCache(%s) error(%+v)", xstr.JoinInts(uids), err)
		}
		operAction := fmt.Sprintf("新增%s挂件 %d天", pi.Name, day)
		if _, err = s.d.AddPendantOperLog(context.Background(), oid, uids, pid, operAction); err != nil {
			log.Error("s.d.AddPendantOperLog(%d,%s,%s) error(%+v)", oid, xstr.JoinInts(uids), operAction, err)
		}
	})
	for _, uid := range uids {
		tid := uid
		s.addAsyn(func() {
			if err = s.accNotify(context.Background(), tid, model.AccountNotifyUpdatePendant); err != nil {
				log.Error("s.accNotify(%d) error(%+v)", tid, err)
			}
			if err = s.d.SetPendantPointCache(context.Background(), tid, pid); err != nil {
				log.Error("s.d.SetPendantPointCache(%d,%d) error(%+v)", tid, pid, err)
			}
		})
	}
	return
}

func (s *Service) diffSlice(sliceOne, sliceTwo []int64) (res []int64) {
	for _, ov := range sliceOne {
		inSlice := func(ov int64, sliceTwo []int64) bool {
			for _, tv := range sliceTwo {
				if tv == ov {
					return true
				}
			}
			return false
		}(ov, sliceTwo)
		if !inSlice {
			res = append(res, ov)
		}
	}
	return
}

// PendantOperlog pendant operactlog .
func (s *Service) PendantOperlog(c context.Context, pn, ps int) (opers []*model.PendantOperLog, pager *model.Pager, err error) {
	var total int64
	pager = &model.Pager{
		PN: pn,
		PS: ps,
	}
	if total, err = s.d.PendantOperationLogTotal(c); err != nil {
		err = errors.Wrap(err, "s.d.PendantOperationLogTotal()")
		return
	}
	if total <= 0 {
		return
	}
	pager.Total = total
	var uids []int64
	if opers, uids, err = s.d.PendantOperLog(c, pn, ps); err != nil {
		err = errors.Wrapf(err, "s.d.PendantOperLog(%d,%d)", pn, ps)
		return
	}
	var accInfoMap map[int64]*accmdl.Info
	if accInfoMap, err = s.fetchInfos(c, uids, _fetchInfoTimeout); err != nil {
		log.Error("service.fetchInfos(%v, %v) error(%v)", xstr.JoinInts(uids), _fetchInfoTimeout, err)
		err = nil
	}
	for _, v := range opers {
		if accInfo, ok := accInfoMap[v.UID]; ok {
			v.Action = fmt.Sprintf("给用户(%s) %s", accInfo.Name, v.Action)
		}
		if operName, ok := s.Managers[v.OID]; ok {
			v.OperName = operName
		}
	}
	return
}
