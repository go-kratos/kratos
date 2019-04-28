package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/relation/model/i64b"
	sml "go-common/app/service/main/relation/model"
	"go-common/library/database/sql"
)

const (
	_shard        = 500
	_tagUserShard = 500
	// following
	_getRelationSQL     = "SELECT r.attribute,r.mtime,t.tag FROM user_relation_mid_%03d AS r join user_relation_tag_user_%03d AS t ON t.mid=r.mid AND t.fid=r.fid WHERE r.mid=? AND r.fid=? AND r.status=0 "
	_UserSetAchieveFlag = "INSERT INTO user_addit (mid,achieve_flags) VALUES (?,?) ON DUPLICATE KEY UPDATE achieve_flags=achieve_flags|VALUES(achieve_flags)"
)

func hit(id int64) int64 {
	return id % _shard
}

func tagUserHit(id int64) int64 {
	return id % _tagUserShard
}

// UserRelation get user relation attr.
func (d *Dao) UserRelation(c context.Context, mid, fid int64) (f *sml.Following, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getRelationSQL, hit(mid), tagUserHit(mid)), mid, fid)
	f = new(sml.Following)
	var ttag i64b.Int64Bytes
	if err = row.Scan(&f.Attribute, &f.MTime, &ttag); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			f = nil
		}
		return
	}
	f.Mid = fid
	f.Tag = []int64(ttag)
	for _, id := range f.Tag {
		if id == -10 {
			f.Special = 1
		}
	}
	return
}

// UserSetAchieveFlag is
func (d *Dao) UserSetAchieveFlag(ctx context.Context, mid int64, flag uint64) (int64, error) {
	res, err := d.db.Exec(ctx, _UserSetAchieveFlag, mid, flag)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
