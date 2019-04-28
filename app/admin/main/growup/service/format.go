package service

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"go-common/app/admin/main/growup/model"
)

var (
	_accountState = map[int]string{
		1: "未申请",
		2: "待审核",
		3: "已签约",
		4: "已驳回",
		5: "主动退出",
		6: "被动退出",
		7: "封禁",
	}

	_avCategory = map[int]string{
		0:   "默认",
		1:   "动画",
		3:   "音乐",
		129: "舞蹈",
		4:   "游戏",
		36:  "科技",
		160: "生活",
		119: "鬼畜",
		155: "时尚",
		23:  "电影",
		11:  "电视剧",
		13:  "番剧",
		167: "国创",
		165: "广告",
		5:   "娱乐",
		177: "纪录片",
		181: "影视",
	}

	_cmCategory = map[int]string{
		1:  "游戏",
		2:  "动画",
		3:  "生活",
		16: "轻小说",
		28: "影视",
		29: "兴趣",
	}
)

// FormatCSV format to csv data
func FormatCSV(records [][]string) (data []byte, err error) {
	buf := new(bytes.Buffer)

	// add utf bom
	if len(records) > 0 {
		buf.WriteString("\xEF\xBB\xBF")
	}

	w := csv.NewWriter(buf)
	err = w.WriteAll(records)
	if err != nil {
		return
	}

	data = buf.Bytes()
	return
}

func formatBlacklist(blacklist []*model.Blacklist) (data [][]string) {
	if len(blacklist) <= 0 {
		return
	}
	data = make([][]string, len(blacklist)+1)
	data[0] = []string{"稿件id", "加入时间", "加入原因", "累计收入", "UP主UID", "UP主昵称", "业务类型"}
	for i := 0; i < len(blacklist); i++ {
		list := blacklist[i]
		var ctype, reasonStr string
		switch list.CType {
		case 0:
			ctype = "视频"
		case 1:
			ctype = "音频"
		case 2:
			ctype = "专栏"
		}

		switch list.Reason {
		case 1:
			reasonStr = "停止收入"
		case 2:
			reasonStr = "违规扣除"
		case 3:
			reasonStr = "私单"
		case 4:
			reasonStr = "绿洲商单"
		}

		data[i+1] = []string{
			strconv.FormatInt(list.AvID, 10),
			time.Unix(int64(list.CTime), 0).Format("2006-01-02 15:04:05"),
			reasonStr,
			strconv.FormatFloat(float64(list.Income)*0.01, 'f', 3, 64),
			strconv.FormatInt(list.MID, 10),
			list.Nickname,
			ctype,
		}
	}
	return
}

func formatUpInfo(ups []*model.UpInfo, states []int64, typ int) (data [][]string) {
	if len(ups) <= 0 {
		return
	}
	data = make([][]string, len(ups)+1)
	switch typ {
	case _video:
		data[0] = []string{"UID", "昵称", "原创稿件数", "稿件播放量", "稿件量", "分区", "粉丝数"}
	case _column:
		data[0] = []string{"UID", "昵称", "稿件数", "阅读量", "分区", "粉丝数"}
	case _bgm:
		data[0] = []string{"UID", "昵称", "素材量", "视频播放量", "素材使用量", "粉丝数"}
	}
	if len(states) == 0 {
		data[0] = append(data[0], "账号状态")
	} else {
		switch states[0] {
		case 2:
			data[0] = append(data[0], "申请时间")
		case 3:
			data[0] = append(data[0], "申请时间")
			data[0] = append(data[0], "签约时间")
		case 4:
			data[0] = append(data[0], "驳回时间")
			data[0] = append(data[0], "到期时间")
		case 5:
			data[0] = append(data[0], "退出时间")
			data[0] = append(data[0], "到期时间")
		case 6:
			data[0] = append(data[0], "退出时间")
		case 7:
			data[0] = append(data[0], "封禁时间")
			data[0] = append(data[0], "恢复时间")

		}
	}
	for i := 1; i <= len(ups); i++ {
		up := ups[i-1]
		data[i] = []string{
			strconv.FormatInt(up.MID, 10),
			up.Nickname,
		}
		switch typ {
		case _video:
			data[i] = append(data[i], strconv.Itoa(up.OriginalArchiveCount))
			data[i] = append(data[i], strconv.Itoa(up.TotalPlayCount))
			data[i] = append(data[i], strconv.Itoa(up.Avs))
			data[i] = append(data[i], _avCategory[up.MainCategory])
		case _column:
			data[i] = append(data[i], strconv.Itoa(up.ArticleCount))
			data[i] = append(data[i], strconv.Itoa(up.TotalViewCount))
			data[i] = append(data[i], _cmCategory[up.MainCategory])
		case _bgm:
			data[i] = append(data[i], strconv.Itoa(up.BGMs))
			data[i] = append(data[i], strconv.Itoa(up.BgmPlayCount))
			data[i] = append(data[i], strconv.Itoa(up.BgmApplyCount))
		}

		data[i] = append(data[i], strconv.Itoa(up.Fans))

		if len(states) == 0 {
			data[i] = append(data[i], _accountState[up.AccountState])
		} else {
			switch states[0] {
			case 2:
				data[i] = append(data[i], up.ApplyAt.Time().Format("2006-01-02"))
			case 3:
				data[i] = append(data[i], up.ApplyAt.Time().Format("2006-01-02"))
				data[i] = append(data[i], up.SignedAt.Time().Format("2006-01-02"))
			case 4:
				data[i] = append(data[i], up.RejectAt.Time().Format("2006-01-02"))
				data[i] = append(data[i], up.ExpiredIn.Time().Format("2006-01-02"))
			case 5:
				data[i] = append(data[i], up.QuitAt.Time().Format("2006-01-02"))
				data[i] = append(data[i], up.ExpiredIn.Time().Format("2006-01-02"))
			case 6:
				data[i] = append(data[i], up.DismissAt.Time().Format("2006-01-02"))
			case 7:
				data[i] = append(data[i], up.ForbidAt.Time().Format("2006-01-02"))
				data[i] = append(data[i], up.ExpiredIn.Time().Format("2006-01-02"))
			}
		}
	}
	return
}
