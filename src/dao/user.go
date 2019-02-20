package dao

import (
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/pkg/auth"
	"gopkg.in/go-playground/validator.v9"
)

const (
	TableName = "user"
)

// 创建用户
func (d *Dao) CreateUser(data *model.User) bool {
	return d.Create(TableName, &data)
}

// 更新用户
func (d *Dao) UpdateUser(id uint64, data map[string]interface{}) bool {
	ok := d.Update(TableName, data, &model.QueryParam{
		Where: []model.WhereParam{
			model.WhereParam{
				Field:   "id",
				Prepare: id,
			},
		},
	})
	return ok
}

// 删除用户
func (d *Dao) DeleteUser(id uint64) bool {
	user := model.User{}
	user.ID = id
	ok := d.DeleteById(user.TableName(), &user)
	return ok
}

// 获取用户列表
func (d *Dao) ListUser(query *model.QueryParam) ([]*model.User, bool) {
	data := make([]*model.User, 0)
	ok := d.FindMulti(TableName, &data, query)
	return data, ok
}

// 根据用户某一条件，统计用户数据条数
func (d *Dao) FindUserTotal(query *model.QueryParam) (int, bool) {
	var count int
	ok := d.Count(TableName, &count, query)
	return count, ok
}

// 根据用户ID, 获取用户数据
func (d *Dao) FindUser(id uint64) (*model.User, bool) {
	data := &model.User{}
	ok := d.FindById(TableName, &data, id)
	return data, ok
}

// 根据用户某一条件，获取用户数据
func (d *Dao) FindUserOne(query *model.QueryParam) (*model.User, bool) {
	data := &model.User{}
	ok := d.FindOne(TableName, &data, query)
	return data, ok
}

// Compare with the plain text password.
func (d *Dao) Compare(data *model.User, pwd string) (err error) {
	err = auth.Compare(data.Password, pwd)
	return
}

// Encrypt the user password.
func (d *Dao) Encrypt(data *model.User) (err error) {
	data.Password, err = auth.Encrypt(data.Password)
	return
}

// Validate the fields.
func (d *Dao) Validate(data *model.User) error {
	validate := validator.New()
	return validate.Struct(data)
}
