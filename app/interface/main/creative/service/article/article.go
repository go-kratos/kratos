package article

import (
	"context"
	artMdl "go-common/app/interface/main/creative/model/article"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
)

// SubArticle submit article.
func (s *Service) SubArticle(c context.Context, mid int64, art *artMdl.ArtParam, ak, ck, ip string) (aid int64, err error) {
	identified, _ := s.acc.IdentifyInfo(c, mid, 0, ip)
	if err = s.acc.CheckIdentify(identified); err != nil {
		log.Error("s.acc.IdentifyInfo mid(%d),ip(%s)", mid, ip)
		return
	}
	aid, err = s.art.AddArticle(c, art)
	return
}

// UpdateArticle update article.
func (s *Service) UpdateArticle(c context.Context, mid int64, art *artMdl.ArtParam, ak, ck, ip string) (err error) {
	identified, _ := s.acc.IdentifyInfo(c, mid, 0, ip)
	if err = s.acc.CheckIdentify(identified); err != nil {
		log.Error("s.acc.IdentifyInfo mid(%d),ip(%s)", mid, ip)
		return
	}
	return s.art.UpdateArticle(c, art)
}

// DelArticle delete article.
func (s *Service) DelArticle(c context.Context, aid, mid int64, ip string) (err error) {
	if err = s.art.DelArticle(c, aid, mid, ip); err != nil {
		log.Error("s.art.DelArticle(%d) error(%v)", aid, err)
		return
	}
	return
}

// View get article detail.
func (s *Service) View(c context.Context, aid, mid int64, ip string) (res *artMdl.Meta, err error) {
	var art *article.Article
	if art, err = s.art.Article(c, aid, mid, ip); err != nil {
		return
	}
	res = &artMdl.Meta{
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

// Articles get article list.
func (s *Service) Articles(c context.Context, mid int64, pn, ps, sort, group, category int, ip string) (arts *artMdl.ArtList, err error) {
	var res *article.CreationArts
	res, err = s.art.Articles(c, mid, pn, ps, sort, group, category, ip)
	if err != nil || res == nil || res.Articles == nil || len(res.Articles) <= 0 {
		if err != nil {
			log.Error("s.art.Articles(%d) res(%v) error(%v)", mid, res, err)
		}
		return
	}
	ms := make([]*artMdl.Meta, 0, len(res.Articles))
	for _, v := range res.Articles {
		id := strconv.FormatInt(v.ID, 10)
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
			EditURL:         "https://member.bilibili.com/article-text/mobile?aid=" + id + "&type=2",
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
		if m.State == 0 {
			m.ViewURL = "https://www.bilibili.com/read/cv" + id
			m.IsPreview = 0
		} else { //预览
			m.ViewURL = "https://www.bilibili.com/read/preview/" + id
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
		ms = append(ms, m)
	}
	arts = &artMdl.ArtList{}
	arts.Articles = ms
	arts.Page = res.Page
	arts.Type = res.Type
	return
}

// Categories get article category.
func (s *Service) Categories(c context.Context) (*article.Categories, error) {
	return s.art.Categories(c, "")
}

// WithDrawArticle withdraw article.
func (s *Service) WithDrawArticle(c context.Context, aid, mid int64, ip string) (err error) {
	if err = s.art.WithDrawArticle(c, aid, mid, ip); err != nil {
		log.Error("s.art.WithdrawArticle(%d,%d) error(%v)", aid, mid, err)
	}
	return
}

// ArticleUpCover article upload cover.
func (s *Service) ArticleUpCover(c context.Context, fileType string, body []byte) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		log.Error("s.ArticleUpCover file not exist")
		return
	}
	if len(body) > s.c.BFS.MaxFileSize {
		log.Error("s.ArticleUpCover too max file size")
		err = ecode.FileTooLarge
		return
	}
	url, err = s.bfs.Upload(c, fileType, body)
	if err != nil {
		log.Error("s.bfs.Upload error(%v)", err)
	}
	return
}

// IsAuthor checks that whether user has permission to write article.
func (s *Service) IsAuthor(c context.Context, mid int64, ip string) (ok bool, err error) {
	if ok, err = s.art.IsAuthor(c, mid, ip); err != nil {
		log.Error("s.art.IsAuthor(%v)", err)
	}
	return
}

// RemainCount article up limit.
func (s *Service) RemainCount(c context.Context, mid int64, ip string) (rc int, err error) {
	rc, err = s.art.RemainCount(c, mid, ip)
	return
}

// ArticleCapture article capture.
func (s *Service) ArticleCapture(c context.Context, url string) (loc string, size int, err error) {
	loc, size, err = s.bfs.Capture(c, url)
	if err != nil {
		log.Error("s.bfs.Capture error(%v)", err)
	}
	return
}

// ArticleStat get article base data.
func (s *Service) ArticleStat(c context.Context, mid int64, ip string) (stat article.UpStat, err error) {
	stat, err = s.art.ArticleStat(c, mid, ip)
	return
}
