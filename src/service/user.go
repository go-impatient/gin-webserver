package service

import (
	"github.com/moocss/go-webserver/src/schema/user"
)

//func GetUserById(id uint64) (*user.User, error) {
//	user := &user.User{}
//	user.ID = id
//
//	detail, err := UserDetail(user);
//	if  err != nil {
//		return nil, err
//	}
//
//	return detail, nil
//}
//
//func Delete(id uint64) error {
//	user := &user.User{}
//	user.ID = id
//
//	detail, err := UserDetail(user);
//	if err != nil {
//		return err
//	}
//
//	ok := userModel.Delete(detail.ID)
//	if !ok {
//		return errors.New("用户删除失败")
//	}
//	return nil
//}
//
//
//func UserDetail(data *user.User) (*user.User, error) {
//	if data.ID == 0 {
//		return nil, errors.New("用户ID不能为空")
//	}
//
//	detail, ok := userModel.Get(data.ID)
//	if !ok {
//		return nil, errors.New("获取用户详情数据失败")
//	}
//	if detail.ID == 0 {
//		return nil, errors.New("用户不存在")
//	}
//
//	return detail, nil
//}
//
//func GetUserByName(data *user.User) (*user.User, error) {
//	if data.Username == "" {
//		return nil, errors.New("用户名不能为空")
//	}
//	detail, ok := userModel.GetOne(model.QueryParam{
//		Where: []model.WhereParam{
//			model.WhereParam{
//				Field: "username",
//				Prepare: data.Username,
//			},
//		},
//	})
//	if !ok {
//		return nil, errors.New("获取用户详情数据失败")
//	}
//
//	return detail, nil
//}
//
//func CreateOrUpdate(data *user.User) error {
//	var ok bool
//	user := &user.User{
//		Username: data.Username,
//		Password: data.Password,
//		Email:    data.Email,
//	}
//
//	// Validate the data.
//	if err := userModel.Validate(user); err != nil {
//		return errors.New("错误字段")
//	}
//
//
//	if user.Password != "" {
//		// Encrypt the user password.
//		if err := userModel.Encrypt(user); err != nil {
//			return errors.New("密码加密失败")
//		}
//	}
//
//	if user.ID > 0 {
//		updateData := map[string]interface{}{
//			"username": user.Username,
//			"password": user.Password,
//			"email": user.Email,
//		}
//		ok = userModel.Update(user.ID, updateData)
//	} else {
//		ok = userModel.Create(user)
//	}
//	if !ok {
//		return errors.New("更新用户数据失败")
//	}
//	return nil
//}

func (service *defaultService) ShowUser(string) (*user.User, error) {
	panic("implement me")
}

func (service *defaultService) DeleteUser(*user.User) error {
	panic("implement me")
}
