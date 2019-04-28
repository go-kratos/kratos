package dao

import (
	"strings"

	"go-common/app/admin/main/growup/model"
)

// ListBlacklist find blacklist by query
func (d *Dao) ListBlacklist(query string, from, limit int, sort string) (list []*model.Blacklist, total int, err error) {
	err = d.db.Table("av_black_list").Where(query).Count(&total).Error
	if err != nil {
		return
	}
	if strings.HasPrefix(sort, "-") {
		sort = strings.TrimPrefix(sort, "-")
		sort = sort + " " + "desc"
	}
	err = d.db.Table("av_black_list").Order(sort).Offset(from).Where(query).Limit(limit).Find(&list).Error
	return
}

// GetAvIncomeStatis get av total income
func (d *Dao) GetAvIncomeStatis(query string) (avIncome []*model.AvIncomeStatis, err error) {
	err = d.db.Table("av_income_statis").Where(query).Find(&avIncome).Error
	return
}

// UpdateBlacklist update blacklist
func (d *Dao) UpdateBlacklist(avID int64, ctype int, update map[string]interface{}) (err error) {
	return d.db.Table("av_black_list").Where("av_id = ? AND ctype = ? AND is_delete = 0", avID, ctype).Updates(update).Error
}
