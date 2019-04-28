package common

const (
	BitWiseBVC = 1
	BitWiseKS  = 2
	BitWiseQN  = 4
	BitWiseTC  = 8
	BitWiseWS  = 16

	WSSrc  = 1
	TXYSrc = 16
	BVCSrc = 32
	JSSrc  = 68
	QNSrc  = 70

	WSName  = "ws"
	TXYName = "txy"
	BVCName = "bvc"
	JSName  = "js"
	QNName  = "qn"

	WsChinaName  = "网宿"
	TXYChinaName = "腾讯云"
	BVCChinaName = "视频云"
	JSChinaName  = "金山"
	QNChinaName  = "七牛"

	AREAIDCHICKEN    = 80
	ChickenAttention = 1000
)

// CdnMapSrc cdn&src映射
var CdnMapSrc = map[string]int8{
	WSName:  WSSrc,
	TXYName: TXYSrc,
	BVCName: BVCSrc,
	JSName:  JSSrc,
	QNName:  QNSrc,
}

// SrcMapBitwise 新老src直接的映射关系
var SrcMapBitwise = map[int8]int64{
	WSSrc:  BitWiseWS,
	TXYSrc: BitWiseTC,
	BVCSrc: BitWiseBVC,
	JSSrc:  BitWiseKS,
	QNSrc:  BitWiseQN,
}

// BitwiseMapSrc 新老src直接的映射关系
var BitwiseMapSrc = map[int64]int8{
	BitWiseWS:  WSSrc,
	BitWiseTC:  TXYSrc,
	BitWiseBVC: BVCSrc,
	BitWiseKS:  JSSrc,
	BitWiseQN:  QNSrc,
}

// ChinaNameMapBitwise
var ChinaNameMapBitwise = map[string]int64{
	WsChinaName:  BitWiseWS,
	TXYChinaName: BitWiseTC,
	BVCChinaName: BitWiseBVC,
	JSChinaName:  BitWiseKS,
	QNChinaName:  BitWiseQN,
}

// NameMapBitwise
var NameMapBitwise = map[int64]string{
	BitWiseWS:  WsChinaName,
	BitWiseTC:  TXYChinaName,
	BitWiseBVC: BVCChinaName,
	BitWiseKS:  JSChinaName,
	BitWiseQN:  QNChinaName,
}

// CdnBitwiseMap 位运算对应表
var CdnBitwiseMap = map[string]int64{
	BVCName: BitWiseBVC,
	JSName:  BitWiseKS,
	QNName:  BitWiseQN,
	TXYName: BitWiseTC,
	WSName:  BitWiseWS,
}

// BitwiseMapName 新src对应关系
var BitwiseMapName = map[int64]string{
	BitWiseWS:  WSName,
	BitWiseTC:  TXYName,
	BitWiseBVC: BVCName,
	BitWiseKS:  JSName,
	BitWiseQN:  QNName,
}

// NewLinkMap 新origin
var NewLinkMap = map[int64]map[string]string{
	BitWiseBVC: {
		"newLink":      "http://live-schedule.acgvideo.com/live-upbvc?up_rtmp=",
		"newLinkHttps": "https://live-schedule.acgvideo.com/live-upbvc?up_rtmp=",
	},
	BitWiseTC: {
		"newLink":      "http://tcdns.myqcloud.com:8086/bilibili_redirect?up_rtmp=",
		"newLinkHttps": "https://tcdns.myqcloud.com/bilibili_redirect?up_rtmp=",
	},
	BitWiseKS: {
		"newLink":      "http://cwwshdns.ksyun.com/a?up_rtmp=",
		"newLinkHttps": "https://cwwshdns.ksyun.com/a?up_rtmp=",
	},
	BitWiseWS: {
		"newLink":      "http://sdkbilibili.wscdns.com/bilibili?up_rtmp=",
		"newLinkHttps": "https://sdkbilibili.wscdns.com?up_rtmp=",
	},
	BitWiseQN: {
		"newLink":      "http://pili-ipswitch.qiniuapi.com/v1/bilibili/publish?up_rtmp=",
		"newLinkHttps": "https://pili-ipswitch.qiniuapi.com/v1/bilibili/publish?up_rtmp=",
	},
}

// 流线路名称
var LineName = map[int64]string{
	BitWiseBVC: "线路二",
	BitWiseTC:  "线路一",
	BitWiseKS:  "线路四",
	BitWiseWS:  "线路三",
	BitWiseQN:  "线路五",
	0:          "默认路线",
}

// CDNSalt cdn&salt map
var CDNSalt = map[string]string{
	"bvc": "bvc_1701101740",
	"js":  "js_1703271720",
	"qn":  "qn_1703271730",
	"txy": "txy_1610171720",
	"ws":  "ws_1608121700",
}
