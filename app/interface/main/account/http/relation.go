package http

import (
	"strconv"

	"go-common/app/interface/main/account/model"
	mrl "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

var (
	_defaultPS int64 = 50
)

// modify modify user relation.
func modify(c *bm.Context) {
	var (
		err      error
		act, fid int64
		src      uint64
		params   = c.Request.Form
		mid, _   = c.Get("mid")
		actStr   = params.Get("act")
		fidStr   = params.Get("fid")
		srcStr   = params.Get("re_src")
		ua       = c.Request.Header.Get("User-Agent")
		referer  = c.Request.Header.Get("Referer")
		sid      string
		realIP   = metadata.String(c, metadata.RemoteIP)
	)
	if act, err = strconv.ParseInt(actStr, 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if src, err = strconv.ParseUint(srcStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Warn("relation infoc get sid failed error(%v)", err)
	} else {
		sid = sidCookie.Value
	}
	buvid := c.Request.Header.Get("Buvid")
	if buvid == "" {
		buvidCookie, _ := c.Request.Cookie("buvid3")
		if buvidCookie != nil {
			buvid = buvidCookie.Value
		}
	}
	ric := map[string]string{
		"ip":         realIP,
		"User-Agent": ua,
		"sid":        sid,
		"buvid":      buvid,
		"Referer":    referer,
	}
	c.JSON(nil, relationSvc.Modify(c, mid.(int64), fid, int8(act), uint8(src), ric))
}

func batchModify(c *bm.Context) {
	var (
		err    error
		act    int64
		fids   []int64
		src    uint64
		params = c.Request.Form
		// res      = c.Result()
		mid, _  = c.Get("mid")
		actStr  = params.Get("act")
		fidsStr = params.Get("fids")
		srcStr  = params.Get("re_src")
		ua      = c.Request.Header.Get("User-Agent")
		referer = c.Request.Header.Get("Referer")
		sid     string
		realIP  = metadata.String(c, metadata.RemoteIP)
	)
	if act, err = strconv.ParseInt(actStr, 10, 8); err != nil {
		// res["code"] = ecode.RequestErr
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fids, err = xstr.SplitInts(fidsStr); err != nil || len(fids) <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if src, err = strconv.ParseUint(srcStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Warn("relation infoc get sid failed error(%v)", err)
	} else {
		sid = sidCookie.Value
	}
	buvid := c.Request.Header.Get("Buvid")
	if buvid == "" {
		buvidCookie, _ := c.Request.Cookie("buvid3")
		if buvidCookie != nil {
			buvid = buvidCookie.Value
		}
	}
	ric := map[string]string{
		"ip":         realIP,
		"User-Agent": ua,
		"sid":        sid,
		"buvid":      buvid,
		"Referer":    referer,
	}
	// res["code"] = relationSvc.Modify(c, mid.(int64), fid, int8(act), uint8(src), ric)
	c.JSON(relationSvc.BatchModify(c, mid.(int64), fids, int8(act), uint8(src), ric))
}

// relation get relation between mid and fid.
func relation(c *bm.Context) {
	var (
		err    error
		fid    int64
		f      *mrl.Following
		params = c.Request.Form
		fidStr = params.Get("fid")
		mid, _ = c.Get("mid")
	)
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if f, err = relationSvc.Relation(c, mid.(int64), fid); err != nil {
		log.Error("relationSvc.Relation(%d, %d) error(%v)", mid, fid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(f, nil)
}

// relations get relations between mid and fids.
func relations(c *bm.Context) {
	var (
		err     error
		fm      map[int64]*mrl.Following
		fids    []int64
		params  = c.Request.Form
		fidsStr = params.Get("fids")
		mid, _  = c.Get("mid")
	)
	if fids, err = xstr.SplitInts(fidsStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("xstr.SplitInts(fids %v) err(%v)", fidsStr, err)
		return
	}
	if fm, err = relationSvc.Relations(c, mid.(int64), fids); err != nil {
		log.Error("relationSvc.Relations(%d, %v) error(%v)", mid, fids, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(fm, nil)
}

// followings get user's following list.
func followings(c *bm.Context) {
	var (
		err        error
		mid, vmid  int64
		pn, ps     int64
		self       bool
		followings []*model.Following
		params     = c.Request.Form
		vmidStr    = params.Get("vmid")
		psStr      = params.Get("ps")
		pnStr      = params.Get("pn")
		order      = params.Get("order")
		version    uint64
		versionStr = params.Get("re_version")
		crc32v     uint32
		total      int
	)
	midS, ok := c.Get("mid")
	if ok {
		mid = midS.(int64)
	} else {
		mid = 0
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	self = mid == vmid
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn <= 0 {
		pn = 1
	}
	if !self && pn > 5 {
		c.JSON(nil, ecode.RelFollowingGuestLimit)
		return
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ps <= 0 || ps > _defaultPS {
		ps = _defaultPS
	}
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if order != "asc" {
		order = "desc"
	}
	if followings, crc32v, total, err = relationSvc.Followings(c, vmid, mid, pn, ps, version, order); err != nil {
		log.Error("relationSvc.Followings(%d) error(%v)", vmid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       followings,
		"total":      total,
	}, nil)
}

// followers get user's follower list.
func followers(c *bm.Context) {
	var (
		err        error
		mid, vmid  int64
		pn, ps     int64
		self       bool
		fs         []*model.Following
		params     = c.Request.Form
		vmidStr    = params.Get("vmid")
		psStr      = params.Get("ps")
		pnStr      = params.Get("pn")
		version    uint64
		total      int
		versionStr = params.Get("re_version")
		crc32v     uint32
	)
	midS, ok := c.Get("mid")
	if ok {
		mid = midS.(int64)
	} else {
		mid = 0
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	self = mid == vmid
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn <= 0 {
		pn = 1
	}
	if !self && pn > 5 {
		c.JSON(nil, ecode.RelFollowingGuestLimit)
		return
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ps <= 0 || ps > _defaultPS {
		ps = _defaultPS
	}
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if fs, crc32v, total, err = relationSvc.Followers(c, vmid, mid, pn, ps, version); err != nil {
		log.Error("relationSvc.Followers(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       fs,
		"total":      total,
	}, nil)
}

// friends get user's friends list: follow eachother.
func friends(c *bm.Context) {
	var (
		err        error
		mid, _     = c.Get("mid")
		followings []*model.Following
		params     = c.Request.Form
		version    uint64
		versionStr = params.Get("re_version")
		crc32v     uint32
	)
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if followings, crc32v, err = relationSvc.Friends(c, mid.(int64), version); err != nil {
		log.Error("relationSvc.Followings(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       followings,
	}, nil)
}

// Blacks get user's black list.
func blacks(c *bm.Context) {
	var (
		err        error
		blacks     []*model.Following
		mid, _     = c.Get("mid")
		params     = c.Request.Form
		version    uint64
		pn, ps     int64
		total      int
		pnStr      = params.Get("pn")
		psStr      = params.Get("ps")
		versionStr = params.Get("re_version")
		crc32v     uint32
	)
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn <= 0 {
		pn = 1
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ps <= 0 || ps > _defaultPS {
		ps = _defaultPS
	}
	if blacks, crc32v, total, err = relationSvc.Blacks(c, mid.(int64), version, pn, ps); err != nil {
		log.Error("relationSvc.Blacks(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       blacks,
		"total":      total,
	}, nil)
}

// whispers get user's whisper list.
func whispers(c *bm.Context) {
	var (
		err        error
		pn, ps     int64
		version    uint64
		crc32v     uint32
		whispers   []*model.Following
		mid, _     = c.Get("mid")
		params     = c.Request.Form
		psStr      = params.Get("ps")
		pnStr      = params.Get("pn")
		versionStr = params.Get("re_version")
	)
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn <= 0 {
		pn = 1
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ps <= 0 || ps > _defaultPS {
		ps = _defaultPS
	}
	if whispers, crc32v, err = relationSvc.Whispers(c, mid.(int64), pn, ps, version); err != nil {
		log.Error("relationSvc.Whispers(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       whispers,
	}, nil)
}

// stat get user's follower list.
func stat(c *bm.Context) {
	var (
		err       error
		mid, vmid int64
		self      bool
		st        *mrl.Stat
		params    = c.Request.Form
		vmidStr   = params.Get("vmid")
	)
	midS, ok := c.Get("mid")
	if ok {
		mid = midS.(int64)
	} else {
		mid = 0
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	self = mid == vmid
	if st, err = relationSvc.Stat(c, vmid, self); err != nil {
		log.Error("relationSvc.Followers(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(st, nil)
}

// stat get user's follower list.
func stats(c *bm.Context) {
	var (
		err     error
		params  = c.Request.Form
		midsStr = params.Get("mids")
	)
	mids, err := xstr.SplitInts(midsStr)
	if err != nil || len(mids) > 20 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sts, err := relationSvc.Stats(c, mids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(sts, nil)
}

// tag 单个标签
func tag(c *bm.Context) {
	var (
		err      error
		pn, ps   int64
		mid, _   = c.Get("mid")
		params   = c.Request.Form
		tagIDStr = params.Get("tagid")
		tagid    int64
		psStr    = params.Get("ps")
		pnStr    = params.Get("pn")
		ti       []*model.Tag
	)
	if tagIDStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagid, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			log.Error("pn parse")
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn <= 0 {
		pn = 1
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ps <= 0 || ps > _defaultPS {
		ps = _defaultPS
	}
	if ti, err = relationSvc.Tag(c, mid.(int64), tagid, pn, ps); err != nil {
		log.Error("relationSvc.Tag(%d).tag(%d) error(%v)", mid, tagid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(ti, nil)
}

// tags 列表：标签-计数
func tags(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		tc     []*mrl.TagCount
	)
	if tc, err = relationSvc.Tags(c, mid.(int64)); err != nil {
		log.Error("relationSvc.Tags(%d). error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(tc, nil)
}

// mobileTags 移动端 列表：标签-计数
func mobileTags(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		tc     map[string][]*mrl.TagCount
	)
	if tc, err = relationSvc.MobileTags(c, mid.(int64)); err != nil {
		log.Error("relationSvc.Tags(%d). error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(tc, nil)
}

// tagUser 用户-fid 标签列表
func tagUser(c *bm.Context) {
	var (
		err    error
		fid    int64
		mid, _ = c.Get("mid")
		params = c.Request.Form
		fidStr = params.Get("fid")
		tc     map[int64]string
	)
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tc, err = relationSvc.UserTag(c, mid.(int64), fid); err != nil {
		log.Error("relationSvc.UserTag(%d).fid(%d) error(%v)", mid, fid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(tc, nil)
}

// tagCreate create tag.
func tagCreate(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		params = c.Request.Form
		tagStr = params.Get("tag")
		cres   int64
	)
	if tagStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cres, err = relationSvc.CreateTag(c, mid.(int64), tagStr); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"tagid": cres,
	}, nil)
}

// tagUpdate update tag.
func tagUpdate(c *bm.Context) {
	var (
		err      error
		mid, _   = c.Get("mid")
		params   = c.Request.Form
		tagIDStr = params.Get("tagid")
		tagID    int64
		newStr   = params.Get("name")
	)
	if tagIDStr == "" || newStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil || tagID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.UpdateTag(c, mid.(int64), tagID, newStr))
}

// tagDel del tag.
func tagDel(c *bm.Context) {
	var (
		err      error
		mid, _   = c.Get("mid")
		params   = c.Request.Form
		tagIDStr = params.Get("tagid")
		tagID    int64
	)
	if tagIDStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil || tagID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.DelTag(c, mid.(int64), tagID))
}

// tagsAddUsers tags add users.
func tagsAddUsers(c *bm.Context) {
	var (
		mid, _    = c.Get("mid")
		params    = c.Request.Form
		tagidsStr = params.Get("tagids")
		fidsStr   = params.Get("fids")
	)
	if tagidsStr == "" || fidsStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.TagsAddUsers(c, mid.(int64), tagidsStr, fidsStr))
}

// tagsCopyUsers tags copy users.
func tagsCopyUsers(c *bm.Context) {
	var (
		mid, _    = c.Get("mid")
		params    = c.Request.Form
		tagidsStr = params.Get("tagids")
		fidsStr   = params.Get("fids")
	)
	if tagidsStr == "" || fidsStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.TagsCopyUsers(c, mid.(int64), tagidsStr, fidsStr))
}

// tagsMoveUsers tags move users.
func tagsMoveUsers(c *bm.Context) {
	var (
		mid, _          = c.Get("mid")
		params          = c.Request.Form
		beforeTagIdsStr = params.Get("beforeTagids")
		afterTagIdsStr  = params.Get("afterTagids")
		fidsStr         = params.Get("fids")
	)
	if beforeTagIdsStr == "" || afterTagIdsStr == "" || fidsStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bid, err := strconv.ParseInt(beforeTagIdsStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.TagsMoveUsers(c, mid.(int64), bid, afterTagIdsStr, fidsStr))
}

func prompt(c *bm.Context) {
	mid, _ := c.Get("mid")
	arg := new(mrl.ArgPrompt)
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Mid = mid.(int64)
	b, err := relationSvc.Prompt(c, arg)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"prompt": b,
	}, nil)
}

func closePrompt(c *bm.Context) {
	mid, _ := c.Get("mid")
	arg := new(mrl.ArgPrompt)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(nil, relationSvc.ClosePrompt(c, arg))
}

func addSpecial(c *bm.Context) {
	mid, _ := c.Get("mid")
	arg := new(mrl.ArgFollowing)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(nil, relationSvc.AddSpecial(c, arg))
}
func delSpecial(c *bm.Context) {
	mid, _ := c.Get("mid")
	arg := new(mrl.ArgFollowing)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(nil, relationSvc.DelSpecial(c, arg))
}
func special(c *bm.Context) {
	mid, _ := c.Get("mid")
	c.JSON(relationSvc.Special(c, mid.(int64)))
}

// recommend get global recommend upper.
// deprecated
func recommend(c *bm.Context) {
	c.JSON([]interface{}{}, nil)
}

func recommendFollowlistEmpty(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &model.ArgRecommend{}
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Device = dev.(*bm.Device)
	arg.Mid = mid.(int64)
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)
	c.JSON(relationSvc.RecommendFollowlistEmpty(c, arg))
}

func recommendAnswerOK(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &model.ArgRecommend{}
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Device = dev.(*bm.Device)
	arg.Mid = mid.(int64)
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)
	c.JSON(relationSvc.RecommendAnswerOK(c, arg))
}

func recommendTagSuggest(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &model.ArgTagSuggestRecommend{}
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Device = dev.(*bm.Device)
	arg.Mid = mid.(int64)
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)
	c.JSON(relationSvc.RecommendTagSuggest(c, arg))
}

func recommendTagSuggestDetail(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &model.ArgTagSuggestRecommend{}
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.TagName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Device = dev.(*bm.Device)
	arg.Mid = mid.(int64)
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)
	c.JSON(relationSvc.RecommendTagSuggestDetail(c, arg))
}

// unread check unread status, for the 'show red point' function.
func unread(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		show   bool
	)
	disbaleAutoReset := c.Request.Form.Get("disableautoreset") == "1"
	if show, err = relationSvc.Unread(c, mid.(int64), disbaleAutoReset); err != nil {
		log.Error("relationSvc.Unread(%d) err(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"show": show,
	}, nil)
}

func unreadReset(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	c.JSON(nil, relationSvc.ResetUnread(c, mid.(int64)))
}

func unreadCount(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		count  int64
	)
	disbaleAutoReset := c.Request.Form.Get("disableautoreset") == "1"
	if count, err = relationSvc.UnreadCount(c, mid.(int64), disbaleAutoReset); err != nil {
		log.Error("relationSvc.UnreadCount(%d) err(%v)", mid, err)
		return
	}
	c.JSON(map[string]interface{}{
		"count": count,
	}, nil)
}

func unreadCountReset(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	c.JSON(nil, relationSvc.ResetUnreadCount(c, mid.(int64)))
}

func achieveGet(ctx *bm.Context) {
	arg := new(model.ArgAchieveGet)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	mid, _ := ctx.Get("mid")
	arg.Mid = mid.(int64)
	if arg.Award != "10k" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(relationSvc.AchieveGet(ctx, arg))
}

func achieve(ctx *bm.Context) {
	arg := new(model.ArgAchieve)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(relationSvc.Achieve(ctx, arg))
}

func followerNotifySetting(c *bm.Context) {
	var mid, _ = c.Get("mid")
	c.JSON(relationSvc.FollowerNotifySetting(c, mid.(int64)))
}

func enableFollowerNotify(c *bm.Context) {
	var mid, _ = c.Get("mid")
	c.JSON(nil, relationSvc.EnableFollowerNotify(c, mid.(int64)))
}

func disableFollowerNotify(c *bm.Context) {
	var mid, _ = c.Get("mid")
	c.JSON(nil, relationSvc.DisableFollowerNotify(c, mid.(int64)))
}

func sameFollowings(c *bm.Context) {
	arg := new(model.ArgSameFollowing)
	if err := c.Bind(arg); err != nil {
		return
	}
	mid, _ := c.Get("mid")
	arg.Mid = mid.(int64)
	if arg.Order != "asc" {
		arg.Order = "desc"
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 || arg.PS > _defaultPS {
		arg.PS = _defaultPS
	}

	followings, crc32v, total, err := relationSvc.SameFollowings(c, arg)
	if err != nil {
		log.Error("relationSvc.SameFollowings(%+v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"re_version": uint64(crc32v),
		"list":       followings,
		"total":      total,
	}, nil)
}
