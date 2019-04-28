package creative

import (
	"context"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/log"
)

const (
	// select
	_getTypeSQL = "SELECT typeid as id, parent_id as parent, lang, name, `desc`, notice, app_notice, intro_original, intro_copy FROM archive_type ORDER BY parent_id asc, id asc"
)

// Types get all Types from creative database.
func (d *Dao) Types(c context.Context) (tops []*archive.Type, langs map[string][]*archive.Type, typeMap map[int16]*archive.Type, err error) {
	rows, err := d.getTypeStmt.Query(c)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	topL := make(map[string]map[int16]*archive.Type)
	topCh := make(map[int16]*archive.Type)
	topEn := make(map[int16]*archive.Type)
	topJp := make(map[int16]*archive.Type)
	langs = make(map[string][]*archive.Type)
	typeMap = make(map[int16]*archive.Type)
	for rows.Next() {
		t := &archive.Type{}
		if err = rows.Scan(&t.ID, &t.Parent, &t.Lang, &t.Name, &t.Desc, &t.Notice, &t.AppNotice, &t.Original, &t.IntroCopy); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if archive.ForbidTopTypesForAll(t.ID) {
			continue
		}
		t.Descapp = t.Desc
		t.Show = true
		if archive.CopyrightForCreatorAdd(t.ID) {
			t.CopyRight = 2
		}
		if t.Lang == "ch" {
			typeMap[t.ID] = t
		}
		if t.Parent == 0 {
			nt := &archive.Type{}
			*nt = *t
			switch t.Lang {
			case "ch":
				topCh[t.ID] = t
				topL[t.Lang] = topCh
				tops = append(tops, nt)
			case "en":
				topEn[t.ID] = t
				topL[t.Lang] = topEn
			case "jp":
				topJp[t.ID] = t
				topL[t.Lang] = topJp
			}
			langs[t.Lang] = append(langs[t.Lang], t)
			continue
		}
		if _, ok := topL[t.Lang][t.Parent]; !ok {
			continue
		}
		if t.Lang == topL[t.Lang][t.Parent].Lang {
			topL[t.Lang][t.Parent].Children = append(topL[t.Lang][t.Parent].Children, t)
		}
	}
	return
}
