package task

import (
	"sort"
	"sync"
	"time"

	"go-common/app/job/main/videoup-report/model/utils"
)

const (
	// UpperSuperWhite 优质>10w粉
	UpperSuperWhite = int8(1)
	// UpperSuperBlack 高危>10w粉
	UpperSuperBlack = int8(2)
	// UpperWhite 优质
	UpperWhite = int8(3)
	// UpperBigWhite 优质>1w粉
	UpperBigWhite = int8(4)
	// UpperBigNormal 普通>1w粉
	UpperBigNormal = int8(5)
	// UpperSuperNormal 普通>10w粉
	UpperSuperNormal = int8(6)
	// UpperBlack 高危
	UpperBlack = int8(7)
	// WConfMid 按照mid配置权重
	WConfMid = int8(0)
	// WConfTaskID 按照taskid配置权重
	WConfTaskID = int8(1)
	// WConfType 按照分区配置权重
	WConfType = int8(2)
	// WConfUpFrom 按照投稿来源配置权重
	WConfUpFrom = int8(3)

	//UpperTypeWhite 优质
	UpperTypeWhite int8 = 1
	//UpperTypeBlack 高危
	UpperTypeBlack int8 = 2
	//UpperTypePGC 生产组
	UpperTypePGC int8 = 3
	//UpperTypeUGCX don't know
	UpperTypeUGCX int8 = 3
	//UpperTypePolitices 时政
	UpperTypePolitices int8 = 5
	//UpperTypeEnterprise 企业
	UpperTypeEnterprise int8 = 7
	//UpperTypeSigned 签约
	UpperTypeSigned int8 = 15
)

var (
	// TaskCountTH 插队任务阈值
	TaskCountTH = 2000
	// SuperUpperTH 粉丝数阈值
	SuperUpperTH = int64(100000)
	// BigUpperTH 粉丝数阈值
	BigUpperTH = int64(10000)
	// WLVConf 各个权重等级具体的配置数值
	WLVConf = &WeightValueConf{
		MaxWeight:  int64(200000), //最大权重值
		MinWeight:  int64(-510),
		SubRelease: int64(18), //指派再释放的任务
		//特殊任务参数
		Slv1: int64(8),  // 普通用户>=1W粉
		Slv2: int64(10), // 普通用户>=10W粉
		Slv3: int64(12), // 优质用户<1W粉
		Slv4: int64(15), // 优质用户>=1W粉
		Slv5: int64(18), // 优质用户>=10W粉
		Slv6: int64(6),  // 高危用户>=10W粉
		Slv7: int64(0),  // 其他高危
		//普通任务参数
		Nlv1:   int64(3),
		Nlv2:   int64(6),
		Nlv3:   int64(9),
		Nlv4:   int64(12),
		Nlv5:   int64(0),
		Nsum9:  int64(0),  // 等待9分钟总和 3*0
		Nsum15: int64(6),  // 等待15分钟总和 2*3
		Nsum27: int64(30), // 等待27分钟总和 6 + 4*6
		Nsum45: int64(84), // 等待45分钟总和 30 + 6*9
		//定时任务参数
		Tlv1:   int64(3),
		Tlv2:   int64(9),
		Tlv3:   int64(21),
		Tlv4:   int64(0),
		Tsum2h: int64(120),
		Tsum1h: int64(300),
	}
)

//WeightValueConf 可配置的权重
type WeightValueConf struct {
	MaxWeight  int64 `json:"maxweight"`
	SubRelease int64 `json:"subrelease"`
	MinWeight  int64 `json:"minweight"`
	Slv1       int64 `json:"slv1"`
	Slv2       int64 `json:"slv2"`
	Slv3       int64 `json:"slv3"`
	Slv4       int64 `json:"slv4"`
	Slv5       int64 `json:"slv5"`
	Slv6       int64 `json:"slv6"`
	Slv7       int64 `json:"slv7"`
	Nlv1       int64 `json:"nlv1"`
	Nlv2       int64 `json:"nlv2"`
	Nlv3       int64 `json:"nlv3"`
	Nlv4       int64 `json:"nlv4"`
	Nlv5       int64 `json:"nlv5"`
	Nsum9      int64 `json:"-"`
	Nsum15     int64 `json:"-"`
	Nsum27     int64 `json:"-"`
	Nsum45     int64 `json:"-"`
	Tlv1       int64 `json:"tlv1"`
	Tlv2       int64 `json:"tlv2"`
	Tlv3       int64 `json:"tlv3"`
	Tlv4       int64 `json:"tlv4"`
	Tsum2h     int64 `json:"-"`
	Tsum1h     int64 `json:"-"`
}

//WeightConfig task_weight_config记录结构
type WeightConfig struct {
	ID       int64
	Mid      int64
	TaskID   int64
	Rule     int8
	Weight   int64
	Ctime    time.Time
	Mtime    time.Time
	UserName string
	Desc     string
}

//WeightParams 审核任务权重的相关参数
type WeightParams struct {
	TaskID    int64            `json:"taskid"`
	Weight    int64            `json:"weight"` //权重总值
	State     int8             `json:"state"`  //任务状态
	Mid       int64            `json:"mid"`
	Special   int8             `json:"special"` //特殊任务
	Ctime     utils.FormatTime `json:"ctime"`   //任务生成时间
	Ptime     utils.FormatTime `json:"ptime"`   //定时发布时间
	CfItems   []*ConfigItem    `json:"cfitems,omitempty"`
	Fans      int64            `json:"fans"`     //粉丝数
	AccFailed bool             `json:"accfaild"` //账号查询是否失败

	UpGroups []int8 `json:"ugs"`    //分组
	UpFrom   int8   `json:"upfrom"` //来源
	TypeID   int16  `json:"typeid"` //分区
}

// ConfigItem task weight config item
type ConfigItem struct {
	ID     int64            `json:"id"`
	Radio  int8             `json:"radio"`
	CID    int64            `json:"cid"` // config id 四种配置通用
	Uname  string           `json:"user,omitempty"`
	Rule   int8             `json:"rule"`
	Weight int64            `json:"weight"`
	Mtime  utils.FormatTime `json:"mtime"`
	Desc   string           `json:"desc,omitempty"`
	Bt     utils.FormatTime `json:"et"`
	Et     utils.FormatTime `json:"bt"`
}

//WeightLog 权重变更记录
type WeightLog struct {
	TaskID  int64            `json:"taskid"`
	Mid     int64            `json:"mid"`     //用户id
	Weight  int64            `json:"weight"`  //任务权重总和
	CWeight int64            `json:"cweight"` //配置权重
	NWeight int64            `json:"nweight"` //普通任务
	SWeight int64            `json:"sweight"` //特殊任务
	TWeight int64            `json:"tweight"` //定时任务
	Uptime  utils.FormatTime `json:"uptime"`  //更新时间
	CfItems []*ConfigItem    `json:"cfitems,omitempty"`
}

// JumpList 插队同步的任务
type JumpList struct {
	l     []*WeightLog
	min   int64
	count int
	mux   sync.RWMutex
}

// NewJumpList New JumpList
func NewJumpList() *JumpList {
	return &JumpList{
		l:     []*WeightLog{},
		min:   -1,
		count: 0,
	}
}

// PUSH 添加
func (jl *JumpList) PUSH(item *WeightLog) {
	jl.mux.Lock()
	defer jl.mux.Unlock()
	if jl.count == TaskCountTH { //队列满了
		if item.Weight > jl.min { //剔除最小的
			jl.l = jl.l[1:jl.count]
			jl.min = jl.l[0].Weight
			jl.count--
		} else {
			return
		}
	}
	inx := sort.SearchInts(jl.List(), int(item.Weight))
	switch {
	case inx == 0: //头部
		jl.l = append([]*WeightLog{item}, jl.l...)
		jl.min = item.Weight
	case inx == jl.count: //尾部
		jl.l = append(jl.l, item)
	default:
		rear := append([]*WeightLog{}, jl.l[inx:]...)
		jl.l = append(jl.l[:inx], item)
		jl.l = append(jl.l, rear...)
	}
	jl.count++
}

// POP 读取
func (jl *JumpList) POP() (item *WeightLog) {
	jl.mux.Lock()
	defer jl.mux.Unlock()
	if jl.count > 1 {
		item = jl.l[jl.count-1]
		jl.l = jl.l[:jl.count-1]
		jl.count--
		jl.min = jl.l[jl.count-1].Weight
		return
	} else if jl.count == 1 {
		item = jl.l[0]
		jl.l = []*WeightLog{}
		jl.count = 0
		jl.min = -1
		return
	}
	return nil
}

// Reset 重置
func (jl *JumpList) Reset() {
	jl.mux.Lock()
	defer jl.mux.Unlock()
	jl.l = []*WeightLog{}
	jl.min = -1
	jl.count = 0
}

// List 待更新权重的任务
func (jl *JumpList) List() []int {
	arr := []int{}
	for _, jw := range jl.l {
		arr = append(arr, int(jw.Weight))
	}
	return arr
}
