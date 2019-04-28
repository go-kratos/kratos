package ad

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-view/model"
)

type Ad struct {
	RequestID  string                       `json:"request_id,omitempty"`
	AdsControl json.RawMessage              `json:"ads_control,omitempty"`
	AdsInfo    map[int64]map[int64]*AdsInfo `json:"ads_info,omitempty"`
	ClientIP   string                       `json:"-"`
}

type AdsInfo struct {
	Index     int     `json:"index,omitempty"`
	AdInfo    *AdInfo `json:"ad_info,omitempty"`
	IsAd      bool    `json:"is_ad,omitempty"`
	CmMark    int     `json:"cm_mark,omitempty"`
	CardIndex int     `json:"card_index,omitempty"`
}

type AdInfo struct {
	CreativeID      int64 `json:"creative_id,omitempty"`
	CreativeType    int64 `json:"creative_type,omitempty"`
	CreativeContent *struct {
		Title       string `json:"title,omitempty"`
		Desc        string `json:"description,omitempty"`
		ButtonTitle string `json:"button_title,omitempty"`
		VideoID     int64  `json:"video_id,omitempty"`
		UserName    string `json:"username,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
		ImageMD5    string `json:"image_md5,omitempty"`
		LogURL      string `json:"log_url,omitempty"`
		LogMD5      string `json:"log_md5,omitempty"`
		URL         string `json:"url,omitempty"`
		ClickURL    string `json:"click_url,omitempty"`
		ShowURL     string `json:"show_url,omitempty"`
	} `json:"creative_content,omitempty"`
	AdCb      string          `json:"ad_cb,omitempty"`
	CardType  int             `json:"card_type,omitempty"`
	Extra     json.RawMessage `json:"extra,omitempty"`
	Resource  int64           `json:"-"`
	Source    int64           `json:"-"`
	RequestID string          `json:"-"`
	IsAd      bool            `json:"-"`
	CmMark    int             `json:"-"`
	Index     int             `json:"-"`
	IsAdLoc   bool            `json:"-"`
	CardIndex int             `json:"-"`
	ClientIP  string          `json:"-"`
	// ad
	URI     string `json:"-"`
	Param   string `json:"-"`
	Goto    string `json:"-"`
	View    int    `json:"-"`
	Danmaku int    `json:"-"`
}

func (ad *Ad) Convert(resource int64) (ads []*AdInfo, aids []int64) {
	if ad == nil {
		return
	}
	if adsInfo, ok := ad.AdsInfo[resource]; ok {
		ads = make([]*AdInfo, 0, len(adsInfo))
		for source, info := range adsInfo {
			var adInfo *AdInfo
			if info != nil {
				if info.AdInfo != nil {
					adInfo = info.AdInfo
					adInfo.RequestID = ad.RequestID
					adInfo.Resource = resource
					adInfo.Source = source
					adInfo.IsAd = info.IsAd
					adInfo.IsAdLoc = true
					adInfo.CmMark = info.CmMark
					adInfo.Index = info.Index
					adInfo.CardIndex = info.CardIndex
					adInfo.ClientIP = ad.ClientIP
					// http://info.bilibili.co/pages/viewpage.action?pageId=6227100
					switch adInfo.CardType {
					case 6:
						adInfo.Goto = model.GotoAv
						if adInfo.CreativeContent != nil {
							adInfo.Param = strconv.FormatInt(int64(adInfo.CreativeContent.VideoID), 10)
							if adInfo.CreativeContent.VideoID > 0 {
								aids = append(aids, adInfo.CreativeContent.VideoID)
							}
						} else {
							adInfo.Param = "0"
						}
						adInfo.URI = model.FillURI(adInfo.Goto, adInfo.Param, nil)
					default:
						adInfo.Goto = model.GotoWeb
						if adInfo.CreativeContent != nil {
							adInfo.Param = adInfo.CreativeContent.URL
						}
						adInfo.URI = model.FillURI(adInfo.Goto, adInfo.Param, nil)
					}
				} else {
					adInfo = &AdInfo{RequestID: ad.RequestID, Resource: resource, Source: source, IsAdLoc: true, IsAd: info.IsAd, CmMark: info.CmMark, Index: info.Index, CardIndex: info.CardIndex, ClientIP: ad.ClientIP}
				}
			}
			if adInfo != nil {
				ads = append(ads, adInfo)
			}
		}
	}
	return
}

type AdInfos []*AdInfo

func (a AdInfos) Len() int           { return len(a) }
func (a AdInfos) Less(i, j int) bool { return int64(a[i].Index) < int64(a[j].Index) }
func (a AdInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
