package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkDMXML(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := svr.DMXML(context.TODO(), 1, 1221)
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

func TestDMDistribution(t *testing.T) {
	res, err := svr.DMDistribution(context.TODO(), 1, 1221, 5)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestAsyncAddDM(t *testing.T) {
	dm := model.DM{
		ID: 123,
		Content: &model.Content{
			ID: 123,
		},
	}
	if err := svr.asyncAddDM(&dm, 1, 1); err != nil {
		t.Error(err)
	}
}

func TestDMList(t *testing.T) {
	var (
		oid   int64 = 5
		dmids       = []int64{719149905, 719149906, 719149907}
	)
	Convey("test dm list", t, func() {
		res, err := svr.dmList(context.TODO(), 1, oid, dmids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
