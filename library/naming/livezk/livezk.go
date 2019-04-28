package livezk

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strings"
	"time"

	"go-common/library/log"
	"go-common/library/naming"
	xtime "go-common/library/time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	basePath = "/live/service"
	scheme   = "grpc"
)

// Zookeeper Server&Client settings.
type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout xtime.Duration
}

// New new live zookeeper registry
func New(config *Zookeeper) (naming.Registry, error) {
	lz := &livezk{
		zkConfig: config,
	}
	var err error
	lz.zkConn, lz.zkEvent, err = zk.Connect(config.Addrs, time.Duration(config.Timeout))
	if err != nil {
		go lz.eventproc()
	}
	return lz, err
}

type zkIns struct {
	Group       string `json:"group"`
	LibVersion  string `json:"lib_version"`
	StartupTime string `json:"startup_time"`
}

func newZkInsData(ins *naming.Instance) ([]byte, error) {
	zi := &zkIns{
		// TODO group support
		Group:       "default",
		LibVersion:  ins.Version,
		StartupTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(zi)
}

// livezk live service zookeeper registry
type livezk struct {
	zkConfig *Zookeeper
	zkConn   *zk.Conn
	zkEvent  <-chan zk.Event
}

var _ naming.Registry = &livezk{}

func (l *livezk) Register(ctx context.Context, ins *naming.Instance) (cancel context.CancelFunc, err error) {
	nodePath := path.Join(l.zkConfig.Root, basePath, ins.AppID)
	if err = l.createAll(nodePath); err != nil {
		return
	}
	var rpc string
	for _, addr := range ins.Addrs {
		u, ue := url.Parse(addr)
		if ue == nil && u.Scheme == scheme {
			rpc = u.Host
			break
		}
	}
	if rpc == "" {
		err = errors.New("no GRPC addr")
		return
	}

	dataPath := path.Join(nodePath, rpc)
	data, err := newZkInsData(ins)
	if err != nil {
		return nil, err
	}
	_, err = l.zkConn.Create(dataPath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return nil, err
	}
	return func() {
		l.unregister(dataPath)
	}, nil
}

func (l *livezk) Close() error {
	l.zkConn.Close()
	return nil
}

func (l *livezk) createAll(nodePath string) (err error) {
	seps := strings.Split(nodePath, "/")
	lastPath := "/"
	ok := false
	for _, part := range seps {
		if part == "" {
			continue
		}
		lastPath = path.Join(lastPath, part)
		if ok, _, err = l.zkConn.Exists(lastPath); err != nil {
			return err
		} else if ok {
			continue
		}
		if _, err = l.zkConn.Create(lastPath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return
		}
	}
	return
}

func (l *livezk) eventproc() {
	for event := range l.zkEvent {
		// TODO handle zookeeper event
		log.Info("zk event: err: %s, path: %s, server: %s, state: %s, type: %s",
			event.Err, event.Path, event.Server, event.State, event.Type)
	}
}

func (l *livezk) unregister(dataPath string) error {
	return l.zkConn.Delete(dataPath, -1)
}
