package dao

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/push/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Task(t *testing.T) {
	Convey("add task", t, WithDao(func(d *Dao) {
		t := &model.Task{
			Job:         1378138219873,
			Type:        model.TaskTypeBusiness,
			APPID:       1,
			BusinessID:  1,
			Platform:    []int{1, 2, 3},
			Title:       "test_tile",
			Summary:     "test_summary",
			LinkType:    model.LinkTypeBrowser,
			LinkValue:   "https://www.bilibili.com",
			Build:       map[int]*model.Build{2: {Build: 100, Condition: "gt"}},
			Sound:       1,
			Vibration:   1,
			PassThrough: 1,
			Progress:    new(model.Progress),
			MidFile:     "xxx.txt",
			PushTime:    xtime.Time(1500000000),
			ExpireTime:  xtime.Time(1600000000),
			Status:      model.TaskStatusPrepared,
		}
		c := context.Background()
		taskID, err := d.AddTask(c, t)
		taskIDString := strconv.FormatInt(taskID, 10)
		So(err, ShouldBeNil)

		Convey("task info", func() {
			task, err := d.Task(c, taskIDString)
			So(err, ShouldBeNil)
			task.ID = ""
			So(task, ShouldResemble, t)
		})

		Convey("update task progress", func() {
			p := &model.Progress{TokenTotal: 100}
			err := d.UpdateTaskProgress(c, taskIDString, p)
			So(err, ShouldBeNil)

			task, err := d.Task(c, taskIDString)
			So(err, ShouldBeNil)
			So(task, ShouldNotBeEmpty)
			So(task.Progress.TokenTotal, ShouldEqual, 100)
		})

		Convey("update task status", func() {
			err := d.UpdateTaskStatus(c, taskIDString, model.TaskStatusDone)
			So(err, ShouldBeNil)

			task, err := d.Task(c, taskIDString)
			So(err, ShouldBeNil)
			So(task, ShouldNotBeEmpty)
			So(task.Status, ShouldEqual, model.TaskStatusDone)
		})

		Convey("tx tokens by platform", func() {
			tx, _ := d.BeginTx(context.Background())
			_, err := d.TxTaskByPlatform(tx, model.PlatformIPad)
			So(err, ShouldBeNil)
			err = tx.Commit()
			So(err, ShouldBeNil)
		})
	}))
}

func Test_Business(t *testing.T) {
	Convey("get businesses", t, WithDao(func(d *Dao) {
		res, err := d.Businesses(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		fmt.Println(res[1])
	}))
}

func Test_Setting(t *testing.T) {
	Convey("setting", t, WithDao(func(d *Dao) {
		c := context.Background()
		mid := int64(910819)

		err := d.SetSetting(c, mid, model.Settings)
		So(err, ShouldBeNil)

		res, err := d.Setting(c, mid)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, model.Settings)
	}))
}

func Test_Report(t *testing.T) {
	r := &model.Report{
		APPID:        model.APPIDBBPhone,
		PlatformID:   model.PlatformIPhone,
		Mid:          910819,
		Buvid:        "b",
		DeviceToken:  strconv.FormatInt(time.Now().UnixNano(), 10),
		Build:        2233,
		TimeZone:     8,
		NotifySwitch: model.SwitchOn,
		DeviceBrand:  "OPPO",
		DeviceModel:  "OPPO R9st",
		OSVersion:    "6.0.1",
	}
	c := context.Background()
	Convey("report", t, WithDao(func(d *Dao) {

		id, err := d.AddReport(c, r)
		So(err, ShouldBeNil)
		r.ID = id

		rt, err := d.Report(c, r.DeviceToken)
		So(err, ShouldBeNil)
		So(rt, ShouldResemble, r)

		_, err = d.ReportsByMid(c, r.Mid)
		So(err, ShouldBeNil)

		_, err = d.ReportsByMids(c, []int64{r.Mid})
		So(err, ShouldBeNil)

		rows, err := d.DelReport(c, r.DeviceToken)
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)

		res, err := d.ReportsByID(context.TODO(), []int64{1, 2, 3})
		So(err, ShouldBeNil)
		// t.Logf("ReportsByID res(%v)", res)
		So(len(res), ShouldBeGreaterThan, 0)

		res1, err := d.Reports(context.TODO(), []string{"742381013eb5fb21e003479d041369481ca861d41a9e489abe9d44c27dd43d74", "cidViiN2cwpUdlrQXXPJlyk47N69WDje3PA1+ISCGIA="})
		So(err, ShouldBeNil)
		t.Log(len(res1))
	}))
}

func Test_UpdateReport(t *testing.T) {
	Convey("update report", t, WithDao(func(d *Dao) {
		ctx := context.Background()
		r := &model.Report{
			APPID:        model.APPIDBBPhone,
			PlatformID:   model.PlatformIPhone,
			Mid:          910819,
			Buvid:        "b",
			DeviceToken:  "dt",
			Build:        2233,
			TimeZone:     8,
			NotifySwitch: model.SwitchOn,
			DeviceBrand:  "OPPO",
			DeviceModel:  "OPPO R9st",
			OSVersion:    "6.0.1",
		}

		_, err := d.db.Exec(context.Background(), "delete from push_reports where token_hash=?", model.HashToken(r.DeviceToken))
		So(err, ShouldBeNil)

		id, err := d.AddReport(ctx, r)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
		r.ID = id

		rt, err := d.Report(ctx, r.DeviceToken)
		So(err, ShouldBeNil)
		So(rt, ShouldResemble, r)

		rt.APPID = 2
		rt.PlatformID = 3
		rt.NotifySwitch = model.SwitchOff
		rt.Mid = 123
		rt.Buvid = "buvidxxxx"
		rt.Build = 1000000
		rt.OSVersion = "x.x.x"

		err = d.UpdateReport(ctx, rt)
		So(err, ShouldBeNil)

		rt2, err := d.Report(ctx, r.DeviceToken)
		So(err, ShouldBeNil)
		So(rt2, ShouldResemble, rt)
		So(rt2, ShouldNotResemble, r)
	}))
}

func Test_Callback(t *testing.T) {
	Convey("add callback", t, WithDao(func(d *Dao) {
		cb := &model.Callback{
			Task:     "task123",
			APP:      model.APPIDBBPhone,
			Platform: model.PlatformXiaomi,
			Mid:      91221505,
			Pid:      model.MobiAndroid,
			Token:    "token",
			Buvid:    "buvid",
			Click:    1,
			Extra: &model.CallbackExtra{
				Status: 2,
			},
		}
		err := d.AddCallback(context.TODO(), cb)
		So(err, ShouldBeNil)
	}))
}

func Test_ReportByID(t *testing.T) {
	Convey("report by id", t, WithDao(func(d *Dao) {
		r, err := d.ReportByID(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		t.Logf("reportByID res(%+v)", r)
	}))
}
