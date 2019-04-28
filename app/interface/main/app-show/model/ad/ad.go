package ad

import (
	"strconv"

	"go-common/app/interface/main/app-show/model/banner"
	"go-common/library/log"
)

type ADRequest struct {
	RequestID string                         `json:"request_id"`
	ADIndexs  map[string]map[string]*ADIndex `json:"ads_info"`
}

type ADIndex struct {
	Index  int     `json:"index"`
	Info   *ADInfo `json:"ad_info"`
	IsAd   bool    `json:"is_ad"`
	CmMark int     `json:"cm_mark"`
}

type ADInfo struct {
	CreativeID      int `json:"creative_id"`
	CreativeType    int `json:"creative_type"`
	CreativeContent struct {
		Title    string `json:"title"`
		Desc     string `json:"description"`
		VideoID  int64  `json:"video_id"`
		UserName string `json:"username"`
		ImageURL string `json:"image_url"`
		ImageMD5 string `json:"image_md5"`
		LogURL   string `json:"log_url"`
		LogMD5   string `json:"log_md5"`
		URL      string `json:"url"`
		ClickURL string `json:"click_url"`
		ShowURL  string `json:"show_url"`
	} `json:"creative_content"`
	AdCb string `json:"ad_cb"`
}

// RecoverBanner
func (adr *ADRequest) ConvertBanner(ip, mobiApp string, build int) (banners map[int]map[int]*banner.Banner) {
	banners = map[int]map[int]*banner.Banner{}
	for resIdStr, sAdis := range adr.ADIndexs {
		resId, _ := strconv.Atoi(resIdStr)
		if len(sAdis) == 0 {
			log.Info("mobi_app:%v-build:%v-resource:%v-is_ad_loc:%v", mobiApp, build, resId, false)
			continue
		}
		for sidStr, adi := range sAdis {
			sid, _ := strconv.Atoi(sidStr)
			var bnnr = &banner.Banner{
				IsAdLoc:     true,
				IsAd:        adi.IsAd,
				IsAdReplace: false,
				CmMark:      adi.CmMark,
				Rank:        adi.Index,
				SrcId:       sid,
				RequestId:   adr.RequestID,
				ClientIp:    ip,
			}
			if adInfo := adi.Info; adInfo != nil {
				bnnr.IsAdReplace = true
				bnnr.CreativeId = adInfo.CreativeID
				bnnr.AdCb = adInfo.AdCb
				bnnr.ShowUrl = adInfo.CreativeContent.ShowURL
				bnnr.ClickUrl = adInfo.CreativeContent.ClickURL
				bnnr.Title = adInfo.CreativeContent.Title
				bnnr.Image = adInfo.CreativeContent.ImageURL
				bnnr.Hash = adInfo.CreativeContent.ImageMD5
				bnnr.URI = adInfo.CreativeContent.URL
				bnnr.Channel = "*"
			}
			if _, ok := banners[resId]; ok {
				banners[resId][bnnr.Rank] = bnnr
			} else {
				banners[resId] = map[int]*banner.Banner{
					bnnr.Rank: bnnr,
				}
			}
			log.Info("mobi_app:%v-build:%v-source:%v-resource:%v-is_ad_loc:%v", mobiApp, build, sid, resId, true)
		}
	}
	return
}
