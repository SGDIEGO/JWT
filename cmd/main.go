package main

import (
	"database/sql"
	"log"

	"github.com/SGDIEGO/JWT/internal/core/handlers"
	"github.com/SGDIEGO/JWT/internal/core/routes"
	"github.com/SGDIEGO/JWT/internal/repositories"
	"github.com/SGDIEGO/JWT/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:diegoasg04@(127.0.0.1:3306)/users")
	if err != nil {
		log.Panic(err)
	}

	server := gin.Default()

	// Dependency Inyection
	UserRepo := repositories.NewUserRepo(db)
	UserService := services.NewUserService(UserRepo)
	UserHandler := handlers.NewHomeHandler(UserService)
	HomeRouter := routes.NewHomeRouter(server, *UserHandler)

	server.Use(cors.Default())
	server.Static("/public", "./public")
	server.LoadHTMLGlob("./www/html/**/*")

	// Routers
	HomeRouter.ServeRouter("/")

	server.Run(":3000")
}
