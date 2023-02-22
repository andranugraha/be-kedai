package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepo "kedai/backend/be-kedai/internal/domain/location/repository"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	userCache "kedai/backend/be-kedai/internal/domain/user/cache"

	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	userService "kedai/backend/be-kedai/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()
	redis := connection.GetCache()

	cityRepo := locationRepo.NewCityRepository(&locationRepo.CityRConfig{
		DB: db,
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

	userRepo := userRepo.NewUserRepository(&userRepo.UserRConfig{
		DB: db,
	})

	userCache := userCache.NewUserCache(&userCache.UserCConfig{
		RDC: redis,
	})

	userService := userService.NewUserService(&userService.UserSConfig{
		Repository: userRepo,
		Redis:      userCache,
	})

	userHandler := userHandler.New(&userHandler.HandlerConfig{
		UserService:   userService,
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
