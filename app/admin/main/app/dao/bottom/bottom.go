package bottom

import (
	"context"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/bottom"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_deleteSQL  = `DELETE FROM bottom_entry WHERE id=?`
	_upStateSQL = `UPDATE bottom_entry SET state=?,mtime=? WHERE id=?`
	_updateSQL  = `UPDATE bottom_entry SET name=?,logo=?,rank=?,action=?,param=?,mtime=? WHERE id=?`
	_insertSQL  = `INSERT INTO bottom_entry (name,logo,rank,action,param,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?)`
)

// Dao bottom dao
type Dao struct {
	db *gorm.DB
}

// New new a bottom dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: orm.NewMySQL(c.ORM.Show),
	}
	d.initORM()
	return
}

// initORM init
func (d *Dao) initORM() {
	d.db.LogMode(true)
}

// Bottoms select all
func (d *Dao) Bottoms(ctx context.Context) (res []*bottom.Bottom, err error) {
	if err = d.db.Order("rank DESC", true).Find(&res).Error; err != nil {
		log.Error("d.db.Order err(%v)", err)
		return nil, err
	}
	return
}

// BottomByID select bottom by id
func (d *Dao) BottomByID(ctx context.Context, id int64) (rows *bottom.Bottom, err error) {
	rows = &bottom.Bottom{}
	if err = d.db.Where("id = ?", id).First(rows).Error; err != nil {
		log.Error("d.db.Where err(%v)", err)
		return nil, err
	}
	return
}

// Insert insert
func (d *Dao) Insert(ctx context.Context, a *bottom.Param, now time.Time) (err error) {
	if err = d.db.Exec(_insertSQL, a.Name, a.Logo, a.Rank, a.Action, a.Param, a.State, now, now).Error; err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	return
}

// Update update
func (d *Dao) Update(ctx context.Context, a *bottom.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateSQL, a.Name, a.Logo, a.Rank, a.Action, a.Param, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	return
}

// UpdateByID update by id
func (d *Dao) UpdateByID(ctx context.Context, a *bottom.Param, now time.Time) (err error) {
	if err = d.db.Exec(_upStateSQL, a.State, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	return
}

// Delete delete by id
func (d *Dao) Delete(ctx context.Context, id int64) (err error) {
	if err = d.db.Exec(_deleteSQL, id).Error; err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	return
}

// Close close db
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
