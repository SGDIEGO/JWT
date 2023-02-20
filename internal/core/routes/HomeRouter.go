package routes

import (
	"github.com/SGDIEGO/JWT/internal/core/handlers"
	"github.com/gin-gonic/gin"
)

type HomeRouter struct {
	Server   *gin.Engine
	Handlers handlers.HomeHandler
}

func NewHomeRouter(server *gin.Engine, handlers handlers.HomeHandler) *HomeRouter {
	return &HomeRouter{
		Server:   server,
		Handlers: handlers,
	}
}

func (hR *HomeRouter) ServeRouter(url string) {
	router := hR.Server.Group(url)
	{
		router.GET("/", hR.Handlers.Index)
		router.GET("/users", hR.Handlers.ShowUsers)

		router.GET("/register", hR.Handlers.RegisterGet)
		router.POST("/register", hR.Handlers.RegisterUser)
		router.GET("/login", hR.Handlers.LoginGet)
		router.POST("/login", hR.Handlers.LoginUser)

		router.GET("/logout", hR.Handlers.LogOut)
	}
}
