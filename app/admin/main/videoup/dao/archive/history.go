package archive

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

const (
	_archistoryByAIDSQL = "SELECT `id`, `aid`, `title`, `tag`, `content`, `cover`, `mid`, `ctime` FROM `archive_edit_history` WHERE `aid`=? AND `ctime`>? ORDER BY `id` DESC;"
	_archistoryByIDSQL  = "SELECT `id`, `aid`, `title`, `tag`, `content`, `cover`, `mid`, `ctime` FROM `archive_edit_history` WHERE `id`=? LIMIT 1"
)

//HistoryByAID 根据aid获取稿件的用户编辑历史
func (d *Dao) HistoryByAID(c context.Context, aid int64, stime time.Time) (hs []*archive.ArcHistory, err error) {
	hs = []*archive.ArcHistory{}
	rows, err := d.db.Query(c, _archistoryByAIDSQL, aid, stime)
	if err != nil {
		log.Error("HistoryByAID d.db.Query(aid(%d)) error(%v)", aid, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		h := &archive.ArcHistory{}
		if err = rows.Scan(&h.ID, &h.AID, &h.Title, &h.Tag, &h.Content, &h.Cover, &h.MID, &h.CTime); err != nil {
			log.Error("HistoryByAID rows.Scan(aid(%d)) error(%v)", aid, err)
			return
		}
		hs = append(hs, h)
	}

	return
}

//HistoryByID 根据id获取一条稿件的用户编辑历史
func (d *Dao) HistoryByID(c context.Context, id int64) (h *archive.ArcHistory, err error) {
	h = &archive.ArcHistory{}
	row := d.db.QueryRow(c, _archistoryByIDSQL, id)
	if err = row.Scan(&h.ID, &h.AID, &h.Title, &h.Tag, &h.Content, &h.Cover, &h.MID, &h.CTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Error("HistoryByID row.Scan(id(%d)) error(%v)", id, err)
	}

	return
}
