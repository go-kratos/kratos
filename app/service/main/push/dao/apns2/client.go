// Package apns2 is a go Apple Push Notification Service (APNs) provider that
// allows you to send remote notifications to your iOS, tvOS, and OS X
// apps, using the new APNs HTTP/2 network protocol.
package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"go-common/library/log"
	"go-common/library/stat"
	"go-common/library/stat/prom"

	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
)

const (
	// HostDevelopment dev host.
	HostDevelopment = "https://api.development.push.apple.com"
	// HostProduction pro host.
	HostProduction = "https://api.push.apple.com"
	// StatusCodeSuccess success.
	StatusCodeSuccess = 200
	// StatusCodeBadReq bad req.
	StatusCodeBadReq = 400
	// StatusCodeCerErr There was an error with the certificate.
	StatusCodeCerErr = 403
	// StatusCodeMethodErr The request used a bad :method value. Only POST requests are supported.
	StatusCodeMethodErr = 405
	// StatusCodeNotForTopic The device token is not form the topic.
	StatusCodeNotForTopic = 400
	// StatusCodeNoActive The device token is no longer active for the topic.
	StatusCodeNoActive = 410
	// StatusCodePayloadTooLarge  The notification payload was too large.
	StatusCodePayloadTooLarge = 413
	// StatusCodeTooManyReq The server received too many requests for the same device token.
	StatusCodeTooManyReq = 429
	// StatusCodeServerErr Internal server error
	StatusCodeServerErr = 500
	// StatusCodeServerUnavailable The server is shutting down and unavailable.
	StatusCodeServerUnavailable = 503
)

// DefaultHost is a mutable var for testing purposes
var DefaultHost = HostDevelopment

// Client represents a connection with the APNs
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
	BoundID     string
	Stats       stat.Stat
}

// func init() {
// 	proxy.RegisterDialerType("http", func(*url.URL, proxy.Dialer) (proxy.Dialer, error) {
// 		return &net.Dialer{}, nil
// 	})
// }

// NewClient returns a new Client with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewClient(certificate tls.Certificate, timeout time.Duration) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	// transport := &http.Transport{
	// 	TLSClientConfig: tlsConfig,
	// 	Proxy: func(_ *http.Request) (*url.URL, error) {
	// 		return url.Parse("http://10.28.10.11:80")
	// 	},
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 		DualStack: true,
	// 	}).DialContext,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }
	return &Client{
		HTTPClient:  &http.Client{Transport: transport, Timeout: timeout},
		Certificate: certificate,
		Host:        DefaultHost,
		Stats:       prom.HTTPClient,
	}
}

// NewClientWithProxy returns a new Client with sock5 proxy.
func NewClientWithProxy(certificate tls.Certificate, timeout time.Duration, proxyAddr string) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	return &Client{
		HTTPClient:  &http.Client{Transport: proxyTransport(proxyAddr, tlsConfig, timeout), Timeout: timeout},
		Certificate: certificate,
		Host:        DefaultHost,
		Stats:       prom.HTTPClient,
	}
}

func proxyTransport(proxyAddr string, config *tls.Config, timeout time.Duration) *http2.Transport {
	return &http2.Transport{
		DialTLS: func(network, addr string, cfg *tls.Config) (nc net.Conn, err error) {
			dialer := &net.Dialer{Timeout: timeout / 2}
			var proxyDialer proxy.Dialer
			if proxyDialer, err = proxy.SOCKS5("tcp", proxyAddr, nil, dialer); err != nil {
				log.Error("proxy.SOCKS5(%s) error(%v)", proxyAddr, err)
				return nil, err
			}
			// u, _ := url.Parse("http://10.28.10.11:80")
			// proxyDialer, err = proxy.FromURL(u, dialer)
			var conn net.Conn
			if conn, err = proxyDialer.Dial(network, addr); err != nil {
				log.Error("proxyDialer.Dial(%s,%s) error(%v)", network, addr, err)
				if conn, err = dialer.Dial(network, addr); err != nil {
					log.Error("dialer.Dial(%s,%s) error(%v)", network, addr, err)
					return nil, err
				}
			}
			tlsConn := tls.Client(conn, cfg)
			if err = tlsConn.Handshake(); err != nil {
				log.Error("tlsConn.Handshake() error(%v)", err)
				return nil, err
			}
			if !cfg.InsecureSkipVerify {
				if err = tlsConn.VerifyHostname(cfg.ServerName); err != nil {
					log.Error("tlsConn.VerifyHostname(%s) error(%v)", cfg.ServerName, err)
					return nil, err
				}
			}
			state := tlsConn.ConnectionState()
			if state.NegotiatedProtocol != http2.NextProtoTLS {
				err = fmt.Errorf("http2: unexpected ALPN protocol(%s) expect(%s)", state.NegotiatedProtocol, http2.NextProtoTLS)
				return nil, err
			}
			if !state.NegotiatedProtocolIsMutual {
				err = errors.New("http2: could not negotiate protocol mutually")
				return nil, err
			}
			return tlsConn, nil
		},
		TLSClientConfig: config,
	}
}

// Development sets the Client to use the APNs development push endpoint.
func (c *Client) Development() *Client {
	c.Host = HostDevelopment
	return c
}

// Production sets the Client to use the APNs production push endpoint.
func (c *Client) Production() *Client {
	c.Host = HostProduction
	return c
}

// Push sends a Notification to the APNs gateway. If the underlying http.Client
// is not currently connected, this method will attempt to reconnect
// transparently before sending the notification.
func (c *Client) Push(deviceToken string, payload *Payload, overTime int64) (response *Response, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.Host, int64(time.Since(now)/time.Millisecond))
			log.Info("apns stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.Host, "failed")
			}
		}()
	}
	var (
		req   *http.Request
		res   *http.Response
		t     = time.NewTimer(c.HTTPClient.Timeout)
		errCh = make(chan error, 1)
		url   = fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)
	)
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(payload.Marshal())); err != nil {
		log.Error("http.NewRequest(%s) error(%v)", url, err)
		return
	}
	req.Header.Set("apns-topic", c.BoundID)
	req.Header.Set("apns-expiration", strconv.FormatInt(overTime, 10))
	req.Header.Set("apns-collapse-id", payload.TaskID)
	go func() {
		res, err = c.HTTPClient.Do(req)
		errCh <- err
	}()
	select {
	case <-t.C:
		err = errors.New("http.Do timeout")
		return
	case err = <-errCh:
		if err != nil {
			log.Error("c.HTTPClient.Do() error(%v)", err)
			return
		}
	}
	defer res.Body.Close()
	response = &Response{StatusCode: res.StatusCode, ApnsID: res.Header.Get("apns-id")}
	var bs []byte
	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return
	} else if len(bs) == 0 {
		return
	}
	if e := json.Unmarshal(bs, &response); e != nil {
		if e != io.EOF {
			log.Error("json decode body(%s) error(%v)", string(bs), e)
		}
	}
	return
}

// MockPush mock push.
func (c *Client) MockPush(deviceToken string, payload *Payload, overTime int64) (response *Response, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.Host, int64(time.Since(now)/time.Millisecond))
			// log.Info("mock apns stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.Host, "apple push mock")
			}
		}()
	}
	time.Sleep(200 * time.Millisecond)
	response = &Response{StatusCode: StatusCodeSuccess}
	return
}
