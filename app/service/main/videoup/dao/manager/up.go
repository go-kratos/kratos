package manager

import (
	"context"
	"database/sql"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_upsWithGroup = "SELECT ups.id,mid,up_group.id as group_id ,up_group.short_tag as group_tag,up_group.name as group_name,ups.note,ups.ctime FROM ups INNER JOIN up_group on  ups.type=up_group.id"
)

// UpSpecial load all ups with group info
func (d *Dao) UpSpecial(c context.Context) (ups []*archive.Up, err error) {
	rows, err := d.managerDB.Query(c, _upsWithGroup)
	if err != nil {
		log.Error("d.tpsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	var note sql.NullString
	for rows.Next() {
		up := &archive.Up{}
		if err = rows.Scan(&up.ID, &up.Mid, &up.GroupID, &up.GroupTag, &up.GroupName, &note, &up.CTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		up.Note = note.String
		ups = append(ups, up)
	}
	return
}
