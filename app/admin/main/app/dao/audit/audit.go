package audit

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/audit"
	"go-common/library/database/orm"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_initAuditKey = "audit_key_%v_%v"
	_insertSQL    = `INSERT INTO audit (build,remark,mobi_app,ctime,mtime) VALUES (?,?,?,?,?)`
	_updateSQL    = `UPDATE audit SET build=?,remark=?,mobi_app=?,mtime=? WHERE id=?`
	_delSQL       = `DELETE FROM audit WHERE id=?`
)

// Dao is audit dao.
type Dao struct {
	// db
	db *gorm.DB
}

// New new a audit dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: orm.NewMySQL(c.ORM.Show),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.db.LogMode(true)
}

// Audits select All
func (d *Dao) Audits(ctx context.Context) (res []*audit.Audit, err error) {
	if err := d.db.Order("id DESC", true).Find(&res).Error; err != nil {
		log.Error("d.db.OrderBy error(%v)", err)
		return nil, err
	}
	return
}

// AuditByID by id
func (d *Dao) AuditByID(ctx context.Context, id int64) (rows *audit.Audit, err error) {
	rows = &audit.Audit{}
	if err := d.db.Where("id = ?", id).First(rows).Error; err != nil {
		log.Error("d.db.Where error(%v)", err)
		return nil, err
	}
	return
}

// AuditExist by build mobi_app
func (d *Dao) AuditExist(ctx context.Context, a *audit.Param) (res map[string]*audit.Audit, err error) {
	res = map[string]*audit.Audit{}
	row := &audit.Audit{}
	if err := d.db.Where("mobi_app = ? AND build = ?", a.MobiApp, a.Build).First(row).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			log.Error("d.db.Where error(%v)", err)
			return nil, err
		}
	}
	key := fmt.Sprintf(_initAuditKey, row.MobiApp, row.Build)
	res[key] = row
	return
}

// Insert insert audit
func (d *Dao) Insert(ctx context.Context, a *audit.Param, now time.Time) (err error) {
	if err = d.db.Exec(_insertSQL, a.Build, a.Remark, a.MobiApp, now, now).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Update update audit
func (d *Dao) Update(ctx context.Context, a *audit.Param, now time.Time) (err error) {
	if err = d.db.Exec(_updateSQL, a.Build, a.Remark, a.MobiApp, now, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Del delete audit by id
func (d *Dao) Del(ctx context.Context, id int64) (err error) {
	if err = d.db.Exec(_delSQL, id).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// PingDB check connection of db , mc.
func (d *Dao) PingDB(c context.Context) (err error) {
	if d.db != nil {
		err = d.db.DB().PingContext(c)
	}
	return
}
