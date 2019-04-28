package server

import (
	"crypto/tls"
	"io"
	"net"
	"time"

	"go-common/app/service/main/broadcast/libs/bytes"
	itime "go-common/app/service/main/broadcast/libs/time"
	"go-common/app/service/main/broadcast/libs/websocket"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

// InitWebsocketV1 listen all tcp.bind and start accept connections.
func InitWebsocketV1(s *Server, addrs []string, accept int) (err error) {
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
		log.Info("start ws listen: \"%s\"", bind)
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptWebsocketV1(s, listener)
		}
	}
	return
}

// InitWebsocketWithTLSV1 .
func InitWebsocketWithTLSV1(s *Server, addrs []string, certFile, privateFile string, accept int) (err error) {
	var (
		bind     string
		listener net.Listener
		cert     tls.Certificate
	)
	cert, err = tls.LoadX509KeyPair(certFile, privateFile)
	if err != nil {
		log.Error("Error loading certificate. error(%v)", err)
		return
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	for _, bind = range addrs {
		if listener, err = tls.Listen("tcp", bind, tlsCfg); err != nil {
			log.Error("net.ListenTCP(\"tcp\", \"%s\") error(%v)", bind, err)
			return
		}
		log.Info("start wss listen: \"%s\"", bind)
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptWebsocketWithTLSV1(s, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptWebsocketV1(s *Server, lis *net.TCPListener) {
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
		if err = conn.SetKeepAlive(s.c.TCP.Keepalive); err != nil {
			log.Error("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(s.c.TCP.Rcvbuf); err != nil {
			log.Error("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(s.c.TCP.Sndbuf); err != nil {
			log.Error("conn.SetWriteBuffer() error(%v)", err)
			return
		}
		go serveWebsocketV1(s, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptWebsocketWithTLSV1(server *Server, lis net.Listener) {
	var (
		conn net.Conn
		err  error
		r    int
	)
	for {
		if conn, err = lis.Accept(); err != nil {
			// if listener close then return
			log.Error("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		go serveWebsocketV1(server, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

func serveWebsocketV1(server *Server, conn net.Conn, r int) {
	var (
		// timer
		tr = server.round.Timer(r)
		rp = server.round.Reader(r)
		wp = server.round.Writer(r)
	)
	server.serveWebsocketV1(conn, rp, wp, tr)
}

// TODO linger close?
func (s *Server) serveWebsocketV1(conn net.Conn, rp, wp *bytes.Pool, tr *itime.Timer) {
	var (
		err    error
		roomID string
		hb     time.Duration // heartbeat
		p      *model.Proto
		b      *Bucket
		trd    *itime.TimerData
		rb     = rp.Get()
		ch     = NewChannel(s.c.ProtoSection.CliProto, s.c.ProtoSection.SvrProto)
		rr     = &ch.Reader
		wr     = &ch.Writer
		ws     *websocket.Conn // websocket
		req    *websocket.Request
		rpt    *Report
	)
	// reader
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	// handshake
	trd = tr.Add(time.Duration(s.c.ProtoSection.HandshakeTimeout), func() {
		conn.SetDeadline(time.Now().Add(time.Millisecond))
		conn.Close()
	})
	// websocket
	if req, err = websocket.ReadRequest(rr); err != nil || req.RequestURI != "/sub" {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		if err != io.EOF {
			log.Error("http.ReadRequest(rr) error(%v)", err)
		}
		return
	}
	// writer
	wb := wp.Get()
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	if ws, err = websocket.Upgrade(conn, rr, wr, req); err != nil {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		wp.Put(wb)
		if err != io.EOF {
			log.Error("websocket.NewServerConn error(%v)", err)
		}
		return
	}
	ch.V1 = true
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	// must not setadv, only used in auth
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, roomID, ch.Mid, hb, rpt, err = s.authWebsocketV1(ws, p, ch.IP); err == nil {
			b = s.Bucket(ch.Key)
			err = b.Put(roomID, ch)
		}
	}
	if err != nil {
		if err != io.EOF && err != websocket.ErrMessageClose {
			log.Error("key: %s ip: %s handshake failed error(%v)", ch.Key, conn.RemoteAddr().String(), err)
		}
		ws.Close()
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
	go s.dispatchWebsocketV1(ch.Key, ws, wp, wb, ch)
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadWebsocketV1(ws); err != nil {
			break
		}
		if p.Operation == model.OpHeartbeat {
			tr.Set(trd, hb)
			p.Operation = model.OpHeartbeatReply
		} else {
			if err = s.Operate(p, ch, b); err != nil {
				break
			}
		}
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		log.Error("key: %s server tcp failed error(%v)", ch.Key, err)
	}
	b.Del(ch)
	tr.Del(trd)
	ws.Close()
	ch.Close()
	rp.Put(rb)
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
func (s *Server) dispatchWebsocketV1(key string, ws *websocket.Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
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
					if err = p.WriteWebsocketHeartV1(ws, online); err != nil {
						goto failed
					}
				} else {
					if err = p.WriteWebsocketV1(ws); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			// server send
			if err = p.WriteWebsocketV1(ws); err != nil {
				goto failed
			}
		}
		// only hungry flush response
		if err = ws.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		log.Error("key: %s dispatch tcp error(%v)", key, err)
	}
	ws.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't blocking Signal
	for !finish {
		finish = (ch.Ready() == model.ProtoFinish)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *Server) authWebsocketV1(ws *websocket.Conn, p *model.Proto, ip string) (key, roomID string, userID int64, heartbeat time.Duration, rpt *Report, err error) {
	if err = p.ReadWebsocketV1(ws); err != nil {
		return
	}
	if p.Operation != model.OpAuth {
		err = ErrOperation
		return
	}
	if userID, roomID, key, rpt, err = s.NoAuth(int16(p.Ver), p.Body, ip); err != nil {
		return
	}
	heartbeat = _clientHeartbeat
	p.Body = nil
	p.Operation = model.OpAuthReply
	if err = p.WriteWebsocketV1(ws); err != nil {
		return
	}
	err = ws.Flush()
	return
}
