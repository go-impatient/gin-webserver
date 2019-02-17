package user

import (
	validator "gopkg.in/go-playground/validator.v9"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/schema/user"
	"github.com/moocss/go-webserver/src/pkg/auth"
)

const (
	TableName = "user"
)

// 创建用户
func Create(data *user.User) bool {
	return model.Create(TableName, &data)
}

// 更新用户
func Update(id uint64, data map[string]interface{}) bool {
	ok := model.Update(TableName, data, model.QueryParam{
		Where: []model.WhereParam{
			model.WhereParam{
				Field: "id",
				Prepare: id,
			},
		},
	})
	return ok
}

// 删除用户
func Delete(id uint64) bool {
	user := user.User{}
	user.ID = id
	ok := model.DeleteById(user.TableName(), &user)
	return ok
}

// 获取用户列表
func List(query model.QueryParam) ([]*user.User, bool) {
	data := make([]*user.User, 0)
	ok := model.GetMulti(TableName, &data, query)
	return data, ok
}

// 根据用户某一条件，统计用户数据条数
func Total(query model.QueryParam) (int, bool) {
	var count int
	ok := model.Count(TableName, &count, query)
	return count, ok
}

// 根据用户ID, 获取用户数据
func Get(id uint64) (*user.User, bool) {
	data := &user.User{}
	ok := model.GetById(TableName, &data, id)
	return data, ok
}

// 根据用户某一条件，获取用户数据
func GetOne(query model.QueryParam) (*user.User, bool) {
	data := &user.User{}
	ok := model.GetOne(TableName, &data, query)
	return data, ok
}

// Compare with the plain text password.
func Compare(data *user.User, pwd string) (err error) {
	err = auth.Compare(data.Password, pwd)
	return
}

// Encrypt the user password.
func Encrypt(data *user.User) (err error) {
	data.Password, err = auth.Encrypt(data.Password)
	return
}

// Validate the fields.
func Validate(data *user.User) error {
	validate := validator.New()
	return validate.Struct(data)
}