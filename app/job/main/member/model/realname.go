package model

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/member/conf"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	realnameSalt = "biliidentification@#$%^&*()(*&^%$#"
)

//RealnamePersonMessage is.
type RealnamePersonMessage struct {
	MID          int64  `json:"mid"`
	Realname     string `json:"realname"`
	IdentifyCard string `json:"identify_card"`
}

//RealnameApplyMessage is.
type RealnameApplyMessage struct {
	ID               int    `json:"id"`
	MID              int64  `json:"mid"`
	Realname         string `json:"realname"`
	Type             int    `json:"type"`
	CardDataCanal    string `json:"card_data"`
	CardForSearch    string `json:"card_for_search"`
	FrontIMG         int    `json:"front_img"`
	BackIMG          int    `json:"back_img"`
	FrontIMG2        int    `json:"front_img2"`
	ApplyTimeUnix    int64  `json:"apply_time"`
	Operater         string `json:"operater"`
	OperaterTimeUnix int64  `json:"operater_time"`
	Status           int8   `json:"status"`
	Remark           string `json:"remark"`
	RemarkStatus     int8   `json:"remark_status"`
}

// CardMD5 is.
func (r *RealnameApplyMessage) CardMD5() (res string) {
	return cardReMD5(r.CardData(), r.CardType(), r.Country())
}

// cardReMD5 is.
func cardReMD5(encrypedCard string, cardType int, country int) (res string) {
	card, err := cardDecrypt([]byte(encrypedCard))
	if err != nil {
		log.Error("cardNewMD5 decrypt err : %+v", err)
		return
	}
	return cardMD5(string(card), cardType, country)
}

func cardMD5(card string, cardType int, country int) (res string) {
	var (
		lowerCode = strings.ToLower(card)
		key       = fmt.Sprintf("%s_%s_%d_%d", realnameSalt, lowerCode, cardType, country)
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

func cardDecrypt(data []byte) (text []byte, err error) {
	var (
		decryptedData []byte
		size          int
	)
	decryptedData = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	if size, err = base64.StdEncoding.Decode(decryptedData, data); err != nil {
		err = errors.Wrapf(err, "base decode %s", data)
		return
	}
	if text, err = rsaDecrypt(decryptedData[:size]); err != nil {
		err = errors.Wrapf(err, "rsa decrypt %s , data : %s", decryptedData, data)
		return
	}
	return
}

func rsaDecrypt(text []byte) (content []byte, err error) {
	block, _ := pem.Decode(conf.Conf.RealnameRsaPriv)
	if block == nil {
		err = errors.New("private key error")
		return
	}
	var (
		privateKey *rsa.PrivateKey
	)
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		err = errors.WithStack(err)
		return
	}
	if content, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, text); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//ApplyTime is.
func (r *RealnameApplyMessage) ApplyTime() (t time.Time) {
	t = time.Unix(r.ApplyTimeUnix, 0)
	return
}

//OperaterTime is.
func (r *RealnameApplyMessage) OperaterTime() (t time.Time) {
	t = time.Unix(r.OperaterTimeUnix, 0)
	return
}

//CardType is.
func (r *RealnameApplyMessage) CardType() (t int) {
	if r.Type > 6 {
		return 6
	}
	return r.Type
}

//Country is.
func (r *RealnameApplyMessage) Country() (t int) {
	if r.Type > 6 {
		return r.Type
	}
	return 0
}

//CardData is.
func (r *RealnameApplyMessage) CardData() (data string) {
	if r.CardDataCanal == "" {
		log.Warn("card data empty (+v)", r)
		return ""
	}
	bytes, err := base64.StdEncoding.DecodeString(r.CardDataCanal)
	if err != nil {
		err = errors.Wrapf(err, "decode (%+v) failed", r)
		log.Error("+v", err)
		return ""
	}
	return string(bytes)
}

//RealnameApplyImgMessage is.
type RealnameApplyImgMessage struct {
	ID         int    `json:"id"`
	IMGData    string `json:"img_data"`
	AddTimeStr string `json:"add_time"`
	AddTimeDB  time.Time
}

//AddTime is.
func (r *RealnameApplyImgMessage) AddTime() (t time.Time) {
	if r.AddTimeDB.IsZero() {
		var err error
		if t, err = time.ParseInLocation("2006-01-02 15:04:05", r.AddTimeStr, time.Local); err != nil {
			log.Error("%+v", err)
			t = time.Now()
		}
		return
	}
	return r.AddTimeDB
}

// RealnameApplyStatus is.
type RealnameApplyStatus int8

const (
	// RealnameApplyStatusPending is.
	RealnameApplyStatusPending RealnameApplyStatus = iota
	// RealnameApplyStatusPass is.
	RealnameApplyStatusPass
	// RealnameApplyStatusBack is.
	RealnameApplyStatusBack
	// RealnameApplyStatusNone is.
	RealnameApplyStatusNone
)

// IsPass return is apply passed
func (r RealnameApplyStatus) IsPass() bool {
	switch r {
	case RealnameApplyStatusPass:
		return true
	default:
		return false
	}
}

// RealnameChannel is
type RealnameChannel int8

// RealnameChannel enum
const (
	RealnameChannelMain RealnameChannel = iota
	RealnameChannelAlipay
)

const (
	// RealnameCountryChina is.
	RealnameCountryChina = 0
	// RealnameCardTypeIdentity is.
	RealnameCardTypeIdentity = 0
)

// RealnameInfo is user realname status info
type RealnameInfo struct {
	ID       int64               `json:"id"`
	MID      int64               `json:"mid"`
	Channel  RealnameChannel     `json:"channel"`
	Realname string              `json:"realname"`
	Country  int                 `json:"country"`
	CardType int                 `json:"card_type"`
	Card     string              `json:"card"`
	CardMD5  string              `json:"card_md5"`
	Status   RealnameApplyStatus `json:"status"`
	Reason   string              `json:"reason"`
	CTime    time.Time           `json:"ctime"`
	MTime    time.Time           `json:"mtime"`
}

// DecryptedCard is
func (r *RealnameInfo) DecryptedCard() (string, error) {
	raw, err := cardDecrypt([]byte(r.Card))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// RealnameApply is user realname apply info from default channel.
type RealnameApply struct {
	ID           int       `json:"id"`
	MID          int64     `json:"mid"`
	Realname     string    `json:"realname"`
	Country      int16     `json:"country"`
	CardType     int8      `json:"card_type"`
	CardNum      string    `json:"card_num"`
	CardMD5      string    `json:"card_md5"`
	HandIMG      int       `json:"hand_img"`
	FrontIMG     int       `json:"front_img"`
	BackIMG      int       `json:"back_img"`
	Status       int       `json:"status"`
	Operator     string    `json:"operator"`
	OperatorID   int64     `json:"operator_id"`
	OperatorTime time.Time `json:"operator_time"`
	Remark       string    `json:"remark"`
	RemarkStatus int8      `json:"remark_status"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"mtime"`
}

// RealnameAlipayApply is user alipay apply info from alipay channle.
type RealnameAlipayApply struct {
	ID       int64               `json:"id"`
	MID      int64               `json:"mid"`
	Realname string              `json:"realname"`
	Card     string              `json:"card"`
	IMG      string              `json:"img"`
	Status   RealnameApplyStatus `json:"status"`
	Reason   string              `json:"reason"`
	Bizno    string              `json:"bizno"`
	CTime    time.Time           `json:"ctime"`
	MTime    time.Time           `json:"mtime"`
}
