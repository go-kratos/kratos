package redis

import (
	"context"
	"errors"
)

type Pipeliner interface {
	// Send writes the command to the client's output buffer.
	Send(commandName string, args ...interface{})

	// Exec executes all commands and get replies.
	Exec(ctx context.Context) (rs *Replies, err error)
}

var (
	ErrNoReply = errors.New("redis: no reply in result set")
)

type pipeliner struct {
	pool *Pool
	cmds []*cmd
}

type Replies struct {
	replies []*reply
}

type reply struct {
	reply interface{}
	err   error
}

func (rs *Replies) Next() bool {
	return len(rs.replies) > 0
}

func (rs *Replies) Scan() (reply interface{}, err error) {
	if !rs.Next() {
		return nil, ErrNoReply
	}
	reply, err = rs.replies[0].reply, rs.replies[0].err
	rs.replies = rs.replies[1:]
	return
}

type cmd struct {
	commandName string
	args        []interface{}
}

func (p *pipeliner) Send(commandName string, args ...interface{}) {
	p.cmds = append(p.cmds, &cmd{commandName: commandName, args: args})
	return
}

func (p *pipeliner) Exec(ctx context.Context) (rs *Replies, err error) {
	n := len(p.cmds)
	if n == 0 {
		return &Replies{}, nil
	}
	c := p.pool.Get(ctx)
	defer c.Close()
	for len(p.cmds) > 0 {
		cmd := p.cmds[0]
		p.cmds = p.cmds[1:]
		if err := c.Send(cmd.commandName, cmd.args...); err != nil {
			p.cmds = p.cmds[:0]
			return nil, err
		}
	}
	if err = c.Flush(); err != nil {
		p.cmds = p.cmds[:0]
		return nil, err
	}
	rps := make([]*reply, 0, n)
	for i := 0; i < n; i++ {
		rp, err := c.Receive()
		rps = append(rps, &reply{reply: rp, err: err})
	}
	rs = &Replies{
		replies: rps,
	}
	return
}
