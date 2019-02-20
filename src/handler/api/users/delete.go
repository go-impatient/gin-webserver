package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/pkg/errno"
	"github.com/moocss/go-webserver/src/service"
)

func HandleDelete(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("id"))
		user := &model.User{}
		user.ID = uint64(userId)

		if err := userService.DeleteUser(user); err != nil {
			model.SendResult(c, errno.ErrDatabase, nil)
			return
		}

		model.SendResult(c, nil, nil)
	}
}
