package weeklyhonor

import (
	"bytes"
	"context"
	"encoding/binary"
	"reflect"
	"testing"

	"go-common/library/database/hbase.v2"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestWeeklyhonorreverseString(t *testing.T) {
	convey.Convey("reverseString", t, func(ctx convey.C) {
		var (
			s = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := reverseString(s)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, p1)
			})
		})
	})
}

func TestWeeklyhonorhonorRowKey(t *testing.T) {
	convey.Convey("honorRowKey", t, func(ctx convey.C) {
		var (
			id   = int64(0)
			date = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := honorRowKey(id, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, p1)
			})
		})
	})
}

func TestWeeklyhonorHonorStat(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(0)
		date       = "20181112"
		val1 int32 = 500
		val2 int32 = -2
		cs         = getMockCells(val1, val2)
	)
	guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
		res := &hrpc.Result{
			Cells: cs,
		}
		return res, nil
	})
	defer guard.Unpatch()
	convey.Convey("honorStat", t, func() {
		stat, err := d.HonorStat(c, mid, date)
		convey.So(err, convey.ShouldBeNil)
		convey.So(stat, convey.ShouldNotBeNil)
		convey.So(stat.Play, convey.ShouldEqual, val1)
		convey.So(stat.PlayInc, convey.ShouldEqual, val1)
		convey.So(stat.Rank0, convey.ShouldEqual, val1)
		convey.So(stat.FansInc, convey.ShouldEqual, val2)
	})
}

func getMockCells(v1, v2 int32) []*hrpc.Cell {
	s1 := make([]byte, 0)
	buf1 := bytes.NewBuffer(s1)
	s2 := make([]byte, 0)
	buf2 := bytes.NewBuffer(s2)

	binary.Write(buf1, binary.BigEndian, v1)
	binary.Write(buf2, binary.BigEndian, v2)
	cells := []*hrpc.Cell{
		{
			Qualifier: []byte("play"),
			Value:     buf1.Bytes(),
			Family:    []byte("f"),
		},
		{
			Qualifier: []byte("play_last_w"),
			Value:     buf1.Bytes(),
			Family:    []byte("f"),
		},
		{
			Qualifier: []byte("play_inc"),
			Value:     buf1.Bytes(),
			Family:    []byte("f"),
		},
		{
			Qualifier: []byte("fans_inc"),
			Value:     buf2.Bytes(),
			Family:    []byte("f"),
		},
		{
			Qualifier: []byte("rank168"),
			Value:     buf1.Bytes(),
			Family:    []byte("r"),
		},
		{
			Qualifier: []byte("rank0"),
			Value:     buf1.Bytes(),
			Family:    []byte("r"),
		},
		{
			Qualifier: []byte("rank1"),
			Value:     buf1.Bytes(),
			Family:    []byte("r"),
		},
		{
			Qualifier: []byte("rank3"),
			Value:     buf1.Bytes(),
			Family:    []byte("r"),
		},
		{
			Qualifier: []byte("rank4"),
			Value:     buf1.Bytes(),
			Family:    []byte("r"),
		},
	}

	return cells
}
