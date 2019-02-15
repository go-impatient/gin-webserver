package model

// User represents a registered user.
type User struct {
	Model
	Username string `gorm:"column:username;not null" json:"username" binding:"required" validate:"min=1,max=32"`
	Password string `gorm:"column:password;not null" json:"password" binding:"required" validate:"min=5,max=128"`
}

func (u *User) TableName() string {
	return "tb_users"
}
