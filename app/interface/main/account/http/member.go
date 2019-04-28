package http

import (
	"io/ioutil"
	"net/http"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func identifyInfo(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		resData struct {
			Status model.IdentifyStatus `json:"identification"`
		}
		status int8
		err    error
	)
	if status, err = realnameSvc.Status(c, mid.(int64)); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	switch status {
	case 1:
		resData.Status = model.IdentifyOK
	case 0:
		resData.Status = model.IdentifyNotOK
	}
	c.JSON(resData, nil)
}

func submitOffical(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	params := new(model.OfficialApply)
	if err := c.Bind(params); err != nil {
		return
	}
	c.JSON(nil, memberSvc.SubmitOfficial(c, mid.(int64), params))
}

func officialDoc(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(memberSvc.OfficialDoc(c, mid.(int64)))
}

func officialConditions(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(memberSvc.OfficialConditions(c, mid.(int64)))
}

func uploadImage(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid := midI.(int64)
	log.Infov(c, log.KV("log", "account-interface: upload image"), log.KV("mid", mid))

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	if ftype != "image/jpeg" && ftype != "image/jpg" && ftype != "image/png" && ftype != "image/gif" {
		log.Error("account-interface: file type not allow file type(%s, mid: %v)", ftype, mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	url, err := memberSvc.UploadImage(c, ftype, bs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url":  url,
		"size": len(bs),
	}, nil)
}

func mobileVerify(c *bm.Context) {
	midI, _ := c.Get("mid")
	arg := &model.ArgMobileVerify{}
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.Country == 0 {
		arg.Country = 86
	}
	c.JSON(nil, memberSvc.MobileVerify(c, midI.(int64), arg.Mobile, arg.Country))
}

func officialPermission(ctx *bm.Context) {
	resp := &model.OfficialPermissionResponse{
		DeniedRoles: []int8{},
		Metadata:    map[string]interface{}{},
	}
	ctx.JSON(resp, nil)
}

func monthlyOfficialSubmittedTimes(ctx *bm.Context) {
	midI, _ := ctx.Get("mid")
	ctx.JSON(memberSvc.MonthlyOfficialSubmittedTimes(ctx, midI.(int64)), nil)
}

func officialAutoFillDoc(ctx *bm.Context) {
	midI, _ := ctx.Get("mid")
	ctx.JSON(memberSvc.OfficialAutoFillDoc(ctx, midI.(int64)))
}
