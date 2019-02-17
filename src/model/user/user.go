package model

import (
	"encoding/json"
	"time"
)

// 创建一个User数据模型
type User struct {
	Base

	Username string `gorm:"type:varchar(100);column:username;not null" json:"username"`
	Password string `gorm:"type:varchar(50);column:password;not null" json:"password"`
	Email    string `gorm:"type:varchar(100);column:email;unique;not null;" json:"email"`
}

// TableName, 获取User表名称
func UserTableName() string {
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
	return jsonMarshal(u)
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
	return jsonMarshal(u)
}


func Create(data *schema.User) bool {
	return model.Create(schema.UserTableName(), data)
}

func Update(id uint64, data map[string]interface{}) bool {
	ok := model.Update(schema.UserTableName(), data, model.QueryParam{
		Where: []model.WhereParam{
			model.WhereParam{
				Field: "id",
				Prepare: id,
			},
		},
	})
	return ok
}

func List(query model.QueryParam) ([]schema.User, bool) {
	var data []schema.User
	ok := model.GetMulti(schema.UserTableName(), &data, query)
	return data, ok
}

func Total(query model.QueryParam) (int, bool) {
	var count int
	ok := model.Count(schema.UserTableName(), &count, query)
	return count, ok
}

func Get(id uint64) (schema.User, bool) {
	var data schema.User
	ok := model.GetByPk(schema.UserTableName(), &data, id)
	return data, ok
}

func GetOne(query model.QueryParam) (schema.User, bool) {
	var data schema.User
	ok := model.GetOne(schema.UserTableName(), &data, query)
	return data, ok
}

func Delete(id uint64) bool {
	user := schema.User{}
	user.ID = id
	ok := model.DeleteByPk(schema.UserTableName(), &user)
	return ok
}
