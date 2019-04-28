package http

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/model/article"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func webArticlePre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
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
	mf, err := accSvc.MyInfo(c, mid, ip, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	categories, err := artSvc.Categories(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	rc, _ := artSvc.RemainCount(c, mid, ip)
	c.JSON(map[string]interface{}{
		"categories": categories,
		"myinfo":     mf,
		"toplimit":   rc,
	}, nil)
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
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
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
	var (
		aid   int64
		actID int64
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
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSvc.ParseParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// params
	art := &article.ArtParam{
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
	}
	// submit
	id, err := artSvc.SubArticle(c, mid, art, "", ck, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": id,
	}, nil)
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
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
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
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSvc.ParseParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// params
	art := &article.ArtParam{
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
	}
	c.JSON(nil, artSvc.UpdateArticle(c, mid, art, "", ck, ip))
}

func webDelArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
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
	c.JSON(artSvc.DelArticle(c, aid, mid, ip), nil)
}

func webArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
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
	art, err := artSvc.View(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(art, nil)
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
	arts, err := artSvc.Articles(c, mid, int(pn), int(ps), sort, group, category, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"artlist": arts,
	}, nil)
}

func webWithDrawArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
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
	c.JSON(nil, artSvc.WithDrawArticle(c, aid, mid, ip))
}

func webArticleUpCover(c *bm.Context) {
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
	log.Info("webArticleUpCover mid(%d)", mid)
	cover := c.Request.Form.Get("cover")
	c.Request.Form.Del("cover") // NOTE: make sure write log concise
	ss := strings.Split(cover, ",")
	if len(ss) != 2 || len(ss[1]) == 0 {
		log.Error("webArticleUpCover format error mid(%d)|cover(%s)|coverSlice(%s)", mid, cover, ss)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		log.Error("webArticleUpCover base64.StdEncoding.DecodeString(%s)|mid(%d)|error(%v)", ss[1], mid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	if ftype != "image/jpeg" && ftype != "image/jpg" && ftype != "image/png" && ftype != "image/gif" {
		log.Error("webArticleUpCover file type not allow file type(%s)|mid(%d)", ftype, mid)
		c.JSON(nil, ecode.CreativeArticleImageTypeErr)
		return
	}
	url, err := artSvc.ArticleUpCover(c, ftype, bs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url":  url,
		"size": len(bs),
	}, nil)
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
	dynamicIntrosStr := params.Get("dynamic_intro")
	originImageURLs := params.Get("origin_image_urls")
	ip := metadata.String(c, metadata.RemoteIP)
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
	var (
		did int64
		err error
	)
	if aidStr != "" {
		did, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil || did <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	words, _ := strconv.ParseInt(wordsStr, 10, 64)
	artParam, err := artSvc.ParseDraftParam(c, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// params
	art := &article.ArtParam{
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
	}
	// add draft
	id, err := artSvc.AddDraft(c, mid, art)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": id,
	}, nil)
}

func webDeleteDraft(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
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
	c.JSON(nil, artSvc.DelDraft(c, aid, mid, ip))
}

func webDraft(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
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
	art, err := artSvc.Draft(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(art, nil)
}

func webDraftList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	ip := metadata.String(c, metadata.RemoteIP)
	// check
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
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	arts, err := artSvc.Drafts(c, mid, pn, ps, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"artlist": arts,
	}, nil)
}

func webAuthor(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
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
	isAuthor, _ := artSvc.IsAuthor(c, mid, ip)
	c.JSON(map[string]interface{}{
		"mid":       mid,
		"is_author": isAuthor,
	}, nil)
}

func webArticleCapture(c *bm.Context) {
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
	imgURL, size, err := artSvc.ArticleCapture(c, originURL)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url":  imgURL,
		"size": size,
	}, nil)
}
