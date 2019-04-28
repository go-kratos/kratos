package model

import (
	"html/template"
	"strconv"
	"time"

	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
)

var encodeTbl = []int64{
	0xf3a97cb, 0x8aed379, 0xedf369a, 0x5c82647, 0xcaf6987, 0xad28536, 0xe5f2a7b, 0x72e85df,
	0xac3d972, 0xca97fe5, 0xbcf5473, 0x85ad732, 0xcd6b324, 0xd549a72, 0xe72ab89, 0xfa6dc53,
	0xa8e752b, 0xa73d25f, 0xcad8296, 0xb35f689, 0x7ce594b, 0x59ca743, 0xc4ab2d7, 0x9c8adf6,
	0x93746c2, 0x6cea579, 0xcd36b75, 0x64a973e, 0xfa49c56, 0xb45f2d9, 0x72a56f8, 0x43d6fa9,
	0x354cf2a, 0x26bf5d8, 0x39f64ad, 0xa4fd326, 0x39247d6, 0xec67f95, 0xed8c9b4, 0x29637db,
	0xefabc54, 0xa9ed87c, 0xc2864ea, 0xf32d475, 0x53b6897, 0xe7f94b8, 0x7a4cfd2, 0x9a82e65,
	0x369b7a4, 0x2cae6d4, 0xc7fba36, 0xd3e7846, 0xd324ba5, 0x7c56f24, 0x598c3af, 0x39fd4ae,
	0x3b6472c, 0x2f9a8be, 0x9fcab42, 0x8f34aeb, 0x9e8b372, 0x8c42b9e, 0xf9b574c, 0x7c693fa,
	0x245fc67, 0x823f4ce, 0x957f84d, 0xe529a87, 0xb625ead, 0xbd4f6a9, 0x863ca52, 0xd762cef,
	0x8d6c479, 0xbc4f579, 0xa486fdc, 0xcd6f289, 0xda3b629, 0x4fce523, 0x2e8db97, 0xc3bf769,
	0x9c64d7f, 0x52db6f7, 0x95cdf8e, 0xc872fe9, 0x964de53, 0x2bef897, 0xb7a962c, 0x38d72be,
	0x26fa89c, 0x58b742e, 0xa3bd967, 0x3cae942, 0x4d3fb9c, 0xaf59ed3, 0x6f8379d, 0x2bf46d7,
	0xcdbe243, 0x3754bf9, 0x82f9dc5, 0x8a46ef5, 0x5d48ac9, 0x9e6ca3d, 0xfec5a3b, 0x57dafe3,
	0x82ed7a9, 0xbc3d687, 0x89ecbf7, 0x738549a, 0x928746c, 0x9cb7e83, 0xc85f9a7, 0x2947c8e,
	0xba689fd, 0xebc4893, 0xa62cf7e, 0xa8e3cb5, 0xe47589d, 0x792edaf, 0x4635c2d, 0xa2c6bfe,
	0xc456daf, 0x2d65f47, 0xf9ce625, 0x74a8b62, 0x9d728f5, 0x3e4a29d, 0x62a589c, 0x83cb629,
	0xce5b6d3, 0x2fda9ce, 0x87af3bc, 0x837a695, 0xf935da4, 0x48b6ea2, 0x52dc4e9, 0x82a537b,
	0xe23456f, 0x6cbdafe, 0x97bf34d, 0x4c72ad8, 0xa5c4982, 0x8afb76e, 0x895fca6, 0x85abd24,
	0xae2475d, 0xf3c5eb8, 0xb4d2ef5, 0xbda463e, 0xf392a5e, 0x7a9fd58, 0xead48f6, 0x8a62537,
	0x6c35ba2, 0x7589e24, 0xd24ec93, 0x6bc42a5, 0x34d9f87, 0xed3578f, 0x87452fa, 0x5439fca,
	0x29b37c8, 0x8fe4c3b, 0x4c5368d, 0x58acf9b, 0x69c3ad2, 0xaf3827b, 0x328e46b, 0xbef7ca9,
	0xda592c4, 0x45f7db2, 0xcb65a3d, 0x4578ec3, 0xc6deab9, 0xb689edc, 0x4aed59f, 0x25b9af7,
	0x9b6d48f, 0x6de79bf, 0x249fa5e, 0x269a7ef, 0xd9e62a7, 0xb9a86d2, 0x539b72c, 0x8fa9ebc,
	0xec397f5, 0xdbac4e2, 0x938e6fd, 0xe8a734f, 0xe4b8d7f, 0x84cd9b3, 0x75c6ef4, 0x956378c,
	0x43f2d78, 0x74e9253, 0x25dbef4, 0xb7e26f9, 0x93b2c6d, 0x2faeb76, 0x3b278de, 0x6b5948f,
	0x4967358, 0x49f3a7e, 0x7596ec4, 0x98cabf5, 0x95c638e, 0x6d258b7, 0x97e8b3f, 0x5ab7823,
	0x53b6a89, 0xa3bc579, 0xac45d36, 0xcea9b28, 0x98f2356, 0xd694a2e, 0xf732e8a, 0xe7463d5,
	0xf5ec9a8, 0x6dba984, 0xc798e5a, 0x6e9382b, 0xeac3249, 0x5238b9a, 0xd632eaf, 0xa92b685,
	0xbcae435, 0x9726fd8, 0x3fcbea4, 0x5e9da23, 0xb93a4f7, 0x327d84c, 0x5db932f, 0x86274de,
	0xa54bd72, 0x63f2ed5, 0x6d37285, 0xb4fe7c9, 0x549a6b3, 0x3b592ec, 0x73d456b, 0x49253b7,
	0x2da9b8c, 0xb85642d, 0x37489ca, 0x726fe3b, 0x4ce6ad2, 0x376becd, 0x6f43bec, 0xf96dba4,
	0xebc8d72, 0xf59b4ca, 0x263547f, 0xabcd87e, 0x3fd25ae, 0xc6f4b38, 0x36cd978, 0x6e94a37,
}

// Catalog catalog map
var Catalog = map[int8]string{
	1: "system",
	2: "bangumi",
	3: "news",
}

// Sid sid string.
type Sid string

// Valid valid sid.
func (sid Sid) Valid() (b bool) {
	var (
		interval int64
		msec     int64
		_xorKey  int64
		hKey     int64
		_dec     int64
		_ts      int64
		offset   int
		ms       int64
		err      error
	)
	if len(sid) != 8 {
		return false
	}
	if interval, err = strconv.ParseInt(string(sid)[0:6], 36, 64); err != nil {
		return false
	}
	if msec, err = strconv.ParseInt(string(sid)[6:8], 36, 64); err != nil {
		return false
	}
	_xorKey = encodeTbl[msec%256]
	hKey = 1 << uint(28+msec%3)
	_dec = (interval ^ _xorKey)
	if (_dec & 0x70000000) != hKey {
		return false
	}
	_ts = (_dec^hKey)*1000 + msec
	_, offset = time.Now().Zone()
	ms = (time.Now().UnixNano() / int64(time.Millisecond)) + (int64(offset)/60)*60000 // GMT Timestamp
	return _ts <= (ms - 1388534400000 + 300000)
}

// Create create sid.
func (sid Sid) Create() (re Sid) {
	var (
		offset    int
		ms        int64
		_interval int64
		msec      int64
		msesStr   string
		_xorKey   int64
		hKey      int64
		tsEncode  string
	)
	_, offset = time.Now().Zone()
	ms = (time.Now().UnixNano() / int64(time.Millisecond)) - (int64(offset)/60)*60000 // GMT Timestamp
	_interval = ms/1000 - 1388534400
	msec = ms % 1000
	_xorKey = encodeTbl[msec%256]
	hKey = 1 << uint(28+msec%3)
	tsEncode = strconv.FormatInt(((hKey | _interval) ^ _xorKey), 36)
	for len(tsEncode) < 6 {
		tsEncode = "0" + tsEncode
	}
	msesStr = strconv.FormatInt(msec, 36)
	if len(msesStr) < 2 {
		msesStr = "0" + msesStr
	}
	sid = Sid(tsEncode + msesStr)
	return sid
}

// Item item struct.
type Item struct {
	Content    string `xml:",cdata"`
	Tooltip    string `xml:"tooltip,attr"`
	Bgcolor    string `xml:"bgcolor,attr"`
	Catalog    string `xml:"catalog,attr"`
	ResourceID string `xml:"resourceid,attr"`
	SrcID      string `xml:"srcid,attr"`
	ID         string `xml:"id,attr"`
}

// Carousel carousel struct.
type Carousel struct {
	Backcolor string
	Fontcolor string
	Hint      string
	Link      string
	Content   string
	Category  string
}

// Player player struct.
type Player struct {
	IP      string
	Zoneid  int64
	Country string
	Isp     string // 运营商暂时不需要
	Login   bool
	Time    int64
	ZoneIP  string
	// user info
	Name           string
	User           int64
	UserHash       string
	Money          string
	Face           string
	IsAdmin        bool
	Upermission    string
	Level          string
	LevelInfo      template.HTML
	Vip            template.HTML
	OfficialVerify template.HTML
	BlockTime      int64
	// archive info
	Aid            int64
	Typeid         int32
	Maxlimit       int
	Click          int
	FwClick        int32
	Duration       string
	Arctype        string
	APermission    bool
	SuggestComment bool
	Chatid         int64
	Vtype          string
	Oriurl         string
	Pid            int64
	AllowBp        bool
	Bottom         int8
	Acceptguest    bool
	Acceptaccel    bool
	Cache          bool
	CacheDispatch  bool
	BrTCP          string
	BrWs           string
	BrWss          string
	DefaultDm      int8
	//progress
	LastPlayTime int64
	LastCid      int64
	Role         string
	// has next page
	HasNext     int8
	OnlineCount int64
	// dm mask
	MaskNew template.HTML
	// subtitle
	Subtitle template.HTML
	// player icon
	PlayerIcon template.HTML
	// view points
	ViewPoints template.HTML
}

// Progress progress struct.
type Progress struct {
	Cid int64 `json:"cid"`
	Pro int64 `json:"pro"`
}

// Policy policy struct.
type Policy struct {
	ID        int64  `json:"id"`
	Des       string `json:"description"`
	Type      string `json:"type"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Mtime     string `json:"mtime"`
	StartTime time.Time
	EndTime   time.Time
	MtimeTime time.Time
	Items     []*Pitem `json:"items"`
}

// Pitem pitem struct
type Pitem struct {
	ID      int64  `json:"item_id"`
	Data    string `json:"data"`
	Comment string `json:"comment"`
	ExtData string `json:"ext_data"`
	Ver     int64  `json:"ver"`
}

// Param param struct.
type Param struct {
	Name  string
	Value string
}

// BlockTime block time struct
type BlockTime struct {
	BlockStatus    int   `json:"block_status"`
	BlockedForever bool  `json:"blocked_forever"`
	BlockedEnd     int64 `json:"blocked_end"`
}

// Videoshot player video shot struct
type Videoshot struct {
	*archive.Videoshot
	Index []uint16 `json:"index,omitempty"`
}

// PlayURLToken playurl token.
type PlayURLToken struct {
	From  string `json:"from"`
	Ts    int64  `json:"ts"`
	Aid   int64  `json:"aid"`
	Cid   int64  `json:"cid"`
	Mid   int64  `json:"mid"`
	VIP   int    `json:"vip"`
	SVIP  int    `json:"svip"`
	Owner int    `json:"owner"`
	Fcs   string `json:"fcs"`
	Token string `json:"token"`
}

// VIPInfo vip info.
type VIPInfo struct {
	Type          int32  `json:"vipType"`
	DueDate       int64  `json:"vipDueDate"`
	DueRemark     string `json:"dueRemark"`
	AccessStatus  int32  `json:"accessStatus"`
	VipStatus     int32  `json:"vipStatus"`
	VipStatusWarn string `json:"vipStatusWarn"`
}

// Official official.
type Official struct {
	Type int8   `json:"type"`
	Desc string `json:"desc"`
}

// Points is
type Points struct {
	Type    int    `json:"type"`
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Content string `json:"content"`
}

const (
	// TpWithUinfo tpl with user info.
	TpWithUinfo = `<ip>{{.IP}}</ip>
<zoneid>{{.Zoneid}}</zoneid>
<zoneip>{{.ZoneIP}}</zoneip>
<country>{{.Country}}</country>
<login>{{.Login}}</login>
<time>{{.Time}}</time>
<name>{{.Name}}</name>
<user>{{.User}}</user>
<user_hash>{{.UserHash}}</user_hash>
<money>{{.Money}}</money>
<face>{{.Face}}</face>
<isadmin>{{.IsAdmin}}</isadmin>
<permission>{{.Upermission}}</permission>
<level>{{.Level}}</level>
<level_info>{{.LevelInfo}}</level_info>
<vip>{{.Vip}}</vip>
<official_verify>{{.OfficialVerify}}</official_verify>
<block_time>{{.BlockTime}}</block_time>
<lastplaytime>{{.LastPlayTime}}</lastplaytime>
<lastcid>{{.LastCid}}</lastcid>
<aid>{{.Aid}}</aid>
<typeid>{{.Typeid}}</typeid>
<vtype>{{.Vtype}}</vtype>
<oriurl>{{.Oriurl}}</oriurl>
<suggest_comment>{{.SuggestComment}}</suggest_comment>
<server>chat.bilibili.com</server>
<maxlimit>{{.Maxlimit}}</maxlimit>
<click>{{.Click}}</click>
<fw_click>{{.FwClick}}</fw_click>
<chatid>{{.Chatid}}</chatid>
<pid>{{.Pid}}</pid>
<duration>{{.Duration}}</duration>
<arctype>{{.Arctype}}</arctype>
<allow_bp>{{.AllowBp}}</allow_bp>
<bottom>{{.Bottom}}</bottom>
<shot>false</shot>
<sinapi>1</sinapi>
<acceptguest>{{.Acceptguest}}</acceptguest>
<acceptaccel>{{.Acceptaccel}}</acceptaccel>
<cache>{{.Cache}}</cache>
<broadcast_tcp>{{.BrTCP}}</broadcast_tcp>
<broadcast_ws>{{.BrWs}}</broadcast_ws>
<broadcast_wss>{{.BrWss}}</broadcast_wss>
<default_dm>{{.DefaultDm}}</default_dm>
<dm_host>0://comment.bilibili.com,1://comment.bilibili.com/rc</dm_host>
<role>{{.Role}}</role>
<has_next>{{.HasNext}}</has_next>
<online_count>{{.OnlineCount}}</online_count>
<dm_mask></dm_mask>
<mask_new>{{.MaskNew}}</mask_new>
<subtitle>{{.Subtitle}}</subtitle>
<player_icon>{{.PlayerIcon}}</player_icon>
<view_points>{{.ViewPoints}}</view_points>
`

	// TpWithNoUinfo tpl without user info.
	TpWithNoUinfo = `<ip>{{.IP}}</ip>
<zoneip>{{.ZoneIP}}</zoneip>
<zoneid>{{.Zoneid}}</zoneid>
<country>{{.Country}}</country>
<login>{{.Login}}</login>
<time>{{.Time}}</time>
<lastplaytime>0</lastplaytime>
<lastcid>0</lastcid>
<aid>{{.Aid}}</aid>
<typeid>{{.Typeid}}</typeid>
<vtype>{{.Vtype}}</vtype>
<oriurl>{{.Oriurl}}</oriurl>
<suggest_comment>{{.SuggestComment}}</suggest_comment>
<server>chat.bilibili.com</server>
<maxlimit>{{.Maxlimit}}</maxlimit>
<click>{{.Click}}</click>
<fw_click>{{.FwClick}}</fw_click>
<chatid>{{.Chatid}}</chatid>
<pid>{{.Pid}}</pid>
<duration>{{.Duration}}</duration>
<arctype>{{.Arctype}}</arctype>
<allow_bp>{{.AllowBp}}</allow_bp>
<bottom>{{.Bottom}}</bottom>
<shot>false</shot>
<sinapi>1</sinapi>
<acceptguest>{{.Acceptguest}}</acceptguest>
<acceptaccel>{{.Acceptaccel}}</acceptaccel>
<cache>{{.Cache}}</cache>
<broadcast_tcp>{{.BrTCP}}</broadcast_tcp>
<broadcast_ws>{{.BrWs}}</broadcast_ws>
<broadcast_wss>{{.BrWss}}</broadcast_wss>
<default_dm>{{.DefaultDm}}</default_dm>
<dm_host>0://comment.bilibili.com,1://comment.bilibili.com/rc</dm_host>
<role>0</role>
<has_next>{{.HasNext}}</has_next>
<online_count>{{.OnlineCount}}</online_count>
<dm_mask></dm_mask>
<mask_new>{{.MaskNew}}</mask_new>
<subtitle>{{.Subtitle}}</subtitle>
<player_icon>{{.PlayerIcon}}</player_icon>
<view_points>{{.ViewPoints}}</view_points>
`
)

// View .
type View struct {
	*arcmdl.Arc
	Pages []*arcmdl.Page `json:"pages"`
}
