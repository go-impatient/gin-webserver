package service

import (
	"errors"
	"github.com/moocss/go-webserver/src/model"
)

func (service *defaultService) FindUser(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	detail, ok := service.dao.FindUserOne(&model.QueryParam{
		Where: []model.WhereParam{
			{
				Field:   "username",
				Prepare: username,
			},
		},
	})
	if !ok {
		return nil, errors.New("获取用户详情数据失败")
	}

	return detail, nil
}

func (service *defaultService) FindUserById(id uint64) (*model.User, error) {
	detail, err := service.UserDetail(id)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (service *defaultService) UserDetail(id uint64) (*model.User, error) {
	if id < 0 {
		return nil, errors.New("用户ID不能为空")
	}

	detail, ok := service.dao.FindUser(id)
	if !ok {
		return nil, errors.New("获取用户详情数据失败")
	}

	return detail, nil
}

func (service *defaultService) DeleteUser(data *model.User) error {
	panic("implement me")
}

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
//
