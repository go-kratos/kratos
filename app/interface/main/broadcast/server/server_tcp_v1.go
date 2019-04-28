package server

import (
	"io"
	"net"
	"time"

	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/libs/bytes"
	itime "go-common/app/service/main/broadcast/libs/time"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

// InitTCPV1 listen all tcp.bind and start accept connections.
func InitTCPV1(server *Server, addrs []string, accept int) (err error) {
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
			go acceptTCPV1(server, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptTCPV1(server *Server, lis *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			log.Error("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			time.Sleep(time.Second)
			continue
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
		go serveTCPV1(server, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

func serveTCPV1(s *Server, conn *net.TCPConn, r int) {
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
	s.serveTCPV1(conn, rp, wp, tr)
}

// TODO linger close?
func (s *Server) serveTCPV1(conn *net.TCPConn, rp, wp *bytes.Pool, tr *itime.Timer) {
	var (
		err    error
		roomID string
		hb     time.Duration // heartbeat
		p      *model.Proto
		b      *Bucket
		trd    *itime.TimerData
		rpt    *Report
		rb     = rp.Get()
		wb     = wp.Get()
		ch     = NewChannel(s.c.ProtoSection.CliProto, s.c.ProtoSection.SvrProto)
		rr     = &ch.Reader
		wr     = &ch.Writer
	)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	// handshake
	trd = tr.Add(time.Duration(s.c.ProtoSection.HandshakeTimeout), func() {
		conn.Close()
	})
	ch.V1 = true
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	// must not setadv, only used in auth
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, roomID, ch.Mid, hb, rpt, err = s.authTCPV1(rr, wr, p, ch.IP); err == nil {
			b = s.Bucket(ch.Key)
			err = b.Put(roomID, ch)
		}
	}
	if err != nil {
		if err != io.EOF {
			log.Error("key: %s ip: %s handshake failed error(%v)", ch.Key, conn.RemoteAddr().String(), err)
		}
		conn.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, hb)
	var online int32
	if ch.Room != nil {
		online = ch.Room.OnlineNum()
	}
	report(actionConnect, rpt, online)
	// hanshake ok start dispatch goroutine
	go s.dispatchTCPV1(ch.Key, conn, wr, wp, wb, ch)
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadTCPV1(rr); err != nil {
			break
		}
		if p.Operation == model.OpHeartbeat {
			tr.Set(trd, hb)
			p.Body = nil
			p.Operation = model.OpHeartbeatReply
		} else {
			if err = s.Operate(p, ch, b); err != nil {
				break
			}
		}
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF {
		log.Error("key: %s server tcp failed error(%v)", ch.Key, err)
	}
	b.Del(ch)
	tr.Del(trd)
	rp.Put(rb)
	conn.Close()
	ch.Close()
	//if err = s.Disconnect(context.Background(), ch.Mid, roomID); err != nil {
	//	log.Error("key: %s operator do disconnect error(%v)", ch.Key, err)
	//}
	if ch.Room != nil {
		online = ch.Room.OnlineNum()
	}
	report(actionDisconnect, rpt, online)
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *Server) dispatchTCPV1(key string, conn *net.TCPConn, wr *bufio.Writer, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
	var (
		err    error
		finish bool
		online int32
	)
	for {
		var p = ch.Ready()
		switch p {
		case model.ProtoFinish:
			finish = true
			goto failed
		case model.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					err = nil // must be empty error
					break
				}
				if p.Operation == model.OpHeartbeatReply {
					if ch.Room != nil {
						online = ch.Room.OnlineNum()
					}
					if err = p.WriteTCPHeartV1(wr, online); err != nil {
						goto failed
					}
				} else {
					if err = p.WriteTCPV1(wr); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			// server send
			if err = p.WriteTCPV1(wr); err != nil {
				goto failed
			}
		}
		// only hungry flush response
		if err = wr.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil {
		log.Error("key: %s dispatch tcp error(%v)", key, err)
	}
	conn.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't blocking Signal
	for !finish {
		finish = (ch.Ready() == model.ProtoFinish)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *Server) authTCPV1(rr *bufio.Reader, wr *bufio.Writer, p *model.Proto, ip string) (key, roomID string, userID int64, heartbeat time.Duration, rpt *Report, err error) {
	if err = p.ReadTCPV1(rr); err != nil {
		return
	}
	if p.Operation != model.OpAuth {
		log.Warn("auth operation not valid: %d", p.Operation)
		err = ErrOperation
		return
	}
	if userID, roomID, key, rpt, err = s.NoAuth(int16(p.Ver), p.Body, ip); err != nil {
		return
	}
	heartbeat = _clientHeartbeat
	p.Body = nil
	p.Operation = model.OpAuthReply
	if err = p.WriteTCPV1(wr); err != nil {
		return
	}
	err = wr.Flush()
	return
}
