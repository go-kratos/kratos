package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/Dai0522/go-hash/murmur3"
)

// GetShareURL .
func (s *Service) GetShareURL(ctx context.Context, mid int64, device *bm.Device, req *v1.ShareRequest) (*v1.ShareResponse, error) {
	_, err := s.dao.VideoBase(ctx, mid, req.Svid)
	if err != nil {
		log.Warnw(ctx, "log", "get video base fail", "svid", req.Svid)
		return nil, err
	}

	token := s.dao.GetUserShareToken(ctx, mid)
	if token == "" {
		hash := murmur3.NewWithSeed(uint32(time.Now().Unix()))
		str := fmt.Sprintf("%d:%s", mid, buvid(device))
		token = toHex(hash.Murmur3_128([]byte(str)))
		s.dao.SetUserShareToken(ctx, mid, token)
	}

	var url, params []*v1.Tuple
	params = append(params, &v1.Tuple{
		Key: "mid",
		Val: strconv.Itoa(int(mid)),
	}, &v1.Tuple{
		Key: "svid",
		Val: strconv.Itoa(int(req.Svid)),
	}, &v1.Tuple{
		Key: "token",
		Val: token,
	})

	url = append(url, &v1.Tuple{
		Key: "1",
		Val: "https://bbq.bilibili.com/video/?id={svid}&token={token}",
	})
	url = append(url, &v1.Tuple{
		Key: "2",
		Val: "https://bbq.bilibili.com/user/?id={mid}&token={token}",
	})

	return &v1.ShareResponse{
		URL:    url,
		Params: params,
	}, nil
}

func toHex(b []byte) string {
	var res string
	pattern := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "e", "f", "g"}
	for _, v := range b {
		res += pattern[v&15]
		res += pattern[(v>>4)&15]
	}

	return res
}

// ShareCallback .
func (s *Service) ShareCallback(ctx context.Context, mid int64, device *bm.Device, args *v1.ShareCallbackRequest) (resp *v1.ShareCallbackResponse, err error) {
	// 增加分享数
	share := int32(0)
	if args.Svid != int64(0) {
		share, err = s.dao.IncrVideoStatisticsShare(ctx, args.Svid)
		if err != nil {
			log.Errorv(ctx, log.KV("log", err))
		}
	}

	resp = &v1.ShareCallbackResponse{
		ShareCount: share,
	}

	return
}
