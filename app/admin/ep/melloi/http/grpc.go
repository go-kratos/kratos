package http

import (
	"net/http"
	"path/filepath"

	"go-common/app/admin/ep/melloi/model"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func getProto(c *bm.Context) {
	v := new(struct {
		Path string `form:"path"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	path, fileName := filepath.Split(v.Path)
	var (
		res = make(map[string]interface{})
		err error
	)
	if res, err = srv.ProtoParsing(path, fileName); err != nil {
		log.Error("parser grpc error(%v)", err)
		c.JSON(nil, err)
		return
	}

	if err = srv.CreateGRPCImportDir(res, path); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func createDependencyPath(c *bm.Context) {
	protoPath := &model.ProtoPathModel{}

	if err := c.BindWith(protoPath, binding.JSON); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.CreateProtoImportDir(protoPath))
}

func grpcQuickStart(c *bm.Context) {
	var (
		userName *http.Cookie
		qsReq    = &model.GRPCQuickStartRequest{}
		err      error
		cookie   string
	)

	cookie = c.Request.Header.Get("Cookie")
	if err = c.BindWith(qsReq, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	if userName, err = c.Request.Cookie("username"); err != nil {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}
	c.JSON(srv.GRPCQuickStart(c, qsReq, userName.Value, cookie))
}

func saveGrpc(c *bm.Context) {
	var (
		qsReq = &model.GRPCQuickStartRequest{}
		err   error
	)

	if err = c.BindWith(qsReq, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(nil, srv.SaveGRPCQuickStart(c, qsReq))
}

func runGrpc(c *bm.Context) {
	var (
		grpc     = model.GRPC{}
		err      error
		userName *http.Cookie
		cookie   string
	)

	cookie = c.Request.Header.Get("Cookie")
	if err = c.BindWith(&grpc, binding.JSON); nil != err {
		c.JSON(nil, err)
		return
	}
	if userName, err = c.Request.Cookie("username"); err != nil {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}
	c.JSON(srv.GRPCRunByModel(c, &grpc, userName.Value, cookie))
}

func runGrpcByScriptId(c *bm.Context) {
	var (
		userName   *http.Cookie
		err        error
		grpcExeReq = &model.GRPCExecuteScriptRequest{}
		cookie     string
	)

	cookie = c.Request.Header.Get("Cookie")
	if err = c.BindWith(grpcExeReq, binding.JSON); err != nil {
		return
	}
	if userName, err = c.Request.Cookie("username"); err != nil {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}
	c.JSON(srv.GRPCRunByScriptID(c, grpcExeReq, userName.Value, cookie))
}

func grpcAddScript(c *bm.Context) {
	grpcReq := model.GRPCAddScriptRequest{}
	if err := c.BindWith(&grpcReq, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.GRPCAddScript(c, &grpcReq))
}

func queryGrpc(c *bm.Context) {
	qgr := model.QueryGRPCRequest{}

	if err := c.BindWith(&qgr, binding.Form); err != nil {
		c.JSON(nil, err)
	}
	if err := qgr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	sessionID, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryGrpc(c, sessionID.Value, &qgr))
}

func deleteGrpc(c *bm.Context) {
	v := new(struct {
		ID int `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.DeleteGrpc(v.ID))
}

func updateGrpc(c *bm.Context) {
	grpc := model.GRPC{}
	if err := c.BindWith(&grpc, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}

	qgr := &model.GRPCAddScriptRequest{}
	qgr.ThreadsSum = grpc.ThreadsSum
	qgr.RampUp = grpc.RampUp
	qgr.Loops = grpc.Loops
	qgr.LoadTime = grpc.LoadTime
	qgr.HostName = grpc.HostName
	qgr.Port = grpc.Port
	qgr.ServiceName = grpc.ServiceName
	qgr.ProtoClassName = grpc.ProtoClassName
	qgr.PkgPath = grpc.PkgPath
	qgr.RequestType = grpc.RequestType
	qgr.ResponseType = grpc.ResponseType
	qgr.ScriptPath = grpc.ScriptPath
	qgr.RequestMethod = grpc.RequestMethod
	qgr.RequestContent = grpc.RequestContent
	qgr.TaskName = grpc.TaskName
	qgr.ParamFilePath = grpc.ParamFilePath
	qgr.ParamNames = grpc.ParamNames
	qgr.ParamDelimiter = grpc.ParamDelimiter
	g, err := srv.CreateJmx(c, qgr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	grpc.JmxPath = g.JmxPath
	grpc.JtlLog = g.JtlLog
	grpc.JmxLog = g.JmxLog

	c.JSON(nil, srv.UpdateGrpc(&grpc))
}

func queryGrpcSnap(c *bm.Context) {
	v := new(struct {
		ID int `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.QueryGRPCSnapByID(v.ID))
}
