package http

import (
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/model"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Infos Infos.
type Infos struct {
	Tids   []int64  `json:"tids"`
	TNames []string `json:"tnames"`
}

// check tag name for up
func checkTagName(c *bm.Context) {
	var (
		err    error
		tag    *model.Tag
		params = c.Request.Form
	)
	name := params.Get("tag_name")
	mid, _ := c.Get("mid")
	if name, err = svr.CheckName(name); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tag, err = svr.CheckTag(c, mid.(int64), name, time.Now()); err != nil {
		log.Error("tagSvr.InfoByName(%s) error(%v)", name, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(tag, nil)
}

func info(c *bm.Context) {
	var (
		err   error
		mid   int64
		tag   *model.Tag
		param = new(struct {
			ID   int64  `form:"tag_id"`
			Name string `form:"tag_name"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.ID <= 0 && param.Name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if param.ID > 0 {
		if tag, err = svr.InfoByID(c, mid, param.ID); err != nil {
			c.JSON(nil, err)
			return
		}
		c.JSON(tag, nil)
		return
	}
	if len([]rune(param.Name)) > 30 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tag, err = svr.InfoByName(c, mid, param.Name); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tag, nil)
}

func infos(c *bm.Context) {
	var (
		err    error
		mid    int64
		data   []*model.Tag
		params = c.Request.Form
	)
	v := new(Infos)
	vStr := params.Get("data")
	if len(vStr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(vStr), &v); err != nil {
		log.Error("json.Unmarshal vStr failed, error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if len(v.Tids) == 0 && len(v.TNames) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.Tids) > 0 {
		if len(v.Tids) > model.MaxTagNum {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if data, err = svr.MinfoByIDs(c, mid, v.Tids); err != nil {
			log.Error("tagSvr.MinfoByIDs(%v) error(%v)", v.Tids, err)
			c.JSON(nil, err)
			return
		}
	} else {
		if len(v.TNames) > model.MaxTagNum {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		for _, name := range v.TNames {
			if len([]rune(name)) > model.TnameMaxLen || name == "" {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		if data, err = svr.MinfoByNames(c, mid, v.TNames); err != nil {
			log.Error("tagSvr.MinfoByNames(%v) error(%v)", v.TNames, err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(data, nil)
}

func mInfo(c *bm.Context) {
	var (
		err   error
		mid   int64
		data  []*model.Tag
		param = new(struct {
			IDs   []int64  `form:"tag_id,split"`
			Names []string `form:"tag_name,split"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if len(param.IDs) == 0 && len(param.Names) == 0 {
		log.Error("info() name == nil or tid == nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if len(param.IDs) != 0 {
		if len(param.IDs) > conf.Conf.Tag.MaxSelTagNum {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if data, err = svr.MinfoByIDs(c, mid, param.IDs); err != nil {
			c.JSON(nil, err)
			return
		}
		c.JSON(data, nil)
		return
	}
	if len(param.Names) > conf.Conf.Tag.MaxSelTagNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, name := range param.Names {
		if len([]rune(name)) > model.TnameMaxLen || name == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if data, err = svr.MinfoByNames(c, mid, param.Names); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func detail(c *bm.Context) {
	var (
		err      error
		tid, mid int64
		pn, ps   int
		params   = c.Request.Form
	)
	tidStr := params.Get("tag_id")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Tag.MaxArcsPageSize {
		ps = conf.Conf.Tag.MaxArcsPageSize
	}
	detail, err := svr.Detail(c, tid, mid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(detail, nil)
}

func hotTags(c *bm.Context) {
	var (
		err     error
		rid     int64
		mid     int64
		hotType int64
		data    []*model.HotTags
		params  = c.Request.Form
	)
	ridStr := params.Get("rid")
	hotTypeStr := params.Get("type")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if hotTypeStr != "" {
		if hotType, err = strconv.ParseInt(hotTypeStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if hotType < 0 || hotType > 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = svr.HotTags(c, mid, rid, hotType); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(data) == 0 {
		data = []*model.HotTags{}
	}
	c.JSON(data, nil)
}

func similarTags(c *bm.Context) {
	var (
		err    error
		rid    int64
		tid    int64
		data   []*model.SimilarTag
		params = c.Request.Form
	)
	ridStr := params.Get("rid")
	tidStr := params.Get("tid")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = svr.SimilarTags(c, rid, tid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func changeSim(c *bm.Context) {
	var (
		err    error
		tid    int64
		data   []*model.SimilarTag
		params = c.Request.Form
	)
	tidStr := params.Get("tag_id")
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = svr.ChangeSim(c, tid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func addActivityTag(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
	)
	name := params.Get("tag_name")
	if name, err = svr.CheckName(name); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.AddActivityTag(c, name, time.Now()); err != nil {
		log.Error("tagSvr.AddActivityTag(%s) error(%v)", name, err)
	}
	c.JSON(nil, err)
}

func recommandTag(c *bm.Context) {
	var (
		err  error
		data map[int64]map[string][]*rpcModel.UploadTag
	)
	if data, err = svr.RecommandTag(c); err != nil {
		log.Error("svr.RecommandTag() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func synonymTag(c *bm.Context) {
	data, err := svr.TagGroup(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func hotMap(c *bm.Context) {
	data, err := svr.HotMap(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func prids(c *bm.Context) {
	data, err := svr.Prids(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func tagBatchInfo(c *bm.Context) {
	var (
		err    error
		mid    int64
		tnames []string
		data   []*model.Tag
		params = c.Request.Form
	)
	v := new(Infos)
	vStr := params.Get("tnames")
	if len(vStr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(vStr), &v); err != nil {
		log.Error("json.Unmarshal vStr failed, error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if len(v.TNames) > model.MaxTagNum || len(v.TNames) <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, name := range v.TNames {
		var tname string
		if tname, err = svr.CheckName(name); err == nil {
			tnames = append(tnames, tname)
		}
	}
	if data, err = svr.TopicTags(c, mid, tnames, time.Now()); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
