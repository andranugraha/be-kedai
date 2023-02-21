package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepo "kedai/backend/be-kedai/internal/domain/location/repository"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"

	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	userService "kedai/backend/be-kedai/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	cityRepo := locationRepo.NewCityRepository(&locationRepo.CityRConfig{
		DB: connection.GetDB(),
	})
	cityService := locationService.NewCityService(&locationService.CitySConfig{
		CityRepo: cityRepo,
	})

	locHandler := locationHandler.New(&locationHandler.Config{
		CityService: cityService,
	})

	walletRepo := userRepo.NewWalletRepository(&userRepo.WalletRConfig{
		DB: connection.GetDB(),
	})
	walletService := userService.NewWalletService(&userService.WalletSConfig{
		WalletRepo: walletRepo,
	})
	userHandler := userHandler.New(&userHandler.HandlerConfig{
		WalletService: walletService,
	})

	return NewRouter(&RouterConfig{
		LocationHandler: locHandler,
		UserHandler:     userHandler,
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
