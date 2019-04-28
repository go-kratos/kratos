package article

import (
	"context"
	"go-common/library/log"
	"strconv"
	"strings"

	artMdl "go-common/app/interface/main/creative/model/article"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
)

// Articles get article list.
func (d *Dao) Articles(c context.Context, mid int64, pn, ps, sort, group, category int, ip string) (res *model.CreationArts, err error) {
	var arg = &model.ArgCreationArts{
		Mid:      mid,
		Pn:       pn,
		Ps:       ps,
		Sort:     sort,
		Group:    group,
		Category: category,
		RealIP:   ip,
	}
	if res, err = d.art.CreationUpperArticles(c, arg); err != nil {
		log.Error("d.art.CreationUpperArticles (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// Categories list all category contain child.
func (d *Dao) Categories(c context.Context, ip string) (res *model.Categories, err error) {
	var arg = &model.ArgIP{
		RealIP: ip,
	}
	if res, err = d.art.Categories(c, arg); err != nil {
		log.Error("d.art.Categories (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// CategoriesMap list all category.
func (d *Dao) CategoriesMap(c context.Context, ip string) (res map[int64]*model.Category, err error) {
	var arg = &model.ArgIP{
		RealIP: ip,
	}
	if res, err = d.art.CategoriesMap(c, arg); err != nil {
		log.Error("d.art.CategoriesMap (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// Article get article detail.
func (d *Dao) Article(c context.Context, aid, mid int64, ip string) (res *model.Article, err error) {
	var arg = &model.ArgAidMid{
		Aid:    aid,
		Mid:    mid,
		RealIP: ip,
	}
	if res, err = d.art.CreationArticle(c, arg); err != nil {
		log.Error("d.art.CreationArticle (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	if res == nil || res.Meta == nil {
		return
	}
	log.Info("d.art.CreationArticle id (%d) words (%d)", res.Meta.ID, len(res.Content))
	return
}

// AddArticle add article.
func (d *Dao) AddArticle(c context.Context, art *artMdl.ArtParam) (id int64, err error) {
	var arg = &model.ArgArticle{
		Aid:             art.AID,
		Mid:             art.MID,
		Category:        art.Category,
		State:           art.State,
		Reprint:         art.Reprint,
		TemplateID:      art.TemplateID,
		Title:           art.Title,
		BannerURL:       art.BannerURL,
		Content:         art.Content,
		Summary:         art.Summary,
		RealIP:          art.RealIP,
		Words:           art.Words,
		DynamicIntro:    art.DynamicIntro,
		ImageURLs:       art.ImageURLs,
		OriginImageURLs: art.OriginImageURLs,
		ActivityID:      art.ActivityID,
	}
	if art.Tags != "" {
		arg.Tags = strings.Split(art.Tags, ",")
	} else {
		arg.Tags = []string{}
	}
	log.Info("d.art.AddArticle id (%d) words (%d) ImageURLs (%s) OriginImageURLs (%s)", arg.Aid, len(arg.Content), art.ImageURLs, art.OriginImageURLs)
	if id, err = d.art.AddArticle(c, arg); err != nil {
		arg.Content = ""
		log.Error("d.art.AddArticle (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// UpdateArticle update article.
func (d *Dao) UpdateArticle(c context.Context, art *artMdl.ArtParam) (err error) {
	var arg = &model.ArgArticle{
		Aid:             art.AID,
		Mid:             art.MID,
		Category:        art.Category,
		State:           art.State,
		Reprint:         art.Reprint,
		TemplateID:      art.TemplateID,
		Title:           art.Title,
		BannerURL:       art.BannerURL,
		Content:         art.Content,
		Summary:         art.Summary,
		RealIP:          art.RealIP,
		Words:           art.Words,
		DynamicIntro:    art.DynamicIntro,
		ImageURLs:       art.ImageURLs,
		OriginImageURLs: art.OriginImageURLs,
	}
	if art.Tags != "" {
		arg.Tags = strings.Split(art.Tags, ",")
	} else {
		arg.Tags = []string{}
	}
	log.Info("d.art.UpdateArticle id (%d) words (%d) ImageURLs (%s) OriginImageURLs (%s)", arg.Aid, len(arg.Content), art.ImageURLs, art.OriginImageURLs)
	if err = d.art.UpdateArticle(c, arg); err != nil {
		arg.Content = ""
		log.Error("d.art.UpdateArticle (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// DelArticle delete article.
func (d *Dao) DelArticle(c context.Context, aid, mid int64, ip string) (err error) {
	var arg = &model.ArgAidMid{
		Aid:    aid,
		Mid:    mid,
		RealIP: ip,
	}
	if err = d.art.DelArticle(c, arg); err != nil {
		log.Error("d.art.AddArticle (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// WithDrawArticle withdraw  article.
func (d *Dao) WithDrawArticle(c context.Context, aid, mid int64, ip string) (err error) {
	var arg = &model.ArgAidMid{
		Aid:    aid,
		Mid:    mid,
		RealIP: ip,
	}
	if err = d.art.CreationWithdrawArticle(c, arg); err != nil {
		log.Error("d.art.CreationWithdrawArticle (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// IsAuthor checks that whether user has permission to write article.
func (d *Dao) IsAuthor(c context.Context, mid int64, ip string) (res bool, err error) {
	var arg = &model.ArgMid{
		Mid:    mid,
		RealIP: ip,
	}
	if res, err = d.art.IsAuthor(c, arg); err != nil {
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			log.Error("d.art.IsAuthor (%v) error(%v)", arg, err)
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// RemainCount article up limit.
func (d *Dao) RemainCount(c context.Context, mid int64, ip string) (res int, err error) {
	var arg = &model.ArgMid{
		Mid:    mid,
		RealIP: ip,
	}
	if res, err = d.art.ArticleRemainCount(c, arg); err != nil {
		log.Error("d.art.ArticleRemainCount (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// ArticleStat article stat
func (d *Dao) ArticleStat(c context.Context, mid int64, ip string) (res model.UpStat, err error) {
	arg := &model.ArgMid{Mid: mid, RealIP: ip}
	if res, err = d.art.CreationUpStat(c, arg); err != nil {
		log.Error("d.art.UpStat(%+v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// ThirtyDayArticle thirty day article
func (d *Dao) ThirtyDayArticle(c context.Context, mid int64, ip string) (res []*model.ThirtyDayArticle, err error) {
	arg := &model.ArgMid{Mid: mid, RealIP: ip}
	if res, err = d.art.CreationUpThirtyDayStat(c, arg); err != nil {
		log.Error("d.art.CreationUpThirtyDayStat(%+v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// ArticleMetas batch get articles by aids.
func (d *Dao) ArticleMetas(c context.Context, aids []int64, ip string) (res map[int64]*model.Meta, err error) {
	arg := &model.ArgAids{Aids: aids, RealIP: ip}
	if res, err = d.art.ArticleMetas(c, arg); err != nil {
		log.Error("d.art.ArticleMetas(%+v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	log.Info("d.art.ArticleMetas aids(%v)", aids)
	return
}
