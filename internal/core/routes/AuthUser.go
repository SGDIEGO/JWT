package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SGDIEGO/JWT/internal/domains"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRouter struct {
	JwtKey         []byte
	UsersPermitted map[string]string
}

func NewAuthRouter() *AuthRouter {
	return &AuthRouter{
		JwtKey: []byte("my_jwt_key"),
		UsersPermitted: map[string]string{
			"user1": "password1",
			"user2": "password2",
		},
	}
}

func (au *AuthRouter) SingIn(ctx *gin.Context) {}

func (au *AuthRouter) LogIn(ctx *gin.Context) {

	fmt.Println(ctx.Param("id"))

	var userInfo domains.UserRegister

	// ERROR
	if err := ctx.ShouldBind(&userInfo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": "Pendejo",
		})
		return
	}

	isRegister, ok := au.UsersPermitted[userInfo.Username]

	if !ok || isRegister != userInfo.Password {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
		return
	}

	// GET USER INFO
	expirationTime := time.Now().Add(5 * time.Minute) // 5 minutes expiration
	userClaim := &domains.UserClaims{
		Username: userInfo.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	} //Info for token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim) // Create token
	tokenString, err := token.SignedString(au.JwtKey)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	cookie, _ := ctx.Request.Cookie("token")

	ctx.JSON(200, gin.H{
		"Data":   userInfo,
		"Cookie": cookie.Value,
	})
}

func (au *AuthRouter) LogOut(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
