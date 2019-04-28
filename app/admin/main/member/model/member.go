package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	xtime "go-common/library/time"
)

const (
	_expMulti = 100
	level1    = 1
	level2    = 200
	level3    = 1500
	level4    = 4500
	level5    = 10800
	level6    = 28800
	levelMax  = -1

	_URLNoFace = "http://static.hdslb.com/images/member/noface.gif"

	// ManagerLogID manager log id.
	ManagerLogID = 121
	//FaceCheckLogID is.
	FaceCheckLogID = 161

	// bfs facepri bucket
	_facepriKeyID     = "8923aff2e1124bb2"
	_facepriKeySecret = "b237e8927823cc2984aee980123cb0"
)

// base audit type const.
const (
	BaseAuditType = iota
	BaseAuditTypeFace
	BaseAuditTypeSign
	BaseAuditTypeName
)

// Base is.
type Base struct {
	Mid      int64      `json:"mid" gorm:"column:mid"`
	Name     string     `json:"name" gorm:"column:name"`
	Sex      int64      `json:"sex" gorm:"column:sex"`
	Face     string     `json:"face" gorm:"column:face"`
	Sign     string     `json:"sign" gorm:"column:sign"`
	Rank     int64      `json:"rank" gorm:"column:rank"`
	Birthday xtime.Time `json:"birthday" gorm:"column:birthday"`
	CTime    xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime    xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// RandFaceURL get face URL
func (b *Base) RandFaceURL() {
	if b.Face == "" {
		b.Face = _URLNoFace
		return
	}
	b.Face = fmt.Sprintf("http://i%d.hdslb.com%s", rand.Int63n(3), b.Face)
}

// Detail is.
type Detail struct {
	Mid      int64      `json:"mid" gorm:"column:mid"`
	Birthday xtime.Time `json:"birthday" gorm:"column:birthday"`
	Place    int64      `json:"place" gorm:"column:place"`
	Marital  int64      `json:"marital" gorm:"column:marital"`
	Dating   int64      `json:"dating" gorm:"column:dating"`
	Tags     string     `json:"tags" gorm:"column:tags"`
	CTime    xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime    xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// Exp is.
type Exp struct {
	Mid     int64      `json:"mid" gorm:"column:mid"`
	Exp     int64      `json:"exp" gorm:"column:exp"`
	Flag    uint32     `json:"flag" gorm:"column:flag"`
	AddTime xtime.Time `json:"addtime" gorm:"column:addtime"`
	CTime   xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime   xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// Moral is.
type Moral struct {
	Mid             int64      `json:"mid" gorm:"column:mid"`
	Moral           int64      `json:"moral" gorm:"column:moral"`
	Added           int64      `json:"added" gorm:"column:added"`
	Deducted        int64      `json:"deducted" gorm:"column:deducted"`
	LastRecoverDate xtime.Time `json:"last_recover_date" gorm:"colum:last_recover_date"`
	CTime           xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime           xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// UserAddit is.
type UserAddit struct {
	ID             int64      `json:"id" gorm:"column:id"`
	Mid            int64      `json:"mid" gorm:"column:mid"`
	FaceReject     int64      `json:"face_reject" gorm:"colum:face_reject"`
	ViolationCount int64      `json:"violation_count" gorm:"colum:violation_count"`
	CTime          xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime          xtime.Time `json:"mtime" gorm:"column:mtime"`
	Remark         string     `json:"remark" gorm:"column:remark"`
}

// Level is.
type Level struct {
	CurrentLevel int32 `json:"current_level"`
	CurrentMin   int32 `json:"current_min"`
	CurrentExp   int32 `json:"current_exp"`
	NextExp      int32 `json:"next_exp"`
}

// Profile is.
type Profile struct {
	Base     Base      `json:"base"`
	Detail   Detail    `json:"detail"`
	Exp      Exp       `json:"exp"`
	Level    Level     `json:"level"`
	Moral    Moral     `json:"moral"`
	Official Official  `json:"official"`
	Coin     Coin      `json:"coin"`
	Addit    UserAddit `json:"addit"`
	Realanme Realname  `json:"realname"`
}

// Coin is.
type Coin struct {
	Coins float64 `json:"coins"`
}

// UserLog is.
type UserLog struct {
	Mid     int64             `json:"mid"`
	IP      string            `json:"ip"`
	TS      int64             `json:"ts"`
	Content map[string]string `json:"content"`
}

// FaceRecord is.
type FaceRecord struct {
	ID         int64      `json:"id"`
	Mid        int64      `json:"mid"`
	ModifyTime xtime.Time `json:"modify_time"`
	ApplyTime  xtime.Time `json:"apply_time"`
	NewFace    string     `json:"new_face"`
	OldFace    string     `json:"old_face"`
	Operator   string     `json:"operator"`
	Status     int8       `json:"status"`
}

// BaseReview is.
type BaseReview struct {
	Base
	Addit UserAddit  `json:"addit"`
	Logs  []AuditLog `json:"logs"`
}

// AddExpMsg is.
type AddExpMsg struct {
	Event string `json:"event,omitempty"`
	Mid   int64  `json:"mid,omitempty"`
	IP    string `json:"ip,omitempty"`
	Ts    int64  `json:"ts,omitempty"`
}

// BuildFaceURL is.
func BuildFaceURL(raw string) string {
	if raw == "" {
		return _URLNoFace
	}
	ori, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if ori.Path == "/images/member/noface.gif" {
		return _URLNoFace
	}
	if strings.HasPrefix(ori.Path, "/bfs/facepri") {
		token := authorize(_facepriKeyID, _facepriKeySecret, "GET", "facepri", filepath.Base(ori.Path), time.Now().Unix())
		p := url.Values{}
		p.Set("token", token)
		ori.RawQuery = p.Encode()
	}
	if ori.Hostname() == "" {
		ori.Host = fmt.Sprintf("i%d.hdslb.com", rand.Int63n(3))
		ori.Scheme = "http"
	}
	return ori.String()
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket, filename string, expire int64) string {
	content := fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, filename, expire)
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s:%s:%d", key, signature, expire)
}

// BuildFaceURL is.
func (fr *FaceRecord) BuildFaceURL() {
	fr.NewFace = BuildFaceURL(fr.NewFace)
	fr.OldFace = BuildFaceURL(fr.OldFace)
}

// ParseStatus is.
func ParseStatus(s string) int8 {
	st, _ := strconv.ParseInt(s, 10, 8)
	return int8(st)
}

// ParseLogTime is.
func ParseLogTime(ts string) (xt xtime.Time, err error) {
	var (
		t time.Time
	)
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local); err != nil {
		return
	}
	xt.Scan(t)
	return
}

// ParseApplyTime is.
func ParseApplyTime(ts string) xtime.Time {
	ti, _ := strconv.ParseInt(ts, 10, 64)
	return xtime.Time(ti)
}

// NewProfile is.
func NewProfile() *Profile {
	return &Profile{}
}

// FaceRecordList is
type FaceRecordList []*FaceRecord

// Filter is
func (frl FaceRecordList) Filter(con func(*FaceRecord) bool) FaceRecordList {
	res := make(FaceRecordList, 0)
	for _, fr := range frl {
		if con(fr) {
			res = append(res, fr)
		}
	}
	return res
}

// Paginate is
func (frl FaceRecordList) Paginate(skip int, size int) FaceRecordList {
	if skip > len(frl) {
		skip = len(frl)
	}

	end := skip + size
	if end > len(frl) {
		end = len(frl)
	}

	return frl[skip:end]
}

// FromExp is.
func (lv *Level) FromExp(e *Exp) {
	exp := e.Exp / _expMulti
	switch {
	case exp < level1:
		lv.CurrentLevel = 0
		lv.CurrentMin = 0
		lv.NextExp = level1
	case exp < level2:
		lv.CurrentLevel = 1
		lv.CurrentMin = level1
		lv.NextExp = level2
	case exp < level3:
		lv.CurrentLevel = 2
		lv.CurrentMin = level2
		lv.NextExp = level3
	case exp < level4:
		lv.CurrentLevel = 3
		lv.CurrentMin = level3
		lv.NextExp = level4
	case exp < level5:
		lv.CurrentLevel = 4
		lv.CurrentMin = level4
		lv.NextExp = level5
	case exp < level6:
		lv.CurrentLevel = 5
		lv.CurrentMin = level5
		lv.NextExp = level6
	default:
		lv.CurrentLevel = 6
		lv.CurrentMin = level6
		lv.NextExp = levelMax
	}
	lv.CurrentExp = int32(exp)
}
