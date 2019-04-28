package dao

import (
	"context"
	"strings"
	"time"

	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
)

const (
	_selEmojiPackageSQL    = "SELECT id,name,url,remark,sort,state from emoji_package order by sort"
	_addEmojiPackageSQL    = "INSERT INTO emoji_package (name, url, sort, state, remark, ctime, mtime ) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_upEmojiPackageSQL     = "UPDATE emoji_package SET `name`=?, url=?, remark=?, state=?, mtime=? WHERE id=?"
	_upEmojiPackageSortSQL = "UPDATE emoji_package SET sort=?, mtime=? where id=?"
	_delEmojiPackageSQL    = "DELETE from emoji_package where id=?"
)

// EmojiPackageList get all emoji_package
func (d *Dao) EmojiPackageList(c context.Context) (packs []*model.EmojiPackage, err error) {
	rows, err := d.db.Query(c, _selEmojiPackageSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		pack := &model.EmojiPackage{}
		if err = rows.Scan(&pack.ID, &pack.Name, &pack.URL, &pack.Remark, &pack.Sort, &pack.State); err != nil {
			return
		}
		packs = append(packs, pack)
	}
	err = rows.Err()
	return
}

// CreateEmojiPackage insert a emoji_package into db
func (d *Dao) CreateEmojiPackage(c context.Context, name string, url string, sort int32, remark string, state int32) (id int64, err error) {
	result, err := d.db.Exec(c, _addEmojiPackageSQL, name, url, sort, state, remark, time.Now(), time.Now())
	if err != nil {
		return
	}
	return result.LastInsertId()
}

// UpEmojiPackageSort udpate emojipack sort
func (d *Dao) UpEmojiPackageSort(tx *xsql.Tx, ids string) (err error) {
	idx := strings.Split(ids, ",")
	for sort, id := range idx {
		_, err = tx.Exec(_upEmojiPackageSortSQL, sort, time.Now(), id)
		if err != nil {
			return
		}
	}
	return nil
}

//UpEmojiPackage update emojipack by id
func (d *Dao) UpEmojiPackage(c context.Context, name string, url string, remark string, state int32, id int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _upEmojiPackageSQL, name, url, remark, state, time.Now(), id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

// DelEmojiPackage delete emoji_package by id
func (d *Dao) DelEmojiPackage(c context.Context, id int64) (idx int64, err error) {
	result, err := d.db.Exec(c, _delEmojiPackageSQL, id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}
