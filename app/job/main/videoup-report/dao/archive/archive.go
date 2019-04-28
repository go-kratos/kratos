package archive

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/log"
)

const (
	_archiveByAidSQL = "SELECT id,mid,typeid,copyright,author,title,cover,content,duration,round,attribute,access,state,tag,pubtime,ctime,mtime FROM archive WHERE id=? "
	_arcNoteSQL      = "SELECT coalesce(note,'') from archive where id=?"
	_upOriginalAids  = "SELECT " +
		"a.id,a.mid,a.typeid,a.copyright,a.author,a.title,a.cover,a.content,a.duration,a.round,a.attribute,a.access,a.state,a.tag,a.pubtime,a.ctime,a.mtime " +
		"FROM archive AS a LEFT JOIN archive_delay as delay ON delay.aid = a.id " +
		"WHERE a.mid=? AND a.copyright=? AND a.ctime>=? AND a.ctime<? AND (a.state >= 0 OR delay.state >= 0)"
	_upArcTagSQL = "UPDATE archive SET tag=? WHERE id=?"
)

// ArchiveByAid get archive by aid
func (d *Dao) ArchiveByAid(c context.Context, aid int64) (arc *archive.Archive, err error) {
	row := d.db.QueryRow(c, _archiveByAidSQL, aid)
	arc = &archive.Archive{}
	if err = row.Scan(&arc.ID, &arc.Mid, &arc.TypeID, &arc.Copyright, &arc.Author, &arc.Title, &arc.Cover, &arc.Desc, &arc.Duration,
		&arc.Round, &arc.Attribute, &arc.Access, &arc.State, &arc.Tag, &arc.PTime, &arc.CTime, &arc.MTime); err != nil {
		log.Error("row.Scan error(%v)", err)
	}
	return
}

//ArchiveNote 稿件审核的备注字段，可能为NIL
func (d *Dao) ArchiveNote(c context.Context, aid int64) (note string, err error) {
	if err = d.db.QueryRow(c, _arcNoteSQL, aid).Scan(&note); err != nil {
		log.Error("ArchiveNote db.row.Scan error(%v), aid(%d)", err, aid)
	}

	return
}

// ExcitationArchivesByTime 获取Up主过审的自制稿件
func (d *Dao) ExcitationArchivesByTime(c context.Context, mid int64, bt, et time.Time) (archives []*archive.Archive, err error) {
	archives = []*archive.Archive{}
	if mid < 1 {
		err = fmt.Errorf("wrong mid(%d)", mid)
		return
	}
	rows, err := d.db.Query(c, _upOriginalAids, mid, archive.CopyrightOriginal, bt, et)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%v,%v) error(%v)", _upOriginalAids, mid, archive.CopyrightOriginal, bt, et)
		return
	}
	defer rows.Close()
	for rows.Next() {
		arc := &archive.Archive{}
		if err = rows.Scan(&arc.ID, &arc.Mid, &arc.TypeID, &arc.Copyright, &arc.Author, &arc.Title, &arc.Cover, &arc.Desc, &arc.Duration,
			&arc.Round, &arc.Attribute, &arc.Access, &arc.State, &arc.Tag, &arc.PTime, &arc.CTime, &arc.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		archives = append(archives, arc)
	}
	return
}

//UpTag update archive tag
func (d *Dao) UpTag(c context.Context, aid int64, tags string) (rows int64, err error) {
	res, err := d.db.Exec(c, _upArcTagSQL, tags, aid)
	if err != nil {
		log.Error("d.UpTag.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
