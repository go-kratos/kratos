package notify

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/snluu/uuid"
	"go-common/app/infra/notify/model"
	"go-common/library/log"
	"go-common/library/net/rpc/liverpc"
	"os"
	"strconv"
	"sync"
)

var (
	errLiverpcInvalidParams = errors.New("liverpc callback without version or cmd params")
	errStrconvVersion       = errors.New("liverpc post try strconv version param failed")
	errCallRaw              = errors.New("liverpc callRaw failed")
	errEmptyClient          = errors.New("liverpc client empty")
	_liverpcCaller          = "notify-liverpc-client"
	_defaultGroup           = "default"
)

// LiverpcClients liverpc clients
type LiverpcClients struct {
	rwLock  sync.RWMutex
	group   string
	clients map[string]*liverpc.Client
}

// LiveMsgValue live databus message.value
type LiveMsgValue struct {
	Topic        string         `json:"topic"`
	MsgID        interface{}    `json:"msg_id"`
	MsgKey       interface{}    `json:"msg_key"`
	MsgContent   interface{}    `json:"msg_content"`
	Timestamp    float64        `json:"timestamp"`
	CallerHeader *LiveMsgHeader `json:"caller_header"`
}

// LiveMsgHeader live databus message.value.caller_header
type LiveMsgHeader struct {
	TraceId     string `json:"trace_id"`
	Caller      string `json:"caller"`
	SourceGroup string `json:"source_group"`
	Group       string `json:"group"`
}

// NewClients new clients
func newLiverpcClients(w *model.Watcher) *LiverpcClients {
	lrcs := &LiverpcClients{
		group:   getGroup(),
		clients: make(map[string]*liverpc.Client),
	}
	for _, cb := range w.Callbacks {
		if cb.URL.Schema == model.LiverpcSchema {
			lrcs.getClient(cb.URL)
		}
	}
	return lrcs
}

// Post do callback with liverpc client
func (lrcs *LiverpcClients) Post(ctx context.Context, notifyURL *model.NotifyURL, msg string) (err error) {
	var (
		reply  *liverpc.Reply
		v      int
		header *liverpc.Header
		body   = make(map[string]string)
		client *liverpc.Client
	)
	// parse params, alarm is must, need retry
	version := notifyURL.Query.Get("version")
	cmd := notifyURL.Query.Get("cmd")
	if version == "" || cmd == "" {
		err = errLiverpcInvalidParams
		log.Error("LiverpcClient.Post callback params error(%v), url(%v), msg(%s)", err, notifyURL, msg)
		return
	}
	v, err = strconv.Atoi(version)
	if err != nil {
		log.Error("LiverpcClient.Post callback params error(%v), strconv version error, url(%v), msg(%s)", err, notifyURL, msg)
		err = errStrconvVersion
		return
	}

	// parse live message
	// if parse error, live message format invalid, fast fail and no need retry, alarm is must
	header, body, err = lrcs.formatLiveMsg(msg)
	if err != nil {
		log.Error("LiverpcClient.Post format message error(%v), url(%v), msg(%s)", err, notifyURL, msg)
		return nil
	}

	// prepare client and args, if no client currently, should retry and alarm
	client = lrcs.getClient(notifyURL)
	if client == nil {
		log.Error("LiverpcClient.Post empty client, url(%v), msg(%s)", notifyURL, msg)
		err = errEmptyClient
		return
	}
	args := &liverpc.Args{
		Header: header,
		Body:   body,
	}

	// do call
	reply, err = client.CallRaw(ctx, v, cmd, args)
	if err != nil {
		log.Error("LiverpcClient.Post CallRaw error(%v), url(%v), msg(%s)", err, notifyURL, msg)
		err = errCallRaw
		return
	}
	log.Info("LiverpcClient.Post CallRaw reply(%v), url(%v), msg(%s)", reply, notifyURL, msg)
	return
}

// getClient get client by appID
func (lrcs *LiverpcClients) getClient(notifyURL *model.NotifyURL) *liverpc.Client {
	lrcs.rwLock.RLock()
	c, ok := lrcs.clients[notifyURL.Host]
	lrcs.rwLock.RUnlock()
	if !ok {
		c = lrcs.newClient(notifyURL)
		lrcs.setClient(notifyURL.Host, c)
	}
	return c
}

// setClient set client
func (lrcs *LiverpcClients) setClient(appID string, client *liverpc.Client) {
	lrcs.rwLock.Lock()
	defer lrcs.rwLock.Unlock()
	lrcs.clients[appID] = client
}

// newClient create new liverpc client
func (lrcs *LiverpcClients) newClient(notifyURL *model.NotifyURL) *liverpc.Client {
	conf := &liverpc.ClientConfig{
		AppID: notifyURL.Host,
	}
	// one can just appoint addr by params
	if notifyURL.Query.Get("addr") != "" {
		conf.Addr = notifyURL.Query.Get("addr")
	}
	return liverpc.NewClient(conf)
}

// formatLiveMsg format live databus callback params
func (lrcs *LiverpcClients) formatLiveMsg(msg string) (header *liverpc.Header, body map[string]string, err error) {
	var bizMsg, msgContent string
	header = &liverpc.Header{}
	value := new(LiveMsgValue)
	err = json.Unmarshal([]byte(msg), value)
	if err != nil {
		log.Error("LiverpcClient.formatLiveMsg unmarshal msg error(%v), msg(%s), value(%v)", err, msg, value)
		return
	}
	if value.MsgContent != nil {
		var m1, m2 []byte
		m1, err = json.Marshal(value.MsgContent)
		if err != nil {
			log.Error("LiverpcClient.formatLiveMsg unmarshal MsgContent error(%v), msg(%s), value(%v)", err, msg, value)
			return
		}
		bm := map[string]interface{}{
			"topic":     value.Topic,
			"msg_id":    value.MsgID,
			"msg_key":   value.MsgKey,
			"timestamp": value.Timestamp,
		}
		m2, err = json.Marshal(bm)
		if err != nil {
			log.Error("LiverpcClient.formatLiveMsg unmarshal bizMsg error(%v), msg(%s), value(%v)", err, msg, value)
			return
		}
		msgContent = string(m1)
		bizMsg = string(m2)

		// caller header
		callerHeader := value.CallerHeader
		if callerHeader != nil {
			header.TraceId = callerHeader.TraceId
			if callerHeader.SourceGroup != "" {
				header.SourceGroup = callerHeader.SourceGroup
			} else {
				header.SourceGroup = callerHeader.Group
			}
		}
	} else {
		msgContent = msg
	}

	// supplement liverpc header
	header.Caller = _liverpcCaller
	if header.TraceId == "" {
		header.TraceId = uuid.Rand().Hex()
	}
	if header.SourceGroup == "" {
		header.SourceGroup = lrcs.group
	}
	// body
	body = map[string]string{
		"msg":         bizMsg,
		"msg_content": msgContent,
	}
	return
}

// getGroup get message source group
func getGroup() (g string) {
	g = os.Getenv("group")
	if g == "" {
		g = _defaultGroup
	}
	return
}
