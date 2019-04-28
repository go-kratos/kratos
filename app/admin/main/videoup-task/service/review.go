package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

// ErrTaskMISS .
var ErrTaskMISS = fmt.Errorf("任务缓存查找失败")

// ReviewForm 复审表单
func (s *Service) ReviewForm(c context.Context, tid int64) (form *model.SubmitForm, err error) {
	return s.dao.ReviewForm(c, tid)
}

// CheckReview 检查任务复审
func (s *Service) CheckReview(c context.Context, form *model.SubmitForm) (isReview bool, err error) {
	var (
		v    *model.Video
		attr int32
		tx   *sql.Tx
		rows int64
		tp   *model.TaskPriority
	)

	attr, err = s.dao.VideoAttribute(c, form.CID)
	if err != nil {
		log.Error("CheckReview VideoAttribute(aid%d,cid%d) miss(%v)", form.AID, form.CID, err)
		return false, ErrTaskMISS
	}
	v = &model.Video{Attribute: attr}

	tp, err = s.getReviewParams(c, form)
	if err != nil || tp == nil {
		log.Info("CheckReview(%d) 不需要复审(%+v)", form.TaskID, tp)
		return false, err
	}

	s.SyncRC(c)
	ck := s.reviewCache.Check(c, tp, form.UID)
	if !ck {
		log.Info("CheckReview(%d) 不需要复审(%+v)", form.TaskID, tp)
		return false, nil
	}

	if _, err = s.dao.InReviewForm(c, form, form.UID, form.Uname); err != nil {
		return false, err
	}

	if tx, err = s.dao.BeginTran(c); err != nil {
		return false, err
	}

	if rows, err = s.dao.TxUpTaskByID(tx, form.TaskID, map[string]interface{}{"state": model.TypeReview}); err != nil {
		tx.Rollback()
		return false, err
	}
	if rows > 0 {
		if _, err = s.dao.TxAddTaskHis(tx, 0, model.ActionSubmit, form.TaskID, form.CID, form.UID, 0, form.Status, "TaskReview"); err != nil {
			tx.Rollback()
			return false, err
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return false, err
	}

	// log
	attrs := map[uint]int32{
		model.AttrBitNoRank:      form.Norank,
		model.AttrBitNoDynamic:   form.Noindex,
		model.AttrBitNoRecommend: form.NoRecommend,
		model.AttrBitNoSearch:    form.Nosearch,
		model.AttrBitOverseaLock: form.OverseaBlock,
		model.AttrBitPushBlog:    form.PushBlog,
	}

	var conts []string
	const template = "[%s]从[%s]设为[%s]"
	var yesOrNo = map[int32]string{model.AttrYes: "是", model.AttrNo: "否"}
	for bit, attr := range attrs {
		v.AttrSet(attr, bit)
		if attr == 1 {
			conts = append(conts, fmt.Sprintf(template, model.BitDesc(bit), yesOrNo[^attr&1], yesOrNo[attr]))
			log.Info("vid(%d) update video bit(%d) bitdesc(%s) attrs(%d)", form.ID, bit, model.BitDesc(bit), attr)
		}
	}

	vp := &model.VideoParam{
		ID:       form.ID,
		Aid:      form.AID,
		Mid:      form.MID,
		RegionID: tp.TypeID,
		Status:   form.Status,
		Cid:      form.CID,
		Title:    form.Eptitle,
		Desc:     form.Description,

		UID:       form.UID,
		TaskID:    form.TaskID,
		Oname:     form.Uname,
		TagID:     form.TID,
		Reason:    form.Reason,
		ReasonID:  form.ReasonID,
		Note:      form.Note,
		Attribute: v.Attribute,
	}
	oper := &model.VideoOper{Aid: form.AID, UID: form.UID, Vid: form.ID, Attribute: vp.Attribute, Status: form.Status, Remark: form.Note}
	operConts := append([]string{fmt.Sprintf("初审提交")}, conts...)
	operConts = append(operConts, s.diffVideoOper(vp)...)
	oper.Content = strings.Join(operConts, "，")
	s.addVideoOper(c, oper)

	s.sendVideoLog(c, vp, oper.Content)

	return true, nil
}

// ListReviewConfs 配置列表
func (s *Service) ListReviewConfs(c context.Context, unames, bt, et, sort string, pn, ps int64) (rcs []*model.ReviewConf, count int64, err error) {
	var uids []int64

	if len(unames) > 0 {
		res, _ := s.dao.Uids(c, strings.Split(unames, ","))
		for _, uid := range res {
			uids = append(uids, uid)
		}
	}

	rcs, count, err = s.dao.ListConfs(c, uids, bt, et, sort, pn, ps)
	for _, v := range rcs {
		if v.Bt.TimeValue().IsZero() {
			v.Bt = ""
		}
		if v.Et.TimeValue().IsZero() {
			v.Et = ""
		}

		if len(v.Uids) > 0 {
			if unames, _ := s.dao.Unames(c, v.Uids); len(unames) > 0 {
				for _, uname := range unames {
					v.Unames = append(v.Unames, uname)
				}
			}
		}
	}
	return
}

// AddReviewConf 添加配置
func (s *Service) AddReviewConf(c context.Context, rc *model.ReviewConf) (err error) {
	if len(rc.Types) > 0 {
		stypes, _ := xstr.SplitInts(s.tarnsType(c, xstr.JoinInts(rc.Types)))
		rc.Types = stypes
	}

	if _, err = s.dao.InReviewConf(c, rc); err != nil {
		log.Error("s.AddReviewConf(%+v) error(%v)", rc, err)
		return err
	}
	s.SyncRC(c)
	return
}

// EditReviewConf 修改配置
func (s *Service) EditReviewConf(c context.Context, rc *model.ReviewConf) (err error) {
	if len(rc.Types) > 0 {
		stypes, _ := xstr.SplitInts(s.tarnsType(c, xstr.JoinInts(rc.Types)))
		rc.Types = stypes
	}

	if _, err = s.dao.UpReviewConf(c, rc); err != nil {
		log.Error("s.EditReviewConf(%+v) error(%v)", rc, err)
		return err
	}
	s.SyncRC(c)
	return
}

// DelReviewConf 删除配置
func (s *Service) DelReviewConf(c context.Context, id int) (err error) {
	if _, err = s.dao.DelReviewConf(c, id); err != nil {
		log.Error("s.DelReviewConf(%d) error(%v)", id, err)
		return err
	}
	s.SyncRC(c)
	return
}

// SyncRC sync from db
func (s *Service) SyncRC(c context.Context) {
	rcs, err := s.dao.ReviewConfs(context.TODO())
	if err != nil {
		log.Error("loadRC error(%v)", err)
		return
	}

	if len(rcs) > 0 {
		s.reviewCache.Mux.Lock()
		defer s.reviewCache.Mux.Unlock()
		s.reviewCache.MRC = make(map[int64]*model.ReviewConf)
		for _, item := range rcs {
			s.reviewCache.MRC[item.ID] = item
		}
	}
}

func (s *Service) loadRC() {
	s.SyncRC(context.TODO())
}

func (s *Service) loadRCproc() {
	for {
		time.Sleep(3 * time.Minute)
		s.SyncRC(context.TODO())
	}
}

func (s *Service) getReviewParams(c context.Context, form *model.SubmitForm) (tp *model.TaskPriority, err error) {
	t, err := s.dao.TaskByID(c, form.TaskID)
	if err != nil || t == nil {
		return nil, ErrTaskMISS
	}
	if t.State != model.TypeDispatched {
		log.Info("CheckReview(%d) 不需要复审 state(%d)", form.TaskID, t.State)
		return nil, nil
	}

	mp, err := s.dao.GetWeightRedis(c, []int64{form.TaskID})
	if err != nil || len(mp) == 0 {
		if mp, err = s.dao.GetWeightDB(c, []int64{form.TaskID}); err != nil || len(mp) == 0 {
			log.Error("GetWeightDB(%d) miss", form.TaskID)
			return nil, ErrTaskMISS
		}
	}
	if _, ok := mp[form.TaskID]; !ok {
		log.Error("mp(%d) miss", form.TaskID)
		return nil, ErrTaskMISS
	}
	tp = mp[form.TaskID]

	// 补充复审判断的参数
	if tp.TypeID == 0 {
		s.setReviewParams(c, form.MID, form.AID, tp)
	}
	return
}

func (s *Service) setReviewParams(c context.Context, mid, aid int64, tp *model.TaskPriority) {
	typeid, upfrom, err := s.dao.ArchiveParam(c, aid)
	if err == nil {
		tp.TypeID = typeid
		tp.UpFrom = upfrom
	}
	tp.UpGroups = s.getSpecial(mid)
}
