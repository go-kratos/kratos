package medal

import (
	"context"
	"fmt"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_sharding = 10

	_insOwnerSQL               = "INSERT INTO medal_owner_%s(mid,nid) VALUES (?,?)"
	_updateOwnerSQL            = "UPDATE medal_owner_%s SET is_activated=? WHERE mid=? AND nid=?"
	_updateOwnerUnallSQL       = "UPDATE medal_owner_%s SET is_activated=0 WHERE mid=? AND nid!=?"
	_selInfoAllSQL             = "SELECT id,name,description,image,image_small,cond,gid,level,level_rank,sort FROM medal_info ORDER BY sort ASC,gid ASC,level ASC"
	_selOwnerByMidSQL          = "SELECT id,mid,nid,is_activated,ctime,mtime FROM medal_owner_%s WHERE mid=? AND is_del=0 ORDER BY ctime DESC"
	_selInfoByNidSQL           = "SELECT name FROM medal_info WHERE id=? AND is_online=1"
	_selInstalledOwnerBYMidSQL = "SELECT nid FROM medal_owner_%s WHERE mid=? AND is_activated=1 AND is_del=0 LIMIT 1"
	_countOwnerBYNidMidSQL     = "SELECT COUNT(*) FROM medal_owner_%s WHERE mid=? AND nid=?"
	_OwnerBYNidMidSQL          = "SELECT id,mid,nid,is_activated,ctime,mtime FROM medal_owner_%s WHERE mid=? AND nid=?"
	_selGroupAllSQL            = "SELECT id,name,pid,rank FROM medal_group WHERE is_online=1 ORDER BY pid ASC,rank ASC"
)

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%d", id%_sharding)
}

// AddMedalOwner insert into medal_owner.
func (d *Dao) AddMedalOwner(c context.Context, mid, nid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_insOwnerSQL, d.hit(mid)), mid, nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// InstallMedalOwner update medal_owner set is_activated=1.
func (d *Dao) InstallMedalOwner(c context.Context, mid, nid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateOwnerSQL, d.hit(mid)), model.OwnerInstall, mid, nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UninstallMedalOwner update medal_owner set is_activated=0.
func (d *Dao) UninstallMedalOwner(c context.Context, mid, nid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateOwnerSQL, d.hit(mid)), model.OwnerUninstall, mid, nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UninstallAllMedalOwner uninst all medal_owner set is_activated=0.
func (d *Dao) UninstallAllMedalOwner(c context.Context, mid, nid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateOwnerUnallSQL, d.hit(mid)), mid, nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// MedalInfoAll retun all medal_info where is_online=1.
func (d *Dao) MedalInfoAll(c context.Context) (res map[int64]*model.MedalInfo, err error) {
	res = make(map[int64]*model.MedalInfo)
	rows, err := d.db.Query(c, _selInfoAllSQL)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := new(model.MedalInfo)
		if err = rows.Scan(&info.ID, &info.Name, &info.Description, &info.Image, &info.ImageSmall, &info.Condition, &info.GID, &info.Level, &info.LevelRank, &info.Sort); err != nil {
			err = errors.WithStack(err)
			return
		}
		info.Build()
		res[info.ID] = info
	}
	err = rows.Err()
	return
}

// MedalOwnerByMid return medal_owner by mid.
func (d *Dao) MedalOwnerByMid(c context.Context, mid int64) (res []*model.MedalOwner, err error) {
	res = make([]*model.MedalOwner, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_selOwnerByMidSQL, d.hit(mid)), mid)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.MedalOwner)
		if err = rows.Scan(&r.ID, &r.MID, &r.NID, &r.IsActivated, &r.CTime, &r.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// MedalInfoByNid return medal_info by nid.
func (d *Dao) MedalInfoByNid(c context.Context, nid int64) (res *model.MedalInfo, err error) {
	res = &model.MedalInfo{}
	rows := d.db.QueryRow(c, _selInfoByNidSQL, nid)
	if err = rows.Scan(&res.Name); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "InstalledOwnerBYMid")
			return
		}
		err = nil
	}
	return
}

// ActivatedOwnerByMid retun nid of medal_owner by mid where is_activated=1.
func (d *Dao) ActivatedOwnerByMid(c context.Context, mid int64) (nid int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selInstalledOwnerBYMidSQL, d.hit(mid)), mid)
	if err = row.Scan(&nid); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "InstalledOwnerBYMid")
			return
		}
		err = nil
	}
	return
}

// CountOwnerBYNidMid retun number of medal_owner by mid and nid.
func (d *Dao) CountOwnerBYNidMid(c context.Context, mid, nid int64) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_countOwnerBYNidMidSQL, d.hit(mid)), mid, nid)
	if err = row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "CountOwnerBYNidMid")
			return
		}
		count = 0
		err = nil
	}
	return
}

// OwnerBYNidMid retun  medal_owner by mid and nid.
func (d *Dao) OwnerBYNidMid(c context.Context, mid, nid int64) (res *model.MedalOwner, err error) {
	res = &model.MedalOwner{}
	row := d.db.QueryRow(c, fmt.Sprintf(_OwnerBYNidMidSQL, d.hit(mid)), mid, nid)
	if err = row.Scan(&res.ID, &res.MID, &res.NID, &res.IsActivated, &res.CTime, &res.MTime); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "OwnerBYNidMid")
			return
		}
		res = nil
		err = nil
	}
	return
}

// MedalGroupAll retun all medal_group where is_online=1.
func (d *Dao) MedalGroupAll(c context.Context) (res []*model.MedalGroup, err error) {
	rows, err := d.db.Query(c, _selGroupAllSQL)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := new(model.MedalGroup)
		if err = rows.Scan(&info.ID, &info.Name, &info.PID, &info.Rank); err != nil {
			err = errors.WithStack(err)
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return
}
