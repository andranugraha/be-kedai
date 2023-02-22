package server

import (
	"kedai/backend/be-kedai/config"
<<<<<<< HEAD
=======
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
<<<<<<< HEAD
	UserHandler *userHandler.Handler
=======
	LocationHandler *locationHandler.Handler
	UserHandler     *userHandler.Handler
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
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
<<<<<<< HEAD
		users := v1.Group("/users")
		{
			authenticated := users.Group("", middleware.JWTAuthorization)
			{
				wishlists := authenticated.Group("/wishlists")
				{
					wishlists.POST("", cfg.UserHandler.AddUserWishlist)
				}
			}
=======

		user := v1.Group("/users")
		{
			user.GET("", middleware.JWTAuthorization, cfg.UserHandler.GetSession, cfg.UserHandler.GetUserByID)
			user.POST("/register", cfg.UserHandler.UserRegistration)
			user.POST("/login", cfg.UserHandler.UserLogin)
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
		}

<<<<<<< HEAD
		r.Static("/docs", "swagger")

	}
=======
>>>>>>> e4fd8db74c2d1f5d9ac94cf1de0592b0a77f3219
	return r
}
