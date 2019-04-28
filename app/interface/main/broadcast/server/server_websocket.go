package server

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"strings"
	"time"

	iModel "go-common/app/interface/main/broadcast/model"
	"go-common/app/service/main/broadcast/libs/bytes"
	itime "go-common/app/service/main/broadcast/libs/time"
	"go-common/app/service/main/broadcast/libs/websocket"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// InitWebsocket listen all tcp.bind and start accept connections.
func InitWebsocket(server *Server, addrs []string, accept int) (err error) {
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
			go acceptWebsocket(server, listener)
		}
	}
	return
}

// InitWebsocketWithTLS init websocket with tls.
func InitWebsocketWithTLS(server *Server, addrs []string, certFile, privateFile string, accept int) (err error) {
	var (
		bind     string
		listener net.Listener
		cert     tls.Certificate
		certs    []tls.Certificate
	)
	certFiles := strings.Split(certFile, ",")
	privateFiles := strings.Split(privateFile, ",")
	for i := range certFiles {
		cert, err = tls.LoadX509KeyPair(certFiles[i], privateFiles[i])
		if err != nil {
			log.Error("Error loading certificate. error(%v)", err)
			return
		}
		certs = append(certs, cert)
	}
	tlsCfg := &tls.Config{Certificates: certs}
	tlsCfg.BuildNameToCertificate()
	for _, bind = range addrs {
		if listener, err = tls.Listen("tcp", bind, tlsCfg); err != nil {
			log.Error("net.ListenTCP(\"tcp\", \"%s\") error(%v)", bind, err)
			return
		}
		log.Info("start wss listen: \"%s\"", bind)
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptWebsocketWithTLS(server, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptWebsocket(server *Server, lis *net.TCPListener) {
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
		go serveWebsocket(server, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptWebsocketWithTLS(server *Server, lis net.Listener) {
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
		go serveWebsocket(server, conn, r)
		if r++; r == _maxInt {
			r = 0
		}
	}
}

func serveWebsocket(s *Server, conn net.Conn, r int) {
	var (
		// timer
		tr = s.round.Timer(r)
		rp = s.round.Reader(r)
		wp = s.round.Writer(r)
	)
	if s.c.Broadcast.Debug {
		// ip addr
		lAddr := conn.LocalAddr().String()
		rAddr := conn.RemoteAddr().String()
		log.Info("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
	}
	s.ServeWebsocket(conn, rp, wp, tr)
}

// ServeWebsocket .
func (s *Server) ServeWebsocket(conn net.Conn, rp, wp *bytes.Pool, tr *itime.Timer) {
	var (
		err     error
		accepts []int32
		rid     string
		white   bool
		p       *model.Proto
		b       *Bucket
		trd     *itime.TimerData
		lastHB  = time.Now()
		rb      = rp.Get()
		ch      = NewChannel(s.c.ProtoSection.CliProto, s.c.ProtoSection.SvrProto)
		rr      = &ch.Reader
		wr      = &ch.Writer
		ws      *websocket.Conn // websocket
		req     *websocket.Request
	)
	// reader
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	// handshake
	step := 0
	trd = tr.Add(time.Duration(s.c.ProtoSection.HandshakeTimeout), func() {
		conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
		conn.Close()
		log.Error("key: %s remoteIP: %s step: %d ws handshake timeout", ch.Key, conn.RemoteAddr().String(), step)
	})
	// websocket
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	step = 1
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
	step = 2
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
	// must not setadv, only used in auth
	step = 3
	md := metadata.MD{
		metadata.RemoteIP: ch.IP,
	}
	ctx := metadata.NewContext(context.Background(), md)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Mid, ch.Key, rid, ch.Platform, accepts, err = s.authWebsocket(ctx, ws, p, req.Header.Get("Cookie")); err == nil {
			ch.Watch(accepts...)
			b = s.Bucket(ch.Key)
			err = b.Put(rid, ch)
			if s.c.Broadcast.Debug {
				log.Info("websocket connnected key:%s mid:%d proto:%+v", ch.Key, ch.Mid, p)
			}
		}
	}
	step = 4
	if err != nil {
		ws.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		if err != io.EOF && err != websocket.ErrMessageClose {
			log.Error("key: %s remoteIP: %s step: %d ws handshake failed error(%v)", ch.Key, conn.RemoteAddr().String(), step, err)
		}
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, _clientHeartbeat)
	white = whitelist.Contains(ch.Mid)
	if white {
		whitelist.Printf("key: %s[%s] auth\n", ch.Key, rid)
	}
	// hanshake ok start dispatch goroutine
	step = 5
	reportCh(actionConnect, ch)
	go s.dispatchWebsocket(ws, wp, wb, ch)
	serverHeartbeat := s.RandServerHearbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if white {
			whitelist.Printf("key: %s start read proto\n", ch.Key)
		}
		if err = p.ReadWebsocket(ws); err != nil {
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
			if now := time.Now(); now.Sub(lastHB) > serverHeartbeat {
				if err = s.Heartbeat(ctx, ch.Mid, ch.Key); err == nil {
					lastHB = now
				} else {
					err = nil
				}
			}
			if s.c.Broadcast.Debug {
				log.Info("websocket heartbeat receive key:%s, mid:%d", ch.Key, ch.Mid)
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
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose && !strings.Contains(err.Error(), "closed") {
		log.Error("key: %s server ws failed error(%v)", ch.Key, err)
	}
	b.Del(ch)
	tr.Del(trd)
	ws.Close()
	ch.Close()
	rp.Put(rb)
	if err = s.Disconnect(ctx, ch.Mid, ch.Key); err != nil {
		log.Error("key: %s operator do disconnect error(%v)", ch.Key, err)
	}
	if white {
		whitelist.Printf("key: %s disconnect error(%v)\n", ch.Key, err)
	}
	reportCh(actionDisconnect, ch)
	if s.c.Broadcast.Debug {
		log.Info("websocket disconnected key: %s mid:%d", ch.Key, ch.Mid)
	}
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *Server) dispatchWebsocket(ws *websocket.Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
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
			log.Info("key:%s dispatch msg:%s", ch.Key, p.Body)
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
					if err = p.WriteWebsocketHeart(ws); err != nil {
						goto failed
					}
				} else {
					if err = p.WriteWebsocket(ws); err != nil {
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
			if err = p.WriteWebsocket(ws); err != nil {
				goto failed
			}
			if white {
				whitelist.Printf("key: %s write server proto%v\n", ch.Key, p)
			}
			if s.c.Broadcast.Debug {
				log.Info("websocket sent a message key:%s mid:%d proto:%+v", ch.Key, ch.Mid, p)
			}
		}
		if white {
			whitelist.Printf("key: %s start flush \n", ch.Key)
		}
		// only hungry flush response
		if err = ws.Flush(); err != nil {
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
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		log.Error("key: %s dispatch ws error(%v)", ch.Key, err)
	}
	ws.Close()
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
func (s *Server) authWebsocket(ctx context.Context, ws *websocket.Conn, p *model.Proto, cookie string) (mid int64, key string, rid string, platform string, accepts []int32, err error) {
	for {
		if err = p.ReadWebsocket(ws); err != nil {
			return
		}
		if p.Operation == model.OpAuth {
			break
		} else {
			log.Error("ws request operation(%d) not auth", p.Operation)
		}
	}
	if mid, key, rid, platform, accepts, err = s.Connect(ctx, p, cookie); err != nil {
		return
	}
	p.Body = []byte(`{"code":0,"message":"ok"}`)
	p.Operation = model.OpAuthReply
	if err = p.WriteWebsocket(ws); err != nil {
		return
	}
	err = ws.Flush()
	return
}
