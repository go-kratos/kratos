package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
	"testing"
)

var (
	gs = model.GRPCSnap{
		ID:             1,
		GRPCID:         1,
		TaskName:       "test",
		Department:     "test",
		Project:        "ep",
		APP:            "melloi",
		Active:         1,
		HostName:       "172.22.33.22",
		Port:           9000,
		ServiceName:    "get",
		ProtoClassName: "code",
		PkgPath:        "book",
		RequestType:    "reqmsg",
		RequestMethod:  "reqmethod",
		RequestContent: "reqcontent",
		ResponseType:   "returnmsg",
		ScriptPath:     "sp",
		JarPath:        "jp",
		JmxPath:        "jmxpat",
		JmxLog:         "hnxkig",
		JtlLog:         "ht.lgo",
		ThreadsSum:     1,
		RampUp:         1,
		Loops:          -1,
		LoadTime:       1,
		UpdateBy:       "hujianping",
	}
)

func Test_GrpcSnap(t *testing.T) {
	Convey("query grpc snap", t, func() {
		_, err := s.QueryGRPCSnapByID(gs.ID)
		So(err, ShouldBeNil)
	})

	Convey("update grpc snap", t, func() {
		err := s.UpdateGRPCSnap(&gs)
		So(err, ShouldBeNil)
	})
	Convey("create grpc snap", t, func() {
		err := s.CreateGRPCSnap(&gs)
		So(err, ShouldBeNil)
	})

	Convey("delete grpc snap", t, func() {
		err := s.DeleteGRPCSnap(gs.ID)
		So(err, ShouldBeNil)
	})
}
