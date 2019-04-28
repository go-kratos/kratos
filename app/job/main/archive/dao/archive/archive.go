package archive

import (
	"context"
	"database/sql"

	"go-common/app/job/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_arcSQL = "SELECT id,mid,typeid,copyright,author,title,cover,content,tag,duration,round,attribute,access,state,reject_reason,pubtime,ctime,mtime,forward FROM archive WHERE id=?"
)

// Archive get a archive by avid.
func (d *Dao) Archive(c context.Context, aid int64) (a *archive.Archive, err error) {
	var reason sql.NullString
	row := d.db.QueryRow(c, _arcSQL, aid)
	a = &archive.Archive{}
	if err = row.Scan(&a.ID, &a.Mid, &a.TypeID, &a.Copyright, &a.Author, &a.Title, &a.Cover, &a.Content, &a.Tag, &a.Duration,
		&a.Round, &a.Attribute, &a.Access, &a.State, &reason, &a.PubTime, &a.CTime, &a.MTime, &a.Forward); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	a.Reason = reason.String
	return
}
