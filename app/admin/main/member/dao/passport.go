package dao

import (
	"context"
	"go-common/app/admin/main/member/model"
	"net/url"
	"strconv"
	"sync"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const (
	_updateUname = "/intranet/acc/updateUname"
	_queryByMids = "/intranet/acc/queryByMids"
)

// UpdateUname is.
func (d *Dao) UpdateUname(ctx context.Context, mid int64, name string) error {
	ip := metadata.String(ctx, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("uname", name)

	var res struct {
		Code int `json:"code"`
	}

	if err := d.passportClient.Post(ctx, d.upUnameURL, ip, params, &res); err != nil {
		return err
	}

	if res.Code != 0 {
		log.Error("Failed to update uname(%+v) code(%+v)", params, res.Code)
		return parsePassportEcode(res.Code)
	}
	return nil
}

func parsePassportEcode(pecode int) error {
	switch pecode {
	case -618:
		return ecode.UpdateUnameRepeated
	case -617:
		return ecode.UpdateUnameHadLocked
	case -601:
		return ecode.UpdateUnameFormat
	}
	log.Error("Unrecognized passport ecode: %d", pecode)
	return ecode.Int(pecode)
}

// PassportQueryByMids is.
func (d *Dao) PassportQueryByMids(ctx context.Context, mids []int64) (map[int64]*model.PassportQueryByMidResult, error) {
	ip := metadata.String(ctx, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mids", xstr.JoinInts(mids))
	var res struct {
		Code  int                                        `json:"code"`
		Cards map[string]*model.PassportQueryByMidResult `json:"cards"`
	}

	if err := d.passportClient.Get(ctx, d.queryByMidsURL, ip, params, &res); err != nil {
		return nil, err
	}

	if res.Code != 0 {
		log.Error("Failed to QueryByMid(%+v) code(%+v)", params, res.Code)
		return nil, ecode.Int(res.Code)
	}

	result := make(map[int64]*model.PassportQueryByMidResult, len(res.Cards))
	for _, card := range res.Cards {
		result[card.Mid] = card
	}

	return result, nil
}

// PassportQueryByMidsChunked is
func (d *Dao) PassportQueryByMidsChunked(ctx context.Context, mids []int64, chunkSize int) (map[int64]*model.PassportQueryByMidResult, error) {
	chunkedMids := func() [][]int64 {
		divided := make([][]int64, 0)
		for i := 0; i < len(mids); i += chunkSize {
			end := i + chunkSize
			if end > len(mids) {
				end = len(mids)
			}
			divided = append(divided, mids[i:end])
		}
		return divided
	}()

	lock := sync.Mutex{}
	result := make(map[int64]*model.PassportQueryByMidResult, len(mids))
	wg := sync.WaitGroup{}
	for _, chunk := range chunkedMids {
		wg.Add(1)
		go func(chunk []int64) {
			defer wg.Done()
			res, err := d.PassportQueryByMids(ctx, chunk)
			if err != nil {
				log.Error("Failed to get passport query by mids: %+v: %+v", chunk, err)
				return
			}
			lock.Lock()
			for k, v := range res {
				result[k] = v
			}
			lock.Unlock()
		}(chunk)
	}
	wg.Wait()

	return result, nil
}
