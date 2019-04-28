package base

import (
	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/dao"
	bsn "go-common/app/job/main/search/dao/business"
)

// Base .
type Base struct {
	D *dao.Dao
	C *conf.Config
}

// NewBase .
func NewBase(c *conf.Config) (b *Base) {
	b = &Base{
		C: c,
		D: dao.New(c),
	}
	b.D.AppPool = b.newAppPool(b.D)
	return
}

// newAppPool .
func (b *Base) newAppPool(d *dao.Dao) (pool map[string]dao.App) {
	pool = make(map[string]dao.App)
	for k, v := range d.BusinessPool {
		switch v.IncrWay {
		case "single":
			pool[k] = dao.NewAppSingle(d, k)
		case "multiple":
			pool[k] = dao.NewAppMultiple(d, k)
		case "dtb":
			pool[k] = dao.NewAppDatabus(d, k)
		case "multipleDtb":
			pool[k] = dao.NewAppMultipleDatabus(d, k)
		case "business":
			switch k {
			case "archive_video":
				pool[k] = bsn.NewAvr(d, k, b.C)
			case "avr_archive", "avr_video":
				pool[k] = bsn.NewAvrArchive(d, k)
			case "log_audit", "log_user_action":
				pool[k] = bsn.NewLog(d, k)
			case "dm_date":
				pool[k] = bsn.NewDmDate(d, k)
			case "aegis_resource":
				pool[k] = bsn.NewAegisResource(d, k, b.C)
			}
		default:
			// to do other thing
		}
	}
	//fmt.Println("strace:app-pool>", pool)
	return
}
