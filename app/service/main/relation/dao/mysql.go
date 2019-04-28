package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/model/i64b"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_shard        = 500
	_statShard    = 50
	_tagShard     = 100
	_tagUserShard = 500
	_defaultTag   = "默认分组"
	_specialTag   = "特别关注"
	// relation
	_getFollowingSQL   = "SELECT fid,attribute,mtime FROM user_relation_mid_%03d WHERE mid=? AND status=0"
	_getFollowingInSQL = "SELECT fid,attribute,mtime FROM user_relation_mid_%03d WHERE mid=? AND status=0 AND fid IN (%s)"
	_addFollowingSQL   = "INSERT INTO user_relation_mid_%03d (mid,fid,attribute,source,status, ctime,mtime) VALUES (?,?,?,?,0,?,?) ON DUPLICATE KEY UPDATE attribute=attribute|?,source=?,status=0,mtime=?"
	_setFollowingSQL   = "UPDATE user_relation_mid_%03d SET attribute=?,source=?,status=?,mtime=? WHERE mid=? AND fid=?"
	_getFollowersSQL   = "SELECT mid,attribute,mtime FROM user_relation_fid_%03d WHERE fid=? AND status=0 AND attribute IN (2,6) ORDER BY mtime DESC LIMIT 1000"
	_addFollowersSQL   = "INSERT INTO user_relation_fid_%03d (mid,fid,attribute,source,status,ctime,mtime) VALUES (?,?,?,?,0,?,?) ON DUPLICATE KEY UPDATE attribute=attribute|?,source=?,status=0,mtime=?"
	_setFollowersSQL   = "UPDATE user_relation_fid_%03d SET attribute=?,source=?,status=?,mtime=? WHERE fid=? AND mid=?"
	_getRelationSQL    = "SELECT attribute FROM user_relation_mid_%03d WHERE mid=? AND fid=? AND status=0 LIMIT 1"
	_addStatIgnoreSQL  = "INSERT IGNORE INTO user_relation_stat_%02d (mid,following,whisper,black,follower,ctime,mtime) VALUES (?,?,?,?,?,?,?)"
	_addStatSQL        = "UPDATE user_relation_stat_%02d SET following=following+?,whisper=whisper+?,black=black+?,follower=follower+?,mtime=? WHERE mid=?"
	_setStatSQL        = "UPDATE user_relation_stat_%02d SET following=?,whisper=?,black=?,follower=?,mtime=? WHERE mid=?"
	_getStatSQL        = "SELECT mid,following,whisper,black,follower, ctime,mtime from user_relation_stat_%02d where mid=?"
	_getTxStatSQL      = "SELECT mid,following,whisper,black,follower,ctime,mtime from user_relation_stat_%02d where mid=? FOR UPDATE"
	// relation tag table
	_getTagsSQL    = "SELECT id,name,status,mtime FROM user_relation_tag_%02d WHERE mid=? AND status=0"
	_addTagSQL     = "INSERT INTO user_relation_tag_%02d (mid,name,status,ctime,mtime) VALUES (?,?,0,?,?)"
	_delTagSQL     = "DELETE FROM user_relation_tag_%02d WHERE id=? AND mid=?"
	_setTagNameSQL = "UPDATE user_relation_tag_%02d SET name=?,mtime=? WHERE id=?"
	// relation tag user table
	_getTagUserSQL      = "SELECT fid,tag FROM user_relation_tag_user_%03d WHERE mid=?"
	_getUsersTagSQL     = "SELECT fid,tag FROM user_relation_tag_user_%03d WHERE mid=? AND fid IN(%s)"
	_getTagsByMidFidSQL = "SELECT fid,tag,mtime FROM user_relation_tag_user_%03d WHERE mid=? and fid=?"
	_addTagUserSQL      = "INSERT INTO user_relation_tag_user_%03d (mid,fid,tag,ctime,mtime) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE tag=?"
	_setTagUserSQL      = "UPDATE user_relation_tag_user_%03d SET tag=?,mtime=? WHERE mid=? AND fid=?"
	_delTagUserSQL      = "DELETE FROM user_relation_tag_user_%03d WHERE mid=? AND fid=?"
	// relation monitor
	_loadMonitorSQL = "SELECT mid FROM user_relation_monitor"
	_addMonitorSQL  = "INSERT IGNORE INTO user_relation_monitor (mid,ctime,mtime) VALUES (?,?,?)"
	_delMonitorSQL  = "DELETE FROM user_relation_monitor WHERE mid = ?"
	// relation achieve
	_hasReachAchieve = "SELECT count(1) FROM user_addit WHERE mid=? AND achieve_flags>=?"
	// follower notify
	_getFollowerNotifySettingSQL = "SELECT disable_follower_notify FROM user_addit where mid=?"
	_enableFollowerNotifySQL     = "INSERT INTO user_addit (mid, disable_follower_notify) VALUES(?, 0) ON DUPLICATE KEY UPDATE disable_follower_notify=0"
	_disableFollowerNotifySQL    = "INSERT INTO user_addit (mid, disable_follower_notify) VALUES(?, 1) ON DUPLICATE KEY UPDATE disable_follower_notify=1"
)

func hit(id int64) int64 {
	return id % _shard
}

func statHit(id int64) int64 {
	return id % _statShard
}

func tagHit(id int64) int64 {
	return id % _tagShard
}

func tagUserHit(id int64) int64 {
	return id % _tagUserShard
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// Followings get user's following list.
func (d *Dao) Followings(c context.Context, mid int64) (res []*model.Following, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getFollowingSQL, hit(mid)), mid); err != nil {
		log.Error("d.query[%s].hit(%d).mid(%d) error(%v)", _getFollowingSQL, hit(mid), mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Following)
		if err = rows.Scan(&r.Mid, &r.Attribute, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// FollowingsIn get user's following list by in fids
func (d *Dao) FollowingsIn(c context.Context, mid int64, fids []int64) (res []*model.Following, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getFollowingInSQL, hit(mid), xstr.JoinInts(fids)), mid); err != nil {
		log.Error("d.query[%d].Query(%d) error(%v)", hit(mid), mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Following)
		if err = rows.Scan(&r.Mid, &r.Attribute, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// TxAddFollowing add following by transaction.
func (d *Dao) TxAddFollowing(c context.Context, tx *sql.Tx, mid, fid int64, mask uint32, source uint8, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addFollowingSQL, hit(mid)), mid, fid, mask, source, now, now, mask, source, now); err != nil {
		log.Error("add following: tx.Exec(%v, %d, %d, %d) error(%v)", mid, fid, mask, source, err)
		return
	}
	return res.RowsAffected()
}

// TxSetFollowing set following by transaction.
func (d *Dao) TxSetFollowing(c context.Context, tx *sql.Tx, mid, fid int64, attribute uint32, source uint8, status int, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_setFollowingSQL, hit(mid)), attribute, source, status, now, mid, fid); err != nil {
		log.Error("tx.Exec(%d, %d, %d, %d, %d) error(%v)", mid, fid, attribute, source, status, err)
		return
	}
	return res.RowsAffected()
}

// Followers get user's latest 1000 followers(attribute = AttrFollowing), order by mtime desc.
func (d *Dao) Followers(c context.Context, mid int64) (res []*model.Following, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getFollowersSQL, hit(mid)), mid); err != nil {
		log.Error("d.query(%s).hit(%d).mid(%d) error(%v)", _getFollowersSQL, hit(mid), mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Following)
		if err = rows.Scan(&r.Mid, &r.Attribute, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// TxAddFollower add follower by transaction.
func (d *Dao) TxAddFollower(c context.Context, tx *sql.Tx, mid, fid int64, mask uint32, source uint8, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addFollowersSQL, hit(fid)), mid, fid, mask, source, now, now, mask, source, now); err != nil {
		log.Error("add follower: tx.Exec(%v, %d, %d, %d), error(%v)", mid, fid, mask, source, err)
		return
	}
	return res.RowsAffected()
}

// TxSetFollower set follower by transaction.
func (d *Dao) TxSetFollower(c context.Context, tx *sql.Tx, mid, fid int64, attribute uint32, source uint8, status int, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_setFollowersSQL, hit(fid)), attribute, source, status, now, fid, mid); err != nil {
		log.Error("tx.Exec(%d, %d, %d, %d, %d) error(%v)", mid, fid, attribute, source, status, err)
		return
	}
	return res.RowsAffected()
}

// Stat get stat.
func (d *Dao) Stat(c context.Context, mid int64) (stat *model.Stat, err error) {
	var row = d.db.QueryRow(c, fmt.Sprintf(_getStatSQL, statHit(mid)), mid)
	stat = new(model.Stat)
	if err = row.Scan(&stat.Mid, &stat.Following, &stat.Whisper, &stat.Black, &stat.Follower, &stat.CTime, &stat.MTime); err != nil {
		if err == sql.ErrNoRows {
			stat = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// TxStat get stat for update by transaction.
func (d *Dao) TxStat(c context.Context, tx *sql.Tx, mid int64) (stat *model.Stat, err error) {
	row := tx.QueryRow(fmt.Sprintf(_getTxStatSQL, statHit(mid)), mid)
	stat = new(model.Stat)
	if err = row.Scan(&stat.Mid, &stat.Following, &stat.Whisper, &stat.Black, &stat.Follower, &stat.CTime, &stat.MTime); err != nil {
		if err == sql.ErrNoRows {
			stat = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// AddStat try add stat.
func (d *Dao) AddStat(c context.Context, mid int64, stat *model.Stat, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_addStatIgnoreSQL, statHit(mid)), mid, stat.Following, stat.Whisper, stat.Black, stat.Follower, now, now); err != nil {
		log.Error("d.db.Exec(%s, %d, %v, %v) error(%v)", _addStatIgnoreSQL, mid, stat, now, err)
		return
	}
	return res.RowsAffected()
}

// TxAddStat add params stat to stat by transaction.
func (d *Dao) TxAddStat(c context.Context, tx *sql.Tx, mid int64, stat *model.Stat, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addStatSQL, statHit(mid)), stat.Following, stat.Whisper, stat.Black, stat.Follower, now, mid); err != nil {
		log.Error("tx.Exec(%s, %d, %v, %v) error(%v)", _addStatSQL, mid, stat, now, err)
		return
	}
	return res.RowsAffected()
}

// TxSetStat set stat to params stat by transaction.
func (d *Dao) TxSetStat(c context.Context, tx *sql.Tx, mid int64, stat *model.Stat, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_setStatSQL, statHit(mid)), stat.Following, stat.Whisper, stat.Black, stat.Follower, now, mid); err != nil {
		log.Error("tx.Exec(%s, %d, %v, %v) error(%v)", _setStatSQL, mid, stat, now, err)
		return
	}
	return res.RowsAffected()
}

// Relation get relation between mid and fid.
func (d *Dao) Relation(c context.Context, mid, fid int64) (attr uint32, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getRelationSQL, hit(mid)), mid, fid)
	if err = row.Scan(&attr); err != nil {
		if err == sql.ErrNoRows {
			attr = model.AttrNoRelation
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// LoadMonitor load all mids into redis set.
func (d *Dao) LoadMonitor(c context.Context) (mids []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _loadMonitorSQL); err != nil {
		log.Error("d.Query.Exec(%s) error(%v)", _loadMonitorSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("row.Scan() error(%v)", err)
			mids = nil
			return
		}
		mids = append(mids, mid)
	}
	err = rows.Err()
	return
}

// AddMonitor add mid to monitor table
func (d *Dao) AddMonitor(c context.Context, mid int64, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addMonitorSQL, mid, now, now); err != nil {
		log.Error("d.AddMonitor.Exec(%s, %d) error(%v)", _addMonitorSQL, mid, err)
		return
	}
	return res.RowsAffected()
}

// DelMonitor del mid from monitor table
func (d *Dao) DelMonitor(c context.Context, mid int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _delMonitorSQL, mid); err != nil {
		log.Error("d.DelMonitor.Exec(%s, %d)", _delMonitorSQL, mid)
		return
	}
	return res.RowsAffected()
}

// TxDelTagUser delete tag user record.
func (d *Dao) TxDelTagUser(c context.Context, tx *sql.Tx, mid, fid int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_delTagUserSQL, tagUserHit(mid)), mid, fid); err != nil {
		log.Error("tx.Exec(%d, %d) error(%v)", mid, fid, err)
		return
	}
	return res.RowsAffected()
}

// Tags get tags list.
func (d *Dao) Tags(c context.Context, mid int64) (res map[int64]*model.Tag, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getTagsSQL, tagHit(mid)), mid); err != nil {
		log.Error("d.getTagsStmt[%d].Query(%d) error(%sv)", tagHit(mid), mid, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Tag)
	for rows.Next() {
		r := new(model.Tag)
		if err = rows.Scan(&r.Id, &r.Name, &r.Status, &r.MTime); err != nil {
			log.Error("d.getTagsStmt[%d].Query(%d) row.Scan() error(%v)", tagHit(mid), mid, err)
			res = nil
			return
		}
		res[r.Id] = r
	}
	res[0] = &model.Tag{Id: 0, Name: _defaultTag}
	res[-10] = &model.Tag{Id: -10, Name: _specialTag}
	err = rows.Err()
	return
}

// AddTag add tag.
func (d *Dao) AddTag(c context.Context, mid, fid int64, tag string, now time.Time) (lastID int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_addTagSQL, tagHit(mid)), mid, tag, now, now); err != nil {
		log.Error("d.db.Exec(%s, %d, %d, %s) error(%v)", _addTagSQL, mid, fid, tag, err)
		return
	}
	return res.LastInsertId()
}

// DelTag del tag.
func (d *Dao) DelTag(c context.Context, mid, id int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_delTagSQL, tagHit(mid)), id, mid); err != nil {
		log.Error("d.db.Exec(%s, %d, %d) error(%v)", _delTagSQL, mid, id, err)
		return
	}
	return res.RowsAffected()
}

// SetTagName update tag name info.
func (d *Dao) SetTagName(c context.Context, id, mid int64, name string, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_setTagNameSQL, tagHit(mid)), name, now, id); err != nil {
		log.Error("d.db.Exec(%s, %s, %d, %d)", _setTagNameSQL, name, mid, id)
		return
	}
	return res.RowsAffected()
}

// TagUserByMidFid get tagIds by mid and fid.
func (d *Dao) TagUserByMidFid(c context.Context, mid, fid int64) (tag *model.TagUser, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getTagsByMidFidSQL, tagUserHit(mid)), mid, fid)
	var tids i64b.Int64Bytes
	tag = new(model.TagUser)
	if err = row.Scan(&tag.Fid, &tids, &tag.MTime); err != nil {
		if err == sql.ErrNoRows {
			tag = nil
			err = nil
		} else {
			log.Error("d.getTagStmt[%d].Query(%d) row.Scan() error(%v)", tagUserHit(mid), mid, err)
		}
		return
	}
	tag.Tag = tids
	return
}

// UsersTags users tag by fids.
func (d *Dao) UsersTags(c context.Context, mid int64, fid []int64) (tags map[int64]*model.TagUser, err error) {
	row, err := d.db.Query(c, fmt.Sprintf(_getUsersTagSQL, tagUserHit(mid), xstr.JoinInts(fid)), mid)
	if err != nil {
		return
	}
	defer row.Close()
	tags = make(map[int64]*model.TagUser)
	for row.Next() {
		tag := new(model.TagUser)
		var tids i64b.Int64Bytes
		if err = row.Scan(&tag.Fid, &tids); err != nil {
			return
		}
		tag.Tag = tids
		tags[tag.Fid] = tag
	}
	return
}

// UserTag user tag
func (d *Dao) UserTag(c context.Context, mid int64) (tags map[int64][]int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getTagUserSQL, tagUserHit(mid)), mid); err != nil {
		log.Error("d.Query[%d].Query(%d) error(%sv)", tagUserHit(mid), mid, err)
		return
	}
	defer rows.Close()
	tags = make(map[int64][]int64)
	for rows.Next() {
		var (
			tids i64b.Int64Bytes
			fid  int64
		)
		if err = rows.Scan(&fid, &tids); err != nil {
			log.Error("d.Scan[%d].Query(%d) row.Scan() error(%v)", tagUserHit(mid), mid, err)
			tags = nil
			return
		}
		tags[fid] = tids
	}
	err = rows.Err()
	return
}

// SetTagUser setTagUser info.
func (d *Dao) SetTagUser(c context.Context, mid, fid int64, tag i64b.Int64Bytes, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_setTagUserSQL, tagUserHit(mid)), tag.Bytes(), now, mid, fid); err != nil {
		log.Error("d.db.Exec(%s, %d, %d, %v) error (%v)", _setTagUserSQL, mid, fid, tag.Bytes(), err)
		return
	}
	return res.RowsAffected()
}

// AddTagUser update tag name info.
func (d *Dao) AddTagUser(c context.Context, mid, fid int64, tag []int64, now time.Time) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_addTagUserSQL, tagUserHit(mid)), mid, fid, i64b.Int64Bytes(tag), now, now, i64b.Int64Bytes(tag)); err != nil {
		log.Error("d.db.Exec(%s, %d, %d, %v) error (%v)", _addTagUserSQL, mid, fid, tag, err)
		return
	}
	return res.RowsAffected()
}

// HasReachAchieve is
func (d *Dao) HasReachAchieve(c context.Context, mid int64, achieve model.AchieveFlag) bool {
	row := d.db.QueryRow(c, _hasReachAchieve, mid, uint64(achieve))
	count := 0
	if err := row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			return false
		}
		log.Warn("Failed to check has reach achieve: mid: %d, achieve: %d, error: %+v", mid, achieve, err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

// FollowerNotifySetting get follower-notify setting
// 这里返回用户通知开关的状态（和数据库存储的状态值相反）
func (d *Dao) FollowerNotifySetting(c context.Context, mid int64) (bool, error) {
	row := d.db.QueryRow(c, _getFollowerNotifySettingSQL, mid)
	var disableFollowerNotify bool
	if err := row.Scan(&disableFollowerNotify); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan() error(%v)", err)
		}
		return true, nil
	}
	if disableFollowerNotify {
		return false, nil
	}
	return true, nil
}

// EnableFollowerNotify enable follower-notify setting
func (d *Dao) EnableFollowerNotify(c context.Context, mid int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _enableFollowerNotifySQL, mid); err != nil {
		log.Error("enable follower-notify: tx.Exec(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// DisableFollowerNotify disable follower-notify setting
func (d *Dao) DisableFollowerNotify(c context.Context, mid int64) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _disableFollowerNotifySQL, mid); err != nil {
		log.Error("diabel follower-notify: tx.Exec(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}
