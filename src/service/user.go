package service

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/moocss/gin-webserver/src/model"
)

func (service *defaultService) FindUser(username string) (*model.User, error) {
	if govalidator.IsNull(username) {
		return nil, errors.New("用户名不能为空")
	}
	if !govalidator.StringLength(username, "1", "255") {
		return nil, errors.New("用户名不能超过1到255个字节")
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
	if govalidator.IsInt(string(id)) {
		return nil, errors.New("用户ID不是数字")
	}

	detail, ok := service.dao.FindUser(id)

	if govalidator.IsNull(detail.Username) {
		return nil, errors.New("用户不存在")
	}

	if !ok {
		return nil, errors.New("获取用户详情数据失败")
	}

	return detail, nil
}

func (service *defaultService) DeleteUser(data *model.User) error {
	panic("implement me")
}

func (service *defaultService) CreateUser(data *model.User) error {
	//service.mtx.Lock()
	//defer service.mtx.Unlock()

	panic("")
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
