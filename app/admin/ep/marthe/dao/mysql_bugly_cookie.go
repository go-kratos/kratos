package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertCookie Insert Cookie.
func (d *Dao) InsertCookie(buglyCookie *model.BuglyCookie) error {
	return pkgerr.WithStack(d.db.Create(buglyCookie).Error)
}

// UpdateCookie Update cookie.
func (d *Dao) UpdateCookie(buglyCookie *model.BuglyCookie) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyCookie{}).Updates(buglyCookie).Error)
}

// UpdateCookieStatus Update Cookie Status.
func (d *Dao) UpdateCookieStatus(id int64, status int) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyCookie{}).Where("id = ?", id).Update("status", status).Error)
}

// UpdateCookieUsageCount Update Cookie Usage Count.
func (d *Dao) UpdateCookieUsageCount(id int64, usageCount int) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyCookie{}).Where("id = ?", id).Update("usage_count", usageCount).Error)
}

// QueryCookieByStatus Query Cookie By Status.
func (d *Dao) QueryCookieByStatus(status int) (buglyCookies []*model.BuglyCookie, err error) {
	err = pkgerr.WithStack(d.db.Where("status=?", status).Order("ctime desc").Find(&buglyCookies).Error)
	return
}

// QueryCookieByQQAccount Query Cookie By QQ Account.
func (d *Dao) QueryCookieByQQAccount(qqAccount int) (buglyCookie *model.BuglyCookie, err error) {
	buglyCookie = &model.BuglyCookie{}
	if err = d.db.Where("qq_account=?", qqAccount).First(buglyCookie).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// FindCookies Find Cookies.
func (d *Dao) FindCookies(req *model.QueryBuglyCookiesRequest) (total int64, buglyCookies []*model.BuglyCookie, err error) {
	gDB := d.db.Model(&model.BuglyCookie{})

	if req.QQAccount != 0 {
		gDB = gDB.Where("qq_account=?", req.QQAccount)
	}

	if req.Status != 0 {
		gDB = gDB.Where("status=?", req.Status)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&buglyCookies).Error)
	return
}
