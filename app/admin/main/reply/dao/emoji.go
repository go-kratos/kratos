package dao

import (
	"context"
	"strings"
	"time"

	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
)

const (
	//state 0-上线 1-下线
	_selEmojiSQL       = "SELECT id,package_id,name,url,sort,state,remark from emoji order by package_id,sort"
	_selEmojiByPidSQL  = "SELECT id,package_id,name,url,sort,state,remark from emoji where package_id=? order by sort"
	_selEmojiByNameSQL = "SELECT id from emoji where name=? "
	_insertEmojiSQL    = "INSERT INTO emoji (package_id, name, url, sort, state, remark, ctime, mtime ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_upEmojiSortSQL    = "UPDATE emoji SET sort=?, mtime=? where id=?"
	_upEmojiStateSQL   = "UPDATE emoji SET state=?, mtime=? WHERE id=?"
	_upEmojiSQL        = "UPDATE emoji SET name=?, remark=?, url=?, mtime=? WHERE id=?"
	_delEmojiByIDSQL   = "DELETE from emoji where id=?"
	_delEmojiByPidSQL  = "DELETE from emoji where package_id=?"
)

// EmojiList get all emoji
func (d *Dao) EmojiList(c context.Context) (emojis []*model.Emoji, err error) {
	rows, err := d.db.Query(c, _selEmojiSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		emo := &model.Emoji{}
		if err = rows.Scan(&emo.ID, &emo.PackageID, &emo.Name, &emo.URL, &emo.Sort, &emo.State, &emo.Remark); err != nil {
			return
		}
		emojis = append(emojis, emo)
	}
	err = rows.Err()
	return
}

// EmojiListByPid get emoji by package_id
func (d *Dao) EmojiListByPid(c context.Context, pid int64) (emojis []*model.Emoji, err error) {
	rows, err := d.db.Query(c, _selEmojiByPidSQL, pid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		emo := &model.Emoji{}
		if err = rows.Scan(&emo.ID, &emo.PackageID, &emo.Name, &emo.URL, &emo.Sort, &emo.State, &emo.Remark); err != nil {
			return
		}
		emojis = append(emojis, emo)
	}
	err = rows.Err()
	return
}

// EmojiByName get emoji by name
func (d *Dao) EmojiByName(c context.Context, name string) (emojis []*model.Emoji, err error) {
	rows, err := d.db.Query(c, _selEmojiByNameSQL, name)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		emo := &model.Emoji{}
		if err = rows.Scan(&emo.ID); err != nil {
			return
		}
		emojis = append(emojis, emo)
	}
	err = rows.Err()
	return
}

// CreateEmoji insert a emoji into db
func (d *Dao) CreateEmoji(c context.Context, pid int64, name string, url string, sort int32, state int32, remark string) (id int64, err error) {
	result, err := d.db.Exec(c, _insertEmojiSQL, pid, name, url, sort, state, remark, time.Now(), time.Now())
	if err != nil {
		return
	}
	return result.LastInsertId()
}

// UpEmojiSort udpate emoji sort
func (d *Dao) UpEmojiSort(tx *xsql.Tx, ids string) (err error) {
	idx := strings.Split(ids, ",")
	for sort, id := range idx {
		_, err = tx.Exec(_upEmojiSortSQL, sort, time.Now(), id)
		if err != nil {
			return
		}
	}
	return nil
}

// UpEmojiStateByID update emoji state
func (d *Dao) UpEmojiStateByID(c context.Context, state int32, id int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _upEmojiStateSQL, state, time.Now(), id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

// UpEmoji update emoji name and remark by id
func (d *Dao) UpEmoji(c context.Context, name string, remark string, url string, id int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _upEmojiSQL, name, remark, url, time.Now(), id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

// DelEmojiByID delete emoji by id
func (d *Dao) DelEmojiByID(c context.Context, id int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _delEmojiByIDSQL, id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

// DelEmojiByPid delete emoji by package_id
func (d *Dao) DelEmojiByPid(c context.Context, pid int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _delEmojiByPidSQL, pid)
	if err != nil {
		return
	}
	return result.RowsAffected()
}
