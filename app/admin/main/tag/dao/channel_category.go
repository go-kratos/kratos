package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_setChannelCategoryString = "(%d,%q,%d,%d,%d,%q)"

	_countChannelCategorySQL   = "SELECT count(id) FROM channel_type"
	_channelCategorySQL        = "SELECT id,name,`order`,state,attr,ctime,mtime FROM channel_type WHERE id=?"
	_channelCategoryByNameSQL  = "SELECT id,name,`order`,state,attr,ctime,mtime FROM channel_type WHERE name=?"
	_channelCategorieSQL       = "SELECT id,name,`order`,state,attr,ctime,mtime FROM channel_type"
	_insertChannelCategorySQL  = "INSERT INTO channel_type(name,`order`,state,attr) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE `order`=?,state=?,attr=?"
	_upChannelCategoryStateSQL = "UPDATE channel_type SET state=? WHERE id=?"
	_upChannelCategorySQL      = "UPDATE channel_type SET attr=? WHERE id=?"
	_setChannelCategorieSQL    = "REPLACE INTO channel_type (id,name,`order`,state,attr,ctime) VALUES %s"
)

func setChannelCategory(categories []*model.ChannelCategory) string {
	var sqls []string
	for _, v := range categories {
		sqls = append(sqls, fmt.Sprintf(_setChannelCategoryString, v.ID, v.Name, v.Order, v.State, v.Attr, v.CTime.Time().Format(_timeLayoutFormat)))
	}
	return strings.Join(sqls, ",")
}

// CountChannelCategory count channel category.
func (d *Dao) CountChannelCategory(c context.Context) (count int32, err error) {
	row := d.db.QueryRow(c, _countChannelCategorySQL)
	if err = row.Scan(&count); err != nil {
		log.Error("d.dao.CountChannelCategory() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ChannelCategory get a channel category.
func (d *Dao) ChannelCategory(c context.Context, id int64) (res *model.ChannelCategory, err error) {
	res = new(model.ChannelCategory)
	row := d.db.QueryRow(c, _channelCategorySQL, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Order, &res.State, &res.Attr, &res.CTime, &res.MTime); err != nil {
		log.Error("d.dao.ChannelCategory(%d) error(%v)", id, err)
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
	}
	return
}

// ChannelCategoryByName get a channel category by name.
func (d *Dao) ChannelCategoryByName(c context.Context, name string) (res *model.ChannelCategory, err error) {
	res = new(model.ChannelCategory)
	row := d.db.QueryRow(c, _channelCategoryByNameSQL, name)
	if err = row.Scan(&res.ID, &res.Name, &res.Order, &res.State, &res.Attr, &res.CTime, &res.MTime); err != nil {
		log.Error("d.dao.ChannelCategory(%s) error(%v)", name, err)
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
	}
	return
}

// ChannelCategories get all channel category.
func (d *Dao) ChannelCategories(c context.Context) (res []*model.ChannelCategory, err error) {
	rows, err := d.db.Query(c, _channelCategorieSQL)
	if err != nil {
		log.Error("d.dao.ChannelCategories() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.ChannelCategory, 0)
	for rows.Next() {
		t := &model.ChannelCategory{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Order, &t.State, &t.Attr, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelCategories scan() error(%v)", err)
			return
		}
		res = append(res, t)
	}
	return
}

// InsertChannelCategory insert channel category.
func (d *Dao) InsertChannelCategory(c context.Context, category *model.ChannelCategory) (id int64, err error) {
	res, err := d.db.Exec(c, _insertChannelCategorySQL, category.Name, category.Order, category.State, category.Attr, category.Order, category.State, category.Attr)
	if err != nil {
		log.Error("d.dao.InsertChannelCategory(%v) error(%v)", category, err)
		return
	}
	return res.LastInsertId()
}

// StateChannelCategory state channel category.
func (d *Dao) StateChannelCategory(c context.Context, id int64, state int32) (affect int64, err error) {
	res, err := d.db.Exec(c, _upChannelCategoryStateSQL, state, id)
	if err != nil {
		log.Error("d.dao.StateChannelCategory(%d,%d) error(%v)", id, state, err)
		return
	}
	return res.RowsAffected()
}

// UpdateChannelCategories update channel categories .
func (d *Dao) UpdateChannelCategories(c context.Context, categories []*model.ChannelCategory) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_setChannelCategorieSQL, setChannelCategory(categories)))
	if err != nil {
		log.Error("d.dao.UpdateChannelCategories(%v) error(%v)", categories, err)
		return
	}
	return res.RowsAffected()
}

// TxUpChannelCategoryAttr .
func (d *Dao) TxUpChannelCategoryAttr(tx *sql.Tx, category *model.ChannelCategory) (affect int64, err error) {
	res, err := tx.Exec(_upChannelCategorySQL, category.Attr, category.ID)
	if err != nil {
		log.Error("d.TxUpChannelCategoryAttr(%v) error(%v)", category, err)
		return
	}
	return res.RowsAffected()
}
