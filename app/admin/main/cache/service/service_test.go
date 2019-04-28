package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/cache/conf"
	"go-common/app/admin/main/cache/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	os.Exit(m.Run())
}
func TestCluster(t *testing.T) {
	Convey("test cluster ", t, func() {
		req := &model.ClusterReq{
			PN: 1,
			PS: 10,
		}
		resp, err := svr.Clusters(context.TODO(), req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		t.Logf("resp %v", resp.Clusters[0])
	})
}

func TestAddCluster(t *testing.T) {
	Convey("test add cluster ", t, func() {
		req := &model.AddClusterReq{
			Type:             "memcache",
			AppID:            "test",
			HashMethod:       "sha1",
			HashDistribution: "ketama",
		}
		_, err := svr.AddCluster(context.TODO(), req)
		So(err, ShouldBeNil)
	})
}

func TestSearchCluster(t *testing.T) {
	Convey("test  search cluster", t, func() {
		req := &model.ClusterReq{
			AppID: "test",
		}
		resp, err := svr.Cluster(context.TODO(), req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		t.Logf("search resp %+v", resp)
	})
}

func TestModifyCluster(t *testing.T) {
	Convey("test add cluster nodes", t, func() {
		req := &model.ModifyClusterReq{
			ID:     1,
			Action: 1,
			Nodes:  `[{"addr":"11","alias":"test1","weight":1}]`,
		}
		_, err := svr.ModifyCluster(context.TODO(), req)
		So(err, ShouldBeNil)
	})
	Convey("test get cluster detail", t, func() {
		req := &model.ClusterDtlReq{
			ID: 1,
		}
		resp, err := svr.ClusterDtl(context.TODO(), req)
		So(err, ShouldBeNil)
		So(len(resp.Nodes), ShouldEqual, 2)
	})

}
