package conf

import (
	"errors"
	"flag"
	"go-common/library/net/rpc/warden"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	eCode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log            *log.Config
	BM             *bm.ServerConfig
	Verify         *verify.Config
	Tracer         *trace.Config
	Redis          *redis.Config
	MemCache       *memcache.Config
	MySQL          *sql.Config
	ECode          *eCode.Config
	LiveRpc        map[string]*liverpc.ClientConfig
	HttpClient     *bm.ClientConfig
	SkyHorseGray   map[string]bool
	SkyHorseStatus bool
	RpcTimeout     map[string]int64
	Bvc            map[string]string
	ChunkSize      map[string]int64
	DummyUid       map[string]int64
	AccountRPC     *rpc.ClientConfig
	XuserClient    *warden.ClientConfig
	LiveGray       map[string]bool
	Warden         *warden.ClientConfig
	AppConf        map[string]string
}

// ErrLogStrut ...
// 自定义ErrLog结构
type ErrLogStrut struct {
	Code       int64
	Msg        string
	ErrDesc    string
	ErrType    string
	URLName    string
	RPCTimeout int64
	ChunkSize  int64
	ChunkNum   int64
	ErrorPtr   *error
}

// ChunkCallInfo rpc调用配置
type ChunkCallInfo struct {
	ParamsName string
	URLName    string
	ChunkSize  int64
	ChunkNum   int64
	RPCTimeout int64
}

const (
	// EmptyResultEn 返回结果集为空
	EmptyResultEn = "got_empty_result"
	// EmptyResult 返回结果集为空
	EmptyResult = "调用直播服务返回data为空"
	// GetStatusInfoByUfos 获取房间信息
	GetStatusInfoByUfos = "room/v1/Room/get_status_info_by_uids"
	// TargetsWithMedal 获取房间信息
	TargetsWithMedal = "fans_medal/v1/FansMedal/targetsWithMedal"
	// GetRoomID 获取房间信息
	GetRoomID = "room/v2/Room/room_id_by_uid_multi"
	// Record 获取房间信息
	Record = "live_data/v1/Record/get"
	// GetPkIdsByRoomIds 获取房间信息
	GetPkIdsByRoomIds = "av/v1/Pk/getPkIdsByRoomIds"
	// RoomPendent 获取房间信息
	RoomPendent = "room/v1/RoomPendant/getPendantByIds"
	// RoomNews 获取房间信息
	RoomNews = "/room_ex/v1/RoomNews/multiGet"
	// RelationGiftInfo 获取房间信息
	// RelationGiftInfo = "/relation/v1/BaseInfo/getGiftInfo"
	// AccountGRPC ...   主站grpc用户信息
	AccountGRPC = "Cards3"
	// LiveUserExpGRPC ...
	// 直播用户经验grpc
	LiveUserExpGRPC = "xuserExp"
)

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	client.Watch("app-interface.toml")
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}

// GetTimeout implementation
// 获取超时
func GetTimeout(k string, def int64) (timeout int64) {
	if t, ok := Conf.RpcTimeout[k]; ok {
		timeout = t
	} else {
		timeout = def
	}
	return
}

// GetDummyUidConf implementation
// 获取模拟配置开关
func GetDummyUidConf() (config int64) {
	if t, ok := Conf.DummyUid["enable"]; ok {
		config = t
	} else {
		config = 0
	}
	return
}

// GetChunkSize implementation
// 获取模GetChunkSize
func GetChunkSize(k string, def int64) (timeout int64) {
	if t, ok := Conf.ChunkSize[k]; ok {
		timeout = t
	} else {
		timeout = def
	}
	return
}

// CheckReturn ...
// 检查返回
func CheckReturn(err error, code int64, msg string, urlName string,
	rpcTimeout int64, chunkSize int64, chunkNum int64) (errLog *ErrLogStrut, success bool) {
	errInfo := ErrLogStrut{}
	errInfo.URLName = urlName
	errInfo.RPCTimeout = rpcTimeout
	errInfo.ChunkSize = chunkSize
	errInfo.ChunkNum = chunkNum
	success = true
	if err != nil {
		errInfo.Code = 1003000
		errInfo.Msg = ""
		errInfo.ErrDesc = "liveRpc调用失败"
		errInfo.ErrType = "LiveRpcFrameWorkCallError"
		errInfo.ErrorPtr = &err
		success = false
	} else if code != 0 {
		errInfo.Code = code
		errInfo.Msg = msg
		errInfo.ErrDesc = "调用直播服务" + urlName + "出错"
		errInfo.ErrType = "CallLiveRpcCodeError"
		success = false
	}
	errLog = &errInfo
	return
}
