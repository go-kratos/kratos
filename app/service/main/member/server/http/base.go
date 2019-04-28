package http

import (
	"encoding/json"
	"strconv"
	"strings"

	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/time"
)

func base(ctx *bm.Context) {
	var (
		err error
		mid int64
		// baseInfo *model.BaseInfo
		params = ctx.Request.Form
		midStr = params.Get("mid")
		// res      = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if baseInfo, err = memberSvc.BaseInfo(c, mid); err != nil {
	// 	log.Error("relationSvc.BaseInfo(%d) error(%v)", mid, err)
	// 	res["code"] = err
	// 	return
	// }
	// res["data"] = baseInfo
	ctx.JSON(memberSvc.BaseInfo(ctx, mid))
}

func member(ctx *bm.Context) {
	params := ctx.Request.Form
	// res := c.Result()
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid <= 0 {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	// mb, err := memberSvc.Member(c, mid)
	// if err != nil {
	// 	log.Error("Failed to memberSvc.Member(%d): %+v", mid, err)
	// 	res["code"] = err
	// 	return
	// }
	// res["data"] = mb
	ctx.JSON(memberSvc.Member(ctx, mid))
}

func batchBase(ctx *bm.Context) {
	var (
		err  error
		mid  int64
		mids []int64
		// binfo   map[int64]*model.BaseInfo
		params  = ctx.Request.Form
		midsStr = params.Get("mids")
		// res     = c.Result()
	)
	for _, str := range strings.Split(midsStr, ",") {
		if mid, err = strconv.ParseInt(str, 10, 64); err != nil {
			// res["code"] = ecode.RequestErr
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		mids = append(mids, mid)
	}
	// if binfo, err = memberSvc.BatchBaseInfo(c, mids); err != nil {
	// 	log.Error("memberSvc.BaseInfo(%d) error(%v)", mid, err)
	// 	res["code"] = err
	// 	return
	// }
	// res["data"] = binfo
	ctx.JSON(memberSvc.BatchBaseInfo(ctx, mids))
}

func setSign(ctx *bm.Context) {
	var (
		err      error
		mid      int64
		params   = ctx.Request.Form
		midStr   = params.Get("mid")
		usersign = params.Get("user_sign")
		// res      = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if usersign == "" {
	// 	res["code"] = ecode.RequestErr
	// 	return
	// }
	// 获取用户状态逻辑 status判断
	// if err := memberSvc.SetSign(c, mid, usersign); err != nil {
	// 	log.Error("memberSvc.SetSign(%d) error(%v)", mid, err)
	// 	res["code"] = ecode.ServerErr
	// 	return
	// }
	ctx.JSON(nil, memberSvc.SetSign(ctx, mid, usersign))
}

func setName(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
		name   = params.Get("name")
		// res    = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if name == "" {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if err := memberSvc.SetName(c, mid, name); err != nil {
	// 	log.Error("memberSvc.SetUname(%d) error(%v)", mid, err)
	// 	res["code"] = ecode.ServerErr
	// 	return
	// }
	ctx.JSON(nil, memberSvc.SetName(ctx, mid, name))
}

func setRank(ctx *bm.Context) {
	var (
		err     error
		mid     int64
		rank    int64
		params  = ctx.Request.Form
		midStr  = params.Get("mid")
		rankStr = params.Get("rank")
		// res     = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if rank, err = strconv.ParseInt(rankStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if err := memberSvc.SetRank(c, mid, rank); err != nil {
	// 	log.Error("relationSvc.SetRank(%d) error(%v)", mid, err)
	// 	res["code"] = ecode.ServerErr
	// 	return
	// }
	ctx.JSON(nil, memberSvc.SetRank(ctx, mid, rank))
}

// setSex set sex.
func setSex(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		sex    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
		sexStr = params.Get("sex")
		// res    = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if sex, err = strconv.ParseInt(sexStr, 10, 8); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if err = memberSvc.SetSex(c, mid, sex); err != nil {
	// 	log.Error("memberSvc.SetSex(%d, %d) error(%v)", mid, sex, err)
	// 	res["code"] = ecode.ServerErr
	// }
	ctx.JSON(nil, memberSvc.SetSex(ctx, mid, sex))
}

// setBirthday set Birthday.
func setBirthday(ctx *bm.Context) {
	var (
		err         error
		mid         int64
		birthdayTs  int64
		birthday    time.Time
		params      = ctx.Request.Form
		midStr      = params.Get("mid")
		birthdayStr = params.Get("birthday")
		// res         = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if birthdayTs, err = strconv.ParseInt(birthdayStr, 10, 32); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	birthday = time.Time(birthdayTs)
	// if err = memberSvc.SetBirthday(c, mid, birthday); err != nil {
	// 	log.Error("memberSvc.SetBirthday(%d, %d) error(%v)", mid, birthday, err)
	// 	res["code"] = ecode.ServerErr
	// }
	ctx.JSON(nil, memberSvc.SetBirthday(ctx, mid, birthday))
}

// setFace set face.
func setFace(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
		face   = params.Get("face")
		// res    = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if err = memberSvc.SetFace(c, mid, face); err != nil {
	// 	log.Error("memberSvc.SetFace(%d, %d) error(%v)", mid, face, err)
	// 	res["code"] = ecode.ServerErr
	// }
	ctx.JSON(nil, memberSvc.SetFace(ctx, mid, face))
}

func setBase(ctx *bm.Context) {
	var (
		err         error
		mid         int64
		rank        int64
		sex         int64
		birthday    int64
		params      = ctx.Request.Form
		midStr      = params.Get("mid")
		rankStr     = params.Get("rank")
		face        = params.Get("face")
		birthdayStr = params.Get("birthday")
		name        = params.Get("name")
		sign        = params.Get("user_sign")
		sexStr      = params.Get("sex")
		// res         = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if len(rankStr) != 0 {
		if rank, err = strconv.ParseInt(rankStr, 10, 64); err != nil {
			// res["code"] = ecode.RequestErr
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if len(sexStr) != 0 {
		if sex, err = strconv.ParseInt(sexStr, 10, 64); err != nil {
			// res["code"] = ecode.RequestErr
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if len(birthdayStr) != 0 {
		if birthday, err = strconv.ParseInt(birthdayStr, 10, 64); err != nil {
			// res["code"] = ecode.RequestErr
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	b := &model.BaseInfo{Mid: mid, Face: face, Sex: sex, Birthday: time.Time(birthday), Name: name, Sign: sign, Rank: rank}
	// if err := memberSvc.SetBase(c, b); err != nil {
	// 	log.Error("memberSvc.SetBase(%d) error(%v)", mid, err)
	// 	res["code"] = ecode.ServerErr
	// 	return
	// }
	ctx.JSON(nil, memberSvc.SetBase(ctx, b))
}

func updateMorals(ctx *bm.Context) {
	var (
		err error
		// morals map[int64]int64
	)
	// res := c.Result()
	arg := &model.ArgUpdateMorals{}
	if err = ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	// morals, err = memberSvc.UpdateMorals(c, arg)
	// if err != nil {
	// 	res["code"] = err
	// 	return
	// }
	// res["data"] = morals
	ctx.JSON(memberSvc.UpdateMorals(ctx, arg))
}

func updateMoral(ctx *bm.Context) {
	arg := &model.ArgUpdateMoral{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, memberSvc.UpdateMoral(ctx, arg))
}

func undoMoral(ctx *bm.Context) {
	arg := &model.ArgUndo{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, memberSvc.UndoMoral(ctx, arg.LogID, arg.Remark, arg.Operator))
}

// cacheDel delete user cache.
func cacheDel(ctx *bm.Context) {
	var (
		mid    int64
		action string
		ak     string
		sd     string
		err    error
	)
	// res := c.Result()
	query := ctx.Request.Form
	midStr := query.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	action = query.Get("modifiedAttr")
	ak = query.Get("access_token")
	sd = query.Get("session")
	memberSvc.DelCache(ctx, mid, action, ak, sd)
	// res["code"] = ecode.OK
	ctx.JSON(nil, nil)
}

// addPropertyReview add user property update review.
func addPropertyReview(ctx *bm.Context) {
	arg := &model.ArgAddPropertyReview{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	form := ctx.Request.Form
	extra := form.Get("extra")
	if extra != "" {
		extraData := map[string]interface{}{}
		if err := json.Unmarshal([]byte(extra), &extraData); err != nil {
			log.Error("Failed to Unmarshal extra: %+v, error: %+v", extra, err)
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		arg.Extra = extraData
	}
	ctx.JSON(nil, memberSvc.AddPropertyReview(ctx, arg))
}
