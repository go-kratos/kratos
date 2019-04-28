package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/spy/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	statID    int64 = 3
	statState int8  = 1
	statCount int64 = 1
	statIsdel int8  = 1
	statMid   int64 = 1
	statType  int8  = 2
	pn              = 1
	ps              = 8
)

// go test  -test.v -test.run TestDB
func TestDB(t *testing.T) {
	Convey(" UpdateStatState ", t, WithMysql(func(d *Dao) {
		_, err := d.UpdateStatState(context.TODO(), statState, statID)
		So(err, ShouldBeNil)
	}))
	Convey(" UpdateStatQuantity ", t, WithMysql(func(d *Dao) {
		_, err := d.UpdateStatQuantity(context.TODO(), statCount, statID)
		So(err, ShouldBeNil)
	}))
	Convey(" DeleteStat ", t, WithMysql(func(d *Dao) {
		_, err := d.DeleteStat(context.TODO(), statIsdel, statID)
		So(err, ShouldBeNil)
	}))
	Convey(" Statistics ", t, WithMysql(func(d *Dao) {
		stat, err := d.Statistics(context.TODO(), statID)
		So(err, ShouldBeNil)
		fmt.Println("stat", stat)
	}))
	Convey(" LogList ", t, WithMysql(func(d *Dao) {
		logs, err := d.LogList(context.TODO(), statID, model.UpdateStat)
		So(err, ShouldBeNil)
		for _, l := range logs {
			fmt.Println(l.Context)
		}
	}))
	Convey(" StatListByMid ", t, WithMysql(func(d *Dao) {
		stat, err := d.StatListByMid(context.TODO(), statMid, pn, ps)
		So(err, ShouldBeNil)
		for _, s := range stat {
			fmt.Println(s.EventID)
		}
	}))
	Convey(" StatListByID ", t, WithMysql(func(d *Dao) {
		stat, err := d.StatListByID(context.TODO(), statID, statType, pn, ps)
		So(err, ShouldBeNil)
		for _, s := range stat {
			fmt.Println(s.EventID)
		}
	}))
	Convey(" StatCountByMid ", t, WithMysql(func(d *Dao) {
		stat, err := d.StatCountByMid(context.TODO(), statMid)
		So(err, ShouldBeNil)
		fmt.Println("count ", stat)
	}))
	Convey(" StatCountByID ", t, WithMysql(func(d *Dao) {
		stat, err := d.StatCountByID(context.TODO(), statID, statType)
		So(err, ShouldBeNil)
		fmt.Println("count ", stat)
	}))
}
