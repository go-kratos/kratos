package archive

import (
	"context"
	"fmt"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_nameByIdsSQL = "SELECT id,description FROM archive_tag WHERE id IN (%s)"
)

// TagNameMap get audit tag id and name map
func (d *Dao) TagNameMap(c context.Context, ids []int64) (nameMap map[int64]string, err error) {
	nameMap = make(map[int64]string)
	if len(ids) == 0 {
		return
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_nameByIdsSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tag := struct {
			ID   int64
			Name string
		}{}
		if err = rows.Scan(&tag.ID, &tag.Name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		nameMap[tag.ID] = tag.Name
	}
	return
}
