package server

import (
	"github.com/samuel/go-zookeeper/zk"
	"go-common/library/log"
	"strings"
	"time"
)

type ZkClient struct {
	conn     *zk.Conn
	address  []string
	timeout  time.Duration
	dialTime time.Time
	closed   bool
	stopper  chan struct{}
}

func NewZkClient(addrs []string, timeout time.Duration) (*ZkClient, error) {
	if timeout <= 0 {
		timeout = time.Second * 5
	}
	c := &ZkClient{
		address: addrs,
		timeout: timeout,
		stopper: make(chan struct{}),
	}
	if err := c.Reset(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ZkClient) GetTimeout() time.Duration {
	return c.timeout
}

func (c *ZkClient) RecursiveCreate(path string) error {
	if path == "" || path == "/" {
		return nil
	}
	if exists, _, err := c.conn.Exists(path); err != nil {
		return err
	} else if exists {
		return nil
	}
	if err := c.RecursiveCreate(path[0:strings.LastIndex(path, "/")]); err != nil {
		return err
	}
	_, err := c.conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (c *ZkClient) Reset() error {
	c.dialTime = time.Now()
	conn, events, err := zk.Connect(c.address, c.timeout)
	if err != nil {
		return err
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.conn = conn
	go func() {
		for ev := range events {
			if ev.Err == nil {
				log.V(2).Info("[ZooKeeper]Event Info:%+v", ev)
			} else {
				log.Error("[ZooKeeper]Event Error:%+v", ev)
			}
		}
	}()
	return nil
}

func (c *ZkClient) CreateEphemeralNode(path string, node string, data []byte) (string, error) {
	if err := c.RecursiveCreate(path); err != nil {
		return "", err
	}
	nodePath := strings.Join([]string{path, node}, "/")
	path, err := c.conn.Create(nodePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return path, err
}

func (c *ZkClient) CreatePersistNode(path string, node string, data []byte) (string, error) {
	if err := c.RecursiveCreate(path); err != nil {
		return "", err
	}
	nodePath := strings.Join([]string{path, node}, "/")
	path, err := c.conn.Create(nodePath, data, 0, zk.WorldACL(zk.PermAll))
	return path, err
}

func (c *ZkClient) Exists(path string, node string) (bool, int32, error) {
	fullPath := strings.Join([]string{path, node}, "/")
	exists, stat, err := c.conn.Exists(fullPath)
	return exists, stat.Version, err
}

func (c *ZkClient) SetNodeData(path string, node string, data []byte, version int32) error {
	var err error
	fullPath := strings.Join([]string{path, node}, "/")
	_, err = c.conn.Set(fullPath, data, version)
	return err
}

func (c *ZkClient) GetNodeData(path string, node string) ([]byte, int32, error) {
	var err error
	fullPath := strings.Join([]string{path, node}, "/")
	data, stat, err := c.conn.Get(fullPath)
	return data, stat.Version, err
}

func (c *ZkClient) DeleteNode(path string, node string) error {
	var err error
	fullPath := strings.Join([]string{path, node}, "/")
	exists, stat, err := c.conn.Exists(fullPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	err = c.conn.Delete(fullPath, stat.Version)
	return err
}

func (c *ZkClient) Close() {
	if c.closed {
		return
	}
	c.closed = true

	if c.conn != nil {
		c.conn.Close()
	}
	close(c.stopper)
}

func (c *ZkClient) GetChildren(node string) ([]string, error) {
	children, _, err := c.conn.Children(node)
	return children, err
}

func (c *ZkClient) GetChildrenWithData(node string) (map[string]string, error) {
	children, _, err := c.conn.Children(node)
	result := make(map[string]string)
	for _, child := range children {
		if data, _, e := c.conn.Get(strings.Join([]string{node, child}, "/")); e == nil {
			result[child] = string(data)
		} else {
			log.Error("[ZookeeperClient]GetChildrenWithData:get child:%s failed, err:%s", child, e.Error())
		}
	}
	return result, err
}

func (c *ZkClient) GetData(path string) (string, error) {
	data, _, err := c.conn.Get(path)
	return string(data), err
}

func (c *ZkClient) WatchChildren(path string) (map[string]struct{}, <-chan zk.Event, error) {
	if exists, _, err := c.conn.Exists(path); err != nil {
		return nil, nil, err
	} else if !exists {
		return nil, nil, zk.ErrNoNode
	}
	children, _, event, err := c.conn.ChildrenW(path)
	if err != nil {
		return nil, nil, err
	}
	result := make(map[string]struct{})
	for _, child := range children {
		result[child] = struct{}{}
	}
	return result, event, nil
}

func (c *ZkClient) WatchChildrenWithData(node string) (map[string]string, <-chan zk.Event, error) {
	if exists, _, err := c.conn.Exists(node); err != nil {
		return nil, nil, err
	} else if !exists {
		return nil, nil, zk.ErrNoNode
	}

	children, _, event, err := c.conn.ChildrenW(node)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]string)
	for _, child := range children {
		if data, _, e := c.conn.Get(strings.Join([]string{node, child}, "/")); e == nil {
			result[child] = string(data)
		} else {
			return nil, nil, e
		}
	}
	return result, event, nil
}

func (c *ZkClient) WatchData(path string) ([]byte, <-chan zk.Event, error) {
	data, _, event, err := c.conn.GetW(path)
	return data, event, err
}

func (c *ZkClient) ZooKeeperPath(args ...string) string {
	return strings.Join(args, "/")
}
