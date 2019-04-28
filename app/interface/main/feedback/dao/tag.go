package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/feedback/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_selTagBySsnID = "SELECT session_tag.session_id,tag.id,tag.name,tag.platform,tag.type FROM tag INNER JOIN session_tag ON tag.id =session_tag.tag_id  WHERE session_id IN (%s)"
	_selTagID      = "SELECT tag_id,session_id FROM session_tag WHERE session_id IN (%s)"
	_inSsnTag      = "INSERT INTO session_tag (session_id,tag_id,ctime) VALUES (?,?,?)"
)

// TagBySsnID get tag by ssnID.
func (d *Dao) TagBySsnID(c context.Context, sids []int64) (tagMap map[int64][]*model.Tag, err error) {
	rows, err := d.dbMs.Query(c, fmt.Sprintf(_selTagBySsnID, xstr.JoinInts(sids)))
	if err != nil {
		log.Error("d.dbMs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	tagMap = make(map[int64][]*model.Tag)
	for rows.Next() {
		var (
			tag = &model.Tag{}
			sid int64
		)
		if err = rows.Scan(&sid, &tag.ID, &tag.Name, &tag.Platform, &tag.Type); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if tag.Type == 0 {
			tagMap[sid] = append(tagMap[sid], tag)
		}
	}
	return
}
