package audit

import (
	"fmt"
	"testing"
	"time"

	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

func getCID(db *sql.DB, isPGC bool, double bool) (cid int64, err error) {
	var query string
	if isPGC {
		if double {
			query = "SELECT cid FROM tv_content WHERE is_deleted = 0 GROUP BY cid HAVING COUNT(1) >1 LIMIT 1"
		} else {
			query = "SELECT cid FROM tv_content WHERE is_deleted = 0 LIMIT 1"
		}
	} else {
		if double {
			query = "SELECT cid FROM ugc_video WHERE deleted = 0 GROUP BY cid HAVING COUNT(1) > 1 LIMIT 1"
		} else {
			query = "SELECT cid FROM ugc_video WHERE deleted = 0 LIMIT 1"
		}
	}
	if err = db.QueryRow(ctx, query).Scan(&cid); err != nil {
		fmt.Println("No Ready Video : ", err)
		return
	}
	fmt.Println("pick CID ", cid)
	return
}

func tryDouble(isPGC bool) (cid int64, err error) {
	var errDouble error
	if cid, errDouble = getCID(d.db, isPGC, true); errDouble != nil {
		fmt.Println("Double Err ", errDouble, " Will Use Single")
		if cid, err = getCID(d.db, isPGC, false); err != nil {
			fmt.Println("Single Err ", err, " Will Fail")
			return
		}
	}
	return
}

func TestDao_UgcCID(t *testing.T) {
	Convey("TestDao_UgcCID", t, WithDao(func(d *Dao) {
		cid, errTry := tryDouble(false)
		if errTry != nil {
			fmt.Println("No Cid Ready")
			return
		}
		cont, err := d.UgcCID(ctx, cid)
		So(err, ShouldBeNil)
		So(len(cont), ShouldBeGreaterThan, 0)
		fmt.Println(cont)
	}))
}

func TestDao_PgcCID(t *testing.T) {
	Convey("TestDao_PgcCID", t, WithDao(func(d *Dao) {
		cid, errTry := tryDouble(true)
		if errTry != nil {
			fmt.Println("No Cid Ready")
			return
		}
		cont, err := d.PgcCID(ctx, cid)
		So(err, ShouldBeNil)
		So(len(cont), ShouldBeGreaterThan, 0)
		fmt.Println(cont)
	}))
}

func TestDao_UgcTranscode(t *testing.T) {
	Convey("TestDao_UgcTranscode", t, WithDao(func(d *Dao) {
		cid, errTry := tryDouble(false)
		if errTry != nil {
			fmt.Println("No Cid Ready")
			return
		}
		cont, _ := d.UgcCID(ctx, cid)
		if len(cont) == 0 {
			fmt.Println("UgcCid Error!")
			return
		}
		fmt.Println("To Update ", cont)
		err := d.UgcTranscode(ctx, cont, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_PgcTranscode(t *testing.T) {
	Convey("TestDao_PgcTranscode", t, WithDao(func(d *Dao) {
		cid, errTry := tryDouble(true)
		if errTry != nil {
			fmt.Println("No Cid Ready")
			return
		}
		cont, _ := d.PgcCID(ctx, cid)
		if len(cont) == 0 {
			fmt.Println("PgcCid Error!")
			return
		}
		fmt.Println("To Update ", cont)
		err := d.PgcTranscode(ctx, cont, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_ApplyPGC(t *testing.T) {
	Convey("TestDao_ApplyPGC", t, WithDao(func(d *Dao) {
		cid, errTry := tryDouble(true)
		if errTry != nil {
			fmt.Println("No Cid Ready")
			return
		}
		cont, _ := d.PgcCID(ctx, cid)
		if len(cont) == 0 {
			fmt.Println("PgcCid Error!")
			return
		}
		fmt.Println("To Update ", cont)
		err := d.ApplyPGC(ctx, cont, time.Now().Unix())
		So(err, ShouldBeNil)
	}))
}
