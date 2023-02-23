package server

import (
	"kedai/backend/be-kedai/config"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserHandler     *userHandler.Handler
	LocationHandler *locationHandler.Handler
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
		user := v1.Group("/users")
		{
			user.POST("/register", cfg.UserHandler.UserRegistration)
			user.POST("/login", cfg.UserHandler.UserLogin)
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
					wishlists.GET("/:productId", cfg.UserHandler.GetUserWishlist)
					wishlists.POST("", cfg.UserHandler.AddUserWishlist)
					wishlists.DELETE("/:productId", cfg.UserHandler.RemoveUserWishlist)
				}
			}
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)

		}
		r.Static("/docs", "swagger")

		return r
	}
}
