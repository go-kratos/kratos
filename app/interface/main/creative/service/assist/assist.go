package assist

import (
	"context"

	"go-common/app/interface/main/creative/model/assist"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
)

// Assists get assists
func (s *Service) Assists(c context.Context, mid int64, ip string) (assists []*assist.Assist, err error) {
	var (
		mainAssists      []*assist.Assist
		liveAssists      []*assist.LiveAssist
		rightsMap        = make(map[int64]*assist.Rights, 20)
		mainMap          = make(map[int64]*assist.Assist, 10)
		assistMids       []int64
		liveAssistMids   []int64
		liveAssistsTotal map[int64]map[int8]map[int8]int
	)
	// main assists
	if mainAssists, err = s.assist.Assists(c, mid, ip); err != nil {
		log.Error("s.assist.Assists err(%v) | mid(%d), ip(%s)", err, mid, ip)
		return
	}
	// live assists
	if liveAssists, err = s.assist.LiveAssists(c, mid, ip); err != nil {
		log.Error("s.assist.Assists err(%v) | mid(%d), ip(%s)", err, mid, ip)
		return
	}
	// rights info
	for _, ass := range mainAssists {
		rightsMap[ass.AssistMid] = &assist.Rights{Main: 1}
		mainMap[ass.AssistMid] = ass
	}
	for _, ass := range liveAssists {
		if _, ok := rightsMap[ass.AssistMid]; !ok {
			rightsMap[ass.AssistMid] = &assist.Rights{Live: 1}
		} else {
			rightsMap[ass.AssistMid].Live = 1
		}
		liveAssistMids = append(liveAssistMids, ass.AssistMid)
	}
	// get live assist total
	if len(liveAssistMids) > 0 {
		liveAssistsTotal, _ = s.assist.Stat(c, mid, liveAssistMids, ip)
	}
	// get user info
	for mid := range rightsMap {
		assistMids = append(assistMids, mid)
	}
	users, err := s.acc.Infos(c, assistMids, ip)
	if err != nil {
		log.Error("s.acc.Infos() err(%v)", err)
		return
	}
	for _, ast := range mainAssists {
		if _, ok := users[ast.AssistMid]; ok {
			ast.AssistAvatar = users[ast.AssistMid].Face
			ast.AssistName = users[ast.AssistMid].Name
			ast.Banned, _ = s.userBanned(c, ast.AssistMid)
			ast.Rights = rightsMap[ast.AssistMid]
			assists = append(assists, ast)
		}
	}
	for _, ass := range liveAssists {
		if _, ok := mainMap[ass.AssistMid]; !ok {
			if _, ok := users[ass.AssistMid]; ok {
				ast := &assist.Assist{}
				ast.AssistMid = ass.AssistMid
				ast.AssistAvatar = users[ass.AssistMid].Face
				ast.AssistName = users[ass.AssistMid].Name
				ast.CTime = ass.CTime
				ast.MTime = ass.CTime
				ast.Banned, _ = s.userBanned(c, ass.AssistMid)
				ast.Rights = rightsMap[ass.AssistMid]
				if _, ok := liveAssistsTotal[ass.AssistMid]; ok {
					ast.Total = liveAssistsTotal[ast.AssistMid]
				}
				assists = append(assists, ast)
			}
		}
	}
	return
}

// AssistLogs get all assistlog by Mid, and format by typeid
func (s *Service) AssistLogs(c context.Context, mid, assistMid, pn, ps, stime, etime int64, ip string) (assistLogs []*assist.AssistLog, pager map[string]int64, err error) {
	if assistLogs, pager, err = s.assist.AssistLogs(c, mid, assistMid, pn, ps, stime, etime, ip); err != nil {
		log.Error("s.creative.Assist err(%v) | mid(%d),assistMid(%d),pn(%d),ps(%d),stime(%d),etime(%d),ip(%s)",
			err, mid, assistMid, pn, ps, stime, etime, ip)
		return
	}
	var assistMids []int64
	for _, ass := range assistLogs {
		assistMids = append(assistMids, ass.AssistMid)
	}
	users, err := s.acc.Infos(c, assistMids, ip)
	if err != nil {
		log.Error("s.acc.Profiles() err(%v)", err)
		return
	}
	for _, assist := range assistLogs {
		if _, ok := users[assist.AssistMid]; ok {
			assist.AssistAvatar = users[assist.AssistMid].Face
			assist.AssistName = users[assist.AssistMid].Name
		}
	}
	return
}

// AddAssist add mid to any assist
func (s *Service) AddAssist(c context.Context, mid, assistMid int64, main, live int8, ip, ak, ck string) (err error) {
	if main == 1 {
		if err = s.addAssist(c, mid, assistMid, ip, ak, ck); err != nil {
			return
		}
	}
	if live == 1 {
		if err = s.liveAddAssist(c, mid, assistMid, ak, ck, ip); err != nil {
			return
		}
	}
	return
}

func (s *Service) addAssist(c context.Context, mid, assistMid int64, ip, ak, ck string) (err error) {
	var (
		card *account.Card
	)
	if card, err = s.acc.Card(c, mid, ip); err != nil {
		log.Error("s.assist.AddAssist err(%v) | mid(%d), assistMid(%d), ip(%s)", err, mid, assistMid, ip)
		return
	}
	if err = s.assist.AddAssist(c, mid, assistMid, ip, card.Name); err != nil {
		log.Error("s.assist.AddAssist err(%v) | mid(%d), assistMid(%d), ip(%s)", err, mid, assistMid, ip)
		return
	}
	return
}

// DelAssist delete all the assist
func (s *Service) DelAssist(c context.Context, mid, assistMid int64, ip, ak, ck string) (err error) {
	isMainAss, _ := s.assist.Info(c, mid, assistMid, ip)
	if isMainAss == 1 {
		err = s.delAssist(c, mid, assistMid, ip, ak, ck)
	}
	isLiveAss, _ := s.LiveCheckAssist(c, mid, assistMid, ip)
	if isLiveAss == 1 {
		err = s.liveDelAssist(c, mid, assistMid, ck, ip)
	}
	return
}

// DelAssist delete the assist
func (s *Service) delAssist(c context.Context, mid, assistMid int64, ip, ak, ck string) (err error) {
	var (
		card *account.Card
	)
	if card, err = s.acc.Card(c, mid, ip); err != nil {
		log.Error("s.assist.AddAssist err(%v) | mid(%d), assistMid(%d), ip(%s)", err, mid, assistMid, ip)
		return
	}
	if err = s.assist.DelAssist(c, mid, assistMid, ip, card.Name); err != nil {
		log.Error("s.assist.DelAssist err(%v) | mid(%d), assistMid(%d), ip(%s)", err, mid, assistMid, ip)
		return
	}
	return
}

// prepare for Live revoc log
func (s *Service) preRevocLogForLive(c context.Context, assistLog *assist.AssistLog) (err error) {
	if assistLog.Action == 8 && assistLog.Type == 3 {
		var (
			objID   int64
			revcLog = &assist.AssistLog{}
		)
		if objID, err = strconv.ParseInt(assistLog.ObjectID, 10, 64); err != nil {
			log.Error("Atoi err(%v) | mid(%d), assistMid(%d), assistLog.ObjectID(%s)", err, assistLog.Mid, assistLog.AssistMid, assistLog.ObjectID)
			return
		}
		revcLog, err = s.assist.AssistLogObj(c, assistLog.Type, 9, assistLog.Mid, objID)
		if err == ecode.AssistLogNotExist {
			return nil
		}
		if err != nil {
			log.Error("s.assist.AssistLogObj err(%v) | mid(%d), assistMid(%d), logID(%d)", err, assistLog.Mid, assistLog.AssistMid, assistLog.ID)
			return
		}
		if revcLog != nil {
			err = ecode.CreativeAssistLogAlreadyRevoc
			return
		}
		return
	}
	return
}

// RevocAssistLog cancel this asssist action
func (s *Service) RevocAssistLog(c context.Context, mid, assistMid, logID int64, ck, ip string) (err error) {
	var assistLog *assist.AssistLog
	if assistLog, err = s.assist.AssistLog(c, mid, assistMid, logID, ip); err != nil {
		log.Error("s.assist.AssistLog err(%v) | mid(%d), assistMid(%d), logID(%d), ip(%s)", err, mid, assistMid, logID, ip)
		return
	}
	if err = s.preRevocLogForLive(c, assistLog); err != nil {
		log.Error("s.preRevocLogForLive err(%v) | mid(%d), assistMid(%d), logID(%d), ip(%s)", err, mid, assistMid, logID, ip)
		return
	}
	if err = s.revoc(c, assistLog, ck, ip); err != nil {
		log.Error("s.revoc err(%v) | assistLog(%+v), ip(%s)", err, assistLog, ip)
		return
	}
	if err = s.assist.RevocAssistLog(c, mid, assistMid, logID, ip); err != nil {
		log.Error("s.assist.CancelAssistLog err(%v) | mid(%d), assistMid(%d), logID(%d), ip(%s)", err, mid, assistMid, logID, ip)
		return
	}
	return
}

func (s *Service) userBanned(c context.Context, mid int64) (banned int8, err error) {
	var card *account.Card
	if card, err = s.acc.Card(c, mid, ""); err != nil {
		err = nil
		return
	}
	if card.Silence == 1 {
		banned = 1
		return
	}
	return
}

// SetAssist get user assist rights about liveRoom
func (s *Service) SetAssist(c context.Context, mid, assistMid int64, main, live int8, ip, ak, ck string) (err error) {
	isMainAss, _ := s.assist.Info(c, mid, assistMid, ip)
	if isMainAss == 0 && main == 1 {
		err = s.addAssist(c, mid, assistMid, ip, ak, ck)
	} else if isMainAss == 1 && main == 0 {
		err = s.delAssist(c, mid, assistMid, ip, ak, ck)
	}
	if err != nil {
		return
	}
	isLiveAss, _ := s.LiveCheckAssist(c, mid, assistMid, ip)
	if live == 1 && isLiveAss == 0 {
		err = s.liveAddAssist(c, mid, assistMid, ak, ck, ip)
	} else if live == 0 && isLiveAss == 1 {
		err = s.liveDelAssist(c, mid, assistMid, ck, ip)
	}
	return
}
