package archive

import (
	"context"
	"go-common/app/interface/main/creative/model/tag"
	"go-common/library/ecode"
	"go-common/library/log"
)

// TagCheck fn
func (s *Service) TagCheck(c context.Context, mid int64, tagName string) (code int, msg string) {
	var (
		err error
		t   *tag.Tag
	)
	for _, act := range s.p.ActVideoAllCache {
		if act.Tags == tagName {
			code = 1
			msg = "当前tag为活动专用，不允许自定义添加，请先删除之后重选对应的活动"
			return
		}
	}
	if t, err = s.tag.TagCheck(c, mid, tagName); err != nil {
		log.Error("s.tag.TagCheck mid(%d)|tagName(%s)|err(%v)", mid, tagName, err)
		code = 1
		msg = ecode.String(err.Error()).Message()
		return
	}
	if t != nil && (t.State == tag.TagStateDel ||
		t.State == tag.TagStateHide ||
		t.Type == tag.OfficailActiveTag) {
		code = 1
		msg = "非法Tag不允许Up主使用"
		return
	}

	return
}
