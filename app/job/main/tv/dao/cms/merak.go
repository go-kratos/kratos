package cms

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// MerakNotify send notify
func (d *Dao) MerakNotify(ctx context.Context, title, content string) (err error) {
	var (
		cfg  = d.conf.Cfg.Merak
		sign string
		req  *http.Request
		body []byte
	)
	params := map[string]string{
		"Action":    "CreateWechatMessage",
		"PublicKey": cfg.Key,
		"UserName":  strings.Join(cfg.Names, ","),
		"Title":     title,
		"Content":   content,
		"TreeId":    "",
	}
	if sign, err = MerakSign(params, cfg.Secret); err != nil {
		log.Error("MerakNotify Failed to sign params: %+v: %+v", params, err)
		return err
	}
	params["Signature"] = sign
	if body, err = json.Marshal(params); err != nil {
		log.Error("MerakNotify Json %v, Err %v", params, err)
		return
	}
	if req, err = http.NewRequest(http.MethodPost, cfg.Host, bytes.NewReader(body)); err != nil {
		log.Error("MerakNotify NewRequest Err %v, Host %v", err, cfg.Host)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := struct {
		Action  string   `json:"Action"`
		RetCode int      `json:"RetCode"`
		Data    []string `json:"Data"`
	}{}
	if err = d.client.Do(ctx, req, &res); err != nil {
		return
	}
	if res.RetCode != 0 {
		err = errors.Errorf("Merak error: %d", res.RetCode)
		log.Error("Failed to send notify by merak with params: %+v: %+v", string(body), err)
		return
	}
	return
}

// MerakSign is used to sign for merak wechat msg
func MerakSign(params map[string]string, secret string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for _, k := range keys {
		buf.WriteString(k + params[k])
	}
	h := sha1.New()
	if _, err := h.Write(buf.Bytes()); err != nil {
		return "", errors.WithStack(err)
	}
	if _, err := h.Write([]byte(secret)); err != nil {
		return "", errors.WithStack(err)
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
