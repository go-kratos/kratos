package http

import (
	"context"
	"net/http"
	"strconv"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/net/http/blademaster"
)

func mcnSignEntry(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNSignEntryReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNSignEntryReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.McnSignEntry(cont, arg.(*model.MCNSignEntryReq))
		},
		"mcnSignEntry")(c)
}

func mcnSignList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNSignStateReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnSignList(cont, arg.(*model.MCNSignStateReq))
		},
		"mcnSignList")(c)
}

func mcnSignOP(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNSignStateOpReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNSignStateOpReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.McnSignOP(cont, arg.(*model.MCNSignStateOpReq))
		},
		"mcnSignOP")(c)
}

func mcnUPReviewList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNUPStateReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnUPReviewList(cont, arg.(*model.MCNUPStateReq))
		},
		"mcnUPReviewList")(c)
}

func mcnUPOP(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNUPStateOpReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNUPStateOpReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.McnUPOP(cont, arg.(*model.MCNUPStateOpReq))
		},
		"mcnUPOP")(c)
}

func mcnPermitOP(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNSignPermissionReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNSignPermissionReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.McnPermitOP(cont, arg.(*model.MCNSignPermissionReq))
		},
		"McnPermitOP")(c)
}

func mcnUPPermitList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNUPPermitStateReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnUPPermitList(cont, arg.(*model.MCNUPPermitStateReq))
		},
		"McnUPPermitList")(c)
}

func mcnUPPermitOP(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNUPPermitOPReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNUPPermitOPReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.McnUPPermitOP(cont, arg.(*model.MCNUPPermitOPReq))
		},
		"McnUPPermitOP")(c)
}

func mcnList(c *blademaster.Context) {
	httpGetWriterByExport(
		new(model.MCNListReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNList(cont, arg.(*model.MCNListReq))
		},
		"mcnList")(c)
}

func mcnPayEdit(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNPayEditReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNPayEditReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNPayEdit(cont, arg.(*model.MCNPayEditReq))
		},
		"mcnPayEdit")(c)
}

// func mcnPayEdit(c *blademaster.Context) {
// 	httpPostFunCheckCookie(
// 		new(model.MCNPayEditReq),
// 		func(cont context.Context, arg interface{}) (res interface{}, err error) {
// 			var uids,name *http.Cookie
// 			args := arg.(*model.MCNPayEditReq)
// 			if name, err = c.Request.Cookie("username"); err == nil {
// 				args.UserName = name.Value
// 			}
// 			if uids, err = c.Request.Cookie("uid"); err == nil {
// 				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
// 					return
// 				}
// 			}
// 			return nil, srv.MCNPayEdit(cont, arg.(*model.MCNPayEditReq))
// 		},
// 		"mcnPayEdit")(c)
// }

func mcnPayStateEdit(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNPayStateEditReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNPayStateEditReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNPayStateEdit(cont, arg.(*model.MCNPayStateEditReq))
		},
		"mcnPayStateEdit")(c)
}

func mcnStateEdit(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNStateEditReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNStateEditReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNStateEdit(cont, arg.(*model.MCNStateEditReq))
		},
		"mcnStateEdit")(c)
}

func mcnRenewal(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNRenewalReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNRenewalReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNRenewal(cont, arg.(*model.MCNRenewalReq))
		},
		"mcnRenewal")(c)
}

func mcnInfo(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNInfoReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNInfo(cont, arg.(*model.MCNInfoReq))
		},
		"mcnInfo")(c)
}

func mcnUPList(c *blademaster.Context) {
	httpGetWriterByExport(
		new(model.MCNUPListReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNUPList(cont, arg.(*model.MCNUPListReq))
		},
		"mcnUPList")(c)
}

func mcnUPStatEdit(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNUPStateEditReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNUPStateEditReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNUPStateEdit(cont, arg.(*model.MCNUPStateEditReq))
		},
		"mcnUPStatEdit")(c)
}

func mcnCheatList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNCheatListReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNCheatList(cont, arg.(*model.MCNCheatListReq))
		},
		"mcnCheatList")(c)
}

func mcnCheatUPList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNCheatUPListReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNCheatUPList(cont, arg.(*model.MCNCheatUPListReq))
		},
		"mcnCheatUPList")(c)
}

func mcnImportUPInfo(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNImportUPInfoReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNImportUPInfo(cont, arg.(*model.MCNImportUPInfoReq))
		},
		"mcnImportUPInfo")(c)
}

func mcnImportUPRewardSign(c *blademaster.Context) {
	httpPostJSONCookie(
		new(model.MCNImportUPRewardSignReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			var uids, name *http.Cookie
			args := arg.(*model.MCNImportUPRewardSignReq)
			if name, err = c.Request.Cookie("username"); err == nil {
				args.UserName = name.Value
			}
			if uids, err = c.Request.Cookie("uid"); err == nil {
				if args.UID, err = strconv.ParseInt(uids.Value, 10, 64); err != nil {
					return
				}
			}
			return nil, srv.MCNImportUPRewardSign(cont, arg.(*model.MCNImportUPRewardSignReq))
		},
		"mcnImportUPRewardSign")(c)
}

func mcnIncreaseList(c *blademaster.Context) {
	httpGetFunCheckCookie(
		new(model.MCNIncreaseListReq),
		func(cont context.Context, arg interface{}) (res interface{}, err error) {
			return srv.MCNIncreaseList(cont, arg.(*model.MCNIncreaseListReq))
		},
		"mcnIncreaseList")(c)
}
