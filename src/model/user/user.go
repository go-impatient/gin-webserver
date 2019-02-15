package user

import (
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/schema"
)

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
	user.BaseModel.ID = id
	ok := model.DeleteByPk(schema.UserTableName(), &user)
	return ok
}
