package model

import (
	"context"
	"sync"
	"time"

	"go-common/library/log"
)

// ReviewConf 复审配置
type ReviewConf struct {
	ID       int64   `json:"id" form:"id"`
	Types    []int64 `json:"types" form:"types,split"` // 分区
	mtp      map[int16]struct{}
	UpFroms  []int64 `json:"upfroms" form:"upfroms,split"` // 投稿来源
	muf      map[int8]struct{}
	UpGroups []int64 `json:"upgroups" form:"upgroups,split"` // 用户组
	mug      map[int8]struct{}
	Uids     []int64  `json:"uids" form:"uids,split"` // 指定uid
	Unames   []string `json:"unames"`                 // 指定uid
	muid     map[int64]struct{}
	FansLow  int64      `json:"fanslow" form:"fanslow"`   // 粉丝数最低值
	FansHigh int64      `json:"fanshigh" form:"fanshigh"` // 粉丝数最高
	Bt       FormatTime `json:"bt" form:"bt"`
	Et       FormatTime `json:"et" form:"et"`
	State    int8       `json:"state" form:"state"`
	UID      int64      `json:"uid"`
	Uname    string     `json:"uname"`
	Desc     string     `json:"desc" form:"desc"`
	Mt       FormatTime `json:"mt"`
}

// Refresh refresh
func (r *ReviewConf) Refresh() {
	mtp := make(map[int16]struct{})
	muf := make(map[int8]struct{})
	mug := make(map[int8]struct{})
	muid := make(map[int64]struct{})

	for _, tp := range r.Types {
		mtp[int16(tp)] = struct{}{}
	}
	for _, uf := range r.UpFroms {
		muf[int8(uf)] = struct{}{}
	}
	for _, ug := range r.UpGroups {
		mug[int8(ug)] = struct{}{}
	}
	for _, uid := range r.Uids {
		muid[uid] = struct{}{}
	}

	r.mtp = mtp
	r.muf = muf
	r.mug = mug
	r.muid = muid
}

// SubmitForm  form
type SubmitForm struct {
	Status       int16  `json:"status" form:"status"`
	ID           int64  `json:"id" form:"id"`
	CID          int64  `json:"cid" form:"cid"`
	AID          int64  `json:"aid" form:"aid"`
	MID          int64  `json:"mid" form:"mid"`
	Eptitle      string `json:"eptitle,omitempty" form:"eptitle"`
	Description  string `json:"description,omitempty" form:"description"`
	Note         string `json:"note,omitempty" form:"note"`
	ReasonID     int64  `json:"reason_id,omitempty" form:"reason_id"`
	Reason       string `json:"reason,omitempty" form:"reason"`
	TID          int64  `json:"tid,omitempty" form:"tid"`
	Norank       int32  `json:"norank" form:"norank"`
	Noindex      int32  `json:"noindex" form:"noindex"`
	PushBlog     int32  `json:"push_blog" form:"push_blog"`
	NoRecommend  int32  `json:"norecommend" form:"norecommend"`
	Nosearch     int32  `json:"nosearch" form:"nosearch"`
	OverseaBlock int32  `json:"oversea_block" form:"oversea_block"`
	Encoding     int8   `json:"encoding" form:"encoding"`
	TaskID       int64  `json:"task_id" form:"task_id"`
	UID          int64  `json:"uid" form:"uid"`
	Uname        string `json:"uname" form:"uname"`
}

// ReviewCache 快速判断配置项是否命中
type ReviewCache struct {
	MRC map[int64]*ReviewConf
	Mux sync.RWMutex
}

// NewRC 复审配置
func NewRC() *ReviewCache {
	rc := &ReviewCache{}
	rc.MRC = make(map[int64]*ReviewConf)
	return rc
}

// Check 检查配置是否命中
func (rc *ReviewCache) Check(c context.Context, opt *TaskPriority, uid int64) bool {
	rc.Mux.RLock()
	defer rc.Mux.RUnlock()

	if len(rc.MRC) == 0 {
		log.Info("ReviewCache empty")
		return false
	}

	log.Info("ReviewCache opt(%+v) uid(%d),", opt, uid)
	for id, item := range rc.MRC {
		log.Info("ReviewCache config(%+v)", item)
		if item.State != 0 {
			continue
		}

		bt := item.Bt.TimeValue()
		et := item.Et.TimeValue()
		if bt.After(time.Now()) || (!et.IsZero() && et.Before(time.Now())) {
			continue
		}

		if len(item.mtp) > 0 {
			if _, ok := item.mtp[opt.TypeID]; !ok {
				continue
			}
		}

		if len(item.muf) > 0 {
			if _, ok := item.muf[opt.UpFrom]; !ok {
				continue
			}
		}

		if len(item.mug) > 0 {
			var hit bool
			for _, ug := range opt.UpGroups {
				if _, ok := item.mug[ug]; ok {
					hit = true
					break
				}
			}
			if !hit {
				continue
			}
		}

		if len(item.muid) > 0 {
			if _, ok := item.muid[uid]; !ok {
				continue
			}
		}

		if item.FansHigh > 0 {
			if opt.Fans < item.FansLow || opt.Fans > item.FansHigh {
				continue
			}
		}

		log.Info("ReviewCache task(%d) hit config(%d)", opt.TaskID, id)
		return true
	}

	return false
}
