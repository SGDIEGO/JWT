package repositories

import (
	"database/sql"
	"errors"

	"github.com/SGDIEGO/JWT/internal/domains"
	"github.com/SGDIEGO/JWT/internal/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) ports.UserRepositoryI {
	return &UserRepository{
		db: db,
	}
}

func (uR *UserRepository) GetUsers() (*[]domains.Users, error) {
	var users []domains.Users
	var userGet domains.Users

	var queryGet = "SELECT * FROM `users`.`user`"
	result, err := uR.db.Query(queryGet)

	if err != nil {
		return nil, errors.New("error in database")
	}

	for result.Next() {
		err := result.Scan(&userGet.Userid, &userGet.Name, &userGet.Email, &userGet.Date, &userGet.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, userGet)
	}

	return &users, nil
}

func (uR *UserRepository) GetUserByName(userName string) (*domains.Users, error) {

	var userGet domains.Users

	var queryGetUser = "SELECT `name`,`email`,`password`,`date` FROM `users`.`user` WHERE name=?"

	err := uR.db.QueryRow(queryGetUser, userName).Scan(&userGet.Name, &userGet.Email, &userGet.Password, &userGet.Date)

	if err != nil {
		return nil, err
	}

	return &userGet, nil
}

func (uR *UserRepository) SaveUser(user *domains.Users) error {

	var querySave = "INSERT INTO `users`.`user`(`user_id`,`name`,`email`,`password`,`date`) VALUES (NULL,?,?,?, ?);"

	result, err := uR.db.Exec(querySave, user.Name, user.Email, user.Password, user.Date)

	if err != nil {
		return errors.New(err.Error())
	}

	if _, err = result.LastInsertId(); err != nil {
		return err
	}

	return nil
}
