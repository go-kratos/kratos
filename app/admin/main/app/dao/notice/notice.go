package notice

import (
	"context"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/notice"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_updateStateSQL = `UPDATE notice SET state=?,mtime=? WHERE id=?`
	_updateBuildSQL = `UPDATE notice SET build=?,conditions=?,mtime=? WHERE id=?`
	_updateSQL      = `UPDATE notice SET plat=?,title=?,content=?,url=?,ef_time=?,ex_time=?,area=?,mtime=?,type=? WHERE id=?`
	_insertSQL      = `INSERT INTO notice (plat,title,content,url,ef_time,ex_time,build,conditions,area,state,ctime,mtime,type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
)

// Dao is notice dao.
type Dao struct {
	db *gorm.DB
}

// New new a notice dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: orm.NewMySQL(c.ORM.Show),
	}
	d.initORM()
	return
}

//initORM init
func (d *Dao) initORM() {
	d.db.LogMode(true)
}

// Notices select all
func (d *Dao) Notices(ctx context.Context) (res []*notice.Notice, err error) {
	if err = d.db.Order("id DESC", true).Find(&res).Error; err != nil {
		log.Error("d.db.Order error(%v)", err)
		return nil, err
	}
	return
}

// NoticeByID select by id
func (d *Dao) NoticeByID(ctx context.Context, id int64) (rows *notice.Notice, err error) {
	rows = &notice.Notice{}
	if err = d.db.Where("id = ?", id).First(rows).Error; err != nil {
		log.Error("d.db.Where error(%v)", err)
		return nil, err
	}
	return
}

// Insert insert notice
func (d *Dao) Insert(ctx context.Context, a *notice.Param, now time.Time) (err error) {
	if err = d.db.Exec(_insertSQL, a.Plat, a.Title, a.Content, a.URL, a.Eftime, a.Extime,
		a.Build, a.Conditions, a.Area, a.State, now, now, a.Type).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Update update notice
func (d *Dao) Update(ctx context.Context, a *notice.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateSQL, a.Plat, a.Title, a.Content, a.URL, a.Eftime, a.Extime,
		a.Area, now, a.Type, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// UpdateBuild update build and conditions
func (d *Dao) UpdateBuild(ctx context.Context, a *notice.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateBuildSQL, a.Build, a.Conditions, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// UpdateState update release or not
func (d *Dao) UpdateState(ctx context.Context, a *notice.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateStateSQL, a.State, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

//Close close connect
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
