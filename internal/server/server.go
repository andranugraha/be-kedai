package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	return NewRouter(&RouterConfig{})
}

func Init() {
	r := createRouter()
	err := r.Run()
	if err != nil {
		log.Println("error while running server", err)
		return
	}
}
