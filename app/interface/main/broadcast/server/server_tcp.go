package server

import (
	"context"
	"io"
	"net"
	"strings"
	"time"

	iModel "go-common/app/interface/main/broadcast/model"
	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/libs/bytes"
	itime "go-common/app/service/main/broadcast/libs/time"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// InitTCP listen all tcp.bind and start accept connections.
func InitTCP(server *Server, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			log.Error("net.ResolveTCPAddr(\"tcp\", \"%s\") error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			log.Error("net.ListenTCP(\"tcp\", \"%s\") error(%v)", bind, err)
			return
		}
		log.Info("start tcp listen: \"%s\"", bind)
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptTCP(server, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptTCP(server *Server, lis *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			log.Error("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(server.c.TCP.Keepalive); err != nil {
			log.Error("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(server.c.TCP.Rcvbuf); err != nil {
			log.Error("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(server.c.TCP.Sndbuf); err != nil {
			log.Error("conn.SetWriteBuffer() error(%v)", err)
			return
		}
		go serveTCP(server, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

func serveTCP(s *Server, conn *net.TCPConn, r int) {
	var (
		// timer
		tr = s.round.Timer(r)
		rp = s.round.Reader(r)
		wp = s.round.Writer(r)
		// ip addr
		lAddr = conn.LocalAddr().String()
		rAddr = conn.RemoteAddr().String()
	)
	if s.c.Broadcast.Debug {
		log.Info("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
	}
	s.ServeTCP(conn, rp, wp, tr)
}

// ServeTCP .
func (s *Server) ServeTCP(conn *net.TCPConn, rp, wp *bytes.Pool, tr *itime.Timer) {
	var (
		err     error
		rid     string
		accepts []int32
		white   bool
		p       *model.Proto
		b       *Bucket
		trd     *itime.TimerData
		lastHb  = time.Now()
		rb      = rp.Get()
		wb      = wp.Get()
		ch      = NewChannel(s.c.ProtoSection.CliProto, s.c.ProtoSection.SvrProto)
		rr      = &ch.Reader
		wr      = &ch.Writer
	)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	// handshake
	step := 0
	trd = tr.Add(time.Duration(s.c.ProtoSection.HandshakeTimeout), func() {
		conn.Close()
		log.Error("key: %s remoteIP: %s step: %d tcp handshake timeout", ch.Key, conn.RemoteAddr().String(), step)
	})
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	// must not setadv, only used in auth
	step = 1
	md := metadata.MD{
		metadata.RemoteIP: ch.IP,
	}
	ctx := metadata.NewContext(context.Background(), md)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Mid, ch.Key, rid, ch.Platform, accepts, err = s.authTCP(ctx, rr, wr, p); err == nil {
			ch.Watch(accepts...)
			b = s.Bucket(ch.Key)
			err = b.Put(rid, ch)
			if s.c.Broadcast.Debug {
				log.Info("tcp connnected key:%s mid:%d proto:%+v", ch.Key, ch.Mid, p)
			}
		}
	}
	step = 2
	if err != nil {
		conn.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		log.Error("key: %s handshake failed error(%v)", ch.Key, err)
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, _clientHeartbeat)
	white = whitelist.Contains(ch.Mid)
	if white {
		whitelist.Printf("key: %s[%s] auth\n", ch.Key, rid)
	}
	step = 3
	reportCh(actionConnect, ch)
	// hanshake ok start dispatch goroutine
	go s.dispatchTCP(conn, wr, wp, wb, ch)
	serverHeartbeat := s.RandServerHearbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if white {
			whitelist.Printf("key: %s start read proto\n", ch.Key)
		}
		if err = p.ReadTCP(rr); err != nil {
			break
		}
		if white {
			whitelist.Printf("key: %s read proto:%v\n", ch.Key, p)
		}
		if p.Operation == model.OpHeartbeat {
			tr.Set(trd, _clientHeartbeat)
			p.Body = nil
			p.Operation = model.OpHeartbeatReply
			// last server heartbeat
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				if err = s.Heartbeat(ctx, ch.Mid, ch.Key); err == nil {
					lastHb = now
				} else {
					err = nil
				}
			}
			if s.c.Broadcast.Debug {
				log.Info("tcp heartbeat receive key:%s, mid:%d", ch.Key, ch.Mid)
			}
			step++
		} else {
			if err = s.Operate(p, ch, b); err != nil {
				break
			}
		}
		if white {
			whitelist.Printf("key: %s process proto:%v\n", ch.Key, p)
		}
		ch.CliProto.SetAdv()
		ch.Signal()
		if white {
			whitelist.Printf("key: %s signal\n", ch.Key)
		}
	}
	if white {
		whitelist.Printf("key: %s server tcp error(%v)\n", ch.Key, err)
	}
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		log.Error("key: %s server tcp failed error(%v)", ch.Key, err)
	}
	b.Del(ch)
	tr.Del(trd)
	rp.Put(rb)
	conn.Close()
	ch.Close()
	if err = s.Disconnect(ctx, ch.Mid, ch.Key); err != nil {
		log.Error("key: %s operator do disconnect error(%v)", ch.Key, err)
	}
	if white {
		whitelist.Printf("key: %s disconnect error(%v)\n", ch.Key, err)
	}
	reportCh(actionDisconnect, ch)
	if s.c.Broadcast.Debug {
		log.Info("tcp disconnected key: %s mid:%d", ch.Key, ch.Mid)
	}
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *Server) dispatchTCP(conn *net.TCPConn, wr *bufio.Writer, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
	var (
		err    error
		finish bool
		online int32
		white  = whitelist.Contains(ch.Mid)
	)
	if s.c.Broadcast.Debug {
		log.Info("key: %s start dispatch tcp goroutine", ch.Key)
	}
	for {
		if white {
			whitelist.Printf("key: %s wait proto ready\n", ch.Key)
		}
		var p = ch.Ready()
		if white {
			whitelist.Printf("key: %s proto ready\n", ch.Key)
		}
		if s.c.Broadcast.Debug {
			log.Info("key:%s dispatch msg:%v", ch.Key, *p)
		}
		switch p {
		case model.ProtoFinish:
			if white {
				whitelist.Printf("key: %s receive proto finish\n", ch.Key)
			}
			if s.c.Broadcast.Debug {
				log.Info("key: %s wakeup exit dispatch goroutine", ch.Key)
			}
			finish = true
			goto failed
		case model.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					err = nil // must be empty error
					break
				}
				if white {
					whitelist.Printf("key: %s start write client proto%v\n", ch.Key, p)
				}
				if p.Operation == model.OpHeartbeatReply {
					if ch.Room != nil {
						online = ch.Room.OnlineNum()
						b := map[string]interface{}{"room": map[string]interface{}{"online": online, "room_id": ch.Room.ID}}
						p.Body = iModel.Message(b, nil)
					}
					if err = p.WriteTCPHeart(wr); err != nil {
						goto failed
					}
				} else {
					if err = p.WriteTCP(wr); err != nil {
						goto failed
					}
				}
				if white {
					whitelist.Printf("key: %s write client proto%v\n", ch.Key, p)
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			if white {
				whitelist.Printf("key: %s start write server proto%v\n", ch.Key, p)
			}
			// server send
			if err = p.WriteTCP(wr); err != nil {
				goto failed
			}
			if white {
				whitelist.Printf("key: %s write server proto%v\n", ch.Key, p)
			}
			if s.c.Broadcast.Debug {
				log.Info("tcp sent a message key:%s mid:%d proto:%+v", ch.Key, ch.Mid, p)
			}
		}
		if white {
			whitelist.Printf("key: %s start flush \n", ch.Key)
		}
		// only hungry flush response
		if err = wr.Flush(); err != nil {
			break
		}
		if white {
			whitelist.Printf("key: %s flush\n", ch.Key)
		}
	}
failed:
	if white {
		whitelist.Printf("key: %s dispatch tcp error(%v)\n", ch.Key, err)
	}
	if err != nil {
		log.Error("key: %s dispatch tcp error(%v)", ch.Key, err)
	}
	conn.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't blocking Signal
	for !finish {
		finish = (ch.Ready() == model.ProtoFinish)
	}
	if s.c.Broadcast.Debug {
		log.Info("key: %s dispatch goroutine exit", ch.Key)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *Server) authTCP(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *model.Proto) (mid int64, key string, rid string, platform string, accepts []int32, err error) {
	if err = p.ReadTCP(rr); err != nil {
		log.Error("authTCP.ReadTCP(key:%v).err(%v)", key, err)
		return
	}
	if p.Operation != model.OpAuth {
		log.Warn("auth operation not valid: %d", p.Operation)
		err = ErrOperation
		return
	}
	if mid, key, rid, platform, accepts, err = s.Connect(ctx, p, ""); err != nil {
		log.Error("authTCP.Connect(key:%v).err(%v)", key, err)
		return
	}
	p.Body = []byte(`{"code":0,"message":"ok"}`)
	p.Operation = model.OpAuthReply
	if err = p.WriteTCP(wr); err != nil {
		log.Error("authTCP.WriteTCP(key:%v).err(%v)", key, err)
		return
	}
	err = wr.Flush()
	return
}
