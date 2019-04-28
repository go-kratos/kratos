package dao

import (
	"testing"

	"go-common/app/admin/ep/merlin/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	imageAdd = model.Image{
		Name:        "docker-reg.bilibili.co/zccdebian:1.0",
		Status:      1,
		OS:          "debian",
		Version:     "1.0 64位",
		Description: "base",
		CreatedBy:   "ut",
	}
	imageUpdate = model.Image{
		Name:        "docker-reg.bilibili.co/zccdebian:1.0update",
		Status:      2,
		OS:          "debian update",
		Version:     "1.0 64位 update",
		Description: "base update",
		CreatedBy:   "ut update",
		UpdatedBy:   "ut",
	}
)

func Test_Image(t *testing.T) {
	Convey("add image", t, func() {
		err := d.AddImage(&imageAdd)
		So(err, ShouldBeNil)
	})
	Convey("get image", t, func() {
		images, err := d.Images()
		So(err, ShouldBeNil)
		So(len(images), ShouldBeGreaterThan, 0)
	})
	Convey("update image", t, func() {
		imageUpdate.ID = imageAdd.ID
		err := d.UpdateImage(&imageUpdate)
		So(err, ShouldBeNil)
	})
	Convey("delete image", t, func() {
		imageDel := model.Image{ID: imageAdd.ID}
		err := d.DelImage(imageDel.ID)
		So(err, ShouldBeNil)
	})
}
