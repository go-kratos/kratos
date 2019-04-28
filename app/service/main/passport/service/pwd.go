package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/main/passport/model"
)

// HistoryPwdCheck history pwd check
func (s *Service) HistoryPwdCheck(c context.Context, param *model.HistoryPwdCheckParam) (result string, err error) {
	var matchResult = make(map[string]int64, 3)
	var historyPwds []*model.HistoryPwd
	if historyPwds, err = s.d.HistoryPwds(c, param.Mid); err != nil {
		return
	}
	pwds := strings.Split(param.Pwd, ",")
	for _, history := range historyPwds {
		for _, pwd := range pwds {
			saltPwd := getSaltPwd(pwd, history.OldSalt)
			if saltPwd == history.OldPwd {
				matchResult[pwd] = 1
			}
		}
	}
	for _, pwd := range pwds {
		result = result + "," + strconv.FormatInt(matchResult[pwd], 10)
	}
	if len(result) > 0 {
		result = result[1:]
	}
	return
}

func getSaltPwd(pwd string, salt string) string {
	return md5Hex(fmt.Sprintf("%s>>BiLiSaLt<<%s", pwd, salt))
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
