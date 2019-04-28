package article

import (
	"context"

	artMdl "go-common/app/interface/main/creative/model/article"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddDraft add draft.
func (s *Service) AddDraft(c context.Context, mid int64, art *artMdl.ArtParam) (aid int64, err error) {
	return s.art.AddDraft(c, art)
}

// DelDraft delete draft.
func (s *Service) DelDraft(c context.Context, aid, mid int64, ip string) (err error) {
	if err = s.art.DelDraft(c, aid, mid, ip); err != nil {
		log.Error("s.art.DelArticle(%d) error(%v)", aid, err)
	}
	return
}

// Draft get draft.
func (s *Service) Draft(c context.Context, aid, mid int64, ip string) (res *artMdl.Meta, err error) {
	var df *article.Draft
	if df, err = s.art.Draft(c, aid, mid, ip); err != nil {
		return
	}
	if df == nil || df.Article == nil {
		err = ecode.CreativeArticleNotExist
		return
	}
	res = &artMdl.Meta{
		ID:              df.Article.ID,
		Category:        df.Article.Category,
		Title:           df.Article.Title,
		Content:         df.Article.Content,
		Summary:         df.Article.Summary,
		BannerURL:       df.Article.BannerURL,
		TemplateID:      df.Article.TemplateID,
		State:           df.Article.State,
		Author:          df.Article.Author,
		Stats:           df.Article.Stats,
		Reprint:         df.Article.Reprint,
		Reason:          df.Article.Reason,
		PTime:           df.Article.PublishTime,
		CTime:           df.Article.Ctime,
		MTime:           df.Article.Mtime,
		DynamicIntro:    df.Article.Dynamic,
		ImageURLs:       df.ImageURLs,
		OriginImageURLs: df.OriginImageURLs,
	}
	if res.ImageURLs == nil {
		res.ImageURLs = []string{}
	}
	if res.OriginImageURLs == nil {
		res.OriginImageURLs = []string{}
	}
	res.Tags = df.Tags
	if len(df.Tags) == 0 {
		res.Tags = []string{}
	}
	return
}

// Drafts get draft list.
func (s *Service) Drafts(c context.Context, mid int64, pn, ps int, ip string) (dls *artMdl.DraftList, err error) {
	var res *article.Drafts
	res, err = s.art.Drafts(c, mid, pn, ps, ip)
	if err != nil || res == nil || res.Drafts == nil || len(res.Drafts) <= 0 {
		if err != nil {
			log.Error("s.art.Drafts(%d) res(%v) error(%v)", mid, res, err)
		}
		return
	}
	ms := make([]*artMdl.Meta, 0, len(res.Drafts))
	for _, v := range res.Drafts {
		m := &artMdl.Meta{
			ID:              v.ID,
			Category:        v.Category,
			Title:           v.Title,
			Summary:         v.Summary,
			BannerURL:       v.BannerURL,
			TemplateID:      v.TemplateID,
			State:           v.State,
			Reprint:         v.Reprint,
			Reason:          v.Reason,
			PTime:           v.PublishTime,
			Author:          v.Author,
			Stats:           v.Stats,
			CTime:           v.Ctime,
			MTime:           v.Mtime,
			DynamicIntro:    v.Dynamic,
			ImageURLs:       v.ImageURLs,
			OriginImageURLs: v.OriginImageURLs,
		}
		if m.ImageURLs == nil {
			m.ImageURLs = []string{}
		}
		if m.OriginImageURLs == nil {
			m.OriginImageURLs = []string{}
		}
		m.Tags = v.Tags
		if len(v.Tags) == 0 {
			m.Tags = []string{}
		}
		ms = append(ms, m)
	}
	dls = &artMdl.DraftList{
		DraftURL: s.c.H5Page.Draft,
	}
	dls.Drafts = ms
	dls.Page = res.Page
	return
}
