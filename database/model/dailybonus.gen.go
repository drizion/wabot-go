// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDailyBonu = "DailyBonus"

// DailyBonu mapped from table <DailyBonus>
type DailyBonu struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	BotUserID string    `gorm:"column:botUserId;not null" json:"botUserId"`
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName DailyBonu's table name
func (*DailyBonu) TableName() string {
	return TableNameDailyBonu
}
