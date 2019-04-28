package http

import (
	"encoding/base64"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func userOperatorIP(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	operator := params.Get("operator")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	if operator == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ipStr := metadata.String(c, metadata.RemoteIP)
	ip := model.InetAtoN(ipStr)
	switch operator {
	case "unicom":
		res["data"] = unicomSvc.UserUnciomIP(ip, ipStr, usermob, mobiApp, build, time.Now())
	case "mobile":
		res["data"] = mobileSvc.IsMobileIP(ip, ipStr, usermob)
	default:
		c.JSON(nil, ecode.RequestErr)
		return
	}
	returnDataJSON(c, res, nil)
}

func mOperatorIP(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	operator := params.Get("operator")
	ipStr := params.Get("ip")
	ip := model.InetAtoN(ipStr)
	switch operator {
	case "unicom":
		var usermobStr string
		if usermob != "" {
			var (
				_aesKey = []byte("9ed226d9")
			)
			bs, err := base64.StdEncoding.DecodeString(usermob)
			if err != nil {
				log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermob, err)
				c.JSON(nil, ecode.RequestErr)
				return
			}
			bs, err = unicomSvc.DesDecrypt(bs, _aesKey)
			if err != nil {
				log.Error("unicomSvc.DesDecrypt error(%v)", err)
				c.JSON(nil, ecode.RequestErr)
				return
			}
			if len(bs) > 32 {
				usermobStr = string(bs[:32])
			} else {
				usermobStr = string(bs)
			}
		}
		res["data"] = unicomSvc.UserUnciomIP(ip, ipStr, usermobStr, "missevan", 0, time.Now())
	case "mobile":
		res["data"] = mobileSvc.IsMobileIP(ip, ipStr, usermob)
	default:
		c.JSON(nil, ecode.RequestErr)
		return
	}
	returnDataJSON(c, res, nil)
}
