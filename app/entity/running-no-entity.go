package entity

import "time"

type RunningNoEntity struct {
	Name           string `gorm:"primaryKey"`
	Value          int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RunningNoEntity) TableName() string {
	return "running_no"
}
