package model

import (
	"encoding/json"
	"time"
)

type Base struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"deletedAt"`
}

// 将数据序列化为JSON字符串
func JsonMarshal(v interface{}) (str string) {
	res, err := json.Marshal(v)
	if err != nil {
		str = ""
	}
	return string(res)
}