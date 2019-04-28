package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/growup/model"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// AccountInfos get account infos
func (d *Dao) AccountInfos(c context.Context, mids []int64) (infos map[int64]*model.ActUpInfo, err error) {
	if len(mids) == 0 {
		return
	}
	infos = make(map[int64]*model.ActUpInfo)
	results := new(model.AccountInfosResult)
	uv := url.Values{}
	uv.Set("mids", xstr.JoinInts(mids))
	if err = d.httpRead.Get(c, d.c.Host.AccountURI, "", uv, results); err != nil {
		return
	}
	if results.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(results.Code), fmt.Sprintf("search account failed: %s?%s", d.c.Host.AccountURI, uv.Get("mids")))
		return
	}
	for mid, account := range results.Data {
		infos[mid] = &model.ActUpInfo{Nickname: account.Name, Face: account.Face}
	}
	return
}

// UpBusinessInfos get business infos
func (d *Dao) UpBusinessInfos(c context.Context, mid int64) (identify *model.UpIdentify, err error) {
	identify = new(model.UpIdentify)
	results := new(model.UperInfosResult)
	uv := url.Values{}
	uv.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.httpRead.Get(c, d.c.Host.UperURI, "", uv, results); err != nil {
		return
	}
	if results.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(results.Code), fmt.Sprintf("search uper failed: %s?%s", d.c.Host.UperURI, uv.Get("mid")))
		return
	}
	identify = results.Data["identify"]
	return
}

// Card get account.
func (d *Dao) Card(c context.Context, mid int64) (res *account.Card, err error) {
	var arg = &account.ArgMid{
		Mid: mid,
	}
	if res, err = d.acc.Card3(c, arg); err != nil {
		log.Error("d.acc.Card3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	return
}

// ProfileWithStat get account.
func (d *Dao) ProfileWithStat(c context.Context, mid int64) (res *account.ProfileStat, err error) {
	var arg = &account.ArgMid{
		Mid: mid,
	}
	if res, err = d.acc.ProfileWithStat3(c, arg); err != nil {
		log.Error("d.acc.ProfileWithStat3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	return
}
