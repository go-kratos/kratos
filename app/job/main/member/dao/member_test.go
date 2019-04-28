package dao

import (
	"context"
	"database/sql"
	"fmt"

	"reflect"
	"testing"
	"time"

	"go-common/app/job/main/member/model"
	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_SetBaseInfo(t *testing.T) {
	r := &model.BaseInfo{Mid: 100, Sex: 1, Name: "lgs", Face: "face", Sign: "sign", Rank: 3}
	Convey("SetBaseInfo", t, func() {
		Convey("SetBaseInfo success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetBaseInfo(context.TODO(), r)
			So(err, ShouldBeNil)
		})
		Convey("SetBaseInfo failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetBaseInfo(context.TODO(), r)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_UpdateSign(t *testing.T) {
	Convey("SetSign", t, func() {
		Convey("SetSign success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetSign(context.TODO(), 1, "my sign")
			So(err, ShouldBeNil)
		})
		Convey("SetSign failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetSign(context.TODO(), 1, "my sign")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_SetName(t *testing.T) {
	Convey("SetName", t, func() {
		Convey("SetName success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetName(context.TODO(), 1, "myname")
			So(err, ShouldBeNil)
		})
		Convey("SetName failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetName(context.TODO(), 1, "myname")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_UpdateFace(t *testing.T) {
	Convey("SetFaceRank", t, func() {
		Convey("SetFaceRank success", func() {
			err := d.SetFace(context.TODO(), int64(1), "http://face")
			So(err, ShouldBeNil)
		})
		Convey("SetFaceRank failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetFace(context.TODO(), int64(1), "http://face")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_SelMoral(t *testing.T) {
	Convey("SelMoral", t, func() {
		Convey("SelMoral success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelMoral, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return nil
			})
			defer monkey.UnpatchAll()
			moral, err := d.SelMoral(context.TODO(), int64(8))
			So(err, ShouldBeNil)
			So(moral, ShouldNotBeNil)
		})
		Convey("SelMoral error", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelMoral, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return fmt.Errorf("row.Scan error")
			})
			defer monkey.UnpatchAll()
			moral, err := d.SelMoral(context.TODO(), int64(8))
			So(err, ShouldNotBeNil)
			So(moral, ShouldNotBeNil)
		})
		Convey("SelMoral no record", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelMoral, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return sql.ErrNoRows
			})
			defer monkey.UnpatchAll()
			moral, err := d.SelMoral(context.TODO(), int64(8))
			So(err, ShouldBeNil)
			So(moral, ShouldBeNil)
		})
	})
}

func TestDao_IncrMoral(t *testing.T) {
	Convey("IncrMoral", t, func() {
		Convey("AuditQueuingFace success", func() {
			_, err := d.RecoverMoral(context.TODO(), 2, 100, 1, time.Now().Format("2006-01-02"))
			So(err, ShouldBeNil)
		})
		Convey("AuditQueuingFace failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			_, err := d.RecoverMoral(context.TODO(), 2, 100, 1, time.Now().Format("2006-01-02"))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_AuditQueuingFace(t *testing.T) {
	Convey("AuditQueuingFace", t, func() {
		Convey("AuditQueuingFace success", func() {
			err := d.AuditQueuingFace(context.TODO(), "/bfs/face/7e68723b9d3664ac3773c1f3c26d5e2bfabc0f22.jpg", "", 0)
			So(err, ShouldBeNil)
		})
		Convey("AuditQueuingFace failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.AuditQueuingFace(context.TODO(), "/bfs/face/7e68723b9d3664ac3773c1f3c26d5e2bfabc0f22.jpg", "", 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_BaseInfo(t *testing.T) {
	Convey("BaseInfo", t, func() {
		Convey("BaseInfo success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelBaseInfo, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return nil
			})
			defer monkey.UnpatchAll()
			r, err := d.BaseInfo(context.TODO(), 100)
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)
		})
		Convey("BaseInfo error", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelBaseInfo, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return fmt.Errorf("row.Scan error")
			})
			defer monkey.UnpatchAll()
			r, err := d.BaseInfo(context.TODO(), 11)
			So(err, ShouldNotBeNil)
			So(r, ShouldBeNil)
		})
		Convey("BaseInfo no record", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.accdb.QueryRow(context.TODO(), _SelBaseInfo, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return sql.ErrNoRows
			})
			defer monkey.UnpatchAll()
			r, err := d.BaseInfo(context.TODO(), 11)
			So(err, ShouldBeNil)
			So(r, ShouldBeNil)
		})
	})

}

// TestDao_RecoverMoral  update
func TestDao_RecoverMoral(t *testing.T) {
	Convey("RecoverMoral", t, func() {
		Convey("RecoverMoral success", func() {
			_, err := d.RecoverMoral(context.TODO(), 121, 1, 1, time.Now().Format("2006-01-02"))
			So(err, ShouldBeNil)
		})
		Convey("RecoverMoral failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			r, err := d.RecoverMoral(context.TODO(), 121, 1, 1, time.Now().Format("2006-01-02"))
			So(err, ShouldNotBeNil)
			So(r, ShouldEqual, 0)
		})
	})
}

func TestDao_SetSign(t *testing.T) {
	Convey("SetSign", t, func() {
		Convey("SetSign success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetSign(context.TODO(), 100, fmt.Sprintf("sign%s", time.Minute.String()))
			So(err, ShouldBeNil)
		})
		Convey("SetSign failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetSign(context.TODO(), 100, fmt.Sprintf("sign%s", time.Minute.String()))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_SetOfficial(t *testing.T) {
	Convey("SetOfficial", t, func() {
		Convey("SetOfficial success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetOfficial(context.TODO(), 100, 1, fmt.Sprintf("title%s", time.Minute.String()), "it's a test")
			So(err, ShouldBeNil)
		})
		Convey("SetOfficial failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetOfficial(context.TODO(), 100, 1, fmt.Sprintf("title%s", time.Minute.String()), "it's a test")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_SetFace(t *testing.T) {
	Convey("SetFace", t, func() {
		err := d.SetFace(context.TODO(), 100, fmt.Sprintf("face%s", time.Minute.String()))
		So(err, ShouldBeNil)
	})

	Convey("SetFace", t, func() {
		Convey("SetFace success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.SetFace(context.TODO(), 100, fmt.Sprintf("face%s", time.Minute.String()))
			So(err, ShouldBeNil)
		})
		Convey("SetFace failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.SetFace(context.TODO(), 100, fmt.Sprintf("face%s", time.Minute.String()))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDao_InitBase(t *testing.T) {
	Convey("InitBase", t, func() {
		Convey("InitBase success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.InitBase(context.TODO(), 1)
			So(err, ShouldBeNil)
		})
		Convey("InitBase failed", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("cannot connect database error")
			})
			defer monkey.UnpatchAll()
			err := d.InitBase(context.TODO(), 1)
			So(err, ShouldNotBeNil)
		})
	})
}
