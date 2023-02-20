package server

import (
	"kedai/backend/be-kedai/connection"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepository "kedai/backend/be-kedai/internal/domain/user/repository"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
	"log"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()

	userRepo := userRepository.NewUserRepository(&userRepository.UserRConfig{
		DB: db,
	})

	userService := userService.NewUserService(&userService.UserSConfig{
		Repository: userRepo,
	})

	userHandler := userHandler.New(&userHandler.HandlerConfig{
		UserService: userService,
	})

	return NewRouter(&RouterConfig{
		UserHandler: *userHandler,
	})
}

func Init() {
	r := createRouter()
	err := r.Run()
	if err != nil {
		log.Println("error while running server", err)
		return
	}
}
