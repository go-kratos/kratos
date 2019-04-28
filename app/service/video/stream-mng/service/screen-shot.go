package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	url2 "net/url"
	"time"
)

// GetSingleScreeShot 得到一个房间一个时间段内截图
func (s *Service) GetSingleScreeShot(c context.Context, rid int64, startTime int64, endTime int64, channel string) ([]string, error) {
	resp := []string{}
	// todo 测试房间映射

	streamName, origin, err := s.dao.OriginUpStreamInfo(c, rid)
	if err != nil {
		return nil, err
	}

	for i := startTime; i <= endTime; i += 15 {
		resp = append(resp, s.getCapture(streamName, i, common.BitwiseMapSrc[origin], channel))
	}
	return resp, nil
}

// GetMultiScreenShot 多个房间一个时间戳截图
func (s *Service) GetMultiScreenShot(c context.Context, rids []int64, ts int64, channel string) (map[int64]string, error) {
	infos, err := s.dao.MultiStreamInfo(c, rids)

	res := map[int64]string{}
	if err != nil {
		return res, err
	}

	if infos == nil {
		return res, fmt.Errorf("查询不到数据")
	}

	for id, item := range infos {
		for _, v := range item.List {
			if v.Type == 1 {
				var or int64
				if v.Origin != 0 {
					or = v.Origin
				} else {
					or = v.DefaultUpStream
				}
				res[id] = s.getCapture(v.StreamName, ts, common.BitwiseMapSrc[or], channel)
				break
			}
		}
	}

	return res, nil
}

// GetOriginScreenShotPic 多个房间一个时间戳原始截图
func (s *Service) GetOriginScreenShotPic(c context.Context, rids []int64, ts int64) (map[int64]string, error) {
	infos, err := s.dao.MultiStreamInfo(c, rids)

	res := map[int64]string{}
	if err != nil {
		return res, err
	}

	if infos == nil {
		return res, fmt.Errorf("查询不到数据")
	}

	for id, item := range infos {
		for _, v := range item.List {
			if v.Type == 1 {
				var or int64
				if v.Origin != 0 {
					or = v.Origin
				} else {
					or = v.DefaultUpStream
				}
				res[id] = s.getOriginCapture(v.StreamName, ts, common.BitwiseMapSrc[or])
				break
			}
		}
	}

	return res, nil
}

// getOriginCapture 原始图片地址
func (s *Service) getOriginCapture(streamName string, t int64, src int8) string {
	tu := time.Unix(t, 0)
	dataStr := tu.Format("200601021504")
	secondStr := fmt.Sprintf("%02d", (tu.Second()/15)*15)

	ts := fmt.Sprintf("%s%s", dataStr, secondStr)
	url := ""

	switch src {
	case common.WSSrc:
		url = "http://live-jk-img.hdslb.net/original_live-%s--%s.jpg"
		url = fmt.Sprintf(url, streamName, ts)
	case common.TXYSrc:
		tsd := time.Unix(t, 0).Format("2006-01-02")
		url = fmt.Sprintf("http://bilibilitest-1252693259.cosgz.myqcloud.com/%s-%s/%s.png", tsd, streamName, ts)
	case common.JSSrc:
		url = fmt.Sprintf("https://ks3-cn-beijing.ksyun.com/live-image/record/live-js/%s/%s.png", streamName, ts)
	case common.QNSrc:
		url = fmt.Sprintf("http://qn.static.acgvideo.com/%s/%s.jpg", streamName, ts)
	case common.BVCSrc:
		path := fmt.Sprintf("/liveshotraw/%s_%s.png", streamName, ts)

		params := make(url2.Values)
		params.Set("deadline", fmt.Sprintf("%d", time.Now().Unix()+5184000))
		url = fmt.Sprintf("http://upos-sz.acgvideo.com%s", s.getBVCSign(path, params))
	default:
	}

	return url
}

// getCapture
func (s *Service) getCapture(streamName string, t int64, src int8, channel string) string {
	// ts date('YmdHi', $iTS) . sprintf('%02d', (int)(date('s', $iTS) / 15) * 15);
	tu := time.Unix(t, 0)
	dataStr := tu.Format("200601021504")
	secondStr := fmt.Sprintf("%02d", (tu.Second()/15)*15)

	ts := fmt.Sprintf("%s%s", dataStr, secondStr)
	tsd := time.Unix(t, 0).Format("2006-01-02")

	url := ""

	switch src {
	case common.WSSrc:
		url = "http://live-jk-img.hdslb.net/live-%s--%s.jpg"
		if channel == "checkup_yellow" {
			url = "http://live.pic.bilibili.8686c.com/live-%s--%s.jpg"
		}
		url = fmt.Sprintf(url, streamName, ts)
	case common.TXYSrc:
		url = fmt.Sprintf("http://bilibilitest-1252693259.cosgz.myqcloud.com/%s-%s/%s.jpg", tsd, streamName, ts)
	case common.JSSrc:
		url = fmt.Sprintf("https://ks3-cn-beijing.ksyun.com/live-image/record/live-js/%s/%s.jpg", streamName, ts)
	case common.QNSrc:
		url = fmt.Sprintf("http://qn.static.acgvideo.com/%s/%s.jpg?imageView2/0/h/293/format/jpg", streamName, ts)
	case common.BVCSrc:
		path := fmt.Sprintf("/liveshot/%s_%s.jpg", streamName, ts)

		params := make(url2.Values)
		params.Set("deadline", fmt.Sprintf("%d", time.Now().Unix()+5184000))
		url = fmt.Sprintf("http://upos-sz-office.acgvideo.com%s", s.getBVCSign(path, params))
	default:
	}

	return url
}

// getBVCSign bvc签名
func (s *Service) getBVCSign(path string, params url2.Values) string {
	upsigSecret := "20170607920cbd5211831ce2a97066a8b544fa7b"
	toSign := fmt.Sprintf("%s?%s", path, params.Encode())

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", toSign, upsigSecret)))
	cipherStr := h.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)

	return fmt.Sprintf("%s&upsig=%s", toSign, md5Str)
}
