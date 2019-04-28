package election

import (
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"go-common/library/log"
	"path"
	"sort"
	"time"
)

var (
	//ErrNotInit init fail
	ErrNotInit = errors.New("init first")
	//ErrFailConn fail to connect
	ErrFailConn = errors.New("fail to connect zk")
)

//ZkElection election by zk
type ZkElection struct {
	dir        string
	servers    []string
	conn       *zk.Conn
	timeout    time.Duration
	RootPath   string
	NodePath   string
	MasterPath string
	IsMaster   bool
	// wait for this channel, true means master, false means follower
	C          <-chan bool
	leaderChan chan bool
	running    bool
}

// New create new ZkElection
func New(servers []string, dir string, timeout time.Duration) *ZkElection {
	var z = &ZkElection{}
	z.servers = servers
	z.dir = dir
	z.timeout = timeout
	z.running = true
	return z
}

//Init init the elections
//	dir is root path for election, if dir = "/project", then, election will use "/project/election" as election path
//		and node would be "/project/election/n_xxxxxxxx"
//	if error happens, the election would work
func (z *ZkElection) Init() (err error) {

	z.conn, _, err = zk.Connect(z.servers, z.timeout)
	if err != nil {
		log.Error("fail connect zk, err=%s", err)
		return
	}
	z.RootPath = path.Join(z.dir, "election")
	exist, _, err := z.conn.Exists(z.RootPath)
	if err != nil {
		log.Error("fail to check path, path=%s, err=%v", z.RootPath, err)
		return
	}
	if !exist {
		_, err = z.conn.Create(z.RootPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Error("create election fail, path=%s, err=%v", z.RootPath, err)
			return
		}
	}

	var pathPrefix = path.Join(z.RootPath, "/n_")
	z.NodePath, err = z.conn.Create(pathPrefix, []byte(""), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Error("fail create node path, path=%s, err=%s", pathPrefix, err)
		return
	}
	return
}

//Elect elect for leader
// wait for the chan, if get a true, mean you are the leader
// if you already a leader, a false means you are kicked
func (z *ZkElection) Elect() (err error) {
	if z.conn == nil {
		return ErrNotInit
	}

	if z.leaderChan != nil {
		return
	}

	z.leaderChan = make(chan bool)
	z.C = z.leaderChan
	go func() {
		defer close(z.leaderChan)
		var ch, err = z.watchChange()
		if err != nil {
			log.Error("fail to watch, err=%s", err)
			return
		}

		for z.running {
			event := <-ch
			switch event.Type {
			case zk.EventNodeDeleted:
				ch, err = z.watchChange()
				if err != nil {
					log.Error("fail to watch, try next time in seconds, err=%s", err)
					time.Sleep(time.Second * 5)
				}
			case zk.EventNotWatching:
				log.Warn("receive not watching event, event=%+v", event)
				if event.State == zk.StateDisconnected {
					log.Info("reinit zk")
					if err = z.Init(); err != nil {
						log.Error("err init zk again, sleep 5s, err=%v", err)
						time.Sleep(5*time.Second)
					} else {
						ch, err = z.watchChange()
						if err != nil {
							log.Error("fail to watch, err=%s", err)
							return
						}
					}
				}
			}
			log.Info("zk event, event=%+v", event)
		}
		log.Error("exit election proc")

	}()
	return
}

func (z *ZkElection) watchChange() (ch <-chan zk.Event, err error) {

	children, _, e := z.conn.Children(z.RootPath)
	if err != nil {
		err = e
		log.Error("get children error, err=%+v\n", err)
		return
	}
	if len(children) == 0 {
		log.Warn("no child get from root: %s", z.RootPath)
		return
	}
	var min string
	if len(children) > 0 {
		sort.SliceStable(children, func(i, j int) bool {
			return children[i] < children[j]
		})
		min = children[0]
		z.MasterPath = min
	}
	var nodebase = path.Base(z.NodePath)
	for i, v := range children {
		if v == nodebase {
			if min == nodebase {
				log.Info("this is master, node=%s", z.NodePath)
				z.IsMaster = true
				z.leaderChan <- true
				_, _, ch, err = z.conn.GetW(z.NodePath)
			} else {
				log.Info("master is %s", min)
				prev := children[i-1]
				var preNode = path.Join(z.RootPath, prev)
				_, _, ch, err = z.conn.GetW(preNode)
				if err != nil {
					log.Error("watch node fail, node=%s, err=%s", preNode, err)
				}
				z.IsMaster = false
				z.leaderChan <- false
			}
		} else {
			log.Warn("v=%s, not same with base, base=%s", v, nodebase)
		}
	}
	log.Info("watchChange, len(children)=%d", len(children))
	return
}

//Close close the election
func (z *ZkElection) Close() {
	z.running = false
}
