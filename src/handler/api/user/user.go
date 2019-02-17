package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/schema/result"
	userService "github.com/moocss/go-webserver/src/service/user"
	"github.com/moocss/go-webserver/src/pkg/errno"
)

func Get(c *gin.Context)  {
	userId, _ := strconv.Atoi(c.Param("id"))

	user, err := userService.GetUserById(uint64(userId))

	if err != nil {
		result.SendResult(c, errno.ErrUserNotFound, nil)
		return
	}

	result.SendResult(c, nil, user)
}

func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := userService.Delete(uint64(userId)); err != nil {
		result.SendResult(c, errno.ErrDatabase, nil)
		return
	}

	result.SendResult(c, nil, nil)
}