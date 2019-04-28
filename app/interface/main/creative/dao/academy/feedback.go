package academy

import (
	"context"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/log"
)

const (
	// insert
	_inFbSQL = "INSERT IGNORE INTO academy_feedback (category, course, suggest, ctime, mtime, mid) VALUES (?,?,?,?,?,?)"
)

// AddFeedBack  add academy_feedback.
func (d *Dao) AddFeedBack(c context.Context, f *academy.FeedBack, mid int64) (id int64, err error) {
	res, err := d.db.Exec(c, _inFbSQL, f.Category, f.Course, f.Suggest, f.CTime, f.MTime, mid)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}
