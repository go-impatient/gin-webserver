package user

import (
	"github.com/moocss/go-webserver/src/schema"
)

type Token struct {
	schema.Base

	Raw    string `gorm:"type:varchar(100);column:username;not null" json:"username"`
	User   User   `gorm:"save_associations:false"" json:"user"`
	UserID uint   `json:"user_id"`
}
