package archive

import (
	"context"
	xsql "database/sql"
	"fmt"
	"reflect"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/database/sql"
	"testing"

	"github.com/bouk/monkey"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Archive(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Archive
	)
	Convey("Tool", t, WithDao(func(d *Dao) {
		sub, err = d.Archive(c, 23333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}
func Test_UpperArcStateMap(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("UpperArcStateMap", t, WithDao(func(d *Dao) {
		_, err = d.UpperArcStateMap(c, 23333)
		So(err, ShouldBeNil)
	}))
}

func Test_UpDelayRound(t *testing.T) {
	var (
		c        = context.TODO()
		err      error
		mt1, mt2 time.Time
	)
	Convey("UpDelayRound", t, WithDao(func(d *Dao) {
		_, err = d.UpDelayRound(c, mt1, mt2)
		So(err, ShouldBeNil)
	}))
}
func TestDao_TxUpState(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpState", t, func(ctx C) {
		_, err := d.TxUpState(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpAccess(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpAccess", t, func(ctx C) {
		_, err := d.TxUpAccess(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpRound(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpRound", t, func(ctx C) {
		_, err := d.TxUpRound(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpAttr(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpAttr", t, func(ctx C) {
		_, err := d.TxUpAttr(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpCover(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpCover", t, func(ctx C) {
		_, err := d.TxUpCover(tx, 2333, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_UpCover(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("UpCover", t, func(ctx C) {
		_, err := d.UpCover(c, 2333, "")
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpArcDuration(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpArcDuration", t, func(ctx C) {
		_, err := d.TxUpArcDuration(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpAttrBit(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpAttrBit", t, func(ctx C) {
		_, err := d.TxUpAttrBit(tx, 2333, 0, 1)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpPTime(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		tm    time.Time
	)
	Convey("TxUpPTime", t, func(ctx C) {
		_, err := d.TxUpPTime(tx, 2333, tm)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func Test_ArchiveAddict(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Addit
	)
	Convey("Addit", t, WithDao(func(d *Dao) {
		sub, err = d.Addit(c, 23333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_Delay(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Delay
	)
	Convey("Delay", t, WithDao(func(d *Dao) {
		sub, err = d.Delay(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_NowDelays(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub []*archive.Delay
	)
	Convey("NowDelays", t, WithDao(func(d *Dao) {
		sub, err = d.NowDelays(c, time.Now())
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_Forbid(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.ForbidAttr
	)
	Convey("Forbid", t, WithDao(func(d *Dao) {
		sub, err = d.Forbid(c, 23333)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	}))
}

func Test_TrackPassed(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int64
	)
	Convey("TrackPassed", t, WithDao(func(d *Dao) {
		sub, err = d.GetFirstPassByAID(c, 233)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

func Test_TypeMapping(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("TypeMapping", t, WithDao(func(d *Dao) {
		_, err = d.TypeMapping(c)
		So(err, ShouldBeNil)
	}))
}

/*func Test_AddRdsCovers(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		cvs = []*archive.Cover{
			&archive.Cover{
				Filename: "sssss",
			},
		}
	)
	Convey("AddRdsCovers", t, WithDao(func(d *Dao) {
		_, err = d.AddRdsCovers(c, cvs)
		So(err, ShouldBeNil)
	}))
}*/
func Test_DBus(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("DBus", t, WithDao(func(d *Dao) {
		_, err = d.DBus(c, "sss", "sssew", 888)
		So(err, ShouldBeNil)
	}))
}
func Test_UpDBus(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("UpDBus", t, WithDao(func(d *Dao) {
		_, err = d.UpDBus(c, "sss", "sssew", 888, 33)
		So(err, ShouldBeNil)
	}))
}

func Test_DelAdminDelay(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("DelAdminDelay", t, WithDao(func(d *Dao) {
		_, err = d.DelAdminDelay(c, 2333)
		So(err, ShouldBeNil)
	}))
}
func Test_DelDelayByIds(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("DelDelayByIds", t, WithDao(func(d *Dao) {
		_, err = d.DelDelayByIds(c, []int64{2333})
		So(err, ShouldBeNil)
	}))
}
func Test_FirstPassCount(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("FirstPassCount", t, WithDao(func(d *Dao) {
		_, err = d.FirstPassCount(c, []int64{2333})
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpForbid(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		af    = &archive.ForbidAttr{
			Aid: 2333,
		}
	)
	Convey("TxUpForbid", t, func(ctx C) {
		_, err := d.TxUpForbid(tx, af)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoXState(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoXState", t, func(ctx C) {
		_, err := d.TxUpVideoXState(tx, "sssss", 4)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoPlayurl(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoPlayurl", t, func(ctx C) {
		_, err := d.TxUpVideoPlayurl(tx, "sssss", "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVDuration(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVDuration", t, func(ctx C) {
		_, err := d.TxUpVDuration(tx, "sssss", 4)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoFilesize(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoFilesize", t, func(ctx C) {
		_, err := d.TxUpVideoFilesize(tx, "sssss", 4)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoResolutionsAndDimensions(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoResolutionsAndDimensions", t, func(ctx C) {
		_, err := d.TxUpVideoResolutionsAndDimensions(tx, "sssss", "", "1280,720,0")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoFailCode(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoFailCode", t, func(ctx C) {
		_, err := d.TxUpVideoFailCode(tx, "sssss", -2)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpRelationStatus(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpRelationStatus", t, func(ctx C) {
		_, err := d.TxUpRelationStatus(tx, 2333, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_NewVideos(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("NewVideos", t, func(ctx C) {
		_, err := d.NewVideos(c, 2333)
		So(err, ShouldBeNil)
	})
}

func TestDao_ValidAidByCid(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("ValidAidByCid", t, func(ctx C) {
		_, err := d.ValidAidByCid(c, 2333)
		So(err, ShouldBeNil)
	})
}
func TestDao_Stat(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("Stat", t, func(ctx C) {
		_, err := d.Stat(c, 2333)
		So(err, ShouldBeNil)
	})
}
func TestDao_TypeNaming(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("TypeNaming", t, func(ctx C) {
		_, err := d.TypeNaming(c)
		So(err, ShouldBeNil)
	})
}
func TestDao_VdoWithArcCntCapable(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("VdoWithArcCntCapable", t, func(ctx C) {
		_, err := d.VdoWithArcCntCapable(c, 2333)
		So(err, ShouldBeNil)
	})
}
func TestDao_VideoCountCapable(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("VideoCountCapable", t, func(ctx C) {
		_, err := d.VideoCountCapable(c, 2333)
		So(err, ShouldBeNil)
	})
}

func TestDao_Videos(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("Videos", t, func(ctx C) {
		_, err := d.Videos(c, 2333)
		So(err, ShouldBeNil)
	})
}

func TestDao_Video(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("Video", t, func(ctx C) {
		_, err := d.Video(c, "sssss")
		So(err, ShouldBeNil)
	})
}
func TestDao_TranVideoOper(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TranVideoOper", t, func(ctx C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
			func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
				return nil, fmt.Errorf("tx.Exec Error")
			})
		defer guard.Unpatch()
		_, err := d.TranVideoOper(c, tx, 2333, 111, 0, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

/*func TestDao_AICover(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("AICover", t, func(ctx C) {
		_, err := d.AICover(c, "sssss")
		So(err, ShouldBeNil)
	})
}*/
