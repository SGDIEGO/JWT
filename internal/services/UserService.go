package services

import (
	"time"

	"github.com/SGDIEGO/JWT/internal/domains"
	"github.com/SGDIEGO/JWT/internal/ports"
)

type UserService struct {
	UserRepository ports.UserRepositoryI
}

func NewUserService(userRepo ports.UserRepositoryI) ports.UserServiceI {
	return &UserService{
		UserRepository: userRepo,
	}
}

func (uS *UserService) GetUsers() (*[]domains.Users, error) {

	users, err := uS.UserRepository.GetUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uS *UserService) GetUserByName(userName string) (*domains.Users, error) {
	return uS.UserRepository.GetUserByName(userName)
}
func (uS *UserService) CreateUser(user *domains.Users) error {

	user.Email = user.Name + "@gmail.com"
	user.Date = time.Now().Day()

	return uS.UserRepository.SaveUser(user)
}
