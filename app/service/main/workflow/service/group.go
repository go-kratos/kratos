package service

import (
	"context"

	"go-common/app/service/main/workflow/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// DeleteGroup .
func (s *Service) DeleteGroup(c context.Context, p *model.DeleteGroupParams) (err error) {
	tx := s.dao.DB.Begin()
	if tx.Error != nil {
		log.Error("s.DeleteGroup error(%v)", tx.Error)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("service tx rollback error(%v)", err1)
			}
		}
	}()
	var gid int64
	if err = tx.Table("workflow_group").Select("id").Where("business=? AND oid=? AND eid=?", p.Business, p.OID, p.EID).Row().Scan(&gid); err != nil {
		log.Error("service select gid err business(%d) oid(%d) eid(%d) error(%v)", p.Business, p.OID, p.EID, err)
		if err == sql.ErrNoRows {
			err = ecode.WkfGroupNotFound
		}
		return
	}
	if err = tx.Table("workflow_group").Where("id=?", gid).Update("state", model.StateDelete).Error; err != nil {
		log.Error("service delete group state error(%v)", err)
		return
	}
	challs := []*model.Challenge3{}
	if err = tx.Table("workflow_chall").Where("gid=?", gid).Find(&challs).Error; err != nil {
		log.Error("service find chall by gid(%d) error(%v)", gid, err)
		return
	}
	cids := []int64{}
	for _, c := range challs {
		cids = append(cids, c.ID)
	}
	if err = tx.Table("workflow_chall").Where("id IN (?)", cids).Update("dispatch_state", model.StateDelete).Error; err != nil {
		log.Error("service update chall state ids(%v) error(%v)", cids, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("service delete group commit error(%v)", err)
	}
	return
}

// PublicRefereeGroup .
func (s *Service) PublicRefereeGroup(c context.Context, prgp *model.PublicRefereeGroupParams) (err error) {
	tx := s.dao.DB.Begin()
	if tx.Error != nil {
		log.Error("s.PublicRefereeGroup error(%v)", tx.Error)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("service tx rollback error(%v)", err1)
			}
		}
	}()
	var gid int64
	if err = tx.Table("workflow_group").Select("id").Where("business=? AND oid=? AND eid=?", prgp.Business, prgp.Oid, prgp.Eid).Row().Scan(&gid); err != nil {
		log.Error("service select gid err business(%d) oid(%d) eid(%d) error(%v)", prgp.Business, prgp.Oid, prgp.Eid, err)
		if err == sql.ErrNoRows {
			err = ecode.WkfGroupNotFound
		}
		return
	}
	if err = tx.Table("workflow_group").Where("id=?", gid).Update("state", model.StatePublicReferee).Error; err != nil {
		log.Error("service public refree group state error(%v)", err)
		return
	}

	if err = tx.Table("workflow_chall").Where("gid=?", gid).Update("dispatch_state", model.StatePublicReferee).Error; err != nil {
		log.Error("service update chall state by gid(%d) error(%v)", gid, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.Error("service delete group commit error(%v)", err)
	}
	return
}
