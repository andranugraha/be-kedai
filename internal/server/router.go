package server

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/server/middleware"

	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	productHandler "kedai/backend/be-kedai/internal/domain/product/handler"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserHandler     *userHandler.Handler
	LocationHandler *locationHandler.Handler
	ProductHandler  *productHandler.Handler
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = config.Origin
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsCfg.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(corsCfg))

	v1 := r.Group("/v1")
	{
		v1.Static("/docs", "swagger")

		user := v1.Group("/users")
		{
			user.POST("/register", cfg.UserHandler.UserRegistration)
			user.POST("/login", cfg.UserHandler.UserLogin)
			user.POST("/google-login", cfg.UserHandler.UserLoginWithGoogle)
			userAuthenticated := user.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				userAuthenticated.GET("", cfg.UserHandler.GetUserByID)
				wallet := userAuthenticated.Group("/wallets")
				{
					wallet.GET("", cfg.UserHandler.GetWalletByUserID)
					wallet.POST("", cfg.UserHandler.RegisterWallet)
				}
				wishlists := userAuthenticated.Group("/wishlists")
				{
					wishlists.GET("", cfg.UserHandler.GetUserWishlists)
					wishlists.GET("/:productId", cfg.UserHandler.GetUserWishlist)
					wishlists.POST("", cfg.UserHandler.AddUserWishlist)
					wishlists.DELETE("/:productId", cfg.UserHandler.RemoveUserWishlist)
				}
				carts := userAuthenticated.Group("/carts")
				{
					carts.POST("", cfg.UserHandler.CreateCartItem)
					carts.GET("", cfg.UserHandler.GetAllCartItem)
				}
			}
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)

		}

		product := v1.Group("/products")
		{
			product.GET("/recommendation", cfg.ProductHandler.GetRecommendation)
			category := product.Group("/categories")
			{
				category.GET("", cfg.ProductHandler.GetCategories)
			}
		}
	}

	return r
}
