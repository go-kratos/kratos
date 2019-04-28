package service

import (
	"context"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"

	strip "github.com/grokify/html-strip-tags-go"
)

var (
	_zeroWidthReg = regexp.MustCompile(`[\x{200b}]+`)
	_nocharReg    = []*regexp.Regexp{
		// regexp.MustCompile(`[\p{Hangul}]+`),  // kr
		regexp.MustCompile(`[\p{Tibetan}]+`), // tibe
		regexp.MustCompile(`[\p{Arabic}]+`),  // arabic
	}
	_chineseReg = regexp.MustCompile(`[\p{Han}]+`) // chinese
)

func (s *Service) allowRepeat(c context.Context, mid int64, title string) (ok bool) {
	log.Info("allowRepeat check start | mid(%d) title(%s).", mid, title)
	exist, _ := s.dao.SubmitCache(c, mid, title)
	log.Info("allowRepeat from cache | mid(%d) title(%s) exist(%d).", mid, title, exist)
	if !exist {
		log.Info("allowRepeat not exist | mid(%d) title(%s)", mid, title)
		s.dao.AddSubmitCache(c, mid, title)
		log.Info("allowRepeat add cache | mid(%d) title(%s).", mid, title)
		ok = true
		return
	}
	dao.PromInfo("creation:禁止重复标题")
	return
}

func (s *Service) preMust(c context.Context, art *model.Article) (err error) {
	var ok bool
	if art.Title, ok = s.checkTitle(art.Title); !ok || art.Title == "" {
		log.Error("s.checkTitle mid(%d) art.Title(%s) title contains illegal char or is empty", art.Author.Mid, art.Title)
		err = ecode.CreativeArticleTitleErr
		return
	}
	if art.Content, ok = s.checkContent(art.Content); !ok {
		log.Error("s.checkContent mid(%d) content too long", art.Author.Mid)
		err = ecode.CreativeArticleContentErr
		return
	}
	if !s.allowCategory(art.Category.ID) {
		log.Error("s.allowCategory mid(%d) art.Category(%d) not exists", art.Author.Mid, art.Category)
		err = ecode.CreativeArticleCategoryErr
		return
	}
	if !s.allowReprints(int8(art.Reprint)) {
		log.Error("s.allowReprints mid(%d) art.Reprint(%d) illegal reprint", art.Author.Mid, art.Reprint)
		err = ecode.CreativeArticleReprintErr
		return
	}
	if !s.allowTID(int8(art.TemplateID)) {
		log.Error("s.allowTID mid(%d) art.TemplateID(%d) illegal reprint", art.Author.Mid, art.TemplateID)
		err = ecode.CreativeArticleTIDErr
		return
	}
	if !model.ValidTemplate(art.TemplateID, art.ImageURLs) {
		err = ecode.ArtCreationTplErr
		return
	}
	if !s.allowTag(art.Tags) {
		log.Error("s.allowTag mid(%d) art.Tags(%s) tag name or number too large", art.Author.Mid, art.Tags)
		err = ecode.CreativeArticleTagErr
	}
	if art.Dynamic, ok = s.allowDynamicIntro(art.Dynamic); !ok {
		log.Error("s.checkDynamicIntro mid(%d) art.DynamicIntro(%s) title contains illegal char", art.Author.Mid, art.Dynamic)
		err = ecode.CreativeDynamicIntroErr
		return
	}
	return
}

func (s *Service) checkTitle(title string) (ct string, ok bool) {
	title = strings.TrimSpace(title)
	allCount := utf8.RuneCountInString(title)
	enCount := utf8.RuneCountInString(_chineseReg.ReplaceAllString(title, ""))
	chineseCount := allCount - enCount
	if chineseCount*2+enCount > 80 {
		return
	}
	for _, reg := range _nocharReg {
		if reg.MatchString(title) {
			return
		}
	}
	ct = _zeroWidthReg.ReplaceAllString(title, "")
	if utf8.RuneCountInString(ct) <= 0 {
		return
	}
	ok = true
	return
}

func (s *Service) contentStripSize(content string) (count int) {
	stripped := strip.StripTags(content)
	stripped = html.UnescapeString(stripped)
	stripped = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, stripped)
	stripped = strings.Replace(stripped, "\u200B", "", -1)
	stripped = strings.Replace(stripped, "\u00a0", "", -1)
	// utf16 size
	for _, r := range stripped {
		count += len(utf16.Encode([]rune{rune(r)}))
	}
	// 图片计算为一个字
	offset := strings.Count(content, "<img") - strings.Count(stripped, "<img")
	count += offset
	return
}

func (s *Service) checkContent(content string) (ct string, ok bool) {
	ct = strings.TrimSpace(content)
	ct = _zeroWidthReg.ReplaceAllString(ct, "")
	if len(ct) > s.c.Article.MaxContentSize {
		return
	}
	size := s.contentStripSize(ct)
	if size < s.c.Article.MinContentLength || size > s.c.Article.MaxContentLength {
		return
	}
	ok = true
	return
}

func (s *Service) preArticleCheck(c context.Context, art *model.Article) (err error) {
	if !s.allowRepeat(c, art.Author.Mid, art.Title) {
		err = ecode.CreativeArticleCanNotRepeat
		return
	}
	if err = s.preMust(c, art); err != nil {
		return
	}
	return
}

func (s *Service) allowCategory(cid int64) (ok bool) {
	_, ok = s.categoriesMap[cid]
	return
}

func (s *Service) allowReprints(cp int8) (ok bool) {
	ok = model.InReprints(cp)
	return
}

func (s *Service) allowTID(tid int8) (ok bool) {
	ok = model.InTemplateID(tid)
	return
}

func (s *Service) allowTag(tags []*model.Tag) (ok bool) {
	if (len(tags) > 12) || (len(tags) == 0) {
		return
	}
	for _, tag := range tags {
		if _zeroWidthReg.MatchString(tag.Name) {
			return
		}
		if (utf8.RuneCountInString(tag.Name) > 20) || (tag.Name == "") {
			return
		}
	}
	return true
}

//allowDynamicIntro 移动端动态推荐语，选填，不能超过233字.
func (s *Service) allowDynamicIntro(dynamicIntro string) (ct string, ok bool) {
	ct = strings.TrimSpace(dynamicIntro)
	if utf8.RuneCountInString(ct) > 233 {
		return
	}
	ok = true
	return
}

func (s *Service) preDraftCheck(c context.Context, art *model.Draft) (err error) {
	if art.Title == "" {
		art.Title = "无标题"
	}
	var ok bool
	if art.Title, ok = s.checkTitle(art.Title); !ok {
		log.Error("s.checkTitle mid(%d) art.Title(%s) title contains illegal char or is empty", art.Author.Mid, art.Title)
		err = ecode.CreativeArticleTitleErr
	}
	return
}

// ParseParam  parse article param which type is int.
func (s *Service) ParseParam(c context.Context, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs string) (art *model.ArtParam, err error) {
	var (
		category     int64
		tid, reprint int
	)
	category, err = strconv.ParseInt(categoryStr, 10, 64)
	if err != nil || category <= 0 { //文章要求必须传大于0的分类
		err = ecode.CreativeArticleCategoryErr
		return
	}
	tid, err = strconv.Atoi(tidStr)
	if err != nil || tid < 0 {
		err = ecode.CreativeArticleTIDErr
		return
	}
	reprint, err = strconv.Atoi(reprintStr)
	if err != nil || reprint < 0 {
		err = ecode.CreativeArticleReprintErr
		return
	}
	imgs, oimgs, err := ParseImageURLs(imageURLs, originImageURLs)
	if err != nil {
		return
	}
	art = &model.ArtParam{
		Category:        category,
		TemplateID:      int32(tid),
		Reprint:         int32(reprint),
		ImageURLs:       imgs,
		OriginImageURLs: oimgs,
	}
	return
}

// ParseDraftParam  parse draft param which type is int.
func (s *Service) ParseDraftParam(c context.Context, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs string) (art *model.ArtParam, err error) {
	var (
		category     int64
		tid, reprint int
	)
	if categoryStr != "" {
		category, err = strconv.ParseInt(categoryStr, 10, 64)
		if err != nil || category < 0 {
			err = ecode.CreativeArticleCategoryErr
			return
		}
	}
	if tidStr != "" {
		tid, err = strconv.Atoi(tidStr)
		if err != nil || tid < 0 {
			err = ecode.CreativeArticleTIDErr
			return
		}
	}
	if reprintStr != "" {
		reprint, err = strconv.Atoi(reprintStr)
		if err != nil || reprint < 0 {
			err = ecode.CreativeArticleReprintErr
			return
		}
	}
	imgs, oimgs, err := ParseImageURLs(imageURLs, originImageURLs)
	if err != nil {
		return
	}
	art = &model.ArtParam{
		Category:        category,
		TemplateID:      int32(tid),
		Reprint:         int32(reprint),
		ImageURLs:       imgs,
		OriginImageURLs: oimgs,
	}
	return
}

//ParseImageURLs parse img urls to []string.
func ParseImageURLs(imageURLs, originImageURLs string) (imgs, oimgs []string, err error) {
	if originImageURLs == "" {
		originImageURLs = imageURLs
	}
	imgs = strings.Split(imageURLs, ",")
	oimgs = strings.Split(originImageURLs, ",")
	if len(imgs) != len(oimgs) {
		err = ecode.CreativeArticleImageURLsErr
	}
	return
}
