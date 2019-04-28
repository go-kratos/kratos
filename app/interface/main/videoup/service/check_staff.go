package service

import (
	"context"
	"fmt"
	"go-common/app/interface/main/videoup/model/archive"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"strings"
	"sync"
	"unicode/utf8"
)

// checkAddStaff 新增稿件：检查联合投稿参数
func (s *Service) checkAddStaff(c context.Context, ap *archive.ArcParam, mid int64, ip string) (err error) {
	var (
		_logs  []string
		_logs2 []string
	)
	defer func() {
		logStr := strings.Join(_logs, "\n")
		if err != nil {
			log.Error("s.checkAddStaff ap: %+v \nmid: %d \nlogs\n:%s", ap, mid, logStr)
		} else {
			log.Info("s.checkAddStaff ap: %+v \nmid: %d \nlogs\n%s", ap, mid, logStr)
		}
	}()
	ap.HandleStaff = false
	//新增稿件时，如果没填写Staffs，则不是联合投稿，跳过检查
	if len(ap.Staffs) == 0 {
		ap.Staffs = make([]*archive.Staff, 0)
		_logs = append(_logs, "INFO:非联合投稿，忽略")
		return
	}
	//如果不是自制稿件，不允许传Staffs
	if ap.Copyright != archive.CopyrightOriginal && len(ap.Staffs) != 0 {
		_logs = append(_logs, "ERR：非自制稿件，不允许传Staffs。")
		err = ecode.VideoupStaffCopyright
		return
	}
	//非全量的情况下，检查当前UP主是在白名单中
	_logs = append(_logs, fmt.Sprintf("INFO：灰度开关 %v。分区配置 %v 注意：这里打出的数据可能不是最新的，有可能后面重新Load了", s.staffGary, s.staffTypeCache))
	if err = s.checkStaffGray(ap.TypeID, mid); err != nil {
		_logs = append(_logs, fmt.Sprintf("INFO：UP主(%d)不在灰度名单中。", mid))
		return
	}
	_logs2, err = s.checkStaffData(c, ap.TypeID, ap.Staffs, []*archive.Staff{}, ip)
	_logs = append(_logs, _logs2...)
	if err != nil {
		return
	}
	//检查拉黑情况
	_logs2, _, err = s.checkStaffRelation(c, mid, ap.Staffs, ip)
	_logs = append(_logs, _logs2...)
	if err != nil {
		return
	}

	ap.HandleStaff = true
	_logs = append(_logs, "INFO:通过")
	return
}

// checkEditStaff 稿件编辑：检查联合投稿相关参数
func (s *Service) checkEditStaff(c context.Context, ap *archive.ArcParam, mid int64, a *archive.Archive, ip string) (err error) {
	var (
		cStaffs    []*archive.Staff //当前稿件的联合投稿人
		diffStaffs []*archive.Staff //变更的联合投稿人
		_logs      []string
		_logs2     []string
	)
	defer func() {
		logStr := strings.Join(_logs, "\n")
		if err != nil {
			log.Error("s.checkEditStaff ap: %+v \nmid: %d \narchive:%+v \nlogs:\n%v", ap, mid, a, logStr)
		} else {
			log.Info("s.checkEditStaff ap: %+v \nmid: %d \narchive:%+v \nlogs\n%v", ap, mid, a, logStr)
		}
	}()
	ap.HandleStaff = false
	//如果UP主不在白名单中，不修改Staff数据
	_logs = append(_logs, fmt.Sprintf("INFO：灰度开关 %v，分区配置 %v 注意：这里打出的数据可能不是最新的，有可能后面重新Load了", s.staffGary, s.staffTypeCache))
	if err = s.checkStaffGray(ap.TypeID, mid); err != nil {
		_logs = append(_logs, fmt.Sprintf("INFO:UP主(%d)不在灰度中", mid))
		//如果UP主不在白名单中，则将staff置空，不修改staff
		ap.Staffs = []*archive.Staff{}
		ap.HandleStaff = false
		err = nil
		return
	}
	//获取当前稿件的Staff信息
	if cStaffs, err = s.arc.ApplyStaffs(c, a.Aid, ip); err != nil {
		_logs = append(_logs, fmt.Sprintf("ERR:获取线上Staffs失败，err:%v", err))
		return
	}
	if len(ap.Staffs) == 0 { //删除Staffs
		_logs2, _, err = s.checkStaffRelation(c, mid, cStaffs, ip)
		_logs = append(_logs, _logs2...)
		if err != nil {
			return
		}
		_logs = append(_logs, "INFO:UP主删除所有Staff")
		ap.HandleStaff = true
		return
	}
	//preEdit()的时候已经加上了以下"移区"、"换类型"的判断逻辑，这里冗余一下
	//联合投稿不允许自制稿件改成非自制
	if a.Copyright == archive.CopyrightOriginal && ap.Copyright != archive.CopyrightOriginal && a.AttrVal(archive.AttrBitStaff) == archive.AttrYes {
		_logs = append(_logs, "ERROR:联合投稿不允许自制稿件改成非自制")
		err = ecode.VideoupStaffChangeCopyright
		return
	}
	_logs2, err = s.checkStaffData(c, ap.TypeID, ap.Staffs, cStaffs, ip)
	_logs = append(_logs, _logs2...)
	if err != nil {
		return
	}
	//只验证修改过的Staff信息
	changes, _logs2 := s.getStaffChanges(cStaffs, ap.Staffs)
	_logs = append(_logs, _logs2...)
	_logs = append(_logs, fmt.Sprintf("INFO:当前Staffs：%+v 提交Staffs：%+v 修改的Staffs：%+v", cStaffs, ap.Staffs, changes))
	for _, v := range changes {
		for _, m := range v {
			diffStaffs = append(diffStaffs, m)
		}
	}
	_logs2, _, err = s.checkStaffRelation(c, mid, diffStaffs, ip)
	_logs = append(_logs, _logs2...)
	if err != nil {
		return
	}
	ap.HandleStaff = true
	_logs = append(_logs, "INFO:通过")
	return
}

// checkStaffGray
func (s *Service) checkStaffGray(typeid int16, mid int64) (err error) {
	if s.staffGary {
		if _, ok := s.staffUps[mid]; !ok {
			log.Error("当前Up主(%d)不在联合投稿白名单中。", mid)
			err = ecode.VideoupStaffAuth
			return
		}
		var ok1, ok2 bool
		_, ok1 = s.staffTypeCache[typeid]
		_, ok2 = s.staffTypeCache[0]
		if ok1 && s.staffTypeCache[typeid].MaxStaff <= 0 {
			log.Error("当前分区(%d)在联合投稿黑名单中。", typeid)
			err = ecode.VideoupStaffTypeNotExists
			return
		}
		if !ok1 && !ok2 {
			log.Error("当前分区(%d)不在联合投稿白名单中。", typeid)
			err = ecode.VideoupStaffTypeNotExists
			return
		}
	}
	return
}

// checkStaffData 检查联合投稿人的数据格式（全量验证，不管UP主有没有编辑，都会走这个验证）
func (s *Service) checkStaffData(c context.Context, typeid int16, staffs, cStaffs []*archive.Staff, ip string) (_logs []string, err error) {
	var (
		titles              []string
		cards               map[int64]*accapi.Card
		staffMids           []int64
		staffMap, cStaffMap map[int64]string
		maxStaff            int
	)
	staffMap = make(map[int64]string)
	cStaffMap = make(map[int64]string)
	if len(cStaffs) != 0 {
		for _, v := range cStaffs {
			cStaffMap[v.Mid] = v.Title
		}
	}
	//检查分区的配置
	if conf, ok := s.staffTypeCache[typeid]; ok {
		maxStaff = conf.MaxStaff
	} else if conf, ok := s.staffTypeCache[0]; ok {
		maxStaff = conf.MaxStaff
	} else {
		_logs = append(_logs, fmt.Sprintf("ERR：分区%d不在配置里。配置：%v", typeid, s.staffTypeCache))
		err = ecode.VideoupStaffTypeNotExists
		return
	}
	if maxStaff == 0 {
		_logs = append(_logs, fmt.Sprintf("ERR：分区%d是黑名单。配置：%v", typeid, s.staffTypeCache))
		err = ecode.VideoupStaffTypeNotExists
		return
	}
	//检查Staff数量
	if len(staffs) > maxStaff {
		_logs = append(_logs, fmt.Sprintf("ERR：Staff数量超限。最多：%d；传递：%d", maxStaff, len(staffs)))
		err = ecode.VideoupStaffCountLimit
		return
	}
	//检查Staff职能、Mid
	for i, v := range staffs {
		staffs[i].Title = strings.TrimSpace(v.Title)
		v.Title = staffs[i].Title
		if v.Mid == 0 {
			_logs = append(_logs, "ERR：Staff Mid为0。")
			err = ecode.VideoupStaffMidInvalid
			return
		}
		tl := utf8.RuneCountInString(v.Title)
		if tl < 2 {
			_logs = append(_logs, fmt.Sprintf("ERR：职能(%v)长度不合法，长度：%d。", v.Title, tl))
			err = ecode.VideoupStaffTitleShort
			return
		}
		if tl > 4 {
			_logs = append(_logs, fmt.Sprintf("ERR：职能(%v)长度不合法，长度：%d。", v.Title, tl))
			err = ecode.VideoupStaffTitleLength
			return
		}
		if !_staffNameReg.MatchString(v.Title) {
			_logs = append(_logs, fmt.Sprintf("ERR：职能(%v)字符不合法。", v.Title))
			err = ecode.VideoupStaffTitleChar
			return
		}
		//不校验未修改的职能
		if cTitle, ok := cStaffMap[v.Mid]; !ok || cTitle != v.Title {
			titles = append(titles, v.Title)
		}
		staffMap[v.Mid] = v.Title
		staffMids = append(staffMids, v.Mid)
	}
	if len(staffMap) != len(staffs) {
		_logs = append(_logs, fmt.Sprintf("ERR：Staff存在重复。传递：%v；去重后：%v", staffs, staffMap))
		err = ecode.VideoupStaffMidRepeat
		return
	}
	//职能名称敏感词
	_, hit, err := s.filter.VideoMultiFilter(c, titles, ip)
	if err != nil {
		_logs = append(_logs, fmt.Sprintf("ERR：职能敏感词接口失败。error：%v", err))
		return
	}
	if len(hit) > 0 {
		_logs = append(_logs, fmt.Sprintf("ERR：职能存在敏感词。敏感词：%v", hit))
		err = ecode.VideoupStaffTitleFilter
		return
	}
	//Staff Mid合法性检查
	if cards, err = s.acc.Cards(c, staffMids, ip); err != nil {
		_logs = append(_logs, fmt.Sprintf("ERR：Staff账号信息获取失败。error：%v", err))
		return
	}
	for _, v := range staffMids {
		if _, ok := cards[v]; !ok {
			_logs = append(_logs, fmt.Sprintf("ERR：Staff Mid(%d)不存在。", v))
			err = ecode.VideoupStaffMidInvalid
			return
		}
	}
	return
}

// checkStaffRelation 检查联合投稿人拉黑关系
func (s *Service) checkStaffRelation(c context.Context, mid int64, staffs []*archive.Staff, ip string) (_logs []string, blocked []*archive.Staff, err error) {
	var (
		mids     []int64
		rels     map[int64]int //relations
		pass     = true
		staffMap map[int64]*archive.Staff
		cards    map[int64]*accapi.Card
	)
	blocked = make([]*archive.Staff, 0)
	staffMap = make(map[int64]*archive.Staff)
	for _, v := range staffs {
		mids = append(mids, v.Mid)
		staffMap[v.Mid] = v
	}
	if rels, err = s.FRelations(c, mid, mids, ip); err != nil {
		_logs = append(_logs, fmt.Sprintf("ERR：获取拉黑信息失败，error:%v", err))
		return
	}
	for k, v := range rels {
		if v >= 128 {
			_logs = append(_logs, fmt.Sprintf("ERR：Staff(%d)在黑名单中", k))
			pass = false
			blocked = append(blocked, staffMap[k])
			continue
		}
	}
	cards, err = s.acc.Cards(c, mids, ip)
	if !pass {
		if err != nil {
			_logs = append(_logs, fmt.Sprintf("ERR：账号(%v)信息获取失败 error:%v", mids, err))
			err = ecode.Errorf(ecode.VideoupStaffBlocked, ecode.VideoupStaffBlocked.Message(), "")
			return
		}
		var bNames []string
		for _, v := range blocked {
			if _, ok := cards[v.Mid]; ok {
				bNames = append(bNames, cards[v.Mid].Name)
			} else {
				_logs = append(_logs, fmt.Sprintf("ERR：账号(%d)信息获取失败 error:%v", v.Mid, err))
			}
		}
		err = ecode.Errorf(ecode.VideoupStaffBlocked, ecode.VideoupStaffBlocked.Message(), strings.Join(bNames, "、"))
	}
	for _, staff := range staffs {
		if _, ok := cards[staff.Mid]; !ok {
			_logs = append(_logs, fmt.Sprintf("ERR：账号(%d)信息获取失败 error:%v", staff.Mid, err))
			err = ecode.Errorf(ecode.CreativeAccServiceErr, "参与者(%d)信息获取失败", staff.Mid)
			return
		}
		if cards[staff.Mid].Silence == 1 {
			_logs = append(_logs, fmt.Sprintf("ERR：Staff Mid(%d)被封禁。", staff.Mid))
			err = ecode.Errorf(ecode.VideoupStaffUpSilence, ecode.VideoupStaffUpSilence.Message(), cards[staff.Mid].Name)
			return
		}
	}

	return
}

// getStaffChanges 获取联合投稿人的变更
func (s *Service) getStaffChanges(oS, nS []*archive.Staff) (changes map[string]map[int64]*archive.Staff, _logs []string) {
	var (
		allS = make([]*archive.Staff, 0)
		oMap = make(map[int64]*archive.Staff)
		nMap = make(map[int64]*archive.Staff)
	)
	str := ""
	for _, v := range oS {
		str += fmt.Sprintf("%d %s;", v.Mid, v.Title)
		oMap[v.Mid] = v
	}
	_logs = append(_logs, " 原Staffs："+str)
	str = ""
	for _, v := range nS {
		str += fmt.Sprintf("%d %s;", v.Mid, v.Title)
		nMap[v.Mid] = v
	}
	_logs = append(_logs, " 提交Staffs："+str)
	changes = make(map[string]map[int64]*archive.Staff)
	changes["add"] = make(map[int64]*archive.Staff)
	changes["edit"] = make(map[int64]*archive.Staff)
	changes["del"] = make(map[int64]*archive.Staff)
	allS = append(allS, oS...)
	allS = append(allS, nS...)
	for _, v := range allS {
		if _, ok := oMap[v.Mid]; !ok {
			changes["add"][v.Mid] = v
		} else if _, ok := nMap[v.Mid]; !ok {
			changes["del"][v.Mid] = v
		} else if oMap[v.Mid].Title != nMap[v.Mid].Title {
			changes["edit"][v.Mid] = v
		}
	}
	return
}

// FRelations 获取用户与mid的关系（Relations的反向）
func (s *Service) FRelations(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	var (
		g, ctx = errgroup.WithContext(c)
		sm     sync.RWMutex
	)

	res = make(map[int64]int)
	for _, v := range fids {
		g.Go(func() error {
			var r map[int64]int
			if r, err = s.acc.Relations(ctx, v, []int64{mid}, ip); err != nil {
				return err
			}
			sm.Lock()
			res[v] = r[mid]
			sm.Unlock()
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("s.FRelations(%d,%v) error(%v)", mid, fids, err)
	}
	return
}

// checkStaffMoveType 联合投稿不允许移区和转载类型（有申请的Staff时，都不允许移区）
func (s *Service) checkStaffMoveType(c context.Context, ap *archive.ArcParam, a *archive.Archive, ip string) (err error) {
	var (
		_logs []string
	)
	defer func() {
		if err != nil {
			log.Error("s.checkStaffMoveType ap: %+v, archive:%+v logs(%v)", ap, a, _logs)
		} else {
			log.Info("s.checkStaffMoveType ap: %+v, archive:%+v logs(%v)", ap, a, _logs)
		}
	}()
	//如果没发生移区和修改转载类型，则直接通过
	if ap.TypeID == a.TypeID && a.Copyright == ap.Copyright {
		_logs = append(_logs, "INFO:没有修改分区和转载类型")
		return
	}
	var (
		cStaffs []*archive.Staff //当前稿件的联合投稿人
	)
	if cStaffs, err = s.arc.ApplyStaffs(c, a.Aid, ip); err != nil {
		_logs = append(_logs, fmt.Sprintf("ERR:获取Staff失败。error:%v", err))
		log.Error("checkStaffMoveType() 获取线上Staffs失败，err:%v", err)
		return
	}
	if len(cStaffs) != 0 {
		_logs = append(_logs, fmt.Sprintf("ERR: 不允许操作。当前Staffs(%v)", cStaffs))
		err = ecode.VideoupStaffChangeTypeCopyright
		return
	}
	_logs = append(_logs, "INFO:通过")
	return
}
