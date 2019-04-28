package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/log"
)

// Available .
func Available(db *gorm.DB) *gorm.DB {
	return db.Where("disable_time = '0000-00-00 00:00:00'")
}

// Disable .
func Disable(db *gorm.DB) *gorm.DB {
	return db.Not("disable_time = '0000-00-00 00:00:00'")
}

func state(state string) func(db *gorm.DB) *gorm.DB {
	state = strings.TrimSpace(state)
	if state == net.StateAvailable {
		return func(db *gorm.DB) *gorm.DB {
			return Available(db)
		}
	} else if state == net.StateDisable {
		return func(db *gorm.DB) *gorm.DB {
			return Disable(db)
		}
	} else {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
}

func pager(ps int64, pn int64, sort string) func(db *gorm.DB) *gorm.DB {
	offset := ps * (pn - 1)
	order := fmt.Sprintf("id %s", sort)
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(ps).Offset(offset).Order(order)
	}
}

// UpdateFields .
func (d *Dao) UpdateFields(c context.Context, db *gorm.DB, table string, id int64, fields map[string]interface{}) (err error) {
	if db == nil {
		db = d.orm
	}
	if err = db.Table(table).Where("id=?", id).Updates(fields).Error; err != nil {
		log.Error("UpdateFlow(%s,%d) error(%v) changed(%+v)", table, id, err, fields)
	}
	return
}

// AddItem .
func (d *Dao) AddItem(c context.Context, db *gorm.DB, n interface{}) (err error) {
	if db == nil {
		db = d.orm
	}
	if err = db.Create(n).Error; err != nil {
		log.Error("AddItem error(%v) (%+v)", err, n)
	}
	return
}

// ColumnMapString .
func (d *Dao) ColumnMapString(c context.Context, table string, column string, ids []int64, where string) (result map[int64]string, err error) {
	var (
		rows  *sql.Rows
		id    int64
		value string
		slt   = fmt.Sprintf("`id`,`%s`", column)
	)

	result = map[int64]string{}
	db := d.orm.Table(table).Select(slt).Where("id in (?)", ids)
	where = strings.TrimSpace(where)
	if where != "" {
		db = db.Where(where)
	}
	rows, err = db.Rows()
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error("ColumnMapString(%s, %s) rows error(%v) ids(%v)", table, column, err, ids)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &value); err != nil {
			log.Error("ColumnMapString(%s, %s) rows.scan error(%v) ids(%v)", table, column, err, ids)
			return
		}
		result[id] = value
	}

	return
}
