// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSiteUser = "SiteUser"

// SiteUser mapped from table <SiteUser>
type SiteUser struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	BotUserID string    `gorm:"column:botUserId;not null" json:"botUserId"`
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName SiteUser's table name
func (*SiteUser) TableName() string {
	return TableNameSiteUser
}
