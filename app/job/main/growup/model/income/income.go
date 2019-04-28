package income

import (
	"go-common/library/time"
)

// AvIncome av income
type AvIncome struct {
	ID          int64
	AvID        int64
	MID         int64
	TagID       int64
	IsOriginal  int
	UploadTime  time.Time
	PlayCount   int64
	TotalIncome int64
	Income      int64
	TaxMoney    int64
	Date        time.Time
	BaseIncome  int64
}

// ColumnIncome column income
type ColumnIncome struct {
	ID          int64
	ArticleID   int64
	Title       string
	MID         int64
	TagID       int64
	ViewCount   int64
	Income      int64
	TotalIncome int64
	TaxMoney    int64
	UploadTime  time.Time
	Date        time.Time
	BaseIncome  int64
}

// BgmIncome sid + date: unique key
type BgmIncome struct {
	AID              int64
	SID              int64
	MID              int64
	CID              int64
	TaxMoney         int64
	Income           int64
	TotalIncome      int64
	Date             time.Time
	BaseIncome       int64
	DailyTotalIncome int64
}

// UpBusinessIncome av or column or bgm's middle-data-structure
type UpBusinessIncome struct {
	MID         int64
	Income      int64
	BaseIncome  int64
	Percent     float64
	Tax         int64
	PlayCount   int64
	AvCount     int64
	ColumnCount int64
	BgmCount    map[int64]bool
	ViewCount   int64
	Business    int // 1.视频 2.专栏 3.素材
}

// UpIncome up income
type UpIncome struct {
	ID                int64
	MID               int64
	AvCount           int64
	PlayCount         int64
	AvIncome          int64
	AudioIncome       int64
	ColumnCount       int64
	ColumnIncome      int64
	BgmIncome         int64
	BgmCount          int64
	AvTax             int64
	ColumnTax         int64
	BgmTax            int64
	TaxMoney          int64
	Income            int64
	TotalIncome       int64
	AvBaseIncome      int64
	ColumnBaseIncome  int64
	BgmBaseIncome     int64
	BaseIncome        int64
	AvTotalIncome     int64
	ColumnTotalIncome int64
	BgmTotalIncome    int64
	Date              time.Time
	IsDeleted         int
	DBState           int
}

// AvIncomeStat av income stat
type AvIncomeStat struct {
	AvID        int64
	MID         int64
	TagID       int64
	IsOriginal  int
	UploadTime  time.Time
	TotalIncome int64
	CTime       time.Time
	IsDeleted   int
	DataState   int // 1: insert 2: update
}

// ColumnIncomeStat column income stat
type ColumnIncomeStat struct {
	ArticleID   int64
	Title       string
	TagID       int64
	MID         int64
	UploadTime  time.Time
	TotalIncome int64
	DataState   int
}

// BgmIncomeStat bgm income stat
type BgmIncomeStat struct {
	SID         int64
	TotalIncome int64
	DataState   int
}

// UpIncomeStat up income stat
type UpIncomeStat struct {
	MID               int64
	TotalIncome       int64
	AvTotalIncome     int64
	ColumnTotalIncome int64
	BgmTotalIncome    int64
	IsDeleted         int
	DataState         int // 1: insert 2: update
}

// UpAccount up account
type UpAccount struct {
	MID                   int64
	HasSignContract       int
	State                 int
	TotalIncome           int64
	TotalUnwithdrawIncome int64
	TotalWithdrawIncome   int64
	IncIncome             int64
	LastWithdrawTime      time.Time
	Version               int64
	AllowanceState        int
	Nickname              string
	WithdrawDateVersion   string
	IsDeleted             int
	DataState             int // 1: insert 2: update
}
