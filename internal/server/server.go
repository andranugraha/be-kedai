package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepo "kedai/backend/be-kedai/internal/domain/location/repository"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	cityRepo := locationRepo.NewCityRepository(&locationRepo.CityRConfig{
		DB: connection.GetDB(),
	})
	cityService := locationService.NewCityService(&locationService.CitySConfig{
		CityRepo: cityRepo,
	})

	locHandler := locationHandler.New(&locationHandler.Config{
		CityService: cityService,
	})

	return NewRouter(&RouterConfig{
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
