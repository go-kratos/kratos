package model

import (
	"encoding/json"
	"errors"
	xtime "go-common/library/time"
	"strconv"
	"strings"
)

// EmbedMid is
type EmbedMid struct {
	Mid int64 `json:"mid"`
}

// PassportSummary is
type PassportSummary struct {
	EmbedMid
	TelStatus   int64  `json:"tel_status"`
	CountryID   int64  `json:"country_id"`
	JoinIP      string `json:"join_ip"`
	JoinTime    string `json:"join_time"`
	EmailSuffix string `json:"email_suffix"`
	OriginType  int64  `json:"origin_type"`
	RegType     int64  `json:"reg_type"`
}

// RelationStat is
type RelationStat struct {
	EmbedMid
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}

// BlockSummary is
type BlockSummary struct {
	EmbedMid
	BlockStatus int64  `json:"block_status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// MemberBase is
type MemberBase struct {
	EmbedMid
	Birthday string `json:"birthday"`
	Face     string `json:"face"`
	Name     string `json:"name"`
	Rank     int64  `json:"rank"`
	Sex      int64  `json:"sex"`
	Sign     string `json:"sign"`
}

// MemberOfficial is
type MemberOfficial struct {
	EmbedMid
	Role        int64  `json:"role"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// MemberExp is
type MemberExp struct {
	EmbedMid
	Exp int64 `json:"exp"`
}

// AccountSummary is
type AccountSummary struct {
	MemberBase

	Exp          *MemberExp       `json:"exp"`
	Official     *MemberOfficial  `json:"official"`
	RelationStat *RelationStat    `json:"relation_stat"`
	Block        *BlockSummary    `json:"block"`
	Passport     *PassportSummary `json:"passport"`
}

// NewAccountSummary is
func NewAccountSummary() *AccountSummary {
	return &AccountSummary{
		Exp:          &MemberExp{},
		Official:     &MemberOfficial{},
		Passport:     &PassportSummary{},
		RelationStat: &RelationStat{},
		Block:        &BlockSummary{},
	}
}

func (sum *AccountSummary) String() string {
	b, _ := json.Marshal(sum)
	return string(b)
}

// Key is
func (e *EmbedMid) Key() (string, error) {
	if e.Mid == 0 {
		return "", errors.New("Empty mid")
	}
	return MidKey(e.Mid), nil
}

// Marshal is
func (b *MemberBase) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"birthday": []byte(b.Birthday),
		"face":     []byte(b.Face),
		"mid":      []byte(strconv.FormatInt(b.Mid, 10)),
		"name":     []byte(b.Name),
		"rank":     []byte(strconv.FormatInt(b.Rank, 10)),
		"sex":      []byte(strconv.FormatInt(b.Sex, 10)),
		"sign":     []byte(b.Sign),
	}
	return data, nil
}

// Marshal is
func (o *MemberOfficial) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"official.mid":         []byte(strconv.FormatInt(o.Mid, 10)),
		"official.role":        []byte(strconv.FormatInt(o.Role, 10)),
		"official.title":       []byte(o.Title),
		"official.description": []byte(o.Description),
	}
	return data, nil
}

// Marshal is
func (e *MemberExp) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"exp.mid": []byte(strconv.FormatInt(e.Mid, 10)),
		"exp.exp": []byte(strconv.FormatInt(e.Exp, 10)),
	}
	return data, nil
}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func rpad(s string, c string, l int) string {
	dt := l - len(s)
	if dt <= 0 {
		return s
	}
	return s + strings.Repeat(c, dt)
}

// MidKey is
func MidKey(mid int64) string {
	ms := strconv.FormatInt(mid, 10)
	return rpad(reverse(ms), "0", 16)
}

// Marshal is
func (r *RelationStat) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"relation.mid":       []byte(strconv.FormatInt(r.Mid, 10)),
		"relation.following": []byte(strconv.FormatInt(r.Following, 10)),
		"relation.whisper":   []byte(strconv.FormatInt(r.Whisper, 10)),
		"relation.black":     []byte(strconv.FormatInt(r.Black, 10)),
		"relation.follower":  []byte(strconv.FormatInt(r.Follower, 10)),
	}
	return data, nil
}

// Marshal is
func (b *BlockSummary) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"block.mid":          []byte(strconv.FormatInt(b.Mid, 10)),
		"block.block_status": []byte(strconv.FormatInt(b.BlockStatus, 10)),
		"block.start_time":   []byte(b.StartTime),
		"block.end_time":     []byte(b.EndTime),
	}
	return data, nil
}

// Marshal is
func (p *PassportSummary) Marshal() (map[string][]byte, error) {
	data := map[string][]byte{
		"passport.mid":          []byte(strconv.FormatInt(p.Mid, 10)),
		"passport.tel_status":   []byte(strconv.FormatInt(p.TelStatus, 10)),
		"passport.country_id":   []byte(strconv.FormatInt(p.CountryID, 10)),
		"passport.join_ip":      []byte(p.JoinIP),
		"passport.join_time":    []byte(p.JoinTime),
		"passport.email_suffix": []byte(p.EmailSuffix),
		"passport.reg_type":     []byte(strconv.FormatInt(p.RegType, 10)),
		"passport.origin_type":  []byte(strconv.FormatInt(p.OriginType, 10)),
	}
	return data, nil
}

// Date convert timestamp to date
func Date(in xtime.Time) string {
	return in.Time().Format("2006-01-02")
}

// Datetime convert timestamp to date time
func Datetime(in xtime.Time) string {
	return in.Time().Format("2006-01-02 15:04:05")
}
