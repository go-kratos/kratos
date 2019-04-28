package model

import (
	"go-common/app/interface/main/creative/model/data"
	"go-common/library/time"
)

// const
const (
	// ReprintForbid 禁止转载.
	ReprintForbid = int8(0)
	// ReprintAllow 允许规范转载.
	ReprintAllow = int8(1)
	// NoImage 无图.
	///NoImage = int8(1)
	// HeadImage 头图.
	HeadImage = int8(4)
)

var (
	_reprint = map[int8]int8{
		ReprintForbid: ReprintForbid,
		ReprintAllow:  ReprintAllow,
	}
	_tid = map[int8]int8{
		TemplateText:         0,
		TemplateSingleImg:    1,
		TemplateMultiImg:     3,
		TemplateSingleBigImg: 1,
	}
)

// InReprints check reprint in all reprints.
func InReprints(rp int8) (ok bool) {
	_, ok = _reprint[rp]
	return
}

// InTemplateID check tid in all tids.
func InTemplateID(tid int8) (ok bool) {
	_, ok = _tid[tid]
	return
}

// ValidTemplate checks template id & images count.
func ValidTemplate(tid int32, imgs []string) bool {
	var images []string
	for _, image := range imgs {
		if image != "" {
			images = append(images, image)
		}
	}
	return len(images) == int(_tid[int8(tid)])
}

// UpStat for bigdata article up stat
type UpStat struct {
	View      int64 `json:"view"`
	Reply     int64 `json:"reply"`
	Like      int64 `json:"like"`
	Coin      int64 `json:"coin"`
	Fav       int64 `json:"fav"`
	Share     int64 `json:"share"`
	PreView   int64 `json:"-"`
	PreReply  int64 `json:"-"`
	PreLike   int64 `json:"-"`
	PreCoin   int64 `json:"-"`
	PreFav    int64 `json:"-"`
	PreShare  int64 `json:"-"`
	IncrView  int64 `json:"incr_view"`
	IncrReply int64 `json:"incr_reply"`
	IncrLike  int64 `json:"incr_like"`
	IncrCoin  int64 `json:"incr_coin"`
	IncrFav   int64 `json:"incr_fav"`
	IncrShare int64 `json:"incr_share"`
}

// ThirtyDayArticle for article 30 days data.
type ThirtyDayArticle struct {
	Category  string            `json:"category"`
	ThirtyDay []*data.ThirtyDay `json:"thirty_day"`
}

// ArtParam param  for article info input.
type ArtParam struct {
	AID             int64    `json:"aid"`
	MID             int64    `json:"mid"`
	Category        int64    `json:"category"`
	State           int32    `json:"state"`
	Reprint         int32    `json:"reprint"`
	TemplateID      int32    `json:"tid"`
	Title           string   `json:"title"`
	BannerURL       string   `json:"banner_url"`
	Content         string   `json:"content"`
	Summary         string   `json:"summary"`
	Tags            string   `json:"tags"`
	ImageURLs       []string `json:"image_urls"`
	OriginImageURLs []string `json:"origin_image_urls"`
	RealIP          string   `json:"-"`
	Action          int      `json:"action"`
	Words           int64    `json:"words"`
	DynamicIntro    string   `json:"dynamic_intro"`
	ActivityID      int64    `json:"activity_id"`
	ListID          int64    `json:"list_id"`
	MediaID         int64    `json:"media_id"`
	Spoiler         int32    `json:"spoiler"`
}

// CreativeMeta  article detail.
type CreativeMeta struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Summary         string    `json:"summary"`
	BannerURL       string    `json:"banner_url"`
	Reason          string    `json:"reason"`
	TemplateID      int32     `json:"template_id"`
	State           int32     `json:"state"`
	Reprint         int32     `json:"reprint"`
	ImageURLs       []string  `json:"image_urls"`
	OriginImageURLs []string  `json:"origin_image_urls"`
	Tags            []string  `json:"tags"`
	Category        *Category `json:"category"`
	Author          *Author   `json:"author"`
	Stats           *Stats    `json:"stats"`
	PTime           time.Time `json:"publish_time"`
	CTime           time.Time `json:"ctime"`
	MTime           time.Time `json:"mtime"`
	ViewURL         string    `json:"view_url"`
	EditURL         string    `json:"edit_url"`
	IsPreview       int       `json:"is_preview"`
	DynamicIntro    string    `json:"dynamic_intro"`
	List            *List     `json:"list"`
	MediaID         int64     `json:"media_id"`
	Spoiler         int32     `json:"spoiler"`
	EditTimes       int       `json:"edit_times"`
	PreViewURL      string    `json:"pre_view_url"`
}

// CreativeArtList article for list.
type CreativeArtList struct {
	Articles []*CreativeMeta   `json:"articles"`
	Type     *CreationArtsType `json:"type"`
	Page     *ArtPage          `json:"page"`
}

// CreativeDraftList draft list.
type CreativeDraftList struct {
	Drafts   []*CreativeMeta `json:"drafts"`
	Page     *ArtPage        `json:"page"`
	DraftURL string          `json:"draft_url"`
}

// ExtMsg .
type ExtMsg struct {
	Tags []*Tag `json:"tags"`
}
