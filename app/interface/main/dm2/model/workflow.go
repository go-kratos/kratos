package model

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	_subtitleReportTagReasonID = 5
)

// WorkFlowTagListResp .
type WorkFlowTagListResp struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []*WorkFlowTag `json:"data"`
}

// WorkFlowTag .
type WorkFlowTag struct {
	Bid   int64  `json:"bid"`
	TagID int64  `json:"tag_id"`
	Rid   int64  `json:"rid"`
	Name  string `json:"name"`
}

// CommonResponse .
type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WorkFlowAppealAddReq .
type WorkFlowAppealAddReq struct {
	Business       int64                   // 14
	LanCode        int64                   // 语言code
	Rid            int64                   // workflow rid
	SubtitleID     int64                   // 字幕id
	Score          int32                   // 举报人得分
	Tid            int64                   // workflow tag id
	Oid            int64                   // 视频cid
	Aid            int64                   // 稿件id
	Mid            int64                   // 举报人mid
	BusinessTypeID int32                   // 分区id
	BusinessTitle  string                  // 举报内容
	BusinessMid    int64                   // 被举报人mid
	Description    string                  // 投诉的具体描述
	Extra          *WorkFlowAppealAddExtra // 附加信息
}

// WorkFlowAppealAddExtra .
type WorkFlowAppealAddExtra struct {
	SubtitleStatus int64  `json:"subtitle_status"`
	SubtitleURL    string `json:"subtitle_url"`
	ArchiveName    string `json:"arcvhive_name"`
}

// Params .
func (w *WorkFlowAppealAddReq) Params() (params url.Values) {
	var (
		err error
		bs  []byte
	)
	params = url.Values{}
	params.Set("business", strconv.FormatInt(w.Business, 10))
	params.Set("fid", strconv.FormatInt(w.LanCode, 10))
	params.Set("rid", strconv.FormatInt(w.Rid, 10))
	params.Set("eid", strconv.FormatInt(w.SubtitleID, 10))
	params.Set("score", strconv.Itoa(int(w.Score)))
	params.Set("tid", strconv.FormatInt(w.Tid, 10))
	params.Set("oid", strconv.FormatInt(w.Oid, 10))
	params.Set("aid", strconv.FormatInt(w.Aid, 10))
	params.Set("mid", strconv.FormatInt(w.Mid, 10))
	if w.Tid == _subtitleReportTagReasonID {
		params.Set("description", w.Description)
	}
	params.Set("business_typeid", strconv.Itoa(int(w.BusinessTypeID)))
	params.Set("business_title", w.BusinessTitle)
	params.Set("business_mid", strconv.FormatInt(w.BusinessMid, 10))
	if w.Extra != nil {
		if bs, err = json.Marshal(w.Extra); err == nil {
			params.Set("business_extra", string(bs))
		}
	}
	return
}
