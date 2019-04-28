package like

import (
	"context"
	"fmt"

	l "go-common/app/interface/main/activity/model/like"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_contentSQL = "select id,message,ip,plat,device,ctime,mtime,image,reply,link,ex_name from like_content where id in (%s)"
)

// RawLikeContent .
func (dao *Dao) RawLikeContent(c context.Context, ids []int64) (res map[int64]*l.LikeContent, err error) {
	rows, err := dao.db.Query(c, fmt.Sprintf(_contentSQL, xstr.JoinInts(ids)))
	if err != nil {
		err = errors.Wrap(err, "dao.db.Query()")
		return
	}
	defer rows.Close()
	res = make(map[int64]*l.LikeContent, len(ids))
	for rows.Next() {
		t := &l.LikeContent{}
		if err = rows.Scan(&t.ID, &t.Message, &t.IP, &t.Plat, &t.Device, &t.Ctime, &t.Mtime, &t.Image, &t.Reply, &t.Link, &t.ExName); err != nil {
			err = errors.Wrapf(err, "rows.Scan()")
			return
		}
		res[t.ID] = t
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, " rows.Err()")
	}
	return
}
