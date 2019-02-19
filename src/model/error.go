package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// The Error contains error relevant information.
type Error struct {
	// The general error message
	Error string `json:"error"`

	// The http error code.
	ErrorCode int `json:"error_code"`

	// The http error code.
	ErrorDescription string `json:"error_description"`
}

func SendError(c *gin.Context) {
	c.JSON(http.StatusNotFound, &Error{
		Error:            http.StatusText(http.StatusNotFound),
		ErrorCode:        http.StatusNotFound,
		ErrorDescription: "page not found",
	})
}

