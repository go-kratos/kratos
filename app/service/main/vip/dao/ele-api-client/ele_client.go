package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	xhttp "net/http"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	pkgerr "github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	_contentTypeJSON = "application/json"
)

//EleClient client is http client, for third ele server.
type EleClient struct {
	client *bm.Client
	conf   *Config
}

// Config is http client conf.
type Config struct {
	*App
}

// App bilibili intranet authorization.
type App struct {
	Key    string
	Secret string
}

// NewEleClient new a http client.
func NewEleClient(c *Config, client *bm.Client) *EleClient {
	cl := new(EleClient)
	cl.conf = c
	// check appkey
	if c.Key == "" || c.Secret == "" {
		panic("http client must config appkey and appsecret")
	}
	cl.client = client
	return cl
}

// Get a json req http get.
func (cl *EleClient) Get(c context.Context, host, path string, args interface{}, res interface{}) (err error) {
	req, err := cl.newRequest(xhttp.MethodGet, host, path, args)
	if err != nil {
		return
	}
	return cl.client.Do(c, req, res)
}

// Post a json req http post.
func (cl *EleClient) Post(c context.Context, host, path string, args interface{}, res interface{}) (err error) {
	req, err := cl.newRequest(xhttp.MethodPost, host, path, args)
	if err != nil {
		return
	}
	return cl.client.Do(c, req, res)
}

// IsSuccess check ele api is success.
func IsSuccess(message string) bool {
	return message == "ok"
}

// newRequest new http request with  host, path, method, ip, values and headers, without sign.
func (cl *EleClient) newRequest(method, host, path string, args interface{}) (req *xhttp.Request, err error) {
	consumerKey := cl.conf.Key
	nonce := UUID4() //TODO uuid 有问题？
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := eleSign(consumerKey, nonce, timestamp, path, cl.conf.Secret)
	params := map[string]interface{}{}
	params["consumer_key"] = consumerKey
	params["nonce"] = nonce
	params["timestamp"] = timestamp
	params["sign"] = sign
	params["args"] = args
	url := host + path
	marshal, err := json.Marshal(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "marshal:%v", params)
		return
	}
	rj := string(marshal)
	log.Info("ele_client req method(%s) url(%s) rj(%s)", method, url, rj)
	req, err = xhttp.NewRequest(method, url, strings.NewReader(rj))
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s", url+" "+rj)
		return
	}
	req.Header.Set("Content-Type", _contentTypeJSON)
	return
}

func eleSign(consumerKey, nonce, timestamp, path, secret string) string {
	var b bytes.Buffer
	b.WriteString(path)
	b.WriteString("&")
	b.WriteString("consumer_key=")
	b.WriteString(consumerKey)
	b.WriteString("&nonce=")
	b.WriteString(nonce)
	b.WriteString("&timestamp=")
	b.WriteString(timestamp)
	return computeHmac256(b, secret)
}

func computeHmac256(b bytes.Buffer, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(b.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// UUID4 is generate uuid
func UUID4() string {
	return uuid.NewV4().String()
}
