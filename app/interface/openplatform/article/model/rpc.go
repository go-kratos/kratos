package model

import (
	xtime "go-common/library/time"
	"time"
)

// ArgArticle .
type ArgArticle struct {
	Action          int
	Aid             int64
	Category        int64
	Title           string
	Summary         string
	BannerURL       string
	TemplateID      int32
	State           int32
	Mid             int64
	Reprint         int32
	ImageURLs       []string
	OriginImageURLs []string
	Tags            []string
	Content         string
	Words           int64
	DynamicIntro    string
	ActivityID      int64
	ListID          int64
	RealIP          string
	MediaID         int64
	Spoiler         int32
}

// ArgAid .
type ArgAid struct {
	Aid    int64
	RealIP string
}

// ArgPtime .
type ArgPtime struct {
	Aid     int64
	PubTime int64
	RealIP  string
}

// ArgAidMid .
type ArgAidMid struct {
	Aid    int64
	Mid    int64
	RealIP string
}

// ArgAids .
type ArgAids struct {
	Aids   []int64
	RealIP string
}

// ArgMid .
type ArgMid struct {
	Mid    int64
	RealIP string
}

// ArgMidAids .
type ArgMidAids struct {
	Mid    int64
	Aids   []int64
	RealIP string
}

// ArgCreationArts .
type ArgCreationArts struct {
	Mid      int64
	Sort     int
	Group    int
	Category int
	Pn, Ps   int
	RealIP   string
}

// ArgStats .
type ArgStats struct {
	*Stats
	Aid int64
}

// ArgIP .
type ArgIP struct {
	RealIP string
}

// ArgUpsArts .
type ArgUpsArts struct {
	Mids   []int64
	Pn, Ps int
	RealIP string
}

// ArgUpArts .
type ArgUpArts struct {
	Mid    int64
	Pn, Ps int
	Sort   int
	RealIP string
}

// ArgRecommends .
type ArgRecommends struct {
	Cid    int64
	Sort   int
	Aids   []int64
	Pn, Ps int
	RealIP string
}

// ArgUpDraft .
type ArgUpDraft struct {
	Mid    int64
	Pn, Ps int
	RealIP string
}

// ArgAidCid .
type ArgAidCid struct {
	Aid    int64
	Cid    int64
	RealIP string
}

// ArgAidContent .
type ArgAidContent struct {
	Aid     int64
	Content string
	RealIP  string
}

// ArgFav .
type ArgFav struct {
	Mid    int64
	Pn, Ps int
	RealIP string
}

// ArgAuthor .
type ArgAuthor struct {
	Mid    int64
	RealIP string
}

// ArgSort .
type ArgSort struct {
	Aid     int64
	Changed [][2]int64
	RealIP  string
}

// ArgNewArt .
type ArgNewArt struct {
	PubTime int64
	RealIP  string
}

// TransformArticle .
func TransformArticle(arg *ArgArticle) *Article {
	a := &Article{
		Meta: &Meta{
			ID:              arg.Aid,
			Category:        &Category{ID: arg.Category},
			Title:           arg.Title,
			Summary:         arg.Summary,
			BannerURL:       arg.BannerURL,
			TemplateID:      arg.TemplateID,
			State:           arg.State,
			Author:          &Author{Mid: arg.Mid},
			Reprint:         arg.Reprint,
			ImageURLs:       arg.ImageURLs,
			OriginImageURLs: arg.OriginImageURLs,
			Words:           arg.Words,
			Dynamic:         arg.DynamicIntro,
			Media:           &Media{MediaID: arg.MediaID, Spoiler: arg.Spoiler},
		},
		Content: arg.Content,
	}
	for _, t := range arg.Tags {
		a.Tags = append(a.Tags, &Tag{Name: t})
	}
	return a
}

// TransformDraft .
func TransformDraft(arg *ArgArticle) *Draft {
	return &Draft{
		Article: &Article{
			Meta: &Meta{
				ID:              arg.Aid,
				Author:          &Author{Mid: arg.Mid},
				Category:        &Category{ID: arg.Category},
				Title:           arg.Title,
				Summary:         arg.Summary,
				BannerURL:       arg.BannerURL,
				TemplateID:      arg.TemplateID,
				Reprint:         arg.Reprint,
				ImageURLs:       arg.ImageURLs,
				OriginImageURLs: arg.OriginImageURLs,
				Mtime:           xtime.Time(time.Now().Unix()),
				Dynamic:         arg.DynamicIntro,
				Media:           &Media{MediaID: arg.MediaID, Spoiler: arg.Spoiler},
			},
			Content: arg.Content,
		},
		Tags:   arg.Tags,
		ListID: arg.ListID,
	}
}

// ArgForce force update
type ArgForce struct {
	Force  bool
	RealIP string
}
