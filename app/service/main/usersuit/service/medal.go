package service

import (
	"context"
	"fmt"
	"sort"

	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_MedalGetTitle   = "滋···叮！%s勋章合成完毕！"
	_MedalGetContext = "恭喜！你已成功领取“%s”勋章～想赶紧把它装备起来吗？请猛戳 ﾟ∀ﾟ)σ #{成就勋章}{\"https://account.bilibili.com/site/nameplate.html\"}"
)

var (
	medalLevel  = map[int32]string{model.Level1: "普通勋章", model.Level2: "高级勋章", model.Level3: "稀有勋章"}
	_emptyOwner = make([]*model.MedalOwner, 0)
)

// MedalInfo return medal info by mid and nid.
func (s *Service) MedalInfo(c context.Context, mid, nid int64) (res *model.MedalInfo, err error) {
	var (
		nos []*model.MedalOwner
	)
	if nos, err = s.OwnerMedal(c, mid); err != nil {
		log.Error("s.medalDao.MedalOwnersCache(%d) error(%v)", mid, err)
		err = nil
	}
	if _, ok := s.medalInfoAll[nid]; !ok {
		err = ecode.MedalNotFound
		return
	}
	res = s.medalInfoAll[nid]
	for _, v := range nos {
		if v.NID == nid {
			res.IsGet = model.IsGet
			break
		}
	}
	return
}

// MedalGet get medal.
func (s *Service) MedalGet(c context.Context, mid, nid int64) (err error) {
	var (
		nos []*model.MedalOwner
		mi  *model.MedalInfo
		ok  bool
	)
	mia := s.medalInfoAll
	if mi, ok = mia[nid]; !ok {
		err = ecode.MedalNotFound
		return
	}
	if nos, err = s.OwnerMedal(c, mid); err != nil {
		log.Error("s.medalDao.MedalOwnersCache(%d) error(%v)", mid, err)
		return
	}
	for _, no := range nos {
		if no.NID == nid {
			err = ecode.MedalHasGet
			return
		}
	}
	if err = s.medalDao.AddMedalOwner(c, mid, nid); err != nil {
		return
	}
	s.medalDao.DelMedalOwnersCache(c, mid)
	s.medalDao.SetPopupCache(c, mid, nid)
	s.medalDao.SetRedPointCache(c, mid, nid)
	s.medalDao.SendMsg(c, mid, fmt.Sprintf(_MedalGetTitle, mi.Name), fmt.Sprintf(_MedalGetContext, mi.Name))
	return
}

// MedalCheck check user is get medal.
func (s *Service) MedalCheck(c context.Context, mid, nid int64) (res *model.MedalCheck, err error) {
	res = &model.MedalCheck{}
	mos, _ := s.OwnerMedal(c, mid)
	if len(mos) != 0 {
		for _, mo := range mos {
			if mo.NID == nid {
				res.Has = model.IsGet
				res.Info = mo
				return
			}
		}
		res.Info = struct{}{}
		res.Has = model.NotGet
		return
	}
	mo, err := s.medalDao.OwnerBYNidMid(c, mid, nid)
	if err != nil {
		return
	}
	if mo == nil {
		res.Info = struct{}{}
		res.Has = model.NotGet
		return
	}
	res.Info = mo
	res.Has = model.IsGet
	return
}

// MedalHomeInfo return user mdeal home info.
func (s *Service) MedalHomeInfo(c context.Context, mid int64) (res []*model.MedalHomeInfo, err error) {
	var (
		nos          []*model.MedalOwner
		mia          = s.medalInfoAll
		activatedNid int64
	)
	res = make([]*model.MedalHomeInfo, 0)
	if nos, err = s.OwnerMedal(c, mid); err != nil {
		log.Error("s.medalDao.MedalOwnersCache(%d) error(%v)", mid, err)
		return
	}
	if activatedNid, err = s.ActivatedMedalID(c, mid); err != nil {
		log.Error("s.medalDao.MedalActivatedCache(%d) error(%v)", mid, err)
		return
	}
	if (len(nos)) == 0 {
		return
	}
	if len(nos) > 4 {
		nos = nos[0:4]
	}
	for _, no := range nos {
		mhi := &model.MedalHomeInfo{}
		mhi.NID = no.NID
		if _, ok := mia[no.NID]; ok {
			mhi.Name = mia[no.NID].Name
			mhi.Description = mia[no.NID].Description
			mhi.Image = mia[no.NID].Image
			mhi.Level = mia[no.NID].LevelDesc
			mhi.ImageSmall = mia[no.NID].ImageSmall
		}
		if activatedNid == no.NID {
			mhi.IsActivated = model.OwnerInstall
		}
		res = append(res, mhi)
	}
	return
}

// MedalUserInfo return medal user info.
func (s *Service) MedalUserInfo(c context.Context, mid int64, ip string) (res *model.MedalUserInfo, err error) {
	var (
		levelInfo *memmdl.LevelInfo
		memInfo   *memmdl.BaseInfo
		nid       int64
		trueLove  bool
	)
	if levelInfo, err = s.memRPC.Level(c, &memmdl.ArgMid2{Mid: mid, RealIP: ip}); err != nil || levelInfo == nil {
		log.Error("s.memRPC.Level(%+v) error(%v)", &memmdl.ArgMid2{Mid: mid, RealIP: ip}, err)
		return
	}
	if memInfo, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil || memInfo == nil {
		log.Error("s.memRPC.Base(%+v) error(%v)", &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}, err)
		return
	}
	if nid, err = s.ActivatedMedalID(c, mid); err != nil {
		log.Error("MedalUserInfo s.medalDao.MedalActivatedCache(%d) error(%v)", mid, err)
		err = nil
	}
	res = &model.MedalUserInfo{
		Name:  memInfo.Name,
		Face:  memInfo.Face,
		Level: levelInfo.Cur,
		NID:   nid,
	}
	if mi, ok := s.medalInfoAll[nid]; ok {
		res.ImageSmall = mi.ImageSmall
	}
	if trueLove, err = s.medalDao.GetWearedfansMedal(c, mid, 2); err != nil {
		log.Warn("s.medalDao.GetWearedfansMedal(%d), err(%+v)", mid, err)
	}
	res.TrueLove = trueLove
	return
}

// MedalInstall install or uninstall medal.
func (s *Service) MedalInstall(c context.Context, mid, nid int64, isActivated int8) (err error) {
	var (
		nos   []*model.MedalOwner
		isGet bool
	)
	if nos, err = s.OwnerMedal(c, mid); err != nil {
		log.Error("s.medalDao.MedalOwnersCache(%d) error(%v)", mid, err)
		return
	}
	for _, no := range nos {
		if no.NID == nid {
			isGet = true
			break
		}
	}
	if !isGet {
		err = ecode.MedalNotGet
		return
	}
	switch isActivated {
	case model.OwnerUninstall:
		if err = s.medalDao.UninstallMedalOwner(c, mid, nid); err != nil {
			err = errors.WithStack(err)
			return
		}
		s.medalDao.DelMedalActivatedCache(c, mid)

	case model.OwnerInstall:
		if err = s.medalDao.UninstallAllMedalOwner(c, mid, nid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if err = s.medalDao.InstallMedalOwner(c, mid, nid); err != nil {
			err = errors.WithStack(err)
			return
		}
		s.medalDao.SetMedalActivatedCache(c, mid, nid)
	default:
		log.Error("unkonw action(%d)", isActivated)
	}
	s.addNotify(func() {
		s.accNotify(context.Background(), mid, model.AccountNotifyUpdateMedal)
	})
	return
}

// MedalPopup return medal popup.
func (s *Service) MedalPopup(c context.Context, mid int64) (res *model.MedalPopup, err error) {
	nid, err := s.medalDao.PopupCache(c, mid)
	if err != nil {
		return
	}
	if nid <= 0 {
		return
	}
	mi, ok := s.medalInfoAll[nid]
	if !ok {
		return
	}
	res = &model.MedalPopup{
		NID:   mi.ID,
		Name:  mi.Name,
		Image: mi.Image,
	}
	s.medalDao.DelPopupCache(c, mid)
	return
}

// MedalMyInfo return medal my info.
func (s *Service) MedalMyInfo(c context.Context, mid int64) (res []*model.MedalMyInfos, err error) {
	var nos []*model.MedalOwner
	res = make([]*model.MedalMyInfos, 0)
	if nos, err = s.OwnerMedal(c, mid); err != nil {
		log.Error("s.medalDao.MedalOwnersCache(%d) error(%v)", mid, err)
		return
	}
	if len(nos) == 0 {
		return
	}
	mia := s.medalInfoAll
	activatedNid, _ := s.ActivatedMedalID(c, mid)
	popNid, _ := s.medalDao.PopupCache(c, mid)
	s.medalDao.DelPopupCache(c, mid)
	list1 := make([]*model.MedalMyInfo, 0)
	list2 := make([]*model.MedalMyInfo, 0)
	list3 := make([]*model.MedalMyInfo, 0)
	for _, no := range nos {
		mi, ok := mia[no.NID]
		if !ok {
			return
		}
		li := &model.MedalMyInfo{}
		if activatedNid == no.NID {
			li.IsActivated = 1
		}
		if popNid == no.NID {
			li.IsNewGet = 1
		}
		li.MedalInfo = mi
		li.GetTime = no.CTime
		switch mi.Level {
		case model.Level1:
			list1 = append(list1, li)
		case model.Level2:
			list2 = append(list2, li)
		case model.Level3:
			list3 = append(list3, li)
		}
	}
	for i := int32(1); i <= 3; i++ {
		re := &model.MedalMyInfos{}
		re.Name = medalLevel[i]
		switch i {
		case model.Level1:
			re.List = list1
		case model.Level2:
			re.List = list2
		case model.Level3:
			re.List = list3
		}
		if len(re.List) > 0 {
			re.Count = int32(len(re.List))
			res = append(res, re)
		}
	}
	return
}

// MedalAllInfo return medal all info.
func (s *Service) MedalAllInfo(c context.Context, mid int64) (res *model.MedalAllInfos, err error) {
	rnid, _ := s.medalDao.RedPointCache(c, mid)
	anid, _ := s.ActivatedMedalID(c, mid)
	res = &model.MedalAllInfos{
		List:         make(map[int64]*model.MedalCategoryInfo),
		RedPoint:     rnid > 0,
		HasActivated: anid,
	}
	var getMap = make(map[int64]bool)
	mos, _ := s.OwnerMedal(c, mid)
	for _, mo := range mos {
		res.HasGet = append(res.HasGet, mo.NID)
		getMap[mo.NID] = true
	}
	var allGids = make(map[int64]*model.MedalGroup)
	for i := 0; i < len(s.medalGroupAll); i++ {
		allGids[s.medalGroupAll[i].ID] = s.medalGroupAll[i]
	}
	var medals = make([][]*model.MedalInfo, len(s.medalInfoAll))
	for _, mi := range s.medalInfoAll {
		medals[mi.GID] = append(medals[mi.GID], mi)
	}
	var tmpPidNotZeros = make(map[int64][]*model.MedalItemInfo)
	var pids []int64
	for _, agi := range allGids {
		if agi.PID == 0 {
			pids = append(pids, agi.ID)
			continue
		}
		if medals[agi.ID] == nil {
			continue
		}
		ms := medals[agi.ID]
		sort.Slice(ms, func(i, j int) bool {
			return ms[i].Sort < ms[j].Sort
		})
		left := ms[len(ms)-1]
		tmpPidNotZero := &model.MedalItemInfo{
			Left:  left,
			Count: int32(len(ms)),
			Right: ms,
		}
		for i := 0; i < len(ms); i++ {
			if _, ok := getMap[ms[i].ID]; !ok {
				tmpPidNotZero.Left = ms[i]
				break
			}
		}
		tmpPidNotZeros[agi.PID] = append(tmpPidNotZeros[agi.PID], tmpPidNotZero)
	}
	for _, pid := range pids {
		var medalCategory = &model.MedalCategoryInfo{}
		if category, ok := tmpPidNotZeros[pid]; ok {
			medalCategory.Count = int32(len(category))
			medalCategory.Data = category
		} else {
			var pidIsZeros = make([]*model.MedalItemInfo, 0)
			for _, m := range medals[pid] {
				var pidIsZero = &model.MedalItemInfo{Left: m, Count: 1}
				pidIsZeros = append(pidIsZeros, pidIsZero)
			}
			medalCategory.Count = int32(len(pidIsZeros))
			medalCategory.Data = pidIsZeros
		}
		medalCategory.Name = allGids[pid].Name
		res.List[int64(allGids[pid].Rank)] = medalCategory
	}
	return
}

func (s *Service) loadMedal() {
	medalInfoAll, err := s.medalDao.MedalInfoAll(context.TODO())
	if err != nil {
		log.Error("s.medalDao.MedalInfoAll error(%v)", err)
		return
	}
	medalGroupAll, err := s.medalDao.MedalGroupAll(context.TODO())
	if err != nil {
		log.Error("s.medalDao.MedalGroupAll error(%v)", err)
		return
	}
	s.medalInfoAll = medalInfoAll
	s.medalGroupAll = medalGroupAll
}

// MedalActivated get user actived medal.
func (s *Service) MedalActivated(c context.Context, mid int64) (res *model.MedalInfo, err error) {
	var nid int64
	if nid, err = s.ActivatedMedalID(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	mia := s.medalInfoAll
	for id, m := range mia {
		if id == nid {
			res = m
			break
		}
	}
	return
}

// MedalActivatedMulti Multi get the user activated medal info(at most 50).
func (s *Service) MedalActivatedMulti(c context.Context, mids []int64) (res map[int64]*model.MedalInfo, err error) {
	if len(mids) > 50 {
		err = ecode.RequestErr
		return
	}
	res = make(map[int64]*model.MedalInfo)
	var nids = make(map[int64]int64)
	var miss = make([]int64, 0)
	if nids, miss, err = s.medalDao.MedalsActivatedCache(c, mids); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, mismid := range miss {
		mismid := mismid
		var (
			nid      int64
			canCache = true
		)
		if nid, err = s.medalDao.ActivatedOwnerByMid(c, mismid); err != nil {
			err = errors.WithStack(err)
			return
		}
		nids[mismid] = nid
		if canCache {
			s.addCache(func() {
				s.medalDao.SetMedalActivatedCache(context.Background(), mismid, nid)
			})
		}
	}
	for mid, nid := range nids {
		if mi, ok := s.medalInfoAll[nid]; ok {
			res[mid] = mi
		}
	}
	return
}

// ActivatedMedalID get user activated medal id.
func (s *Service) ActivatedMedalID(c context.Context, mid int64) (nid int64, err error) {
	var (
		notFound bool
		canCache = true
	)
	if nid, notFound, err = s.medalDao.MedalActivatedCache(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if notFound {
		if nid, err = s.medalDao.ActivatedOwnerByMid(c, mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if canCache {
			s.addCache(func() {
				s.medalDao.SetMedalActivatedCache(context.Background(), mid, nid)
			})
		}
	}
	return
}

// OwnerMedal get user Owner medal.
func (s *Service) OwnerMedal(c context.Context, mid int64) (res []*model.MedalOwner, err error) {
	var (
		notFound bool
		canCache = true
	)
	if res, notFound, err = s.medalDao.MedalOwnersCache(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if notFound {
		if res, err = s.medalDao.MedalOwnerByMid(c, mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(res) == 0 {
			res = _emptyOwner
		}
		if canCache {
			s.addCache(func() {
				s.medalDao.SetMedalOwnersache(context.Background(), mid, res)
			})
		}
	}
	return
}
