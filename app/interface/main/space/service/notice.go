package service

import (
	"context"
	"encoding/json"
	"html/template"
	"strings"

	"go-common/app/interface/main/space/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const _noticeTable = "member_up_notice"

// Notice get notice.
func (s *Service) Notice(c context.Context, mid int64) (res string, err error) {
	if _, ok := s.noNoticeMids[mid]; ok {
		return
	}
	var notice *model.Notice
	if notice, err = s.dao.Notice(c, mid); err != nil {
		return
	}
	if notice.IsForbid == _noticeForbid {
		notice.Notice = ""
	}
	res = template.HTMLEscapeString(notice.Notice)
	return
}

// SetNotice set notice.
func (s *Service) SetNotice(c context.Context, mid int64, notice string) (err error) {
	var (
		info    *accmdl.Profile
		preData *model.Notice
	)
	if info, err = s.realName(c, mid); err != nil {
		return
	}
	if info.Silence == _silenceForbid {
		err = ecode.UserDisabled
		return
	}
	if preData, err = s.dao.Notice(c, mid); err != nil {
		return
	}
	if notice == preData.Notice {
		err = ecode.NotModified
		return
	}
	if err = s.dao.SetNotice(c, mid, notice); err != nil {
		log.Error("s.dao.SetNotice(%d,%s) error(%v)", mid, notice, err)
		return
	}
	s.cache.Do(c, func(c context.Context) {
		s.dao.AddCacheNotice(c, mid, &model.Notice{Notice: notice})
	})
	return
}

// ClearCache del match and object cache
func (s *Service) ClearCache(c context.Context, msg string) (err error) {
	var m struct {
		Table string `json:"table"`
		Old   struct {
			Mid      int64  `json:"mid"`
			Notice   string `json:"notice"`
			IsForbid int    `json:"is_forbid"`
		} `json:"old,omitempty"`
		New struct {
			Mid      int64  `json:"mid"`
			Notice   string `json:"notice"`
			IsForbid int    `json:"is_forbid"`
		} `json:"new,omitempty"`
	}
	if err = json.Unmarshal([]byte(msg), &m); err != nil || m.Table == "" {
		log.Error("ClearCache json.Unmarshal msg(%s) error(%v)", msg, err)
		return
	}
	log.Info("ClearCache json.Unmarshal msg(%s)", msg)
	if strings.HasPrefix(m.Table, _noticeTable) && (m.Old.IsForbid != m.New.IsForbid || m.Old.Notice != m.New.Notice) {
		if err = s.dao.DelCacheNotice(c, m.New.Mid); err != nil {
			log.Error("s.dao.DelCacheNotice mid(%d) error(%v)", m.New.Mid, err)
		}
	}
	return
}
