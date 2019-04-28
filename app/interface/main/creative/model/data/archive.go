package data

import "go-common/app/interface/main/creative/model/archive"

// ArchiveStat for archive stat.
type ArchiveStat struct {
	Play  int64 `json:"play"`
	Dm    int64 `json:"dm"`
	Reply int64 `json:"reply"`
	Like  int64 `json:"like"`
	Coin  int64 `json:"coin"`
	Elec  int64 `json:"elec"`
	Fav   int64 `json:"fav"`
	Share int64 `json:"share"`
}

// ArchiveSource for archive source
type ArchiveSource struct {
	Mainsite int64 `json:"mainsite"`
	Outsite  int64 `json:"outsite"`
	Mobile   int64 `json:"mobile"`
	Others   int64 `json:"others"`
	WebPC    int64 `json:"-"`
	WebH5    int64 `json:"-"`
	IOS      int64 `json:"-"`
	Android  int64 `json:"-"`
}

// ArchiveGroup for archive group.
type ArchiveGroup struct {
	Fans  int64 `json:"fans"`
	Guest int64 `json:"guest"`
}

// ArchiveArea for archive area.
type ArchiveArea struct {
	Location string `json:"location"`
	Count    int64  `json:"count"`
}

// ArchiveData for  single  archive stats.
type ArchiveData struct {
	ArchiveStat   *ArchiveStat           `json:"stat"`
	ArchiveSource *ArchiveSource         `json:"source"`
	ArchiveGroup  *ArchiveGroup          `json:"group"`
	ArchivePlay   *ArchivePlay           `json:"play"`
	ArchiveAreas  []*ArchiveArea         `json:"area"`
	Videos        []*archive.SimpleVideo `json:"videos,omitempty"`
}

// UpBaseStat for up base.
type UpBaseStat struct {
	View  int64 `json:"view"`
	Reply int64 `json:"reply"`
	Dm    int64 `json:"dm"`
	Fans  int64 `json:"fans"`
	Fav   int64 `json:"fav"`
	Like  int64 `json:"like"`
	Share int64 `json:"share"`
	Coin  int64 `json:"coin"`
	Elec  int64 `json:"elec"`
}

// ViewerBase for up base data analysis.
type ViewerBase struct {
	Male         int64 `json:"male"`
	Female       int64 `json:"female"`
	AgeOne       int64 `json:"age_one"`
	AgeTwo       int64 `json:"age_two"`
	AgeThree     int64 `json:"age_three"`
	AgeFour      int64 `json:"age_four"`
	PlatPC       int64 `json:"plat_pc"`
	PlatH5       int64 `json:"plat_h5"`
	PlatOut      int64 `json:"plat_out"`
	PlatIOS      int64 `json:"plat_ios"`
	PlatAndroid  int64 `json:"plat_android"`
	PlatOtherApp int64 `json:"plat_other_app"`
}

// ViewerActionHour for up action data analysis.
type ViewerActionHour struct {
	View     map[int]int `json:"view"`
	Reply    map[int]int `json:"reply"`
	Dm       map[int]int `json:"danmu"`
	Elec     map[int]int `json:"elec"`
	Contract map[int]int `json:"contract"`
}

// Trend for up trend data analysis.
type Trend struct {
	Ty  map[int]int64
	Tag map[int]int64
}

// UpDataIncrMeta for Play/Dm/Reply/Fav/Share/Elec/Coin incr.
type UpDataIncrMeta struct {
	Incr        int            `json:"-"`
	TopAIDList  map[int]int64  `json:"-"`
	TopIncrList map[int]int    `json:"-"`
	Rank        map[int]int    `json:"-"`
	TyRank      map[string]int `json:"-"`
}

const (
	//Play 播放相关.
	Play = int8(1)
	//Dm 弹幕相关.
	Dm = int8(2)
	//Reply 评论相关.
	Reply = int8(3)
	//Share 分享相关.
	Share = int8(4)
	//Coin 投币相关.
	Coin = int8(5)
	//Fav 收藏相关.
	Fav = int8(6)
	//Elec 充电相关.
	Elec = int8(7)
	//Like 点赞相关.
	Like = int8(8)
)

var (
	typeNameMap = map[int8]string{
		Play:  "play",
		Dm:    "dm",
		Reply: "reply",
		Share: "share",
		Coin:  "coin",
		Fav:   "fav",
		Elec:  "elec",
		Like:  "like",
	}
)

//IncrTy return incr data type.
func IncrTy(ty int8) (val string, ok bool) {
	val, ok = typeNameMap[ty]
	return
}

// ArchiveMaxStat 获取单个稿件最多播放、评论、弹幕。。。
type ArchiveMaxStat struct {
	PlayV        int64 `family:"f" qualifier:"play_v" json:"play_v"`
	PlayA        int64 `family:"f" qualifier:"play_a" json:"play_a"`
	CoinV        int64 `family:"f" qualifier:"coin_v" json:"coin_v"`
	CoinA        int64 `family:"f" qualifier:"coin_a" json:"coin_a"`
	LikeV        int64 `family:"f" qualifier:"like_v" json:"like_v"`
	LikeA        int64 `family:"f" qualifier:"like_a" json:"like_a"`
	ReplyV       int64 `family:"f" qualifier:"reply_v" json:"reply_v"`
	ReplyA       int64 `family:"f" qualifier:"reply_a" json:"reply_a"`
	ShareV       int64 `family:"f" qualifier:"share_v" json:"share_v"`
	ShareA       int64 `family:"f" qualifier:"share_a" json:"share_a"`
	FavV         int64 `family:"f" qualifier:"fav_v" json:"fav_v"`
	FavA         int64 `family:"f" qualifier:"fav_a" json:"fav_a"`
	DmV          int64 `family:"f" qualifier:"dm_v" json:"dm_v"`
	DmA          int64 `family:"f" qualifier:"dm_a" json:"dm_a"`
	FromPhoneNum int64 `family:"f" qualifier:"from_phone_num" json:"from_phone_num"`
}
