package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"go-common/app/admin/main/usersuit/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// Medal medal .
func (s *Service) Medal(c context.Context) (res []*model.MedalInfo, err error) {
	nps, err := s.d.Medal(c)
	if err != nil {
		log.Error("s.d.Medal error(%v)", err)
	}
	mgs, err := s.d.MedalGroup(c)
	if err != nil {
		log.Error("s.d.MedalGroup error(%v)", err)
	}
	res = make([]*model.MedalInfo, 0)
	for _, np := range nps {
		re := &model.MedalInfo{}
		re.Medal = np
		if _, ok := mgs[np.GID]; ok {
			re.GroupName = mgs[np.GID].Name
			if mgs[np.GID].PID != 0 {
				re.ParentGroupName = mgs[mgs[np.GID].PID].Name
			}
		}
		res = append(res, re)
	}
	return
}

// MedalView .
func (s *Service) MedalView(c context.Context, id int64) (res *model.MedalInfo, err error) {
	res = &model.MedalInfo{}
	res.Medal, err = s.d.MedalByID(c, id)
	if err != nil {
		log.Error("s.d.MedalByID(%d) error(%v)", id, err)
		return
	}
	mgs, err := s.d.MedalGroup(c)
	if err != nil {
		log.Error("s.d.MedalGroup error(%v)", err)
		return
	}
	if _, ok := mgs[res.Medal.GID]; ok {
		res.GroupName = mgs[res.Medal.GID].Name
		if _, ok1 := mgs[mgs[res.Medal.GID].PID]; ok && ok1 {
			res.ParentGroupName = mgs[mgs[res.Medal.GID].PID].Name
		}
	}
	return
}

func getImagePath(raw string) string {
	uri, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return uri.Path
}

// AddMedal add medal .
func (s *Service) AddMedal(c context.Context, np *model.Medal) (err error) {
	np.Image = getImagePath(np.Image)
	np.ImageSmall = getImagePath(np.ImageSmall)
	if _, err := s.d.AddMedal(c, np); err != nil {
		log.Error("s.d.AddMedal error(%v)", err)
	}
	return
}

// UpMedal update medal .
func (s *Service) UpMedal(c context.Context, id int64, np *model.Medal) (err error) {
	np.Image = getImagePath(np.Image)
	np.ImageSmall = getImagePath(np.ImageSmall)
	if _, err := s.d.UpMedal(c, id, np); err != nil {
		log.Error("s.d.UpMedal error(%v)", err)
	}
	return
}

// MedalGroup medal group .
func (s *Service) MedalGroup(c context.Context) (res map[int64]*model.MedalGroup, err error) {
	res, err = s.d.MedalGroup(c)
	if err != nil {
		log.Error("s.MedalGroup error(%v)", err)
	}
	return
}

// MedalGroupInfo medal group all info include parent group name .
func (s *Service) MedalGroupInfo(c context.Context) (res []*model.MedalGroup, err error) {
	res, err = s.d.MedalGroupInfo(c)
	if err != nil {
		log.Error("s.MedalGroupInfo error(%v)", err)
	}
	return
}

// MedalGroupParent medal group all info include parent group name .
func (s *Service) MedalGroupParent(c context.Context) (res []*model.MedalGroup, err error) {
	res, err = s.d.MedalGroupParent(c)
	if err != nil {
		log.Error("s.MedalGroupParent error(%v)", err)
	}
	return
}

// MedalGroupByGid nameplate by gid .
func (s *Service) MedalGroupByGid(c context.Context, id int64) (ng *model.MedalGroup, err error) {
	if ng, err = s.d.MedalGroupByID(c, id); err != nil {
		log.Error("s.MedalGroupByID error(%v)", err)
	}
	return
}

// MedalGroupAdd add medal group.
func (s *Service) MedalGroupAdd(c context.Context, ng *model.MedalGroup) (err error) {
	if _, err := s.d.MedalGroupAdd(c, ng); err != nil {
		log.Error("s.MedalGroupAdd error(%v)", err)
	}
	return
}

// MedalGroupUp update medal group.
func (s *Service) MedalGroupUp(c context.Context, id int64, ng *model.MedalGroup) (err error) {
	if _, err := s.d.MedalGroupUp(c, id, ng); err != nil {
		log.Error("s.MedalGroupUp error(%v)", err)
	}
	return
}

// MedalOwner medal onwer .
func (s *Service) MedalOwner(c context.Context, mid int64) (no []*model.MedalMemberMID, err error) {
	if no, err = s.d.MedalOwner(c, mid); err != nil {
		log.Error("s.d.MedalOwner error(%+v)", err)
	}
	return
}

// MedalOwnerAdd medal owner add .
func (s *Service) MedalOwnerAdd(c context.Context, mid, nid int64, title, msg string, oid int64) (err error) {
	count, err := s.d.CountOwnerBYNidMid(c, mid, nid)
	if count > 0 || err != nil {
		err = ecode.MedalHasGet
		return
	}
	if _, err = s.d.MedalOwnerAdd(c, mid, nid); err != nil {
		log.Error("s.MedalOwnerAdd(mid:%d nid:%d) error(%v)", mid, nid, err)
		return
	}
	if err = s.d.DelMedalOwnersCache(c, mid); err != nil {
		log.Error("s.DelMedalOwnersCache(mid:%d) error(%v)", mid, err)
		err = nil
	}
	var (
		mids   []int64
		ismsg  bool
		action string
	)
	mids = append(mids, mid)
	log.Error("MedalOwnerAdd title(%+v) msg(%+v)", title, msg)
	if title != "" && msg != "" {
		ismsg = true
		if err = s.d.SendSysMsg(c, mids, title, msg, ""); err != nil {
			log.Error("MedalOwnerAdd(mid:%d nid:%d title:%s msg:%s) SendSysMsg error(%+v)", mid, nid, title, msg, err)
			err = nil
		}
	}
	if oid > 0 {
		mi, err := s.d.MedalByID(c, nid)
		if err != nil {
			log.Error("MedalByID(id:%d) error(%+v)", nid, err)
			err = nil
		}
		action = fmt.Sprintf("激活勋章:%s,", mi.Name)
		log.Error("MedalOwnerAdd ismsg(id:%d %+v) error(%+v)", nid, ismsg, err)
		if ismsg {
			action += fmt.Sprintf("并发送消息title:%s msg:%s", title, msg)
		} else {
			action += fmt.Sprintf("并没有发送消息")
		}
		s.d.AddMedalOperLog(c, oid, mid, nid, action)
	}
	return
}

// MedalOwnerAddList .
func (s *Service) MedalOwnerAddList(c context.Context, mid int64) (res []*model.MedalMemberAddList, err error) {
	if res, err = s.d.MedalAddList(c, mid); err != nil {
		log.Error("s.d.MedalOwnerAddList(%d) error(%v)", mid, err)
	}
	return
}

// MedalOwnerUpActivated update medal owner is_activated.
func (s *Service) MedalOwnerUpActivated(c context.Context, mid, nid int64) (err error) {
	if _, err = s.d.MedalOwnerUpActivated(c, mid, nid); err != nil {
		log.Error("s.d.UpMedalOwnerActivated(mid:%d nid:%d) error(%v)", mid, nid, err)
	}
	if _, err = s.d.MedalOwnerUpNotActivated(c, mid, nid); err != nil {
		log.Error("s.d.UpMedalOwnerNotActivated(mid:%d nid:%d) error(%v)", mid, nid, err)
	}
	if err = s.d.SetMedalActivatedCache(c, mid, nid); err != nil {
		log.Error("s.d.DelMedalActivatedCache(mid:%d nid:%d) error(%v)", mid, nid, err)
		err = nil
	}
	s.addAsyn(func() {
		if err = s.accNotify(context.Background(), mid, model.AccountNotifyUpdateMedal); err != nil {
			log.Error("s.accNotify(%d) error(%+v)", mid, err)
			return
		}
	})
	return
}

// MedalOwnerDel update medal owner is_del .
func (s *Service) MedalOwnerDel(c context.Context, mid, nid int64, isDel int8, title, msg string) (err error) {
	if _, err = s.d.MedalOwnerDel(c, mid, nid, isDel); err != nil {
		log.Error("s.d.MedalOwnerDel error(%v)", err)
	}
	if err = s.d.DelMedalOwnersCache(c, mid); err != nil {
		log.Error("s.d.DelMedalOwnersCache(%d) error(%v)", mid, err)
		err = nil
	}
	var mids []int64
	mids = append(mids, mid)
	if title != "" && msg != "" {
		if err = s.d.SendSysMsg(c, mids, title, msg, ""); err != nil {
			log.Error("MedalOwnerDel(mid:%d nid:%d title:%s msg:%s) SendSysMsg error(%+v)", mid, nid, title, msg, err)
			err = nil
		}
	}
	return
}

// ReadCsv read csv file
func (s *Service) ReadCsv(f multipart.File, h *multipart.FileHeader) (rs [][]string, err error) {
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			log.Error("upload question ReadCsv error(%v)", err)
		}
		if len(record) == 1 {
			rs = append(rs, record)
		}
	}
	return
}

// BatchAdd medal bacth add.
func (s *Service) BatchAdd(c context.Context, nid int64, f multipart.File, h *multipart.FileHeader) (msg string, err error) {
	if h != nil && !strings.HasSuffix(h.Filename, ".csv") {
		msg = "not csv file."
		return
	}
	rs, err := s.ReadCsv(f, h)
	if len(rs) == 0 || len(rs) > model.MaxCount {
		msg = "file size count is 0 or more than " + strconv.FormatInt(model.MaxCount, 10)
		return
	}
	for _, r := range rs {
		mid, err := strconv.ParseInt(r[0], 10, 64)
		if err == nil {
			if err = s.MedalOwnerAdd(c, mid, nid, "", "", 0); err != nil {
				log.Error("s.d.MedalOwnerAdd(mid:%d nid:%d) error(%v)", mid, nid, err)
			}
		}

	}
	return
}

// MedalOperlog medal operactlog .
func (s *Service) MedalOperlog(c context.Context, mid int64, pn, ps int) (opers []*model.MedalOperLog, pager *model.Pager, err error) {
	var total int64
	pager = &model.Pager{
		PN: pn,
		PS: ps,
	}
	if total, err = s.d.MedalOperationLogTotal(c, mid); err != nil {
		err = errors.Wrap(err, "s.d.MedalOperationLogTotal()")
		return
	}
	if total <= 0 {
		return
	}
	pager.Total = total
	var uids []int64
	if opers, uids, err = s.d.MedalOperLog(c, mid, pn, ps); err != nil {
		err = errors.Wrapf(err, "s.d.MedalOperLog(%d,%d,%d)", mid, pn, ps)
		return
	}
	var accInfoMap map[int64]*accmdl.Info
	if accInfoMap, err = s.fetchInfos(c, uids, _fetchInfoTimeout); err != nil {
		log.Error("service.fetchInfos(%v, %v) error(%v)", xstr.JoinInts(uids), _fetchInfoTimeout, err)
		err = nil
	}
	for _, v := range opers {
		if accInfo, ok := accInfoMap[v.MID]; ok {
			v.Action = fmt.Sprintf("给用户(%s) %s", accInfo.Name, v.Action)
		}
		if operName, ok := s.Managers[v.OID]; ok {
			v.OperName = operName
		}
	}
	return
}
