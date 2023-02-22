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
	LocationHandler *locationHandler.Handler
	UserHandler     *userHandler.Handler
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
			authenticated := user.Group("", middleware.JWTAuthorization)
			{
				wallet := authenticated.Group("/wallets")
				{
					wallet.POST("", cfg.UserHandler.RegisterWallet)
				}
			}
			user.POST("/register", cfg.UserHandler.UserRegistration)
			user.POST("/login", cfg.UserHandler.UserLogin)
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)
		}

	}

	return r
}
