package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_upLevel = map[int]string{
		0: "光腚激大励",
		1: "奶瓶激大励",
		2: "幼稚园激大励",
		3: "小学生激大励",
		4: "中学生激大励",
		5: "社会人激大励",
		6: "筋肉激大励",
	}

	_downshiftUps = map[int64]bool{
		148834382: true,
		3647169:   true,
		100238104: true,
		42932619:  true,
		331782289: true,
		40305856:  true,
		19419058:  true,
		4397552:   true,
		234891605: true,
		310562646: true,
		11793131:  true,
		316737962: true,
		5326066:   true,
		669622:    true,
		312663583: true,
		34028169:  true,
		288220100: true,
		344838590: true,
		25729281:  true,
		96506011:  true,
		312508662: true,
		8888364:   true,
		238674714: true,
		337061658: true,
		321230684: true,
		337376429: true,
		85738972:  true,
		10094840:  true,
		291595967: true,
		300676862: true,
		298484242: true,
		40426408:  true,
		2336206:   true,
		78839625:  true,
		291682397: true,
		300576557: true,
		61086273:  true,
		318078497: true,
		291051169: true,
		27174089:  true,
		274284147: true,
		306198844: true,
	}
)

// UpYear up year
func (s *Service) UpYear(c context.Context, mid int64) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-up-year:%d", mid)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.upYear(c, mid)
	if err != nil {
		log.Error("s.upYear error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) upYear(c context.Context, mid int64) (data interface{}, err error) {
	up := new(struct {
		Name      string     `json:"name"`
		Face      string     `json:"face"`
		IsJoin    bool       `json:"is_join"`
		SignedAt  xtime.Time `json:"signed_at"`
		FirstTime xtime.Time `json:"first_time"`
		HasIncome bool       `json:"has_income"`
		Title     string     `json:"title"`
		Level     int        `json:"level"`
		TagIncome []int64    `json:"tag_income"`
	})
	defer func() {
		data = up
	}()

	// has signed
	up.SignedAt, err = s.getUpFirstSignedAt(c, mid)
	if err != nil {
		log.Error("s.dao.getUpFirstSignedAt error(%v)", err)
		return
	}
	if up.SignedAt == 0 {
		up.IsJoin = false
		return
	}
	up.IsJoin = true
	// first income
	up.FirstTime, err = s.dao.GetFirstUpIncome(c, mid)
	if err != nil {
		log.Error("s.dao.GetFirstUpIncome error(%v)", err)
		return
	}
	if up.FirstTime != 0 {
		up.HasIncome = true
	}

	earliestTime := xtime.Time(time.Date(2018, 2, 1, 0, 0, 0, 0, time.Local).Unix())
	if up.SignedAt < earliestTime {
		up.SignedAt = earliestTime
	}
	if up.FirstTime < earliestTime {
		up.FirstTime = earliestTime
	}

	// tag income
	var totalIncome int64
	up.TagIncome, totalIncome, err = s.dao.GetUpYearTag(c, mid)
	if err != nil {
		log.Error("s.dao.GetUpYearTag error(%v)", err)
		return
	}

	upInfos, err := s.dao.AccountInfos(c, []int64{mid})
	if err != nil {
		log.Error("s.dao.AccountInfos error(%v)", err)
		return
	}
	if info, ok := upInfos[mid]; ok {
		up.Name = info.Nickname
		up.Face = info.Face
	}

	switch {
	case totalIncome == 0:
		up.Level = 0
	case totalIncome > 0 && totalIncome <= 10000:
		up.Level = 1
	case totalIncome > 10000 && totalIncome <= 100000:
		up.Level = 2
	case totalIncome > 100000 && totalIncome <= 1000000:
		up.Level = 3
	case totalIncome > 1000000 && totalIncome <= 5000000:
		up.Level = 4
	case totalIncome > 5000000 && totalIncome <= 10000000:
		up.Level = 5
	case totalIncome > 10000000:
		up.Level = 6
	}
	if up.Level > 0 && _downshiftUps[mid] {
		up.Level--
	}
	up.Title = _upLevel[up.Level]
	return
}

func (s *Service) getUpFirstSignedAt(c context.Context, mid int64) (signedAt xtime.Time, err error) {
	video, err := s.dao.GetUpSignedAt(c, "up_info_video", mid)
	if err != nil {
		return
	}
	signedAt = video

	column, err := s.dao.GetUpSignedAt(c, "up_info_column", mid)
	if err != nil {
		return
	}
	if column != 0 && (signedAt > column || signedAt == 0) {
		signedAt = column
	}

	bgm, err := s.dao.GetUpSignedAt(c, "up_info_bgm", mid)
	if err != nil {
		return
	}
	if bgm != 0 && (signedAt > bgm || signedAt == 0) {
		signedAt = bgm
	}
	return
}
