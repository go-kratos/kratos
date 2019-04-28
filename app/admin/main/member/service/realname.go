package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go-common/app/admin/main/member/conf"
	"go-common/app/admin/main/member/model"
	memmdl "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

// consts
const (
	_512KiloBytes = 512 * 1024
)

// RealnameList .
func (s *Service) RealnameList(ctx context.Context, arg *model.ArgRealnameList) (list []*model.RespRealnameApply, total int, err error) {
	list = make([]*model.RespRealnameApply, 0)
	switch arg.Channel {
	case model.ChannelMain:
		return s.realnameMainList(ctx, arg)
	case model.ChannelAlipay:
		return s.realnameAlipayList(ctx, arg)
	}
	return
}

func (s *Service) realnameMainList(ctx context.Context, arg *model.ArgRealnameList) (list []*model.RespRealnameApply, total int, err error) {
	var (
		dl []*model.DBRealnameApply
	)
	if arg.Card != "" {
		// 如果查询证件号, 则通过证件号MD5 realname_info 查询到 对应的 mids
		var (
			infos []*model.DBRealnameInfo
			mids  []int64
			md5   = cardMD5(arg.Card, arg.DBCardType(), arg.Country)
		)
		log.Info("realnameMainList card: %s, md5: %s", arg.Card, md5)
		if infos, err = s.dao.RealnameInfoByCardMD5(ctx, md5, arg.State.DBStatus(), model.ChannelMain.DBChannel()); err != nil {
			return
		}
		log.Info("realnameMainList infos : %+v", infos)
		if len(infos) <= 0 {
			return
		}
		for _, i := range infos {
			log.Info("realnameMainList info: %+v", i)
			mids = append(mids, i.MID)
		}
		log.Info("realnameMainList mids: %+v", mids)
		if dl, total, err = s.dao.RealnameMainList(ctx, mids, arg.DBCardType(), arg.DBCountry(), arg.OPName, arg.TSFrom, arg.TSTo, arg.DBState(), arg.PN, arg.PS, arg.IsDesc); err != nil {
			return
		}
	} else {
		var (
			mids []int64
		)
		if arg.MID > 0 {
			mids = append(mids, arg.MID)
		}
		if dl, total, err = s.dao.RealnameMainList(ctx, mids, arg.DBCardType(), arg.DBCountry(), arg.OPName, arg.TSFrom, arg.TSTo, arg.DBState(), arg.PN, arg.PS, arg.IsDesc); err != nil {
			return
		}
	}
	var (
		midMap = make(map[int64]int) // map[mid]count
		mids   []int64
		imgIDs []int64
	)
	// 审核 db 数据解析进 list
	for _, d := range dl {
		midMap[d.MID] = 0
		var (
			r = &model.RespRealnameApply{}
		)
		r.ParseDBMainApply(d)
		imgIDs = append(imgIDs, d.HandIMG, d.FrontIMG, d.BackIMG)
		list = append(list, r)
	}
	// 没有数据则返回
	if len(midMap) <= 0 {
		return
	}
	// 获取实名申请次数
	for mid := range midMap {
		if midMap[mid], err = s.dao.RealnameApplyCount(ctx, mid); err != nil {
			return
		}
		mids = append(mids, mid)
	}
	// 获取mid的昵称 & 等级信息
	var (
		memsArg = &memmdl.ArgMemberMids{
			Mids: mids,
		}
		memMap map[int64]*memmdl.Member
		imgMap map[int64]*model.DBRealnameApplyIMG
	)
	if memMap, err = s.memberRPC.Members(ctx, memsArg); err != nil {
		err = errors.WithStack(err)
		return
	}
	// 获取证件照信息
	if imgMap, err = s.dao.RealnameApplyIMG(ctx, imgIDs); err != nil {
		return
	}
	for _, ra := range list {
		if mem, ok := memMap[ra.MID]; ok {
			ra.ParseMember(mem)
		}
		for _, id := range ra.IMGIDs {
			if img, ok := imgMap[id]; ok {
				ra.ParseDBApplyIMG(img.IMGData)
			}
		}
		ra.Times = midMap[ra.MID]
	}
	return
}

func (s *Service) realnameAlipayList(ctx context.Context, arg *model.ArgRealnameList) (list []*model.RespRealnameApply, total int, err error) {
	var (
		dl []*model.DBRealnameAlipayApply
	)
	if arg.Card != "" {
		// 如果查询证件号, 则通过证件号MD5 realname_info 查询到 对应的 mids
		var (
			infos []*model.DBRealnameInfo
			mids  []int64
			md5   = cardMD5(arg.Card, arg.DBCardType(), arg.Country)
		)
		log.Info("realnameAlipayList card: %s, md5: %s", arg.Card, md5)
		if infos, err = s.dao.RealnameInfoByCardMD5(ctx, md5, arg.State.DBStatus(), model.ChannelAlipay.DBChannel()); err != nil {
			return
		}
		log.Info("realnameAlipayList infos : %+v", infos)
		if len(infos) <= 0 {
			return
		}
		for _, i := range infos {
			log.Info("realnameAlipayList info: %+v", i)
			mids = append(mids, i.MID)
		}
		log.Info("realnameAlipayList mids: %+v", mids)
		if dl, total, err = s.dao.RealnameAlipayList(ctx, mids, 0, 0, arg.State.DBStatus(), arg.PN, arg.PS, arg.IsDesc); err != nil {
			return
		}
	} else {
		var (
			mids []int64
		)
		if arg.MID > 0 {
			mids = append(mids, arg.MID)
		}
		if dl, total, err = s.dao.RealnameAlipayList(ctx, mids, arg.TSFrom, arg.TSTo, arg.State.DBStatus(), arg.PN, arg.PS, arg.IsDesc); err != nil {
			return
		}
	}
	log.Info("realnameAlipayList dl: %+v, total: %d", dl, total)
	var (
		midMap = make(map[int64]int)
	)
	// append to list
	for _, d := range dl {
		midMap[d.MID] = 0
		var (
			r = &model.RespRealnameApply{}
		)
		r.ParseDBAlipayApply(d)
		list = append(list, r)
	}
	if len(midMap) <= 0 {
		return
	}
	var mids []int64
	for mid := range midMap {
		if midMap[mid], err = s.dao.RealnameApplyCount(ctx, mid); err != nil {
			return
		}
		mids = append(mids, mid)
	}
	var (
		memsArg = &memmdl.ArgMemberMids{
			Mids: mids,
		}
		memMap map[int64]*memmdl.Member
	)
	if memMap, err = s.memberRPC.Members(ctx, memsArg); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, ra := range list {
		if mem, ok := memMap[ra.MID]; ok {
			ra.ParseMember(mem)
		}
		ra.Times = midMap[ra.MID]
	}
	return
}

func cardMD5(card string, cardType int, country int) (res string) {
	if card == "" || cardType < 0 || country < 0 {
		return
	}
	var (
		lowerCode = strings.ToLower(card)
		key       = fmt.Sprintf("%s_%s_%d_%d", model.RealnameSalt, lowerCode, cardType, country)
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

// RealnamePendingList .
func (s *Service) RealnamePendingList(ctx context.Context, arg *model.ArgRealnamePendingList) (list []*model.RespRealnameApply, total int, err error) {
	var (
		larg = &model.ArgRealnameList{
			Channel: arg.Channel,
			State:   model.RealnameApplyStatePending,
			TSFrom:  time.Now().Add(-time.Hour * 24 * 7).Unix(),
			PS:      arg.PS,
			PN:      arg.PN,
		}
	)
	return s.RealnameList(ctx, larg)
}

// RealnameAuditApply .
func (s *Service) RealnameAuditApply(ctx context.Context, arg *model.ArgRealnameAuditApply, adminName string, adminID int64) (err error) {
	var (
		mid int64
	)
	// 1. check the apply state
	switch arg.Channel {
	case model.ChannelMain:
		var apply *model.DBRealnameApply
		if apply, err = s.dao.RealnameMainApply(ctx, arg.ID); err != nil {
			return
		}
		if apply.IsPassed() {
			return
		}
		mid = apply.MID
	case model.ChannelAlipay:
		var apply *model.DBRealnameAlipayApply
		if apply, err = s.dao.RealnameAlipayApply(ctx, arg.ID); err != nil {
			return
		}
		if apply.Status == model.RealnameApplyStateNone.DBStatus() || apply.Status == model.RealnameApplyStateRejective.DBStatus() {
			return
		}
		mid = apply.MID
	}

	var (
		state      = 0
		msgTitle   = ""
		msgContent = ""
		mc         = "2_2_1"
		expNotify  = false
	)
	switch arg.Action {
	case model.RealnameActionPass:
		state = model.RealnameApplyStatePassed.DBStatus()
		msgTitle = "您提交的实名认证已审核通过"
		msgContent = "恭喜，您提交的实名认证已通过审核"
		expNotify = true
	case model.RealnameActionReject:
		state = model.RealnameApplyStateRejective.DBStatus()
		msgTitle = "您提交的实名认证未通过审核"
		msgContent = fmt.Sprintf(`抱歉，您提交的实名认证未通过审核，驳回原因：%s。请修改后重新提交实名认证。`, arg.Reason)
	default:
		err = ecode.RequestErr
		return
	}

	// 2. do something
	switch arg.Channel {
	case model.ChannelMain:
		if err = s.dao.UpdateOldRealnameApply(ctx, arg.ID, state, adminName, adminID, time.Now(), arg.Reason); err != nil {
			return
		}
	case model.ChannelAlipay:
		if err = s.dao.UpdateRealnameAlipayApply(ctx, arg.ID, adminID, adminName, state, arg.Reason); err != nil {
			return
		}
		if err = s.dao.UpdateRealnameInfo(ctx, mid, state, arg.Reason); err != nil {
			return
		}
	}
	go func() {
		if err := s.dao.RawMessage(context.Background(), mc, msgTitle, msgContent, []int64{mid}); err != nil {
			log.Error("%+v", err)
		}
		if expNotify {
			expMsg := &model.AddExpMsg{
				Event: "identify",
				Mid:   mid,
				IP:    metadata.String(ctx, metadata.RemoteIP),
				Ts:    time.Now().Unix(),
			}
			if err := s.dao.PubExpMsg(ctx, expMsg); err != nil {
				log.Error("%+v", err)
			}
		}
	}()
	return
}

// RealnameReasonList .
func (s *Service) RealnameReasonList(ctx context.Context, arg *model.ArgRealnameReasonList) (list []string, total int, err error) {
	return s.dao.RealnameReasonList(ctx)
}

// RealnameSetReason .
func (s *Service) RealnameSetReason(ctx context.Context, arg *model.ArgRealnameSetReason) (err error) {
	return s.dao.UpdateRealnameReason(ctx, arg.Reasons)
}

// RealnameSearchCard .
func (s *Service) RealnameSearchCard(ctx context.Context, cards []string, cardType int, country int) (data map[string]int64, err error) {
	var (
		hashmap = make(map[string]string) //map[hash]card
		hashes  = make([]string, 0)
		list    []*model.DBRealnameInfo
	)
	for _, card := range cards {
		hash := cardMD5(card, cardType, country)
		hashmap[hash] = card
		hashes = append(hashes, hash)
	}
	if list, err = s.dao.RealnameSearchCards(ctx, hashes); err != nil {
		return
	}
	data = make(map[string]int64)
	for _, l := range list {
		if rawCode, ok := hashmap[l.CardMD5]; ok {
			data[rawCode] = l.MID
		}
	}
	return
}

// RealnameUnbind is.
func (s *Service) RealnameUnbind(ctx context.Context, mid int64, adminName string, adminID int64) (err error) {
	var (
		info *model.DBRealnameInfo
	)
	if info, err = s.dao.RealnameInfo(ctx, mid); err != nil {
		return
	}
	if info == nil {
		err = ecode.RealnameAlipayApplyInvalid
		return
	}
	if info.Status != model.RealnameApplyStatePassed.DBStatus() {
		return
	}
	if err = s.dao.UpdateRealnameInfo(ctx, mid, model.RealnameApplyStateRejective.DBStatus(), "管理后台解绑"); err != nil {
		return
	}
	switch info.Channel {
	case model.ChannelMain.DBChannel():
		if err = s.dao.RejectRealnameMainApply(ctx, mid, adminName, adminID, "管理后台解绑"); err != nil {
			return
		}
	case model.ChannelAlipay.DBChannel():
		if err = s.dao.RejectRealnameAlipayApply(ctx, mid, adminName, adminID, "管理后台解绑"); err != nil {
			return
		}
	default:
		log.Warn("Failed to reject realname apply: unrecognized channel: %+v", info)
	}

	go func() {
		r := &report.ManagerInfo{
			Uname:    adminName,
			UID:      adminID,
			Business: model.RealnameManagerLogID,
			Type:     0,
			Oid:      mid,
			Action:   model.LogActionRealnameUnbind,
			Ctime:    time.Now(),
		}
		if err = report.Manager(r); err != nil {
			log.Error("Send manager log failed : %+v , report : %+v", err, r)
			err = nil
			return
		}
		log.Info("Send manager log success report : %+v", r)
	}()
	return
}

// RealnameImage return img
func (s *Service) RealnameImage(ctx context.Context, token string) ([]byte, error) {
	filePath := fmt.Sprintf("%s/%s.txt", conf.Conf.Realname.DataDir, token)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Info("file : %s , not found", filePath)
		return nil, ecode.RequestErr
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()
	img, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.mainCryptor.IMGDecrypt(img)
}

// FetchRealnameImage is
func (s *Service) FetchRealnameImage(ctx context.Context, token string) ([]byte, error) {
	img, err := s.dao.GetRealnameImageCache(ctx, asIMGData(token))
	if err == nil && len(img) > 0 {
		return img, nil
	}
	if err != nil {
		log.Warn("Failed to get realname image from cache: %s: %+v", token, err)
	}

	img, err = s.RealnameImage(ctx, token)
	if err != nil {
		return nil, err
	}

	if len(img) <= _512KiloBytes {
		return img, nil
	}

	striped, err := StripImage(img)
	if err != nil {
		log.Warn("Failed to strip image: %+v", err)
		return img, nil
	}
	return striped, nil
}

// RealnameImagePreview return preview img
func (s *Service) RealnameImagePreview(ctx context.Context, token string, borderSize uint) (data []byte, err error) {
	var (
		src []byte
	)
	if src, err = s.RealnameImage(ctx, token); err != nil {
		return
	}
	if len(src) == 0 {
		return
	}
	var (
		img                 image.Image
		imgWidth, imgHeight int
		imgFormat           string
		sr                  = bytes.NewReader(src)
	)
	if img, imgFormat, err = image.Decode(sr); err != nil {
		log.Warn("Failed to decode image: %+v, return origin image data directly", err)
		return src, nil
	}
	imgWidth, imgHeight = img.Bounds().Dx(), img.Bounds().Dy()
	log.Info("Decode img : %s , format : %s , width : %d , height : %d ", token, imgFormat, imgWidth, imgHeight)
	if imgFormat != "png" && imgFormat != "jpg" && imgFormat != "jpeg" {
		return
	}
	if imgWidth > imgHeight {
		img = resize.Resize(borderSize, 0, img, resize.Lanczos3)
	} else {
		img = resize.Resize(0, borderSize, img, resize.Lanczos3)
	}
	var (
		bb bytes.Buffer
		bw = bufio.NewWriter(&bb)
	)
	switch imgFormat {
	case "jpg", "jpeg":
		if err = jpeg.Encode(bw, img, nil); err != nil {
			err = errors.WithStack(err)
			return
		}
	case "png":
		if err = png.Encode(bw, img); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	data = bb.Bytes()
	return
}

// RealnameExcel export user realname info
func (s *Service) RealnameExcel(ctx context.Context, mids []int64) ([]*model.RealnameExport, error) {
	infos, err := s.dao.BatchRealnameInfo(ctx, mids)
	if err != nil {
		log.Warn("Failed to get realname info with mids: %+v: %+v", mids, err)
		// keep an empty infos
		infos = make(map[int64]*model.DBRealnameInfo)
	}

	pinfos, err := s.dao.PassportQueryByMidsChunked(ctx, mids, 100)
	if err != nil {
		log.Warn("Failed to get passport query by mids: %+v: %+v", mids, err)
		// keep an empty infos
		pinfos = make(map[int64]*model.PassportQueryByMidResult)
	}

	res := make([]*model.RealnameExport, 0, len(mids))
	for _, mid := range mids {
		export := &model.RealnameExport{
			Mid: mid,
		}

		// passport
		func() {
			p, ok := pinfos[mid]
			if !ok {
				log.Warn("Failed to get passport info with mid: %d", mid)
				return
			}
			export.UserID = p.Userid
			export.Uname = p.Name
			export.Tel = p.Tel
		}()

		// realname
		func() {
			info, ok := infos[mid]
			if !ok {
				log.Warn("Failed to get realname info with mid: %d", mid)
				return
			}
			export.Realname = info.Realname
			export.CardType = info.CardType

			cardDecode, err := model.CardDecrypt(info.Card)
			if err != nil {
				log.Error("Failed to decrypt card: %s: %+v", info.Card, err)
				return
			}
			export.CardNum = cardDecode
		}()

		res = append(res, export)
	}
	return res, nil
}

// RealnameSubmit is
func (s *Service) RealnameSubmit(ctx context.Context, arg *model.ArgRealnameSubmit) error {
	encryptedCardNum, err := s.realnameCrypto.CardEncrypt([]byte(arg.CardNum))
	if err != nil {
		return err
	}

	_ = func() error {
		front := &model.DBRealnameApplyIMG{IMGData: asIMGData(arg.FrontImageToken)}
		if err := s.dao.AddRealnameIMG(ctx, front); err != nil {
			return err
		}
		back := &model.DBRealnameApplyIMG{IMGData: asIMGData(arg.BackImageToken)}
		if err := s.dao.AddRealnameIMG(ctx, back); err != nil {
			return err
		}
		apply := &model.DBRealnameApply{
			MID:          arg.Mid,
			Realname:     arg.Realname,
			Country:      arg.Country,
			CardType:     arg.CardType,
			CardNum:      string(encryptedCardNum),
			CardMD5:      cardMD5(arg.CardNum, int(arg.CardType), int(arg.Country)),
			FrontIMG:     front.ID,
			BackIMG:      back.ID,
			Status:       model.RealnameApplyStatePassed.DBStatus(),
			Operator:     arg.Operator,
			OperatorID:   arg.OperatorID,
			OperatorTime: time.Now(),
		}
		if arg.HandImageToken != "" {
			hand := &model.DBRealnameApplyIMG{IMGData: asIMGData(arg.HandImageToken)}
			if err := s.dao.AddRealnameIMG(ctx, hand); err != nil {
				return err
			}
			apply.HandIMG = hand.ID
		}
		if err := s.dao.AddRealnameApply(ctx, apply); err != nil {
			return err
		}
		info := &model.DBRealnameInfo{
			MID:      apply.MID,
			Channel:  model.ChannelMain.DBChannel(),
			Realname: apply.Realname,
			Country:  apply.Country,
			CardType: apply.CardType,
			Card:     apply.CardNum,
			CardMD5:  apply.CardMD5,
			Status:   model.RealnameApplyStatePassed.DBStatus(),
			Reason:   fmt.Sprintf("管理后台提交：%s", arg.Remark),
		}
		return s.dao.SubmitRealnameInfo(ctx, info)
	}

	toOld := func() error {
		front := &model.DeDeIdentificationCardApplyImg{IMGData: asIMGData(arg.FrontImageToken)}
		if err := s.dao.AddOldRealnameIMG(ctx, front); err != nil {
			return err
		}
		back := &model.DeDeIdentificationCardApplyImg{IMGData: asIMGData(arg.BackImageToken)}
		if err := s.dao.AddOldRealnameIMG(ctx, back); err != nil {
			return err
		}
		apply := &model.DeDeIdentificationCardApply{
			MID:           arg.Mid,
			Realname:      arg.Realname,
			Type:          arg.CardType,
			CardData:      string(encryptedCardNum),
			CardForSearch: cardMD5(arg.CardNum, int(arg.CardType), int(arg.Country)),
			FrontImg:      front.ID,
			BackImg:       back.ID,
			ApplyTime:     int32(time.Now().Unix()),
			Operator:      arg.Operator,
			OperatorTime:  int32(time.Now().Unix()),
			Status:        int8(model.RealnameApplyStatePassed.DBStatus()),
			Remark:        fmt.Sprintf("管理后台提交：%s", arg.Remark),
		}
		if arg.HandImageToken != "" {
			hand := &model.DeDeIdentificationCardApplyImg{IMGData: asIMGData(arg.HandImageToken)}
			if err := s.dao.AddOldRealnameIMG(ctx, hand); err != nil {
				return err
			}
			apply.FrontImg2 = hand.ID
		}
		if err := s.dao.AddOldRealnameApply(ctx, apply); err != nil {
			return err
		}
		info := &model.DBRealnameInfo{
			MID:      apply.MID,
			Channel:  model.ChannelMain.DBChannel(),
			Realname: apply.Realname,
			Country:  arg.Country,
			CardType: arg.CardType,
			Card:     apply.CardData,
			CardMD5:  apply.CardForSearch,
			Status:   model.RealnameApplyStatePassed.DBStatus(),
			Reason:   fmt.Sprintf("管理后台提交：%s", arg.Remark),
		}
		return s.dao.SubmitRealnameInfo(ctx, info)
	}

	if err := toOld(); err != nil {
		return err
	}

	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.RealnameManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   model.LogActionRealnameSubmit,
		Ctime:    time.Now(),
	})

	return nil
}

// RealnameFileUpload is
func (s *Service) RealnameFileUpload(ctx context.Context, mid int64, data []byte) (src string, err error) {
	var (
		md5Engine = md5.New()
		hashMID   string
		hashRand  string
		fileName  string
		dirPath   string
		dateStr   string
	)
	md5Engine.Write([]byte(strconv.FormatInt(mid, 10)))
	hashMID = hex.EncodeToString(md5Engine.Sum(nil))
	md5Engine.Reset()
	md5Engine.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	md5Engine.Write([]byte(strconv.FormatInt(mrand.Int63n(1000000), 10)))
	hashRand = hex.EncodeToString(md5Engine.Sum(nil))
	fileName = fmt.Sprintf("%s_%s.txt", hashMID[:6], hashRand)
	dateStr = time.Now().Format("20060102")
	dirPath = fmt.Sprintf("%s/%s/", s.c.Realname.DataDir, dateStr)

	var (
		dataFile      *os.File
		writeFileSize int
		encrptedData  []byte
	)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if encrptedData, err = s.mainCryptor.IMGEncrypt(data); err != nil {
		err = errors.WithStack(err)
		return
	}
	if dataFile, err = os.Create(dirPath + fileName); err != nil {
		err = errors.Wrapf(err, "create file %s failed", dirPath+fileName)
		return
	}
	defer dataFile.Close()
	if writeFileSize, err = dataFile.Write(encrptedData); err != nil {
		err = errors.Wrapf(err, "write file %s size %d failed", dirPath+fileName, len(encrptedData))
		return
	}
	if writeFileSize != len(encrptedData) {
		err = errors.Errorf("Write file data to %s , expected %d actual %d", dirPath+fileName, len(encrptedData), writeFileSize)
		return
	}
	src = fmt.Sprintf("%s/%s", dateStr, strings.TrimSuffix(fileName, ".txt"))
	return
}

func asIMGData(imgToken string) string {
	return model.RealnameImgPrefix + imgToken + model.RealnameImgSuffix
}

func asIMGToken(IMGData string) string {
	token := strings.TrimPrefix(IMGData, "/idenfiles/")
	token = strings.TrimSuffix(token, ".txt")
	return token
}

// StripImage is
func StripImage(raw []byte) ([]byte, error) {
	i, format, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out := &bytes.Buffer{}
	switch format {
	case "jpg", "jpeg":
		if err := jpeg.Encode(out, i, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return nil, errors.WithStack(err)
		}
	default:
		return nil, errors.Errorf("Unsupported type: %s", format)
	}
	return out.Bytes(), nil
}
