package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	productRepoImport "kedai/backend/be-kedai/internal/domain/product/repository"
	productServiceImport "kedai/backend/be-kedai/internal/domain/product/service"
	userHandlerImport "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepoImport "kedai/backend/be-kedai/internal/domain/user/repository"
	userServiceImport "kedai/backend/be-kedai/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
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
