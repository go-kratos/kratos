package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/job/main/account-summary/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	queryByMidURI = "/intranet/acc/queryByMid"

	// _AsoAccountByMid          = `SELECT mid,email FROM aso_account WHERE mid=?`
	// _AsoAccountInfoByMid      = `SELECT mid,join_ip,join_time FROM aso_account_info%d WHERE mid=?`
	_AsoAccountRegOriginByMid = `SELECT mid,origintype,regtype FROM aso_account_reg_origin_%d WHERE mid=?`
	// _AllCountryID  = `SELECT id,code FROM aso_country_code`
)

// PassportProfile is
func (d *Dao) PassportProfile(ctx context.Context, mid int64) (*model.PassportProfile, error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))

	var res struct {
		Code int                    `json:"code"`
		Data *model.PassportProfile `json:"data"`
	}
	if err := d.httpClient.Get(ctx, d.c.Host.Passport+queryByMidURI, "", params, &res); err != nil {
		log.Error("Failed to query by mid: %+v: %+v", params, err)
		return nil, err
	}
	if res.Code != 0 {
		log.Error("Failed to query by mid with code: %+v: %d", params, res.Code)
		return nil, ecode.Int(res.Code)
	}
	return res.Data, nil
}

// AsoAccountRegOrigin is
func (d *Dao) AsoAccountRegOrigin(ctx context.Context, mid int64) (*model.AsoAccountRegOrigin, error) {
	row := d.PassportDB.QueryRow(ctx, fmt.Sprintf(_AsoAccountRegOriginByMid, mid%20), mid)
	origin := new(model.AsoAccountRegOrigin)
	if err := row.Scan(&origin.Mid, &origin.OriginType, &origin.RegType); err != nil {
		return nil, err
	}
	return origin, nil
}
