package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/usersuit/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestService_MedalAllInfo
func TestService_MedalAllInfo(t *testing.T) {
	var (
		res = &model.MedalAllInfos{}
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.MedalAllInfo(context.Background(), 1)
		for _, re := range res.List {
			fmt.Printf("Name:%+v \n", re.Name)
			fmt.Printf("Count:%+v \n", re.Count)
			fmt.Printf("re:%+v \n", re.Data)
		}
		fmt.Printf("err:%+v \n", err)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalActivatedMulti(t *testing.T) {
	var (
		res = make(map[int64]*model.MedalInfo)
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.MedalActivatedMulti(context.Background(), []int64{650460, 4780461, 14082, 631053})
		for k, re := range res {
			fmt.Printf("key:%v re:%+v \n", k, re)
		}
		fmt.Printf("err:%+v \n", err)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
