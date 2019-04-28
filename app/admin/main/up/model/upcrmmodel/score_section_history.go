package upcrmmodel

import "go-common/library/time"

//ScoreSectionHistory score section
type ScoreSectionHistory struct {
	ID           uint32    `gorm:"column:id"`
	GenerateDate time.Time `gorm:"column:generate_date"`
	ScoreType    int       `gorm:"column:score_type"`
	Section0     int       `gorm:"column:section_0"`
	Section1     int       `gorm:"column:section_1"`
	Section2     int       `gorm:"column:section_2"`
	Section3     int       `gorm:"column:section_3"`
	Section4     int       `gorm:"column:section_4"`
	Section5     int       `gorm:"column:section_5"`
	Section6     int       `gorm:"column:section_6"`
	Section7     int       `gorm:"column:section_7"`
	Section8     int       `gorm:"column:section_8"`
	Section9     int       `gorm:"column:section_9"`
	CTime        time.Time `gorm:"column:ctime"`
	Mtime        time.Time `gorm:"column:mtime"`
}
