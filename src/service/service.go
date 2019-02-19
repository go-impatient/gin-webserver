package service

import (
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/dao"
	"github.com/moocss/go-webserver/src/model"
)

// TokenService interface
type TokenService interface {
	ShowToken(string) (*model.Token, error)
	//CreateToken(*model.Token) error
	//DeleteToken(*model.Token) error
	//ListTokensByUser(string) ([]*model.Token, error)
}

// UserService interface
type UserService interface {
	// ShowUserById(uint64) (*model.User, error)
	ShowUser(string) (*model.User, error)
	//ShowUserByToken(string) (*model.User, error)
	//CreateUser(*model.User) error
	//UpdateUser(*model.User) error
	DeleteUser(*model.User) error
	//ChangeUserPassword(*model.User) error
	//Login(string, string) (*model.User, bool, error)
	//ListUsers() ([]*model.User, error)
}

// Service interface combines all concrete model services
type Service interface {
	TokenService
	UserService
}

type defaultService struct{
	d 	*dao.Dao
	cfg *config.Config
}

// New constructs a new service layer
func New(cfg *config.Config) Service {
	return &defaultService{
		d: dao.New(cfg),
		cfg: cfg,
	}
}
