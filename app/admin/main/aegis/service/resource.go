package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/resource"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// TxAddResource .
func (s *Service) TxAddResource(ormTx *gorm.DB, b *resource.Resource, r *resource.Result) (rid int64, err error) {
	if rid, err = s.gorm.TxAddResource(ormTx, b, r); err != nil {
		log.Error("s.dao.TxAddResource error(%v)", err)
	}
	return
}

// TxUpdateResource .
func (s *Service) TxUpdateResource(tx *gorm.DB, rid int64, res map[string]interface{}, rscRes *resource.Result) (err error) {
	if err = s.gorm.TxUpdateResult(tx, rid, res, rscRes); err != nil {
		return
	}

	return s.gorm.TxUpdateResource(tx, rid, res)
}

// CheckResource .
func (s *Service) CheckResource(c context.Context, opt *model.AddOption) (rid int64, err error) {
	r, err := s.gorm.ResByOID(c, opt.BusinessID, opt.OID)
	if err != nil || r == nil {
		return
	}
	if r != nil {
		rid = r.ID
	}
	return
}

// syncResource 同步到业务方
func (s *Service) syncResource(c context.Context, opt *model.SubmitOptions, mid int64, restoken *net.TokenPackage) (err error) {
	if opt.ExtraData == nil {
		opt.ExtraData = make(map[string]interface{})
	}
	opt.ExtraData["mid"] = mid

	if opt.Forbid != nil {
		if opt.Forbid.Duration > 0 {
			opt.Forbid.Duration = opt.Forbid.Duration * 60 * 60 * 24
		}
	}

	act, ropt, err := s.formatSubmitRequest(c, opt, restoken)
	if err != nil || act == nil {
		log.Error("syncResource formatSubmitRequest opt(%+v) restoken(%+v) error(%v)", opt, restoken, err)
		return
	}

	// TODO: 偶发错误(接口超时,业务方对应错误码等)，重试3次
	var codes = s.getConfig(c, opt.BusinessID, business.TypeTempCodes)
	for i := 0; i < 3; i++ {
		var code int
		code, err = s.http.SyncResource(c, act, ropt)
		if err == nil {
			break
		}
		if err != nil && code == 0 { //http错误不重试了，免得错误日志太多
			break
		}
		if len(codes) > 0 && !strings.Contains(","+codes+",", fmt.Sprintf(",%d,", code)) {
			return
		}
	}

	return err
}

// ResourceRes .
func (s *Service) ResourceRes(c context.Context, args *resource.Args) (res *resource.Res, err error) {
	if res, err = s.gorm.ResourceRes(c, args.RID); err != nil {
		log.Error("s.gorm.ResourceRes error(%v)", err)
	}
	return
}

// TxDelResource .
func (s *Service) TxDelResource(c context.Context, ormTx *gorm.DB, bizid int64, rids []int64) (delState int, err error) {
	delState = -10010
	var (
		ststr string
	)

	if ststr = s.getConfig(c, bizid, business.TypeDeleteState); len(ststr) != 0 {
		if delState, err = strconv.Atoi(ststr); err != nil {
			log.Error("TxDelResource strconv.Atoi error(%v)", err)
			return
		}
	}
	err = s.gorm.TxUpdateState(ormTx, rids, delState)
	return
}

func (s *Service) formatSubmitRequest(c context.Context, opt *model.SubmitOptions, restoken *net.TokenPackage) (act *model.Action, ropt map[string]interface{}, err error) {
	//1. 读取业务回调配置
	cfg := s.getConfig(c, opt.BusinessID, business.TypeCallback)
	if len(cfg) == 0 {
		err = errors.New("empty submit config")
		log.Error("FormatSubmitRequest(%d) error(%v)", opt.BusinessID, cfg)
		return
	}
	mapcfg := make(map[string]*model.Action)
	if err = json.Unmarshal([]byte(cfg), &mapcfg); err != nil {
		log.Error("FormatSubmitRequest(%d) json.Unmarshal error(%v)", opt.BusinessID, cfg)
		return
	}

	var ok bool
	//2. 根据流程结果判断，调用哪个回调
	if restoken.HitAudit {
		act, ok = mapcfg["audit"]
	} else {
		act, ok = mapcfg["noaudit"]
	}
	if !ok || act == nil {
		err = errors.New("empty submit config[noaudit]")
		log.Error("FormatSubmitRequest(%d) error(%v)", opt.BusinessID, cfg)
		return
	}

	ropt = make(map[string]interface{})

	ot := reflect.TypeOf(*opt)
	ov := reflect.ValueOf(*opt)
	for k, v := range act.Params {
		model.SubReflect(ot, ov, k, strings.Split(v.Value, "."), v.Default, ropt)
	}

	//3. 根据回调配置，读取需要的参数进行拼接
	for k, v := range restoken.Values {
		ropt[k] = v
	}

	return
}
