package http

import (
	"strconv"

	"go-common/app/interface/main/reply/conf"
	model "go-common/app/interface/main/reply/model/reply"
	xmodel "go-common/app/interface/main/reply/model/xreply"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func subReplyByCursor(ctx *bm.Context) {
	params := ctx.Request.Form
	oid, err := strconv.ParseInt(params.Get("oid"), 10, 64)
	if err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	var build int64
	appStr := params.Get("mobi_app")
	buildStr := params.Get("build")
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(build %s) err(%v)", buildStr, err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
	}

	otyp, err := strconv.ParseInt(params.Get("type"), 10, 8)
	if err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	var size int64
	if params.Get("size") != "" {
		size, err = strconv.ParseInt(params.Get("size"), 10, 64)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if size <= 0 || size > int64(conf.Conf.Reply.MaxPageSize) {
		size = int64(conf.Conf.Reply.MaxPageSize)
	}
	var (
		rootID  int64
		replyID int64
		cursor  *model.Cursor
	)

	if params.Get("rpid") != "" {
		// jump subReply
		replyID, err = strconv.ParseInt(params.Get("rpid"), 10, 64)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		rootID, cursor, err = rpSvr.NewSubCursorByReplyID(ctx, oid, int8(otyp), replyID, int(size), model.OrderASC)
		if err != nil {
			ctx.JSON(nil, err)
			return
		}
	} else {
		var maxID, minID int64
		if params.Get("min_id") != "" {
			minID, err = strconv.ParseInt(params.Get("min_id"), 10, 64)
			if err != nil {
				ctx.JSON(nil, ecode.RequestErr)
				return
			}
		}
		if params.Get("max_id") != "" {
			maxID, err = strconv.ParseInt(params.Get("max_id"), 10, 64)
			if err != nil {
				ctx.JSON(nil, ecode.RequestErr)
				return
			}
		}
		cursor, err = model.NewCursor(maxID, minID, int(size), model.OrderASC)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if rootID <= 0 {
		rootID, err = strconv.ParseInt(params.Get("root"), 10, 64)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
	}
	var (
		showFolded bool
		scene      int64
	)

	// 消息中心折叠的评论也要给加回来，这里约定他们传一个scene=1
	if params.Get("scene") != "" {
		scene, _ = strconv.ParseInt(params.Get("scene"), 10, 64)
	}
	if model.ShouldShowFolded(params.Get("mobi_app"), build, scene) {
		showFolded = true
	}
	// 这里老版本折叠评论也要显示
	cursorParams := &model.CursorParams{
		IP:         metadata.String(ctx, metadata.RemoteIP),
		Oid:        oid,
		RootID:     rootID,
		OTyp:       int8(otyp),
		HTMLEscape: params.Get("mobi_app") == "",
		Cursor:     cursor,
		ShowFolded: showFolded,
	}

	if m, ok := ctx.Get("mid"); ok {
		cursorParams.Mid = m.(int64)
	}
	cursorRes, err := rpSvr.GetSubReplyListByCursor(ctx, cursorParams)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	var blacklist, assist int
	if mid := cursorParams.Mid; mid > 0 {
		if !rpSvr.IsWhiteAid(cursorParams.Oid, int8(cursorParams.OTyp)) {
			if rpSvr.RelationBlocked(ctx, cursorRes.Subject.Mid, mid) {
				blacklist = 1
			}
			if ok, _ := rpSvr.CheckAssist(ctx, cursorRes.Subject.Mid, mid); ok {
				assist = 1
			}
		}
	}
	var config xmodel.ReplyConfig
	config.ShowFloor = 1
	if !rpSvr.ShowFloor(cursorParams.Oid, cursorParams.OTyp) {
		config.ShowFloor = 0
	}
	rootReply := cursorRes.Roots[0]
	if showFolded {
		rootReply.Folder.HasFolded = false
	}
	rpSvr.EmojiReplaceI(appStr, build, rootReply)
	data := map[string]interface{}{
		"assist":    assist,
		"blacklist": blacklist,
		"upper": map[string]interface{}{
			"mid": cursorRes.Subject.Mid,
		},
		"root": rootReply,
		"cursor": map[string]interface{}{
			"all_count": rootReply.RCount,
			"max_id":    cursorRes.CursorRangeMax,
			"min_id":    cursorRes.CursorRangeMin,
			"size":      len(rootReply.Replies),
		},
		"config": config,
	}
	ctx.JSON(data, err)
}

func replyByCursor(ctx *bm.Context) {
	params := ctx.Request.Form
	buvid := ctx.Request.Header.Get("buvid")
	oid, err := strconv.ParseInt(params.Get("oid"), 10, 64)
	if err != nil {
		log.Warn("%v", err)
		err = ecode.RequestErr
		ctx.JSON(nil, err)
		return
	}

	otyp, err := strconv.ParseInt(params.Get("type"), 10, 8)
	if err != nil {
		log.Warn("%v", err)
		err = ecode.RequestErr
		ctx.JSON(nil, err)
		return
	}

	var sort int64
	if params.Get("sort") != "" {
		sort, err = strconv.ParseInt(params.Get("sort"), 10, 8)
		if err != nil {
			log.Warn("%v", err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
	}

	plat := int64(model.PlatWeb)
	if params.Get("plat") != "" {
		plat, err = strconv.ParseInt(params.Get("plat"), 10, 8)
		if err != nil {
			log.Warn("%v", err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
	}

	var build int64
	if params.Get("build") != "" {
		build, err = strconv.ParseInt(params.Get("build"), 10, 64)
		if err != nil {
			log.Warn("%v", err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
	}

	var size int64
	if params.Get("size") != "" {
		size, err = strconv.ParseInt(params.Get("size"), 10, 32)
		if err != nil {
			log.Warn("%v", err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
	}
	if size <= 0 || size > int64(conf.Conf.Reply.MaxPageSize) {
		size = int64(conf.Conf.Reply.MaxPageSize)
	}

	var (
		replyID int64
		cursor  *model.Cursor
	)

	if params.Get("rpid") != "" {
		// jump root reply
		replyID, err = strconv.ParseInt(params.Get("rpid"), 10, 64)
		if err != nil {
			log.Warn("%v", err)
			err = ecode.RequestErr
			ctx.JSON(nil, err)
			return
		}
		cursor, err = rpSvr.NewCursorByReplyID(ctx, oid, int8(otyp), replyID, int(size), model.OrderDESC)
	} else {
		var maxID, minID int64
		if params.Get("min_id") != "" {
			minID, err = strconv.ParseInt(params.Get("min_id"), 10, 64)
			if err != nil {
				log.Warn("%v", err)
				err = ecode.RequestErr
				ctx.JSON(nil, err)
				return
			}
		}
		if params.Get("max_id") != "" {
			maxID, err = strconv.ParseInt(params.Get("max_id"), 10, 64)
			if err != nil {
				log.Warn("%v", err)
				err = ecode.RequestErr
				ctx.JSON(nil, err)
				return
			}
		}
		cursor, err = model.NewCursor(maxID, minID, int(size), model.OrderDESC)
	}
	if err != nil {
		log.Warn("%v", err)
		err = ecode.RequestErr
		ctx.JSON(nil, err)
		return
	}
	appStr := params.Get("mobi_app")
	var showFolded bool
	// 对于根评论 scene传0, 只需要做版本的兼容即可
	if model.ShouldShowFolded(params.Get("mobi_app"), build, 0) {
		showFolded = true
	}
	cursorParams := &model.CursorParams{
		IP:         metadata.String(ctx, metadata.RemoteIP),
		Oid:        oid,
		OTyp:       int8(otyp),
		Sort:       int8(sort),
		HTMLEscape: appStr == "",
		Cursor:     cursor,
		HotSize:    5,
		Mid:        metadata.Int64(ctx, metadata.Mid),
		ShowFolded: showFolded,
	}

	if m, ok := ctx.Get("mid"); ok {
		cursorParams.Mid = m.(int64)
	}

	cursorRes, err := rpSvr.GetRootReplyListByCursor(ctx, cursorParams)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	var blacklist, assist int
	if mid := cursorParams.Mid; mid > 0 {
		if !rpSvr.IsWhiteAid(cursorParams.Oid, int8(cursorParams.OTyp)) {
			if rpSvr.RelationBlocked(ctx, cursorRes.Subject.Mid, mid) {
				blacklist = 1
			}
			if ok, _ := rpSvr.CheckAssist(ctx, cursorRes.Subject.Mid, mid); ok {
				assist = 1
			}
		}
	}
	rpSvr.EmojiReplace(int8(plat), build, cursorRes.Roots...)
	rpSvr.EmojiReplaceI(appStr, build, cursorRes.Roots...)
	cursorRes.Roots = rpSvr.FilDelReply(cursorRes.Roots)
	m := map[string]interface{}{
		"assist":    assist,
		"blacklist": blacklist,
		"replies":   cursorRes.Roots,
		"upper": map[string]interface{}{
			"mid": cursorRes.Subject.Mid,
		},
		"cursor": map[string]interface{}{
			"all_count": cursorRes.Subject.ACount,
			"max_id":    cursorRes.CursorRangeMax,
			"min_id":    cursorRes.CursorRangeMin,
			"size":      len(cursorRes.Roots),
		},
	}
	if cursorRes.Header != nil {
		rpSvr.EmojiReplace(int8(plat), build, cursorRes.Header.Hots...)
		rpSvr.EmojiReplace(int8(plat), build, cursorRes.Header.TopAdmin)
		rpSvr.EmojiReplace(int8(plat), build, cursorRes.Header.TopUpper)
		rpSvr.EmojiReplaceI(appStr, build, cursorRes.Header.Hots...)
		rpSvr.EmojiReplaceI(appStr, build, cursorRes.Header.TopAdmin)
		rpSvr.EmojiReplaceI(appStr, build, cursorRes.Header.TopUpper)
		cursorRes.Header.Hots = rpSvr.FilDelReply(cursorRes.Header.Hots)
		showEntry, showAdmin, showFloor := 1, 1, 1
		if config, _ := rpSvr.GetReplyLogConfig(ctx, cursorRes.Subject, 1); config != nil {
			showEntry = int(config.ShowEntry)
			showAdmin = int(config.ShowAdmin)
		}
		if !rpSvr.ShowFloor(cursorParams.Oid, cursorParams.OTyp) {
			showFloor = 0
		}
		m["config"] = map[string]int{
			"showentry": showEntry,
			"showadmin": showAdmin,
			"showfloor": showFloor,
		}
		if cursorRes.Subject.RCount <= 20 && len(cursorRes.Header.Hots) > 0 {
			cursorRes.Header.Hots = cursorRes.Header.Hots[:0]
		}
		m["hots"] = cursorRes.Header.Hots
		m["notice"] = rpSvr.RplyNotice(ctx, int8(plat), build, appStr, buvid)
		m["top"] = map[string]interface{}{
			"admin": cursorRes.Header.TopAdmin,
			"upper": cursorRes.Header.TopUpper,
		}
	}
	ctx.JSON(m, err)
}

func dialogByCursor(c *bm.Context) {
	var mid int64
	v := new(struct {
		Oid      int64  `form:"oid" validate:"required"`
		Type     int8   `form:"type" validate:"required"`
		Root     int64  `form:"root" validate:"required"`
		Dialog   int64  `form:"dialog" validate:"required"`
		Size     int    `form:"size" validate:"min=1"`
		MinFloor int64  `form:"min_floor"`
		MaxFloor int64  `form:"max_floor"`
		Plat     int64  `form:"plat"`
		Build    int64  `form:"build"`
		MobiApp  string `form:"mobi_app"`
	})
	// buvid := c.Request.Header.Get("buvid")
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	if v.Size > conf.Conf.Reply.MaxPageSize {
		v.Size = conf.Conf.Reply.MaxPageSize
	}
	if m, ok := c.Get("mid"); ok {
		mid = m.(int64)
	}
	cursor, err := model.NewCursor(v.MaxFloor, v.MinFloor, v.Size, model.OrderASC)
	if err != nil {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rps, dialogCursor, dialogMeta, err := rpSvr.DialogByCursor(c, mid, v.Oid, v.Type, v.Root, v.Dialog, cursor)
	if err != nil {
		log.Error("rpSvr.DialogByCursor error (%v)", err)
		c.JSON(nil, err)
		return
	}
	rpSvr.EmojiReplaceI(v.MobiApp, v.Build, rps...)
	var config xmodel.ReplyConfig
	config.ShowFloor = 1
	if !rpSvr.ShowFloor(v.Oid, v.Type) {
		config.ShowFloor = 0
	}
	data := map[string]interface{}{
		"cursor":  dialogCursor,
		"dialog":  dialogMeta,
		"replies": rps,
		"config":  config,
	}
	c.JSON(data, nil)
}
