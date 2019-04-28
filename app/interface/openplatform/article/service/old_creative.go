package service

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CreativeSubArticle submit model.
func (s *Service) CreativeSubArticle(c context.Context, mid int64, art *model.ArtParam, ak, ck, ip string) (aid int64, err error) {
	identified, _ := s.IdentifyInfo(c, mid, 0, ak, ck, ip)
	if err = s.CheckIdentify(identified); err != nil {
		log.Error("s.accountRPC.IdentifyInfo mid(%d),ip(%s)", mid, ip)
		return
	}
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
		ListID:          art.ListID,
		MediaID:         art.MediaID,
		Spoiler:         art.Spoiler,
	}
	if art.Tags != "" {
		arg.Tags = strings.Split(art.Tags, ",")
	} else {
		arg.Tags = []string{}
	}
	var a = model.TransformArticle(arg)
	if aid, err = s.AddArticle(c, a, arg.ActivityID, arg.ListID, arg.RealIP); err != nil {
		return
	}
	return
}

// CreativeUpdateArticle update model.
func (s *Service) CreativeUpdateArticle(c context.Context, mid int64, art *model.ArtParam, ak, ck, ip string) (err error) {
	identified, _ := s.IdentifyInfo(c, mid, 0, ak, ck, ip)
	if err = s.CheckIdentify(identified); err != nil {
		log.Error("s.accountRPC.IdentifyInfo mid(%d),ip(%s)", mid, ip)
		return
	}
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
		ListID:          art.ListID,
		MediaID:         art.MediaID,
		Spoiler:         art.Spoiler,
	}
	if art.Tags != "" {
		arg.Tags = strings.Split(art.Tags, ",")
	} else {
		arg.Tags = []string{}
	}
	log.Info("d.art.UpdateArticle id (%d) words (%d) ImageURLs (%s) OriginImageURLs (%s)", arg.Aid, len(arg.Content), art.ImageURLs, art.OriginImageURLs)
	var a = model.TransformArticle(arg)
	if err = s.UpdateArticle(c, a, arg.ActivityID, arg.ListID, arg.RealIP); err != nil {
		return
	}
	return
}

// CreativeAddDraft .
func (s *Service) CreativeAddDraft(c context.Context, mid int64, art *model.ArtParam) (aid int64, err error) {
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
		ListID:          art.ListID,
		MediaID:         art.MediaID,
		Spoiler:         art.Spoiler,
	}
	if art.Tags != "" {
		arg.Tags = strings.Split(art.Tags, ",")
	} else {
		arg.Tags = []string{}
	}
	log.Info("d.art.AddDraft id (%d) words (%d) ImageURLs (%s) OriginImageURLs (%s)", arg.Aid, len(arg.Content), art.ImageURLs, art.OriginImageURLs)
	d := model.TransformDraft(arg)
	aid, err = s.AddArtDraft(c, d)
	return
}

// CheckIdentify fn
func (s *Service) CheckIdentify(identify int) (err error) {
	switch identify {
	case 0:
		err = nil
	case 1:
		err = ecode.UserCheckInvalidPhone
	case 2:
		err = ecode.UserCheckNoPhone
	}
	return
}

// IdentifyInfo .
func (s *Service) IdentifyInfo(c context.Context, mid int64, phoneOnly int8, ak, ck, ip string) (ret int, err error) {
	var (
		mf  *accmdl.Profile
		arg = &accmdl.ArgMid{
			Mid: mid,
		}
	)
	if mf, err = s.accountRPC.Profile3(c, arg); err != nil {
		log.Error("d.acc.MyInfo error(%+v) | mid(%d) ck(%s) ak(%s) ip(%s) arg(%v)", err, mid, ck, ak, ip, arg)
		err = ecode.CreativeAccServiceErr
		return
	}
	//switch for FrontEnd return json format
	ret = s.switchPhoneRet(int(mf.TelStatus))
	if phoneOnly == 1 {
		return
	}
	if mf.TelStatus == 1 || mf.Identification == 1 {
		return 0, err
	}
	return
}

// 0: "已实名认证",
// 1: "根据国家实名制认证的相关要求，您需要换绑一个非170/171的手机号，才能继续进行操作。",
// 2: "根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作。",
func (s *Service) switchPhoneRet(newV int) (oldV int) {
	switch newV {
	case 0:
		oldV = 2
	case 1:
		oldV = 0
	case 2:
		oldV = 1
	}
	return
}

// CreativeArticles creative articles list
func (s *Service) CreativeArticles(c context.Context, mid int64, pn, ps, sort, group, category int, ip string) (arts *model.CreativeArtList, err error) {
	var res *model.CreationArts
	res, err = s.CreationUpperArticlesMeta(c, mid, group, category, sort, pn, ps, ip)
	arts = &model.CreativeArtList{}
	if res != nil {
		arts.Articles = make([]*model.CreativeMeta, 0, len(res.Articles))
		arts.Page = res.Page
		arts.Type = res.Type
	}
	if err != nil {
		log.Error("s.art.Articles(mid:%d) error(%+v)", mid, err)
		return
	}
	if (res == nil) || (res.Articles == nil) || (len(res.Articles) == 0) {
		log.Info("s.art.Articles(mid:%d) res(%v)", mid, res)
		return
	}
	for _, v := range res.Articles {
		id := strconv.FormatInt(v.ID, 10)
		m := &model.CreativeMeta{
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
			EditURL:         "https://member.bilibili.com/article-text/mobile?aid=" + id + "&type=2",
			DynamicIntro:    v.Dynamic,
			ImageURLs:       v.ImageURLs,
			OriginImageURLs: v.OriginImageURLs,
			List:            v.List,
		}
		if m.ImageURLs == nil {
			m.ImageURLs = []string{}
		}
		if m.OriginImageURLs == nil {
			m.OriginImageURLs = []string{}
		}
		switch m.State {
		case 0, 7:
			// 开放浏览, 可编辑
			m.ViewURL = "https://www.bilibili.com/read/cv" + id
			m.IsPreview = 0
			m.EditTimes = s.EditTimes(c, m.ID)
			m.EditURL = "https://member.bilibili.com/article-text/mobile?aid=" + id + "&type=3"
		case 4:
			// 开放浏览
			m.ViewURL = "https://www.bilibili.com/read/cv" + id
			m.IsPreview = 0
		case 5, 6:
			// 开放浏览,重复编辑待审/重复编辑未通过
			m.ViewURL = "https://www.bilibili.com/read/cv" + id
			m.IsPreview = 2
			m.PreViewURL = "https://www.bilibili.com/read/preview/" + id
			var (
				a  *model.Article
				e1 error
			)
			if a, e1 = s.ArticleVersion(c, m.ID); e1 != nil {
				log.Error("s.ArticleVersion(%d) error(%+v)", m.ID, e1)
				continue
			}
			m.Title = a.Title
			m.Reason = a.Reason
			m.Category = a.Category
			m.TemplateID = a.TemplateID
			m.ImageURLs = a.ImageURLs
			m.Summary = a.Summary
			m.Reprint = a.Reprint
			m.BannerURL = a.BannerURL
			m.OriginImageURLs = a.OriginImageURLs
			if m.State == 6 {
				m.EditURL = "https://member.bilibili.com/article-text/mobile?aid=" + id + "&type=3"
				m.EditTimes = s.EditTimes(c, m.ID)
				m.Reason, _ = s.lastReason(c, m.ID, m.State)
			}
		default:
			// 预览
			m.PreViewURL = "https://www.bilibili.com/read/preview/" + id
			m.IsPreview = 1
		}
		tags := []string{}
		m.Tags = tags
		if len(v.Tags) > 0 {
			for _, vv := range v.Tags {
				tags = append(tags, vv.Name)
			}
			m.Tags = tags
		}
		arts.Articles = append(arts.Articles, m)
	}
	return
}

// CreativeDrafts get draft list.
func (s *Service) CreativeDrafts(c context.Context, mid int64, pn, ps int, ip string) (dls *model.CreativeDraftList, err error) {
	var res *model.Drafts
	res, err = s.UpperDrafts(c, mid, pn, ps)
	if err != nil {
		log.Error("s.art.Drafts(mid:%d) error(%+v)", mid, err)
		return
	}
	if res == nil || res.Drafts == nil || len(res.Drafts) <= 0 {
		log.Info("s.art.Drafts(mid:%d) res(%+v)", mid, res)
		return
	}
	ms := make([]*model.CreativeMeta, 0, len(res.Drafts))
	for _, v := range res.Drafts {
		id := strconv.FormatInt(v.ID, 10)
		m := &model.CreativeMeta{
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
			List:            v.List,
			EditURL:         "https://member.bilibili.com/article-text/mobile?aid=" + id,
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
	dls = &model.CreativeDraftList{
		DraftURL: "https://member.bilibili.com/creative/app/article_drafts",
	}
	dls.Drafts = ms
	dls.Page = res.Page
	return
}

// CreativeDraft get draft.
func (s *Service) CreativeDraft(c context.Context, aid, mid int64, ip string) (res *model.CreativeMeta, err error) {
	var df *model.Draft
	if df, err = s.ArtDraft(c, aid, mid); err != nil {
		return
	}
	if df == nil || df.Article == nil {
		err = ecode.CreativeArticleNotExist
		return
	}
	res = &model.CreativeMeta{
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
		List:            df.List,
		MediaID:         df.Article.Media.MediaID,
		Spoiler:         df.Article.Media.Spoiler,
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

// CreativeView get article detail.
func (s *Service) CreativeView(c context.Context, aid, mid int64, ip string) (res *model.CreativeMeta, err error) {
	var art *model.Article
	if art, err = s.CreationArticle(c, aid, mid); err != nil {
		return
	}
	res = &model.CreativeMeta{
		ID:              art.ID,
		Category:        art.Category,
		Title:           art.Title,
		Content:         art.Content,
		Summary:         art.Summary,
		BannerURL:       art.BannerURL,
		TemplateID:      art.TemplateID,
		State:           art.State,
		Reprint:         art.Reprint,
		Reason:          art.Reason,
		PTime:           art.PublishTime,
		Author:          art.Author,
		Stats:           art.Stats,
		CTime:           art.Ctime,
		MTime:           art.Mtime,
		DynamicIntro:    art.Dynamic,
		ImageURLs:       art.ImageURLs,
		OriginImageURLs: art.OriginImageURLs,
		List:            art.List,
		MediaID:         art.Media.MediaID,
		Spoiler:         art.Media.Spoiler,
	}
	if res.State == model.StateOpen || res.State == model.StateReReject || res.State == model.StateRePass {
		res.EditTimes = s.EditTimes(c, res.ID)
	}
	if res.ImageURLs == nil {
		res.ImageURLs = []string{}
	}
	if res.OriginImageURLs == nil {
		res.OriginImageURLs = []string{}
	}
	if len(art.Tags) > 0 {
		var tags []string
		for _, v := range art.Tags {
			tags = append(tags, v.Name)
		}
		res.Tags = tags
	}
	return
}

// CreativeDraftCount count of upper's drafts
func (s *Service) CreativeDraftCount(c context.Context, mid int64) (count int) {
	count, _ = s.dao.CountUpperDraft(c, mid)
	return
}

// ArticleCapture capture a new image.
func (s *Service) ArticleCapture(c context.Context, url string) (loc string, size int, err error) {
	loc, size, err = s.dao.Capture(c, url)
	if err != nil {
		log.Error("s.bfs.Capture error(%v)", err)
	}
	return
}

// SetMediaScore set media score.
func (s *Service) SetMediaScore(c context.Context, score, aid, mediaID, mid int64) (err error) {
	return s.dao.SetScore(c, score, aid, mediaID, mid)
}

// DelMediaScore get media score.
func (s *Service) DelMediaScore(c context.Context, aid, mediaID, mid int64) (err error) {
	return s.dao.DelScore(c, aid, mediaID, mid)
}
