// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameChatgpt = "Chatgpt"

// Chatgpt mapped from table <Chatgpt>
type Chatgpt struct {
	ID        int32     `gorm:"column:id;primaryKey;default:nextval('"Chatgpt_id_seq"" json:"id"`
	From      string    `gorm:"column:from;not null" json:"from"`
	Message   string    `gorm:"column:message;not null" json:"message"`
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName Chatgpt's table name
func (*Chatgpt) TableName() string {
	return TableNameChatgpt
}
