package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/job/main/spy/model"
	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testMysqlMid int64 = 15555180
	size               = 3
)

func Test_Configs(t *testing.T) {
	Convey("Test_Configs info data", t, func() {
		res, err := d.Configs(context.TODO())
		fmt.Println(res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		So(res, ShouldContainKey, model.LimitBlockCount)
		So(res, ShouldContainKey, model.LessBlockScore)
		So(res, ShouldContainKey, model.AutoBlock)
	})
}

func Test_History(t *testing.T) {
	Convey("Test_History get history data", t, func() {
		res, err := d.History(context.TODO(), testMysqlMid)
		fmt.Println(res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_HistoryList(t *testing.T) {
	Convey("Test_HistoryList get history list data", t, func() {
		res, err := d.HistoryList(context.TODO(), testMysqlMid, size)
		fmt.Println(res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_TxUpdateUserState(t *testing.T) {
	Convey("Test_TxUpdateUserState no err", t, func() {
		var (
			c   = context.TODO()
			err error
			tx  *sql.Tx
			ui  = &model.UserInfo{Mid: 15555180}
		)
		tx, err = d.db.Begin(c)
		So(err, ShouldBeNil)
		ui.State = model.StateNormal
		err = d.TxUpdateUserState(context.TODO(), tx, ui)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_TxAddPunishment(t *testing.T) {
	Convey("Test_TxAddPunishment add data", t, func() {
		var (
			c   = context.TODO()
			err error
			tx  *sql.Tx
		)
		tx, err = d.db.Begin(c)
		So(err, ShouldBeNil)
		err = d.TxAddPunishment(context.TODO(), tx, testMysqlMid, 0, "test封禁", 100000000)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run Test_Stat
func Test_Stat(t *testing.T) {
	Convey("AddIncrStatistics", t, func() {
		var (
			c   = context.TODO()
			err error
		)
		stat := &model.Statistics{
			TargetMid: 1,
			TargetID:  1,
			EventID:   1,
			State:     0,
			Type:      1,
			Quantity:  1,
			Ctime:     time.Now(),
		}
		_, err = d.AddIncrStatistics(c, stat)
		So(err, ShouldBeNil)
	})
	Convey("AddStatistics", t, func() {
		var (
			c   = context.TODO()
			err error
		)
		stat := &model.Statistics{
			TargetMid: 1,
			TargetID:  1,
			EventID:   2,
			State:     0,
			Type:      1,
			Quantity:  1,
			Ctime:     time.Now(),
		}
		stat.Quantity = 3
		_, err = d.AddStatistics(c, stat)
		So(err, ShouldBeNil)
	})
}

func Test_TxAddEventHistory(t *testing.T) {
	Convey("test userinfo", t, func() {
		res, err := d.UserInfo(c, 7593623)
		So(err, ShouldBeNil)
		if res == nil {
			_, err = d.db.Exec(c, "insert into spy_user_info_23 (mid) values (7593623)")
			So(err, ShouldBeNil)
		}
	})
	var tx *sql.Tx
	Convey("Test_TxAddEventHistory start", t, func() {
		var err error
		tx, err = d.BeginTran(c)
		So(err, ShouldBeNil)
	})
	ueh := &model.UserEventHistory{Mid: 7593634, BaseScore: 100}
	Convey("Test_TxAddEventHistory", t, func() {

		err := d.TxAddEventHistory(c, tx, ueh)
		So(err, ShouldBeNil)
	})
	Convey("TxUpdateEventScore", t, func() {
		err := d.TxUpdateEventScore(c, tx, 7593623, 100, 100)
		So(err, ShouldBeNil)
	})
	Convey("Test_TxAddEventHistory commit", t, func() {
		err := tx.Commit()
		So(err, ShouldBeNil)
	})
	Convey("AllEvent", t, func() {
		res, err := d.AllEvent(c)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
	Convey("clean data", t, func() {
		d.db.Exec(c, "delete from spy_user_event_history_23 where mid = 7593623")

	})
}

func Test_SecurityLoginCount(t *testing.T) {
	Convey("Test_SecurityLoginCount", t, func() {
		_, err := d.SecurityLoginCount(c, 1, "test", time.Now(), time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_ReBuildMidList(t *testing.T) {
	Convey("Test_ReBuildMidList", t, func() {
		_, err := d.ReBuildMidList(c, 1, 1, time.Now(), time.Now(), 1)
		So(err, ShouldBeNil)
	})
}
