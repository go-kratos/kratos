package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// userFolders return user's all folders
func userFolders(c *bm.Context) {
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	vmidStr := params.Get("vmid")
	oidStr := params.Get("oid")
	var (
		err  error
		mid  int64
		vmid int64
	)
	if midStr != "" {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if vmidStr != "" {
		vmid, err = strconv.ParseInt(vmidStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		if vmid <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else if mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var oid int64
	if oidStr != "" {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err := favSvc.UserFolders(c, int8(typ), mid, vmid, oid, int8(typ))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// folders return multi folders.
func folders(c *bm.Context) {
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	fvmidsStr := params.Get("fvmids")
	var (
		err    error
		mid    int64
		fid    int64
		vmid   int64
		fvmids []*model.ArgFVmid
	)
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if typeStr == "" || fvmidsStr == "" {
		log.Warn("params typeStr(%s) or fvmidsStr(%s) is empty", typeStr, fvmidsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fvmidsArr := strings.Split(fvmidsStr, ",")
	for _, fvmidStr := range fvmidsArr {
		fvmidArr := strings.Split(fvmidStr, "-")
		if len(fvmidArr) != 2 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		fid, err = strconv.ParseInt(fvmidArr[0], 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		vmid, err = strconv.ParseInt(fvmidArr[1], 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		fvmid := &model.ArgFVmid{
			Fid:  fid,
			Vmid: vmid,
		}
		fvmids = append(fvmids, fvmid)
	}
	data, err := favSvc.Folders(c, int8(tp), mid, fvmids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// folderInfo return folders & pages.
func folderInfo(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	vmidStr := params.Get("vmid")
	fidStr := params.Get("fid")
	var (
		err           error
		mid, uid, fid int64
	)
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
			log.Error("mid(%s) need a number > 0  error(%v)", midStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if vmidStr != "" {
		if uid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || uid <= 0 {
			log.Error("vmid(%s) need a number > 0  error(%v)", vmidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mid <= 0 && uid <= 0 {
		log.Warn("method (mid=0 && vmid(%s)) is empty", vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeStr == "" || fidStr == "" {
		log.Warn("params typeStr(%s) or fidStr(%s) is empty", typeStr, fidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err = strconv.ParseInt(fidStr, 10, 64)
	if err != nil || fid == 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := favSvc.Folder(c, int8(tp), mid, uid, fid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// folders add a new folder.
func addFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	name := params.Get("name")
	description := params.Get("description")
	cover := params.Get("cover")
	publicStr := params.Get("public")
	if typeStr == "" || name == "" {
		log.Error("params typeStr(%s) or name(%s) is empty", typeStr, name)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var public int64
	if publicStr != "" {
		public, err = strconv.ParseInt(publicStr, 10, 8)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", publicStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if public > 1 {
			public = 0
		}
	}
	fid, err := favSvc.AddFolder(c, int8(tp), mid.(int64), name, description, cover, int32(public), c.Request.Header.Get("Cookie"), params.Get("access_key"))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]int64{
		"fid": fid,
	}
	c.JSON(data, nil)
}

// updateFolder update folder info.
func updateFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	name := params.Get("name")
	description := params.Get("description")
	cover := params.Get("cover")
	publicStr := params.Get("public")
	fidStr := params.Get("fid")
	if typeStr == "" || fidStr == "" || name == "" {
		log.Error("params typeStr(%s) or fidStr(%s) or name(%s) is empty", typeStr, fidStr, name)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var public int64
	if publicStr != "" {
		public, err = strconv.ParseInt(publicStr, 10, 8)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", publicStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	err = favSvc.UpdateFolder(c, int8(tp), fid, mid.(int64), name, description, cover, int32(public), nil, nil)
	c.JSON(nil, err)
}

// delFolder del one folder.
func delFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
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
	err = favSvc.DelFolder(c, int8(tp), mid.(int64), fid)
	c.JSON(nil, err)
}

// cntUserFolders del one folder.
func cntUserFolders(c *bm.Context) {
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	vmidStr := params.Get("vmid")
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
			log.Error("mid(%s) need a number > 0  error(%v)", midStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	var uid int64
	if vmidStr != "" {
		if uid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || uid <= 0 {
			log.Error("vmid(%s) need a number > 0  error(%v)", vmidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mid <= 0 && uid <= 0 {
		log.Warn("method (mid=0 && vmid(%s)) is empty", vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	count, err := favSvc.CntUserFolders(c, int8(tp), mid, uid)
	if err != nil {
		log.Error("favSvr.IsFaved() err(%v)", err)
		return
	}
	data := map[string]interface{}{"count": count}
	c.JSON(data, nil)
}

// Favorites return all objects in the fid folder.
func Favorites(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	vmidStr := params.Get("vmid")
	fidStr := params.Get("fid")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	order := params.Get("order")
	tvStr := params.Get("tv")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	var (
		err error
		mid int64
		uid int64
		fid int64
		tid int
		tv  int
	)
	if midStr != "" {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
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
	} else if mid <= 0 {
		log.Warn("mid(%d) && vmidStr(%s) is empty", mid, vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tvStr != "" {
		if tv, err = strconv.Atoi(tvStr); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", tvStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tidStr != "" {
		if tid, err = strconv.Atoi(tidStr); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", tidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > conf.Conf.Fav.MaxPagesize || ps <= 0 {
		ps = conf.Conf.Fav.MaxPagesize
	}
	// fav objects
	data, err := favSvc.Favorites(c, int8(tp), mid, uid, fid, tid, tv, pn, ps, keyword, order)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// tlist return all video's type list.
func tlists(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	vmidStr := params.Get("vmid")
	fidStr := params.Get("fid")
	var (
		err error
		mid int64
		uid int64
		fid int64
	)
	if midStr != "" {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
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
	} else if mid <= 0 {
		log.Warn("mid(%d) && vmidStr(%s) is empty", mid, vmidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.Atoi(typeStr)
	if err != nil {
		log.Error("strconv.Aoti(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err := favSvc.Tlists(c, int8(typ), mid, uid, fid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// recentFavs return user's recent favs.
func recentFavs(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	sizeStr := params.Get("size")
	if typeStr == "" || midStr == "" {
		log.Warn("params typeStr(%s) or midStr(%s) is empty", typeStr, midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, err)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", sizeStr, err)
		c.JSON(nil, err)
		return
	}
	data, err := favSvc.RecentFavs(c, int8(typ), mid, size)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// batchFavs return all objects in the fid folder.
func batchFavs(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	if typeStr == "" || midStr == "" {
		log.Warn("params typeStr(%s) or midStr(%s) is empty", typeStr, midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, err)
		return
	}
	data, err := favSvc.BatchFavs(c, int8(typ), mid, conf.Conf.Fav.MaxBatchSize)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// addFav add a object into folder.
func addFav(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	fidsStr := params.Get("fid")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidStr == "" {
		log.Warn("params oid(%s) is empty", oidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var fids []int64
	if fidsStr != "" {
		fids, err = xstr.SplitInts(fidsStr)
		if err != nil {
			log.Warn("xstr.PlitInts(fids:%s) err(%v)", fidsStr, err)
		}
	}
	if len(fids) == 0 {
		fids = []int64{0}
	}
	for _, fid := range fids {
		err = favSvc.AddFav(c, int8(tp), mid.(int64), fid, oid, int8(tp))
	}
	c.JSON(nil, err)
}

// delFav delete a object from folder.
func delFav(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	fidsStr := params.Get("fid")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidStr == "" {
		log.Warn("method oid(%s) is empty", oidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fids, err := xstr.SplitInts(fidsStr)
	if err != nil {
		log.Warn("xstr.SplitInts(fidsStr:%v) err(%v)", fidsStr, err)
	}
	if len(fids) == 0 {
		fids = []int64{0}
	}
	for _, fid := range fids {
		err = favSvc.DelFav(c, int8(tp), mid.(int64), fid, oid, int8(tp))
	}
	c.JSON(nil, err)
}

// multiAddFav add multiple object from folder.
func multiAddFavs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	fidStr := params.Get("fid")
	oidsStr := params.Get("oids")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidsStr == "" {
		log.Warn("method oid(%s) is empty", oidsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oidsArr := strings.Split(oidsStr, ",")
	if len(oidsArr) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	var oid int64
	oids := make([]int64, len(oidsArr))
	for i, oidStr := range oidsArr {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		oids[i] = oid
	}
	err = favSvc.MultiAddFavs(c, int8(tp), mid.(int64), fid, oids)
	c.JSON(nil, err)
}

// multiDelFav delete multiple object from folder.
func multiDelFavs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	fidStr := params.Get("fid")
	oidsStr := params.Get("oids")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidsStr == "" {
		log.Warn("method oid(%s) is empty", oidsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oidsArr := strings.Split(oidsStr, ",")
	if len(oidsArr) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	oids := make([]int64, len(oidsArr))
	var oid int64
	for i, oidStr := range oidsArr {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		oids[i] = oid
	}
	if err != nil {
		log.Warn("xstr.SplitInts(fidsStr:%v) err(%v)", oidsStr, err)
	}
	err = favSvc.MultiDelFavs(c, int8(tp), mid.(int64), fid, oids)
	c.JSON(nil, err)
}

// isFavored detemine object whether or not favored by mid.
func isFavored(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidStr == "" {
		log.Warn("method oid(%s) is empty", oidStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	favored, err := favSvc.IsFavored(c, int8(tp), mid.(int64), oid)
	if err != nil {
		log.Error("favSvr.IsFaved() err(%v)", err)
		return
	}
	data := map[string]interface{}{"favored": favored}
	c.JSON(data, nil)
}

// isFavoreds detemine objects whether or not favored by mid.
func isFavoreds(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	typeStr := params.Get("type")
	oidsStr := params.Get("oids")
	if typeStr == "" {
		log.Warn("params typeStr(%s) is empty", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidsStr == "" {
		log.Warn("method oid(%s) is empty", oidsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", oidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	favoreds, err := favSvc.IsFavoreds(c, int8(typ), mid.(int64), oids)
	if err != nil {
		log.Error("favSvr.IsFaved() err(%v)", err)
		return
	}
	c.JSON(favoreds, nil)
}

// sortFolders sort all user's folders
func sortFolders(c *bm.Context) {
	var (
		fids []int64
		err  error
	)
	params := c.Request.Form
	fidStr := params.Get("fids")
	typeStr := params.Get("type")
	mid, _ := c.Get("mid")
	if fidStr == "" || typeStr == "" {
		log.Error("arg fids or type is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fids, err = xstr.SplitInts(fidStr)
	if err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.SetFolderSort(c, int8(typ), mid.(int64), fids)
	c.JSON(nil, err)
}

// renameFolder rename folder.
func renameFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	fidStr := params.Get("fid")
	typeStr := params.Get("type")
	name := params.Get("name")
	if fidStr == "" || typeStr == "" || name == "" {
		log.Warn("arg fid or type is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(name)) > conf.Conf.Fav.MaxNameLen {
		log.Warn("arg name(%s) is empty or it's length more than %d", name, conf.Conf.Fav.MaxNameLen)
		c.JSON(nil, ecode.FavNameTooLong)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.UpFolderName(c, int8(typ), mid.(int64), fid, name, c.Request.Header.Get("Cookie"), params.Get("access_key"))
	c.JSON(nil, err)
}

// upAttrFolder update folder's attr.
func upAttrFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	fidStr := params.Get("fid")
	pubStr := params.Get("public")
	if fidStr == "" || pubStr == "" || typeStr == "" {
		log.Warn("args fid(%s) or public(%s) or type(%s) is empty", fidStr, pubStr, typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	public, err := strconv.Atoi(pubStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", pubStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = favSvc.UpFolderAttr(c, int8(typ), mid.(int64), fid, int32(public))
	c.JSON(nil, err)
}

// moveFavs move some video into other folder.
func moveFavs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	ofidStr := params.Get("old_fid")
	nfidStr := params.Get("new_fid")
	oidsStr := params.Get("oids")
	if oidsStr == "" || ofidStr == "" || nfidStr == "" || typeStr == "" {
		log.Warn("args oids(%s) old_fid(%s) new_fid(%s) type(%s) is empty", oidsStr, ofidStr, nfidStr, typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ofidStr == nfidStr {
		log.Warn("move favs to the same folder...")
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
	typ, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oidStrs := strings.Split(oidsStr, ",")
	if len(oidStrs) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(oidStrs) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	oids := make([]int64, len(oidStrs))
	var oid int64
	for i, oidStr := range oidStrs {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		oids[i] = oid
	}
	err = favSvc.MoveFavs(c, int8(typ), mid.(int64), ofid, nfid, oids)
	c.JSON(nil, err)
}

// copyFavs copy resources into other folder.
func copyFavs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	omidStr := params.Get("old_mid")
	ofidStr := params.Get("old_fid")
	nfidStr := params.Get("new_fid")
	oidsStr := params.Get("oids")
	if oidsStr == "" || ofidStr == "" || nfidStr == "" || typeStr == "" {
		log.Warn("args oids(%s) old_fid(%s) new_fid(%s) type(%s) is empty", oidsStr, ofidStr, nfidStr, typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ofidStr == nfidStr {
		log.Warn("move favs to the same folder...")
		c.JSON(nil, ecode.FavFolderSame)
		return
	}
	var (
		err  error
		omid int64
	)
	if omidStr != "" {
		omid, err = strconv.ParseInt(omidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", omidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
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
	typ, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oidStrs := strings.Split(oidsStr, ",")
	if len(oidStrs) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(oidStrs) > conf.Conf.Fav.MaxOperationNum {
		c.JSON(nil, ecode.FavMaxOperNum)
		return
	}
	oids := make([]int64, len(oidStrs))
	var oid int64
	for i, oidStr := range oidStrs {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		oids[i] = oid
	}
	err = favSvc.CopyFavs(c, int8(typ), omid, mid.(int64), ofid, nfid, oids)
	c.JSON(nil, err)
}

// inDefaultFolder detemine resource whether or not favored in default folder.
func inDefaultFolder(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	if oidStr == "" || typeStr == "" {
		log.Warn("oid(%s) or type(%s) is empty", oidStr, typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var isIn bool
	isIn, err = favSvc.InDefaultFolder(c, int8(typ), mid.(int64), oid)
	data := map[string]bool{"default": isIn}
	c.JSON(data, err)
}

// userList return all objects in the fid folder.
func userList(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	var (
		err error
		oid int64
	)
	if typeStr == "" || oidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > conf.Conf.Fav.MaxBatchSize || ps <= 0 {
		ps = conf.Conf.Fav.MaxBatchSize
	}
	// fav users
	data, err := favSvc.UserList(c, int8(typ), oid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// oidCount return oid's fav stats.
func oidCount(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	var (
		err error
		oid int64
	)
	if typeStr == "" || oidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// fav users
	data, err := favSvc.OidCount(c, int8(typ), oid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// oidsCount return oid's fav stats.
func oidsCount(c *bm.Context) {
	// params
	params := c.Request.Form
	typeStr := params.Get("type")
	oidsStr := params.Get("oids")
	var err error
	if typeStr == "" || oidsStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Error("xstr.SplitInts(oidsStr:%v) err(%v)", oidsStr, err)
		return
	}
	data, err := favSvc.OidsCount(c, int8(typ), oids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// isCleaned check the clean action's cool down time and access
func isCleaned(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.URL.Query()
	typeStr := params.Get("type")
	fidStr := params.Get("fid")
	if fidStr == "" || typeStr == "" {
		log.Error("arg fids or type is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cleanState, err := favSvc.CleanState(c, int8(typ), mid.(int64), fid)
	if err != nil {
		log.Error("favSvc.IsCleaned(%d,%d) error(%v)", mid, fid, err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{"state": cleanState}
	c.JSON(data, nil)
}

func cleanInvalidFavs(c *bm.Context) {
	mid, _ := c.Get("mid")
	params := c.Request.Form
	typeStr := params.Get("type")
	fidStr := params.Get("fid")
	if fidStr == "" || typeStr == "" {
		log.Error("arg fids or type is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, err := strconv.ParseInt(typeStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = favSvc.CleanInvalidArcs(c, int8(typ), mid.(int64), fid); err != nil {
		log.Error("favSvc.CleanInvalidArcs(%d,%d,%d) error(%v)", typ, mid, fid, err)
	}
	c.JSON(nil, err)
}
