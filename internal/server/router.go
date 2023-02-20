package server

import (
	"kedai/backend/be-kedai/config"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct{
	UserHandler userHandler.Handler
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	userHandler := cfg.UserHandler

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = config.Origin
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsCfg.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(corsCfg))

	r.Static("/docs", "swagger-ui")

	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.UserRegistration)
		}
	}

	return r
}
