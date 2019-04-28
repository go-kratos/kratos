package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/member/api"
	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_realnameSalt = "biliidentification@#$%^&*()(*&^%$#"
	_imgPrefix    = "/idenfiles/"
	_imgSuffix    = ".txt"

	_idCard18Regexp = regexp.MustCompile(`^\d{17}[\d|x]$`)
	_idCard15Regexp = regexp.MustCompile(`^\d{15}$`)
	_genderUnknown  = "unknown"
	_genderMale     = "male"
	_genderFemale   = "female"
)

// RealnameStatus 返回实名状态
func (s *Service) RealnameStatus(c context.Context, mid int64) (status model.RealnameStatus, err error) {
	var info *model.RealnameCacheInfo
	if info, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	if info.Status.IsPass() {
		status = model.RealnameStatusTrue
	} else {
		status = model.RealnameStatusFalse
	}
	return
}

// RealnameApplyStatus 返回实名申请状态
func (s *Service) RealnameApplyStatus(c context.Context, mid int64) (info *model.RealnameApplyStatusInfo, err error) {
	var (
		ci *model.RealnameCacheInfo
	)
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	info = &model.RealnameApplyStatusInfo{
		Status:   ci.Status,
		Remark:   ci.Reason,
		Realname: ci.Realname,
		Card:     ci.RealCard,
	}
	return
}

// RealnameBrief is
func (s *Service) RealnameBrief(c context.Context, mid int64) (brief *model.RealnameBrief, err error) {
	var (
		ci *model.RealnameCacheInfo
	)
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	brief = &model.RealnameBrief{
		Realname: ci.Realname,
		Card:     ci.RealCard,
		CardType: ci.CardType,
	}
	if ci.Status.IsPass() {
		brief.Status = model.RealnameStatusTrue
	} else {
		brief.Status = model.RealnameStatusFalse
	}
	return
}

// RealnameDetail is
func (s *Service) RealnameDetail(c context.Context, mid int64) (detail *model.RealnameDetail, err error) {
	var (
		ci *model.RealnameCacheInfo
	)
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	detail = &model.RealnameDetail{
		RealnameBrief: &model.RealnameBrief{
			Realname: ci.Realname,
			Card:     ci.RealCard,
			CardType: ci.CardType,
		},
		Gender: _genderUnknown,
	}
	if ci.Status.IsPass() {
		detail.Status = model.RealnameStatusTrue
	} else {
		detail.Status = model.RealnameStatusFalse
		return
	}

	// ger gender
	if detail.CardType == model.RealnameCardTypeIdentity {
		if _, detail.Gender, err = ParseIdentity(detail.Card); err != nil {
			log.Info("%+v", err)
			err = nil
			detail.Gender = _genderUnknown
		}
	}
	//get hand img
	if ci.Channel == model.RealnameChannelMain {
		var (
			apply *model.RealnameApply
			img   *model.RealnameApplyImage
		)
		if apply, err = s.realnameApply(c, mid); err != nil {
			return
		}
		if apply.Status.IsPass() && apply.HandIMG > 0 {
			if img, err = s.mbDao.RealnameApplyIMG(c, apply.HandIMG); err != nil {
				return
			}
			detail.HandIMG = generateIMGURL(img.IMGData)
		}
	}
	return
}

func generateIMGURL(token string) string {
	token = strings.TrimPrefix(token, "/idenfiles/")
	token = strings.TrimSuffix(token, ".txt")
	return fmt.Sprintf(conf.Conf.RealnameProperty.IMGURLTemplate, token)
}

// realnameInfo return info always not nil
func (s *Service) realnameInfo(c context.Context, mid int64) (info *model.RealnameCacheInfo, err error) {
	defer func() {
		log.Info("realname info mid : %d , info : %+v", mid, info)
	}()
	var (
		cacheFlag = true
	)
	if info, err = s.mbDao.RealnameInfoCache(c, mid); err != nil {
		log.Error("%+v", err)
		err = nil
		info = nil
		cacheFlag = false
	}
	// cache hit
	if info != nil {
		return
	}
	var dbInfo *model.RealnameInfo
	// cache miss
	if dbInfo, err = s.mbDao.RealnameInfo(c, mid); err != nil {
		return
	}
	if dbInfo == nil {
		info = &model.RealnameCacheInfo{
			RealnameInfo: &model.RealnameInfo{
				MID:    mid,
				Status: model.RealnameApplyStatusNone,
			},
		}
	} else {
		info = &model.RealnameCacheInfo{
			RealnameInfo: dbInfo,
		}
		if info.Status.IsPass() && dbInfo.Card != "" {
			var cardCodeBytes []byte
			if cardCodeBytes, err = s.realnameCryptor.CardDecrypt([]byte(dbInfo.Card)); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			info.RealCard = string(cardCodeBytes)
		}
	}
	if cacheFlag {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			if theErr := s.mbDao.SetRealnameInfo(ctx, mid, info); theErr != nil {
				log.Error("%+v", theErr)
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	}
	return
}

func (s *Service) realnameAlipayApply(c context.Context, mid int64) (info *model.RealnameAlipayApply, err error) {
	if info, err = s.mbDao.RealnameAlipayApply(c, mid); err != nil {
		return
	}
	if info == nil {
		info = &model.RealnameAlipayApply{
			MID:    mid,
			Status: model.RealnameApplyStatusNone,
		}
	}
	return
}

// checkCardDup 检查证件重复
func (s *Service) checkCardDup(ctx context.Context, cardMD5 string) (err error) {
	var (
		rawInfo *model.RealnameInfo
	)
	if rawInfo, err = s.mbDao.RealnameInfoByCard(ctx, cardMD5); err != nil {
		return
	}
	log.Info("checkCardDup cardMD5(%s) rawInfo(%+v)", cardMD5, rawInfo)
	if rawInfo != nil && (rawInfo.Status == model.RealnameApplyStatusPending || rawInfo.Status == model.RealnameApplyStatusPass) {
		err = ecode.RealnameCardBindAlready
		return
	}
	return
}

func (s *Service) realnameApply(c context.Context, mid int64) (info *model.RealnameApply, err error) {
	if info, err = s.mbDao.RealnameApply(c, mid); err != nil {
		return
	}
	if info == nil {
		info = &model.RealnameApply{
			MID:    mid,
			Status: model.RealnameApplyStatusNone,
		}
	}
	return
}

// RealnameTelCapture is
func (s *Service) RealnameTelCapture(c context.Context, mid int64) (capture int, err error) {
	var (
		times int
	)
	if times, err = s.mbDao.RealnameCaptureTimesCache(c, mid); err != nil {
		return
	}
	if times < 0 {
		if err = s.mbDao.SetRealnameCaptureTimes(c, mid, 0); err != nil {
			return
		}
	}
	if times > 5 {
		err = ecode.RealnameCaptureSendTooMany
		return
	}
	capture = rand.Intn(900000) + 100000
	if err = s.mbDao.SendCapture(c, mid, capture); err != nil {
		return
	}
	if err = s.mbDao.SetRealnameCaptureCode(c, mid, capture); err != nil {
		return
	}
	if err = s.mbDao.IncreaseRealnameCaptureTimes(c, mid); err != nil {
		return
	}
	if err = s.mbDao.DeleteRealnameCaptureErrTimes(c, mid); err != nil {
		return
	}
	return
}

// RealnameTelCaptureCheck check a tel capture code from mid
func (s *Service) RealnameTelCaptureCheck(c context.Context, mid int64, captureCode int) (err error) {
	var (
		serverCode int
	)
	if serverCode, err = s.mbDao.RealnameCaptureCodeCache(c, mid); err != nil {
		return
	}
	if serverCode < 0 {
		err = ecode.RealnameCaptureInvalid
		return
	}
	if captureCode != serverCode {
		err = ecode.RealnameCaptureErr
		return
	}
	return
}

// RealnameApply is
func (s *Service) RealnameApply(c context.Context, mid int64, captureCode int, realname string, cardType int8, cardCode string, country int16, handIMGToken, frontIMGToken, backIMGToken string) (err error) {
	if err = s.checkCapture(c, mid, captureCode); err != nil {
		return
	}
	var (
		ci      *model.RealnameCacheInfo
		cardMD5 = s.cardMD5(cardCode, int(cardType), int(country))
	)
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	if ci != nil && (ci.Status == model.RealnameApplyStatusPending || ci.Status == model.RealnameApplyStatusPass) {
		err = ecode.RealnameApplyAlready
		return
	}
	if err = s.checkCardDup(c, cardMD5); err != nil {
		return
	}
	var (
		handIMG = &model.RealnameApplyImage{
			IMGData: _imgPrefix + handIMGToken + _imgSuffix,
			CTime:   time.Now(),
		}
		frontIMG = &model.RealnameApplyImage{
			IMGData: _imgPrefix + frontIMGToken + _imgSuffix,
			CTime:   time.Now(),
		}
		backIMG = &model.RealnameApplyImage{
			IMGData: _imgPrefix + backIMGToken + _imgSuffix,
			CTime:   time.Now(),
		}
	)
	// if handIMG.ID, err = s.mbDao.InsertRealnameApplyImg(c, handIMG); err != nil {
	// 	return
	// }
	// if frontIMG.ID, err = s.mbDao.InsertRealnameApplyImg(c, frontIMG); err != nil {
	// 	return
	// }
	// if backIMG.ID, err = s.mbDao.InsertRealnameApplyImg(c, backIMG); err != nil {
	// 	return
	// }
	if handIMG.IMGData != "" {
		if handIMG.ID, err = s.mbDao.InsertOldRealnameApplyImg(c, handIMG); err != nil {
			return
		}
	}
	if frontIMG.IMGData != "" {
		if frontIMG.ID, err = s.mbDao.InsertOldRealnameApplyImg(c, frontIMG); err != nil {
			return
		}
	}
	if backIMG.IMGData != "" {
		if backIMG.ID, err = s.mbDao.InsertOldRealnameApplyImg(c, backIMG); err != nil {
			return
		}
	}

	var (
		encryptedCard []byte
	)

	// 兼容新老版本 account-interface 证件号
	if strings.HasSuffix(cardCode, "=") {
		encryptedCard = []byte(cardCode)
	} else {
		if cardType == model.RealnameCardTypeIdentity && !s.isIDCard(cardCode) {
			err = ecode.RealnameCardNumErr
			return
		}
		if encryptedCard, err = s.realnameCryptor.CardEncrypt([]byte(cardCode)); err != nil {
			return
		}
	}

	// 4. insert apply
	var (
		daoApply = &model.RealnameApply{
			MID:      mid,
			Realname: realname,
			Country:  country,
			CardType: cardType,
			CardNum:  string(encryptedCard),
			CardMD5:  cardMD5,
			HandIMG:  int(handIMG.ID),
			FrontIMG: int(frontIMG.ID),
			BackIMG:  int(backIMG.ID),
			Status:   model.RealnameApplyStatusPending,
			CTime:    time.Now(),
		}
	)

	if daoApply.ID, err = s.mbDao.InsertOldRealnameApply(c, daoApply); err != nil {
		return
	}
	// if err = s.mbDao.InsertRealnameApply(c, daoApply); err != nil {
	// 	return
	// }
	fanoutErr := s.cache.Do(c, func(ctx context.Context) {
		if theErr := s.mbDao.DeleteRealnameCaptureCode(ctx, mid); theErr != nil {
			log.Error("%+v", err)
		}
		if theErr := s.mbDao.DeleteRealnameInfo(ctx, mid); theErr != nil {
			log.Error("%+v", err)
		}
	})
	if fanoutErr != nil {
		log.Error("fanout do err(%+v)", fanoutErr)
	}
	return
}

func (s *Service) isIDCard(card string) bool {
	if !_idCard15Regexp.MatchString(card) && !_idCard18Regexp.MatchString(card) {
		return false
	}
	return true
}

// 检查手机验证码合法性
func (s *Service) checkCapture(c context.Context, mid int64, capture int) (err error) {
	var (
		serverCode int
		errTimes   int
	)
	// 1. check err times
	if errTimes, err = s.mbDao.RealnameCaptureErrTimesCache(c, mid); err != nil {
		return
	}
	if errTimes > 3 {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			if theErr := s.mbDao.DeleteRealnameCaptureCode(ctx, mid); theErr != nil {
				log.Error("%+v", err)
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
		err = ecode.RealnameCaptureErrTooMany
		return
	}
	// 2. check capture code
	if serverCode, err = s.mbDao.RealnameCaptureCodeCache(c, mid); err != nil {
		return
	}
	if serverCode < 0 {
		err = ecode.RealnameCaptureInvalid
		return
	}
	if capture != serverCode {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			if errTimes < 0 {
				if theErr := s.mbDao.SetRealnameCaptureErrTimes(ctx, mid, 0); theErr != nil {
					log.Error("%+v", theErr)
					return
				}
			}
			if theErr := s.mbDao.IncreaseRealnameCaptureErrTimes(ctx, mid); theErr != nil {
				log.Error("%+v", theErr)
				return
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
		err = ecode.RealnameCaptureErr
		return
	}
	return
}

func (s *Service) cardMD5(card string, cardType int, country int) (md5Hex string) {
	var (
		lowerCode = strings.ToLower(card)
		key       = fmt.Sprintf("%s_%s_%d_%d", _realnameSalt, lowerCode, cardType, country)
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

// RealnameAdult is
func (s *Service) RealnameAdult(c context.Context, mid int64) (t model.RealnameAdultType, err error) {
	var (
		ci *model.RealnameCacheInfo
	)
	t = model.RealnameAdultTypeUnknown
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	if !ci.Status.IsPass() {
		return
	}
	if ci.CardType != model.RealnameCardTypeIdentity || ci.Country != model.RealnameCountryChina {
		return
	}
	var (
		birthday      time.Time
		flag          bool
		cardCodeBytes []byte
	)
	if cardCodeBytes, err = s.realnameCryptor.CardDecrypt([]byte(ci.Card)); err != nil {
		return
	}
	// check birthday
	if birthday, _, err = ParseIdentity(string(cardCodeBytes)); err != nil {
		log.Error("mid : %d , identification card err : %+v", mid, err)
		err = nil
		return
	}
	if flag, err = isAdult(birthday, time.Now()); err != nil {
		log.Error("mid : %d , identification card err : %+v", mid, err)
		err = nil
		return
	}
	if flag {
		t = model.RealnameAdultTypeTrue
	} else {
		t = model.RealnameAdultTypeFalse
	}
	return
}

// ParseIdentity to birthday and gender
func ParseIdentity(id string) (birthday time.Time, gender string, err error) {
	var (
		ystr, mstr, dstr, gstr string
		y, m, d, g             int
	)
	switch len(id) {
	case 15:
		ystr, mstr, dstr = "19"+id[6:8], id[8:10], id[10:12]
		gstr = id[14:15]
	case 18:
		ystr, mstr, dstr = id[6:10], id[10:12], id[12:14]
		gstr = id[16:17]
	default:
		err = errors.Errorf("identity id invalid : %s", id)
		return
	}
	if y, err = strconv.Atoi(ystr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if m, err = strconv.Atoi(mstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if d, err = strconv.Atoi(dstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if g, err = strconv.Atoi(gstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	birthday = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	if g%2 == 1 {
		gender = _genderMale
	} else {
		gender = _genderFemale
	}
	return
}

func isAdult(birthday time.Time, anchor time.Time) (f bool, err error) {
	var (
		ny, nm, nd = anchor.Year(), anchor.Month(), anchor.Day()
		by, bm, bd = birthday.Year(), birthday.Month(), birthday.Day()
	)
	if anchor.Before(birthday) || ny-by > 150 {
		err = errors.Errorf("birthday invalid : %+v %+v", birthday, anchor)
		return
	}
	// (年份差18 && (现月份 > 生日月份 || (现月份 = 生日月份 && 现day > 生日day)))) || (年份差>18)
	f = (ny-by == 18 && (nm > bm || (nm == bm && nd >= bd))) || (ny-by > 18)
	return
}

// RealnameCheck check realname card
func (s *Service) RealnameCheck(c context.Context, mid int64, cardType int8, cardCode string) (result bool, err error) {
	result = false
	if cardCode == "" {
		return
	}
	var (
		cinfo *model.RealnameCacheInfo
	)
	if cinfo, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	if !cinfo.Status.IsPass() {
		return
	}
	if (cardType >= 0) && (cinfo.CardType != int(cardType)) {
		return
	}
	var (
		cardCodeBytes []byte
	)
	if cardCodeBytes, err = s.realnameCryptor.CardDecrypt([]byte(cinfo.Card)); err != nil {
		return
	}
	//证件进行比较时统一转换为大写进行比较
	upperCode := strings.ToUpper(cardCode)
	if strings.ToUpper(string(cardCodeBytes)) != upperCode {
		return
	}
	result = true
	return
}

// alipay

// RealnameAlipayApply 芝麻认证申请
func (s *Service) RealnameAlipayApply(c context.Context, mid int64, captureCode int, realname, cardCode, imgToken string, bizno string) (err error) {
	if err = s.checkCapture(c, mid, captureCode); err != nil {
		return
	}
	var (
		ci      *model.RealnameCacheInfo
		cardMD5 = s.cardMD5(cardCode, model.RealnameCardTypeIdentity, model.RealnameCountryChina)
	)
	if ci, err = s.realnameInfo(c, mid); err != nil {
		return
	}
	if ci != nil && (ci.Status == model.RealnameApplyStatusPending || ci.Status == model.RealnameApplyStatusPass) {
		err = ecode.RealnameApplyAlready
		return
	}
	if err = s.checkCardDup(c, cardMD5); err != nil {
		return
	}
	var (
		cardEncBytes []byte
		cardEnc      string
	)
	if cardEncBytes, err = s.realnameCryptor.CardEncrypt([]byte(cardCode)); err != nil {
		return
	}
	cardEnc = string(cardEncBytes)
	log.Info("cardcode : %s , enc : %s", cardCode, cardEnc)

	var (
		daoApply = &model.RealnameAlipayApply{
			MID:      mid,
			Realname: realname,
			Card:     cardEnc,
			IMG:      _imgPrefix + imgToken + _imgSuffix,
			Status:   model.RealnameApplyStatusPending,
			Bizno:    bizno,
		}
		daoInfo = &model.RealnameInfo{
			MID:      mid,
			Channel:  model.RealnameChannelAlipay,
			Realname: realname,
			Country:  model.RealnameCountryChina,
			CardType: model.RealnameCardTypeIdentity,
			Card:     cardEnc,
			CardMD5:  cardMD5,
			Status:   model.RealnameApplyStatusPending,
		}
		tx *xsql.Tx
	)
	if tx, err = s.mbDao.BeginTran(c); err != nil {
		return
	}
	if err = s.mbDao.UpsertTxRealnameInfo(c, tx, daoInfo); err != nil {
		tx.Rollback()
		return
	}
	if err = s.mbDao.InsertTxRealnameAlipayApply(c, tx, daoApply); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	fanoutErr := s.cache.Do(c, func(ctx context.Context) {
		if theErr := s.mbDao.DeleteRealnameCaptureCode(ctx, mid); theErr != nil {
			log.Error("%+v", theErr)
		}
		if theErr := s.mbDao.DeleteRealnameInfo(ctx, mid); theErr != nil {
			log.Error("%+v", theErr)
		}
	})
	if fanoutErr != nil {
		log.Error("fanout do err(%+v)", fanoutErr)
	}
	return
}

// RealnameAlipayConfirm confirm a alipay apply
func (s *Service) RealnameAlipayConfirm(c context.Context, mid int64, pass bool, reason string) (err error) {
	var (
		apply *model.RealnameAlipayApply
	)
	if apply, err = s.mbDao.RealnameAlipayApply(c, mid); err != nil {
		return
	}
	if apply == nil {
		err = ecode.RealnameAlipayApplyInvalid
		return
	}
	// //不为pending说明有更早的消息已更新状态，对此次confirm不予处理
	// if apply.Status != model.RealnameApplyStatusPending {
	// 	log.Info("realname alipay confirm mid : %d , pass : %t , reason : %s , but already confirmed : %+v", mid, pass, reason, apply)
	// 	return
	// }
	var (
		tx     *xsql.Tx
		status model.RealnameApplyStatus
	)
	if pass {
		status = model.RealnameApplyStatusPass
	} else {
		status = model.RealnameApplyStatusBack
	}
	if tx, err = s.mbDao.BeginTran(c); err != nil {
		return
	}
	if err = s.mbDao.UpdateTxRealnameAlipayApply(c, tx, apply.ID, status, reason); err != nil {
		tx.Rollback()
		return
	}
	if err = s.mbDao.UpdateTxRealnameInfo(c, tx, mid, status, reason); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	fanoutErr := s.cache.Do(c, func(ctx context.Context) {
		if theErr := s.mbDao.DeleteRealnameInfo(ctx, mid); theErr != nil {
			log.Error("%+v", err)
		}
	})
	if fanoutErr != nil {
		log.Error("fanout do err(%+v)", fanoutErr)
	}
	return
}

// RealnameAlipayBizno return a alipay apply bizno
func (s *Service) RealnameAlipayBizno(c context.Context, mid int64) (bizno string, err error) {
	var (
		apply *model.RealnameAlipayApply
	)
	if apply, err = s.realnameAlipayApply(c, mid); err != nil {
		return
	}
	bizno = apply.Bizno
	return
}

// RealnameStrippedInfo is
func (s *Service) RealnameStrippedInfo(ctx context.Context, mid int64) (*api.RealnameStrippedInfoReply, error) {
	info, err := s.realnameInfo(ctx, mid)
	if err != nil {
		return nil, err
	}
	stripped := &api.RealnameStrippedInfoReply{
		Mid:      info.MID,
		Channel:  int8(info.Channel),
		Country:  int16(info.Country),
		CardType: int8(info.CardType),
		Status:   int8(info.Status),
	}

	adultType, err := s.RealnameAdult(ctx, mid)
	if err != nil {
		log.Error("Failed to get realname adult status with mid: %d: %+v", mid, err)
		adultType = model.RealnameAdultTypeUnknown
	}
	stripped.AdultType = int8(adultType)

	return stripped, nil
}

// MidByRealnameCard is
func (s *Service) MidByRealnameCard(ctx context.Context, arg *api.MidByRealnameCardsReq) (*api.MidByRealnameCardReply, error) {
	md5ToCode := make(map[string]string, len(arg.CardCode))
	cardMD5s := make([]string, 0, len(arg.CardCode))
	for _, code := range arg.CardCode {
		hashed := s.cardMD5(code, int(arg.CardType), int(arg.Country))
		cardMD5s = append(cardMD5s, hashed)
		md5ToCode[hashed] = code
	}

	md5ToMid, err := s.mbDao.MidByRealnameCards(ctx, cardMD5s)
	if err != nil {
		return nil, err
	}

	codeToMid := make(map[string]int64, len(md5ToMid))
	for hashed, mid := range md5ToMid {
		code, ok := md5ToCode[hashed]
		if !ok {
			log.Warn("Failed to get code by md5: %s", hashed)
			continue
		}
		codeToMid[code] = mid
	}
	return &api.MidByRealnameCardReply{CodeToMid: codeToMid}, nil
}
