package article

import (
	"context"
	"strconv"
	"strings"

	artMdl "go-common/app/interface/main/creative/model/article"
	"go-common/library/ecode"
)

// ParseParam  parse article param which type is int.
func (s *Service) ParseParam(c context.Context, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs string) (art *artMdl.ArtParam, err error) {
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
	art = &artMdl.ArtParam{
		Category:        category,
		TemplateID:      int32(tid),
		Reprint:         int32(reprint),
		ImageURLs:       imgs,
		OriginImageURLs: oimgs,
	}
	return
}

// ParseDraftParam  parse draft param which type is int.
func (s *Service) ParseDraftParam(c context.Context, categoryStr, reprintStr, tidStr, imageURLs, originImageURLs string) (art *artMdl.ArtParam, err error) {
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
	art = &artMdl.ArtParam{
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
