package v2

import (
	"context"
	"flag"
	"testing"
	"time"

	pb "go-common/app/service/live/resource/api/grpc/v2"
	"go-common/app/service/live/resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	urs *UserResourceService
)

func init() {
	flag.Set("conf", "../../cmd/user_resource.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	urs = NewUserResourceService(conf.Conf)
}

// go test  -test.v -test.run TestAddUserResource
func TestAddUserResource(t *testing.T) {
	Convey("TestAddUserResource", t, func() {
		res, err := urs.Add(context.TODO(), &pb.AddReq{
			ResType: 1,
			Title:   "resource",
			Url:     "http://www.bilibili.com",
			Weight:  0,
			Creator: "liutengda",
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAddUserResource
func TestEditUserResource(t *testing.T) {
	Convey("TestAddUserResource", t, func() {
		res, err := urs.Add(context.TODO(), &pb.AddReq{
			ResType: 1,
			Title:   "resource",
			Url:     "http://www.bilibili.com",
			Weight:  0,
			Creator: "liutengda",
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestEditUserResource", t, func() {
		res, err := urs.Edit(context.TODO(), &pb.EditReq{
			ResType:  1,
			CustomId: 1,
			Title:    time.Now().Format("2006-01-02 15:04:05"),
			Url:      "http://www.bilibili.com",
			Weight:   0,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestEditUserResource", t, func() {
		res, err := urs.Edit(context.TODO(), &pb.EditReq{
			ResType:  1,
			CustomId: 200000,
			Title:    time.Now().Format("2006-01-02 15:04:05"),
			Url:      "http://www.bilibili.com/" + time.Now().Format("2006-01-02 15:04:05"),
			Weight:   10,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeError)
	})

	Convey("TestEditUserResource", t, func() {
		res, err := urs.Edit(context.TODO(), &pb.EditReq{
			ResType:  2000,
			CustomId: 2,
			Title:    time.Now().Format("2006-01-02 15:04:05"),
			Url:      "http://www.bilibili.com/" + time.Now().Format("2006-01-02 15:04:05"),
			Weight:   10,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeError)
	})
}

// go test  -test.v -test.run TestQueryUserResource
func TestQueryUserResource(t *testing.T) {
	Convey("TestAddUserResource", t, func() {
		res, err := urs.Add(context.TODO(), &pb.AddReq{
			ResType: 1,
			Title:   "resource",
			Url:     "http://www.bilibili.com",
			Weight:  0,
			Creator: "liutengda",
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestQueryUserResource", t, func() {
		res, err := urs.Query(context.TODO(), &pb.QueryReq{
			ResType:  1,
			CustomId: 1,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestQueryUserResource", t, func() {
		res, err := urs.Query(context.TODO(), &pb.QueryReq{
			ResType:  1,
			CustomId: 100,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeError)
	})
}

// go test  -test.v -test.run TestListUserResource
func TestListUserResource(t *testing.T) {
	Convey("TestListUserResource", t, func() {
		res, err := urs.List(context.TODO(), &pb.ListReq{
			ResType:  1,
			Page:     1,
			PageSize: 10,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestListUserResource
func TestSetUserResourceStatus(t *testing.T) {
	Convey("TestAddUserResource", t, func() {
		res, err := urs.Add(context.TODO(), &pb.AddReq{
			ResType: 1,
			Title:   "resource",
			Url:     "http://www.bilibili.co",
			Weight:  0,
			Creator: "liutengda",
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestSetUserResourceStatus", t, func() {
		res, err := urs.SetStatus(context.TODO(), &pb.SetStatusReq{
			ResType:  1,
			CustomId: 1,
			Status:   99,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})

	Convey("TestSetUserResourceStatus", t, func() {
		res, err := urs.SetStatus(context.TODO(), &pb.SetStatusReq{
			ResType:  1,
			CustomId: 1,
			Status:   10,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
