package weeklyhonor

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// HonorSub honor subscribe state subscribed
	HonorSub uint8 = iota
	// HonorUnSub honor subscribe state unsubscribed
	HonorUnSub
	layout = "20060102"
	// when edit, add new line below the line you want edit and fill the end time field
	// if Word changed,use new hid,else use the same.
	hStr = `1|S|史诗级选手|本周上过全站排行榜前三名|SSR|20181019|
2|666|惊不惊喜意不意外|本周上过全站排行榜第%d名|SSR|20181019|
3|233|天选之人|本周上过全站排行榜前%d名|SSR|20181019|
4|牛|给大佬递茶|本周上过全站排行榜前%d名|SSR|20181019|
5|史|少年创造奇迹|本周累计总播放量达成了%s的小目标|SSR|20181019|
6|火|「%s」不火天理难容|本周总粉丝数达成%d百万的小目标|SSR|20181019|
7|光|一入动态深似海||SSR|20181019|
8|希|全村人的希望|本周上过全站%s区排行榜前%d名|SR|20181019|20181214
8|希|全村人的希望|本周上过%s区排行榜前%d名|R|20181214|
9|妥|定个小目标，破%d万+|本周累计总播放量达成%d万的小目标|SR|20181019|
10|强|你上辈子大概拯救过世界|本周总粉丝数达成%d万的小目标|SR|20181019|
11|Boom|请收下我的膝盖|本周新增播放量破百万|SR|20181019|
12|囍|此生无悔爱「%s」|本周新增粉丝数破万|SR|20181019|
13|心|给你我的小心心|本周稿件被点赞破五千|SR|20181019|
14|爆|我吹爆这个up主|本周稿件被转发破千|SR|20181019|
15|富|我要让世界知道此人被我承包|本周稿件收到硬币数破三千|SR|20181019|
16|爱|承包这个动态日更up主！||SR|20181019|
17|妙|定个小目标，破%d万+|本周累计总播放量达成%d万的小目标|R|20181019|
18|撩|夭寿啦～这个up好会撩人|本周总粉丝数达成%d万的小目标|R|20181019|
19|星|下一个巨星就是你|本周新增播放量破十万|R|20181019|
20|怦|糟了是心动的感觉|本周新增粉丝数破千|R|20181019|
21|赞|为你打call|本周稿件被点赞破千|R|20181019|
22|鸣|火钳刘明|本周稿件被转发破百|R|20181019|
23|劲|这个up主大概使用了洪荒之力|本周发布稿件数达五个以上|R|20181019|20181214
23|劲|这个up主大概使用了洪荒之力|本周发布稿件数达五个以上|A|20181214|
24|发|下一个大佬就是你|本周稿件收到硬币数破千|R|20181019|
25|call|今天也是可爱的动态up呢！||R|20181019|
26|盯|确认过眼神，是有潜力的up|本周累计总播放量达成%d千的小目标|A|20181019|
27|抱|大大，我是你的腿部挂件|本周总粉丝数达成%d千的小目标|A|20181019|
28|Wow|这只up主有点东西|本周新增播放量破万|A|20181019|
29|嗷|掐指一算，必成大事|本周新增粉丝数破百|A|20181019|
30|仙|给「%s」献上膝盖|本周稿件被点赞破五百|A|20181019|
31|转|前方“转发”高能预警|本周稿件被转发上五十|A|20181019|
32|肝|肝稿模式，启动！|本周发布稿件数达到三个及以上|A|20181019|20181214
32|肝|肝稿模式，启动！|本周发布稿件数达到三个及以上|B|20181214|
33|美|美滋滋，我有个大胆的想法|本周稿件收到硬币数破百|A|20181019|
34|和|本周优秀动态，安排上了！||A|20181019|
35|稳|获得粉丝稳稳的爱||B|20181019|
36|竞|播放1万+还会远吗？||B|20181019|
37|辛|只要有你们，再辛苦的时光都元气满满||B|20181019|
38|秀|加油打气你最棒||B|20181019|
39|A+|被大家的转发种草了||B|20181019|
40|更|爷爷你关注的up主更新了||B|20181019|
41|囤|囤着硬币干大事||B|20181019|
42|奋|秘技！左右横跳发动态||B|20181019|
43|粉|在世界中心呼唤粉丝||B|20181019|
44|进|向着播放破千的目标前进吧||B|20181019|
45|励|粉丝的鼓励你收到了嘛||B|20181019|
46|勤|业精于勤，积少成多||C|20181019|
47|冲|为了更多的赞冲鸭||C|20181019|
48|享|分享是走向世界第一步||C|20181019|
49|币|我也收到硬币啦||C|20181019|
50|等|和我签订契约，成为动态区UP主吧！||C|20181019|
51|孤|生活已经如此艰难，有些事情就不要揭穿了||D|20181019|
52|懵|我好像把我的播放量弄丢了……||D|20181019|
53|槑|呆呆，我的点赞去哪了||D|20181019|
54|空|我的转发空空荡荡…||D|20181019|
55|鸽|鸽了鸽了！投稿是不可能投稿——下周就投稿||D|20181019|
56|静|钱袋静静地躺着||D|20181019|
57|新|来吧！打破零动态惨案！||D|20181019|
58|迟|数据宝宝正在非常努力地奔跑——过会会儿再看啦||E|20181019|
59|炫|疯狂Pick你！！|本周上过全站排行榜前%d名|SR|20181214|`
)

// Honor weeklyhonor info.
type Honor struct {
	HID        int            `json:"hid"`
	MID        int64          `json:"mid"`
	HonorCount int64          `json:"honor_count"`
	Uname      string         `json:"uname"`
	Face       string         `json:"face"`
	Word       string         `json:"word"`
	Text       string         `json:"text"`
	Desc       string         `json:"desc"`
	Priority   string         `json:"priority"`
	ShareToken string         `json:"share_token"`
	RiseStage  *RiseStage     `json:"rise_stage"`
	SubState   uint8          `json:"sub_state"`
	LoveFans   []*accmdl.Info `json:"love_fans"`
	PlayFans   []*accmdl.Info `json:"play_fans"`
	NewArchive *api.Arc       `json:"new_archive"`
	HotArchive *api.Arc       `json:"hot_archive"`
	DateBegin  xtime.Time     `json:"date_begin"`
	DateEnd    xtime.Time     `json:"date_end"`
}

// RiseStage .
type RiseStage struct {
	Play  int `json:"play"`
	Like  int `json:"like"`
	Fans  int `json:"fans"`
	Coin  int `json:"coin"`
	Share int `json:"share"`
}

// HonorWord .
type HonorWord struct {
	ID       int        `json:"id"`
	Word     string     `json:"word"`
	Text     string     `json:"text"`
	Desc     string     `json:"desc"`
	Priority string     `json:"priority"`
	Start    xtime.Time `json:"start"`
	End      xtime.Time `json:"end"`
}

// HMap get hid-honorWord map
func HMap() map[int][]*HonorWord {
	h := make(map[int][]*HonorWord)
	sArr := strings.Split(hStr, "\n")
	for _, v := range sArr {
		ws := strings.Split(v, "|")
		id, err := strconv.Atoi(ws[0])
		if err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", ws[0], err)
		}
		wd := HonorWord{
			ID:       id,
			Word:     ws[1],
			Text:     ws[2],
			Desc:     ws[3],
			Priority: ws[4],
		}
		h[id] = append(h[id], &wd)
		start, err := time.ParseInLocation(layout, ws[5], time.Local)
		if err != nil {
			continue
		}
		wd.Start = xtime.Time(start.Unix())
		end, err := time.ParseInLocation(layout, ws[6], time.Local)
		if err != nil {
			continue
		}
		wd.End = xtime.Time(end.Unix())
	}
	return h
}

// HonorStat up honor stats form hbase.
type HonorStat struct {
	Play         int32 `family:"f" qualifier:"play" json:"play"`
	PlayLastW    int32 `family:"f" qualifier:"play_last_w" json:"play_last_w"`
	Fans         int32 `family:"f" qualifier:"fans" json:"fans"`
	FansLastW    int32 `family:"f" qualifier:"fans_last_w" json:"fans_last_w"`
	PlayInc      int32 `family:"f" qualifier:"play_inc" json:"play_inc"`
	FansInc      int32 `family:"f" qualifier:"fans_inc" json:"fans_inc"`
	LikeInc      int32 `family:"f" qualifier:"like_inc" json:"like_inc"`
	ShInc        int32 `family:"f" qualifier:"sh_inc" json:"sh_inc"`
	CoinInc      int32 `family:"f" qualifier:"coin_inc" json:"coin_inc"`
	AvsInc       int32 `family:"f" qualifier:"avs_inc" json:"avs_inc"`
	DyInc        int32 `family:"f" qualifier:"dy_inc" json:"dy_inc"`
	Act1         int32 `family:"f" qualifier:"act1" json:"act1"`
	Act2         int32 `family:"f" qualifier:"act2" json:"act2"`
	Act3         int32 `family:"f" qualifier:"act3" json:"act3"`
	Dr1          int32 `family:"f" qualifier:"dr1" json:"dr1"`
	Dr2          int32 `family:"f" qualifier:"dr2" json:"dr2"`
	Dr3          int32 `family:"f" qualifier:"dr3" json:"dr3"`
	HottestAvNew int32 `family:"f" qualifier:"hottest_av_new" json:"hottest_av_new"`
	HottestAvInc int32 `family:"f" qualifier:"hottest_av_inc" json:"hottest_av_inc"`
	HottestAvAll int32 `family:"f" qualifier:"hottest_av_all" json:"hottest_av_all"`
	Rank0        int32 `family:"r" qualifier:"rank0" json:"rank0"`
	Rank1        int32 `family:"r" qualifier:"rank1" json:"rank1"`
	Rank3        int32 `family:"r" qualifier:"rank3" json:"rank3"`
	Rank4        int32 `family:"r" qualifier:"rank4" json:"rank4"`
	Rank5        int32 `family:"r" qualifier:"rank5" json:"rank5"`
	Rank36       int32 `family:"r" qualifier:"rank36" json:"rank36"`
	Rank119      int32 `family:"r" qualifier:"rank119" json:"rank119"`
	Rank129      int32 `family:"r" qualifier:"rank129" json:"rank129"`
	Rank155      int32 `family:"r" qualifier:"rank155" json:"rank155"`
	Rank160      int32 `family:"r" qualifier:"rank160" json:"rank160"`
	Rank168      int32 `family:"r" qualifier:"rank168" json:"rank168"`
	Rank181      int32 `family:"r" qualifier:"rank181" json:"rank181"`
}

// HonorLog .
type HonorLog struct {
	ID    int64      `json:"id"`
	MID   int64      `json:"mid"`
	HID   int        `json:"hid"`
	Count int64      `json:"count"`
	CTime xtime.Time `json:"ctime"`
	MTime xtime.Time `json:"mtime"`
}

// GenHonor .
func (hs *HonorStat) GenHonor(mid int64, distinctID int) int {
	ids := hs.priority()
	log.Info("GenHonor mid(%d) ids(%v) lastid(%d)", mid, ids, distinctID)
	if len(ids) == 1 {
		return ids[0]
	}
	var hids []int
	for _, id := range ids {
		if id != distinctID {
			hids = append(hids, id)
		}
	}
	if len(hids) == 1 {
		return hids[0]
	}
	if len(hids) > 1 {
		rand.Seed(mid)
		var mu sync.Mutex
		mu.Lock()
		rnd := rand.Intn(len(hids))
		mu.Unlock()
		return hids[rnd]
	}
	return 0
}

func (hs *HonorStat) priority() []int {
	ids, _ := hs.PrioritySSR()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PrioritySR()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PriorityR()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PriorityA()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PriorityB()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PriorityC()
	if len(ids) > 0 {
		return ids
	}
	ids, _ = hs.PriorityD()
	return ids
}

// PrioritySSR .
func (hs *HonorStat) PrioritySSR() ([]int, *RiseStage) {
	ids := make([]int, 0)
	rs := new(RiseStage)
	hid := 0
	if hs.Rank0 > 0 {
		switch {
		case hs.Rank0 <= 3:
			hid = 1
		case hs.Rank0 == 6 || hs.Rank0 == 66:
			hid = 2
		case hs.Rank0 == 23:
			hid = 3
		case hs.Rank0 <= 10:
			hid = 4
		}
		if hid != 0 {
			ids = append(ids, hid)
		}
	}
	if (hs.Play > 100000000 && hs.PlayLastW < 100000000) || (hs.Play > 10000000 && hs.PlayLastW < 10000000) {
		ids = append(ids, 5)
	}
	if hs.Fans > 1000000 && hs.FansLastW < 1000000 {
		ids = append(ids, 6)
	}
	if hs.DyInc >= 300000 {
		ids = append(ids, 7)
	}
	return ids, rs
}

// PrioritySR .
func (hs *HonorStat) PrioritySR() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	if hs.Rank0 > 0 && hs.Rank0 <= 100 {
		ids = append(ids, 59)
	}
	if hs.Play > 100000 && hs.PlayLastW < 100000 {
		ids = append(ids, 9)
	}
	if hs.Fans > 100000 && hs.FansLastW < 100000 {
		ids = append(ids, 10)
	}
	if hs.PlayInc >= 1000000 {
		ids = append(ids, 11)
		stars.Play = 5
	}
	if hs.FansInc >= 10000 {
		ids = append(ids, 12)
		stars.Fans = 5
	}
	if hs.LikeInc >= 5000 {
		ids = append(ids, 13)
		stars.Like = 5
	}
	if hs.ShInc >= 1000 {
		ids = append(ids, 14)
		stars.Share = 5
	}
	if hs.CoinInc >= 3000 {
		ids = append(ids, 15)
		stars.Coin = 5
	}
	if hs.DyInc >= 50000 && hs.DyInc < 299999 {
		ids = append(ids, 16)
	}
	return ids, stars
}

// PriorityR .
func (hs *HonorStat) PriorityR() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	on, _, _ := hs.PartionRank()
	if on {
		ids = append(ids, 8)
	}
	if hs.Play > 10000 && hs.PlayLastW < 10000 {
		ids = append(ids, 17)
	}
	if hs.Fans > 10000 && hs.FansLastW < 10000 {
		ids = append(ids, 18)
	}
	if hs.PlayInc >= 100000 && hs.PlayInc < 1000000 {
		ids = append(ids, 19)
		stars.Play = 4
	}
	if hs.FansInc >= 1000 && hs.FansInc < 10000 {
		ids = append(ids, 20)
		stars.Fans = 4
	}
	if hs.LikeInc >= 1000 && hs.LikeInc < 5000 {
		ids = append(ids, 21)
		stars.Like = 4
	}
	if hs.ShInc >= 100 && hs.ShInc < 1000 {
		ids = append(ids, 22)
		stars.Share = 4
	}
	if hs.CoinInc >= 1000 && hs.CoinInc < 3000 {
		ids = append(ids, 24)
		stars.Coin = 4
	}
	if hs.DyInc >= 5000 && hs.DyInc < 49999 {
		ids = append(ids, 25)
	}
	return ids, stars
}

// PriorityA .
func (hs *HonorStat) PriorityA() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	if hs.AvsInc >= 5 {
		ids = append(ids, 23)
	}
	if hs.Play > 1000 && hs.PlayLastW < 1000 {
		ids = append(ids, 26)
	}
	if hs.Fans >= 1000 && hs.FansLastW < 1000 {
		ids = append(ids, 27)
	}
	if hs.PlayInc >= 10000 && hs.PlayInc < 100000 {
		ids = append(ids, 28)
		stars.Play = 3
	}
	if hs.FansInc >= 100 && hs.FansInc < 1000 {
		ids = append(ids, 29)
		stars.Fans = 3
	}
	if hs.LikeInc >= 500 && hs.LikeInc < 1000 {
		ids = append(ids, 30)
		stars.Like = 3
	}
	if hs.ShInc >= 50 && hs.ShInc < 100 {
		ids = append(ids, 31)
		stars.Share = 3
	}
	if hs.CoinInc >= 100 && hs.CoinInc < 1000 {
		ids = append(ids, 33)
		stars.Coin = 3
	}
	if hs.DyInc >= 500 && hs.DyInc < 4999 {
		ids = append(ids, 34)
	}
	return ids, stars
}

// PriorityB .
func (hs *HonorStat) PriorityB() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	if hs.AvsInc >= 3 && hs.AvsInc <= 5 {
		ids = append(ids, 32)
	}
	if (hs.Fans > 500 && hs.FansLastW < 500) || (hs.Fans > 600 && hs.FansLastW < 600) || (hs.Fans > 700 && hs.FansLastW < 700) || (hs.Fans > 800 && hs.FansLastW < 800) || (hs.Fans > 900 && hs.FansLastW < 900) {
		ids = append(ids, 35)
	}
	if hs.PlayInc >= 1000 && hs.PlayInc < 10000 {
		ids = append(ids, 36)
		stars.Play = 2
	}
	if hs.FansInc >= 10 && hs.FansInc < 100 {
		ids = append(ids, 37)
		stars.Fans = 2
	}
	if hs.LikeInc >= 30 && hs.LikeInc < 500 {
		ids = append(ids, 38)
		stars.Like = 2
	}
	if hs.ShInc >= 10 && hs.ShInc < 50 {
		ids = append(ids, 39)
		stars.Share = 2
	}
	if hs.AvsInc >= 1 && hs.AvsInc <= 2 {
		ids = append(ids, 40)
	}
	if hs.CoinInc >= 30 && hs.CoinInc < 100 {
		ids = append(ids, 41)
		stars.Coin = 2
	}
	if hs.DyInc >= 100 && hs.DyInc < 499 {
		ids = append(ids, 42)
	}
	if hs.FansInc/10 > 2 && hs.FansLastW/10 < 2 {
		ids = append(ids, 43)
	}
	if hs.PlayInc >= 100 && hs.PlayInc < 1000 {
		ids = append(ids, 44)
		stars.Play = 2
	}
	if hs.FansInc >= 1 && hs.FansInc < 10 {
		ids = append(ids, 45)
		stars.Fans = 2
	}
	return ids, stars
}

// PriorityC .
func (hs *HonorStat) PriorityC() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	if hs.PlayInc > 0 && hs.PlayInc < 100 {
		ids = append(ids, 46)
		stars.Play = 1
	}
	if hs.LikeInc > 0 && hs.LikeInc < 30 {
		ids = append(ids, 47)
		stars.Like = 1
	}
	if hs.ShInc > 0 && hs.ShInc < 10 {
		ids = append(ids, 48)
		stars.Share = 1
	}
	if hs.CoinInc >= 1 && hs.CoinInc < 30 {
		ids = append(ids, 49)
		stars.Coin = 1
	}
	if hs.DyInc >= 1 && hs.DyInc <= 99 {
		ids = append(ids, 50)
	}
	return ids, stars
}

// PriorityD .
func (hs *HonorStat) PriorityD() ([]int, *RiseStage) {
	ids := make([]int, 0)
	stars := new(RiseStage)
	if hs.FansInc <= 0 {
		ids = append(ids, 51)
	}
	if hs.PlayInc == 0 {
		ids = append(ids, 52)
	}
	if hs.LikeInc == 0 {
		ids = append(ids, 53)
	}
	if hs.ShInc == 0 {
		ids = append(ids, 54)
	}
	if hs.AvsInc == 0 {
		ids = append(ids, 55)
	}
	if hs.CoinInc == 0 {
		ids = append(ids, 56)
	}
	if hs.DyInc == 0 {
		ids = append(ids, 57)
	}
	return ids, stars
}

// PartionRank .
func (hs *HonorStat) PartionRank() (bool, string, int32) {
	partions := make(map[int]string)
	partions[1] = "动画"
	partions[168] = "国创相关"
	partions[3] = "音乐"
	partions[129] = "舞蹈"
	partions[4] = "游戏"
	partions[36] = "科技"
	partions[160] = "生活"
	partions[119] = "鬼畜"
	partions[155] = "时尚"
	partions[5] = "娱乐"
	partions[181] = "影视"
	if hs.Rank1 != 0 && hs.Rank1 <= 100 {
		return true, partions[1], hs.Rank1
	}
	if hs.Rank168 != 0 && hs.Rank168 <= 100 {
		return true, partions[168], hs.Rank168
	}
	if hs.Rank3 != 0 && hs.Rank3 <= 100 {
		return true, partions[3], hs.Rank3
	}
	if hs.Rank129 != 0 && hs.Rank129 <= 100 {
		return true, partions[129], hs.Rank129
	}
	if hs.Rank4 != 0 && hs.Rank4 <= 100 {
		return true, partions[4], hs.Rank4
	}
	if hs.Rank36 != 0 && hs.Rank36 <= 100 {
		return true, partions[36], hs.Rank36
	}
	if hs.Rank160 != 0 && hs.Rank160 <= 100 {
		return true, partions[160], hs.Rank160
	}
	if hs.Rank119 != 0 && hs.Rank119 <= 100 {
		return true, partions[119], hs.Rank119
	}
	if hs.Rank155 != 0 && hs.Rank155 <= 100 {
		return true, partions[155], hs.Rank155
	}
	if hs.Rank5 != 0 && hs.Rank5 <= 100 {
		return true, partions[5], hs.Rank5
	}
	if hs.Rank181 != 0 && hs.Rank181 <= 100 {
		return true, partions[181], hs.Rank181
	}
	return false, "", 0
}

// LatestSunday when today=sunday,return today,else return last week's sunday
func LatestSunday() time.Time {
	now := time.Now()
	today, _ := time.ParseInLocation(layout, now.Format(layout), time.Local)
	sunday := today.AddDate(0, 0, int(-now.Weekday()))
	return sunday
}
