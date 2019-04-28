package netutil

import (
	"net"
	"sync"
	"sync/atomic"
)

var (
	// ErrLimitListener listen limit error.
	ErrLimitListener = &LimitListenerError{}
)

// LimitListenerError limit max connections of listener.
type LimitListenerError struct{}

// Temporary is the error temporary.
func (l *LimitListenerError) Temporary() bool { return true }

// Timeout is the error a timeout.
func (l *LimitListenerError) Timeout() bool { return true }

// Error return error message of error.
func (l *LimitListenerError) Error() string { return "LimitListener: limit" }

// LimitListener returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
func LimitListener(l net.Listener, n int32) net.Listener {
	return &limitListener{l, 0, n, make(chan struct{}, n)}
}

type limitListener struct {
	net.Listener
	cur int32
	max int32
	sem chan struct{}
}

func (l *limitListener) acquire() (ok bool) {
	ok = true
	if cur := atomic.AddInt32(&l.cur, 1); cur > l.max {
		select {
		case l.sem <- struct{}{}:
		default:
			ok = false
		}
	}
	return
}

func (l *limitListener) release() {
	if cur := atomic.AddInt32(&l.cur, -1); cur >= l.max {
		<-l.sem
	}
}

func (l *limitListener) Accept() (net.Conn, error) {
	ok := l.acquire()
	c, err := l.Listener.Accept()
	if err != nil {
		l.release()
		return nil, err
	}
	if !ok {
		l.release()
		c.Close()
		return nil, ErrLimitListener
	}
	return &limitListenerConn{Conn: c, release: l.release}, nil
}

type limitListenerConn struct {
	net.Conn
	once    sync.Once
	release func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.once.Do(l.release)
	return err
}
