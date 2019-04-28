package sobot

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/library/log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/service/main/workflow/model/sobot"
)

const (
	_sobotTicketModifyURL = "/ws/updateStatusBilibili/4"
	_sobotAddTicketURL    = "/ws/addCustomerTicketBilibili/4"
	_sobotAddReplyURL     = "/ws/addCustomerReplyInfoBilibili/4"
	_sobotTicketInfoURL   = "/ws/queryTicketReplyByCustomerListBilibili/4"
)

// SobotTicketInfo get ticket into
func (d *Dao) SobotTicketInfo(c context.Context, ticketID int32) (res json.RawMessage, err error) {
	var (
		req *http.Request
	)
	params := url.Values{}
	params.Set("companyId", d.c.HTTPClient.Sobot.Secret)
	params.Set("ticketId", strconv.Itoa(int(ticketID)))
	sign := md5.Sum([]byte(fmt.Sprintf("%s%s%d", d.c.HTTPClient.Sobot.Secret, d.c.HTTPClient.Sobot.Key, ticketID)))
	params.Set("sobotKey", hex.EncodeToString(sign[:]))
	if req, err = http.NewRequest("GET", d.ticketInfoURL+"?"+params.Encode(), nil); err != nil {
		log.Error("http.NewRequest(GET,%s) error(%v)", d.ticketInfoURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.httpSobot.Do(c, req, &res); err != nil {
		log.Error("d.httpSobot.Do() error(%v)", err)
		return
	}
	return
}

// SobotAddTicket add ticket to sobot
func (d *Dao) SobotAddTicket(c context.Context, tp *sobot.TicketParam) (err error) {
	var (
		req *http.Request
		res struct {
			RetCode string `json:"retCode"`
			Item    string `json:"item"`
		}
	)
	params := url.Values{}
	params.Set("fileStr", tp.FileStr)
	params.Set("companyId", d.c.HTTPClient.Sobot.Secret)
	params.Set("customerName", tp.CustomerName)
	params.Set("customerQq", tp.CustomerQQ)
	params.Set("customerNick", tp.CustomerNick)
	params.Set("customerEmail", tp.CustomerEmail)
	params.Set("customerPhone", tp.CustomerPhone)
	params.Set("customerSource", strconv.Itoa(int(tp.CustomerSource)))
	params.Set("ticketId", strconv.Itoa(int(tp.TicketID)))
	params.Set("ticketTitle", tp.TicketTitle)
	params.Set("ticketContent", tp.TicketContent)
	params.Set("ticketLevel", strconv.Itoa(int(tp.TicketLevel)))
	params.Set("ticketStatus", strconv.Itoa(int(tp.TicketStatus)))
	params.Set("ticketFrom", strconv.Itoa(int(sobot.TicketFrom)))
	sign := md5.Sum([]byte(fmt.Sprintf("%s%s%d", d.c.HTTPClient.Sobot.Secret, d.c.HTTPClient.Sobot.Key, tp.TicketID)))
	params.Set("sobotKey", hex.EncodeToString(sign[:]))
	if req, err = http.NewRequest("POST", d.ticketAddURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("http.NewRequest(POST,%s) error(%v)", d.ticketAddURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.httpSobot.Do(c, req, &res); err != nil {
		log.Error("d.httpSobot.Do() error(%v)", err)
		return
	}
	if res.RetCode != sobot.EcodeOK {
		log.Error("d.httpSobot.Do() url(%s) ecode(%s)", d.ticketAddURL, res.RetCode)
		err = errors.New(res.RetCode)
		return
	}
	return
}

// SobotAddReply add reply to sobot
func (d *Dao) SobotAddReply(c context.Context, rp *sobot.ReplyParam) (err error) {
	var (
		req *http.Request
		res struct {
			RetCode string `json:"retCode"`
			Item    string `json:"item"`
		}
	)
	params := url.Values{}
	params.Set("companyId", d.c.HTTPClient.Sobot.Secret)
	params.Set("customerEmail", rp.CustomerEmail)
	params.Set("replyContent", rp.ReplyContent)
	params.Set("ticketId", strconv.Itoa(int(rp.TicketID)))
	params.Set("replyType", strconv.Itoa(int(rp.ReplyType)))
	params.Set("startType", strconv.Itoa(int(rp.StartType)))
	sign := md5.Sum([]byte(fmt.Sprintf("%s%s%d", d.c.HTTPClient.Sobot.Secret, d.c.HTTPClient.Sobot.Key, rp.TicketID)))
	params.Set("sobotKey", hex.EncodeToString(sign[:]))
	if req, err = http.NewRequest("POST", d.replyAddURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("http.NewRequest(POST,%s) error(%v)", d.replyAddURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.httpSobot.Do(c, req, &res); err != nil {
		log.Error("d.httpSobot.Do() error(%v)", err)
		return
	}
	if res.RetCode != sobot.EcodeOK {
		log.Error("d.httpSobot.Do() url(%s) ecode(%s)", d.replyAddURL, res.RetCode)
		err = errors.New(res.RetCode)
		return
	}
	return
}

// SobotTicketModify modify ticket
func (d *Dao) SobotTicketModify(c context.Context, tp *sobot.TicketParam) (err error) {
	// http
	var (
		req *http.Request
		res struct {
			RetCode string `json:"retCode"`
			Item    string `json:"item"`
		}
	)
	params := url.Values{}
	params.Set("companyId", d.c.HTTPClient.Sobot.Secret)
	params.Set("customerEmail", tp.CustomerEmail)
	params.Set("ticketId", strconv.Itoa(int(tp.TicketID)))
	params.Set("ticketFrom", strconv.Itoa(int(sobot.TicketFrom)))
	params.Set("ticketStatus", strconv.Itoa(int(tp.TicketStatus)))
	params.Set("startType", strconv.Itoa(int(tp.StartType)))
	sign := md5.Sum([]byte(fmt.Sprintf("%s%s%d", d.c.HTTPClient.Sobot.Secret, d.c.HTTPClient.Sobot.Key, tp.TicketID)))
	params.Set("sobotKey", hex.EncodeToString(sign[:]))
	if req, err = http.NewRequest("POST", d.ticketModifyURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("http.NewRequest(POST,%s) error(%v)", d.ticketModifyURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.httpSobot.Do(c, req, &res); err != nil {
		log.Error("d.httpSobot.Do() error(%v)", err)
		return
	}
	if res.RetCode != sobot.EcodeOK {
		log.Error("d.httpSobot.Do() url(%s) ecode(%s)", d.ticketModifyURL, res.RetCode)
		err = errors.New(res.RetCode)
		return
	}
	return
}
