// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameEuNunca = "EuNunca"

// EuNunca mapped from table <EuNunca>
type EuNunca struct {
	ID        int32     `gorm:"column:id;primaryKey;default:nextval('"EuNunca_id_seq"" json:"id"`
	Text      string    `gorm:"column:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName EuNunca's table name
func (*EuNunca) TableName() string {
	return TableNameEuNunca
}
