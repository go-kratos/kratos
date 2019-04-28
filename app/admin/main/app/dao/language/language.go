package language

import (
	"context"
	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/language"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_insertSQL = `INSERT INTO language (name, remark) VALUES (?,?)`
	_updateSQL = `UPDATE language SET name=?,remark=? WHERE id=?`
)

// Dao is language dao
type Dao struct {
	db *gorm.DB
}

// New new a language dao
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

// Languages select all
func (d *Dao) Languages(ctx context.Context) (res []*language.Language, err error) {
	if err := d.db.Order("id DESC", true).Find(&res).Error; err != nil {
		log.Error("d.db.Order err(%v)", err)
		return nil, err
	}
	return
}

// LangByID select by id
func (d *Dao) LangByID(ctx context.Context, id int64) (rows *language.Language, err error) {
	rows = &language.Language{}
	if err = d.db.Where("id = ?", id).First(rows).Error; err != nil {
		log.Error("d.db.Where error(%v)", err)
		return nil, err
	}
	return
}

// Insert language date
func (d *Dao) Insert(ctx context.Context, a *language.Param) (err error) {
	if err = d.db.Exec(_insertSQL, a.Name, a.Remark).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// Update date
func (d *Dao) Update(ctx context.Context, a *language.Param) (err error) {
	if err = d.db.Exec(_updateSQL, a.Name, a.Remark, a.ID).Error; err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}
