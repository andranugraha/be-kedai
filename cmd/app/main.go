package main

import (
	"kedai/backend/be-kedai/connection"
	"kedai/backend/be-kedai/internal/server"
	"log"
)

func main() {
	err := connection.ConnectDB()
	if err != nil {
		log.Fatal("couldn't connect to DB:", err.Error())
	}

	err = connection.ConnectCache()
	if err != nil {
		log.Fatal("couldn't connect to Cache:", err.Error())
	}

	connection.ConnectMailer()

	server.Init()
}
