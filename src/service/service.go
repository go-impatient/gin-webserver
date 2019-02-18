package service

import (
	"github.com/moocss/go-webserver/src/schema/user"
)

// TokenService interface
type TokenService interface {
	ShowToken(string) (*user.Token, error)
	//CreateToken(*user.Token) error
	//DeleteToken(*user.Token) error
	//ListTokensByUser(string) ([]*user.Token, error)
}

// UserService interface
type UserService interface {
	// ShowUserById(uint64) (*user.User, error)
	ShowUser(string) (*user.User, error)
	//ShowUserByToken(string) (*user.User, error)
	//CreateUser(*user.User) error
	//UpdateUser(*user.User) error
	DeleteUser(*user.User) error
	//ChangeUserPassword(*user.User) error
	//Login(string, string) (*user.User, bool, error)
	//ListUsers() ([]*user.User, error)
}

// Service interface combines all concrete model services
type Service interface {
	TokenService
	UserService
}

type defaultService struct{}

// NewService constructs a new service layer
func NewService() Service {
	return &defaultService{}
}
