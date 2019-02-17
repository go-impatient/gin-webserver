package schema

import (
	"encoding/json"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// 创建一个User数据模型
type User struct {
	Base

	Username string `gorm:"type:varchar(100);column:username;not null" json:"username"`
	Password string `gorm:"type:varchar(50);column:password;not null" json:"password"`
	Email    string `gorm:"type:varchar(100);column:email;unique;not null;" json:"email"`
}

// TableName, 获取User表名称
func TableName() string {
	return "user"
}

// UserFrom, 解析JSON字符串并返回一个 User 实例
func UserFrom(str string) (*User, error) {
	user := new(User)
	if err := json.Unmarshal([]byte(str), user); err != nil {
		return nil, err
	}
	return user, nil
}

// String, 返回一个为JSON 字符串的用户信息
func (u *User) String() string {
	return JsonMarshal(u)
}

// Result, 返回一个 UserResult 实例
func (u *User) Result() *UserResult {
	return &UserResult{
		ID:      u.ID,
		Username: u.Username,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// 创建一个脱敏的User数据模型
type UserResult struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Email 		string `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// String, 返回一个为JSON 字符串的脱敏用户信息
func (u *UserResult) String() string {
	return JsonMarshal(u)
}

// Validate the fields.
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}