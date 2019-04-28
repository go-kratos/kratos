package income

import (
	"fmt"
	"strconv"

	model "go-common/app/admin/main/growup/model/income"
)

func formatUpWithdraw(ups []*model.UpWithdrawRes, isDeleted int) (data [][]string) {
	if len(ups) <= 0 {
		return
	}
	data = make([][]string, len(ups)+1)
	if isDeleted == 0 {
		data[0] = []string{"UID", "昵称", "已提现", "未提现", "最近提现时间"}
		for i := 1; i <= len(ups); i++ {
			up := ups[i-1]
			data[i] = []string{
				strconv.FormatInt(up.MID, 10),
				up.Nickname,
				up.WithdrawIncome,
				up.UnWithdrawIncome,
				up.WithdrawDate,
			}
		}
	} else {
		data[0] = []string{"UID", "昵称", "禁止提现", "禁止时间", "已提现", "最近提现时间"}
		for i := 1; i <= len(ups); i++ {
			up := ups[i-1]
			data[i] = []string{
				strconv.FormatInt(up.MID, 10),
				up.Nickname,
				up.UnWithdrawIncome,
				up.MTime.Time().Format(_layout),
				up.WithdrawIncome,
				up.WithdrawDate,
			}
		}
	}
	return
}

func formatUpIncomeWithdraw(ups []*model.UpIncomeWithdraw) (data [][]string) {
	if len(ups) <= 0 {
		return
	}
	data = make([][]string, len(ups)+1)
	data[0] = []string{"最近一次提现日期", "UP主ID", "UP主昵称", "已经提现的收入"}
	for i := 1; i <= len(ups); i++ {
		up := ups[i-1]
		data[i] = []string{
			up.MTime.Time().Format("2006-01-02"),
			strconv.FormatInt(up.MID, 10),
			up.Nickname,
			fmt.Sprintf("%.2f", fromYuanToFen(up.WithdrawIncome)),
		}
	}
	return
}

func formatUpIncome(ups []*model.UpIncome) (data [][]string) {
	if len(ups) <= 0 {
		return
	}
	data = make([][]string, 1)
	data[0] = []string{"时间", "UID", "昵称", "新增收入", "稿件数", "基础收入", "额外收入", "违规扣除", "扣税金额", "累计收入"}
	for _, up := range ups {
		data = append(data, []string{
			up.DateFormat,
			strconv.FormatInt(up.MID, 10),
			up.Nickname,
			fmt.Sprintf("%.2f", fromYuanToFen(up.Income)),
			strconv.FormatInt(up.AvCount, 10),
			fmt.Sprintf("%.2f", fromYuanToFen(up.BaseIncome)),
			fmt.Sprintf("%.2f", fromYuanToFen(up.ExtraIncome)),
			fmt.Sprintf("%.2f", fromYuanToFen(up.Breach)),
			fmt.Sprintf("%.2f", fromYuanToFen(up.TaxMoney)),
			fmt.Sprintf("%.2f", fromYuanToFen(up.TotalIncome)),
		})
	}
	return
}

func formatBreach(breachs []*model.AvBreach) (data [][]string) {
	if len(breachs) <= 0 {
		return
	}
	ctype := []string{"视频", "音频", "专栏", "素材"}
	data = make([][]string, 1)
	data[0] = []string{"日期", "稿件id", "稿件类型", "投稿时间", "UID", "up主昵称", "扣除金额", "扣除原因"}
	for _, b := range breachs {
		data = append(data, []string{
			b.CDate.Time().Format(_layout),
			strconv.FormatInt(b.AvID, 10),
			ctype[b.CType],
			b.UploadTime.Time().Format(_layout),
			strconv.FormatInt(b.MID, 10),
			b.Nickname,
			fmt.Sprintf("%.2f", fromYuanToFen(b.Money)),
			b.Reason,
		})
	}
	return
}
