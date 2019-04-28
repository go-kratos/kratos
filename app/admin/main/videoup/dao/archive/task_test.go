package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/model/archive"
	"testing"
	"time"
)

func Test_Weight(t *testing.T) {
	cfg, boolean, err := archive.ParseWeightConf(&archive.WeightConf{
		Radio:  archive.WConfTaskID,
		Ids:    "1,2,3,4,5",
		Rule:   0,
		Weight: 15,
		Desc:   "测试taskid权重配置",
	}, 10086, "cxf")
	if err != nil || cfg == nil || !boolean {
		t.Fatalf("err %+v cfg:%+v bool:%v\n", err, cfg, boolean)
	}

	if err = d.InWeightConf(context.TODO(), cfg); err != nil {
		t.Fatal(err)
	}
}

func Test_MulAddTaskHis(t *testing.T) {
	row, err := d.MulAddTaskHis(context.TODO(), []*archive.TaskForLog{
		&archive.TaskForLog{
			ID:      1,
			Cid:     2,
			Subject: 0,
			Mtime:   time.Now(),
		}, &archive.TaskForLog{
			ID:      2,
			Cid:     4,
			Subject: 1,
			Mtime:   time.Now(),
		},
	}, archive.ActionDispatch, 10086)
	if row != 2 || err != nil {
		t.Fail()
	}
}

func Test_TaskTooksByHalfHour(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.TaskTooksByHalfHour(context.Background(), time.Now().Add(-time.Hour), time.Now())
		So(err, ShouldBeNil)
	}))
}
