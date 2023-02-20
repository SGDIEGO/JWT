package ports

import "github.com/SGDIEGO/JWT/internal/domains"

type UserRepositoryI interface {
	GetUsers() (*[]domains.Users, error)
	GetUserByName(userName string) (*domains.Users, error)
	SaveUser(user *domains.Users) error
}

type UserServiceI interface {
	GetUsers() (*[]domains.Users, error)
	GetUserByName(userName string) (*domains.Users, error)
	CreateUser(user *domains.Users) error
}
