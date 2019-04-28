package income

import (
	"context"
	"fmt"
	"sort"
	"time"

	model "go-common/app/job/main/growup/model/income"
	task "go-common/app/job/main/growup/service"
	"go-common/library/log"

	"golang.org/x/sync/errgroup"
)

// RunAndSendMail run and send email
func (s *Service) RunAndSendMail(c context.Context, date time.Time) (err error) {
	var mailReceivers []string
	var msg string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 3 {
			mailReceivers = v.Addr
		}
	}
	startTime := time.Now().Unix()
	err = s.run(c, date)
	if err != nil {
		msg = err.Error()
		mailReceivers = []string{"shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com"}
	} else {
		msg = fmt.Sprintf("%s 计算完成，耗时%ds", date.Format("2006-01-02"), time.Now().Unix()-startTime)
	}
	emailErr := s.email.SendMail(date, msg, "创作激励每日计算%d年%d月%d日", mailReceivers...)
	if emailErr != nil {
		log.Error("s.email.SendMail error(%v)", emailErr)
	}
	return
}

func (s *Service) run(c context.Context, date time.Time) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskCreativeIncome, date.Format(_layout), err)
	}()

	err = task.GetTaskService().TaskReady(c, date.Format("2006-01-02"), task.TaskAvCharge, task.TaskCmCharge, task.TaskTagRatio, task.TaskBubbleMeta, task.TaskBlacklist, task.TaskBgmSync)
	if err != nil {
		return
	}

	startWeeklyDate = getStartWeeklyDate(date)
	startMonthlyDate = getStartMonthlyDate(date)

	/*################ Serializable Begin  ################*/

	// av charge ratio
	ratios, err := s.ratio.ArchiveChargeRatio(c, int64(_limitSize))
	if err != nil {
		return
	}
	// up charge ratio
	urs, err := s.ratio.UpChargeRatio(c, int64(_limitSize))
	if err != nil {
		return
	}
	// av income statistics
	astat, err := s.income.avIncomeStatSvr.AvIncomeStat(c, int64(_limitSize))
	if err != nil {
		log.Error("s.income.avIncomeStatSvr.AvIncomeStat error(%v) ", err)
		return
	}
	log.Info("get av_income_statis : %d", len(astat))

	// bgm income statistics
	bstat, err := s.income.bgmIncomeStatSvr.BgmIncomeStat(c, int64(_limitSize))
	if err != nil {
		log.Error("s.income.bgmIncomeStatSvr.BgmIncomeStat error(%v) ", err)
		return
	}
	log.Info("get bgm_income_statis : %d", len(bstat))

	// column income statistics
	cstat, err := s.income.columnIncomeStatSvr.ColumnIncomeStat(c, int64(_limitSize))
	if err != nil {
		log.Error("s.income.columnIncomeStatSvr.ColumnIncomeStat error(%v) ", err)
		return
	}
	log.Info("get column_income_statis : %d", len(cstat))

	// up income statistics
	ustat, err := s.income.upIncomeStatSvr.UpIncomeStat(c, int64(_limitSize))
	if err != nil {
		log.Error("s.income.upIncomeStatSvr.UpIncomeStat error(%v) ", err)
		return
	}
	log.Info("get up_income_statis : %d", len(ustat))

	// up accounts
	accs, err := s.income.upAccountSvr.UpAccount(c, int64(_limitSize))
	if err != nil {
		log.Error("s.income.upAccountSvr.UpAccount error(%v)", err)
		return
	}
	log.Info("get up_account : %d", len(accs))

	// bubble meta
	bubbleMeta, err := s.GetBubbleMeta(c)
	if err != nil {
		log.Error("s.GetBubbleMeta error(%v)", err)
		return
	}
	log.Info("get lottery_av_info avids: %d", len(bubbleMeta))
	bubbleRatio := s.avToBubbleRatio(bubbleMeta)

	// av signed ups
	var (
		avFilters     []AvFilter
		columnFilters []ColumnFilter
		bgmFilter     BgmFilter
	)

	//black list
	blacks, err := s.Blacklist(c, 2000)
	if err != nil {
		return
	}
	// black ctype 0: av
	avBlackFilter := avFilter(blacks[0])

	// black ctype 2: column
	columnBlackFilter := columnFilter(blacks[2])

	// business orders
	bos, err := s.GetBusinessOrders(c, 2000)
	if err != nil {
		return
	}
	bosFilter := avFilter(bos)

	// signed up
	signed := make(map[int64]bool)

	signedAv, err := s.Signed(c, "video", 2000)
	if err != nil {
		return
	}
	savf := signedAvFilter(signedAv, date)
	for mid := range signedAv {
		signed[mid] = true
	}

	signedColumn, err := s.Signed(c, "column", 2000)
	if err != nil {
		return
	}
	for mid := range signedColumn {
		signed[mid] = true
	}

	signedBgm, err := s.Signed(c, "bgm", 2000)
	if err != nil {
		return
	}
	for mid := range signedBgm {
		signed[mid] = true
	}

	bgms, err := s.BGMs(c, 2000)
	if err != nil {
		return
	}

	{
		avFilters = append(avFilters, avBlackFilter)
		avFilters = append(avFilters, bosFilter)
		avFilters = append(avFilters, savf)

		bgmFilter = signedBgmFilter(signedBgm, date)

		columnFilters = append(columnFilters, signedColumnFilter(signedColumn, date))
		columnFilters = append(columnFilters, columnBlackFilter)
	}

	/*################ Serializable End  ##################*/

	var (
		readGroup errgroup.Group
		sourceCh  = make(chan []*model.AvCharge, 1000)
		// av
		incomeCh = make(chan []*model.AvCharge, 1000)
		// bgm
		bgmCh = make(chan []*model.AvCharge, 1000)
		// column
		columnSourceCh = make(chan []*model.ColumnCharge, 1000)
		// business income
		businessCh = make(chan map[int64]*model.UpBusinessIncome, 10)
	)

	// get av daily charge and repost to other channels
	readGroup.Go(func() (err error) {
		err = s.avCharge.AvCharges(c, date, sourceCh, bubbleRatio)
		if err != nil {
			log.Error("s.avCharge.AvCharges error(%v)", err)
			return
		}
		log.Info("av_daily_charge finished")
		return
	})

	// get column daily charge
	readGroup.Go(func() (err error) {
		err = s.columnCharges(c, date, columnSourceCh)
		if err != nil {
			log.Error("s.columnCharges error(%v)", err)
			return
		}
		log.Info("column_daily_charge finished")
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(incomeCh)
			close(bgmCh)
		}()

		for charges := range sourceCh {
			incomeCh <- charges
			bgmCh <- charges
		}
		return
	})

	// up and av income compute
	var (
		um map[int64]*model.UpIncome
		am map[int64][]*model.AvIncome
		bm map[int64]map[int64]map[int64]*model.BgmIncome
		cm map[int64][]*model.ColumnIncome
	)

	readGroup.Go(func() (err error) {
		//um, am = s.income.Compute(c, date, incomeCh, urs, ars, ustat, astat, accs, filters, signed)
		var business map[int64]*model.UpBusinessIncome
		am, business = s.income.CalAvIncome(incomeCh, urs[1], ratios[1], avFilters, signed)
		businessCh <- business
		return
	})

	readGroup.Go(func() (err error) {
		var business map[int64]*model.UpBusinessIncome
		bm, business = s.income.CalBgmIncome(bgmCh, bgms, urs[3], ratios[3], avFilters, bgmFilter, blacks[3], signed)
		businessCh <- business
		return
	})

	readGroup.Go(func() (err error) {
		var business map[int64]*model.UpBusinessIncome
		cm, business = s.income.CalColumnIncome(columnSourceCh, urs[2], ratios[2], columnFilters, signed)
		businessCh <- business
		return
	})

	readGroup.Go(func() (err error) {
		um = s.income.CalUpIncome(businessCh, date)
		s.income.IncomeStat(um, am, bm, cm, ustat, astat, bstat, cstat)
		s.income.PurgeUpAccount(date, accs, um)
		return
	})

	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	// security verification
	{
		if len(am) == 0 {
			err = fmt.Errorf("Error: insert 0 av_income")
			return
		}
		if len(bm) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_income")
			return
		}
		if len(cm) == 0 {
			err = fmt.Errorf("Error: insert 0 column_income")
			return
		}
		if len(um) == 0 {
			err = fmt.Errorf("Error: insert 0 up_income")
			return
		}
		if len(astat) == 0 {
			err = fmt.Errorf("Error: insert 0 av_income_statis")
			return
		}
		if len(bstat) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_income_statis")
			return
		}
		if len(cstat) == 0 {
			err = fmt.Errorf("Error: insert 0 column_income_statis")
			return
		}
		if len(ustat) == 0 {
			err = fmt.Errorf("Error: insert 0 up_income_statis")
			return
		}
		if len(accs) == 0 {
			err = fmt.Errorf("Error: insert 0 up_account")
			return
		}
	}

	// persistent
	var writeGroup errgroup.Group

	// av_income
	writeGroup.Go(func() (err error) {
		err = s.income.avIncomeSvr.BatchInsertAvIncome(c, am)
		if err != nil {
			log.Error("s.income.BatchInsertAvIncome error(%v)", err)
			return
		}
		log.Info("insert av_income : %d", len(am))
		return
	})

	// column_income
	writeGroup.Go(func() (err error) {
		err = s.income.columnIncomeSvr.BatchInsertColumnIncome(c, cm)
		if err != nil {
			log.Error("s.income.BatchInsertColumnIncome error(%v)", err)
			return
		}
		log.Info("insert column_income : %d", len(cm))
		return
	})

	// bgm_income
	writeGroup.Go(func() (err error) {
		err = s.income.bgmIncomeSvr.BatchInsertBgmIncome(c, bm)
		if err != nil {
			log.Error("s.income.BatchInsertBgmIncome error(%v)", err)
			return
		}
		log.Info("insert bgm_income : %d", len(bm))
		return
	})

	// up income
	writeGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.BatchInsertUpIncome(c, um)
		if err != nil {
			log.Error("s.income.BatchInsertUpIncome error(%v)", err)
			return
		}
		log.Info("insert up_income : %d", len(um))
		return
	})

	// av_income_statis
	writeGroup.Go(func() (err error) {
		err = s.income.avIncomeStatSvr.BatchInsertAvIncomeStat(c, astat)
		if err != nil {
			log.Error("s.income.BatchInsertAvIncomeStat error(%v)", err)
			return
		}
		log.Info("insert av_income_statis : %d", len(astat))
		return
	})

	// column_income_statis
	writeGroup.Go(func() (err error) {
		err = s.income.columnIncomeStatSvr.BatchInsertColumnIncomeStat(c, cstat)
		if err != nil {
			log.Error("s.income.BatchInsertColumnIncomeStat error(%v)", err)
			return
		}
		log.Info("insert column_income_statis : %d", len(cstat))
		return
	})

	// bgm_income_statis
	writeGroup.Go(func() (err error) {
		err = s.income.bgmIncomeStatSvr.BatchInsertBgmIncomeStat(c, bstat)
		if err != nil {
			log.Error("s.income.BatchInsertBgmIncomeStat error(%v)", err)
			return
		}
		log.Info("insert bgm_income_statis : %d", len(bstat))
		return
	})

	// up_income_statis
	writeGroup.Go(func() (err error) {
		err = s.income.upIncomeStatSvr.BatchInsertUpIncomeStat(c, ustat)
		if err != nil {
			log.Error("s.income.BatchInsertUpIncomeStat error(%v)", err)
			return
		}
		log.Info("insert up_income_statis : %d", len(ustat))
		return
	})

	// up_account batch insert
	writeGroup.Go(func() (err error) {
		err = s.income.upAccountSvr.BatchInsertUpAccount(c, accs)
		if err != nil {
			log.Error("s.income.BatchInsertUpAccount error(%v)", err)
			return
		}
		log.Info("insert up_account : %d", len(accs))
		return
	})

	// up account single update
	writeGroup.Go(func() (err error) {
		err = s.income.upAccountSvr.UpdateUpAccount(c, accs)
		if err != nil {
			log.Error("s.income.UpdateUpAccount error(%v)", err)
			return
		}
		log.Info("update up_account : %d", len(accs))
		return
	})

	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}

func signedBgmFilter(m map[int64]*model.Signed, date time.Time) BgmFilter {
	return func(charge *model.AvCharge, bgm *model.BGM) bool {
		if up, ok := m[bgm.MID]; ok {
			if charge.Date.Time().Before(up.SignedAt.Time()) {
				return true
			}
			if (up.AccountState == 5 || up.AccountState == 6) && up.QuitAt.Time().Before(date.AddDate(0, 0, 1)) {
				return true
			}
		} else {
			return true
		}
		return false
	}
}

func signedAvFilter(m map[int64]*model.Signed, date time.Time) AvFilter {
	return func(charge *model.AvCharge) bool {
		if up, ok := m[charge.MID]; ok {
			if charge.UploadTime.Time().Before(up.SignedAt.Time()) {
				return true
			}
			if (up.AccountState == 5 || up.AccountState == 6) && up.QuitAt.Time().Before(date.AddDate(0, 0, 1)) {
				return true
			}
		} else {
			return true
		}
		return false
	}
}

func signedColumnFilter(m map[int64]*model.Signed, date time.Time) ColumnFilter {
	return func(charge *model.ColumnCharge) bool {
		if up, ok := m[charge.MID]; ok {
			if charge.UploadTime.Time().Before(up.SignedAt.Time()) {
				return true
			}
			if (up.AccountState == 5 || up.AccountState == 6) && up.QuitAt.Time().Before(date.AddDate(0, 0, 1)) {
				return true
			}
		} else {
			return true
		}
		return false
	}
}

func avFilter(m map[int64]bool) AvFilter {
	return func(charge *model.AvCharge) bool {
		return m[charge.AvID]
	}
}

func columnFilter(m map[int64]bool) ColumnFilter {
	return func(charge *model.ColumnCharge) bool {
		return m[charge.ArticleID]
	}
}

func (s *Service) avToBubbleRatio(bubbleMeta map[int64][]int) map[int64]float64 {
	var (
		res         = make(map[int64]float64)
		typeToRatio = make(map[int]float64)
		chooseBType = func(bTypes []int) (bType int) {
			if len(bTypes) == 1 {
				return bTypes[0]
			}
			sort.Slice(bTypes, func(i, j int) bool {
				bti, btj := bTypes[i], bTypes[j]
				if typeToRatio[bti] == typeToRatio[btj] {
					return bti < btj
				}
				return typeToRatio[bti] < typeToRatio[btj]
			})
			return bTypes[0]
		}
	)

	for _, v := range s.conf.Bubble.BRatio {
		typeToRatio[v.BType] = v.Ratio
	}
	for avID, bTypes := range bubbleMeta {
		bType := chooseBType(bTypes)
		res[avID] = typeToRatio[bType]
	}
	return res
}

// BgmFilter av charge filter
type BgmFilter func(*model.AvCharge, *model.BGM) bool

// AvFilter av charge filter
type AvFilter func(*model.AvCharge) bool

// ColumnFilter column charge filter
type ColumnFilter func(*model.ColumnCharge) bool

// ChargeRegulator regulates av charge
type ChargeRegulator func(*model.AvCharge)

// UpdateBusinessIncome ..
func (s *Service) UpdateBusinessIncome(c context.Context, date string) (err error) {
	return s.income.UpdateBusinessIncomeByDate(c, date)
}
