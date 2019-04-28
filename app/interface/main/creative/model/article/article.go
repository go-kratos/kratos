package article

import (
	model "go-common/app/interface/openplatform/article/model"
	"go-common/library/time"
)

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
}

// Meta  article detail.
type Meta struct {
	ID              int64           `json:"id"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	Summary         string          `json:"summary"`
	BannerURL       string          `json:"banner_url"`
	Reason          string          `json:"reason"`
	TemplateID      int32           `json:"template_id"`
	State           int32           `json:"state"`
	Reprint         int32           `json:"reprint"`
	ImageURLs       []string        `json:"image_urls"`
	OriginImageURLs []string        `json:"origin_image_urls"`
	Tags            []string        `json:"tags"`
	Category        *model.Category `json:"category"`
	Author          *model.Author   `json:"author"`
	Stats           *model.Stats    `json:"stats"`
	PTime           time.Time       `json:"publish_time"`
	CTime           time.Time       `json:"ctime"`
	MTime           time.Time       `json:"mtime"`
	ViewURL         string          `json:"view_url"`
	EditURL         string          `json:"edit_url"`
	IsPreview       int             `json:"is_preview"`
	DynamicIntro    string          `json:"dynamic_intro"`
}

// ArtList article for list.
type ArtList struct {
	Articles []*Meta                 `json:"articles"`
	Type     *model.CreationArtsType `json:"type"`
	Page     *model.ArtPage          `json:"page"`
}

// DraftList draft list.
type DraftList struct {
	Drafts   []*Meta        `json:"drafts"`
	Page     *model.ArtPage `json:"page"`
	DraftURL string         `json:"draft_url"`
}
