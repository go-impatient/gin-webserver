package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/pkg/errno"
	"github.com/moocss/go-webserver/src/schema/result"
	"github.com/moocss/go-webserver/src/schema/user"
	"github.com/moocss/go-webserver/src/service"
)

func GetUser(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		user, err := userService.ShowUser(username)

		if err != nil {
			result.SendResult(c, errno.ErrUserNotFound, nil)
			return
		}

		result.SendResult(c, nil, user)
	}
}

func DeleteUser(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("id"))
		user := &user.User{}
		user.ID = uint64(userId)

		if err := userService.DeleteUser(user); err != nil {
			result.SendResult(c, errno.ErrDatabase, nil)
			return
		}

		result.SendResult(c, nil, nil)
	}
}
