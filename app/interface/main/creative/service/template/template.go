package template

import (
	"context"
	"go-common/library/ecode"
	"time"

	"go-common/app/interface/main/creative/model/template"
	"go-common/library/log"
	timex "go-common/library/time"
)

// Templates get templates for archive.
func (s *Service) Templates(c context.Context, mid int64) (tps []*template.Template, err error) {
	if tps, err = s.tpl.Templates(c, mid); err != nil {
		log.Error("s.tem.Templates(%d) error(%v)", mid, err)
		return
	}
	return
}

// AddTemplate add template for archive.
func (s *Service) AddTemplate(c context.Context, mid int64, typeid int16, cp, name, title, tag, content string, now time.Time) (err error) {
	var (
		tpl   *template.Template
		count int64
	)
	if count, err = s.tpl.Count(c, mid); err != nil {
		log.Error("s.tpl.Count(%d) error(%v)", count, err)
		return
	}
	if count >= 5 {
		log.Error("mid(%d) upper limit(%d)", mid, count)
		err = ecode.CreativeTemplateOverMax
		return
	}
	tpl = &template.Template{
		Name:      name,
		Title:     title,
		Tag:       tag,
		Content:   content,
		TypeID:    typeid,
		Copyright: template.Copyright(cp),
		State:     int8(template.StateNormal),
		CTime:     timex.Time(now.Unix()),
		MTime:     timex.Time(now.Unix()),
	}
	_, err = s.tpl.AddTemplate(c, mid, tpl)
	return
}

// UpdateTemplate update template for archive.
func (s *Service) UpdateTemplate(c context.Context, id, mid int64, typeid int16, cp, name, title, tag, content string, now time.Time) (err error) {
	var (
		t, tpl *template.Template
	)
	if t, err = s.tpl.Template(c, id, mid); err != nil {
		log.Error("s.tpl.Template id(%d)  mid(%d) error(%v)", id, mid, err)
		return
	}
	if t == nil {
		err = ecode.NothingFound
		return
	}
	if t.State != 0 {
		err = ecode.CreativeTemplateDeleted
		return
	}
	tpl = &template.Template{
		ID:        id,
		Name:      name,
		Title:     title,
		Tag:       tag,
		Content:   content,
		TypeID:    typeid,
		Copyright: template.Copyright(cp),
		MTime:     timex.Time(now.Unix()),
	}
	_, err = s.tpl.UpTemplate(c, mid, tpl)
	return
}

// DelTemplate delete template.
func (s *Service) DelTemplate(c context.Context, id, mid int64, now time.Time) (err error) {
	var (
		t, tpl *template.Template
	)
	if t, err = s.tpl.Template(c, id, mid); err != nil {
		log.Error("s.tpl.Template id(%d) mid(%d) error(%v)", id, mid, err)
		return
	}
	if t == nil {
		err = ecode.NothingFound
		return
	}
	if t.State != 0 {
		err = ecode.CreativeTemplateDeleted
		return
	}
	tpl = &template.Template{
		ID:    id,
		State: int8(template.StateDel),
		MTime: timex.Time(now.Unix()),
	}
	_, err = s.tpl.DelTemplate(c, mid, tpl)
	return
}
