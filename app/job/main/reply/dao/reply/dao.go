package reply

import (
	"context"

	"go-common/app/job/main/reply/conf"
	"go-common/library/database/sql"
	"go-common/library/queue/databus"
)

// Dao define mysql info
type Dao struct {
	// memcache
	Mc *MemcacheDao
	// mysql
	mysql    *sql.DB
	Admin    *AdminDao
	Content  *ContentDao
	Report   *ReportDao
	Reply    *RpDao
	Subject  *SubjectDao
	Business *BusinessDao
	// redis
	Redis *RedisDao
	// databus
	eventBus *databus.Databus
}

// New new a db and return
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// memchache
		Mc: NewMemcacheDao(c.Memcache),
		// mysql
		mysql: sql.NewMySQL(c.MySQL.Reply),
		// redis
		Redis: NewRedisDao(c.Redis),
		// databus
		eventBus: databus.New(c.Databus.Event),
	}
	d.Admin = NewAdminDao(d.mysql)
	d.Content = NewContentDao(d.mysql)
	d.Reply = NewReplyDao(d.mysql)
	d.Report = NewReportDao(d.mysql)
	d.Subject = NewSubjectDao(d.mysql)
	d.Business = NewBusinessDao(d.mysql)
	return
}

// Ping check db is alive
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.mysql.Ping(c); err != nil {
		return
	}
	if err = d.Redis.Ping(c); err != nil {
		return
	}
	return d.Mc.Ping(c)
}

// Close close all db connection
func (d *Dao) Close() {
	if d.Mc.mc != nil {
		d.Mc.mc.Close()
	}
	if d.Redis.redis != nil {
		d.Redis.redis.Close()
	}
	d.mysql.Close()
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.mysql.Begin(c)
}
