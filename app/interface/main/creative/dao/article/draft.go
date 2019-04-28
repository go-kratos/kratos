package article

import (
	"context"
	artMdl "go-common/app/interface/main/creative/model/article"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"strings"
)

// AddDraft add draft.
func (d *Dao) AddDraft(c context.Context, art *artMdl.ArtParam) (id int64, err error) {
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
	log.Info("d.art.AddDraft id (%d) words (%d) ImageURLs (%s) OriginImageURLs (%s)", arg.Aid, len(arg.Content), art.ImageURLs, art.OriginImageURLs)
	if id, err = d.art.AddArtDraft(c, arg); err != nil {
		arg.Content = ""
		log.Error("d.art.AddArtDraft (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// DelDraft delete draft.
func (d *Dao) DelDraft(c context.Context, aid, mid int64, ip string) (err error) {
	var arg = &model.ArgAidMid{
		Aid:    aid,
		Mid:    mid,
		RealIP: ip,
	}
	if err = d.art.DelArtDraft(c, arg); err != nil {
		log.Error("d.art.DelArtDraft (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}

// Draft get draft detail.
func (d *Dao) Draft(c context.Context, aid, mid int64, ip string) (res *model.Draft, err error) {
	var arg = &model.ArgAidMid{
		Aid:    aid,
		Mid:    mid,
		RealIP: ip,
	}
	if res, err = d.art.ArtDraft(c, arg); err != nil {
		log.Error("d.art.ArtDraft (%v) error(%v)", arg, err)
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

// Drafts get draft list.
func (d *Dao) Drafts(c context.Context, mid int64, pn, ps int, ip string) (res *model.Drafts, err error) {
	var arg = &model.ArgUpDraft{
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		RealIP: ip,
	}
	if res, err = d.art.UpperDrafts(c, arg); err != nil {
		log.Error("d.art.UpperDrafts (%v) error(%v)", arg, err)
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			err = ecode.CreativeArticleRPCErr
		}
	}
	return
}
