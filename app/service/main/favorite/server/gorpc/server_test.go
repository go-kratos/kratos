package rpc

import (
	"context"
	favmdl "go-common/app/service/main/favorite/model"
	favsrv "go-common/app/service/main/favorite/api/gorpc"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.TODO()
)

func WithRPC(f func(client *favsrv.Service)) func() {
	return func() {
		client := favsrv.New2(nil)
		time.Sleep(2 * time.Second)
		f(client)
	}
}

// func Test_Folders(t *testing.T) {
// 	Convey("Folders", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgFolders{
// 			Type: 5,
// 			Mid:  88888894,
// 			FVmids: []*favmdl.ArgFVmid{
// 				{
// 					Fid:  12331,
// 					Vmid: 27515232,
// 				},
// 				{
// 					Fid:  12315,
// 					Vmid: 27515254,
// 				},
// 				{
// 					Fid:  1237,
// 					Vmid: 27515255,
// 				},
// 			},
// 		}
// 		res, err := client.Folders(ctx, arg)
// 		t.Logf("res: %+v, error:%+v", res, err)
// 		So(err, ShouldBeNil)
// 	}))
// }

// func Test_AddFolder(t *testing.T) {
// 	Convey("AddFolder", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgAddFolder{
// 			Type:        5,
// 			Mid:         88888894,
// 			Name:        "1234",
// 			Description: "1234",
// 			Cover:       "1234",
// 			Cookie:      "finger=14bc3c4e; sid=m65qzbpt; LIVE_BUVID=AUTO5915076321203625; fts=1507632123; UM_distinctid=15f05e08d7d1e5-0319506e9267a5-31657c00-1fa400-15f05e08d7e2b3; pgv_pvi=294734848; buvid3=C0E6B232-BC5C-4AEC-879C-1CC61C2336841986infoc; rpdid=kmilkmximpdoswqosospw; tma=136533283.50523868.1508466136624.1508466136624.1508466136624.1; tmd=2.136533283.50523868.1508466136624.; pgv_si=s1585227776; DedeUserID=88888894; DedeUserID__ckMd5=53cf4c9cd3a9e1fb; SESSDATA=6c75f0f3%2C1511521084%2Ce6967cf6; bili_jct=63ab513073169112f9f89b97c1713009; _cnt_pm=0; _cnt_notify=0; _dfcaptcha=db4362e021547cfab3dce5b4718ad573",
// 		}
// 		fid, err := client.AddFolder(ctx, arg)
// 		So(err, ShouldBeNil)
// 		t.Logf("fid: %d", fid)
// 	}))
// }

// func Test_UpdateFolder(t *testing.T) {
// 	Convey("UpdateFolder", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgUpdateFolder{
// 			Type:        5,
// 			Fid:         9,
// 			Mid:         88888894,
// 			Name:        "123",
// 			Description: "",
// 			// Cover:       "123",
// 		}
// 		err := client.UpdateFolder(ctx, arg)
// 		So(err, ShouldBeNil)
// 	}))
// }

// func Test_Add(t *testing.T) {
// 	Convey("Add", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgAdd{
// 			Type: 1,
// 			Mid:  88888894,
// 			Oid:  123,
// 			Fid:  456,
// 		}
// 		err := client.Add(ctx, arg)
// 		So(err, ShouldBeNil)
// 	}))
// }

// func Test_AddVideo(t *testing.T) {
// 	Convey("AddVideo", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgAddVideo{
// 			Mid:    88888894,
// 			Aid:    123,
// 			Fids:   []int64{123, 456},
// 			Cookie: "finger=14bc3c4e; sid=m65qzbpt; LIVE_BUVID=AUTO5915076321203625; fts=1507632123; UM_distinctid=15f05e08d7d1e5-0319506e9267a5-31657c00-1fa400-15f05e08d7e2b3; pgv_pvi=294734848; buvid3=C0E6B232-BC5C-4AEC-879C-1CC61C2336841986infoc; rpdid=kmilkmximpdoswqosospw; tma=136533283.50523868.1508466136624.1508466136624.1508466136624.1; tmd=2.136533283.50523868.1508466136624.; pgv_si=s1585227776; DedeUserID=88888894; DedeUserID__ckMd5=53cf4c9cd3a9e1fb; SESSDATA=6c75f0f3%2C1511521084%2Ce6967cf6; bili_jct=63ab513073169112f9f89b97c1713009; _cnt_pm=0; _cnt_notify=0; _dfcaptcha=db4362e021547cfab3dce5b4718ad573",
// 		}
// 		err := client.AddVideo(ctx, arg)
// 		So(err, ShouldBeNil)
// 	}))
// }

// func Test_Folder(t *testing.T) {
// 	Convey("Add", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgFolder{
// 			Type: 5,
// 			Mid:  0,
// 			Fid:  0,
// 		}
// 		res, err := client.Folder(ctx, arg)
// 		t.Logf("res: %+v", res)
// 		So(err, ShouldBeNil)
// 	}))
// }

// func Test_IsFavedByFid(t *testing.T) {
// 	Convey("IsFavedByFid", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgIsFavedByFid{
// 			Type: 6,
// 			Mid:  88888894,
// 			Fid:  23,
// 			Oid:  11112,
// 		}
// 		res, err := client.IsFavedByFid(ctx, arg)
// 		t.Logf("res: %+v", res)
// 		So(err, ShouldBeNil)
// 	}))
// }

func Test_CntUserFolders(t *testing.T) {
	Convey("CntUserFolders", t, WithRPC(func(client *favsrv.Service) {
		arg := &favmdl.ArgCntUserFolders{
			Type: 6,
			Mid:  88888894,
			Vmid: 88888894,
		}
		res, err := client.CntUserFolders(ctx, arg)
		t.Logf("res: %+v", res)
		So(err, ShouldBeNil)
	}))
}

// func Test_DelFolder(t *testing.T) {
// 	Convey("DelFolder", t, WithRPC(func(client *favsrv.Service) {
// 		arg := &favmdl.ArgDelFolder{
// 			Type: 5,
// 			Mid:  88888894,
// 			Fid:  9,
// 		}
// 		err := client.DelFolder(ctx, arg)
// 		So(err, ShouldBeNil)
// 	}))
// }

func Test_IsFavs(t *testing.T) {
	Convey("IsFavs", t, WithRPC(func(client *favsrv.Service) {
		arg := &favmdl.ArgIsFavs{
			Type: 1,
			Mid:  88888894,
			Oids: []int64{123, 456},
		}
		res, err := client.IsFavs(ctx, arg)
		t.Logf("res: %+v", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}
