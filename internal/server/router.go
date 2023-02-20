package server

import (
	"kedai/backend/be-kedai/config"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
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

	r.Static("/docs", "swagger")

	v1 := r.Group("/v1")
	{
		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)
		}
	}

	return r
}
