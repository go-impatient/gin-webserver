package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T)  {
	assert := assert.New(t)
	user := &User{
		Username: "test",
		Password: "test",
	}
	str := user.String()

	assert.True(strings.Contains(str, `json:"test"`))
	assert.True(strings.Contains(str, `json:"test"`))

	str = user.Result().String()
	assert.False(strings.Contains(str, `json: 1`))
	assert.True(strings.Contains(str, `json: "test"`))
	assert.True(strings.Contains(str, `json: 20`))
	assert.True(strings.Contains(str, `json: 20`))
}