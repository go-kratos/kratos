package wall

import (
	"context"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/wall"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_upstateSQL = `UPDATE wall SET state=?,mtime=? WHERE id=?`
	_insertSQL  = `INSERT INTO wall (name,title,logo,package,size,download,remark,rank,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_updateSQL  = `UPDATE wall SET name=?,title=?,logo=?,package=?,size=?,download=?,remark=?,rank=?,state=?,mtime=? WHERE id=?`
)

// Dao wall dao
type Dao struct {
	db *gorm.DB
}

// New new a wall dao
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

// Walls select all
func (d *Dao) Walls(ctx context.Context) (res []*wall.Wall, err error) {
	if err = d.db.Order("id DESC", true).Find(&res).Error; err != nil {
		log.Error("d.db.Order error(%v)", err)
		return nil, err
	}
	return
}

// WallByID appwall by id
func (d *Dao) WallByID(ctx context.Context, id int64) (rows *wall.Wall, err error) {
	rows = &wall.Wall{}
	if err = d.db.Where("id = ?", id).First(rows).Error; err != nil {
		log.Error("d.db.Where error(%v)", err)
		return nil, err
	}
	return
}

// Insert appwall date
func (d *Dao) Insert(ctx context.Context, a *wall.Param, now time.Time) (err error) {
	if err = d.db.Exec(_insertSQL, a.Name, a.Title, a.Logo, a.Package, a.Size, a.Download,
		a.Remark, a.Rank, a.State, now, now).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Update update appwall
func (d *Dao) Update(ctx context.Context, a *wall.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateSQL, a.Name, a.Title, a.Logo, a.Package, a.Size, a.Download,
		a.Remark, a.Rank, a.State, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// UpdateByID Update state by id
func (d *Dao) UpdateByID(ctx context.Context, a *wall.Param, now time.Time) (err error) {
	if err = d.db.Exec(_upstateSQL, a.State, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Close close connect
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
