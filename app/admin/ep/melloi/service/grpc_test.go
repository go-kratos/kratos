package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
	"testing"
)

var (
	protoPath  = "/data/jmeter-log/test/ep/melloi/test/381016516/"
	protoFile  = "StreamEvent.proto"
	protoModel = model.ProtoPathModel{
		RootPath: "/data/jmeter-log/test/ep/melloi/test/445057856/proto", ExtraPath: "account/service/member",
	}
	scriptID = model.GRPCExecuteScriptRequest{ScriptID: 1}

	grpc = model.GRPC{
		TaskName:       "testGrpc",
		Department:     "test",
		Project:        "ep",
		APP:            "melloi",
		Active:         1,
		HostName:       "172.22.33.22",
		Port:           9000,
		ServiceName:    "Identify",
		ProtoClassName: "Api",
		PkgPath:        "V1",
		RequestType:    "GetCookieInfo",
		RequestMethod:  "getCookieInfo",
		RequestContent: "{\"Cookie\":\"sid:1ers12;SETDATA:a18jds9234js9sfa24jsdf\"}",
		ResponseType:   "Reponse",
		ScriptPath:     "/data/jmeter/log/test/ep/melloi/",
		JarPath:        "/data/jmeter/log/test/ep/melloi/text.jar",
		ThreadsSum:     1,
		RampUp:         1,
		Loops:          -1,
		LoadTime:       100,
		UpdateBy:       "hujianping",
		IsDebug:        0,
	}
	qgr = model.QueryGRPCRequest{
		Executor: "hujianping",
		GRPC:     grpc,
	}
	gasr = model.GRPCAddScriptRequest{
		TaskName:       "testGrpc",
		Department:     "test",
		Project:        "ep",
		APP:            "melloi",
		Active:         1,
		HostName:       "172.22.33.22",
		Port:           9000,
		ServiceName:    "Identify",
		ProtoClassName: "Api",
		PkgPath:        "V1",
		RequestType:    "GetCookieInfo",
		RequestMethod:  "getCookieInfo",
		RequestContent: "{\"Cookie\":\"sid:1ers12;SETDATA:a18jds9234js9sfa24jsdf\"}",
		ResponseType:   "Reponse",
		ScriptPath:     "/data/jmeter/log/test/ep/melloi/",
		JarPath:        "/data/jmeter/log/test/ep/melloi/text.jar",
		ThreadsSum:     1,
		RampUp:         1,
		Loops:          -1,
		LoadTime:       100,
		UpdateBy:       "hujianping",
		IsDebug:        0,
	}
)

func Test_Grpc(t *testing.T) {
	Convey("proto parse", t, func() {
		_, err := s.ProtoParsing(protoPath, protoFile)
		So(err, ShouldBeNil)
	})
	Convey("create proto dependency dir", t, func() {
		err := s.CreateProtoImportDir(&protoModel)
		So(err, ShouldBeNil)
	})
	Convey(" add grpc script", t, func() {
		_, err := s.GRPCAddScript(c, &gasr)
		So(err, ShouldBeNil)
	})
	Convey("create jmx file", t, func() {
		_, err := s.CreateJmx(c, &gasr)
		So(err, ShouldBeNil)
	})
	Convey("run by script", t, func() {
		cookie = "baf4dd3244116f492b71af3532cac03e"
		_, err := s.GRPCRunByScriptID(c, &scriptID, userName, cookie)
		So(err, ShouldBeNil)
	})
	Convey("query grpc", t, func() {
		_, err := s.QueryGrpc(c, "e2df43ed324d20811e8d1be1a9fb36d5", &qgr)
		So(err, ShouldBeNil)
	})
	Convey("run grpc by model", t, func() {
		cookie = "baf4dd3244116f492b71af3532cac03e"
		_, err := s.GRPCRunByModel(c, &grpc, userName, cookie)
		So(err, ShouldBeNil)
	})
	Convey("query grpc by id", t, func() {
		_, err := s.QueryGrpcById(grpc.ID)
		So(err, ShouldBeNil)
	})
	Convey("update grpc", t, func() {
		err := s.UpdateGrpc(&grpc)
		So(err, ShouldBeNil)
	})
	Convey("delete grpc", t, func() {
		err := s.DeleteGrpc(grpc.ID)
		So(err, ShouldBeNil)
	})
	Convey("create jmx file ", t, func() {
		_, err := s.createJmeterFile(&gasr)
		So(err, ShouldBeNil)
	})

}
