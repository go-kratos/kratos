package websocket

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"go-common/app/service/main/broadcast/libs/bufio"
)

// Request request.
type Request struct {
	Method     string
	RequestURI string
	Proto      string
	Host       string
	Header     http.Header

	reader *bufio.Reader
}

// ReadRequest reads and parses an incoming request from b.
func ReadRequest(r *bufio.Reader) (req *Request, err error) {
	var (
		b  []byte
		ok bool
	)
	req = &Request{reader: r}
	if b, err = req.readLine(); err != nil {
		return
	}
	if req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(string(b)); !ok {
		return nil, fmt.Errorf("malformed HTTP request %s", b)
	}
	if req.Header, err = req.readMIMEHeader(); err != nil {
		return
	}
	req.Host = req.Header.Get("Host")
	return req, nil
}

func (r *Request) readLine() ([]byte, error) {
	var line []byte
	for {
		l, more, err := r.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		// Avoid the copy if the first call produced a full line.
		if line == nil && !more {
			return l, nil
		}
		line = append(line, l...)
		if !more {
			break
		}
	}
	return line, nil
}

func (r *Request) readMIMEHeader() (header http.Header, err error) {
	var (
		line []byte
		i    int
		k, v string
	)
	header = make(http.Header, 16)
	for {
		if line, err = r.readLine(); err != nil {
			return
		}
		line = trim(line)
		if len(line) == 0 {
			return
		}
		if i = bytes.IndexByte(line, ':'); i <= 0 {
			err = fmt.Errorf("malformed MIME header line: " + string(line))
			return
		}
		k = string(line[:i])
		// Skip initial spaces in value.
		i++ // skip colon
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		v = string(line[i:])
		header.Add(k, v)
	}
}

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

// trim returns s with leading and trailing spaces and tabs removed.
// It does not assume Unicode or UTF-8.
func trim(s []byte) []byte {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	n := len(s)
	for n > i && (s[n-1] == ' ' || s[n-1] == '\t') {
		n--
	}
	return s[i:n]
}
