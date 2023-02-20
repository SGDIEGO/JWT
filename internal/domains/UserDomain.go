package domains

import (
	"github.com/golang-jwt/jwt/v4"
)

type Users struct {
	Userid   int    `form:"id"`
	Name     string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Date     int    `form:"date"`
}

type UserRegister struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserClaims struct {
	Userid   int    `form:"id"`
	Username string `json:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Date     int    `form:"date"`
	jwt.RegisteredClaims
}
