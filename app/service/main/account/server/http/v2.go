package http

import (
	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	member "go-common/app/service/main/member/model"
	bm "go-common/library/net/http/blademaster"
)

// v2MyInfo
func v2MyInfo(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	ps, err := accSvc.ProfileWithStat(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	info := &V2MyInfo{}
	info.FromProfile(ps)
	c.JSON(info, nil)
}

// V2MyInfo myinfo.
type V2MyInfo struct {
	Mid            int64             `json:"mid"`
	Name           string            `json:"uname"`
	Face           string            `json:"face"`
	Rank           int32             `json:"rank"`
	Scores         int32             `json:"scores"`
	Coins          float64           `json:"coins"`
	Sex            int32             `json:"sex"`
	Sign           string            `json:"sign"`
	JoinTime       int32             `json:"jointime"`
	Spacesta       int32             `json:"spacesta"`
	Active         int32             `json:"active"`
	Silence        int32             `protobuf:"varint,12,opt,name=Silence,proto3" json:"silence"`
	EmailStatus    int32             `protobuf:"varint,13,opt,name=EmailStatus,proto3" json:"email_status"`
	TelStatus      int32             `protobuf:"varint,14,opt,name=TelStatus,proto3" json:"tel_status"`
	Identification int32             `protobuf:"varint,15,opt,name=Identification,proto3" json:"identification"`
	Moral          int32             `protobuf:"varint,16,opt,name=Moral,proto3" json:"moral"`
	Birthday       string            `protobuf:"bytes,17,opt,name=Birthday,proto3" json:"birthday"`
	Telephone      string            `protobuf:"bytes,18,opt,name=Telephone,proto3" json:"telephone"`
	Level          member.LevelInfo  `protobuf:"bytes,19,opt,name=Level" json:"level_info"`
	Pendant        v1.PendantInfo    `protobuf:"bytes,20,opt,name=Pendant" json:"pendant"`
	Nameplate      v1.NameplateInfo  `protobuf:"bytes,21,opt,name=Nameplate" json:"nameplate"`
	Official       model.OldOfficial `json:"official_verify"`
	Vip            struct {
		Type          int32  `protobuf:"varint,1,opt,name=Type,proto3" json:"vipType"`
		DueDate       int64  `protobuf:"varint,2,opt,name=DueDate,proto3" json:"vipDueDate"`
		DueRemark     string `protobuf:"bytes,3,opt,name=DueRemark,proto3" json:"dueRemark"`
		AccessStatus  int32  `protobuf:"varint,4,opt,name=AccessStatus,proto3" json:"accessStatus"`
		VipStatus     int32  `protobuf:"varint,5,opt,name=VipStatus,proto3" json:"vipStatus"`
		VipStatusWarn string `protobuf:"bytes,6,opt,name=VipStatusWarn,proto3" json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromProfile from profile.
func (i *V2MyInfo) FromProfile(c *model.ProfileStat) {
	i.Mid = c.Mid
	i.Name = c.Name
	switch c.Sex {
	case "男":
		i.Sex = 1
	case "女":
		i.Sex = 2
	default:
		i.Sex = 0
	}
	i.Sign = c.Sign
	i.Face = c.Face
	i.Rank = c.Rank
	i.JoinTime = c.JoinTime
	i.Silence = c.Silence
	if c.Silence == 1 {
		i.Spacesta = -2
	}
	i.EmailStatus = c.EmailStatus
	i.TelStatus = c.TelStatus
	if c.EmailStatus == 1 || c.TelStatus == 1 {
		i.Active = 1
	}
	i.Identification = c.Identification
	i.Coins = c.Coins
	i.Moral = c.Moral
	i.Level.Cur = c.Level
	i.Level.Min = c.LevelExp.Min
	i.Level.NowExp = c.LevelExp.NowExp
	i.Level.NextExp = c.LevelExp.NextExp
	i.Pendant = c.Pendant
	i.Nameplate = c.Nameplate
	i.Official = model.CvtOfficial(c.Official)
	i.Vip.Type = c.Vip.Type
	i.Vip.VipStatus = c.Vip.Status
	i.Vip.DueDate = c.Vip.DueDate
}
