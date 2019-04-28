package reply

import (
	"context"

	"go-common/app/interface/main/reply/conf"
	"go-common/library/database/sql"
)

// Dao Dao
type Dao struct {
	// memcache
	Mc *MemcacheDao
	// mysql
	mysql       *sql.DB
	dbSlave     *sql.DB
	Admin       *AdminDao
	Content     *ContentDao
	Report      *ReportDao
	Reply       *RpDao
	Captcha     *CaptchaDao
	Notice      *NoticeDao
	CreditUser  *CreditUserDao
	Subject     *SubjectDao
	Config      *ConfigDao
	BlockStatus *BlockStatusDao
	Business    *BusinessDao
	// redis
	Redis *RedisDao
	// kafka
	Databus *DatabusDao
	Emoji   *EmoDao
}

// New New
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// memchache
		Mc: NewMemcacheDao(c.Memcache),
		// mysql
		mysql:   sql.NewMySQL(c.MySQL.Reply),
		dbSlave: sql.NewMySQL(c.MySQL.ReplySlave),
		// redis
		Redis: NewRedisDao(c.Redis),

		Databus: NewDatabusDao(c.Databus),
	}
	d.Admin = NewAdminDao(d.mysql)
	d.Content = NewContentDao(d.mysql, d.dbSlave)
	d.Reply = NewReplyDao(d.mysql, d.dbSlave)
	d.Report = NewReportDao(d.mysql)
	d.Subject = NewSubjectDao(d.mysql)
	d.Notice = NewNoticeDao(d.mysql)
	d.Config = NewConfigDao(d.mysql)
	d.Captcha = NewCaptchaDao(c.HTTPClient)
	d.CreditUser = NewCreditDao(c)
	d.BlockStatus = NewBlockStatusDao(c)
	d.Business = NewBusinessDao(d.mysql)
	d.Emoji = NewEmojiDao(d.mysql)
	return
}

// Ping Ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.mysql.Ping(c); err != nil {
		return
	}
	if err = d.Redis.Ping(c); err != nil {
		return
	}
	return d.Mc.Ping(c)
}

// Close Close
func (d *Dao) Close() {
	if d.Mc.mc != nil {
		d.Mc.mc.Close()
	}
	if d.Redis.redis != nil {
		d.Redis.redis.Close()
	}
	d.mysql.Close()
}
