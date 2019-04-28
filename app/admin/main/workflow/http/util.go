package http

import (
	"strconv"
	"strings"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

type intsParam struct {
	value string
	p     *[]int64
}
type intParam struct {
	value string
	p     *int64
}

func dealNumsmap(intsparams []*intsParam) error {
	for _, isp := range intsparams {
		if isp.value != "" {
			ids, err := xstr.SplitInts(isp.value)
			if err != nil {
				log.Error("strconv.ParseInt(%s) error(%v)", isp.value, err)
				return err
			}
			*isp.p = ids
		}
	}

	return nil
}

func dealNummap(intparams []*intParam) error {
	for _, ip := range intparams {
		if ip.value != "" {
			id, err := strconv.ParseInt(ip.value, 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt(%s) error(%v)", ip.value, err)
				return err
			}
			*ip.p = id
		}
	}

	return nil
}

// adjustOrder will convert order field from request to search service
func adjustOrder(subject string, order string) string {
	SubOrderFields := map[string]map[string]string{
		"group": {
			"last_time": "lasttime",
		},

		"tag": {
			"count":    "tag_all_num",
			"handling": "tag_todo_num",
		},

		"challenge": {},

		"log": {
			"ctime": "opt_ctime",
		},
	}

	orderFields, ok := SubOrderFields[subject]
	if !ok {
		return order
	}

	field, ok := orderFields[order]
	if !ok {
		return order
	}

	return field
}

// read permissions of an admin
// only use workflow permissions in platform
// like WF_BUSINESS_2_ROUND_11
// round > 10 means feedback flow
// one admin only has one flow permission in a business
func parsePermission(permissions []string) (permissionMap map[int8]int64) {
	//permissionMap map[business]round
	//var regex = `^WF_BUSINESS_[1-9]\d*_ROUND_[1-9]\d*`
	permissionMap = make(map[int8]int64)
	for _, str := range permissions {
		splitStr := strings.Split(str, "_")
		if len(splitStr) == 5 && splitStr[0] == "WF" && splitStr[1] == "BUSINESS" && splitStr[3] == "ROUND" {
			business, err := strconv.ParseInt(splitStr[2], 10, 32)
			if err != nil {
				continue
			}
			round, err := strconv.ParseInt(splitStr[4], 10, 32)
			if err != nil {
				continue
			}
			permissionMap[int8(business)] = round
			//todo: support round??
		}
	}
	// todo: use cache?
	return
}

func adminInfo(ctx *bm.Context) (adminID int64, adminName string) {
	if IUid, ok := ctx.Get("uid"); ok {
		adminID = IUid.(int64)
	}
	if IUName, ok := ctx.Get("username"); ok {
		adminName = IUName.(string)
	}
	return
}
