package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/pkg/errno"
)

// Result represents HTTP response body.
type Result struct {
	Code 			int         `json:"code"` // return code, 0 for succ
	Message  	string      `json:"message"`  // message
	Data 			interface{} `json:"data"` // data object
}

// 输出返回结果
func SendResult(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	c.JSON(http.StatusOK, &Result{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
