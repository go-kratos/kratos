package dao

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/videoup-task/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_arcSQL          = "SELECT id,mid,title,access,attribute,reject_reason,tag,forward,round,state,copyright,cover,content,typeid,pubtime,ctime,mtime FROM archive WHERE id=?"
	_archiveParamSQL = `SELECT a.typeid,addit.up_from FROM archive AS a LEFT JOIN archive_addit AS addit ON a.id=addit.aid WHERE a.id=?`
)

// Archive get archive by aid
func (d *Dao) Archive(c context.Context, aid int64) (a *model.Archive, err error) {
	var (
		row         = d.arcDB.QueryRow(c, _arcSQL, aid)
		reason, tag sql.NullString
	)
	a = &model.Archive{}
	if err = row.Scan(&a.Aid, &a.Mid, &a.Title, &a.Access, &a.Attribute, &reason, &tag, &a.Forward, &a.Round, &a.State,
		&a.Copyright, &a.Cover, &a.Desc, &a.TypeID, &a.PTime, &a.CTime, &a.MTime); err != nil {
		if err == xsql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	a.RejectReason = reason.String
	a.Tag = tag.String
	return
}

// ArchiveParam .
func (d *Dao) ArchiveParam(c context.Context, aid int64) (typeid int16, upfrom int8, err error) {
	if err = d.arcDB.QueryRow(c, _archiveParamSQL, aid).Scan(&typeid, &upfrom); err != nil {
		log.Error("ArchiveParam error(%v)", err)
	}
	return
}
