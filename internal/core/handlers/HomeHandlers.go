package handlers

import (
	"net/http"
	"time"

	"github.com/SGDIEGO/JWT/internal/domains"
	"github.com/SGDIEGO/JWT/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var JWT = []byte("my_JWT_key")

type HomeHandler struct {
	UserService ports.UserServiceI
}

func NewHomeHandler(UserService ports.UserServiceI) *HomeHandler {
	return &HomeHandler{
		UserService: UserService,
	}
}

// GET
func (hH *HomeHandler) Index(ctx *gin.Context) {

	userLogin, isLogin := hH.userLogin(ctx)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
		"user":  userLogin,
		"login": isLogin,
	})
}

// GET
func (hH *HomeHandler) ShowUsers(ctx *gin.Context) {
	users, err := hH.UserService.GetUsers()

	if err != nil {
		ctx.HTML(http.StatusBadRequest, "users.html", gin.H{
			"title":     err.Error(),
			"validData": false,
		})
		return
	}

	ctx.HTML(http.StatusOK, "users.html", gin.H{
		"title":     "USERS",
		"data":      users,
		"validData": true,
	})
}

// GET
func (hH *HomeHandler) RegisterGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{
		"title": "REGISTER",
	})
}

// POST
func (hH *HomeHandler) RegisterUser(ctx *gin.Context) {
	var UserCreated *domains.Users

	// User isnt bind data
	if err := ctx.ShouldBind(&UserCreated); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Verify if user exists
	VerifyUser, _ := hH.UserService.GetUserByName(UserCreated.Name)
	if VerifyUser != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": "User exists",
		})
		return
	}

	// Error creating user
	if err := hH.UserService.CreateUser(UserCreated); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// GET
func (hH *HomeHandler) LoginGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "LOGIN",
	})
}

// POST
func (hH *HomeHandler) LoginUser(ctx *gin.Context) {
	var UserBind *domains.Users

	// User isnt bind data
	if err := ctx.ShouldBind(&UserBind); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Error loading user
	userLogin, err := hH.UserService.GetUserByName(UserBind.Name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Error": "User not register",
		})
		return
	}

	if userLogin.Password != UserBind.Password {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": "Wrong password",
		})
		return
	}

	// JWT
	expirationTime := time.Now().Add(5 * time.Minute)

	userClaim := &domains.UserClaims{
		Userid:   userLogin.Userid,
		Username: userLogin.Name,
		Email:    userLogin.Email,
		Password: userLogin.Password,
		Date:     userLogin.Date,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(JWT)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	ctx.Redirect(http.StatusFound, "/")
	// http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusFound)
}

// GET
func (hH *HomeHandler) userLogin(ctx *gin.Context) (*domains.UserClaims, bool) {
	cookie, err := ctx.Request.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			ctx.Status(http.StatusUnauthorized)
			return nil, false
		}
		ctx.Status(http.StatusBadRequest)
		return nil, false
	}

	getToken := cookie.Value
	claims := &domains.UserClaims{}
	tkn, err := jwt.ParseWithClaims(getToken, claims, func(t *jwt.Token) (interface{}, error) {
		return JWT, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.Status(http.StatusUnauthorized)
			return nil, false
		}
		ctx.Status(http.StatusBadRequest)
		return nil, false
	}
	if !tkn.Valid {
		ctx.Status(http.StatusUnauthorized)
		return nil, false
	}

	return claims, true
}

// GET
func (hH *HomeHandler) LogOut(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	ctx.Redirect(http.StatusFound, "/")
}
