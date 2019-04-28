package http

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/adminlog"
	model "go-common/app/interface/main/reply/model/reply"
	xmodel "go-common/app/interface/main/reply/model/xreply"
	"go-common/app/interface/main/reply/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// regexp utf8 char 0x0e0d~0e4A
var (
	_emptyUnicodeReg = []*regexp.Regexp{
		regexp.MustCompile(`[\x{202e}]+`),  // right-to-left override
		regexp.MustCompile(`[\x{200b}]+`),  // zeroWithChar
		regexp.MustCompile(`[\x{1f6ab}]+`), // no_entry_sign
	}
	re      = regexp.MustCompile(`[\x{0E0D}\x{0E4A}]+`)
	_emojis = regexp.MustCompile(`\[[^\[+][^]]+]`)
	// trim
	returnReg  = regexp.MustCompile(`[\n]{3,}`)
	returnReg2 = regexp.MustCompile(`(\r\n){3,}`)
	spaceReg   = regexp.MustCompile(`[　]{5,}`) // Chinese quanjiao space character
)

func isMobile(params url.Values) bool {
	//return (params.Get("appkey") == "c1b107428d337928" || params.Get("appkey") == "27eb53fc9058f8c3") && len(params.Get("access_key")) == 32
	return params.Get("mobi_app") != ""
}

// replyInfo get reply info by rpID.
func replyInfo(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpStr := params.Get("rpid")
	// check params
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rp, err := rpSvr.Reply(c, oid, int8(tp), rpID)
	if err != nil {
		log.Warn("rpSvr.Reply(%d, %d, %d) error(%s)", oid, tp, rpID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rp, nil)
}

func replyMultiInfo(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpStr := params.Get("rpid")
	// check params
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("xstr.SplitInts(%s) err(%v)", err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpStr)
	if err != nil {
		log.Warn("xstr.SplintInt(%s) err(%v)", rpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("miss oid:%v rpid:%v", oidStr, rpStr)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	data := map[string]*model.Reply{}
	for i := 0; i < len(oids); i++ {
		rp, err := rpSvr.Reply(c, oids[i], int8(tp), rpIDs[i])
		if err != nil {
			log.Warn("rpSvr.Reply(%d, %d, %d) error(%s)", oids[i], tp, rpIDs[i], err)
			continue
		}
		data[strconv.FormatInt(rpIDs[i], 10)] = rp
	}
	c.JSON(data, nil)
}

// reply range subject replies.
func reply(c *bm.Context) {
	var (
		showEntry = int(1)
		showAdmin = int(1)
		showFloor = int(1)
		mid       int64
		sort      int64
		build     int64
		plat      int64
		oid       int64
		curPage   int
		err       error
		nohot     bool
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	otypStr := params.Get("type")
	sortStr := params.Get("sort")
	curPageStr := params.Get("pn")
	perPageStr := params.Get("ps")
	nohotStr := params.Get("nohot")
	platStr := params.Get("plat")
	buildStr := params.Get("build")
	appStr := params.Get("mobi_app")
	buvid := c.Request.Header.Get("buvid")
	if m, ok := c.Get("mid"); ok {
		mid = m.(int64)
	}
	if platStr == "" {
		plat = int64(model.PlatWeb)
	} else {
		plat, err = strconv.ParseInt(platStr, 10, 8)
		if err != nil {
			log.Warn("strconv.ParseInt(platStr %s) err(%v)", platStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(build %s) err(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	oid, err = strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	otyp, err := strconv.ParseInt(otypStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", otypStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if sortStr != "" {
		sort, err = strconv.ParseInt(sortStr, 10, 8)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", sortStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	curPage, err = strconv.Atoi(curPageStr)
	if err != nil || curPage < 1 {
		curPage = 1
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage > conf.Conf.Reply.MaxPageSize || perPage <= 0 {
		perPage = conf.Conf.Reply.MaxPageSize
	}
	htmlEscape := true
	needSndReply := true
	// check android and ios appkey
	if appStr != "" {
		// if mobile, no html escape
		htmlEscape = false
		// if mobile, no return second replies of hot and top reply
		if int8(sort) == model.SortByFloor {
			needSndReply = false
		}
	}
	if nohotStr != "" {
		nohot, _ = strconv.ParseBool(nohotStr)
	}
	pageParams := &model.PageParams{
		Mid:        mid,
		Oid:        oid,
		Type:       int8(otyp),
		Sort:       int8(sort),
		PageNum:    curPage,
		PageSize:   perPage,
		NeedSecond: needSndReply,
		Escape:     htmlEscape,
		NeedHot:    !nohot,
	}
	pageRes, err := rpSvr.RootReplies(c, pageParams)
	if err != nil {
		log.Warn("rpSvr.RootReplies(%+v) error(%v)", pageParams, err)
		c.JSON(nil, err)
		return
	}
	if config, _ := rpSvr.GetReplyLogConfig(c, pageRes.Subject, 1); config != nil {
		showEntry = int(config.ShowEntry)
		showAdmin = int(config.ShowAdmin)
	}
	if !rpSvr.ShowFloor(pageParams.Oid, pageParams.Type) {
		showFloor = int(0)
		if !isMobile(params) {
			rpSvr.ResetFloor(pageRes.Roots...)
			rpSvr.ResetFloor(pageRes.TopAdmin)
			rpSvr.ResetFloor(pageRes.TopUpper)
			rpSvr.ResetFloor(pageRes.Hots...)
		}
	}
	rpSvr.EmojiReplace(int8(plat), build, pageRes.Roots...)
	rpSvr.EmojiReplace(int8(plat), build, pageRes.TopAdmin)
	rpSvr.EmojiReplace(int8(plat), build, pageRes.TopUpper)
	rpSvr.EmojiReplace(int8(plat), build, pageRes.Hots...)
	rpSvr.EmojiReplaceI(appStr, build, pageRes.Roots...)
	rpSvr.EmojiReplaceI(appStr, build, pageRes.TopAdmin)
	rpSvr.EmojiReplaceI(appStr, build, pageRes.TopUpper)
	rpSvr.EmojiReplaceI(appStr, build, pageRes.Hots...)
	pageRes.Roots = rpSvr.FilDelReply(pageRes.Roots)
	pageRes.Hots = rpSvr.FilDelReply(pageRes.Hots)
	data := map[string]interface{}{
		"page": map[string]int{
			"num":    curPage,
			"size":   perPage,
			"count":  pageRes.Total,
			"acount": pageRes.AllCount,
		},
		"config": map[string]int{
			"showentry": showEntry,
			"showadmin": showAdmin,
			"showfloor": showFloor,
		},
		"replies": pageRes.Roots,
		"hots":    pageRes.Hots,
		"upper": map[string]interface{}{
			"mid": pageRes.Subject.Mid,
			"top": pageRes.TopUpper,
		},
		"top":    pageRes.TopAdmin,
		"notice": rpSvr.RplyNotice(c, int8(plat), build, appStr, buvid),
	}
	if mid > 0 {
		if !(rpSvr.IsWhiteAid(oid, int8(otyp))) {
			if rpSvr.RelationBlocked(c, pageRes.Subject.Mid, mid) {
				data["blacklist"] = 1
			}
			if ok, _ := rpSvr.CheckAssist(c, pageRes.Subject.Mid, mid); ok {
				data["assist"] = 1
			}
		}
	}
	c.JSON(data, nil)
}

// replyReply range replies from root reply.
func replyReply(c *bm.Context) {
	var (
		err    error
		mid    int64
		root   int64
		jumpID int64
		escape = true
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rtStr := params.Get("root")
	jumpStr := params.Get("jump")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	midS, ok := c.Get("mid")
	if !ok {
		log.Warn("user no login")
		mid = 0
	} else {
		mid = midS.(int64)
	}
	// check params
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > conf.Conf.Reply.MaxPageSize || ps <= 0 {
		ps = conf.Conf.Reply.MaxPageSize
	}
	if jumpStr != "" {
		if jumpID, err = strconv.ParseInt(jumpStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(jump:%s) error(%v)", jumpStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	} else {
		if root, err = strconv.ParseInt(rtStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", rtStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	// check android and ios appkey
	if isMobile(params) {
		// if mobile, no html escape
		escape = false
	}
	var config xmodel.ReplyConfig
	config.ShowFloor = 1
	if !rpSvr.ShowFloor(oid, int8(tp)) {
		config.ShowFloor = 0
	}
	rs, rtRp, umid, pn, err := rpSvr.SecondReplies(c, mid, oid, root, jumpID, int8(tp), pn, ps, escape)
	if err == nil {
		if config.ShowFloor == 0 && !isMobile(params) {
			rpSvr.ResetFloor(rs...)
			rpSvr.ResetFloor(rtRp)
		}
		data := make(map[string]interface{}, 2)
		data["page"] = map[string]int{
			"num":   pn,
			"size":  ps,
			"count": rtRp.RCount,
		}
		data["upper"] = map[string]int64{
			"mid": umid,
		}
		data["replies"] = rs
		data["root"] = rtRp
		data["config"] = config
		if mid > 0 {
			if !(rpSvr.IsWhiteAid(oid, int8(tp))) {
				if rpSvr.RelationBlocked(c, rtRp.Mid, mid) {
					data["blacklist"] = 1
				}
				if ok, _ := rpSvr.CheckAssist(c, umid, mid); ok {
					data["assist"] = 1
				}
			}
		}
		c.JSON(data, nil)
	} else {
		log.Warn("rpSvr.ReplyReplies(%d, %d, %d) error(%d)", oid, root, tp, err)
		c.JSON(nil, err)
	}
}

func getTopics(c *bm.Context) {
	var (
		err error
		mid int64
	)
	params := c.Request.Form
	midStr, ok := c.Get("mid")
	if ok {
		mid = midStr.(int64)
	} else {
		midStr := params.Get("mid")
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", midStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	msg := params.Get("message")
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if !model.LegalSubjectType(int8(tp)) {
		err = ecode.ReplyIllegalSubType
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if int8(tp) == model.SubTypeDynamic || int8(tp) == model.SubTypeLivePicture || int8(tp) == model.SubTypeLive {
		data := map[string]interface{}{
			"topics": []string{},
		}
		c.JSON(data, nil)
		return
	}
	topics, err := rpSvr.Topics(c, mid, oid, int8(tp), msg)
	if err != nil {
		log.Warn("rpSvr.Topics(%d,%d,%d,%s) error(%v)", mid, oid, tp, msg, err)
		topics = make([]string, 0)
	}
	data := map[string]interface{}{
		"topics": topics,
	}
	c.JSON(data, nil)
}

// addReply add a reply.
func addReply(c *bm.Context) {
	var (
		err          error
		rp           *model.Reply
		captchaURL   string
		ats          []int64
		root, parent int64
		plat         int64
		build        int64
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rtStr := params.Get("root")
	paStr := params.Get("parent")
	atStr := params.Get("at")
	msg := params.Get("message")
	platStr := params.Get("plat")
	device := params.Get("device")
	version := params.Get("version")
	captcha := params.Get("code")
	ak := params.Get("access_key")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	// check params
	msg = strings.TrimSpace(msg)
	msg = re.ReplaceAllString(msg, " ")
	msg = spaceReg.ReplaceAllString(msg, "　　　")
	msg = returnReg.ReplaceAllString(msg, "\n\n\n")
	msg = returnReg2.ReplaceAllString(msg, "\n\n\n")
	// checkout empty
	for _, reg := range _emptyUnicodeReg {
		msg = reg.ReplaceAllString(msg, "")
	}
	// check len
	tmp := _emojis.ReplaceAllString(msg, "")
	tmp = strings.TrimSpace(tmp)
	ml := len([]rune(tmp))
	if conf.Conf.Reply.MaxConLen < ml || ml < conf.Conf.Reply.MinConLen {
		log.Warn("content(%s) length %d, max %d, min %d", msg, ml, conf.Conf.Reply.MaxConLen, conf.Conf.Reply.MinConLen)
		err = ecode.ReplyContentOver
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if rtStr != "" {
		root, err = strconv.ParseInt(rtStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", rtStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if paStr != "" {
		parent, err = strconv.ParseInt(paStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", paStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if !((root == 0 && parent == 0) || (root > 0 && parent > 0)) {
		log.Warn("the wrong root(%d) and parent(%d)", root, parent)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if !model.LegalSubjectType(int8(tp)) {
		err = ecode.ReplyIllegalSubType
		c.JSON(nil, err)
		return
	}
	if atStr != "" {
		ats, err = xstr.SplitInts(atStr)
		if err != nil {
			log.Warn("utils.SplitInts(%s) error(%v)", atStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if len(ats) > 10 {
		log.Warn("too many people to be at len(%d)", len(ats))
		err = ecode.ReplyTooManyAts
		c.JSON(nil, err)
		return
	}
	if platStr != "" {
		plat, err = strconv.ParseInt(platStr, 10, 8)
		if err != nil || !model.CheckPlat(int8(plat)) {
			log.Warn("the wrong plat(%s)", platStr)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if device != "" && len([]rune(device)) > 30 {
		log.Warn("device len(%d)", len([]rune(device)))
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if int8(tp) == model.SubTypeArchive {
		for _, aid := range cnf.Reply.ForbidList {
			if aid == oid {
				err = ecode.ReplyForbidList
				c.JSON(nil, err)
				return
			}
		}
	}
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	var status int
	if status, err = rpSvr.UserBlockStatus(c, mid); err != nil {
		err = ecode.ServerErr
		c.JSON(nil, err)
		return
	}
	if status == service.StatusForbidden {
		// "账号被封停"
		err = ecode.UserDisabled
		c.JSON(nil, err)
		return
	}
	if status == service.StatusNeedContest {
		// 该账号处于封禁中，点击申请答题
		err = ecode.ReplyContestNotPassed
		c.JSON(nil, err)
		return
	}
	if root == 0 && parent == 0 {
		rp, captchaURL, err = rpSvr.AddReply(c, mid, oid, int8(tp), int8(plat), ats, ak, c.Request.Header.Get("Cookie"), captcha, msg, device, version, platform, build, buvid)
	} else {
		rp, captchaURL, err = rpSvr.AddReplyReply(c, mid, oid, root, parent, int8(tp), int8(plat), ats, ak, c.Request.Header.Get("Cookie"), captcha, msg, device, version, platform, build, buvid)
	}
	if err != nil && err != ecode.ReplyMosaicByFilter {
		log.Warn("rpSvr.AddReply or ReplyReply failed mid(%d) oid(%d) error(%d)", mid, oid, err)
		data := map[string]interface{}{
			"need_captcha": (captchaURL != ""),
			"url":          captchaURL,
		}
		c.JSON(data, err)
		return
	}
	data := map[string]interface{}{
		"rpid":       rp.RpID,
		"rpid_str":   strconv.FormatInt(rp.RpID, 10),
		"dialog":     rp.Dialog,
		"dialog_str": strconv.FormatInt(rp.Dialog, 10),
		"root":       rp.Root,
		"root_str":   strconv.FormatInt(rp.Root, 10),
		"parent":     rp.Parent,
		"parent_str": strconv.FormatInt(rp.Parent, 10),
	}
	c.JSON(data, nil)
}

// likeReply do like or cancel like for reply.
func likeReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpStr := params.Get("rpid")
	actStr := params.Get("action")
	ak := params.Get("access_key")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	var build int64
	var err error
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	// check parameters
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	act, err := strconv.ParseInt(actStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", actStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	err = rpSvr.AddAction(c, mid.(int64), oid, rpID, int8(tp), int8(act), ak, c.Request.Header.Get("Cookie"), "like", platform, buvid, build)
	c.JSON(nil, err)
}

// hateReply do hate or cancel hate for reply.
func hateReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpStr := params.Get("rpid")
	actStr := params.Get("action")
	ak := params.Get("access_key")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	var build int64
	var err error
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	// check parameters
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	act, err := strconv.ParseInt(actStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", actStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	err = rpSvr.AddAction(c, mid.(int64), oid, rpID, int8(tp), int8(act), ak, c.Request.Header.Get("Cookie"), "hate", platform, buvid, build)
	c.JSON(nil, err)
}

// jump to specified reply.
func jumpReply(c *bm.Context) {
	var (
		showEntry = int(1)
		showAdmin = int(1)
		mid       int64
		escape    = true
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	psStr := params.Get("ps")
	rppsStr := params.Get("rpps")
	platStr := params.Get("plat")
	buildStr := params.Get("build")
	appStr := params.Get("mobi_app")
	midS, ok := c.Get("mid")
	buvid := c.Request.Header.Get("buvid")
	if !ok {
		log.Warn("user no login")
		mid = 0
	} else {
		mid = midS.(int64)
	}
	// check params
	plat, err := strconv.ParseInt(platStr, 10, 8)
	if err != nil {
		plat = 1 // default pc
	}
	var build int64
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(build %s) err(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpIDStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	ps, _ := strconv.Atoi(psStr)
	if ps > conf.Conf.Reply.MaxPageSize || ps <= 0 {
		ps = conf.Conf.Reply.MaxPageSize
	}
	rpPs, _ := strconv.Atoi(rppsStr)
	if rpPs <= 0 || rpPs > conf.Conf.Reply.MaxPageSize {
		rpPs = 10
	}
	// check android and ios appkey
	if isMobile(params) {
		// if mobile, no html escape
		escape = false
	}
	rs, hots, topAdmin, topUpper, sub, pn, rtpn, total, err := rpSvr.JumpReplies(c, mid, oid, rpID, int8(tp), ps, rpPs, escape)
	if err != nil && err != ecode.ReplyNotExist {
		log.Warn("rpSvr.JumpReplies(%d,%d,%d,%d) error(%d)", oid, tp, rpID, ps, err)
		c.JSON(nil, err)
		return
	}
	if err == ecode.ReplyNotExist {
		var pageRes *model.PageResult
		pn = 1
		pageParams := &model.PageParams{
			Mid: mid, Oid: oid, Type: int8(tp), Sort: model.SortByFloor, PageNum: pn, PageSize: ps, NeedHot: false, NeedSecond: false, Escape: escape,
		}
		if pageRes, err = rpSvr.RootReplies(c, pageParams); err != nil {
			log.Warn("rpSvr.RootReplies(%d,%d,%d,%d) error(%d)", oid, tp, rpID, ps, err)
			c.JSON(nil, err)
			return
		}
		rs = pageRes.Roots
		sub = pageRes.Subject
		hots = pageRes.Hots
		topAdmin = pageRes.TopAdmin
		topUpper = pageRes.TopUpper
		total = pageRes.Total
	}
	if config, _ := rpSvr.GetReplyLogConfig(c, sub, 1); config != nil {
		showEntry = int(config.ShowEntry)
		showAdmin = int(config.ShowAdmin)
	}
	rpSvr.EmojiReplace(int8(plat), build, hots...)
	rpSvr.EmojiReplace(int8(plat), build, rs...)
	rpSvr.EmojiReplace(int8(plat), build, topAdmin)
	rpSvr.EmojiReplace(int8(plat), build, topUpper)
	rpSvr.EmojiReplaceI(appStr, build, hots...)
	rpSvr.EmojiReplaceI(appStr, build, rs...)
	rpSvr.EmojiReplaceI(appStr, build, topAdmin)
	rpSvr.EmojiReplaceI(appStr, build, topUpper)
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":    pn,
		"size":   ps,
		"count":  total,
		"acount": sub.ACount,
		"rt_num": rtpn,
	}
	configValue := map[string]int{
		"showentry": showEntry,
		"showadmin": showAdmin,
	}
	upper := map[string]interface{}{
		"mid": sub.Mid,
		"top": topUpper,
	}
	data["config"] = configValue
	data["upper"] = upper
	data["page"] = page
	data["replies"] = rs
	data["hots"] = hots
	data["top"] = topAdmin
	if notice := rpSvr.RplyNotice(c, int8(plat), build, appStr, buvid); notice != nil {
		data["notice"] = notice
	}
	//NOTE donot need
	// data["need_captcha"] = rpSvr.NeedCaptcha(c, mid)
	c.JSON(data, nil)
}

// showReply show reply by upper.
func showReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidsStr := params.Get("oid")
	rpsStr := params.Get("rpid")
	tpStr := params.Get("type")
	ak := params.Get("access_key")
	// check params
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("oids(%v) not equal roids(%v)", oids, rpIDs)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("miss oid:%s rpid:%s", oidsStr, rpsStr)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		if oids[i] <= 0 {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		if err = rpSvr.Show(c, oids[i], mid.(int64), rpIDs[i], int8(tp), ak, c.Request.Header.Get("Cookie")); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
}

// hideReplys hide reply by upper.
func hideReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidsStr := params.Get("oid")
	rpsStr := params.Get("rpid")
	tpStr := params.Get("type")
	ak := params.Get("access_key")
	// check params
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("oids(%v) not equal roids(%v)", oids, rpIDs)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("miss oid:%s rpid:%s", oidsStr, rpsStr)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		if oids[i] <= 0 {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		if err = rpSvr.Hide(c, oids[i], mid.(int64), rpIDs[i], int8(tp), ak, c.Request.Header.Get("Cookie")); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
}

func replyCounts(c *bm.Context) {
	params := c.Request.Form
	oids, err := xstr.SplitInts(params.Get("oid"))
	if err != nil {
		log.Warn("sxtr.PlintInts(%v) err(%v)", params.Get("oid"), err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) == 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		if oids[i] <= 0 {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	otyp, err := strconv.ParseInt(params.Get("type"), 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%v) err(%v)", params.Get("type"), err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if !model.LegalSubjectType(int8(otyp)) {
		err = ecode.ReplyIllegalSubType
		c.JSON(nil, err)
		return
	}
	counts, err := rpSvr.GetReplyCounts(c, oids, int8(otyp))
	if err != nil {
		log.Warn("rcSvr.GetReplyCount err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(counts, nil)
}

// replyCount get replies count.
func replyMultiCount(c *bm.Context) {
	params := c.Request.Form
	oidsStr := params.Get("oid")
	tpStr := params.Get("type")
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Warn("sxtr.PlintInts(%v) err(%v)", oidsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%v) err(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) == 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		if oids[i] <= 0 {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	counts, err := rpSvr.ReplyCounts(c, oids, int8(tp))
	if err != nil {
		log.Warn("rcSvr.ReplyCount err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(counts, nil)
}

// replyAdminLog fetch pages of replies deleted by upper, reply user or administrators
func replyAdminLog(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	curPageStr := params.Get("pn")
	perPageStr := params.Get("ps")
	// check params
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}

	curPage, err := strconv.Atoi(curPageStr)
	if err != nil || curPage < 1 {
		curPage = 1
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage > conf.Conf.Reply.MaxPageSize || perPage <= 0 {
		perPage = conf.Conf.Reply.MaxPageSize
	}
	if rpSvr.IsBnj(oid, int8(tp)) {
		data := map[string]interface{}{
			"page": map[string]int{
				"num":   curPage,
				"size":  perPage,
				"pages": 0,
				"total": 0,
			},
			"logs":         make([]*adminlog.AdminLog, 0),
			"reply_count":  0,
			"report_count": 0,
		}
		c.JSON(data, nil)
		return
	}
	logs, replyCount, reportCount, pageCount, total, err := rpSvr.PaginateUpperDeletedLogs(c, oid, int(tp), curPage, perPage)
	if err != nil {
		log.Warn("rpSvr.PaginateUpperDeletedLogs(%d, %d) error(%s)", oid, tp, err)
		c.JSON(nil, err)
		return
	}
	if replyCount > reportCount {
		reportCount += replyCount
	}
	data := map[string]interface{}{
		"page": map[string]int{
			"num":   curPage,
			"size":  perPage,
			"pages": int(pageCount),
			"total": int(total),
		},
		"logs":         logs,
		"reply_count":  replyCount,
		"report_count": reportCount,
	}
	c.JSON(data, nil)
}

// replyCount get replies count.
func replyCount(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	// check params
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	count, err := rpSvr.ReplyCount(c, oid, int8(tp))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]int{"count": count}
	c.JSON(data, nil)
}

// delReply delete reply by upper or self.
func delReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidsStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDsStr := params.Get("rpid")
	ak := params.Get("access_key")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	var build int64
	var err error
	// check params
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(build %s) err(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIDs) {
		log.Warn("miss oid:%s rpid:%s", oidsStr, rpIDsStr)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		if oids[i] <= 0 {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		if err = rpSvr.Delete(c, mid.(int64), oids[i], rpIDs[i], int8(tp), ak, c.Request.Header.Get("Cookie"), platform, build, buvid); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
}

// AddTopReply add top reply by upper
func AddTopReply(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	actStr := params.Get("action")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	var build int64
	var err error
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	mid, _ := c.Get("mid")
	// check params
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpIDStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	act, err := strconv.ParseInt(actStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(actStr :%s) err(%v)", actStr, err)
		act = 1
	}
	err = rpSvr.UpperAddTop(c, mid.(int64), oid, rpID, int8(tp), int8(act), platform, build, buvid)
	c.JSON(nil, err)
}

// emojis
func emojis(c *bm.Context) {
	params := c.Request.Form
	appStr := params.Get("mobi_app")
	buildStr := params.Get("build")
	if buildStr != "" {
		build, err := strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(build %s) err(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		if appStr == "android_i" && build > 1125000 && build < 2005000 {
			c.JSON(nil, nil)
			return
		}
	}
	c.JSON(rpSvr.Emojis(c), nil)
}

func replyHots(c *bm.Context) {
	params := c.Request.Form
	oidStr := params.Get("oid")
	typStr := params.Get("type")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil || oid <= 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	typ, err := strconv.ParseInt(typStr, 10, 8)
	if err != nil {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	pn, _ := strconv.Atoi(pnStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps > conf.Conf.Reply.MaxPageSize || ps <= 0 {
		ps = 3
	}
	sub, rps, err := rpSvr.ReplyHots(c, oid, int8(typ), pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	rps = rpSvr.FilDelReply(rps)
	data := map[string]interface{}{
		"page": map[string]int{
			"num":    pn,
			"size":   ps,
			"count":  sub.RCount,
			"acount": sub.ACount,
		},
		"replies": rps,
	}
	c.JSON(data, nil)
}

func dialog(c *bm.Context) {
	var (
		err    error
		mid    int64
		escape = true
	)

	v := new(struct {
		Oid     int64  `form:"oid" validate:"required"`
		Type    int8   `form:"type" validate:"required"`
		Dialog  int64  `form:"dialog" validate:"required"`
		Root    int64  `form:"root" validate:"required"`
		Pn      int    `form:"pn" validate:"min=1"`
		Ps      int    `form:"ps" validate:"min=1"`
		Plat    string `form:"plat"`
		Build   string `form:"build"`
		MobiApp string `form:"mobi_app"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Ps > conf.Conf.Reply.MaxPageSize {
		v.Ps = conf.Conf.Reply.MaxPageSize
	}
	midS, ok := c.Get("mid")
	if !ok {
		log.Warn("user no login")
	} else {
		mid = midS.(int64)
	}
	rps, err := rpSvr.Dialog(c, mid, v.Oid, v.Type, v.Root, v.Dialog, v.Pn, v.Ps, escape)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page": map[string]int{
			"num":   v.Pn,
			"size":  v.Ps,
			"total": len(rps),
		},
		// "dialog":,
		"replies": rps,
	}
	c.JSON(data, nil)

}

func isHotReply(c *bm.Context) {
	v := new(struct {
		Type int8  `form:"type" validate:"required"`
		Oid  int64 `form:"oid" validate:"required"`
		RpID int64 `form:"rpid" validate:"required"`
	})
	var (
		isHot bool
		err   error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if isHot, err = rpSvr.IsHotReply(c, v.Type, v.Oid, v.RpID); err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]bool{
		"isHot": isHot,
	}
	c.JSON(data, nil)
}

func hotsBatch(c *bm.Context) {
	v := new(struct {
		Type int8    `form:"type" validate:"required"`
		Oids []int64 `form:"oids,split" validate:"required"`
		Size int8    `form:"size" default:"1" validate:"min=1"`
		Mid  int64   `form:"mid"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(rpSvr.HotsBatch(c, v.Type, v.Size, v.Oids, v.Mid))
}
