package users

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/pkg/errno"
	"github.com/moocss/go-webserver/src/service"
	"strconv"
)

func HandleFind(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		user, err := userService.FindUser(username)

		if err != nil {
			model.SendResult(c, errno.ErrUserNotFound, nil)
			return
		}

		model.SendResult(c, nil, user)
	}
}

func HandleFindById(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("id"))

		user, err := userService.FindUserById(uint64(userId))

		if err != nil {
			model.SendResult(c, errno.ErrUserNotFound, nil)
			return
		}

		model.SendResult(c, nil, user)
	}
}
