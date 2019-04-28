package archive

import (
	"context"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_descFormatSQL = "SELECT id,typeid,copyright,components,lang,platform FROM archive_desc_format WHERE state=0"
)

// DescFormats get desc_format info.
func (d *Dao) DescFormats(c context.Context) (dfs []*archive.DescFormat, err error) {
	rows, err := d.rddb.Query(c, _descFormatSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		df := &archive.DescFormat{}
		if err = rows.Scan(&df.ID, &df.TypeID, &df.Copyright, &df.Components, &df.Lang, &df.Platform); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		dfs = append(dfs, df)
	}
	return
}
