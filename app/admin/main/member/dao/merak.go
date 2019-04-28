package dao

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

const (
	_privateKey = "37ba757817b4e9c45c7e97f6ed5eee4e1c7bac52"
	_publicKey  = "71f079db59672ecec5b8d6f252c4b59ab2a8a227mainsite@bilibili.com"
)

// MerakNotify send notify
func (d *Dao) MerakNotify(ctx context.Context, title, content string) error {
	params := map[string]string{
		"Action":    "CreateWechatMessage",
		"PublicKey": _publicKey,
		"UserName":  strings.Join(d.c.ReviewNotify.Users, ","),
		"Title":     title,
		"Content":   content,
		"TreeId":    "",
	}
	sign, err := MerakSign(params)
	if err != nil {
		log.Error("Failed to sign params: %+v: %+v", params, err)
		return err
	}
	params["Signature"] = sign
	b, err := json.Marshal(params)
	if err != nil {
		return errors.WithStack(err)
	}
	req, err := http.NewRequest(http.MethodPost, d.merakURL, bytes.NewReader(b))
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := struct {
		Action  string   `json:"Action"`
		RetCode int      `json:"RetCode"`
		Data    []string `json:"Data"`
	}{}
	if err = d.httpClient.Do(ctx, req, &res); err != nil {
		return err
	}
	if res.RetCode != 0 {
		err := errors.Errorf("Merak error: %d", res.RetCode)
		log.Error("Failed to send notify by merak with params: %+v: %+v", string(b), err)
		return err
	}
	return nil
}

// MerakSign is
func MerakSign(params map[string]string) (string, error) {
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
	if _, err := h.Write([]byte(_privateKey)); err != nil {
		return "", errors.WithStack(err)
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
