// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameBotUser = "BotUser"

// BotUser mapped from table <BotUser>
type BotUser struct {
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	AdsCounter int32     `gorm:"column:adsCounter;not null" json:"adsCounter"`
	WaCoins    int32     `gorm:"column:waCoins;not null" json:"waCoins"`
	VerifiedAt time.Time `gorm:"column:verifiedAt" json:"verifiedAt"`
	CreatedAt  time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName BotUser's table name
func (*BotUser) TableName() string {
	return TableNameBotUser
}
