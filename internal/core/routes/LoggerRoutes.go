package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SGDIEGO/JWT/internal/domains"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_word")
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type LoggerRoutes struct{}

func NewLoggerRoutes() *LoggerRoutes {
	return &LoggerRoutes{}
}

func (lg *LoggerRoutes) Signin(c *gin.Context) {

	var creds domains.UserRegister
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		c.Status(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &domains.UserClaims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func (lg *LoggerRoutes) Welcome(ctx *gin.Context) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := ctx.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			ctx.Status(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &domains.UserClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Status(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	ctx.String(200, "Welcome %s!", claims.Username)
}

func (lg *LoggerRoutes) Refresh(ctx *gin.Context) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := ctx.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Status(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &domains.UserClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Status(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func (lg *LoggerRoutes) Logout(ctx *gin.Context) {
	// immediately clear the token cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
