package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddArc(t *testing.T) {
	Convey("should return without err", t, WithService(t, func(svf *Service) {
		err := svf.AddArc(context.TODO(), _mid, _dataAV, 0, _ip)
		So(err, ShouldBeNil)
	}))
}

func Test_DelArc(t *testing.T) {
	Convey("should return without err", t, WithService(t, func(svf *Service) {
		err := svf.DelArc(context.TODO(), _mid, _dataAV, _ip)
		So(err, ShouldBeNil)
	}))
}

func Test_UpArcs(t *testing.T) {
	Convey("should return archives", t, WithService(t, func(svf *Service) {
		res, err := svf.upArcs(context.TODO(), 10, _ip, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))

	Convey("should return blank arcs if not attention ups", t, WithService(t, func(svf *Service) {
		res, err := svf.upArcs(context.TODO(), 10, _ip)
		So(res, ShouldBeEmpty)
		So(err, ShouldBeNil)
	}))

}

func Test_attenUpArcs(t *testing.T) {
	Convey("should return archives", t, WithService(t, func(svf *Service) {
		res, err := svf.attenUpArcs(context.TODO(), 10, _mid, _ip)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))

	Convey("should return only time archives", t, WithService(t, func(svf *Service) {
		res, err := svf.attenUpArcs(context.TODO(), 10, _mid, _ip)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_upsPassed(t *testing.T) {
	Convey("should return archives", t, WithService(t, func(svf *Service) {
		res, err := svf.upsPassed(context.TODO(), []int64{_mid}, _ip)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_archives(t *testing.T) {
	Convey("should return archives", t, WithService(t, func(svf *Service) {
		res, err := svf.archives(context.TODO(), []int64{_dataAV}, _ip)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_archive(t *testing.T) {
	Convey("should return archive", t, WithService(t, func(svf *Service) {
		res, err := svf.archive(context.TODO(), _dataAV, _ip)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_ChangeAuthor(t *testing.T) {
	Convey("should return without err", t, WithService(t, func(svf *Service) {
		c := context.TODO()
		oldArc, _ := svf.archive(c, _dataAV, _ip)
		svf.AddArc(c, oldArc.Author.Mid, _dataAV, 0, _ip)

		err := svf.ChangeAuthor(c, oldArc.Aid, oldArc.Author.Mid, 2, _ip)
		So(err, ShouldBeNil)

		res, _ := svf.dao.UppersCaches(c, []int64{oldArc.Author.Mid}, 0, -1)
		So(res[oldArc.Author.Mid], ShouldBeEmpty)

		res, _ = svf.dao.UppersCaches(c, []int64{2}, 0, -1)
		So(len(res), ShouldNotBeEmpty)
	}))
}
