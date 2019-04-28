package archive

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_vhistoryByHIDSQL = "SELECT `id`, `cid`, `eptitle`, `description`, `filename`, `ctime` FROM `archive_video_edit_history` WHERE `hid`=? ORDER BY `id` ASC;"
)

//VideoHistoryByHID 根据稿件编辑历史id, 获取当时视频的用户编辑历史
func (d *Dao) VideoHistoryByHID(c context.Context, hid int64) (hs []*archive.VideoHistory, err error) {
	hs = []*archive.VideoHistory{}
	rows, err := d.db.Query(c, _vhistoryByHIDSQL, hid)
	if err != nil {
		log.Error("VideoHistoryByHID d.db.Query(hid(%d)) error(%v)", hid, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		h := &archive.VideoHistory{}
		if err = rows.Scan(&h.ID, &h.CID, &h.EpTitle, &h.Description, &h.Filename, &h.CTime); err != nil {
			log.Error("VideoHistoryByHID rows.Scan(hid(%d)) error(%v)", hid, err)
			return
		}
		hs = append(hs, h)
	}

	return
}
