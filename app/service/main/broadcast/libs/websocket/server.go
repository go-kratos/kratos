package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"go-common/app/service/main/broadcast/libs/bufio"
)

var (
	keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	// ErrBadRequestMethod bad request method
	ErrBadRequestMethod = errors.New("bad method")
	// ErrNotWebSocket not websocket protocal
	ErrNotWebSocket = errors.New("not websocket protocol")
	// ErrBadWebSocketVersion bad websocket version
	ErrBadWebSocketVersion = errors.New("missing or bad WebSocket Version")
	// ErrChallengeResponse mismatch challenge response
	ErrChallengeResponse = errors.New("mismatch challenge/response")
)

// Upgrade Switching Protocols
func Upgrade(rwc io.ReadWriteCloser, rr *bufio.Reader, wr *bufio.Writer, req *Request) (conn *Conn, err error) {
	if req.Method != "GET" {
		return nil, ErrBadRequestMethod
	}
	if req.Header.Get("Sec-Websocket-Version") != "13" {
		return nil, ErrBadWebSocketVersion
	}
	if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
		return nil, ErrNotWebSocket
	}
	if !strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") {
		return nil, ErrNotWebSocket
	}
	challengeKey := req.Header.Get("Sec-Websocket-Key")
	if challengeKey == "" {
		return nil, ErrChallengeResponse
	}
	wr.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n")
	wr.WriteString("Sec-WebSocket-Accept: " + computeAcceptKey(challengeKey) + "\r\n\r\n")
	if err = wr.Flush(); err != nil {
		return
	}
	return newConn(rwc, rr, wr), nil
}

func computeAcceptKey(challengeKey string) string {
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
