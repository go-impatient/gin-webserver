package model

type Token struct {
	Base

	Raw    string `gorm:"type:varchar(100);column:username;not null" json:"username"`
	User   User   `gorm:"save_associations:false"" json:"user"`
	UserID uint   `json:"user_id"`
}
