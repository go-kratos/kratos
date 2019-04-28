package dao

import (
	"context"
	"fmt"
	"math/rand"

	"go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/model"
	"go-common/library/sync/errgroup"
)

// Nodes is helper to manage lifecycle of a collection of Nodes.
type Nodes struct {
	nodes    []*Node
	zones    map[string][]*Node
	selfAddr string
}

// NewNodes new nodes and return.
func NewNodes(c *conf.Config) *Nodes {
	nodes := make([]*Node, 0, len(c.Nodes))
	for _, addr := range c.Nodes {
		n := newNode(c, addr)
		n.pRegisterURL = fmt.Sprintf("http://%s%s", c.BM.Inner.Addr, _registerURL)
		nodes = append(nodes, n)
	}
	zones := make(map[string][]*Node)
	for name, addrs := range c.Zones {
		var znodes []*Node
		for _, addr := range addrs {
			n := newNode(c, addr)
			n.otherZone = true
			n.zone = name
			n.pRegisterURL = fmt.Sprintf("http://%s%s", c.BM.Inner.Addr, _registerURL)
			znodes = append(znodes, n)
		}
		zones[name] = znodes
	}
	return &Nodes{
		nodes:    nodes,
		zones:    zones,
		selfAddr: c.BM.Inner.Addr,
	}
}

// Replicate replicate information to all nodes except for this node.
func (ns *Nodes) Replicate(c context.Context, action model.Action, i *model.Instance, otherZone bool) (err error) {
	if len(ns.nodes) == 0 {
		return
	}
	eg, c := errgroup.WithContext(c)
	for _, n := range ns.nodes {
		if !ns.Myself(n.addr) {
			ns.action(c, eg, action, n, i)
		}
	}
	if !otherZone {
		for _, zns := range ns.zones {
			if n := len(zns); n > 0 {
				ns.action(c, eg, action, zns[rand.Intn(n)], i)
			}
		}
	}
	err = eg.Wait()
	return
}

func (ns *Nodes) action(c context.Context, eg *errgroup.Group, action model.Action, n *Node, i *model.Instance) {
	switch action {
	case model.Register:
		eg.Go(func() error {
			n.Register(c, i)
			return nil
		})
	case model.Renew:
		eg.Go(func() error {
			n.Renew(c, i)
			return nil
		})
	case model.Cancel:
		eg.Go(func() error {
			n.Cancel(c, i)
			return nil
		})
	}
}

// ReplicateSet replicate set information to all nodes except for this node.
func (ns *Nodes) ReplicateSet(c context.Context, arg *model.ArgSet, otherZone bool) (err error) {
	if len(ns.nodes) == 0 {
		return
	}
	eg, c := errgroup.WithContext(c)
	for _, n := range ns.nodes {
		if !ns.Myself(n.addr) {
			eg.Go(func() error {
				return n.Set(c, arg)
			})
		}
	}
	if !otherZone {
		for _, zns := range ns.zones {
			if n := len(zns); n > 0 {
				node := zns[rand.Intn(n)]
				eg.Go(func() error {
					return node.Set(c, arg)
				})
			}
		}
	}
	err = eg.Wait()
	return
}

// Nodes returns nodes of local zone.
func (ns *Nodes) Nodes() (nsi []*model.Node) {
	nsi = make([]*model.Node, 0, len(ns.nodes))
	for _, nd := range ns.nodes {
		if nd.otherZone {
			continue
		}
		node := &model.Node{
			Addr:   nd.addr,
			Status: nd.status,
			Zone:   nd.zone,
		}
		nsi = append(nsi, node)
	}
	return
}

// AllNodes returns nodes contain other zone nodes.
func (ns *Nodes) AllNodes() (nsi []*model.Node) {
	nsi = make([]*model.Node, 0, len(ns.nodes))
	for _, nd := range ns.nodes {
		node := &model.Node{
			Addr:   nd.addr,
			Status: nd.status,
			Zone:   nd.zone,
		}
		nsi = append(nsi, node)
	}
	for _, zns := range ns.zones {
		if n := len(zns); n > 0 {
			nd := zns[rand.Intn(n)]
			node := &model.Node{
				Addr:   nd.addr,
				Status: nd.status,
				Zone:   nd.zone,
			}
			nsi = append(nsi, node)
		}
	}
	return
}

// Myself returns whether or not myself.
func (ns *Nodes) Myself(addr string) bool {
	return ns.selfAddr == addr
}

// UP marks status of myself node up.
func (ns *Nodes) UP() {
	for _, nd := range ns.nodes {
		if ns.Myself(nd.addr) {
			nd.status = model.NodeStatusUP
		}
	}
}
