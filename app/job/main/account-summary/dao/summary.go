package dao

import (
	"context"
	"strconv"

	"go-common/app/job/main/account-summary/model"

	"github.com/pkg/errors"
)

const (
	_SummaryTable = "ugc:AccountSum"
	_ColFamily    = "accountsum"
)

// Save is
func (d *Dao) Save(ctx context.Context, key string, data map[string][]byte) error {
	values := map[string]map[string][]byte{
		_ColFamily: data,
	}
	_, err := d.AccountSumHBase.PutStr(ctx, _SummaryTable, key, values)
	return err
}

// GetByKey is
func (d *Dao) GetByKey(ctx context.Context, key string) (*model.AccountSummary, error) {
	res, err := d.AccountSumHBase.GetStr(ctx, _SummaryTable, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sum := model.NewAccountSummary()
	for _, c := range res.Cells {
		v := string(c.Value)
		switch string(c.Qualifier) {
		case "birthday":
			sum.Birthday = v
		case "face":
			sum.Face = v
		case "mid":
			sum.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "name":
			sum.Name = v
		case "rank":
			sum.Rank, _ = strconv.ParseInt(v, 10, 64)
		case "sex":
			sum.Sex, _ = strconv.ParseInt(v, 10, 64)
		case "sign":
			sum.Sign = v
		case "official.role":
			sum.Official.Role, _ = strconv.ParseInt(v, 10, 64)
		case "official.mid":
			sum.Official.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "official.title":
			sum.Official.Title = v
		case "official.description":
			sum.Official.Description = v
		case "exp.mid":
			sum.Exp.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "exp.exp":
			sum.Exp.Exp, _ = strconv.ParseInt(v, 10, 64)
		case "relation.mid":
			sum.RelationStat.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "relation.follower":
			sum.RelationStat.Follower, _ = strconv.ParseInt(v, 10, 64)
		case "relation.following":
			sum.RelationStat.Following, _ = strconv.ParseInt(v, 10, 64)
		case "relation.black":
			sum.RelationStat.Black, _ = strconv.ParseInt(v, 10, 64)
		case "relation.whisper":
			sum.RelationStat.Whisper, _ = strconv.ParseInt(v, 10, 64)
		case "block.mid":
			sum.Block.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "block.block_status":
			sum.Block.BlockStatus, _ = strconv.ParseInt(v, 10, 64)
		case "block.start_time":
			sum.Block.StartTime = v
		case "block.end_time":
			sum.Block.EndTime = v
		case "passport.mid":
			sum.Passport.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "passport.tel_status":
			sum.Passport.TelStatus, _ = strconv.ParseInt(v, 10, 64)
		case "passport.country_id":
			sum.Passport.CountryID, _ = strconv.ParseInt(v, 10, 64)
		case "passport.join_ip":
			sum.Passport.JoinIP = v
		case "passport.join_time":
			sum.Passport.JoinTime = v
		case "passport.email_suffix":
			sum.Passport.EmailSuffix = v
		case "passport.origin_type":
			sum.Passport.OriginType, _ = strconv.ParseInt(v, 10, 64)
		case "passport.reg_type":
			sum.Passport.RegType, _ = strconv.ParseInt(v, 10, 64)
		}
	}
	return sum, nil
}
