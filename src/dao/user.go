package dao

import (
	"github.com/moocss/go-webserver/src/model"
)

// user 数据表
const (
	TableName = "user"
)

// CreateUser 创建用户
func (d *Dao) CreateUser(data *model.User) bool {
	return d.Create(TableName, &data)
}

// UpdateUser 更新用户
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

// DeleteUser 删除用户
func (d *Dao) DeleteUser(id uint64) bool {
	user := model.User{}
	user.ID = id
	ok := d.DeleteById(user.TableName(), &user)
	return ok
}

// ListUser 获取用户列表
func (d *Dao) ListUser(query *model.QueryParam) ([]*model.User, bool) {
	data := make([]*model.User, 0)
	ok := d.FindMulti(TableName, &data, query)
	return data, ok
}

// FindUserTotal 根据用户某一条件，统计用户数据条数
func (d *Dao) FindUserTotal(query *model.QueryParam) (int, bool) {
	var count int
	ok := d.Count(TableName, &count, query)
	return count, ok
}

// FindUser 根据用户ID, 获取用户数据
func (d *Dao) FindUser(id uint64) (*model.User, bool) {
	data := &model.User{}
	ok := d.FindById(TableName, &data, id)
	return data, ok
}

// FindUserOne 根据用户某一条件，获取用户数据
func (d *Dao) FindUserOne(query *model.QueryParam) (*model.User, bool) {
	data := &model.User{}
	ok := d.FindOne(TableName, &data, query)
	return data, ok
}
