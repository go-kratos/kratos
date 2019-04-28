package resource

import (
	"context"
	"time"

	"go-common/app/interface/main/web-show/model/resource"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selAllResSQL    = `SELECT id,platform,name,parent,counter,position FROM resource WHERE state=0 ORDER BY counter desc`
	_selAllAssignSQL = `SELECT id,name,contract_id,resource_id,pic,litpic,url,atype,weight,rule,agency FROM resource_assignment WHERE stime<? AND etime>? AND state=0 ORDER BY weight,stime desc`
	_selDefBannerSQL = `SELECT id,name,contract_id,resource_id,pic,litpic,url,atype,weight,rule FROM default_one WHERE  state=0`
)

func (dao *Dao) initRes() {
	dao.selAllResStmt = dao.db.Prepared(_selAllResSQL)
	dao.selAllAssignStmt = dao.db.Prepared(_selAllAssignSQL)
	dao.selDefBannerStmt = dao.db.Prepared(_selDefBannerSQL)
}

// Resources get resource infos from db
func (dao *Dao) Resources(c context.Context) (rscs []*resource.Res, err error) {
	rows, err := dao.selAllResStmt.Query(c)
	if err != nil {
		log.Error("dao.selAllResStmt query error (%v)", err)
		return
	}
	defer rows.Close()
	rscs = make([]*resource.Res, 0)
	for rows.Next() {
		rsc := &resource.Res{}
		if err = rows.Scan(&rsc.ID, &rsc.Platform, &rsc.Name, &rsc.Parent, &rsc.Counter, &rsc.Position); err != nil {
			PromError("Resources", "rows.scan err(%v)", err)
			return
		}
		rscs = append(rscs, rsc)
	}
	return
}

// Assignment get assigment from db
func (dao *Dao) Assignment(c context.Context) (asgs []*resource.Assignment, err error) {
	rows, err := dao.selAllAssignStmt.Query(c, time.Now(), time.Now())
	if err != nil {
		log.Error("dao.selAllAssignmentStmt query error (%v)", err)
		return
	}
	defer rows.Close()
	asgs = make([]*resource.Assignment, 0)
	for rows.Next() {
		asg := &resource.Assignment{}
		if err = rows.Scan(&asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic, &asg.URL, &asg.Atype, &asg.Weight, &asg.Rule, &asg.Agency); err != nil {
			PromError("Assignment", "rows.scan err(%v)", err)
			return
		}
		asgs = append(asgs, asg)
	}
	return
}

// DefaultBanner set
func (dao *Dao) DefaultBanner(c context.Context) (asg *resource.Assignment, err error) {
	row := dao.selDefBannerStmt.QueryRow(c)
	asg = &resource.Assignment{}
	if err = row.Scan(&asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic, &asg.URL, &asg.Atype, &asg.Weight, &asg.Rule); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			PromError("DefaultBanner", "dao.DefaultBanner.QueryRow error(%v)", err)
		}
	}
	return
}
