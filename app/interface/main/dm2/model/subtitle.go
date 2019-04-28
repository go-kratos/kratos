package model

import (
	"go-common/library/ecode"
)

// SubtitleLocation .
const (
	SubtitleLocationLeftBottom  = uint8(1) //左下角
	SubtitleLocationBottomMid   = uint8(2) //底部居中
	SubtitleLocationRightBottom = uint8(3) //右下角
	SubtitleLocationLeftUp      = uint8(7) //左上角
	SubtitleLocationUpMid       = uint8(8) //顶部居中
	SubtitleLocationRightUp     = uint8(9) //右上角

	SubtitleContentSizeLimit = 300
)

var (
	// SubtitleLocationMap .
	SubtitleLocationMap = map[uint8]struct{}{
		SubtitleLocationLeftBottom:  {},
		SubtitleLocationBottomMid:   {},
		SubtitleLocationRightBottom: {},
		SubtitleLocationLeftUp:      {},
		SubtitleLocationUpMid:       {},
		SubtitleLocationRightUp:     {},
	}
)

// SubtitleStatus .
type SubtitleStatus uint8

// SubtitleStatus
const (
	SubtitleStatusUnknown SubtitleStatus = iota
	SubtitleStatusDraft
	SubtitleStatusToAudit
	SubtitleStatusAuditBack
	SubtitleStatusRemove
	SubtitleStatusPublish
	SubtitleStatusCheckToAudit
	SubtitleStatusCheckPublish
	SubtitleStatusManagerBack
	SubtitleStatusManagerRemove
)

// UpperStatus .
type UpperStatus uint8

// UpperStatus
const (
	UpperStatusUnknow UpperStatus = iota
	UpperStatusUpper
)

// AuthorStatus .
type AuthorStatus uint8

// AuthorStatus
const (
	AuthorStatusUnknow AuthorStatus = iota
	AuthorStatusAuthor
)

// WaveFormStatus .
type WaveFormStatus uint8

//WaveFormStatus
const (
	WaveFormStatusWaitting WaveFormStatus = iota
	WaveFormStatusSuccess
	WaveFormStatusFailed
	WaveFormStatusError // this status need retry
)

// Subtitle .
type Subtitle struct {
	ID            int64          `json:"id"`
	Oid           int64          `json:"oid"`
	Type          int32          `json:"type"`
	Lan           uint8          `json:"lan"`
	Aid           int64          `json:"aid"`
	Mid           int64          `json:"mid"`
	AuthorID      int64          `json:"author_id"`
	UpMid         int64          `json:"up_mid"`
	IsSign        bool           `json:"is_sign"`
	IsLock        bool           `json:"is_lock"`
	Status        SubtitleStatus `json:"status"`
	CheckSum      string         `json:"-"`
	SubtitleURL   string         `json:"subtitle_url"`
	PubTime       int64          `json:"pub_time"`
	RejectComment string         `json:"reject_comment"`
	Mtime         int64          `json:"mtime"`
	Empty         bool           `json:"empty"`
}

// SubtitleShow .
type SubtitleShow struct {
	ID            int64          `json:"id"`
	Oid           int64          `json:"oid"`
	Type          int32          `json:"type"`
	Lan           string         `json:"lan"`
	LanDoc        string         `json:"lan_doc"`
	Mid           int64          `json:"mid"`
	Author        string         `json:"author"`
	Aid           int64          `json:"aid"`
	ArchiveName   string         `json:"archive_name"`
	IsSign        bool           `json:"is_sign"`
	IsLock        bool           `json:"is_lock"`
	Status        SubtitleStatus `json:"status"`
	SubtitleURL   string         `json:"subtitle_url"`
	RejectComment string         `json:"reject_comment"`
	AuthorStatus  AuthorStatus   `json:"author_status"` // 1:作者
	UpperStatus   UpperStatus    `json:"upper_status"`  // 1:up主
}

// SubtitlePub .
type SubtitlePub struct {
	Oid        int64 `json:"oid"`
	Type       int32 `json:"type"`
	Lan        uint8 `json:"lan"`
	SubtitleID int64 `json:"subtitle_id"`
	IsDelete   bool  `json:"is_delete"`
}

// VideoSubtitles .
type VideoSubtitles struct {
	AllowSubmit bool             `json:"allow_submit"`
	Lan         string           `json:"lan"`
	LanDoc      string           `json:"lan_doc"`
	Subtitles   []*VideoSubtitle `json:"subtitles"`
}

// VideoSubtitleCache .
type VideoSubtitleCache struct {
	VideoSubtitles []*VideoSubtitle `json:"video_subtitles"`
}

// VideoSubtitle .
type VideoSubtitle struct {
	ID          int64  `json:"id"`
	Lan         string `json:"lan"`
	LanDoc      string `json:"lan_doc"`
	IsLock      bool   `json:"is_lock"`
	AuthorMid   int64  `json:"author_mid,omitempty"`
	SubtitleURL string `json:"subtitle_url"`
}

// Language .
type Language struct {
	Lan       string       `json:"lan"`
	LanDoc    string       `json:"lan_doc"`
	Pub       *LanguagePub `json:"pub,omitempty"`
	Draft     *LanguageID  `json:"draft,omitempty"`
	Audit     *LanguageID  `json:"audit,omitempty"`
	AuditBack *LanguageID  `json:"audit_back,omitempty"`
}

// LanguagePub .
type LanguagePub struct {
	SubtitleID int64 `json:"subtitle_id"`
	IsLock     bool  `json:"is_lock"`
	IsPub      bool  `json:"is_pub"`
}

// LanguageID .
type LanguageID struct {
	SubtitleID int64 `json:"subtitle_id"`
}

// SubtitlePageResult .
type SubtitlePageResult struct {
	ID  int64 `json:"id"`
	Oid int64 `json:"oid"`
}

// CountSubtitleResult .
type CountSubtitleResult struct {
	Draft     int64
	ToAudit   int64
	AuditBack int64
	Publish   int64
}

// SearchSubtitleResult .
type SearchSubtitleResult struct {
	Page    *SearchPage           `json:"page"`
	Results []*SubtitlePageResult `json:"result"`
}

// SearchSubtitle .
type SearchSubtitle struct {
	ID            int64  `json:"id"`
	Oid           int64  `json:"oid"`
	Aid           int64  `json:"aid"`
	Type          int32  `json:"type"`
	ArchiveName   string `json:"archive_name"`
	VideoName     string `json:"video_name"`
	ArchivePic    string `json:"archive_pic"`
	AuthorID      int64  `json:"author_id"`
	Author        string `json:"author"`
	AuthorPic     string `json:"author_pic"`
	Lan           string `json:"lan"`
	LanDoc        string `json:"lan_doc"`
	Status        int32  `json:"status"`
	IsSign        bool   `json:"is_sign"`
	IsLock        bool   `json:"is_lock"`
	RejectComment string `json:"reject_comment"`
	Mtime         int64  `json:"mtime"`
}

// SearchSubtitleResponse .
type SearchSubtitleResponse struct {
	Page      *SearchPage       `json:"page"`
	Subtitles []*SearchSubtitle `json:"subtitles"`
}

// SearchSubtitleAuthorItem .
type SearchSubtitleAuthorItem struct {
	ID            int64  `json:"id"`
	Oid           int64  `json:"oid"`
	Aid           int64  `json:"aid"`
	Type          int32  `json:"type"`
	ArchiveName   string `json:"archive_name"`
	VideoName     string `json:"video_name"`
	ArchivePic    string `json:"archive_pic"`
	Lan           string `json:"lan"`
	LanDoc        string `json:"lan_doc"`
	Status        int32  `json:"status"`
	IsSign        bool   `json:"is_sign"`
	IsLock        bool   `json:"is_lock"`
	RejectComment string `json:"reject_comment"`
	Mtime         int64  `json:"mtime"`
}

// SearchSubtitleAuthor .
type SearchSubtitleAuthor struct {
	Page         *SearchPage                 `json:"page"`
	Subtitles    []*SearchSubtitleAuthorItem `json:"subtitles"`
	Total        int64                       `json:"total"`
	DraftCount   int64                       `json:"draft_count"`
	AuditCount   int64                       `json:"audit_count"`
	BackCount    int64                       `json:"back_count"`
	PublishCount int64                       `json:"publish_count"`
}

// SearchSubtitleAssit .
type SearchSubtitleAssit struct {
	Page         *SearchPage       `json:"page"`
	Subtitles    []*SearchSubtitle `json:"subtitles"`
	Total        int64             `json:"total"`
	AuditCount   int64             `json:"audit_count"`
	PublishCount int64             `json:"publish_count"`
}

// Subtitle state
const (
	AttrSubtitleClose = uint(1) // 关闭稿件字幕
)

// SubtitleSubject .
type SubtitleSubject struct {
	Aid   int64 `json:"aid"`
	Allow bool  `json:"allow"`
	Attr  int32 `json:"attr"`
	Lan   uint8 `json:"lan"`
	Empty bool  `json:"empty"`
}

// AttrVal return val of subtitle subject'attr
func (s *SubtitleSubject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subtitle subject'attr
func (s *SubtitleSubject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}

// SubtitleItem .
type SubtitleItem struct {
	From     float64 `json:"from"`
	To       float64 `json:"to"`
	Location uint8   `json:"location"`
	Content  string  `json:"content"`
}

// SubtitleBody .
type SubtitleBody struct {
	FontSize        float64         `json:"font_size,omitempty"`
	FontColor       string          `json:"font_color,omitempty"`
	BackgroundAlpha float64         `json:"background_alpha,omitempty"`
	BackgroundColor string          `json:"background_color,omitempty"`
	Stroke          string          `json:"Stroke,omitempty"`
	Bodys           []*SubtitleItem `json:"body"`
}

// CheckItem .
// err 兼容老接口error，等创作中心上线后去掉error返回
func (s *SubtitleBody) CheckItem(duration int64) (detectErrs []*SubtitleDetectError, err error) {
	var (
		maxDuration = float64(duration) / float64(1000)
	)
	maxDuration = maxDuration + 1 // 时间刻度上线兼容1
	for idx, item := range s.Bodys {
		if len(item.Content) > SubtitleContentSizeLimit {
			detectErrs = append(detectErrs, &SubtitleDetectError{
				Line:     int32(idx),
				ErrorMsg: ecode.SubtitleSizeLimit.Message(),
			})
			err = ecode.SubtitleSizeLimit
			continue
		}
		if _, ok := SubtitleLocationMap[item.Location]; !ok {
			detectErrs = append(detectErrs, &SubtitleDetectError{
				Line:     int32(idx),
				ErrorMsg: ecode.SubtitleLocationUnValid.Message(),
			})
			err = ecode.SubtitleSizeLimit
			continue
		}
		if item.From >= item.To {
			detectErrs = append(detectErrs, &SubtitleDetectError{
				Line:     int32(idx),
				ErrorMsg: ecode.SubtitleDuarionMustThanZero.Message(),
			})
			err = ecode.SubtitleSizeLimit
			continue
		}
		if item.From > maxDuration || item.To > maxDuration {
			detectErrs = append(detectErrs, &SubtitleDetectError{
				Line:     int32(idx),
				ErrorMsg: ecode.SubtitleVideoDurationOverFlow.Message(),
			})
			err = ecode.SubtitleSizeLimit
			continue
		}
	}
	return
}

// WaveForm .
type WaveForm struct {
	Oid         int64          `json:"oid"`
	Type        int32          `json:"type"`
	State       WaveFormStatus `json:"state"`
	WaveFromURL string         `json:"wave_form_url"`
	Mtime       int64          `json:"mtime"`
	Empty       bool
}

// WaveFormResp .
type WaveFormResp struct {
	State       WaveFormStatus `json:"state"`
	WaveFromURL string         `json:"wave_form_url"`
}

// SubtitleLans .
type SubtitleLans []*SubtitleLan

// SubtitleLan .
type SubtitleLan struct {
	Code     int64  `json:"-"`
	Lan      string `json:"lan"`
	DocZh    string `json:"doc_zh"`
	DocEn    string `json:"-"`
	IsDelete bool   `json:"-"`
}

// GetByLan .
func (ss SubtitleLans) GetByLan(lan string) (code int64) {
	for _, s := range ss {
		if s.Lan == lan {
			return s.Code
		}
	}
	return 0
}

// GetByID .
func (ss SubtitleLans) GetByID(lanID int64) (lan string, doc string) {
	for _, s := range ss {
		if s.Code == lanID {
			return s.Lan, s.DocZh
		}
	}
	return
}

// SubtitleCheckMsg .
type SubtitleCheckMsg struct {
	SubtitleID int64 `json:"subtitle_id"`
	Oid        int64 `json:"oid"`
}

// FilterCheckResp .
type FilterCheckResp struct {
	Hits map[string]string `json:"hits"`
}

// SubtitleDetectError .
type SubtitleDetectError struct {
	Line     int32  `json:"line"`
	ErrorMsg string `json:"error_msg"`
}
