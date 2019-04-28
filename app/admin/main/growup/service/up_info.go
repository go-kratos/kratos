package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_video         = 0
	_column        = 2
	_bgm           = 3
	_passTitle     = "创作激励计划(%s)申请成功"
	_rejectTitle   = "创作激励计划(%s)申请未通过"
	_passContent   = "恭喜，您的创作激励计划(%s)申请已通过！"
	_rejectContent = "您的创作激励计划(%s)申请被驳回，驳回原因：%s。如有任何疑问，请参看帮助中心说明。"
)

func getUpTable(typ int) (table string, err error) {
	switch typ {
	case _video:
		table = "up_info_video"
	case _column:
		table = "up_info_column"
	case _bgm:
		table = "up_info_bgm"
	default:
		err = fmt.Errorf("up type error")
	}
	return
}

func text(typ int, base string) (text string, err error) {
	switch typ {
	case _video:
		text = fmt.Sprintf(base, "视频")
	case _column:
		text = fmt.Sprintf(base, "专栏")
	case _bgm:
		text = fmt.Sprintf(base, "素材")
	default:
		err = fmt.Errorf("up type error")
	}
	return
}

// AddUp add user to creative pgc不区分业务，所有的都加入pgc
func (s *Service) AddUp(c context.Context, mid int64, accType int) (err error) {
	id, err := s.dao.Blocked(c, mid)
	if err != nil {
		return
	}
	if id != 0 { // blocked == true
		return ecode.GrowupDisabled
	}

	nickname, category, err := s.dao.CategoryInfo(c, mid)
	if err != nil {
		return
	}
	fans, oac, err := s.dao.Stat(c, mid)
	if err != nil {
		return
	}

	err = s.addUpVideo(c, mid, nickname, accType, category, fans, oac)
	if err != nil {
		log.Error("s.addUpVideo mid(%d), error(%v)", mid, err)
		return
	}

	err = s.addUpColumn(c, mid, nickname, accType, category, fans)
	if err != nil {
		log.Error("s.addUpColumn mid(%d), error(%v)", mid, err)
	}

	err = s.addUpBgm(c, mid, nickname, accType, category, fans)
	if err != nil {
		log.Error("s.addUpBgm mid(%d), error(%v)", mid, err)
	}
	return
}

func (s *Service) addUpVideo(c context.Context, mid int64, nickname string, accType, category, fans, oac int) (err error) {
	ups, err := s.dao.UpsVideoInfo(c, fmt.Sprintf("mid = %d and is_deleted = 0", mid))
	if err != nil {
		log.Error("s.dao.UpsInfo mid(%d), error(%v)", mid, err)
		return
	}
	if len(ups) == 1 {
		if ups[0].AccountType == accType {
			return
		}
	}
	up := &model.UpInfo{
		MID:                  mid,
		Nickname:             nickname,
		MainCategory:         category,
		Fans:                 fans,
		AccountType:          accType,
		OriginalArchiveCount: oac,
	}
	if accType == 2 { // if account_type is pgc, sign_type update to first publish
		up.SignType = 2
		up.AccountState = 1
	}
	_, err = s.dao.InsertUpVideo(c, up)
	return
}

func (s *Service) addUpColumn(c context.Context, mid int64, nickname string, accType, category, fans int) (err error) {
	ups, err := s.dao.UpsColumnInfo(c, fmt.Sprintf("mid = %d and is_deleted = 0", mid))
	if err != nil {
		log.Error("s.dao.UpsColumnInfo mids(%d), error(%v)", mid, err)
		return
	}
	if len(ups) == 1 {
		if ups[0].AccountType == accType {
			return
		}
	}

	up := &model.UpInfo{
		MID:          mid,
		Nickname:     nickname,
		MainCategory: category,
		Fans:         fans,
		AccountType:  accType,
	}

	if accType == 2 { // if account_type is pgc, sign_type update to first publish
		up.SignType = 2
		up.AccountState = 1
	}
	_, err = s.dao.InsertUpColumn(c, up)
	return
}

func (s *Service) addUpBgm(c context.Context, mid int64, nickname string, accType, category, fans int) (err error) {
	ups, err := s.dao.UpsBgmInfo(c, fmt.Sprintf("mid=%d AND is_deleted=0", mid))
	if err != nil {
		return
	}
	if len(ups) == 1 {
		if ups[0].AccountType == accType {
			return
		}
	}
	count, err := s.dao.BGMCount(c, mid)
	if err != nil {
		return
	}

	up := &model.UpInfo{
		MID:         mid,
		Nickname:    nickname,
		Fans:        fans,
		AccountType: accType,
		BGMs:        count,
	}
	if accType == 2 {
		up.SignType = 2
		up.AccountState = 1
	}
	_, err = s.dao.InsertBgmUpInfo(c, up)
	return
}

func (s *Service) getBusinessType(c context.Context, mids []int64, state int) (business map[int64][]int, err error) {
	business = make(map[int64][]int)
	// video
	video, err := s.dao.GetUpInfoByState(c, "up_info_video", mids, state)
	if err != nil {
		log.Error("s.dao.GetUpInfoSigned error(%v)", err)
		return
	}
	// column
	column, err := s.dao.GetUpInfoByState(c, "up_info_column", mids, state)
	if err != nil {
		log.Error("s.dao.GetUpInfoSigned error(%v)", err)
		return
	}
	// bgm
	bgm, err := s.dao.GetUpInfoByState(c, "up_info_bgm", mids, state)
	if err != nil {
		log.Error("s.dao.GetUpInfoSigned error(%v)", err)
		return
	}
	for _, mid := range mids {
		business[mid] = make([]int, 0)
		if _, ok := video[mid]; ok {
			business[mid] = append(business[mid], 0)
		}
		if _, ok := column[mid]; ok {
			business[mid] = append(business[mid], 2)
		}
		if _, ok := bgm[mid]; ok {
			business[mid] = append(business[mid], 3)
		}
	}
	return
}

// QueryFromUpInfo query up-info in growup plan
func (s *Service) QueryFromUpInfo(c context.Context, busType int, accType int, states []int64, mid int64, category int, signType int, nickname string, lower int, upper int, from int, limit int, sort string) (ups []*model.UpInfo, total int, err error) {
	table, err := getUpTable(busType)
	if err != nil {
		return
	}
	query := upsInfoStmt(accType, states, mid, category, signType, nickname, lower, upper)
	total, err = s.dao.UpsCount(c, table, query)
	if err != nil {
		return
	}
	sb := sortBy(sort)
	if len(sb) != 0 {
		query += " " + sb
	}
	query += fmt.Sprintf(" LIMIT %d, %d", from, limit)
	switch busType {
	case _video:
		ups, err = s.dao.UpsVideoInfo(c, query)
		if err != nil {
			log.Error("s.dao.UpsInfo mids(%+d), state(%d) error(%v)", mid, states, err)
			return
		}
	case _column:
		ups, err = s.dao.UpsColumnInfo(c, query)
		if err != nil {
			log.Error("s.dao.UpsColumnInfo mids(%+d), state(%d) error(%v)", mid, states, err)
			return
		}
	case _bgm:
		ups, err = s.dao.UpsBgmInfo(c, query)
		if err != nil {
			log.Error("s.dao.UpsBgmInfo mids(%+d), state(%d) error(%v)", mid, states, err)
			return
		}
	}

	if ups == nil {
		ups = make([]*model.UpInfo, 0)
	}

	if len(ups) == 0 {
		return
	}

	mids := make([]int64, 0)
	for _, up := range ups {
		mids = append(mids, up.MID)
	}

	signedType, err := s.getBusinessType(c, mids, 3)
	if err != nil {
		log.Error("s.getBusinessType error(%v)", err)
		return
	}

	// credit scores
	scores, err := s.dao.CreditScores(c, mids)
	if err != nil {
		return
	}

	for _, up := range ups {
		up.SignedType = signedType[up.MID]
		up.CreditScore = scores[up.MID]
	}

	if len(states) == 1 {
		var other map[int64][]int
		other, err = s.getBusinessType(c, mids, int(states[0]))
		if err != nil {
			log.Error("s.getBusinessType error(%v)", err)
			return
		}
		for _, up := range ups {
			up.OtherType = other[up.MID]
		}
	}
	return
}

func sortBy(name string) (sort string) {
	if len(name) == 0 {
		return
	}
	if strings.HasPrefix(name, "-") {
		name = strings.TrimPrefix(name, "-")
		name += " DESC"
	}
	sort = " ORDER BY " + name
	return
}

func upsInfoStmt(accountType int, states []int64, mid int64, category int, signType int, nickname string, lower, upper int) (query string) {
	if accountType > 0 {
		query += "account_type = " + strconv.Itoa(accountType)
		query += " AND "
	}
	if len(states) > 0 {
		query += fmt.Sprintf("account_state IN (%s)", xstr.JoinInts(states))
		query += " AND "
	}
	if mid > 0 {
		query += "mid = " + strconv.FormatInt(mid, 10)
		query += " AND "
	}
	if category != 0 {
		query += "category_id = " + strconv.Itoa(category)
		query += " AND "
	}
	if signType > 0 {
		query += "sign_type = " + strconv.Itoa(signType)
		query += " AND "
	}
	if nickname != "" {
		query += "nickname = " + "\"" + nickname + "\""
		query += " AND "
	}
	query += "fans >= " + strconv.Itoa(lower)
	query += " AND "
	if upper > 0 {
		query += "fans <=" + strconv.Itoa(upper)
		query += " AND "
	}
	query += "is_deleted = 0"
	return
}

// Reject update account state to 4(reject)
func (s *Service) Reject(c context.Context, typ int, mids []int64, reason string, days int) (err error) {
	table, err := getUpTable(typ)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	_, err = s.dao.Reject(c, table, 4, reason, xtime.Time(now), xtime.Time(now+int64(86400*days)), mids)
	if err != nil {
		return
	}
	title, err := text(typ, _rejectTitle)
	if err != nil {
		return
	}
	var content string
	switch typ {
	case _video:
		content = fmt.Sprintf(_rejectContent, "视频", reason)
	case _column:
		content = fmt.Sprintf(_rejectContent, "专栏", reason)
	case _bgm:
		content = fmt.Sprintf(_rejectContent, "素材", reason)
	}
	pushErr := s.msg.Send(c, "1_14_2", title, content, mids, now)
	if pushErr != nil {
		log.Error("reject push error(%v)", pushErr)
	}
	return
}

// Pass update account state to 3(signed) type 0.video 1.audio 2.column 3.bgm,
func (s *Service) Pass(c context.Context, mids []int64, typ int) (err error) {
	if len(mids) == 0 {
		return
	}
	table, err := getUpTable(typ)
	if err != nil {
		return
	}
	ms, err := s.dao.Pendings(c, mids, table)
	if err != nil {
		return
	}
	if len(ms) == 0 {
		return
	}
	upM := make(map[int64]struct{})
	for _, m := range ms {
		upM[m] = struct{}{}
	}
	// check other two business
	if typ != _video {
		err = s.removeUnusualUps(c, "up_info_video", upM, ms)
		if err != nil {
			log.Error("s.removeUnusualUps error(%v)", err)
			return
		}
	}
	if typ != _column {
		err = s.removeUnusualUps(c, "up_info_column", upM, ms)
		if err != nil {
			log.Error("s.removeUnusualUps error(%v)", err)
			return
		}
	}
	if typ != _bgm {
		err = s.removeUnusualUps(c, "up_info_bgm", upM, ms)
		if err != nil {
			log.Error("s.removeUnusualUps error(%v)", err)
			return
		}
	}
	if len(upM) == 0 {
		return
	}
	ms = make([]int64, 0)
	for mid := range upM {
		ms = append(ms, mid)
	}
	_, err = s.dao.Pass(c, table, 3, xtime.Time(time.Now().Unix()), ms)
	if err != nil {
		return
	}
	_, err = s.dao.InsertCreditScore(c, midValues(mids))
	if err != nil {
		return
	}
	title, err := text(typ, _passTitle)
	if err != nil {
		return
	}
	msg, err := text(typ, _passContent)
	if err != nil {
		return
	}
	pushErr := s.msg.Send(c, "1_14_1", title, msg, ms, time.Now().Unix())
	if pushErr != nil {
		log.Error("pass push error(%v)", pushErr)
	}
	// add creative task notify
	s.msg.NotifyTask(c, mids)
	return
}

func (s *Service) removeUnusualUps(c context.Context, table string, upM map[int64]struct{}, ms []int64) (err error) {
	mids, err := s.dao.UnusualUps(c, ms, table)
	if err != nil {
		return
	}
	for _, mid := range mids {
		if _, ok := upM[mid]; ok {
			delete(upM, mid)
		}
	}
	return
}

func midValues(mids []int64) (values string) {
	var buf bytes.Buffer
	for _, mid := range mids {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

// Dismiss update account state to 6
func (s *Service) Dismiss(c context.Context, operator string, typ, oldState int, mid int64, reason string) (err error) {
	var (
		tx      *sql.Tx
		now     = xtime.Time(time.Now().Unix())
		current int
		drows   int64
		crows   int64
	)

	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}

	if current, err = s.dao.CreditScore(c, mid); err != nil {
		tx.Rollback()
		return
	}

	// deduct credit score 10
	var drows1, drows2, drows3 int64
	if drows1, err = s.dao.TxDismiss(tx, "up_info_video", 6, oldState, reason, now, now, mid); err != nil {
		tx.Rollback()
		return
	}

	if drows2, err = s.dao.TxDismiss(tx, "up_info_column", 6, oldState, reason, now, now, mid); err != nil {
		tx.Rollback()
		return
	}

	if drows3, err = s.dao.TxDismiss(tx, "up_info_bgm", 6, oldState, reason, now, now, mid); err != nil {
		tx.Rollback()
		return
	}

	switch typ {
	case _video:
		drows = drows1
	case _column:
		drows = drows2
	case _bgm:
		drows = drows3
	}
	if crows, err = s.txInsertCreditRecord(tx, mid, operator, 6, 10, current-10); err != nil {
		tx.Rollback()
		return
	}
	if drows != crows {
		tx.Rollback()
		return
	}

	if typ == _video {
		if _, err = s.dao.TxUpdateUpSpyState(tx, 6, mid); err != nil {
			tx.Rollback()
			return
		}

		if _, err = s.dao.DelCheatUp(c, mid); err != nil {
			tx.Rollback()
			return
		}
	}

	_, err = s.dao.TxUpdateCreditScore(tx, mid, current-10)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	pushErr := s.msg.Send(c, "1_14_4", "您的创作激励计划参与资格已被取消", fmt.Sprintf(`您好，根据创作激励计划规则，您因%s，已被永久取消参与创作激励计划的资格。如有任何疑问，请联系客服。`, reason), []int64{mid}, time.Now().Unix())
	if pushErr != nil {
		log.Error("dismiss push error(%v)", pushErr)
	}
	return
}

func (s *Service) txInsertCreditRecord(tx *sql.Tx, mid int64, operator string, state, deducted, remaining int) (rows int64, err error) {
	// insert to credit record
	cr := &model.CreditRecord{
		MID:       mid,
		OperateAt: xtime.Time(time.Now().Unix()),
		Operator:  operator,
		Reason:    state,
		Deducted:  deducted,
		Remaining: remaining,
	}
	return s.dao.TxInsertCreditRecord(tx, cr)
}

// Forbid update account state to 7 and add a n days CD
func (s *Service) Forbid(c context.Context, operator string, typ, oldState int, mid int64, reason string, days, second int) (err error) {
	var (
		tx        *sql.Tx
		now       = time.Now().Unix()
		expiredIn = xtime.Time(now + int64(second))
		current   int
		frows     int64
		crows     int64
	)

	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}

	if current, err = s.dao.CreditScore(c, mid); err != nil {
		tx.Rollback()
		return
	}

	var frows1, frows2, frows3 int64
	// deduct credit score 5
	if frows1, err = s.dao.TxForbid(tx, "up_info_video", 7, oldState, reason, xtime.Time(now), expiredIn, mid); err != nil {
		tx.Rollback()
		return
	}
	if frows2, err = s.dao.TxForbid(tx, "up_info_column", 7, oldState, reason, xtime.Time(now), expiredIn, mid); err != nil {
		tx.Rollback()
		return
	}
	if frows3, err = s.dao.TxForbid(tx, "up_info_bgm", 7, oldState, reason, xtime.Time(now), expiredIn, mid); err != nil {
		tx.Rollback()
		return
	}
	switch typ {
	case _video:
		frows = frows1
	case _column:
		frows = frows2
	case _bgm:
		frows = frows3
	}
	// insert credit record
	if crows, err = s.txInsertCreditRecord(tx, mid, operator, 7, 5, current-5); err != nil {
		tx.Rollback()
		return
	}
	if frows != crows {
		tx.Rollback()
		return
	}

	if typ == _video {
		// update up spy state
		if _, err = s.dao.TxUpdateUpSpyState(tx, 7, mid); err != nil {
			tx.Rollback()
			return
		}

		// del up from cheat fans list
		if _, err = s.dao.DelCheatUp(c, mid); err != nil {
			tx.Rollback()
			return
		}
	}

	_, err = s.dao.TxUpdateCreditScore(tx, mid, current-5)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}

	pushErr := s.msg.Send(c, "1_14_3", "您已被暂停参加创作激励计划", fmt.Sprintf(`根据创作激励计划规则，您因为%s原因被停止参与计划%d天，%d年%d月%d日恢复，如有疑问，请联系客服。`, reason, days, expiredIn.Time().Year(), expiredIn.Time().Month(), expiredIn.Time().Day()), []int64{mid}, time.Now().Unix())
	if pushErr != nil {
		log.Error("forbid push error(%v)", pushErr)
	}
	return
}

// Recovery update (video/column) account state to 3(signed)
func (s *Service) Recovery(c context.Context, mid int64) (err error) {
	err = s.UpdateUpAccountState(c, "up_info_video", mid, 3)
	if err != nil {
		log.Error("s.UpdateUpAccountState(video) error(%v)", err)
		return
	}
	err = s.UpdateUpAccountState(c, "up_info_column", mid, 3)
	if err != nil {
		log.Error("s.UpdateUpAccountState(column) error(%v)", err)
		return
	}
	err = s.UpdateUpAccountState(c, "up_info_bgm", mid, 3)
	if err != nil {
		log.Error("s.UpdateUpAccountState(bgm) error(%v)", err)
	}
	return
}

// UpdateUpAccountState update video up account state
func (s *Service) UpdateUpAccountState(c context.Context, table string, mid int64, state int) (err error) {
	_, err = s.dao.UpdateAccountState(c, table, state, mid)
	return
}

// DeleteUp delete up from up_info (update is_deleted = 1)
func (s *Service) DeleteUp(c context.Context, mid int64) (err error) {
	_, err = s.dao.DelUpInfo(c, "up_info_video", mid)
	if err != nil {
		return
	}
	_, err = s.dao.DelUpInfo(c, "up_info_column", mid)
	if err != nil {
		return
	}
	_, err = s.dao.DelUpInfo(c, "up_info_bgm", mid)
	return
}

// Block add to blacklist
func (s *Service) Block(c context.Context, mid int64) (err error) {
	nickname, categoryID, err := s.dao.CategoryInfo(c, mid)
	if err != nil {
		return
	}
	fans, oac, err := s.dao.Stat(c, mid)
	if err != nil {
		return
	}
	applyAt, err := s.dao.ApplyAt(c, mid)
	if err != nil {
		return
	}
	b := &model.Blocked{
		MID:                  mid,
		Nickname:             nickname,
		OriginalArchiveCount: oac,
		MainCategory:         categoryID,
		Fans:                 fans,
		ApplyAt:              applyAt,
	}

	_, err = s.dao.InsertBlocked(c, b)
	if err != nil {
		return
	}
	// if up in table up_info_video, delete
	_, err = s.dao.DelUpInfo(c, "up_info_video", mid)
	if err != nil {
		return
	}
	_, err = s.dao.DelUpInfo(c, "up_info_column", mid)
	if err != nil {
		return
	}
	_, err = s.dao.DelUpInfo(c, "up_info_bgm", mid)
	return
}

// QueryFromBlocked query up-info in black list of growup plan
func (s *Service) QueryFromBlocked(c context.Context, mid int64, category int, nickname string, lower, upper, from, limit int, sort string) (ups []*model.Blocked, total int, err error) {
	query := queryBlockStmt(mid, category, nickname, lower, upper)
	total, err = s.dao.BlockCount(c, query)
	if err != nil {
		return
	}
	sb := sortBy(sort)
	if len(sb) != 0 {
		query += " " + sb
	}
	query += fmt.Sprintf(" LIMIT %d, %d", from, limit)
	ups, err = s.dao.QueryFromBlocked(c, query)
	if err != nil {
		log.Error("s.dao.QueryFromBlocked error(%v)", err)
		return
	}
	if ups == nil {
		ups = make([]*model.Blocked, 0)
	}
	return
}

func queryBlockStmt(mid int64, categoryID int, nickname string, lower int, upper int) (query string) {
	if mid > 0 {
		query += "mid = " + strconv.FormatInt(mid, 10)
		query += " AND "
	}
	if categoryID != 0 {
		query += "category_id = " + strconv.Itoa(categoryID)
		query += " AND "
	}
	if nickname != "" {
		query += "nickname = " + "\"" + nickname + "\""
		query += " AND "
	}
	query += "fans >= " + strconv.Itoa(lower)
	query += " AND "
	query += "is_deleted = 0"
	if upper > 0 {
		query += " AND "
		query += "fans <=" + strconv.Itoa(upper)
	}
	return
}

// DeleteFromBlocked del blocked and recover up info of video
func (s *Service) DeleteFromBlocked(c context.Context, mid int64) (err error) {
	_, err = s.dao.DelFromBlocked(c, mid)
	if err != nil {
		return
	}
	_, err = s.dao.RecUpInfo(c, "up_info_video", mid)
	if err != nil {
		return
	}
	_, err = s.dao.RecUpInfo(c, "up_info_column", mid)
	if err != nil {
		return
	}
	_, err = s.dao.RecUpInfo(c, "up_info_bgm", mid)
	return
}

// DelUpAccount del mid from up_account
func (s *Service) DelUpAccount(c context.Context, mid int64) (err error) {
	_, err = s.dao.DelUpAccount(c, mid)
	return
}

// UpdateUpAccount update up_account
func (s *Service) UpdateUpAccount(c context.Context, mid int64, isDeleted int, withdrawDate string) (err error) {
	_, err = s.dao.UpdateUpAccount(c, mid, isDeleted, withdrawDate)
	return
}

// CreditRecords get credit records by mid
func (s *Service) CreditRecords(c context.Context, mid int64) (crs []*model.CreditRecord, err error) {
	return s.dao.CreditRecords(c, mid)
}

// RecoverCreditScore recover credit score
func (s *Service) RecoverCreditScore(c context.Context, typ int, id, mid int64) (err error) {
	var (
		tx       *sql.Tx
		deducted int
		drows    int64
		urows    int64
	)

	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}

	if deducted, err = s.dao.DeductedScore(c, id); err != nil {
		tx.Rollback()
		return
	}

	// del detucted record by id
	if drows, err = s.dao.TxDelCreditRecord(tx, id); err != nil {
		tx.Rollback()
		return
	}

	// recover credit score
	if urows, err = s.dao.TxRecoverCreditScore(tx, deducted, mid); err != nil {
		return
	}

	if drows != urows {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	return
}

// ExportUps export ups by query
func (s *Service) ExportUps(c context.Context, busType, accType int, states []int64, mid int64, category int, signType int, nickname string, lower int, upper int, from int, limit int, sort string) (res []byte, err error) {
	ups, _, err := s.QueryFromUpInfo(c, busType, accType, states, mid, category, signType, nickname, lower, upper, from, limit, sort)
	if err != nil {
		log.Error("s.QueryFromUpInfo error(%v)", err)
		return
	}

	records := formatUpInfo(ups, states, busType)
	res, err = FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)", err)
	}
	return
}

// UpState get up state
func (s *Service) UpState(c context.Context, mid int64, typ int) (data interface{}, err error) {
	table, err := getUpTable(typ)
	if err != nil {
		return
	}
	state, err := s.dao.GetUpState(c, table, mid)
	if err != nil {
		return
	}
	data = map[string]interface{}{
		"mid":   mid,
		"state": state,
	}
	return
}
