package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/member/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	_officialName         = "user_official"
	_officialDocName      = "user_official_doc"
	_officialDocAdditName = "user_official_doc_addit"
)

// Official is.
func (d *Dao) Official(ctx context.Context, mid int64) (off *model.Official, err error) {
	off = &model.Official{}
	if err = d.member.Table(_officialName).Where("mid=?", mid).Find(off).Error; err != nil {
		err = errors.Wrap(err, "official docs")
	}
	return
}

// OfficialEdit is.
func (d *Dao) OfficialEdit(ctx context.Context, mid int64, role int8, title, desc string) (off *model.Official, err error) {
	off = &model.Official{
		Mid: mid,
	}
	attrs := map[string]interface{}{
		"role":        role,
		"title":       title,
		"description": desc,
	}
	if err = d.member.Table(_officialName).Where("mid=?", mid).Assign(attrs).FirstOrCreate(off).Error; err != nil {
		err = errors.Wrap(err, "official edit")
	}
	return
}

// Officials is.
func (d *Dao) Officials(ctx context.Context, mid int64, roles []int8, stime, etime time.Time, pn, ps int) (offs []*model.Official, total int, err error) {
	where := "role in (?) AND ctime>? AND ctime<?"
	if mid > 0 {
		where = "mid=? AND " + where
		err = d.member.Table(_officialName).Order("ctime DESC").Offset((pn-1)*ps).Limit(ps).Where(where, mid, roles, stime, etime).Find(&offs).Error
		d.member.Table(_officialName).Where(where, mid, roles, stime, etime).Count(&total)
	} else {
		err = d.member.Table(_officialName).Order("ctime DESC").Offset((pn-1)*ps).Limit(ps).Where(where, roles, stime, etime).Find(&offs).Error
		d.member.Table(_officialName).Where(where, roles, stime, etime).Count(&total)
	}
	if err != nil {
		err = errors.Wrap(err, "official docs")
	}
	return
}

// OfficialDoc is.
func (d *Dao) OfficialDoc(ctx context.Context, mid int64) (off *model.OfficialDoc, err error) {
	off = &model.OfficialDoc{OfficialExtra: &model.OfficialExtra{}}
	if err = d.member.Table(_officialDocName).Where("mid=?", mid).First(off).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		err = errors.Wrap(err, "official doc")
		off = nil
		return
	}
	off.ParseExtra()
	return
}

// OfficialDocs is.
func (d *Dao) OfficialDocs(ctx context.Context, mid int64, roles, states []int8, uname string, stime, etime time.Time, pn, ps int) (offs []*model.OfficialDoc, total int, err error) {
	query := d.member.Table(_officialDocName).Order("submit_time DESC")
	if uname != "" {
		query = query.Where("uname=?", uname)
	}
	if mid > 0 {
		query = query.Where("mid=?", mid)
	}
	query = query.Where("role IN (?)", roles).
		Where("state IN (?)", states).
		Where("ctime > ?", stime).
		Where("ctime < ?", etime)

	query.Count(&total)
	if err = query.Offset((pn - 1) * ps).Limit(ps).Find(&offs).Error; err != nil {
		err = errors.Wrap(err, "official docs")
		return
	}

	for _, od := range offs {
		od.ParseExtra()
	}
	return
}

// OfficialDocsByMids is.
func (d *Dao) OfficialDocsByMids(ctx context.Context, mids []int64) (map[int64]*model.OfficialDoc, error) {
	ofl := make([]*model.OfficialDoc, 0, len(mids))
	offs := make(map[int64]*model.OfficialDoc, len(mids))
	if err := d.member.Table(_officialDocName).Where("mid IN (?)", mids).Find(&ofl).Error; err != nil {
		err = errors.Wrap(err, "official docs")
		return nil, err
	}
	for _, of := range ofl {
		offs[of.Mid] = of
	}
	for _, od := range offs {
		od.ParseExtra()
	}
	return offs, nil
}

// OfficialDocAudit is.
func (d *Dao) OfficialDocAudit(ctx context.Context, mid int64, state int8, uname string, isInternal bool, rejectReason string) (err error) {
	ups := map[string]interface{}{
		"state":         state,
		"uname":         uname,
		"is_internal":   isInternal,
		"reject_reason": rejectReason,
	}
	if err = d.member.Table(_officialDocName).Where("mid=?", mid).Updates(ups).Error; err != nil {
		err = errors.Wrap(err, "official doc audit")
	}
	return
}

// OfficialDocEdit is.
func (d *Dao) OfficialDocEdit(ctx context.Context, mid int64, name string, role, state int8, title, desc, extra string, uname string, isInternal bool) (err error) {
	off := &model.OfficialDoc{
		Mid:          mid,
		SubmitSource: "admin",
	}
	ups := map[string]interface{}{
		"state":       state,
		"role":        role,
		"name":        name,
		"title":       title,
		"description": desc,
		"extra":       extra,
		"uname":       uname,
		"is_internal": isInternal,
	}
	if err = d.member.Table(_officialDocName).Where("mid=?", mid).Assign(ups).FirstOrCreate(off).Error; err != nil {
		err = errors.Wrap(err, "official doc audit")
	}
	return
}

// OfficialDocSubmit is.
func (d *Dao) OfficialDocSubmit(ctx context.Context, mid int64, name string, role, state int8, title, desc, extra string, uname string, isInternal bool, submitSource string) (err error) {
	off := &model.OfficialDoc{
		Mid: mid,
	}
	ups := map[string]interface{}{
		"state":         state,
		"role":          role,
		"name":          name,
		"title":         title,
		"description":   desc,
		"extra":         extra,
		"uname":         uname,
		"is_internal":   isInternal,
		"submit_source": submitSource,
	}
	if err = d.member.Table(_officialDocName).Where("mid=?", mid).Assign(ups).FirstOrCreate(off).Error; err != nil {
		err = errors.Wrap(err, "official doc audit")
	}
	return
}

// OfficialDocAddits .
func (d *Dao) OfficialDocAddits(ctx context.Context, property string, vstring string) ([]*model.OfficialDocAddit, error) {
	addits := make([]*model.OfficialDocAddit, 0)
	err := d.member.Table(_officialDocAdditName).Where("property=? and vstring=?", property, vstring).Find(&addits).Error
	if err != nil {
		err = errors.Wrap(err, "find official doc addit error")
		return nil, err
	}
	return addits, nil
}
