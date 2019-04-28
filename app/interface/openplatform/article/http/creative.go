package http

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func lists(c *bm.Context) {
	var (
		mid   int64
		err   error
		novel bool
		list  []*model.CreativeList
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	if novel, list, err = artSrv.CreativeUpLists(c, mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"novel": novel,
		"lists": list,
		"total": len(list),
		"limit": conf.Conf.Article.ListLimit,
	}, nil)
}

func addList(c *bm.Context) {
	var (
		mid int64
		err error
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	req := new(struct {
		Name     string `form:"name" validate:"required"`
		Summary  string `form:"summary" validate:"min=0,max=233"`
		ImageURL string `form:"image_url"`
	})
	if err = c.Bind(req); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if req.ImageURL != "" && !model.CheckBFSImage(req.ImageURL) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.CreativeAddList(c, mid, req.Name, req.Summary, req.ImageURL))
}

func delList(c *bm.Context) {
	var (
		mid, id int64
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	id, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.CreativeDelList(c, mid, id))
}

func updateArticleList(c *bm.Context) {
	var (
		mid int64
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	listID, _ := strconv.ParseInt(c.Request.Form.Get("list_id"), 10, 64)
	if listID < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	articleID, _ := strconv.ParseInt(c.Request.Form.Get("article_id"), 10, 64)
	if articleID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.CreativeUpdateArticleList(c, mid, articleID, listID, true))
}

func listAllArticles(c *bm.Context) {
	var (
		mid, id int64
		err     error
		list    *model.List
		arts    []*model.ListArtMeta
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	id, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, arts, err = artSrv.CreativeListAllArticles(c, mid, id); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"list":     list,
		"articles": arts,
		"total":    len(arts),
		"limit":    conf.Conf.Article.ListArtsLimit,
	}, nil)
}

func updateListArticles(c *bm.Context) {
	var (
		mid  int64
		list *model.List
		err  error
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req := new(struct {
		ListID   int64   `form:"list_id" validate:"min=0"`
		Name     string  `form:"name" validate:"required"`
		Summary  string  `form:"summary" validate:"min=0,max=233"`
		ImageURL string  `form:"image_url"`
		OnlyList bool    `form:"only_list"`
		IDs      []int64 `form:"ids,split"`
	})
	if err = c.Bind(req); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if req.ImageURL != "" && !model.CheckBFSImage(req.ImageURL) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, err = artSrv.CreativeUpdateListArticles(c, req.ListID, req.Name, req.ImageURL, req.Summary, req.OnlyList, mid, req.IDs); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{"list": list}, nil)
}

func canAddArts(c *bm.Context) {
	var (
		mid  int64
		err  error
		arts []*model.ListArtMeta
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	if arts, err = artSrv.CreativeCanAddArticles(c, mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{"articles": arts}, nil)
}

func webSubArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	title := params.Get("title")
	content := params.Get("content")
	summary := params.Get("summary")
	bannerURL := params.Get("banner_url")
	tidStr := params.Get("tid")
	categoryStr := params.Get("category")
	reprintStr := params.Get("reprint")
	tags := params.Get("tags")
	imageURLs := params.Get("image_urls")
	wordsStr := params.Get("words")
	actIDStr := params.Get("act_id")
	scoreStr := params.Get("score")
	mediaIDStr := params.Get("media_id")
	spoilerStr := params.Get("spoiler")
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	// check params
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	var (
		spoiler          int64
		aid              int64
		actID            int64
		mediaID          int64
		score            int64
		err1, err2, err3 error
	)
	if aidStr != "" {
		id, err := strconv.ParseInt(aidStr, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		aid = id
	}
	if actIDStr != "" {
		actid, err := strconv.ParseInt(actIDStr, 10, 64)
		if err != nil || actid <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		actID = actid
	}
	listIDStr := params.Get("list_id")
	var listID int64
	if listIDStr != "" {
		lid, err := strconv.ParseInt(listIDStr, 10, 64)
		if err != nil || lid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		listID = lid
	}
	if mediaIDStr != "" {
		mediaID, err1 = strconv.ParseInt(mediaIDStr, 10, 64)
		if err1 != nil || mediaID < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mediaID > 0 {
		score, err2 = strconv.ParseInt(scoreStr, 10, 64)
		spoiler, err3 = strconv.ParseInt(spoilerStr, 10, 32)
		if err2 != nil || err3 != nil || spoiler < 0 || (score < 1 || score > 10 || score%2 != 0) {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if ok, err := artSrv.LevelRequired(c, mid); err != nil || !ok {
			c.JSON(nil, ecode.ArtLevelFailedErr)
			return
		}
		if id, err := artSrv.MediaArticle(c, mediaID, mid); err != nil || id > 0 {
			c.JSON(nil, ecode.ArtMediaExistedErr)
			return
		}
	}
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSrv.ParseParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// params
	art := &model.ArtParam{
		AID:             aid,
		MID:             mid,
		Title:           title,
		Content:         content,
		Summary:         summary,
		BannerURL:       bannerURL,
		Tags:            tags,
		ImageURLs:       artParam.ImageURLs,
		OriginImageURLs: artParam.OriginImageURLs,
		RealIP:          ip,
		Category:        artParam.Category,
		TemplateID:      artParam.TemplateID,
		Reprint:         artParam.Reprint,
		Words:           words,
		DynamicIntro:    dynamicIntrosStr,
		ActivityID:      actID,
		ListID:          listID,
		MediaID:         mediaID,
		Spoiler:         int32(spoiler),
	}
	// submit
	id, err := artSrv.CreativeSubArticle(c, mid, art, "", ck, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// 番剧评分
	if mediaID > 0 {
		artSrv.SetMediaScore(c, score, id, mediaID, mid)
	}
	c.JSON(map[string]int64{"aid": id}, nil)
}

func webUpdateArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	title := params.Get("title")
	content := params.Get("content")
	summary := params.Get("summary")
	bannerURL := params.Get("banner_url")
	tidStr := params.Get("tid")
	categoryStr := params.Get("category")
	reprintStr := params.Get("reprint")
	tags := params.Get("tags")
	imageURLs := params.Get("image_urls")
	wordsStr := params.Get("words")
	spoilerStr := params.Get("spoiler")
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	scoreStr := params.Get("score")
	mediaIDStr := params.Get("media_id")
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	// check params
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSrv.ParseParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	listIDStr := params.Get("list_id")
	var listID int64
	if listIDStr != "" {
		lid, err := strconv.ParseInt(listIDStr, 10, 64)
		if err != nil || lid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		listID = lid
	}
	var (
		spoiler          int64
		mediaID          int64
		score            int64
		err1, err2, err3 error
	)
	if mediaIDStr != "" {
		mediaID, err1 = strconv.ParseInt(mediaIDStr, 10, 64)
		if err1 != nil || mediaID < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mediaID > 0 {
		score, err2 = strconv.ParseInt(scoreStr, 10, 64)
		spoiler, err3 = strconv.ParseInt(spoilerStr, 10, 32)
		if err2 != nil || err3 != nil || spoiler < 0 || (score < 1 || score > 10 || score%2 != 0) {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if mediaid, err := artSrv.MediaIDByID(c, aid); err != nil || mediaid != mediaID {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	// params
	art := &model.ArtParam{
		AID:             aid,
		MID:             mid,
		Title:           title,
		Content:         content,
		Summary:         summary,
		BannerURL:       bannerURL,
		Tags:            tags,
		ImageURLs:       artParam.ImageURLs,
		OriginImageURLs: artParam.OriginImageURLs,
		RealIP:          ip,
		Category:        artParam.Category,
		TemplateID:      artParam.TemplateID,
		Reprint:         artParam.Reprint,
		Words:           words,
		DynamicIntro:    dynamicIntrosStr,
		ListID:          listID,
		Spoiler:         int32(spoiler),
		MediaID:         mediaID,
	}
	if err = artSrv.CreativeUpdateArticle(c, mid, art, "", ck, ip); err != nil {
		c.JSON(nil, err)
		return
	}
	// 番剧评分
	if mediaID > 0 {
		artSrv.SetMediaScore(c, score, aid, mediaID, mid)
	}
	c.JSON(nil, nil)
}

func webSubmitDraft(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	title := params.Get("title")
	content := params.Get("content")
	summary := params.Get("summary")
	bannerURL := params.Get("banner_url")
	tidStr := params.Get("tid")
	categoryStr := params.Get("category")
	reprintStr := params.Get("reprint")
	tags := params.Get("tags")
	imageURLs := params.Get("image_urls")
	wordsStr := params.Get("words")
	mediaIDStr := params.Get("media_id")
	spoilerStr := params.Get("spoiler")
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	var (
		did        int64
		mediaID    int64
		spoiler    int64
		err1, err2 error
	)
	if aidStr != "" {
		var err error
		did, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil || did <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSrv.ParseDraftParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	listIDStr := params.Get("list_id")
	var listID int64
	if listIDStr != "" {
		var lid int64
		lid, err = strconv.ParseInt(listIDStr, 10, 64)
		if err != nil || lid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		listID = lid
	}
	if mediaIDStr != "" {
		mediaID, err1 = strconv.ParseInt(mediaIDStr, 10, 64)
		spoiler, err2 = strconv.ParseInt(spoilerStr, 10, 32)
		if err1 != nil || err2 != nil || mediaID < 0 || spoiler < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	// params
	art := &model.ArtParam{
		AID:             did,
		MID:             mid,
		Title:           title,
		Content:         content,
		Summary:         summary,
		BannerURL:       bannerURL,
		Tags:            tags,
		ImageURLs:       artParam.ImageURLs,
		OriginImageURLs: artParam.OriginImageURLs,
		RealIP:          ip,
		Category:        artParam.Category,
		TemplateID:      artParam.TemplateID,
		Reprint:         artParam.Reprint,
		Words:           words,
		DynamicIntro:    dynamicIntrosStr,
		ListID:          listID,
		MediaID:         mediaID,
		Spoiler:         int32(spoiler),
	}
	// add draft
	id, err := artSrv.CreativeAddDraft(c, mid, art)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{"aid": id}, nil)
}

func webArticleList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	sortStr := params.Get("sort")
	groupStr := params.Get("group")
	categoryStr := params.Get("category")
	ip := metadata.String(c, metadata.RemoteIP)
	// check
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	sort, err := strconv.Atoi(sortStr)
	if err != nil || sort < 0 {
		sort = 0
	}
	group, err := strconv.Atoi(groupStr)
	if err != nil || group < 0 {
		group = 0
	}
	category, err := strconv.Atoi(categoryStr)
	if err != nil || category < 0 {
		category = 0
	}
	arts, err := artSrv.CreativeArticles(c, mid, int(pn), int(ps), sort, group, category, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{"artlist": arts}, nil)
}

func webDraftList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	ip := metadata.String(c, metadata.RemoteIP)
	// check
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	arts, err := artSrv.CreativeDrafts(c, mid, pn, ps, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{"artlist": arts}, nil)
}

func webDraft(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.CreativeDraft(c, aid, mid, ip))
}

func webArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.CreativeView(c, aid, mid, ip))
}

func creatorArticlePre(c *bm.Context) {
	var (
		isAuthor int
		url      string
		mid      int64
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	ia, _, err := artSrv.IsAuthor(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if ia {
		isAuthor = 1
		url = "https://member.bilibili.com/article-text/mobile"
	} else {
		isAuthor = 0
		url = "https://www.bilibili.com/read/apply/"
	}
	c.JSON(map[string]interface{}{
		"is_author":  isAuthor,
		"reason":     "", // 保持接口不变
		"submit_url": url,
	}, nil)
}

func uploadImage(c *bm.Context) {
	var (
		bs  []byte
		mid int64
	)
	midI, _ := c.Get("mid")
	mid = midI.(int64)
	log.Infov(c, log.KV("log", "creative: upload image"), log.KV("mid", mid))
	dataURI := c.Request.FormValue("file")
	if dataURI != "" {
		dataURI = strings.Split(dataURI, ",")[1]
		bs, _ = base64.StdEncoding.DecodeString(dataURI)
	} else {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		defer file.Close()
		bs, err = ioutil.ReadAll(file)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	ftype := http.DetectContentType(bs)
	if ftype != "image/jpeg" && ftype != "image/jpg" && ftype != "image/png" && ftype != "image/gif" {
		log.Error("creative: file type not allow file type(%s, mid: %v)", ftype, mid)
		c.JSON(nil, ecode.CreativeArticleImageTypeErr)
		return
	}
	url, err := artSrv.ArticleUpCover(c, ftype, bs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url":  url,
		"size": len(bs),
	}, nil)
}

func deleteDraft(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.DelArtDraft(c, aid, mid))
}

func delArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.DelArticle(c, aid, mid))
}

func withdrawArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.CreationWithdrawArticle(c, mid, aid))
}

func draftCount(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	count := artSrv.CreativeDraftCount(c, mid)
	c.JSONMap(map[string]interface{}{"count": count}, nil)
}

func articleCapture(c *bm.Context) {
	params := c.Request.Form
	originURL := params.Get("url")
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if originURL == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("capture mid(%d) origin imageURL (%s)", mid, originURL)
	_, err := url.ParseRequestURI(originURL)
	if err != nil {
		log.Error("capture check url(%s) format error(%v)", originURL, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	imgURL, size, err := artSrv.ArticleCapture(c, originURL)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url":  imgURL,
		"size": size,
	}, nil)
}
