package model

import (
	"encoding/json"
	"time"
	"sync"
)

// User represents a registered user.
type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}


func (u *UserModel) TableName() string {
	return "tb_users"
}

// UserFrom parse JSON string and returns a User intance.
func UserFrom(str string) (*UserModel, error) {
	user := new(UserModel)
	if err := json.Unmarshal([]byte(str), user); err != nil {
		return nil, err
	}
	return user, nil
}

// String returns JSON string with full user info
func (u *UserModel) String() string {
	return jsonMarshal(u)
}

// Result returns UserResult intance
func (u *UserModel) Result() *UserResult {
	return &UserResult{
		ID:      u.ID,
		Username: u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserResult represents desensitized user
type UserResult struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// String returns JSON string with desensitized user info
func (u *UserResult) String() string {
	return jsonMarshal(u)
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserResult
}
