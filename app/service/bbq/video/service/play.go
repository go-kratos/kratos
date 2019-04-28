package service

import (
	"context"
	"go-common/app/service/bbq/common/db/bbq"
	"go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/model/grpc"
	"go-common/library/log"
	"time"
)

// PlayInfo 批量获取playurl(相对地址方法)
func (s *Service) PlayInfo(c context.Context, args *v1.PlayInfoRequest) (res *v1.PlayInfoResponse, err error) {
	var (
		relAddr []string
		bvcUrls map[string]*grpc.VideoKeyItem
		bvcKeys map[int64][]*bbq.VideoBvc
		svids   = args.Svid
	)
	res = new(v1.PlayInfoResponse)
	bvcKeys, err = s.dao.RawSVBvcKey(c, svids)
	if err != nil {
		log.Error("s.dao.RawSVBvcKey err[%v]", err)
	}
	res.List = make(map[int64]*v1.PlayInfo)
	for id, keys := range bvcKeys {
		res.List[id] = &v1.PlayInfo{
			Svid: id,
		}
		for k, v := range keys {
			if k == 0 {
				res.List[id].Quality = int64(v.CodeRate)
			}
			fi := &v1.FileInfo{
				Timelength: v.Duration,
				Filesize:   v.FileSize,
				Path:       v.Path,
			}
			res.List[id].FileInfo = append(res.List[id].FileInfo, fi)
			res.List[id].SupportQuality = append(res.List[id].SupportQuality, int64(v.CodeRate))
			relAddr = append(relAddr, v.Path)
		}
	}
	bvcUrls, err = s.dao.RelPlayURLs(c, relAddr)
	if err != nil {
		log.Error("s.dao.RelPlayURLs err[%v]", err)
		return
	}
	//拼装playurl
	for _, svid := range svids {
		if play, ok := res.List[svid]; ok {
			for fk, f := range play.FileInfo {
				if urls, ok := bvcUrls[f.Path]; ok {
					res.List[svid].ExpireTime = int64(urls.Etime)
					res.List[svid].CurrentTime = time.Now().Unix()
					for _, u := range urls.URL {
						if res.List[svid].FileInfo[fk].Url == "" {
							res.List[svid].FileInfo[fk].Url = u
							if res.List[svid].Url == "" {
								res.List[svid].Url = u
							}
							continue
						}
						if res.List[svid].FileInfo[fk].UrlBc == "" {
							res.List[svid].FileInfo[fk].UrlBc = u
							break
						}
					}
				} else {
					delete(res.List, svid)
					break
				}
			}
			res.List[svid] = play
		}
	}
	return
}
