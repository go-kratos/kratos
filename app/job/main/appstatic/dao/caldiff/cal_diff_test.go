package caldiff

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/job/main/appstatic/model"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_availableRes = "SELECT id FROM resource ORDER BY id DESC"
)

func TestDao_DiffNew(t *testing.T) {
	Convey("TestDao_DiffNew", t, WithDao(func(d *Dao) {
		file, err := d.DiffNew(ctx)
		So(err, ShouldBeNil)
		data, err2 := (json.Marshal(file))
		So(err2, ShouldBeNil)
		fmt.Println(string(data))
	}))
}

func TestDao_DiffRetry(t *testing.T) {
	Convey("TestDao_DiffRetry", t, WithDao(func(d *Dao) {
		file, err := d.DiffRetry(ctx)
		So(err, ShouldBeNil)
		data, err2 := (json.Marshal(file))
		So(err2, ShouldBeNil)
		fmt.Println(string(data))
	}))
}

func TestDao_SaveFile(t *testing.T) {
	Convey("TestDao_SaveFile", t, WithDao(func(d *Dao) {
		err := d.SaveFile(ctx, 1, &model.FileInfo{
			Name: "123",
			Size: 123,
			Type: "1",
			Md5:  "1234",
			URL:  "xxx",
		})
		So(err, ShouldBeNil)
	}))
}

func TestDao_ParseResID(t *testing.T) {
	Convey("TestDao_ParseResID", t, WithDao(func(d *Dao) {
		var r = &model.Resource{}
		if err := d.db.QueryRow(ctx, _availableRes).Scan(&r.ID); err != nil {
			return
		}
		res, err := d.ParseResID(ctx, int(r.ID))
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		fmt.Println(r.ID)
	}))
}

func TestDao_ParseResVer(t *testing.T) {
	Convey("TestDao_ParseResVer", t, WithDao(func(d *Dao) {
		dd, err := d.ParseResVer(ctx, 23, 1)
		So(err, ShouldBeNil)
		fmt.Println(err)
		So(dd, ShouldNotBeNil)
		data, err2 := (json.Marshal(dd))
		So(err2, ShouldBeNil)
		fmt.Println(string(data))
	}))
}

func TestDao_ReadyFile(t *testing.T) {
	Convey("TestDao_ReadyFile", t, WithDao(func(d *Dao) {
		var r = &model.Resource{}
		if err := d.db.QueryRow(ctx, _availableRes).Scan(&r.ID); err != nil {
			return
		}
		dd, err := d.ReadyFile(ctx, int(r.ID), 0) // ftype = 0 full package
		So(err, ShouldBeNil)
		So(dd, ShouldNotBeNil)
		data, _ := json.Marshal(dd)
		fmt.Println(string(data))
	}))
}

func TestDao_UpdateStatus(t *testing.T) {
	Convey("TestDao_UpdateStatus", t, WithDao(func(d *Dao) {
		err := d.UpdateStatus(ctx, 2, 304)
		So(err, ShouldBeNil)
	}))
}
