package dao

import (
	"fmt"
	"testing"
	"time"

	"go-common/app/service/live/xlottery/model"
	xsql "go-common/library/database/sql"
	xtime "go-common/library/time"
)

var db *xsql.DB

func init() {
	c := &xsql.Config{
		Addr:         "172.16.38.117:3312",
		DSN:          "live:oWni@ElNs0P0C(dphdj*F1y4@tcp(172.16.38.117:3312)/live-app?timeout=2000ms&readTimeout=2000ms&writeTimeout=2000ms&parseTime=true&loc=Local&charset=utf8,utf8mb4",
		Active:       10,
		Idle:         5,
		IdleTimeout:  xtime.Duration(time.Minute),
		QueryTimeout: xtime.Duration(time.Minute),
		ExecTimeout:  xtime.Duration(time.Minute),
		TranTimeout:  xtime.Duration(time.Minute),
	}
	db = xsql.NewMySQL(c)
}
func TestDao_InsertSpecialGift(t *testing.T) {
	d := &Dao{
		db: db,
	}
	got, err := d.InsertSpecialGift(&model.SpecialGift{
		UID:        333,
		RoomID:     33,
		GiftID:     33,
		GiftNum:    1,
		CreateTime: time.Now(),
		CustomField: "{\"content\":\"真是无法无天了。\"	}",
	})
	fmt.Println(got, err)
}

func TestDao_FindBeatByBeatIDAndUID(t *testing.T) {
	d := &Dao{
		db: db,
	}
	got, err := d.FindBeatByBeatIDAndUID(123, 10666892)

	fmt.Println(got, err)
}

func TestDao_FindShieldKeyWorkByUID(t *testing.T) {
	d := &Dao{
		db: db,
	}
	got, err := d.FindShieldKeyWorkByUID(12209)
	fmt.Println(got, err)
}
