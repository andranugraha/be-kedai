package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
<<<<<<< HEAD
	productRepoImport "kedai/backend/be-kedai/internal/domain/product/repository"
	productServiceImport "kedai/backend/be-kedai/internal/domain/product/service"
	userHandlerImport "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepoImport "kedai/backend/be-kedai/internal/domain/user/repository"
	userServiceImport "kedai/backend/be-kedai/internal/domain/user/service"
=======
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepo "kedai/backend/be-kedai/internal/domain/location/repository"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepository "kedai/backend/be-kedai/internal/domain/user/repository"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
	userCache "kedai/backend/be-kedai/internal/domain/user/cache"
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
<<<<<<< HEAD
	userWishlistRepo := userRepoImport.NewUserWishlistRepository(&userRepoImport.UserWishlistRConfig{
		DB: connection.GetDB(),
	})
	userRepo := userRepoImport.NewUserRepository(&userRepoImport.UserRConfig{
		DB: connection.GetDB(),
	})
	productRepo := productRepoImport.NewProductRepository(&productRepoImport.ProductRConfig{
		DB: connection.GetDB(),
	})
	userService := userServiceImport.NewUserService(&userServiceImport.UserSConfig{
		Repository: userRepo,
	})

	productService := productServiceImport.NewProductService(&productServiceImport.ProductSConfig{
		ProductRepository: productRepo,
	})

	userWishlistService := userServiceImport.NewUserWishlistService(&userServiceImport.UserWishlistSConfig{
		UserWishlistRepository: userWishlistRepo,
		UserService:            userService,
		ProductService:         productService,
	})

	userHandler := userHandlerImport.NewHandler(&userHandlerImport.HandlerConfig{
		UserWishlistService: userWishlistService,
	})

	return NewRouter(&RouterConfig{
=======
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

	userRepo := userRepository.NewUserRepository(&userRepository.UserRConfig{
		DB: db,
	})

	userCache := userCache.NewUserCache(&userCache.UserCConfig{
		RDC: redis,
	})

	userService := userService.NewUserService(&userService.UserSConfig{
		Repository: userRepo,
		Redis: userCache,
	})

	userHandler := userHandler.New(&userHandler.HandlerConfig{
		UserService: userService,
	})

	return NewRouter(&RouterConfig{
		LocationHandler: locHandler,
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
		UserHandler: userHandler,
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
