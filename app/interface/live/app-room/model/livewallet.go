package model

//LiveWallet 直播钱包数据
type LiveWallet struct {
	Gold            int64 `json:"gold"`
	Silver          int64 `json:"silver"`
	GoldRechargeCnt int64 `json:"gold_recharge_cnt"`
	GoldPayCnt      int64 `json:"gold_pay_cnt"`
	SilverPayCnt    int64 `json:"silver_pay_cnt"`
}
