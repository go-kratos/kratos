package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/share/model"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_shareSQL    = "SELECT oid,tp,share FROM %s WHERE oid=? AND tp=?"
	_addShareSQL = "INSERT INTO %s(oid,tp,share) VALUES(?,?,1) ON DUPLICATE KEY UPDATE share=share+1"
)

func table(oid int64) string {
	return fmt.Sprintf("share_count_%02d", oid%100)
}

// Share get share
func (d *Dao) Share(c context.Context, oid int64, tp int) (share *model.Share, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_shareSQL, table(oid)), oid, tp)
	share = &model.Share{}
	if err = row.Scan(&share.OID, &share.Tp, &share.Count); err != nil {
		if err == sql.ErrNoRows {
			share = nil
			err = nil
		} else {
			err = errors.WithStack(err)
		}
		return
	}
	return
}

// AddShare add share
func (d *Dao) AddShare(c context.Context, oid int64, tp int) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_addShareSQL, table(oid)), oid, tp); err != nil {
		err = errors.Wrapf(err, "d.db.Exec(%s, %d, %d)", fmt.Sprintf(_addShareSQL, table(oid)), oid, tp)
		return
	}
	if _, ok := d.sources[oid]; ok && tp == model.ArchiveTyp {
		if _, err = d.db.Exec(c, fmt.Sprintf(_addShareSQL, table(d.c.Target)), d.c.Target, tp); err != nil {
			return
		}
		var share *model.Share
		if share, err = d.Share(c, d.c.Target, tp); err != nil {
			return
		}
		// pub msg
		if err = d.PubStatShare(c, model.ArchiveMsgTyp, d.c.Target, share.Count); err != nil {
			log.Error("s.dao.PubArchiveShare error(%v)", err)
			err = nil
		}
		d.asyncCache.Save(func() {
			if err = d.SetShareCache(context.Background(), d.c.Target, tp, share.Count); err != nil {
				log.Error("d.SetShareCache error(%v)", err)
				return
			}
		})
	}
	return
}
