package http

import (
	"strconv"
	"strings"

	"go-common/app/interface/main/favorite/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// videoFolders get all favorite folders
func videoFolders(c *bm.Context) {
	var (
		uid       int64
		mid       int64
		vmid      int64
		aid       int64
		err       error
		isSelf    bool
		mediaList bool
		fromWeb   bool
	)
	req := c.Request
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	if req.Form.Get("medialist") == "1" {
		mediaList = true
	}
	params := req.URL.Query()
	app := params.Get("mobi_app")
	build, _ := strconv.ParseInt(params.Get("build"), 10, 64)
	device := params.Get("device")
	if (app == "android" && build >= 5360001 && build <= 5361000) || (app == "iphone" && build == 8300 && device == "phone") {
		mediaList = true
	}
	if app == "" {
		fromWeb = true
	}
	vmidStr := req.Form.Get("vmid")
	if vmidStr != "" {
		uid, err = strconv.ParseInt(vmidStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		if uid <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		isSelf = mid == uid
		vmid = uid
	} else if mid != 0 {
		uid = mid
		isSelf = true
	} else {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := req.Form.Get("aid")
	if aidStr != "" {
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(aid:%s)", aidStr)
			return
		}
	}
	data, err := favSvc.FavFolders(c, mid, vmid, uid, aid, isSelf, mediaList, fromWeb)
	c.JSON(data, err)
}

// addVideoFolder add a folder.
func addVideoFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	name := params.Get("name")
	pubStr := params.Get("public")
	if name == "" || len([]rune(name)) > conf.Conf.Fav.MaxNameLen {
		log.Warn("arg name(%s) is empty or it's length more than %d", name, conf.Conf.Fav.MaxNameLen)
		c.JSON(nil, ecode.FavNameTooLong)
		return
	}
	var (
		pub int64
		err error
	)
	if pubStr != "" {
		if pub, err = strconv.ParseInt(pubStr, 10, 64); err != nil || pub < 0 || pub > 1 {
			pub = 0
		}
	}
	var fid int64
	if fid, err = favSvc.AddFavFolder(c, mid.(int64), name, c.Request.Header.Get("Cookie"), params.Get("access_key"), int32(pub)); err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]int64{
		"fid": fid,
	}
	c.JSON(data, nil)
}

// renameVideoFolder rename folder.
func renameVideoFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	name := params.Get("name")
	if fidStr == "" {
		log.Warn("arg fid is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if name == "" || len([]rune(name)) > conf.Conf.Fav.MaxNameLen {
		log.Warn("arg name(%s) is empty or it's length more than %d", name, conf.Conf.Fav.MaxNameLen)
		c.JSON(nil, ecode.FavNameTooLong)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.UpFavName(c, mid.(int64), fid, name, c.Request.Header.Get("Cookie"), params.Get("access_key"))
	c.JSON(nil, err)
}

// upStateVideoFolder update folder's state.
func upStateVideoFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	pubStr := params.Get("public")
	if fidStr == "" || pubStr == "" {
		log.Warn("method fid(%s) public(%s) is empty", fidStr, pubStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	public, err := strconv.Atoi(pubStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", pubStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.UpFavState(c, mid.(int64), fid, int32(public), c.Request.Header.Get("Cookie"), params.Get("access_key"))
	c.JSON(nil, err)
}

// delVideoFolder delete folder.
func delVideoFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	if fidStr == "" {
		log.Warn("method fid(%s) is empty", fidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.DelVideoFolder(c, mid.(int64), fid)
	c.JSON(nil, err)
}

// sortFavFolders sort all favorite folders
func sortVideoFolders(c *bm.Context) {
	var (
		fids []int64
		err  error
	)
	params := c.Request.Form
	fidStr := params.Get("fids")
	mid, _ := c.Get("mid")
	if fidStr == "" {
		log.Error("arg fids is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fids, err = xstr.SplitInts(fidStr)
	if err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.SetVideoFolderSort(c, mid.(int64), fids)
	c.JSON(nil, err)
}

// favVideo return all videos in the fid folder.
func favVideo(c *bm.Context) {
	// params
	params := c.Request.Form
	vmidStr := params.Get("vmid")
	fidStr := params.Get("fid")
	tidStr := params.Get("tid")
	keywordStr := params.Get("keyword")
	orderStr := params.Get("order")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	var (
		err  error
		mid  int64
		vmid int64
		uid  int64
	)
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	if vmidStr != "" {
		if uid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || uid <= 0 {
			log.Error("vmid(%s) need a number > 0  error(%v)", vmidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		vmid = uid
	} else if mid != 0 {
		uid = mid
	} else {
		log.Warn("mid(%d) && vmidStr(%s)) is empty", mid, vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, _ := strconv.ParseInt(fidStr, 10, 64)
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		tid = 0
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > conf.Conf.Fav.MaxPagesize || ps <= 0 {
		ps = conf.Conf.Fav.MaxPagesize
	}
	// fav video
	data, err := favSvc.FavVideo(c, mid, vmid, uid, fid, keywordStr, orderStr, tid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// tidList return all tids in the fid folder.
func tidList(c *bm.Context) {
	// params
	params := c.Request.Form
	vmidStr := params.Get("vmid")
	fidStr := params.Get("fid")
	var (
		err  error
		mid  int64
		vmid int64
		uid  int64
	)
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	if vmidStr != "" {
		if uid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || uid <= 0 {
			log.Error("vmid(%s) need a number > 0  error(%v)", vmidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		vmid = uid
	} else if mid != 0 {
		uid = mid
	} else {
		log.Warn("mid(%d) && vmidStr(%s)) is empty", mid, vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, _ := strconv.ParseInt(fidStr, 10, 64)
	// fav video
	data, err := favSvc.TidList(c, mid, vmid, uid, fid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// favVideoNewest return the newest videos in the all folder.
func favVideoNewest(c *bm.Context) {
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	params := c.Request.URL.Query()
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > conf.Conf.Fav.MaxPagesize || ps <= 0 {
		ps = conf.Conf.Fav.MaxPagesize
	}
	data, err := favSvc.RecentArcs(c, mid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// addFavVideo add a video into folder.
func addFavVideo(c *bm.Context) {
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	params := c.Request.Form
	fidsStr := params.Get("fid")
	aidStr := params.Get("aid")
	if aidStr == "" {
		log.Warn("params aid(%s) is empty", aidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fids, err := xstr.SplitInts(fidsStr)
	if err != nil {
		log.Warn("xstr.PlitInts(fids:%s) err(%v)", fidsStr, err)
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(fids) == 0 {
		if err = favSvc.AddArc(c, mid, 0, aid, c.Request.Header.Get("Cookie"), params.Get("access_key")); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	for _, fid := range fids {
		if err = favSvc.AddArc(c, mid, fid, aid, c.Request.Header.Get("Cookie"), params.Get("access_key")); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if collector != nil {
		collector.InfoAntiCheat2(c, "", aidStr, strconv.FormatInt(mid, 10), fidsStr, infoc.ItemTypeAv, infoc.ActionFav, "")
	}
	c.JSON(nil, err)
}

// delFavVideo delete a video from folder.
func delFavVideo(c *bm.Context) {
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	params := c.Request.Form
	fidsStr := params.Get("fid")
	aidStr := params.Get("aid")
	if aidStr == "" {
		log.Warn("method aid(%s) is empty", aidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fids, err := xstr.SplitInts(fidsStr)
	if err != nil {
		log.Warn("xstr.SplitInts(fidsStr:%v) err(%v)", fidsStr, err)
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(fids) == 0 {
		err = favSvc.DelArc(c, mid, 0, aid)
		c.JSON(nil, err)
		return
	}
	for _, fid := range fids {
		err = favSvc.DelArc(c, mid, fid, aid)
	}
	c.JSON(nil, err)
}

// moveFavVideos move some video into other folder.
func moveFavVideos(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	ofidStr := params.Get("old_fid")
	nfidStr := params.Get("new_fid")
	aidsStr := params.Get("aids")
	if aidsStr == "" || ofidStr == "" || nfidStr == "" {
		log.Warn("method aids(%s) old_fid(%s) new_fid(%s) is empty", aidsStr, ofidStr, nfidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ofidStr == nfidStr {
		log.Warn("move videos to the same folder...")
		c.JSON(nil, ecode.FavFolderSame)
		return
	}
	ofid, err := strconv.ParseInt(ofidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", ofidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	nfid, err := strconv.ParseInt(nfidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", nfidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidArr := strings.Split(aidsStr, ",")
	if len(aidArr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aidArr) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	aids := make([]int64, len(aidArr))
	var aid int64
	for i, aidStr := range aidArr {
		aid, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		aids[i] = aid
	}
	err = favSvc.MoveArcs(c, mid.(int64), ofid, nfid, aids)
	c.JSON(nil, err)
}

// copyFavVideos move some video into other folder.
func copyFavVideos(c *bm.Context) {
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	params := c.Request.Form
	omidStr := params.Get("old_mid")
	ofidStr := params.Get("old_fid")
	nfidStr := params.Get("new_fid")
	aidsStr := params.Get("aids")
	if aidsStr == "" || ofidStr == "" || nfidStr == "" {
		log.Warn("method aids(%s) old_fid(%s) new_mid(%s) is empty", aidsStr, ofidStr, nfidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ofidStr == nfidStr {
		log.Warn("copy videos to the same folder...")
		c.JSON(nil, ecode.FavFolderSame)
		return
	}
	omid, err := strconv.ParseInt(omidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", omidStr, err)
		omid = mid
	}
	ofid, err := strconv.ParseInt(ofidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", ofidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	nfid, err := strconv.ParseInt(nfidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", nfidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidArr := strings.Split(aidsStr, ",")
	if len(aidArr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aidArr) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	aids := make([]int64, len(aidArr))
	var aid int64
	for i, aidStr := range aidArr {
		aid, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		aids[i] = aid
	}
	err = favSvc.CopyArcs(c, mid, omid, ofid, nfid, aids)
	c.JSON(nil, err)
}

// delVideos delete some video from folder.
func delFavVideos(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	aidsStr := params.Get("aids")
	if aidsStr == "" {
		log.Warn("method aid(%s) is empty", aidsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", fidStr, err)
	}
	aidArr := strings.Split(aidsStr, ",")
	if len(aidArr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aidArr) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	aids := make([]int64, len(aidArr))
	var aid int64
	for i, aidStr := range aidArr {
		aid, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		aids[i] = aid
	}
	err = favSvc.DelArcs(c, mid.(int64), fid, aids)
	c.JSON(nil, err)
}

// isFavoured detemine video whether or not favoured by mid.
func isFavoured(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	faved, count, err := favSvc.IsFaved(c, mid.(int64), aid)
	if err != nil {
		log.Error("favSvr.IsFaved() err(%v)", err)
		return
	}
	data := map[string]interface{}{"favoured": faved, "count": count}
	c.JSON(data, nil)
}

// isFavoureds detemine video whether or not favoured by mid.
func isFavoureds(c *bm.Context) {
	var (
		aids []int64
		err  error
	)
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	aidStr := params.Get("aids")
	if aidStr == "" {
		log.Warn("method aid(%s) is empty", aidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aids, err = xstr.SplitInts(aidStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	faved, _ := favSvc.IsFaveds(c, mid.(int64), aids)
	c.JSON(faved, nil)
}

// inDefaultFav detemine video whether or not favoured in default folder.
func inDefaultFav(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	aidStr := params.Get("aid")
	if aidStr == "" {
		log.Warn("method aid(%s) is empty", aidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var isIn bool
	isIn, err = favSvc.InDef(c, mid.(int64), aid)
	data := map[string]bool{"default": isIn}
	c.JSON(data, err)
}

// isCleaned check the clean action's cool down time and access
func isCleaned(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	fidStr := params.Get("fid")
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cleanState, err := favSvc.CleanState(c, mid.(int64), fid)
	if err != nil {
		log.Error("favSvc.IsCleaned(%d,%d) error(%v)", mid, fid, err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{"state": cleanState}
	c.JSON(data, nil)
}

func cleanInvalidArcs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = favSvc.CleanInvalidArcs(c, mid.(int64), fid); err != nil {
		log.Error("favSvc.CleanInvalidArcs(%d,%d) error(%v)", mid, fid, err)
	}
	c.JSON(nil, err)
}
