package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

//GetTokenList .
func (s *Service) GetTokenList(c context.Context, pm *net.ListTokenParam) (result *net.ListTokenRes, err error) {
	var (
		unames map[int64]string
	)
	if result, err = s.gorm.TokenListWithPager(c, pm); err != nil {
		return
	}
	if len(result.Result) == 0 {
		return
	}

	//get username
	uid := []int64{}
	for _, item := range result.Result {
		if item.UID > 0 {
			uid = append(uid, item.UID)
		}
	}
	if len(uid) == 0 {
		return
	}

	if unames, err = s.http.GetUnames(c, uid); err != nil || len(unames) == 0 {
		log.Error("GetTokenList s.http.GetUnames error(%v) or empty uid(%+v)", err, uid)
		err = nil
	}
	for _, item := range result.Result {
		item.Username = unames[item.UID]
	}
	return
}

//TokenGroupByType .
func (s *Service) TokenGroupByType(c context.Context, netID int64) (result map[string][]*net.Token, err error) {
	var (
		list []*net.Token
	)
	result = make(map[string][]*net.Token)
	if list, err = s.gorm.TokenList(c, []int64{netID}, nil, "", true); err != nil || len(list) == 0 {
		return
	}

	for _, item := range list {
		typeDesc := net.TokenValueTypeDesc[item.Type]
		if _, exist := result[typeDesc]; !exist {
			result[typeDesc] = []*net.Token{item}
			continue
		}
		result[typeDesc] = append(result[typeDesc], item)
	}
	return
}

//TokenByName .
func (s *Service) TokenByName(c context.Context, businessID int64, name string) (result map[string]string, err error) {
	var (
		netID     []int64
		tokenList []*net.Token
	)
	result = map[string]string{}
	if netID, err = s.netIDByBusiness(c, businessID); err != nil {
		log.Error("TokenByName s.netIDByBusiness(%d) error(%v)", businessID, err)
		return
	}
	if len(netID) == 0 {
		return
	}

	if tokenList, err = s.gorm.TokenList(c, netID, nil, name, true); err != nil {
		log.Error("TokenByName s.gorm.TokenList(%v, %s) error(%v) businessid(%d)", netID, name, err, businessID)
		return
	}

	for _, item := range tokenList {
		result[item.Value] = item.ChName
	}
	return
}

//ShowToken .
func (s *Service) ShowToken(c context.Context, id int64) (n *net.Token, err error) {
	return s.gorm.TokenByID(c, id)
}

func (s *Service) checkTokenUnique(c context.Context, netID int64, name string, compare int8, value string) (err error, msg string) {
	var exist *net.Token
	if exist, err = s.gorm.TokenByUnique(c, netID, name, compare, value); err != nil {
		log.Error("checkTokenUnique s.gorm.TokenByUnique(%d,%s,%d,%s) error(%v)", netID, name, compare, value, err)
		return
	}
	if exist != nil {
		err = ecode.AegisUniqueAlreadyExist
		msg = fmt.Sprintf(ecode.AegisUniqueAlreadyExist.Message(), "令牌",
			fmt.Sprintf("%s %s %s", name, net.GetTokenCompare(compare), value))
	}
	return
}

//AddToken .
func (s *Service) AddToken(c context.Context, t *net.Token) (id int64, err error, msg string) {
	if err, msg = s.checkTokenUnique(c, t.NetID, t.Name, t.Compare, t.Value); err != nil {
		return
	}
	if err = s.gorm.AddItem(c, nil, t); err != nil {
		return
	}
	id = t.ID

	//日志
	oper := &model.NetConfOper{
		OID:    t.ID,
		Action: model.LogNetActionNew,
		UID:    t.UID,
		NetID:  t.NetID,
		ChName: t.ChName,
		Diff:   []string{t.FormatLog()},
	}
	s.sendNetConfLog(c, model.LogTypeTokenConf, oper)
	return
}
