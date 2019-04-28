package service

import (
	"bytes"
	"context"
	"image/png"
	"os"
	"testing"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServicePrivileges
func TestServicePrivileges(t *testing.T) {
	Convey("TestServicePrivileges", t, func() {
		res, err := s.Privileges(context.TODO(), 0)
		for _, v := range res {
			t.Logf("%+v", v)
		}
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdatePrivilegeState
func TestUpdatePrivilegeState(t *testing.T) {
	Convey("TestUpdatePrivilegeState", t, func() {
		err := s.UpdatePrivilegeState(context.TODO(), &model.Privilege{
			ID:    4,
			State: 1,
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDeletePrivilege
func TestDeletePrivilege(t *testing.T) {
	Convey("TestDeletePrivilege", t, func() {
		err := s.DeletePrivilege(context.TODO(), 4)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdateOrder
func TestUpdateOrder(t *testing.T) {
	Convey("TestUpdateOrder", t, func() {
		err := s.UpdateOrder(context.TODO(), &model.ArgOrder{
			AID: 5,
			BID: 4,
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceUpload
func TestServiceUpload(t *testing.T) {
	Convey("TestServiceUpload", t, func() {
		file, _ := os.Open("/Users/baihai/Documents/test.png")
		So(file, ShouldNotBeNil)
		buf := new(bytes.Buffer)
		img, _ := png.Decode(file)
		png.Encode(buf, img)
		bys := buf.Bytes()
		t.Logf("conf %+v", s.c.Bfs)
		iconURL, iconGrayURL, webImageURL, appImageURL, err := s.uploadImage(context.TODO(), &model.ArgImage{
			IconBody:         bys,
			IconFileType:     "image/png",
			IconGrayBody:     bys,
			IconGrayFileType: "image/png",
		})
		t.Logf("url iconURL:%s iconGrayURL:%s webImageURL:%s appImageURL:%s", iconURL, iconGrayURL, webImageURL, appImageURL)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceAddPrivilege
func TestServiceAddPrivilege(t *testing.T) {
	Convey("TestServiceAddPrivilege", t, func() {
		file, _ := os.Open("/Users/baihai/Documents/test.png")
		So(file, ShouldNotBeNil)
		buf := new(bytes.Buffer)
		img, _ := png.Decode(file)
		png.Encode(buf, img)
		bys := buf.Bytes()
		t.Logf("conf %+v", s.c.Bfs)
		err := s.AddPrivilege(context.TODO(), &model.ArgAddPrivilege{
			Name:     "超级",
			Title:    "超级标题",
			Explain:  "超级内容",
			Type:     0,
			Operator: "admin",
			WebLink:  "http://www.baidu.com",
			AppLink:  "http://www.baidu.com",
		}, &model.ArgImage{
			IconBody:         bys,
			IconFileType:     "image/png",
			IconGrayBody:     bys,
			IconGrayFileType: "image/png",
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceUpdatePrivilege
func TestServiceUpdatePrivilege(t *testing.T) {
	Convey("TestServiceUpdatePrivilege", t, func() {
		file, _ := os.Open("/Users/baihai/Documents/test.png")
		So(file, ShouldNotBeNil)
		buf := new(bytes.Buffer)
		img, _ := png.Decode(file)
		png.Encode(buf, img)
		bys := buf.Bytes()
		t.Logf("conf %+v", s.c.Bfs)
		err := s.UpdatePrivilege(context.TODO(), &model.ArgUpdatePrivilege{
			ID:       3,
			Name:     "超级",
			Title:    "超级标题",
			Explain:  "超级内容",
			Type:     0,
			Operator: "admin",
			WebLink:  "http://www.baidu.com",
			AppLink:  "http://www.baidu.com",
		}, &model.ArgImage{
			IconBody:         bys,
			IconFileType:     "image/png",
			IconGrayBody:     bys,
			IconGrayFileType: "image/png",
		})
		So(err, ShouldBeNil)
	})
}
