package dao

import (
	"context"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_addRegSQL  = `INSERT INTO tv_pages (page_id,title,index_type,index_tid,rank) VALUES (?,?,?,?,?)`
	_editRegSQL = `UPDATE tv_pages SET title=?,index_type=?,index_tid=? WHERE page_id=?`
	_del        = 0
	_firpid     = 1
)

// RegList region list .
func (d *Dao) RegList(ctx context.Context, arg *model.Param) (res []*model.RegDB, err error) {
	if arg.PageID != "" {
		if err = d.DB.Where("page_id=? AND deleted=?", arg.PageID, _del).Find(&res).Error; err != nil {
			log.Error("d.DB.Where error(%v)", err)
		}
	} else {
		if arg.State == "2" {
			if err = d.DB.Order("rank ASC", true).Where("title LIKE ? AND deleted=?", "%"+arg.Title+"%", _del).Find(&res).Error; err != nil {
				log.Error("d.DB.Where error(%v)", err)
			}
			return
		}
		if err = d.DB.Order("rank ASC", true).Where("title LIKE ? AND deleted=? AND valid=?", "%"+arg.Title+"%", _del, arg.State).Find(&res).Error; err != nil {
			log.Error("d.DB.Order Where error(%v)", err)
		}
	}
	return
}

// AddReg add region .
func (d *Dao) AddReg(ctx context.Context, title, itype, itid, rank string) (err error) {
	var (
		pid int
		res = &model.RegDB{}
	)
	tx := d.DB.Begin()
	if err = d.DB.Order("page_id DESC", true).First(res).Error; err != nil || res == nil {
		if err == ecode.NothingFound {
			if err = d.DB.Exec(_addRegSQL, _firpid, title, itype, itid, rank).Error; err != nil {
				log.Error(" d.DB.Exec error(%v)", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			return
		}
		log.Error(" d.DB.Order error(%v)", err)
		tx.Rollback()
		return
	}
	pid = res.PageID + 1
	if err = d.DB.Exec(_addRegSQL, pid, title, itype, itid, rank).Error; err != nil {
		log.Error(" d.DB.Exec error(%v)", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// EditReg edit region .
func (d *Dao) EditReg(ctx context.Context, pid, title, itype, itid string) (err error) {
	if err = d.DB.Exec(_editRegSQL, title, itype, itid, pid).Error; err != nil {
		log.Error(" d.DB.Exec error(%v)", err)
	}
	return
}

// UpState publish or not .
func (d *Dao) UpState(ctx context.Context, pids []int, state string) (err error) {
	m := map[string]string{"valid": state}
	if err = d.DB.Table("tv_pages").Where("page_id IN (?)", pids).Updates(m).Error; err != nil {
		log.Error(" d.DB.Table.Where.Updates error(%v)", err)
	}
	return
}
