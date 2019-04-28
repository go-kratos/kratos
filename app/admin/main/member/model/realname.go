package model

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/member/conf"
	memmdl "go-common/app/service/main/member/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// realname conf var
var (
	RealnameSalt      = "biliidentification@#$%^&*()(*&^%$#"
	RealnameImgPrefix = "/idenfiles/"
	RealnameImgSuffix = ".txt"

	LogActionRealnameUnbind = "realname_unbind"
	LogActionRealnameBack   = "realname_back"
	LogActionRealnameSubmit = "realname_submit"

	RealnameManagerLogID = 251
)

type realnameChannel string

// RealnameChannle
const (
	ChannelMain   realnameChannel = "main"
	ChannelAlipay realnameChannel = "alipay"
)

func (ch realnameChannel) DBChannel() uint8 {
	switch ch {
	case ChannelMain:
		return 0
	case ChannelAlipay:
		return 1
	}
	return 0
}

type realnameAction string

// RealnameAction
const (
	RealnameActionPass   realnameAction = "pass"
	RealnameActionReject realnameAction = "reject"
)

// RealnameCardType is.
type RealnameCardType string

const (
	cardTypeIdentityCard    RealnameCardType = "identity_card"
	cardTypeForeignPassport RealnameCardType = "foreign_passport"
	//Mainland Travel Permit for Hong Kong and Macao Residents
	cardTypeHongkongMacaoPermit RealnameCardType = "hongkong_macao_travel_permit"
	//Mainland travel permit for Taiwan residents
	cardTypeTaiwanPermit  RealnameCardType = "taiwan_travel_permit"
	cardTypeChinaPassport RealnameCardType = "china_passport"
	//Foreigner's Permanent Residence Card
	cardTypeForeignerPermanentResidenceCard RealnameCardType = "foreigner_permanent_residence_card"
	cardTypeForeignIdentityCard             RealnameCardType = "foreign_identity_card"
)

// RealnameApplyStatus is.
type RealnameApplyStatus string

// RealnameApplyStatus
const (
	RealnameApplyStateAll       RealnameApplyStatus = "all"
	RealnameApplyStatePending   RealnameApplyStatus = "pending"
	RealnameApplyStatePassed    RealnameApplyStatus = "passed"
	RealnameApplyStateRejective RealnameApplyStatus = "rejective"
	RealnameApplyStateNone      RealnameApplyStatus = "none"
)

// DBStatus is.
func (r RealnameApplyStatus) DBStatus() int {
	switch r {
	case RealnameApplyStatePending:
		return 0
	case RealnameApplyStatePassed:
		return 1
	case RealnameApplyStateRejective:
		return 2
	default:
		return -1
	}
}

// ArgRealnameList is.
type ArgRealnameList struct {
	Channel  realnameChannel     `form:"channel" validate:"required"`
	MID      int64               `form:"mid"`
	Card     string              `form:"card"`
	CardType RealnameCardType    `form:"card_type"`
	Country  int                 `form:"country"`
	OPName   string              `form:"op_name"`
	TSFrom   int64               `form:"ts_from"`
	TSTo     int64               `form:"ts_to"`
	State    RealnameApplyStatus `form:"state"`
	PS       int                 `form:"ps"`
	PN       int                 `form:"pn"`
	IsDesc   bool                `form:"is_desc"`
}

// DBCardType return card_type store in db
func (a *ArgRealnameList) DBCardType() int {
	switch a.CardType {
	case cardTypeIdentityCard:
		return 0
	case cardTypeForeignPassport:
		return 1
	case cardTypeHongkongMacaoPermit:
		return 2
	case cardTypeTaiwanPermit:
		return 3
	case cardTypeChinaPassport:
		return 4
	case cardTypeForeignerPermanentResidenceCard:
		return 5
	case cardTypeForeignIdentityCard:
		return 6
	default:
		log.Warn("ArgRealnameList : %+v , unknown CardType", a)
		return -1
	}
}

// DBCountry return country store in db
func (a *ArgRealnameList) DBCountry() int {
	if a.CardType == "" {
		return -1
	}
	return a.Country
}

// DBState return state store in db
func (a *ArgRealnameList) DBState() int {
	switch a.State {
	case RealnameApplyStateAll:
		return -1
	case RealnameApplyStatePending:
		return 0
	case RealnameApplyStatePassed:
		return 1
	case RealnameApplyStateRejective:
		return 2
	case RealnameApplyStateNone:
		return 3
	default:
		log.Warn("ArgRealnameList : %+v , unknown State", a)
		return 0
	}
}

// ArgRealnamePendingList is.
type ArgRealnamePendingList struct {
	Channel realnameChannel `form:"channel" validate:"required"`
	PS      int             `form:"ps"`
	PN      int             `form:"pn"`
}

// ArgRealnameAuditApply is.
type ArgRealnameAuditApply struct {
	ID      int64           `form:"id" validate:"required"`
	Channel realnameChannel `form:"channel" validate:"required"`
	Action  realnameAction  `form:"action" validate:"required"`
	Reason  string          `form:"reason"`
}

// DBChannel return channel store in db
func (a *ArgRealnameAuditApply) DBChannel() int {
	switch a.Channel {
	case ChannelMain:
		return 0
	case ChannelAlipay:
		return 1
	default:
		log.Warn("ArgRealnameAuditApply : %+v , unknown Channel", a)
		return 0
	}
}

// ArgRealnameReasonList is.
type ArgRealnameReasonList struct {
	PS int `form:"ps"`
	PN int `form:"pn"`
}

// ArgRealnameSetReason is.
type ArgRealnameSetReason struct {
	Reasons []string `form:"reasons,split"`
}

// ArgRealnameImage is.
type ArgRealnameImage struct {
	Token string `form:"token" validate:"required"`
}

// ArgRealnameImagePreview is.
type ArgRealnameImagePreview struct {
	ArgRealnameImage
	BorderSize uint `form:"border_size"` // 图片最大边长度（缩放后）
}

// ArgRealnameSearchCard is.
type ArgRealnameSearchCard struct {
	Cards    []string `form:"cards,split" validate:"required"`
	CardType int      `form:"card_type"`
	Country  int      `form:"card_type"`
}

// RespRealnameApply is.
type RespRealnameApply struct {
	ID       int64               `json:"id"`
	Channel  realnameChannel     `json:"channel"`
	MID      int64               `json:"mid"`
	Nickname string              `json:"nickname"`
	Times    int                 `json:"times"`
	CardType RealnameCardType    `json:"card_type"`
	Country  int16               `json:"country"`
	Card     string              `json:"card"`
	Realname string              `json:"realname"`
	Level    int32               `json:"level"`
	IMGIDs   []int64             `json:"-"`
	IMGs     []string            `json:"imgs"`
	State    RealnameApplyStatus `json:"state"`
	OPName   string              `json:"op_name"`
	OPTS     int64               `json:"op_ts"`
	OPReason string              `json:"op_reason"`
	CreateTS int64               `json:"create_ts"`
}

// ParseDBMainApply parse realname_apply from db
func (r *RespRealnameApply) ParseDBMainApply(db *DBRealnameApply) {
	var err error
	r.ID = db.ID
	r.Channel = ChannelMain
	r.MID = db.MID
	r.CardType = r.convertCardType(db.CardType)
	r.Country = db.Country
	if db.CardNum != "" {
		if r.Card, err = CardDecrypt(db.CardNum); err != nil {
			log.Error("%+v", err)
		}
	}
	r.IMGIDs = append(r.IMGIDs, db.HandIMG, db.FrontIMG, db.BackIMG)
	r.Realname = db.Realname
	r.State = r.ParseStatus(db.Status)
	r.OPName = db.Operator
	r.OPTS = db.OperatorTime.Unix()
	r.OPReason = db.Remark
	r.CreateTS = db.CTime.Unix()
}

// ParseDBAlipayApply parse realname_alipay_apply from db
func (r *RespRealnameApply) ParseDBAlipayApply(db *DBRealnameAlipayApply) {
	var err error
	r.ID = db.ID
	r.Channel = ChannelAlipay
	r.MID = db.MID
	r.CardType = cardTypeIdentityCard // identity_card
	r.Country = 0                     // china
	if db.Card != "" {
		if r.Card, err = CardDecrypt(db.Card); err != nil {
			log.Error("%+v", err)
		}
	}
	r.ParseDBApplyIMG(db.IMG)
	r.Realname = db.Realname
	r.State = r.ParseStatus(db.Status)
	r.OPName = "alipay"
	if db.Operator != "" {
		r.OPName = db.Operator
	}
	r.OPTS = db.OperatorTime.Unix()
	r.OPReason = db.Reason
	r.CreateTS = db.CTime.Unix()
}

// ParseDBApplyIMG parse apply_img from db
func (r *RespRealnameApply) ParseDBApplyIMG(token string) {
	r.IMGs = append(r.IMGs, imgURL(token))
}

// ParseMember parse member info from rpc call
func (r *RespRealnameApply) ParseMember(mem *memmdl.Member) {
	r.Nickname = mem.Name
	r.Level = mem.LevelInfo.Cur
}

func imgURL(token string) string {
	token = strings.TrimPrefix(token, "/idenfiles/")
	token = strings.TrimSuffix(token, ".txt")
	return fmt.Sprintf(conf.Conf.Realname.ImageURLTemplate, token)
}

// ParseStatus parse status stored in db
func (r *RespRealnameApply) ParseStatus(status int) (s RealnameApplyStatus) {
	switch status {
	case 0:
		return RealnameApplyStatePending
	case 1:
		return RealnameApplyStatePassed
	case 2:
		return RealnameApplyStateRejective
	default:
		log.Warn("RespRealnameApply parse status err , unknown apply status :%d", status)
		return RealnameApplyStateNone
	}
}

// ConvertCardType convert card_type from db to api
func (r *RespRealnameApply) convertCardType(cardType int8) (t RealnameCardType) {
	switch cardType {
	case 0:
		return cardTypeIdentityCard
	case 1:
		return cardTypeForeignPassport
	case 2:
		return cardTypeHongkongMacaoPermit
	case 3:
		return cardTypeTaiwanPermit
	case 4:
		return cardTypeChinaPassport
	case 5:
		return cardTypeForeignerPermanentResidenceCard
	case 6:
		return cardTypeForeignIdentityCard
	default:
		log.Warn("RespRealnameApply parse card type err , unknown card type :%d", cardType)
		return cardTypeIdentityCard
	}
}

// DBRealnameInfo is.
type DBRealnameInfo struct {
	ID       int64     `gorm:"column:id"`
	MID      int64     `gorm:"column:mid"`
	Channel  uint8     `gorm:"column:channel"`
	Realname string    `gorm:"column:realname"`
	Country  int16     `gorm:"column:country"`
	CardType int8      `gorm:"column:card_type"`
	Card     string    `gorm:"column:card"`
	CardMD5  string    `gorm:"column:card_md5"`
	Status   int       `gorm:"column:status"`
	Reason   string    `gorm:"column:reason"`
	CTime    time.Time `gorm:"column:ctime"`
	MTime    time.Time `gorm:"column:mtime"`
}

// TableName is...
func (d *DBRealnameInfo) TableName() string {
	return "realname_info"
}

// DBRealnameApply is.
type DBRealnameApply struct {
	ID           int64     `gorm:"column:id"`
	MID          int64     `gorm:"column:mid"`
	Realname     string    `gorm:"column:realname"`
	Country      int16     `gorm:"column:country"`
	CardType     int8      `gorm:"column:card_type"`
	CardNum      string    `gorm:"column:card_num"`
	CardMD5      string    `gorm:"column:card_md5"`
	HandIMG      int64     `gorm:"column:hand_img"`
	FrontIMG     int64     `gorm:"column:front_img"`
	BackIMG      int64     `gorm:"column:back_img"`
	Status       int       `gorm:"column:status"`
	Operator     string    `gorm:"column:operator"`
	OperatorID   int64     `gorm:"column:operator_id"`
	OperatorTime time.Time `gorm:"column:operator_time"`
	Remark       string    `gorm:"column:remark"`
	RemarkStatus int8      `gorm:"column:remark_status"`
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
}

// TableName is...
func (d *DBRealnameApply) TableName() string {
	return "realname_apply"
}

// DBRealnameAlipayApply is.
type DBRealnameAlipayApply struct {
	ID           int64     `gorm:"column:id"`
	MID          int64     `gorm:"column:mid"`
	Realname     string    `gorm:"column:realname"`
	Card         string    `gorm:"column:card"`
	IMG          string    `gorm:"column:img"`
	Status       int       `gorm:"column:status"`
	Reason       string    `gorm:"column:reason"`
	Bizno        string    `gorm:"column:bizno"`
	Operator     string    `gorm:"column:operator"`
	OperatorID   int64     `gorm:"column:operator_id"`
	OperatorTime time.Time `gorm:"column:operator_time"`
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
}

// TableName is...
func (d *DBRealnameAlipayApply) TableName() string {
	return "realname_alipay_apply"
}

// IsPassed is...
func (d *DBRealnameApply) IsPassed() bool {
	return d.Status == 1
}

// DBRealnameApplyIMG is.
type DBRealnameApplyIMG struct {
	ID      int64     `gorm:"column:id"`
	IMGData string    `gorm:"column:img_data"`
	CTime   time.Time `gorm:"column:ctime"`
	MTime   time.Time `gorm:"column:mtime"`
}

// TableName ...
func (d *DBRealnameApplyIMG) TableName() string {
	return "realname_apply_img"
}

// DBRealnameConfig ...
type DBRealnameConfig struct {
	ID    int64     `gorm:"column:id"`
	Key   string    `gorm:"column:key"`
	Data  string    `gorm:"column:data"`
	CTime time.Time `gorm:"column:ctime"`
	MTime time.Time `gorm:"column:mtime"`
}

// TableName ...
func (d *DBRealnameConfig) TableName() string {
	return "realname_config"
}

// CardDecrypt is
func CardDecrypt(data string) (text string, err error) {
	var (
		dataBytes     = []byte(data)
		decryptedData []byte
		textBytes     []byte
		size          int
	)
	decryptedData = make([]byte, base64.StdEncoding.DecodedLen(len(dataBytes)))
	if size, err = base64.StdEncoding.Decode(decryptedData, dataBytes); err != nil {
		err = errors.Wrapf(err, "base decode %s", data)
		return
	}
	if textBytes, err = rsaDecrypt(decryptedData[:size]); err != nil {
		err = errors.Wrapf(err, "rsa decrypt %s , data : %s", decryptedData, data)
		return
	}
	text = string(textBytes)
	return
}

func rsaDecrypt(text []byte) (content []byte, err error) {
	block, _ := pem.Decode(conf.Conf.Realname.RsaPriv)
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

// DBRealnameAuditLog is.
type DBRealnameAuditLog struct {
	ID        int64     `gorm:"column:id"`
	MID       int64     `gorm:"column:mid"`
	AdminID   int64     `gorm:"column:admin_id"`
	AdminName string    `gorm:"column:admin_name"`
	Channel   uint8     `gorm:"column:channel"`
	FromState int       `gorm:"column:from_state"`
	ToState   int       `gorm:"column:to_state"`
	CTime     time.Time `gorm:"column:ctime"`
	MTime     time.Time `gorm:"column:mtime"`
}

// Tablename is.
func (d *DBRealnameAuditLog) Tablename() string {
	return "realname_audit_log"
}

// Realname is.
type Realname struct {
	State    RealnameApplyStatus `json:"state"`
	Channel  realnameChannel     `json:"moral"`
	Card     string              `json:"card"`
	CardType int8                `json:"card_type"`
	Country  int16               `json:"country"`
	Realname string              `json:"realname"`
	Images   []string            `json:"images"`
}

// ParseDBApplyIMG parse apply_img from db
func (r *Realname) ParseDBApplyIMG(token string) {
	r.Images = append(r.Images, imgURL(token))
}

// ParseInfo .
func (r *Realname) ParseInfo(info *DBRealnameInfo) {
	switch info.Status {
	case 0:
		r.State = RealnameApplyStatePending
	case 1:
		r.State = RealnameApplyStatePassed
	case 2:
		r.State = RealnameApplyStateRejective
	default:
		log.Warn("Realname status err , unknown info status :%d", info.Status)
		r.State = RealnameApplyStateNone
	}
	switch info.Channel {
	case 0:
		r.Channel = ChannelMain
	case 1:
		r.Channel = ChannelAlipay
	default:
		log.Warn("Realname channel err , unknown info channel :%d", info.Channel)
		r.Channel = ChannelMain
	}
	r.Realname = info.Realname
	r.CardType = info.CardType
	r.Country = info.Country
	var err error
	var maskedCard string
	if info.Card != "" && r.State == RealnameApplyStatePassed {
		if r.Card, err = CardDecrypt(info.Card); err != nil {
			log.Error("%+v", err)
		}
		var (
			cStrs = strings.Split(r.Card, "")
		)
		if len(cStrs) > 0 {
			if len(cStrs) == 1 {
				maskedCard = "*"
			} else if len(cStrs) > 5 {

				maskedCard = cStrs[0] + strings.Repeat("*", len(cStrs)-3) + strings.Join(cStrs[len(cStrs)-2:], "")
			} else {
				maskedCard = cStrs[0] + strings.Repeat("*", len(cStrs)-1)
			}
		}
		r.Card = maskedCard
	}
}

// RealnameExport is.
type RealnameExport struct {
	Mid      int64  `json:"mid" gorm:"column:mid"`
	UserID   string `json:"userid" gorm:"column:userid"`
	Uname    string `json:"uname" gorm:"column:uname"`
	Realname string `json:"realname" gorm:"column:realname"`
	Tel      string `json:"tel" gorm:"column:tel"`
	CardType int8   `json:"card_type" gorm:"column:card_type"`
	CardNum  string `json:"card_num" gorm:"column:card_num"`
}

// PassportQueryByMidResult is.
type PassportQueryByMidResult struct {
	Mid      int64      `json:"mid"`
	Name     string     `json:"name"`
	Userid   string     `json:"userid"`
	Email    string     `json:"email"`
	Tel      string     `json:"tel"`
	Jointime xtime.Time `json:"jointime"`
}

var _cardTypeToString = map[int8]string{
	0: "身份证",
	1: "护照(境外签发)",
	2: "港澳居民来往内地通行证",
	3: "台湾居民来往大陆通行证",
	4: "护照(中国签发)",
	5: "外国人永久居留证",
	6: "其他国家或地区身份证",
}

// CardTypeString is
func CardTypeString(cardType int8) string {
	typeString, ok := _cardTypeToString[cardType]
	if !ok {
		return strconv.FormatInt(int64(cardType), 10)
	}
	return typeString
}
