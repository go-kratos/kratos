package resource

import (
	"go-common/app/service/main/archive/api"
	xtime "go-common/library/time"
)

// OverSeasCountry OverSeas Country
var OverSeasCountry = map[string]int64{
	"澳大利亚":       2,
	"泰国":         4,
	"印度":         5,
	"英国":         6,
	"马来西亚":       8,
	"安哥拉":        9,
	"韩国":         10,
	"俄罗斯":        11,
	"新加坡":        12,
	"菲律宾":        13,
	"越南":         14,
	"法国":         15,
	"波兰":         16,
	"荷兰":         17,
	"德国":         18,
	"西班牙":        19,
	"瑞士":         20,
	"欧盟":         21,
	"丹麦":         22,
	"瑞典":         23,
	"意大利":        24,
	"比利时":        25,
	"爱尔兰":        26,
	"芬兰":         27,
	"匈牙利":        28,
	"希腊":         29,
	"保加利亚":       30,
	"奥地利":        31,
	"阿联酋":        32,
	"捷克":         33,
	"南非":         34,
	"以色列":        35,
	"卡塔尔":        36,
	"乌克兰":        37,
	"哈萨克斯坦":      38,
	"葡萄牙":        39,
	"沙特阿拉伯":      40,
	"伊朗":         41,
	"挪威":         42,
	"加拿大":        43,
	"北美地区":       44,
	"叙利亚":        45,
	"科威特":        46,
	"巴林":         47,
	"黎巴嫩":        48,
	"阿曼":         49,
	"约旦":         50,
	"伊拉克":        51,
	"土耳其":        52,
	"罗马尼亚":       53,
	"印度尼西亚":      54,
	"格鲁吉亚":       55,
	"阿塞拜疆":       56,
	"布隆迪":        57,
	"津巴布韦":       58,
	"赞比亚":        59,
	"刚果(金)":      60,
	"巴勒斯坦":       61,
	"立陶宛":        62,
	"斯洛伐克":       63,
	"塞尔维亚":       64,
	"冰岛":         65,
	"斯洛文尼亚":      66,
	"摩尔多瓦":       67,
	"马其顿":        68,
	"列支敦士登":      69,
	"泽西岛":        70,
	"克罗地亚":       71,
	"根西岛":        72,
	"波斯尼亚和黑塞哥维那": 73,
	"爱沙尼亚":       74,
	"拉脱维亚":       75,
	"智利":         76,
	"秘鲁":         77,
	"巴西":         78,
	"吉尔吉斯斯坦":     79,
	"留尼汪岛":       80,
	"马恩岛":        81,
	"直布罗陀":       82,
	"利比亚":        83,
	"亚美尼亚":       84,
	"也门":         85,
	"白俄罗斯":       86,
	"瓜德罗普":       87,
	"卢森堡":        88,
	"马提尼克岛":      89,
	"圭亚那":        90,
	"科索沃":        91,
	"关岛":         92,
	"多米尼加":       93,
	"墨西哥":        94,
	"委内瑞拉":       95,
	"波多黎各":       97,
	"格林纳达":       98,
	"蒙古":         99,
	"新西兰":        100,
	"孟加拉":        101,
	"巴基斯坦":       102,
	"亚太地区":       103,
	"尼泊尔":        104,
	"巴布亚新几内亚":    105,
	"特立尼达和多巴哥":   106,
	"哥伦比亚":       107,
	"阿根廷":        108,
	"斯里兰卡":       109,
	"埃及":         110,
	"厄瓜多尔":       111,
	"哥斯达黎加":      112,
	"乌拉圭":        113,
	"巴巴多斯":       114,
	"巴哈马":        115,
	"圣卢西亚":       116,
	"拉美地区":       117,
	"托克劳群岛":      118,
	"柬埔寨":        119,
	"马尔代夫":       120,
	"阿富汗":        121,
	"新喀里多尼亚":     122,
	"斐济":         123,
	"瓦利斯和富图纳群岛":  124,
	"尼日利亚":       125,
	"阿尔巴尼亚":      126,
	"乌兹别克斯坦":     127,
	"塞浦路斯":       128,
	"圣马力诺":       129,
	"黑山":         130,
	"塔吉克斯坦":      131,
	"马耳他":        132,
	"百慕大":        133,
	"圣文森特和格林纳丁斯": 134,
	"牙买加":        135,
	"多哥":         136,
	"危地马拉":       137,
	"玻利维亚":       138,
	"几内亚":        139,
	"苏里南":        140,
	"利比里亚":       141,
	"肯尼亚":        142,
	"加纳":         143,
	"坦桑尼亚":       144,
	"塞内加尔":       145,
	"马达加斯加":      146,
	"纳米比亚":       147,
	"科特迪瓦":       148,
	"苏丹":         149,
	"喀麦隆":        150,
	"马拉维":        151,
	"加蓬":         152,
	"马里":         153,
	"贝宁":         154,
	"乍得":         155,
	"博茨瓦纳":       156,
	"佛得角":        157,
	"卢旺达":        158,
	"刚果(布)":      159,
	"乌干达":        160,
	"莫桑比克":       161,
	"冈比亚":        162,
	"莱索托":        163,
	"毛里求斯":       164,
	"非洲地区":       165,
	"阿尔及利亚":      166,
	"斯威士兰":       167,
	"布基纳法索":      168,
	"塞拉利昂":       169,
	"索马里":        170,
	"尼日尔":        171,
	"中非":         172,
	"南苏丹":        173,
	"赤道几内亚":      174,
	"塞舌尔":        175,
	"吉布提":        176,
	"摩洛哥":        177,
	"毛里塔尼亚":      178,
	"科摩罗":        179,
	"英属印度洋领地":    180,
	"开曼群岛":       181,
	"突尼斯":        182,
	"马约特":        183,
	"老挝":         184,
	"缅甸":         185,
	"文莱":         186,
	"瑙鲁":         187,
	"瓦努阿图":       188,
	"不丹":         189,
	"密克罗尼西亚联邦":   190,
	"法属波利尼西亚":    191,
	"东帝汶":        192,
	"汤加":         193,
	"北马里亚纳群岛":    194,
	"格陵兰":        195,
	"英属维尔京群岛":    196,
	"法罗群岛":       197,
	"纽埃岛":        198,
	"福克兰群岛":      199,
	"特克斯和凯科斯群岛":  200,
	"洪都拉斯":       201,
	"库拉索":        202,
	"荷兰加勒比":      203,
	"马绍尔群岛":      204,
	"库克群岛":       205,
	"巴拿马":        206,
	"法属圣马丁":      207,
	"美属维尔京群岛":    208,
	"美属萨摩亚":      209,
	"萨尔瓦多":       210,
	"荷属圣马丁":      211,
	"阿鲁巴":        212,
	"巴拉圭":        213,
	"多米尼克":       214,
	"安提瓜和巴布达":    215,
	"安圭拉":        216,
	"圣基茨和尼维斯":    217,
	"圣皮埃尔和密克隆群岛": 218,
	"土库曼斯坦":      219,
	"奥兰群岛":       220,
	"摩纳哥":        221,
	"法属圭亚那":      222,
	"安道尔":        223,
	"梵蒂冈":        224,
	"海地":         225,
	"共享地址":       226,
	"所罗门群岛":      227,
	"基里巴斯":       228,
	"帕劳":         229,
	"诺福克岛":       230,
	"萨摩亚":        231,
	"阿里云骨干网":     232,
	"本机地址":       233,
	"伯利兹":        234,
	"尼加拉瓜":       235,
	"古巴":         236,
	"圣多美和普林西比":   237,
	"几内亚比绍":      238,
	"本地链路":       239,
	"朝鲜":         240,
	"埃塞俄比亚":      241,
	"厄立特里亚":      242,
	"蒙塞拉特岛":      243,
	"图瓦卢":        244,
	"圣诞岛":        245,
	"圣巴泰勒米岛":     246,
}

// AsgTypePic AsgTypeVideo
const (
	AsgTypePic   = int8(0)
	AsgTypeVideo = int8(1)
	// pgc mobile
	AsgTypeURL     = int8(2)
	AsgTypeBangumi = int8(3)
	AsgTypeLive    = int8(4)
	AsgTypeGame    = int8(5)
	AsgTypeAv      = int8(6)
	AsgTypeTopic   = int8(7)
	// content type
	FromManager = int8(0)
	FromCpm     = int8(1)
)

// Assignment struct
type Assignment struct {
	ID         int        `json:"id"`
	ContractID string     `json:"contract_id"`
	ResID      int        `json:"-"`
	PosNum     int        `json:"pos_num"`
	Name       string     `json:"name"`
	Pic        string     `json:"pic"`
	LitPic     string     `json:"litpic"`
	URL        string     `json:"url"`
	Rule       string     `json:"-"`
	Style      int32      `json:"style"`
	IsAd       bool       `json:"is_ad,omitempty"`
	Archive    *api.Arc   `json:"archive,omitempty"`
	Aid        int64      `json:"-"`
	Weight     int        `json:"-"`
	Atype      int8       `json:"-"`
	MTime      xtime.Time `json:"-"`
	Agency     string     `json:"agency"`
	Label      string     `json:"label"`
	Intro      string     `json:"intro"`
	// cpm
	CreativeType int8       `json:"creative_type"`
	RequestID    string     `json:"request_id,omitempty"`
	CreativeID   int64      `json:"creative_id,omitempty"`
	SrcID        int64      `json:"src_id,omitempty"`
	ShowURL      string     `json:"show_url,omitempty"`
	ClickURL     string     `json:"click_url,omitempty"`
	Area         int8       `json:"area"`
	IsAdLoc      bool       `json:"is_ad_loc"`
	AdCb         string     `json:"ad_cb"`
	Title        string     `json:"title"`
	ServerType   int8       `json:"server_type"`
	CmMark       int8       `json:"cm_mark"`
	IsCpm        bool       `json:"-"`
	STime        xtime.Time `json:"stime"`
	Mid          string     `json:"mid"`
}

// Relation struct
type Relation struct {
	*api.Arc
	// cpm
	RequestID  string `json:"request_id,omitempty"`
	CreativeID int64  `json:"creative_id,omitempty"`
	SrcID      int64  `json:"src_id,omitempty"`
	ShowURL    string `json:"show_url,omitempty"`
	ClickURL   string `json:"click_url,omitempty"`
	Area       int8   `json:"area"`
	IsAdLoc    bool   `json:"is_ad_loc"`
	AdCb       string `json:"ad_cb"`
	ResID      int    `json:"resource_id"`
	IsAd       bool   `json:"is_ad"`
}

// Position struct
type Position struct {
	Pos     []*Loc
	Counter int `json:"-"`
}

// Res struct
type Res struct {
	ID       int    `json:"-"`
	Platform int    `json:"-"`
	Name     string `json:"-"`
	Parent   int    `json:"-"`
	Counter  int    `json:"-"`
	Position int    `json:"-"`
	// ass
	Assignments []*Assignment `json:"-"`
}

// Loc struct
type Loc struct {
	ID     int `json:"-"`
	PosNum int `json:"-"`
}

// ArgRess ArgRess
type ArgRess struct {
	Pf    int     `form:"pf" validate:"min=0"`
	Ids   []int64 `form:"ids,split" validate:"min=1,dive,gte=1"`
	Mid   int64
	Sid   string
	IP    string
	Buvid string
}

// ArgRes ArgRes
type ArgRes struct {
	Pf    int   `form:"pf" validate:"min=0"`
	ID    int64 `form:"id" validate:"min=1"`
	Mid   int64
	Sid   string
	IP    string
	Buvid string
}

// ArgAid ArgAid
type ArgAid struct {
	Aid   int64 `form:"aid" validate:"min=1"`
	Mid   int64
	Sid   string
	IP    string
	Buvid string
}
