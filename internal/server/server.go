package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepo "kedai/backend/be-kedai/internal/domain/location/repository"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	productHandler "kedai/backend/be-kedai/internal/domain/product/handler"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	productService "kedai/backend/be-kedai/internal/domain/product/service"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepository "kedai/backend/be-kedai/internal/domain/user/repository"
	userService "kedai/backend/be-kedai/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()

	cityRepo := locationRepo.NewCityRepository(&locationRepo.CityRConfig{
		DB: db,
	})
	cityService := locationService.NewCityService(&locationService.CitySConfig{
		CityRepo: cityRepo,
	})

	locHandler := locationHandler.New(&locationHandler.Config{
		CityService: cityService,
	})

	userRepo := userRepository.NewUserRepository(&userRepository.UserRConfig{
		DB: db,
	})

	userService := userService.NewUserService(&userService.UserSConfig{
		Repository: userRepo,
	})

	userHandler := userHandler.New(&userHandler.HandlerConfig{
		UserService: userService,
	})

	productRepo := productRepo.NewProductRepository(&productRepo.ProductRConfig{
		DB: db,
	})
	productService := productService.NewProductService(&productService.ProductSConfig{
		Repository: productRepo,
	})
	productHandler := productHandler.New(&productHandler.HandlerConfig{
		ProductService: productService,
	})

	return NewRouter(&RouterConfig{
		LocationHandler: locHandler,
		UserHandler:     userHandler,
		ProductHandler:  productHandler,
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
