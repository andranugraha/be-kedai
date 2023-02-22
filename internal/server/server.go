package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandlerPackage "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepoPackage "kedai/backend/be-kedai/internal/domain/location/repository"
	locationServicePackage "kedai/backend/be-kedai/internal/domain/location/service"
	productRepoPackage "kedai/backend/be-kedai/internal/domain/product/repository"
	productServicePackage "kedai/backend/be-kedai/internal/domain/product/service"
	userCache "kedai/backend/be-kedai/internal/domain/user/cache"

	userHandlerPackage "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepoPackage "kedai/backend/be-kedai/internal/domain/user/repository"
	userServicePackage "kedai/backend/be-kedai/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()
	redis := connection.GetCache()

	productRepo := productRepoPackage.NewProductRepository(&productRepoPackage.ProductRConfig{
		DB: db,
	})
	productService := productServicePackage.NewProductService(&productServicePackage.ProductSConfig{
		ProductRepository: productRepo,
	})

	cityRepo := locationRepoPackage.NewCityRepository(&locationRepoPackage.CityRConfig{
		DB: db,
	})
	cityService := locationServicePackage.NewCityService(&locationServicePackage.CitySConfig{
		CityRepo: cityRepo,
	})

	locHandler := locationHandlerPackage.New(&locationHandlerPackage.Config{
		CityService: cityService,
	})

	walletRepo := userRepoPackage.NewWalletRepository(&userRepoPackage.WalletRConfig{
		DB: connection.GetDB(),
	})
	walletService := userServicePackage.NewWalletService(&userServicePackage.WalletSConfig{
		WalletRepo: walletRepo,
	})

	userRepo := userRepoPackage.NewUserRepository(&userRepoPackage.UserRConfig{
		DB: db,
	})

	userCache := userCache.NewUserCache(&userCache.UserCConfig{
		RDC: redis,
	})

	userService := userServicePackage.NewUserService(&userServicePackage.UserSConfig{
		Repository: userRepo,
		Redis:      userCache,
	})

	userWishlistRepo := userRepoPackage.NewUserWishlistRepository(&userRepoPackage.UserWishlistRConfig{
		DB: db,
	})

	userWishlistService := userServicePackage.NewUserWishlistService(&userServicePackage.UserWishlistSConfig{
		UserWishlistRepository: userWishlistRepo,
		UserService:            userService,
		ProductService:         productService,
	})

	userHandler := userHandlerPackage.New(&userHandlerPackage.HandlerConfig{
		UserService:         userService,
		WalletService:       walletService,
		UserWishlistService: userWishlistService,
	})

	return NewRouter(&RouterConfig{
		UserHandler:     userHandler,
		LocationHandler: locHandler,
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
