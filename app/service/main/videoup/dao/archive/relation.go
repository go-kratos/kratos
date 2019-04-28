package archive

import (
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inRelationSQL = "INSERT INTO archive_video_relation (aid,cid,title,description,index_order) VALUES (?,?,?,?,?)"
)

// TxAddRelation insert archive_video_relation.
func (d *Dao) TxAddRelation(tx *sql.Tx, v *archive.Video) (vid int64, err error) {
	res, err := tx.Exec(_inRelationSQL, v.Aid, v.Cid, v.Title, v.Desc, v.Index)
	if err != nil {
		log.Error("d.inRelation.Exec error(%v)", err)
		return
	}
	vid, err = res.LastInsertId()
	return
}
