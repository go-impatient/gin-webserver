package user

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/gin-webserver/src/service"
)

func HandleFindUser(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
