package income

import (
	"bytes"
	"context"
	"strconv"
	"time"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
	xtime "go-common/library/time"
)

// Income income service
type Income struct {
	avIncomeSvr         *AvIncomeSvr
	upIncomeSvr         *UpIncomeSvr
	avIncomeStatSvr     *AvIncomeStatSvr
	bgmIncomeStatSvr    *BgmIncomeStatSvr
	columnIncomeStatSvr *ColumnIncomeStatSvr
	upIncomeStatSvr     *UpIncomeStatSvr
	upAccountSvr        *UpAccountSvr
	dateStatisSvr       *DateStatis
	bgmIncomeSvr        *BgmIncomeSvr
	columnIncomeSvr     *ColumnIncomeSvr
}

// NewIncome new income service
func NewIncome(batchSize int, dao *incomeD.Dao) *Income {
	return &Income{
		avIncomeSvr:         NewAvIncomeSvr(dao, batchSize),
		bgmIncomeSvr:        NewBgmIncomeSvr(dao, batchSize),
		columnIncomeSvr:     NewColumnIncomeSvr(dao, batchSize),
		upIncomeSvr:         NewUpIncomeSvr(dao, batchSize),
		avIncomeStatSvr:     NewAvIncomeStatSvr(dao, batchSize),
		bgmIncomeStatSvr:    NewBgmIncomeStatSvr(dao, batchSize),
		columnIncomeStatSvr: NewColumnIncomeStatSvr(dao, batchSize),
		upIncomeStatSvr:     NewUpIncomeStatSvr(dao, batchSize),
		upAccountSvr:        NewUpAccountSvr(dao, batchSize),
		dateStatisSvr:       NewDateStatis(dao),
	}
}

// CalAvIncome cal av income
func (p *Income) CalAvIncome(ch chan []*model.AvCharge, urs map[int64]*model.UpChargeRatio, ars map[int64]*model.ArchiveChargeRatio, filters []AvFilter, signed map[int64]bool) (am map[int64][]*model.AvIncome, total map[int64]*model.UpBusinessIncome) {
	am, total = getClassifyAvIncome(ch, ars, filters)
	classifyAvIncome(am, urs, ars, total, signed)
	return
}

// CalColumnIncome cal column income
func (p *Income) CalColumnIncome(ch chan []*model.ColumnCharge, urs map[int64]*model.UpChargeRatio, ars map[int64]*model.ArchiveChargeRatio, filters []ColumnFilter, signed map[int64]bool) (cm map[int64][]*model.ColumnIncome, total map[int64]*model.UpBusinessIncome) {
	cm, total = getClassifyColumnIncome(ch, ars, filters)
	classifyColumnIncome(cm, urs, ars, total, signed)
	return
}

// CalBgmIncome cal bgm income
func (p *Income) CalBgmIncome(ch chan []*model.AvCharge, bgms map[int64][]*model.BGM, urs map[int64]*model.UpChargeRatio, ars map[int64]*model.ArchiveChargeRatio, avf []AvFilter, exclude BgmFilter, black map[int64]bool, signed map[int64]bool) (bm map[int64]map[int64]map[int64]*model.BgmIncome, total map[int64]*model.UpBusinessIncome) {
	bm, total = getClassifyBgmIncome(ch, bgms, avf, exclude, black, ars)
	classifyBgmIncome(bm, urs, ars, total, signed)
	return
}

func trans2AvIncome(charge *model.AvCharge, incr int64) *model.AvIncome {
	return &model.AvIncome{
		AvID:       charge.AvID,
		MID:        charge.MID,
		TagID:      charge.TagID,
		IsOriginal: charge.IsOriginal,
		UploadTime: charge.UploadTime,
		PlayCount:  charge.TotalPlayCount,
		Date:       charge.Date,
		Income:     incr,
		BaseIncome: charge.IncCharge,
	}
}

func av2BusinessIncome(a *model.AvIncome) *model.UpBusinessIncome {
	return &model.UpBusinessIncome{
		MID:       a.MID,
		Income:    a.Income,
		PlayCount: a.PlayCount,
		AvCount:   1,
		Business:  1,
	}
}

// av income before tax
// total key:mid, value:total_income
func getClassifyAvIncome(ch chan []*model.AvCharge, ratio map[int64]*model.ArchiveChargeRatio, filters []AvFilter) (am map[int64][]*model.AvIncome, total map[int64]*model.UpBusinessIncome) {
	am = make(map[int64][]*model.AvIncome)
	total = make(map[int64]*model.UpBusinessIncome)
	for charges := range ch {
	CHARGE:
		for _, charge := range charges {
			for _, filter := range filters {
				if filter(charge) {
					continue CHARGE
				}
			}
			var incr int64
			if r, ok := ratio[charge.AvID]; ok {
				if r.AdjustType == 0 {
					incr = int64(Round(float64(charge.IncCharge)*Div(float64(r.Ratio), float64(100)), 0))
				} else if r.AdjustType == 1 {
					incr = charge.IncCharge
				}
			} else {
				incr = charge.IncCharge
			}
			avIncome := trans2AvIncome(charge, incr)
			if _, ok := am[charge.MID]; ok {
				am[charge.MID] = append(am[charge.MID], avIncome)
			} else {
				am[charge.MID] = []*model.AvIncome{avIncome}
			}

			if business, ok := total[charge.MID]; ok {
				business.PlayCount += avIncome.PlayCount
				business.AvCount++
				business.Income += avIncome.Income
			} else {
				business := av2BusinessIncome(avIncome)
				total[charge.MID] = business
			}
		}
	}
	return
}

func classifyAvIncome(am map[int64][]*model.AvIncome,
	urs map[int64]*model.UpChargeRatio,
	ars map[int64]*model.ArchiveChargeRatio,
	total map[int64]*model.UpBusinessIncome,
	signed map[int64]bool) {

	for mid, business := range total {
		realIncome := business.Income
		if r, ok := urs[mid]; ok {
			if r.AdjustType == 0 {
				// up主浮动调节,分配到视频业务
				realIncome = int64(Round(float64(business.Income)*Div(float64(r.Ratio), float64(100)), 0))
			}
		}
		// up主浮动调节后视频业务的税
		tax := int64(Round(Tax(Div(float64(realIncome), 100))*100, 0))
		// 税后收入
		netIncome := realIncome - tax
		// update up income
		business.Percent = Div(float64(netIncome), float64(business.Income))
		business.Income = netIncome
		business.Tax = tax
	}

	// 计算每个视频浮动调节后收入
	for mid, as := range am {
		business := total[mid]
		// up主的固定调节收入
		var patchIncome int64
		var c bool
		for _, a := range as {
			avIncome := int64(Round(Mul(float64(a.Income), business.Percent), 0))
			a.Income = avIncome
			a.TaxMoney = int64(Round(Mul(float64(business.Tax), Div(float64(avIncome), float64(business.Income))), 0))

			// 以下是计算up主基础收入
			// original raito: 100
			var o int64 = 100
			if upRatio, ok := urs[mid]; ok {
				if upRatio.AdjustType == 0 {
					o = upRatio.Ratio
					c = true
				}
			}
			avBaseIncome := int64(Round(Div(float64(avIncome), Div(float64(o), 100)), 0))
			a.BaseIncome = avBaseIncome
			if r, ok := ars[a.AvID]; ok {
				if r.AdjustType == 0 {
					// 如果视频收入被加倍过,还原
					originIncome := int64(Round(Div(float64(avBaseIncome), Div(float64(r.Ratio), 100)), 0))
					a.BaseIncome = originIncome
					business.BaseIncome += originIncome
					c = true
				}
				if r.AdjustType == 1 {
					// av 固定调节, 更新av的收入和up主的基础收入
					a.Income += r.Ratio
					// 更新up主的av收入, Income在计算BaseIncome后更新
					// upIncome.AvIncome += r.Ratio
					// business.Income += r.Ratio
					patchIncome += r.Ratio
					business.BaseIncome += avBaseIncome
				}
			} else {
				business.BaseIncome += avBaseIncome
			}
		}

		// 如果没有被加倍过, 那么避免误差将BaseIncome直接置为Income
		if !c {
			business.BaseIncome = business.Income
		}
		// 最后加上该up主av的总固定调节收入
		business.Income += patchIncome
	}

	// + up主的固定调节
	for mid, ratio := range urs {
		if _, ok := signed[mid]; !ok {
			continue
		}
		if ratio.AdjustType == 0 {
			continue
		}

		if business, ok := total[mid]; ok {
			business.Income += ratio.Ratio
		} else {
			total[mid] = &model.UpBusinessIncome{
				MID:      mid,
				Income:   ratio.Ratio,
				Business: 1,
			}
		}
	}
}

/*################################################################## BGM ########################################################################*/

func trans2BgmIncome(b *model.BGM, charge int64, date xtime.Time) *model.BgmIncome {
	return &model.BgmIncome{
		AID:    b.AID,
		SID:    b.SID,
		MID:    b.MID,
		CID:    b.CID,
		Income: charge,
		Date:   date,
	}
}

func bgm2BusinessIncome(b *model.BgmIncome) *model.UpBusinessIncome {
	return &model.UpBusinessIncome{
		MID:      b.MID,
		Income:   b.Income,
		Business: 3,
		BgmCount: map[int64]bool{
			b.SID: true,
		},
	}
}

// bgms map[avid][sid]*model.BGM | bm map[mid]map[sid]map[avid]*model.BgmIncome
func getClassifyBgmIncome(ch chan []*model.AvCharge, bgms map[int64][]*model.BGM, filters []AvFilter, exclude BgmFilter, black map[int64]bool, ratio map[int64]*model.ArchiveChargeRatio) (bm map[int64]map[int64]map[int64]*model.BgmIncome, total map[int64]*model.UpBusinessIncome) {
	bm = make(map[int64]map[int64]map[int64]*model.BgmIncome)
	total = make(map[int64]*model.UpBusinessIncome)
	for charges := range ch {
	CHARGE:
		for _, charge := range charges {
			for _, filter := range filters {
				if filter(charge) {
					continue CHARGE
				}
			}

			if bs, ok := bgms[charge.AvID]; ok {
				bgmCharge := int64(Round(Div(Mul(float64(charge.IncCharge), float64(0.3)), float64(len(bs))), 0))

				for _, b := range bs {
					// if av's bgm is own, continue
					if b.MID == charge.MID {
						continue
					}

					if _, ok := black[b.SID]; ok {
						continue
					}

					if exclude(charge, b) {
						continue
					}

					var incr int64
					if r, ok := ratio[b.SID]; ok {
						if r.AdjustType == 0 {
							incr = int64(Round(float64(bgmCharge)*Div(float64(r.Ratio), float64(100)), 0))
						} else if r.AdjustType == 1 {
							incr = bgmCharge
						}
					} else {
						incr = bgmCharge
					}

					// bm map[mid]map[sid]map[avid]*model.BgmIncome
					var bgmIncome *model.BgmIncome
					if sm, ok := bm[b.MID]; ok {
						if am, ok := sm[b.SID]; ok {
							if bgmIncome, ok = am[b.AID]; ok {
								bgmIncome.Income += incr
							} else {
								bgmIncome = trans2BgmIncome(b, incr, charge.Date)
								am[b.AID] = bgmIncome
							}
						} else {
							bgmIncome = trans2BgmIncome(b, incr, charge.Date)
							sm[b.SID] = map[int64]*model.BgmIncome{
								b.AID: bgmIncome,
							}
						}
					} else {
						bgmIncome = trans2BgmIncome(b, incr, charge.Date)
						// am map[avid]*model.BgmIncome
						am := map[int64]*model.BgmIncome{
							b.AID: bgmIncome,
						}
						// sm map[sid]map[avid]*model.BgmIncome
						sm := map[int64]map[int64]*model.BgmIncome{
							b.SID: am,
						}
						// bm map[mid]map[sid]map[avid]*model.BgmIncome
						bm[b.MID] = sm
					}

					if business, ok := total[b.MID]; ok {
						business.Income += incr
						business.BgmCount[b.SID] = true
					} else {
						business := bgm2BusinessIncome(bgmIncome)
						total[b.MID] = business
					}
				}
			}
		}
	}
	return
}

// 算税并分配收入
func classifyBgmIncome(bm map[int64]map[int64]map[int64]*model.BgmIncome,
	urs map[int64]*model.UpChargeRatio,
	ars map[int64]*model.ArchiveChargeRatio,
	total map[int64]*model.UpBusinessIncome,
	signed map[int64]bool) {
	for mid, business := range total {
		if business.Income == 0 {
			delete(total, mid)
			delete(bm, mid)
			continue
		}
		realIncome := business.Income
		if r, ok := urs[mid]; ok {
			if r.AdjustType == 0 {
				// up主浮动调节,分配到视频业务
				realIncome = int64(Round(float64(business.Income)*Div(float64(r.Ratio), float64(100)), 0))
			}
		}
		// up主浮动调节后视频业务的税
		tax := int64(Round(Tax(Div(float64(realIncome), 100))*100, 0))
		// 税后收入
		netIncome := realIncome - tax
		// update up income
		business.Percent = Div(float64(netIncome), float64(business.Income))
		business.Income = netIncome
		business.Tax = tax
	}
	// bm map[mid]map[sid]map[avid]*model.BgmIncome
	for mid, sm := range bm {
		business := total[mid]
		var c bool
		var patchIncome int64
		for sid, bs := range sm {
			var dailyTotalIncome int64
			for _, b := range bs {
				income := int64(Round(Mul(float64(b.Income), business.Percent), 0))
				b.Income = income
				b.TaxMoney = int64(Round(Mul(float64(business.Tax), Div(float64(income), float64(business.Income))), 0))
				dailyTotalIncome += b.Income

				// 以下是计算up主基础收入
				// original raito: 100
				var o int64 = 100
				if upRatio, ok := urs[mid]; ok {
					if upRatio.AdjustType == 0 {
						o = upRatio.Ratio
						c = true
					}
				}

				bgmBaseIncome := int64(Round(Div(float64(income), Div(float64(o), 100)), 0))
				b.BaseIncome = bgmBaseIncome
				if r, ok := ars[sid]; ok {
					if r.AdjustType == 0 {
						// 如果bgm收入被加倍过,还原
						originIncome := int64(Round(Div(float64(bgmBaseIncome), Div(float64(r.Ratio), 100)), 0))
						b.BaseIncome = originIncome
						business.BaseIncome += originIncome
						c = true
					} else {
						business.BaseIncome += bgmBaseIncome
					}
				} else {
					business.BaseIncome += bgmBaseIncome
				}
			}
			if r, ok := ars[sid]; ok {
				if r.AdjustType == 1 {
					patchIncome += r.Ratio
					dailyTotalIncome += r.Ratio
				}
			}
			// update bgm daily total income
			for _, b := range bs {
				b.DailyTotalIncome = dailyTotalIncome
			}
		}

		if !c {
			business.BaseIncome = business.Income
		}
		// 最后加上该up主bgm的总固定调节收入
		business.Income += patchIncome
	}

	// + up主的固定调节
	for mid, ratio := range urs {
		if _, ok := signed[mid]; !ok {
			continue
		}
		if ratio.AdjustType == 0 {
			continue
		}

		if business, ok := total[mid]; ok {
			business.Income += ratio.Ratio
		} else {
			total[mid] = &model.UpBusinessIncome{
				MID:      mid,
				Income:   ratio.Ratio,
				Business: 3,
			}
		}
	}
}

/*###################################################### Column #######################################################*/

func trans2ColumnIncome(charge *model.ColumnCharge, incr int64) (c *model.ColumnIncome) {
	return &model.ColumnIncome{
		ArticleID:  charge.ArticleID,
		Title:      charge.Title,
		MID:        charge.MID,
		TagID:      charge.TagID,
		UploadTime: charge.UploadTime,
		ViewCount:  charge.IncViewCount,
		Date:       charge.Date,
		Income:     incr,
		BaseIncome: charge.IncCharge,
	}
}

// business type 2: column
func column2BusinessIncome(c *model.ColumnIncome) *model.UpBusinessIncome {
	return &model.UpBusinessIncome{
		MID:         c.MID,
		Income:      c.Income,
		ViewCount:   c.ViewCount,
		ColumnCount: 1,
		Business:    2,
	}
}

func getClassifyColumnIncome(ch chan []*model.ColumnCharge, ratio map[int64]*model.ArchiveChargeRatio, filters []ColumnFilter) (cm map[int64][]*model.ColumnIncome, total map[int64]*model.UpBusinessIncome) {
	cm = make(map[int64][]*model.ColumnIncome)
	total = make(map[int64]*model.UpBusinessIncome)

	for charges := range ch {
	CHARGE:
		for _, charge := range charges {
			for _, filter := range filters {
				if filter(charge) {
					continue CHARGE
				}
			}

			var incr int64
			if r, ok := ratio[charge.ArticleID]; ok {
				if r.AdjustType == 0 {
					incr = int64(Round(float64(charge.IncCharge)*Div(float64(r.Ratio), float64(100)), 0))
				} else if r.AdjustType == 1 {
					incr = charge.IncCharge
				}
			} else {
				incr = charge.IncCharge
			}

			columnIncome := trans2ColumnIncome(charge, incr)
			if _, ok := cm[charge.MID]; ok {
				cm[charge.MID] = append(cm[charge.MID], columnIncome)
			} else {
				cm[charge.MID] = []*model.ColumnIncome{columnIncome}
			}

			if business, ok := total[charge.MID]; ok {
				business.ViewCount += columnIncome.ViewCount
				business.ColumnCount++
				business.Income += columnIncome.Income
			} else {
				business := column2BusinessIncome(columnIncome)
				total[charge.MID] = business
			}
		}
	}
	return
}

// calculate column income and archive income
func classifyColumnIncome(cm map[int64][]*model.ColumnIncome,
	urs map[int64]*model.UpChargeRatio,
	ars map[int64]*model.ArchiveChargeRatio,
	total map[int64]*model.UpBusinessIncome,
	signed map[int64]bool) {
	for mid, business := range total {
		realIncome := business.Income
		if r, ok := urs[mid]; ok {
			if r.AdjustType == 0 {
				// up主浮动调节,分配到视频业务
				realIncome = int64(Round(float64(business.Income)*Div(float64(r.Ratio), float64(100)), 0))
			}
		}
		// up主浮动调节后视频业务的税
		tax := int64(Round(Tax(Div(float64(realIncome), 100))*100, 0))
		// 税后收入
		netIncome := realIncome - tax
		// update up income
		business.Percent = Div(float64(netIncome), float64(business.Income))
		business.Income = netIncome
		business.Tax = tax
	}
	for mid, cs := range cm {
		business := total[mid]
		var x bool
		var patchIncome int64
		for _, c := range cs {
			columnIncome := int64(Round(Mul(float64(c.Income), business.Percent), 0))
			c.Income = columnIncome
			c.TaxMoney = int64(Round(Mul(float64(business.Tax), Div(float64(columnIncome), float64(business.Income))), 0))

			// 以下是计算up主基础收入
			// original raito: 100
			var o int64 = 100
			if upRatio, ok := urs[mid]; ok {
				if upRatio.AdjustType == 0 {
					o = upRatio.Ratio
					x = true
				}
			}
			columnBaseIncome := int64(Round(Div(float64(columnIncome), Div(float64(o), 100)), 0))
			c.BaseIncome = columnBaseIncome
			if r, ok := ars[c.ArticleID]; ok {
				if r.AdjustType == 0 {
					// 如果视频收入被加倍过,还原
					originIncome := int64(Round(Div(float64(columnBaseIncome), Div(float64(r.Ratio), 100)), 0))
					c.BaseIncome = originIncome
					business.BaseIncome += originIncome
					x = true
				}
				if r.AdjustType == 1 {
					// column 固定调节, 更新column 的收入和up主的基础收入
					c.Income += r.Ratio
					// 更新up主的column 收入, Income在计算BaseIncome后更新
					patchIncome += r.Ratio
					business.BaseIncome += columnBaseIncome
				}
			} else {
				business.BaseIncome += columnBaseIncome
			}
		}
		if !x {
			business.BaseIncome = business.Income
		}
		// 最后加上该up主column的总固定调节收入
		business.Income += patchIncome
	}

	// + up主的固定调节
	for mid, ratio := range urs {
		if _, ok := signed[mid]; !ok {
			continue
		}
		if ratio.AdjustType == 0 {
			continue
		}

		if business, ok := total[mid]; ok {
			business.Income += ratio.Ratio
		} else {
			total[mid] = &model.UpBusinessIncome{
				MID:      mid,
				Income:   ratio.Ratio,
				Business: 2,
			}
		}
	}
}

/**######################################################## UpIncome ###################################################################**/

// CalUpIncome calculate upIncome
func (p *Income) CalUpIncome(ch chan map[int64]*model.UpBusinessIncome, date time.Time) (m map[int64]*model.UpIncome) {
	defer close(ch)
	m = make(map[int64]*model.UpIncome)
	var finished int
	for bm := range ch {
		for mid, business := range bm {
			if upIncome, ok := m[mid]; ok {
				if business.Business == 1 {
					upIncome.AvCount = business.AvCount
					upIncome.PlayCount = business.PlayCount
					upIncome.AvIncome = business.Income
					upIncome.AvBaseIncome = business.BaseIncome
					upIncome.AvTax = business.Tax
				}
				if business.Business == 2 {
					upIncome.ColumnCount = business.ColumnCount
					upIncome.ColumnIncome = business.Income
					upIncome.ColumnBaseIncome = business.BaseIncome
					upIncome.ColumnTax = business.Tax
				}
				if business.Business == 3 {
					upIncome.BgmIncome = business.Income
					upIncome.BgmBaseIncome = business.BaseIncome
					upIncome.BgmTax = business.Tax
					upIncome.BgmCount = int64(len(business.BgmCount))
				}
				upIncome.TaxMoney += business.Tax
				upIncome.BaseIncome += business.BaseIncome
				upIncome.Income += business.Income
			} else {
				var upIncome *model.UpIncome
				if business.Business == 1 {
					upIncome = &model.UpIncome{
						AvCount:      business.AvCount,
						PlayCount:    business.PlayCount,
						AvIncome:     business.Income,
						AvBaseIncome: business.BaseIncome,
						AvTax:        business.Tax,
					}
				}
				if business.Business == 2 {
					upIncome = &model.UpIncome{
						ColumnCount:      business.ColumnCount,
						ColumnIncome:     business.Income,
						ColumnBaseIncome: business.BaseIncome,
						ColumnTax:        business.Tax,
					}
				}
				if business.Business == 3 {
					upIncome = &model.UpIncome{
						BgmIncome:     business.Income,
						BgmBaseIncome: business.BaseIncome,
						BgmTax:        business.Tax,
						BgmCount:      int64(len(business.BgmCount)),
					}
				}
				upIncome.MID = mid
				upIncome.TaxMoney = business.Tax
				upIncome.BaseIncome = business.BaseIncome
				upIncome.Income = business.Income
				upIncome.Date = xtime.Time(date.Unix())
				m[mid] = upIncome
			}
		}
		finished++
		if finished == 3 {
			break
		}
	}
	return
}

// IncomeStat income statistics
func (p *Income) IncomeStat(
	um map[int64]*model.UpIncome,
	am map[int64][]*model.AvIncome,
	bm map[int64]map[int64]map[int64]*model.BgmIncome,
	cm map[int64][]*model.ColumnIncome,
	ustat map[int64]*model.UpIncomeStat,
	astat map[int64]*model.AvIncomeStat,
	bstat map[int64]*model.BgmIncomeStat,
	cstat map[int64]*model.ColumnIncomeStat) {
	for mid, upIncome := range um {
		// update up income total income
		ut, ok := ustat[mid]
		if !ok {
			ut = &model.UpIncomeStat{
				MID:       upIncome.MID,
				DataState: 1, // insert
			}
			ustat[mid] = ut
		} else {
			ut.DataState = 2 // update
		}

		// update up income statis up total income
		// up av total income
		// up column total income
		// up bgm total income
		ut.TotalIncome += upIncome.Income
		upIncome.TotalIncome = ut.TotalIncome

		ut.AvTotalIncome += upIncome.AvIncome
		upIncome.AvTotalIncome = ut.AvTotalIncome

		ut.ColumnTotalIncome += upIncome.ColumnIncome
		upIncome.ColumnTotalIncome = ut.ColumnTotalIncome

		ut.BgmTotalIncome += upIncome.BgmIncome
		upIncome.BgmTotalIncome = ut.BgmTotalIncome

		// update av income total income
		for _, a := range am[mid] {
			at, ok := astat[a.AvID]
			if !ok {
				at = &model.AvIncomeStat{
					AvID:       a.AvID,
					MID:        a.MID,
					TagID:      a.TagID,
					IsOriginal: a.IsOriginal,
					UploadTime: a.UploadTime,
					DataState:  1, // insert
				}
				astat[a.AvID] = at
			} else {
				at.DataState = 2 // update
			}
			at.TotalIncome += a.Income
			a.TotalIncome = at.TotalIncome
		}

		// update bgmIncome total income
		for sid, am := range bm[mid] {
			bt, ok := bstat[sid]
			if !ok {
				bt = &model.BgmIncomeStat{
					SID:       sid,
					DataState: 1, // insert
				}
				bstat[sid] = bt
			} else {
				bt.DataState = 2 // update
			}
			for _, b := range am {
				bt.TotalIncome += b.DailyTotalIncome
				break
			}
			for _, b := range am {
				b.TotalIncome = bt.TotalIncome
			}
		}

		// update columnIncome total income
		for _, c := range cm[mid] {
			ct, ok := cstat[c.ArticleID]
			if !ok {
				ct = &model.ColumnIncomeStat{
					ArticleID:  c.ArticleID,
					Title:      c.Title,
					TagID:      c.TagID,
					MID:        c.MID,
					UploadTime: c.UploadTime,
					DataState:  1, // insert
				}
				cstat[c.ArticleID] = ct
			} else {
				ct.DataState = 2 // update
			}
			ct.TotalIncome += c.Income
			c.TotalIncome = ct.TotalIncome
		}
	}
}

// PurgeUpAccount purge up account
func (p *Income) PurgeUpAccount(date time.Time, accs map[int64]*model.UpAccount, um map[int64]*model.UpIncome) {
	fd := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	last := fd.AddDate(0, 0, -1).Format("2006-01")

	// TODO check
	for mid, upIncome := range um {
		// update old up account
		if upAccount, ok := accs[mid]; ok {
			upAccount.TotalIncome += upIncome.Income
			if upAccount.WithdrawDateVersion == last {
				upAccount.TotalUnwithdrawIncome += upIncome.Income
			}
			upAccount.DataState = 2
		} else { // append new up account
			accs[mid] = &model.UpAccount{
				MID:                   mid,
				HasSignContract:       1,
				TotalIncome:           upIncome.Income,
				TotalUnwithdrawIncome: upIncome.Income,
				WithdrawDateVersion:   last,
				DataState:             1,
			}
		}
	}
}

/************************************************************ FOR HISTORY ****************************************************************/

// UpdateBusinessIncomeByDate update business income by date
func (p *Income) UpdateBusinessIncomeByDate(c context.Context, date string) (err error) {
	ustat, err := p.upIncomeStatSvr.UpIncomeStat(c, int64(_limitSize))
	if err != nil {
		return
	}
	return p.updateBusinessIncomeByDate(c, date, ustat)
}

func (p *Income) updateBusinessIncomeByDate(c context.Context, date string, ustat map[int64]*model.UpIncomeStat) (err error) {
	us, err := p.businessTotalIncome(c, ustat, date)
	if err != nil {
		return
	}
	err = p.batchUpdateUpIncome(c, us)
	if err != nil {
		return
	}
	return p.batchUpdateUpIncomeStat(c, ustat)
}

// for history data m: map[mid]map[date]*model.UpIncome
func (p *Income) businessTotalIncome(c context.Context, ustat map[int64]*model.UpIncomeStat, date string) (m []*model.UpIncome, err error) {
	var id int64
	for {
		var ups []*model.UpIncome
		ups, err = p.upIncomeSvr.dao.GetUpIncomeTable(c, "up_income", date, id, 2000)
		if err != nil {
			return
		}
		for _, up := range ups {
			if ut, ok := ustat[up.MID]; ok {
				ut.AvTotalIncome += up.AvIncome
				up.AvTotalIncome = ut.AvTotalIncome

				ut.ColumnTotalIncome += up.ColumnIncome
				up.ColumnTotalIncome = ut.ColumnTotalIncome

				ut.BgmTotalIncome += up.BgmIncome
				up.BgmTotalIncome = ut.BgmTotalIncome
				m = append(m, up)
			}
		}
		if len(ups) < 2000 {
			break
		}
		id = ups[len(ups)-1].ID
	}
	return
}

// BatchUpdateUpIncome insert up_income batch
func (p *Income) batchUpdateUpIncome(c context.Context, us []*model.UpIncome) (err error) {
	var (
		buff    = make([]*model.UpIncome, p.upIncomeSvr.batchSize)
		buffEnd = 0
	)

	for _, u := range us {
		buff[buffEnd] = u
		buffEnd++

		if buffEnd >= p.upIncomeSvr.batchSize {
			values := businessValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.upIncomeSvr.dao.FixInsertUpIncome(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := businessValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.upIncomeSvr.dao.FixInsertUpIncome(c, values)
	}
	return
}

func businessValues(us []*model.UpIncome) (values string) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + u.Date.Time().Format(_layout) + "'")
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

// BatchUpdateUpIncome insert up_income batch
func (p *Income) batchUpdateUpIncomeStat(c context.Context, us map[int64]*model.UpIncomeStat) (err error) {
	var (
		buff    = make([]*model.UpIncomeStat, p.upIncomeStatSvr.batchSize)
		buffEnd = 0
	)

	for _, u := range us {
		buff[buffEnd] = u
		buffEnd++

		if buffEnd >= p.upIncomeStatSvr.batchSize {
			values := businessStatValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.upIncomeStatSvr.dao.FixInsertUpIncomeStat(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := businessStatValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.upIncomeStatSvr.dao.FixInsertUpIncomeStat(c, values)
	}
	return
}

func businessStatValues(ustat []*model.UpIncomeStat) (values string) {
	var buf bytes.Buffer
	for _, u := range ustat {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
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
