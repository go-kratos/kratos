package dao

import (
	"fmt"
	"go-common/app/admin/main/filter/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_AiWhite(t *testing.T) {
	Convey("TestDao_AiScore", t, func() {
		res, err := dao.AiWhite(ctx, 1, 1)
		for _, re := range res {
			fmt.Printf("re:%+v \n", re)
		}
		fmt.Printf("err:%+v \n", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_AiWhiteCount(t *testing.T) {
	Convey("TestDao_AiWhiteCount", t, func() {
		res, err := dao.AiWhiteCount(ctx)
		fmt.Printf("res:%+v \n", res)
		fmt.Printf("err:%+v \n", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
func TestDao_InsertAiWhite(t *testing.T) {
	Convey("TestDao_InsertAiWhite", t, func() {
		res, err := dao.InsertAiWhite(ctx, 3)
		fmt.Printf("err:%+v \n", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_EditAiWhite(t *testing.T) {
	Convey("TestDao_EditAiWhite", t, func() {
		res, err := dao.EditAiWhite(ctx, 2, 1)
		fmt.Printf("res:%+v \n", res)
		fmt.Printf("err:%+v \n", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_InsertAiCase(t *testing.T) {
	Convey("TestDao_InsertAiCase", t, func() {
		var aicase = &model.AiCase{}
		aicase.Content = "test"
		aicase.Source = 1
		aicase.Type = 1
		res, err := dao.InsertAiCase(ctx, aicase)
		fmt.Printf("res:%+v \n", res)
		fmt.Printf("err:%+v \n", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
