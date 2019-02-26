package users

import (
	"github.com/moocss/gin-webserver/src/pkg/log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moocss/gin-webserver/src/model"
	"github.com/moocss/gin-webserver/src/pkg/errno"
	"github.com/moocss/gin-webserver/src/service"
)

// HandleFind 控制器， 按照用户名查询
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

// HandleFindById 控制器， 按照用户ID查询
func HandleFindById(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("id"))

		log.Debugf("用户ID: %s", userId)

		user, err := userService.FindUserById(uint64(userId))

		if err != nil {
			model.SendResult(c, errno.ErrUserNotFound, nil)
			return
		}

		model.SendResult(c, nil, user)
	}
}
